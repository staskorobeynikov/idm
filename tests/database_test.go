package tests

import (
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"idm/inner/common"
	"idm/inner/database"
	"testing"
)

const dsn = "host=127.0.0.1 port=5432 user=postgres password=postgres dbname=db_for_tests sslmode=disable"

func TestGetConnectWithInvalidData(t *testing.T) {
	config := common.Config{
		DbDriverName: "postgres",
		Dsn:          "host=127.0.0.1 port=5432 user=p password=p dbname=db_for_tests sslmode=disable",
	}
	defer func() {
		r := recover()
		assert.NotNil(t, r)
		err, ok := r.(error)
		assert.True(t, ok)
		assert.Equal(t, err.Error(), "pq: password authentication failed for user \"p\"")
	}()
	_ = database.ConnectDbWithCfg(config)
}

func TestGetConnectWithValidData(t *testing.T) {
	_ = assert.New(t)
	config := common.Config{
		DbDriverName: "postgres",
		Dsn:          dsn,
	}

	db := database.ConnectDbWithCfg(config)
	defer db.Close()
	assert.NotNil(t, db)
	assert.IsType(t, &sqlx.DB{}, db)
	assert.NoError(t, db.Ping())
}
