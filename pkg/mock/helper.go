package mock

import (
	"database/sql"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/roblesoft/topics/internal/usecase"
	repo "github.com/roblesoft/topics/internal/usecase/repo"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type TestSuiteConnections struct {
	Port    string
	Db      *gorm.DB
	Mock    sqlmock.Sqlmock
	Service *usecase.Service
}

func NewTestSuit(port string) *TestSuiteConnections {

	server := &TestSuiteConnections{Port: port}

	server.SetUpRouter()

	return server
}

func (ts *TestSuiteConnections) SetUpRouter() {
	gin.SetMode(gin.TestMode)
	ts.SetupMockDb()

	userRepo := &repo.UserRepository{Db: ts.Db}
	ts.Service = usecase.NewService(userRepo)
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
