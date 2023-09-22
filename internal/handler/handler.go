package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/KoLLlaka/sobes/internal/logger"
	"github.com/KoLLlaka/sobes/internal/model"
	"github.com/KoLLlaka/sobes/internal/services"
	"github.com/gorilla/mux"
)

const (
	UID = "uid"
)

type server struct {
	peopleService services.PeopleService
	logger        logger.MyLogger
}

// create a new stuct of html server
func NewServer(
	logger logger.MyLogger,
	peopleService services.PeopleService,
) *server {
	return &server{
		peopleService: peopleService,
		logger:        logger,
	}
}

// start html server on host:port
func (s *server) StartServer(host, port string) {
	router := mux.NewRouter()
	router.HandleFunc("/api/peoples", s.getPeoplesHandler).Methods("GET")
	router.HandleFunc("/api/peoples/{"+UID+"}", s.getPeopleHandler).Methods("GET")
	router.HandleFunc("/api/peoples", s.addPeopleHandler).Methods("POST")
	router.HandleFunc("/api/peoples/{"+UID+"}", s.updatePeopleHandler).Methods("PUT")
	router.HandleFunc("/api/peoples/{"+UID+"}", s.deletePeopleHandler).Methods("DELETE")

	s.logger.LogTraceLvl(
		fmt.Sprintf("server start on %s:%s", host, port),
		"handler",
		"StartServer",
		nil,
		nil,
	)
	go http.ListenAndServe(fmt.Sprintf("%s:%s", host, port), router)
}

// handler to recieve list of peoples from db
func (s *server) getPeoplesHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	peoples, err := s.peopleService.GetPeoples()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "server error"})

		return
	}

	s.logger.LogTraceLvl("get list of peoples", "handler", "getPeoplesHandler", nil, peoples)
	json.NewEncoder(w).Encode(peoples)
}

// handler to recieve people from db by uid
func (s *server) getPeopleHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	uid := mux.Vars(r)[UID]

	people, err := s.peopleService.GetPeople(uid)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "wrong uid"})

		return
	}

	s.logger.LogTraceLvl("get people by uid", "handler", "getPeopleHandler", nil, people)
	json.NewEncoder(w).Encode(people)
}

// handler to add new people to db
// recieve uid from db
func (s *server) addPeopleHandler(w http.ResponseWriter, r *http.Request) {
	people := model.MessageToDB{}
	w.Header().Set("Content-Type", "application/json")
	err := json.NewDecoder(r.Body).Decode(&people)
	if err != nil {
		s.logger.LogWarningLvl(
			"Error decoding",
			"handler",
			"addPeopleHandler",
			err,
			r.Body,
		)
	}

	s.peopleService.EnrichPeople(&people)

	uid, err := s.peopleService.AddPeople(people)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "internel server error"})

		return
	}

	s.logger.LogTraceLvl("people add to db by uid", "handler", "addPeopleHandler", nil, uid)
	json.NewEncoder(w).Encode(map[string]string{"uid": uid})
}

// handler to update people on db by uid
func (s *server) updatePeopleHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	uid := mux.Vars(r)[UID]
	people := model.MessageToDB{
		UID: uid,
	}

	err := json.NewDecoder(r.Body).Decode(&people)
	if err != nil {
		s.logger.LogWarningLvl(
			"Error decoding",
			"handler",
			"updatePeopleHandler",
			err,
			r.Body,
		)
	}

	err = s.peopleService.UpdatePeople(people)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "wrong uid"})

		return
	}

	s.logger.LogTraceLvl("people update on db by uid", "handler", "updatePeopleHandler", nil, people)
	json.NewEncoder(w).Encode(people)
}

// handler to delete people from db by uid
func (s *server) deletePeopleHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	uid := mux.Vars(r)[UID]

	err := s.peopleService.DeletePeople(uid)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "server error"})

		return
	}

	s.logger.LogTraceLvl("delete people from db by uid", "handler", "deletePeopleHandler", nil, uid)
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}
