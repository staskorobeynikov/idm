package tests

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"idm/inner/common"
	"idm/inner/database"
	"idm/inner/employee"
	"idm/inner/role"
	"idm/inner/validator"
	"idm/inner/web"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestEmployeeControllerFindWithOffset(t *testing.T) {
	a := assert.New(t)
	var db = database.ConnectDb()
	var clearDatabase = func() {
		db.MustExec("DELETE FROM employee")
	}
	defer func() {
		if r := recover(); r != nil {
			clearDatabase()
		}
	}()
	defer func() {
		db.MustExec("DELETE FROM role")
	}()
	var employeeRepository = employee.NewRepository(db)
	var emplFixture = Fixture{
		employees: employeeRepository,
		db:        db,
	}
	var roleRepository = role.NewRepository(db)
	var roleFixture = NewRoleFixture(roleRepository)
	var newRoleId = roleFixture.Role("Test Name")
	_ = emplFixture.CreateDatabase(db)
	v := validator.New()
	var employeeService = employee.NewService(employeeRepository, v)
	server := web.NewServer()
	var employeeController = employee.NewController(server, employeeService)
	employeeController.RegisterRoutes()
	t.Run("get employees with offset - page 0, size 3", func(t *testing.T) {
		_ = emplFixture.Employee("Test Name", newRoleId)
		_ = emplFixture.Employee("Test Name 1", newRoleId)
		_ = emplFixture.Employee("Test Name 2", newRoleId)
		_ = emplFixture.Employee("Test Name 3", newRoleId)
		_ = emplFixture.Employee("Test Name 4", newRoleId)
		var request = httptest.NewRequest(http.MethodGet, "/api/v1/employees/page?pageNumber=0&pageSize=3", nil)
		resp, err := server.App.Test(request)
		a.Nil(err)
		a.Equal(http.StatusOK, resp.StatusCode)
		a.NotEmpty(resp)
		bytesData, err := io.ReadAll(resp.Body)
		a.Nil(err)
		var responseBody common.Response[employee.PageResponse]
		err = json.Unmarshal(bytesData, &responseBody)
		a.Nil(err)
		a.NotEmpty(responseBody)
		a.Equal(3, len(responseBody.Data.Result))
		clearDatabase()
	})
	t.Run("get employees with offset - page 1, size 3", func(t *testing.T) {
		_ = emplFixture.Employee("Test Name", newRoleId)
		_ = emplFixture.Employee("Test Name 1", newRoleId)
		_ = emplFixture.Employee("Test Name 2", newRoleId)
		_ = emplFixture.Employee("Test Name 3", newRoleId)
		_ = emplFixture.Employee("Test Name 4", newRoleId)
		var request = httptest.NewRequest(http.MethodGet, "/api/v1/employees/page?pageNumber=1&pageSize=3", nil)
		resp, err := server.App.Test(request)
		a.Nil(err)
		a.Equal(http.StatusOK, resp.StatusCode)
		a.NotEmpty(resp)
		bytesData, err := io.ReadAll(resp.Body)
		a.Nil(err)
		var responseBody common.Response[employee.PageResponse]
		err = json.Unmarshal(bytesData, &responseBody)
		a.Nil(err)
		a.NotEmpty(responseBody)
		a.Equal(2, len(responseBody.Data.Result))
		clearDatabase()
	})
	t.Run("get employees with offset - page 2, size 3", func(t *testing.T) {
		_ = emplFixture.Employee("Test Name", newRoleId)
		_ = emplFixture.Employee("Test Name 1", newRoleId)
		_ = emplFixture.Employee("Test Name 2", newRoleId)
		_ = emplFixture.Employee("Test Name 3", newRoleId)
		_ = emplFixture.Employee("Test Name 4", newRoleId)
		var request = httptest.NewRequest(http.MethodGet, "/api/v1/employees/page?pageNumber=2&pageSize=3", nil)
		resp, err := server.App.Test(request)
		a.Nil(err)
		a.Equal(http.StatusOK, resp.StatusCode)
		a.NotEmpty(resp)
		bytesData, err := io.ReadAll(resp.Body)
		a.Nil(err)
		var responseBody common.Response[employee.PageResponse]
		err = json.Unmarshal(bytesData, &responseBody)
		a.Nil(err)
		a.NotEmpty(responseBody)
		a.Equal(0, len(responseBody.Data.Result))
		clearDatabase()
	})
	t.Run("get employees with offset - page -1, size 3", func(t *testing.T) {
		_ = emplFixture.Employee("Test Name", newRoleId)
		_ = emplFixture.Employee("Test Name 1", newRoleId)
		_ = emplFixture.Employee("Test Name 2", newRoleId)
		_ = emplFixture.Employee("Test Name 3", newRoleId)
		_ = emplFixture.Employee("Test Name 4", newRoleId)
		var request = httptest.NewRequest(http.MethodGet, "/api/v1/employees/page?pageNumber=-1&pageSize=3", nil)
		resp, err := server.App.Test(request)
		a.Nil(err)
		a.Equal(http.StatusBadRequest, resp.StatusCode)
		a.NotEmpty(resp)
		bytesData, err := io.ReadAll(resp.Body)
		a.Nil(err)
		var responseBody common.Response[employee.PageResponse]
		err = json.Unmarshal(bytesData, &responseBody)
		a.Nil(err)
		a.NotEmpty(responseBody)
		want := "Key: 'PageRequest.PageNumber' Error:Field validation for 'PageNumber' failed on the 'min' tag"
		a.Equal(want, responseBody.Message)
		clearDatabase()
	})
	t.Run("get employees with offset - page no, size 3", func(t *testing.T) {
		_ = emplFixture.Employee("Test Name", newRoleId)
		_ = emplFixture.Employee("Test Name 1", newRoleId)
		_ = emplFixture.Employee("Test Name 2", newRoleId)
		_ = emplFixture.Employee("Test Name 3", newRoleId)
		_ = emplFixture.Employee("Test Name 4", newRoleId)
		var request = httptest.NewRequest(http.MethodGet, "/api/v1/employees/page?pageNumber=&pageSize=3", nil)
		resp, err := server.App.Test(request)
		a.Nil(err)
		a.Equal(http.StatusOK, resp.StatusCode)
		a.NotEmpty(resp)
		bytesData, err := io.ReadAll(resp.Body)
		a.Nil(err)
		var responseBody common.Response[employee.PageResponse]
		err = json.Unmarshal(bytesData, &responseBody)
		a.Nil(err)
		a.NotEmpty(responseBody)
		a.Equal(3, len(responseBody.Data.Result))
		clearDatabase()
	})
	t.Run("get employees with offset - page 0, size no", func(t *testing.T) {
		_ = emplFixture.Employee("Test Name", newRoleId)
		_ = emplFixture.Employee("Test Name 1", newRoleId)
		_ = emplFixture.Employee("Test Name 2", newRoleId)
		_ = emplFixture.Employee("Test Name 3", newRoleId)
		_ = emplFixture.Employee("Test Name 4", newRoleId)
		var request = httptest.NewRequest(http.MethodGet, "/api/v1/employees/page?pageNumber=0&pageSize=", nil)
		resp, err := server.App.Test(request)
		a.Nil(err)
		a.Equal(http.StatusOK, resp.StatusCode)
		a.NotEmpty(resp)
		bytesData, err := io.ReadAll(resp.Body)
		a.Nil(err)
		var responseBody common.Response[employee.PageResponse]
		err = json.Unmarshal(bytesData, &responseBody)
		a.Nil(err)
		a.NotEmpty(responseBody)
		a.Equal(5, len(responseBody.Data.Result))
		clearDatabase()
	})
}
