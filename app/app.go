package app

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/semihsemih/internship-application-management-api/db"
	"github.com/semihsemih/internship-application-management-api/model"
	"github.com/semihsemih/internship-application-management-api/util"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func CreateCandidate(candidate model.Candidate) (model.Candidate, error) {
	// Check if there is anyone working in the department.
	numberOfAssigneeInDepartment, err := AssigneeCountByDepartment(candidate.Department)
	if err != nil {
		return model.Candidate{}, err
	}
	if numberOfAssigneeInDepartment == 0 {
		return model.Candidate{}, errors.New("New candidates should have an assignee who is working in the department that the candidate is applying to work.")
	}

	// Checks the validity of the candidate's e-mail address.
	if !util.IsEmailValid(candidate.Email) {
		return model.Candidate{}, errors.New("Not a valid email")
	}
	_, err = FindCandidateByEmail(candidate.Email)
	if err == nil {
		return model.Candidate{}, errors.New("This e-mail address is registered in the system.")
	}

	candidate.ID = primitive.NewObjectID().Hex()
	candidate.MeetingCount = 0
	candidate.NextMeeting = nil
	candidate.Status = "Pending"
	candidate.ApplicationDate = time.Now() // Forgotten part =)

	ctx, _ := context.WithTimeout(context.TODO(), 10*time.Millisecond)
	client := db.GetClient()
	result, err := client.Database("Otsimo").Collection("Candidates").InsertOne(ctx, candidate)
	if err != nil {
		log.Println(err)
	}
	candidate.ID = fmt.Sprintf("%v", result.InsertedID)

	return candidate, err
}

func ReadCandidate(_id string) (model.Candidate, error) {
	var c model.Candidate
	client := db.GetClient()
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Millisecond)
	err := client.Database("Otsimo").Collection("Candidates").FindOne(ctx, bson.D{{"_id", _id}}).Decode(&c)
	if err != nil {
		log.Println(err)
	}

	return c, err
}

func DeleteCandidate(_id string) error {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Millisecond)
	// Check if there is any candidate with this id.
	_, err := ReadCandidate(_id)
	if err != nil {
		return errors.New("There is no such candidate.")
	}
	client := db.GetClient()
	_, err = client.Database("Otsimo").Collection("Candidates").DeleteOne(ctx, bson.M{"_id": _id})
	return err
}

func ArrangeMeeting(_id string, nextMeetingTime *time.Time) error {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Millisecond)
	client := db.GetClient()
	// Check if there is any candidate with this id.
	c, err := ReadCandidate(_id)
	if err != nil {
		return errors.New("There is no such candidate.")
	}

	// Update candidate's meeting time
	c.NextMeeting = nextMeetingTime

	// Check candidate's meeting count and change candidate's assignee
	if c.MeetingCount >= 0 && c.MeetingCount < 3 {
		// Every meeting select random assignee
		a, err := GetRandomAssignee(c.Department)
		if err != nil {
			return err
		}
		c.Assignee = a.ID
	} else if c.MeetingCount == 3 {
		// Candidate meets with CEO
		var ceo model.Assignee
		err = client.Database("Otsimo").Collection("Assignees").FindOne(context.TODO(), bson.M{"department": "CEO"}).Decode(&ceo)
		if err != nil {
			return err
		} else {
			c.Assignee = ceo.ID
		}
	}

	_, err = client.Database("Otsimo").Collection("Candidates").UpdateOne(ctx, bson.M{"_id": c.ID}, bson.M{"$set": &c})
	return err
}

func CompleteMeeting(_id string) error {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Millisecond)
	client := db.GetClient()
	// Check if there is any candidate with this id.
	c, err := ReadCandidate(_id)
	if err != nil {
		return errors.New("There is no such candidate.")
	}

	// Meeting completing increases meeting count and change status
	if c.MeetingCount < 4 {
		c.MeetingCount++
		c.Status = "In Progress"
	}
	if c.NextMeeting != nil {
		c.NextMeeting = nil
	}

	_, err = client.Database("Otsimo").Collection("Candidates").UpdateOne(ctx, bson.M{"_id": c.ID}, bson.M{"$set": &c})
	return err
}

