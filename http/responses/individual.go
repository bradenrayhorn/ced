package responses

import "github.com/bradenrayhorn/ced/ced"

type Individual struct {
	ID           ced.ID `json:"id"`
	Name         string `json:"name"`
	Response     bool   `json:"response"`
	HasResponded bool   `json:"has_responded"`
}

func FromIndividual(individual ced.Individual) Individual {
	return Individual{
		ID:           individual.ID,
		Name:         string(individual.Name),
		Response:     individual.Response,
		HasResponded: individual.HasResponded,
	}
}
