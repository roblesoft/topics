package repository

import (
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	entity "github.com/roblesoft/topics/internal/entity"
	"github.com/roblesoft/topics/pkg/dbmock"
	"github.com/stretchr/testify/assert"
)

func TestCreateUser(t *testing.T) {
	var (
		testSuite = dbmock.NewTestSuit()
		userRepo  = &UserRepository{Db: testSuite.Db}
		user      = &entity.User{Username: "test", Password: "test"}
	)
	// Sql transaction expectetation
	testSuite.Mock.ExpectBegin()
	testSuite.Mock.ExpectQuery(regexp.QuoteMeta(
		`INSERT INTO "users" ("created_at","updated_at","deleted_at","username","password") VALUES ($1,$2,$3,$4,$5) RETURNING "id"`)).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	testSuite.Mock.ExpectCommit()

	t.Run("error nil", func(t *testing.T) {
		err := userRepo.Create(user)
		assert.Nil(t, err)
	})
}

func TestGetUserByName(t *testing.T) {
	var (
		testSuite    = dbmock.NewTestSuit()
		userRepo     = &UserRepository{Db: testSuite.Db}
		expectedUser = &entity.User{Username: "test", Password: "test"}
	)

	testSuite.Mock.MatchExpectationsInOrder(false)
	testSuite.Mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT * FROM "users" WHERE username = $1 AND "users"."deleted_at" IS NULL ORDER BY "users"."id" LIMIT 1`)).
		WithArgs(sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id", "username", "password"}).
			AddRow("1", "test", "$2a$14$ajq8Q7fbtFRQvXpdCq7Jcuy.Rx1h/L4J60Otx.gyNLbAYctGMJ9tK"))

	t.Run("Get user instance and error nil", func(t *testing.T) {
		user, err := userRepo.GetUserByUsername("test")
		assert.Equal(t, user.Username, expectedUser.Username)
		assert.Nil(t, err)
	})
}

func TestGetUserByID(t *testing.T) {
	var (
		testSuite    = dbmock.NewTestSuit()
		userRepo     = &UserRepository{Db: testSuite.Db}
		expectedUser = &entity.User{Username: "test", Password: "test"}
	)

	testSuite.Mock.MatchExpectationsInOrder(false)
	testSuite.Mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT * FROM "users" WHERE id = $1 AND "users"."deleted_at" IS NULL ORDER BY "users"."id" LIMIT 1`)).
		WithArgs(sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id", "username", "password"}).
			AddRow("1", expectedUser.Username, "$2a$14$ajq8Q7fbtFRQvXpdCq7Jcuy.Rx1h/L4J60Otx.gyNLbAYctGMJ9tK"))

	t.Run("Get user instance and error nil", func(t *testing.T) {
		user, err := userRepo.GetUserById(1)
		assert.Equal(t, user.Username, expectedUser.Username)
		assert.Nil(t, err)
	})
}
