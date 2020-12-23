package router

import (
	"github.com/gorilla/mux"
	"github.com/semihsemih/internship-application-management-api/api"
)

func Init() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/candidates", api.CreateCandidateHandler).Methods("POST")
	r.HandleFunc("/candidates/{id:[a-zA-Z0-9]*}", api.ReadCandidateHandler).Methods("GET")
	r.HandleFunc("/candidates/{id:[a-zA-Z0-9]*}", api.DeleteCandidateHandler).Methods("DELETE")
	r.HandleFunc("/meetings/arrange", api.ArrangeMeetingHandler).Methods("PUT")
	r.HandleFunc("/meetings/complete/{id:[a-zA-Z0-9]*}", api.CompleteMeetingHandler).Methods("PUT")
	r.HandleFunc("/candidates/deny/{id:[a-zA-Z0-9]*}", api.DenyCandidateHandler).Methods("PUT")
	r.HandleFunc("/candidates/accept/{id:[a-zA-Z0-9]*}", api.AcceptCandidateHandler).Methods("PUT")
	r.HandleFunc("/assignees/name/{name:[a-zA-Z]*}", api.FindAssigneeIDByNameHandler).Methods("GET")
	r.HandleFunc("/assignees/{id:[a-zA-Z0-9]*}/candidates", api.FindAssigneesCandidatesHandler).Methods("GET")

	return r
}
