package externalapi

import "github.com/KoLLlaka/sobes/internal/logger"

type Enrich interface {
	NewPeopleEnrich() PeopleEnrich
}

type enrich struct {
	logger logger.MyLogger
}

// create a new Enrich interface
func NewEnrich(logger logger.MyLogger) Enrich {
	return &enrich{
		logger: logger,
	}
}

// create a new PeopleEnrich interface
func (s *enrich) NewPeopleEnrich() PeopleEnrich {
	return &peopleEnrich{
		logger: s.logger,
	}
}