func DenyCandidate(_id string) error {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Millisecond)
	client := db.GetClient()
	// Check if there is any candidate with this id.
	c, err := ReadCandidate(_id)
	if err != nil {
		return errors.New("There is no such candidate.")
	}

	// Change status
	c.Status = "Denied"
	_, err = client.Database("Otsimo").Collection("Candidates").UpdateOne(ctx, bson.M{"_id": c.ID}, bson.M{"$set": &c})
	return err
}

func AcceptCandidate(_id string) error {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Millisecond)
	client := db.GetClient()

	// Check if there is any candidate with this id.
	c, err := ReadCandidate(_id)
	if err != nil {
		return errors.New("There is no such candidate.")
	}
	// Candidates cannot be accepted before the completion of 4 meetings.
	if c.MeetingCount < 4 {
		return errors.New("There is incomplete meetings.")
	}

	// Change status
	c.Status = "Accepted"
	_, err = client.Database("Otsimo").Collection("Candidates").UpdateOne(ctx, bson.M{"_id": c.ID}, bson.M{"$set": &c})
	return err
}

func FindAssigneeIDByName(name string) string {
	var a model.Assignee
	client := db.GetClient()
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Millisecond)
	err := client.Database("Otsimo").Collection("Assignees").FindOne(ctx, bson.D{{"name", name}}).Decode(&a)
	if err != nil {
		log.Println(err)
	}

	return a.ID
}

// This is bonus feature =)
// Returns candidates by given assignee's id
func FindAssigneesCandidates(_id string) ([]model.Candidate, error) {
	var c []model.Candidate
	var a model.Assignee
	client := db.GetClient()
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Millisecond)
	// Check if there is any assignee with this id.
	err := client.Database("Otsimo").Collection("Assignees").FindOne(ctx, bson.D{{"_id", _id}}).Decode(&a)
	if err != nil {
		return nil, errors.New("There is no such assignee")
	}

	cursor, err := client.Database("Otsimo").Collection("Candidates").Find(ctx, bson.D{{"assignee", _id}})
	err = cursor.All(ctx, &c)
	if err != nil {
		log.Println(err)
	}

	return c, err
}

// Returns number of assignees working in the department
func AssigneeCountByDepartment(department string) (int64, error) {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Millisecond)
	client := db.GetClient()
	itemCount, err := client.Database("Otsimo").Collection("Assignees").CountDocuments(ctx, bson.M{"department": department})
	if err != nil {
		log.Println(err)
	}

	return itemCount, err
}

// Randomly return assignee by given department name
func GetRandomAssignee(department string) (model.Assignee, error) {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client := db.GetClient()
	var assignees []model.Assignee

	cursor, err := client.Database("Otsimo").Collection("Assignees").Aggregate(ctx, mongo.Pipeline{
		bson.D{{"$match", bson.D{{"department", department}}}}, // Given department name
		bson.D{{"$sample", bson.D{{"size", 1}}}},               // How many to choose
	})
	err = cursor.All(ctx, &assignees)
	if err != nil {
		log.Println(err)
	}

	var assignee model.Assignee
	if len(assignees) > 0 {
		assignee = assignees[0]
	} else {
		return model.Assignee{}, errors.New("There is nobody working in this department.")
	}

	return assignee, err
}

// Return candidate by given department name
func FindCandidateByEmail(email string) (model.Candidate, error) {
	client := db.GetClient()
	var c model.Candidate
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Millisecond)
	err := client.Database("Otsimo").Collection("Candidates").FindOne(ctx, bson.D{{"email", email}}).Decode(&c)
	if err != nil {
		log.Println(err)
	}

	return c, err
}
