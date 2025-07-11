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
	t.Setenv("APP_NAME", "idm")
	t.Setenv("APP_VERSION", "1.0.0")
	t.Setenv("SSL_SERT", "ssl.cert")
	t.Setenv("SSL_KEY", "ssl.key")
	t.Setenv("KEYCLOAK_JWK_URL", "http://localhost:9990/realms/")
	config := GetConfig("")
	assert.Equal(t, "postgres", config.DbDriverName)
	assert.Equal(t, dsn, config.Dsn)
}

func TestEnvFileExistHaventVarsThenGetEmptyConfig(t *testing.T) {
	_ = assert.New(t)
	defer func() {
		r := recover()
		assert.NotNil(t, r)
		assert.Equal(t, "config validation error: Key: 'Config.DbDriverName' Error:Field validation for "+
			"'DbDriverName' failed on the 'required' tag\nKey: 'Config.Dsn' Error:Field validation for 'Dsn' failed "+
			"on the 'required' tag\nKey: 'Config.AppName' Error:Field validation for 'AppName' failed on the "+
			"'required' tag\nKey: 'Config.AppVersion' Error:Field validation for 'AppVersion' failed on the "+
			"'required' tag\nKey: 'Config.SslSert' Error:Field validation for 'SslSert' failed on the 'required' tag"+
			"\nKey: 'Config.SslKey' Error:Field validation for 'SslKey' failed on the 'required' tag\n"+
			"Key: 'Config.KeycloakJwkUrl' Error:Field validation for 'KeycloakJwkUrl' failed on the 'required' tag", r)
	}()
	t.Setenv("DB_DRIVER_NAME", "")
	t.Setenv("DB_DSN", "")
	t.Setenv("APP_NAME", "")
	t.Setenv("APP_VERSION", "")
	file := createEnvFile(t, "")
	defer os.Remove(file)
	_ = GetConfig(file)
}

func TestEnvFileExistHaventVarsInEnvFileThenGetValidConfig(t *testing.T) {
	_ = assert.New(t)
	t.Setenv("DB_DRIVER_NAME", "postgres")
	t.Setenv("DB_DSN", dsn)
	t.Setenv("APP_NAME", "idm")
	t.Setenv("APP_VERSION", "1.0.0")
	t.Setenv("SSL_SERT", "ssl.cert")
	t.Setenv("SSL_KEY", "ssl.key")
	t.Setenv("KEYCLOAK_JWK_URL", "http://localhost:9990/realms/")
	file := createEnvFile(t, "")
	defer os.Remove(file)
	config := GetConfig(file)
	assert.Equal(t, "postgres", config.DbDriverName)
	assert.Equal(t, dsn, config.Dsn)
}

func TestEnvFileExistHaveVarsInEnvFileThenGetValidConfig(t *testing.T) {
	_ = assert.New(t)
	file := createEnvFile(t, "DB_DRIVER_NAME=random_driver\nDB_DSN=random_dsn\nAPP_NAME=idm\nAPP_VERSION=1.0.0\n"+
		"SSL_SERT=certs/ssl.cert\nSSL_KEY=certs/ssl.key\nKEYCLOAK_JWK_URL=http://localhost:9990/realms/")
	defer os.Remove(file)
	config := GetConfig(file)
	assert.Equal(t, "random_driver", config.DbDriverName)
	assert.Equal(t, "random_dsn", config.Dsn)
}

func TestEnvFileExistHaveVarsInEnvFileAndEnvVarsThenGetValidConfig(t *testing.T) {
	_ = assert.New(t)
	t.Setenv("DB_DRIVER_NAME", "postgres")
	t.Setenv("DB_DSN", dsn)
	file := createEnvFile(t, "DB_DRIVER_NAME=random_driver\nDB_DSN=random_dsn\nAPP_NAME=idm\nAPP_VERSION=1.0.0"+
		"\nSSL_SERT=certs/ssl.cert\nSSL_KEY=certs/ssl.key\nKEYCLOAK_JWK_URL=http://localhost:9990/realms/")
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
