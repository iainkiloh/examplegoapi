package queries

import (
	"github.com/iainkiloh/examplegoapi/contracts"
)

func CreatePersonQuery(person *contracts.PersonForCreate) (int, error) {
	id := 0
	queryToExecute := `INSERT INTO public.person (name, title) VALUES ($1,$2) RETURNING Id`
	err := db.QueryRow(queryToExecute, person.Name, person.Title).Scan(&id)
	return id, err
}
