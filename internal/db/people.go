package db

import (
	"database/sql"
	"errors"

	"github.com/KoLLlaka/sobes/internal/logger"
	"github.com/KoLLlaka/sobes/internal/model"
	"github.com/google/uuid"
)

type PeopleStore interface {
	AddPeople(model.MessageToDB) (string, error)
	DeletePeople(string) error
	UpdatePeople(model.MessageToDB) error
	GetPeoples() ([]model.MessageToDB, error)
	GetPeople(uid string) (model.MessageToDB, error)
}

type peopleStore struct {
	db     *sql.DB
	logger logger.MyLogger
}

// add new people to db
func (p peopleStore) AddPeople(people model.MessageToDB) (string, error) {
	uid := uuid.New().String()
	people.UID = uid

	query := `
	INSERT INTO
		people
		(uid, name, surname, patronymic, age, gender, nation, error)
	VALUES
		($1, $2, $3, $4, $5, $6, $7, $8);`
	_, err := p.db.Exec(
		query,
		people.UID,
		people.Name,
		people.Surname,
		people.Patronymic,
		people.Age,
		people.Gender,
		people.Nation,
		people.Error,
	)
	if err != nil {
		p.logger.LogErrorLvl("exec error", "db", "AddPeople", err, nil)

		return "", err
	}

	return uid, nil
}

// delete people from db by uid
func (p peopleStore) DeletePeople(uid string) error {
	query := `
	DELETE FROM
		people
	WHERE
		uid = $1;`
	result, err := p.db.Exec(
		query,
		uid,
	)
	if err != nil {
		p.logger.LogErrorLvl("exec error", "db", "DeletePeople", err, nil)

		return err
	}

	count, err := result.RowsAffected()
	if err != nil {
		p.logger.LogErrorLvl("RowsAffected error", "db", "DeletePeople", err, nil)

		return err
	}
	if count == 0 {
		p.logger.LogErrorLvl("wrong uid", "db", "DeletePeople", err, uid)

		return errors.New("wrong uid")
	}

	return nil
}

// update people on db
func (p peopleStore) UpdatePeople(people model.MessageToDB) error {
	query := `
	UPDATE
		people
	SET
		name = $2,
		surname = $3,
		patronymic = $4,
		age = $5,
		gender = $6,
		nation = $7,
		error = $8
	WHERE
		uid = $1;`
	result, err := p.db.Exec(
		query,
		people.UID,
		people.Name,
		people.Surname,
		people.Patronymic,
		people.Age,
		people.Gender,
		people.Nation,
		people.Error,
	)
	if err != nil {
		p.logger.LogErrorLvl("exec error", "db", "UpdatePeople", err, nil)

		return err
	}

	count, err := result.RowsAffected()
	if err != nil {
		p.logger.LogErrorLvl("RowsAffected error", "db", "UpdatePeople", err, nil)

		return err
	}
	if count == 0 {
		p.logger.LogTraceLvl("wrong uid", "db", "UpdatePeople", err, people.UID)

		return errors.New("wrong uid")
	}

	return nil
}

// get all peoples from db
func (p peopleStore) GetPeoples() ([]model.MessageToDB, error) {
	peoples := []model.MessageToDB{}
	query := `
		SELECT
			*
		FROM
			people;`
	rows, err := p.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		people := model.MessageToDB{}
		err := rows.Scan(
			&people.UID,
			&people.Name,
			&people.Surname,
			&people.Patronymic,
			&people.Age,
			&people.Gender,
			&people.Nation,
			&people.Error,
		)
		if err != nil {
			p.logger.LogErrorLvl("db error", "db", "GetPeoples", err, nil)

			continue
		}

		peoples = append(peoples, people)
	}

	return peoples, nil
}

// get people from db by uid
func (p peopleStore) GetPeople(uid string) (model.MessageToDB, error) {
	people := model.MessageToDB{}
	query := `
		SELECT
			*
		FROM
			people
		WHERE
			uid = $1`
	err := p.db.QueryRow(query, uid).Scan(
		&people.UID,
		&people.Name,
		&people.Surname,
		&people.Patronymic,
		&people.Age,
		&people.Gender,
		&people.Nation,
		&people.Error,
	)
	if err != nil {
		p.logger.LogTraceLvl("wrong uid", "db", "GetPeople", err, uid)

		return people, err
	}

	return people, nil
}
