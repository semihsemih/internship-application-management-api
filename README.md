# Internship Application Management Api
internship application management api built to learn Go programming language. Gorilla mux and go mongo driver is used to build the api
## Usage
#### Create Candidate

You can create a candidate by posting a candidate model like the following:
```bash
curl -X POST \
  http://localhost:8080/candidates \
  -H 'content-type: application/json' \
  -d '{
    "first_name" : "Semih",
    "last_name" : "Arslan",
    "email" : "test@test.com",
    "department" : "Development",
    "university" : "Ankara",
    "experience" : false
   }'
```
`ApplicationDate`, `Status`, `MeetingCount`, `NextMeeting` will be set automatically after you create the candidate. `Assignee` field will be set after you have arranged a meeting with this candidate. `Email`, `Department` and `University` fields required.
Also, email format should be example@email.xyz. Otherwise, the api returns bad request and candidate will be not inserted to DB.

#### Read Candidate

You can read a candidate by using its id like the following:
```bash
curl -X GET http://localhost:8080/candidates/5ea980281dafc611002fbc41
```

#### Delete Candidate

You can delete a candidate by using its id like the following:
```bash
curl -X DELETE http://localhost:8080/candidates/5ea980281dafc611002fbc41
```

#### Arrange Meeting

You can arrange a meeting with a candidate by posting candidate id and next meeting time in request body:
```bash
curl -X POST \
  http://localhost:8080/meetings/arrange \
  -H 'accept: application/json' \
  -d '{
	"candidate_id": "5ea980281dafc611002fbc41",
	"next_meeting_time": "2020-05-03T13:40:00.000+00:00"
  }'
```
Both of the `candidate_id` and `next_meeting_time` are required in order to arrange the meeting. After arranging a meeting, a randomly chosen assignee will be assigned to the given candidate according to the department they have applied.

#### Complete Meeting

You can complete meeting by posting the candidate id like the following:
```bash
curl -X POST http://localhost:8080/meetings/complete/5ea980281dafc611002fbc41
```

#### Deny Candidate

You can deny a candidate by using its id like the following:
```bash
curl -X PATCH http://localhost:8080/candidates/deny/5ea980281dafc611002fbc41
```

#### Accept Candidate

You can accept a candidate by using its id like the following:
```bash
curl -X PATCH http://localhost:8080/candidates/accept/5ea980281dafc611002fbc41
```

#### Find Assignee ID by Name

You can find the assignee id by using its name like the following:
```bash
curl -X GET http://localhost:8080/assignees/name/Sercan
```
Please note that this endpoint is case-sensitive. It **will not produce** the same result with sercan as it produced with Sercan

### Bonus Functions

#### Find Assignee's Candidates

You can find the candidates of the assignee by using its id like the following:
```bash
curl -X GET http://localhost:8080/candidates/assigneeId/5ea980281dafc611002fbc41
```

#### Find All Candidates

You can find all candidates that are available in the system like the following:
```bash
curl -X GET http://localhost:8080/candidates
```

#### Create Assignee

You can create an assignee by posting an assignee model like the following:
```bash
curl -X POST \
  http://localhost:8080/assignees \
  -H 'content-type: application/json' \
  -d '{
    "name" : "semih",
    "department" : "Development"
   }'
```
Both of the `name` and `department` are required in order to create an assignee.

#### Find All Assignees

You can find all assignees that are available in the system like the following:
```bash
curl -X GET http://localhost:8080/assignees
```

#### Find All Assignees by Department

You can find all assignees in a department by using department name like the following:
```bash
curl -X GET http://localhost:8080/assignees/department/Design
```
Please note that this endpoint is case-sensitive. It **will not produce** the same result with design as it produced with Design
