package queries

import (
	"strconv"

	"github.com/iainkiloh/examplegoapi/contracts"
)

func GetPersonQuery(id int) (contracts.PersonForFetch, error) {
	queryToExecute := "SELECT * FROM public.person WHERE id=$1"
	dbResult := db.QueryRow(queryToExecute, strconv.Itoa(id))
	response := new(contracts.PersonForFetch)
	err := dbResult.Scan(&response.Id, &response.Name, &response.Title)
	return *response, err
}
