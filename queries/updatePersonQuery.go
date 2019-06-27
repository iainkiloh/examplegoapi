package queries

import (
	"github.com/iainkiloh/examplegoapi/contracts"
)

func UpdatePersonQuery(person *contracts.PersonForUpdate) (int64, error) {
	queryToExecute := `UPDATE public.person SET Name = $1, Title = $2 WHERE id=$3`
	info, err := db.Exec(queryToExecute, person.Name, person.Title, person.Id)
	count, _ := info.RowsAffected()
	return count, err
}
