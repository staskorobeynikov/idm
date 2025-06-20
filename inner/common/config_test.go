package common

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

const dsn = "host=127.0.0.1 port=5432 user=postgres password=postgres dbname=postgres sslmode=disable"

func TestEnvFileNotExistThenGetConfigFromEnvVars(t *testing.T) {
	_ = assert.New(t)
	t.Setenv("DB_DRIVER_NAME", "postgres")
	t.Setenv("DB_DSN", dsn)
	config := GetConfig("")
	assert.Equal(t, "postgres", config.DbDriverName)
	assert.Equal(t, dsn, config.Dsn)
}

func TestEnvFileExistHaventVarsThenGetEmptyConfig(t *testing.T) {
	_ = assert.New(t)
	t.Setenv("DB_DRIVER_NAME", "")
	t.Setenv("DB_DSN", "")
	file := createEnvFile(t, "")
	defer os.Remove(file)
	config := GetConfig(file)
	assert.Equal(t, "", config.Dsn)
	assert.Equal(t, "", config.DbDriverName)
}

func TestEnvFileExistHaventVarsInEnvFileThenGetValidConfig(t *testing.T) {
	_ = assert.New(t)
	t.Setenv("DB_DRIVER_NAME", "postgres")
	t.Setenv("DB_DSN", dsn)
	file := createEnvFile(t, "")
	defer os.Remove(file)
	config := GetConfig(file)
	assert.Equal(t, "postgres", config.DbDriverName)
	assert.Equal(t, dsn, config.Dsn)
}

func TestEnvFileExistHaveVarsInEnvFileThenGetValidConfig(t *testing.T) {
	_ = assert.New(t)
	file := createEnvFile(t, "DB_DRIVER_NAME=random_driver\nDB_DSN=random_dsn")
	defer os.Remove(file)
	config := GetConfig(file)
	assert.Equal(t, "random_driver", config.DbDriverName)
	assert.Equal(t, "random_dsn", config.Dsn)
}

func TestEnvFileExistHaveVarsInEnvFileAndEnvVarsThenGetValidConfig(t *testing.T) {
	_ = assert.New(t)
	t.Setenv("DB_DRIVER_NAME", "postgres")
	t.Setenv("DB_DSN", dsn)
	file := createEnvFile(t, "DB_DRIVER_NAME=random_driver\nDB_DSN=random_dsn")
	defer os.Remove(file)
	config := GetConfig(file)
	assert.Equal(t, "postgres", config.DbDriverName)
	assert.Equal(t, dsn, config.Dsn)
}

func createEnvFile(t *testing.T, s string) string {
	f, err := os.CreateTemp(".", ".env")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := f.WriteString(s); err != nil {
		t.Fatal(err)
	}
	if err := f.Close(); err != nil {
		t.Fatal(err)
	}
	return f.Name()
}
