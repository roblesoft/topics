package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	models "github.com/roblesoft/topics/internal/entity"
	"github.com/roblesoft/topics/pkg/mock"
	"github.com/stretchr/testify/assert"
)

func TestRegisterUser(t *testing.T) {
	var (
		testSuite = mock.NewTestSuit("3000")
		server    = &Server{Port: testSuite.Port, Service: testSuite.Service}
		router    = gin.Default()
	)

	router.POST("/api/v1/users", server.Register)
	server.Router = router

	t.Run("http status ok", func(t *testing.T) {
		user := models.User{Username: "Uriel", Password: "test"}
		// Sql transaction expectetation
		testSuite.Mock.ExpectBegin()
		testSuite.Mock.ExpectQuery(regexp.QuoteMeta(
			`INSERT INTO "users" ("created_at","updated_at","deleted_at","username","password") VALUES ($1,$2,$3,$4,$5) RETURNING "id"`)).
			WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), user.Username, sqlmock.AnyArg()).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
		testSuite.Mock.ExpectCommit()

		// Server request
		var (
			jsonValue, _ = json.Marshal(user)
			w            = httptest.NewRecorder()
			req, _       = http.NewRequest("POST", "/api/v1/users", bytes.NewBuffer(jsonValue))
		)
		server.Router.ServeHTTP(w, req)

		if err := json.Unmarshal(w.Body.Bytes(), &user); err != nil {
			fmt.Println("Can not unmarshal JSON")
		}
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("http status bad request", func(t *testing.T) {
		user := models.User{Username: "", Password: "test"}
		// Sql transaction expectetation
		testSuite.Mock.ExpectBegin()
		testSuite.Mock.ExpectQuery(regexp.QuoteMeta(
			`INSERT INTO "users" ("created_at","updated_at","deleted_at","username","password") VALUES ($1,$2,$3,$4,$5) RETURNING "id"`)).
			WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), user.Username, sqlmock.AnyArg()).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
		testSuite.Mock.ExpectCommit()

		// Server request
		var (
			jsonValue, _ = json.Marshal(user)
			w            = httptest.NewRecorder()
			req, _       = http.NewRequest("POST", "/api/v1/users", bytes.NewBuffer(jsonValue))
		)
		server.Router.ServeHTTP(w, req)

		if err := json.Unmarshal(w.Body.Bytes(), &user); err != nil {
			fmt.Println("Can not unmarshal JSON")
		}
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestLogin(t *testing.T) {
	t.Setenv("TOKEN_HOUR_LIFESPAN", "1")

	var (
		testSuite = mock.NewTestSuit("3000")
		server    = &Server{Port: testSuite.Port, Service: testSuite.Service}
		router    = gin.Default()
	)

	router.POST("/api/v1/users/login", server.Login)
	server.Router = router

	t.Run("http status ok", func(t *testing.T) {
		user := models.User{Username: "newuser", Password: "secret"}
		var (
			jsonValue, _ = json.Marshal(&user)
			w            = httptest.NewRecorder()
			req, _       = http.NewRequest("POST", "/api/v1/users/login", bytes.NewBuffer(jsonValue))
		)
		testSuite.Mock.MatchExpectationsInOrder(false)
		testSuite.Mock.ExpectQuery(regexp.QuoteMeta(
			`SELECT * FROM "users" WHERE username = $1 AND "users"."deleted_at" IS NULL ORDER BY "users"."id" LIMIT 1`)).
			WithArgs(sqlmock.AnyArg()).
			WillReturnRows(sqlmock.NewRows([]string{"id", "username", "password"}).
				AddRow("1", "newuser", "$2a$14$ajq8Q7fbtFRQvXpdCq7Jcuy.Rx1h/L4J60Otx.gyNLbAYctGMJ9tK"))
		server.Router.ServeHTTP(w, req)

		if err := json.Unmarshal(w.Body.Bytes(), &user); err != nil {
			fmt.Println("Can not unmarshal JSON")
		}
		assert.Equal(t, http.StatusOK, w.Code)
	})
}
