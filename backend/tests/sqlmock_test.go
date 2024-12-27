// This file ensures the database query logic functions as expected.
// I'm using sqlmock to simulate database interactions, validating queries and their results without messing up the actual database.

package tests

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestDatabaseQuery(t *testing.T) {

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error initializing sqlmock: %s", err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id", "name", "statement"}).
		AddRow(1, "John Doe", 100)
	mock.ExpectQuery("SELECT \\* FROM users").WillReturnRows(rows)

	var users []struct {
		ID        int
		Name      string
		Statement int
	}
	rowsData, err := db.Query("SELECT * FROM users")
	if err != nil {
		t.Fatalf("Error querying database: %s", err)
	}
	defer rowsData.Close()

	for rowsData.Next() {
		var user struct {
			ID        int
			Name      string
			Statement int
		}
		err := rowsData.Scan(&user.ID, &user.Name, &user.Statement)
		assert.NoError(t, err, "Error scanning row")
		users = append(users, user)
	}

	assert.Len(t, users, 1, "Expected 1 user")
	assert.Equal(t, "John Doe", users[0].Name, "Unexpected user name")

	assert.NoError(t, mock.ExpectationsWereMet(), "SQL expectations were not met")
}
