package dbmock

import (
	"database/sql"

	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type TestSuiteConnections struct {
	Port string
	Db   *gorm.DB
	Mock sqlmock.Sqlmock
}

func NewTestSuit() *TestSuiteConnections {
	service := &TestSuiteConnections{}
	service.SetupMockDb()

	return service
}

func (ts *TestSuiteConnections) SetupMockDb() {
	var (
		db  *sql.DB
		err error
	)

	db, ts.Mock, err = sqlmock.New()

	if err != nil {
		panic(err)
	}

	if db == nil {
		panic(err)
	}

	if ts.Mock == nil {
		panic("mock db is null")
	}

	dialector := postgres.New(postgres.Config{
		DSN:                  "sqlmock_db_0",
		DriverName:           "postgres",
		Conn:                 db,
		PreferSimpleProtocol: true,
	})
	ts.Db, err = gorm.Open(dialector, &gorm.Config{})

	if err != nil {
		panic(err)
	}

	if ts.Db == nil {
		panic("gorm db is null")
	}
}
