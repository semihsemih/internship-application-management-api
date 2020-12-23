package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/semihsemih/internship-application-management-api/app"
	"github.com/semihsemih/internship-application-management-api/model"
)

func CreateCandidateHandler(w http.ResponseWriter, r *http.Request) {
	var c model.Candidate
	// Read request body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
	}
	// Call json.Unmarshal, passing it a []byte of JSON data and a pointer to c
	json.Unmarshal(body, &c)

	createdCandidate, err := app.CreateCandidate(c)
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte(err.Error()))
	} else {
		w.Header().Set("Content-Type", "application/json")
		// Candidate encoded as json
		response, _ := json.Marshal(createdCandidate)
		w.WriteHeader(201)
		w.Write(response)
	}
}

func ReadCandidateHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	// Read candidate's id from route path
	c, err := app.ReadCandidate(vars["id"])
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte(err.Error()))
	} else {
		w.Header().Set("Content-Type", "application/json")
		// Candidate encoded as json
		response, _ := json.Marshal(c)
		w.WriteHeader(200)
		w.Write(response)
	}
}

func DeleteCandidateHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	// Read candidate's id from route path and pass to function
	err := app.DeleteCandidate(vars["id"])
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte(err.Error()))
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(202)
		w.Write([]byte("Candidate Deleted"))
	}
}

func ArrangeMeetingHandler(w http.ResponseWriter, r *http.Request) {
	var meeting struct {
		CandidateId     string     `json:"candidate_id"`
		NextMeetingTime *time.Time `json:"next_meeting_time"`
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
	}
	json.Unmarshal(body, &meeting)
	err = app.ArrangeMeeting(meeting.CandidateId, meeting.NextMeetingTime)
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte(err.Error()))
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		w.Write([]byte("The next meeting has been arranged."))
	}
}

func CompleteMeetingHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	_id := vars["id"]
	err := app.CompleteMeeting(_id)
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte(err.Error()))
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		w.Write([]byte("Meeting complete"))
	}
}

func DenyCandidateHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	_id := vars["id"]
	err := app.DenyCandidate(_id)
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte(err.Error()))
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		w.Write([]byte("Candidate denied."))
	}
}

func AcceptCandidateHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	_id := vars["id"]
	err := app.AcceptCandidate(_id)
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte(err.Error()))
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		w.Write([]byte("Candidate accepted."))
	}
}

func FindAssigneeIDByNameHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := app.FindAssigneeIDByName(vars["name"])
	if id == "" {
		w.WriteHeader(400)
		w.Write([]byte("There is no such user."))
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		w.Write([]byte(id))
	}
}

func FindAssigneesCandidatesHandler(w http.ResponseWriter, r *http.Request) {
	var c []model.Candidate
	vars := mux.Vars(r)
	id := vars["id"]
	c, err := app.FindAssigneesCandidates(id)
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte(err.Error()))
	} else if len(c) == 0 {
		w.WriteHeader(400)
		w.Write([]byte("There is no such candidates"))
	} else {
		response, _ := json.Marshal(c)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		w.Write(response)
	}
}
