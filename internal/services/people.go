package services

import (
	"strings"

	"github.com/KoLLlaka/sobes/internal/db"
	externalapi "github.com/KoLLlaka/sobes/internal/externalAPI"
	"github.com/KoLLlaka/sobes/internal/logger"
	"github.com/KoLLlaka/sobes/internal/model"
)

type PeopleService interface {
	AddPeople(people model.MessageToDB) (string, error)
	DeletePeople(uid string) error
	UpdatePeople(people model.MessageToDB) error
	GetPeoples() ([]model.MessageToDB, error)
	GetPeople(uid string) (model.MessageToDB, error)
	EnrichPeople(people *model.MessageToDB)
}

type peopleService struct {
	store  db.Store
	enrich externalapi.Enrich
	logger *logger.MyLogger
}

// create a new people service
func NewPeopleService(
	store *db.Store,
	enrich *externalapi.Enrich,
	logger *logger.MyLogger,
) PeopleService {
	return &peopleService{
		store:  *store,
		enrich: *enrich,
		logger: logger,
	}
}

// method of service to add people to db interface
func (p *peopleService) AddPeople(people model.MessageToDB) (string, error) {
	return p.store.NewPeopleStore().AddPeople(people)
}

// method of service to delete people from db interface
func (p *peopleService) DeletePeople(uid string) error {
	return p.store.NewPeopleStore().DeletePeople(uid)
}

// method of service to update people on db interface
func (p *peopleService) UpdatePeople(people model.MessageToDB) error {
	return p.store.NewPeopleStore().UpdatePeople(people)
}

// method of service to recieve peoples from db interface
func (p *peopleService) GetPeoples() ([]model.MessageToDB, error) {
	return p.store.NewPeopleStore().GetPeoples()
}

// method of service to recieve people from db interface by uid
func (p *peopleService) GetPeople(uid string) (model.MessageToDB, error) {
	return p.store.NewPeopleStore().GetPeople(uid)
}

// method of service to add age, gender, nation
// and error of this functions to people
func (p *peopleService) EnrichPeople(people *model.MessageToDB) {
	var (
		param     = map[string]string{"name": people.Name}
		errorToDB = []string{}
		err       error
	)
	people.Age, err = p.enrich.NewPeopleEnrich().AddAge(param)
	if err != nil {
		errorToDB = append(errorToDB, err.Error())
	}

	people.Gender, err = p.enrich.NewPeopleEnrich().AddGender(param)
	if err != nil {
		errorToDB = append(errorToDB, err.Error())
	}

	people.Nation, err = p.enrich.NewPeopleEnrich().AddNation(param)
	if err != nil {
		errorToDB = append(errorToDB, err.Error())
	}

	if len(errorToDB) > 0 {
		people.Error = strings.Join(errorToDB, ",")
	}
}
