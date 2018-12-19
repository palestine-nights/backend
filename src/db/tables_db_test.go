package db

import (
	"fmt"
	"testing"
)

import (
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/suite"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
)

/* --- Suite 1 --- */

type GetTablesSuite struct {
	suite.Suite
	Expected []Table
	DB       *sqlx.DB
	Mock     sqlmock.Sqlmock
}

func (suite *GetTablesSuite) SetupTest() {
	db, mock, err := sqlmock.New()

	if err != nil {
		suite.T().Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	suite.Mock = mock
	suite.DB = sqlx.NewDb(db, "sqlmock")
	suite.Expected = []Table{
		{
			ID:          1,
			Places:      2,
			Description: "Fake Table 1",
			Active:      true,
		},
		{
			ID:          3,
			Places:      6,
			Description: "Fake Table 6",
			Active:      true,
		},
	}
}

func (suite *GetTablesSuite) TestGetTables() {
	defer suite.DB.Close()

	rows := sqlmock.NewRows([]string{"id", "places", "description", "active"})

	for _, table := range suite.Expected {
		rows.AddRow(table.ID, table.Places, table.Description, table.Active)
	}

	// Mock SQL query to return tables.
	suite.Mock.ExpectQuery("^SELECT (.+) FROM tables").WillReturnRows(rows)

	// Call tested function.
	tables, err := Table.GetAll(Table{}, suite.DB)

	// Throw if error.
	if err != nil {
		suite.T().Errorf("expected no error, but got %s instead", err)
	}

	// Make sure that all expectations were met.
	if err := suite.Mock.ExpectationsWereMet(); err != nil {
		suite.T().Errorf("there were unfulfilled expections: %s", err)
	}

	suite.Equal(suite.Expected, *tables)
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestGetTablesSuite(t *testing.T) {
	suite.Run(t, new(GetTablesSuite))
}

/* --- Suite 2 --- */

type FindTableSuite struct {
	suite.Suite
	Expected Table
	DB       *sqlx.DB
	Mock     sqlmock.Sqlmock
}

func (suite *FindTableSuite) SetupTest() {
	db, mock, err := sqlmock.New()

	if err != nil {
		suite.T().Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	suite.Mock = mock
	suite.DB = sqlx.NewDb(db, "sqlmock")
	suite.Expected = Table{
		ID:          1,
		Places:      2,
		Description: "Fake Table 1",
		Active:      true,
	}
}

func (suite *FindTableSuite) TestFindTable() {
	defer suite.DB.Close()

	rows := sqlmock.NewRows([]string{"id", "places", "description", "active"})
	rows.AddRow(suite.Expected.ID, suite.Expected.Places, suite.Expected.Description, suite.Expected.Active)

	// Mock SQL query to return tables.
	suite.Mock.ExpectQuery("^SELECT (.+) FROM tables WHERE id = (.+)$").WillReturnRows(rows)

	// Call tested function.
	table, err := Table.Find(Table{}, suite.DB, suite.Expected.ID)

	// Throw if error.
	if err != nil {
		suite.T().Errorf("expected no error, but got %s instead", err)
	}

	// Make sure that all expectations were met.
	if err := suite.Mock.ExpectationsWereMet(); err != nil {
		suite.T().Errorf("there were unfulfilled expections: %s", err)
	}

	suite.Equal(suite.Expected, *table)
}

func (suite *FindTableSuite) TestFindTableError() {
	defer suite.DB.Close()

	suite.Mock.ExpectQuery("^SELECT (.+) FROM tables WHERE id = (.+)$").
		WillReturnError(fmt.Errorf("SQL Error"))

	// Call tested function.
	table, err := Table.Find(Table{}, suite.DB, suite.Expected.ID)
	suite.Equal(err.Error(),"SQL Error")
	suite.Nil(table)

	// Make sure that all expectations were met.
	if err := suite.Mock.ExpectationsWereMet(); err != nil {
		suite.T().Errorf("there were unfulfilled expections: %s", err)
	}
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func
TestFindTableSuite(t *testing.T) {
	suite.Run(t, new(FindTableSuite))
}
