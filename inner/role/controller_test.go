package role

import (
	"encoding/json"
	"github.com/gofiber/fiber/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
	"idm/inner/common"
	"idm/inner/web"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type MockService struct {
	mock.Mock
}

func (svc *MockService) Save(request CreateRequest) (Response, error) {
	args := svc.Called(request)
	return args.Get(0).(Response), args.Error(1)
}

func (svc *MockService) FindById(request IdRequest) (Response, error) {
	args := svc.Called(request)
	return args.Get(0).(Response), args.Error(1)
}

func (svc *MockService) FindAll() ([]Response, error) {
	args := svc.Called()
	return args.Get(0).([]Response), args.Error(1)
}

func (svc *MockService) FindByIds(request IdsRequest) ([]Response, error) {
	args := svc.Called(request)
	return args.Get(0).([]Response), args.Error(1)
}

func (svc *MockService) DeleteById(request IdRequest) error {
	args := svc.Called(request)
	return args.Error(0)
}

func (svc *MockService) DeleteByIds(request IdsRequest) error {
	args := svc.Called(request)
	return args.Error(0)
}

var logger = &common.Logger{Logger: zap.NewNop()}

func TestCreateRole(t *testing.T) {
	var a = assert.New(t)
	t.Run("create role without error", func(t *testing.T) {
		server := web.NewServer(false)
		var svc = new(MockService)
		var controller = NewController(server, svc, logger)
		controller.RegisterRoutes()
		var body = strings.NewReader("{\"name\": \"john doe\"}")
		var request = httptest.NewRequest(fiber.MethodPost, "/api/v1/roles", body)
		request.Header.Add("Content-Type", "application/json")
		svc.On("Save", mock.AnythingOfType("CreateRequest")).Return(Response{Id: int64(123)}, nil)
		resp, err := server.App.Test(request)
		a.Nil(err)
		a.NotEmpty(resp)
		a.Equal(http.StatusOK, resp.StatusCode)
		bytesData, err := io.ReadAll(resp.Body)
		a.Nil(err)
		var responseBody common.Response[int64]
		err = json.Unmarshal(bytesData, &responseBody)
		a.Nil(err)
		a.Equal(int64(123), responseBody.Data)
		a.True(responseBody.Success)
		a.Empty(responseBody.Message)
	})
	t.Run("create role validation error - name required", func(t *testing.T) {
		server := web.NewServer(false)
		var svc = new(MockService)
		var controller = NewController(server, svc, logger)
		controller.RegisterRoutes()
		var body = strings.NewReader("{\"name\": \"\"}")
		var request = httptest.NewRequest(fiber.MethodPost, "/api/v1/roles", body)
		request.Header.Add("Content-Type", "application/json")
		message := "Field validation for 'Name' failed on the 'required' tag"
		svc.On("Save", mock.AnythingOfType("CreateRequest")).Return(Response{}, common.RequestValidationError{
			Message: message,
		})
		resp, err := server.App.Test(request)
		a.Nil(err)
		a.NotEmpty(resp)
		a.Equal(http.StatusBadRequest, resp.StatusCode)
		bytesData, err := io.ReadAll(resp.Body)
		a.Nil(err)
		var responseBody common.Response[int64]
		err = json.Unmarshal(bytesData, &responseBody)
		a.Nil(err)
		a.Equal(message, responseBody.Message)
	})
}

func TestFindRoleById(t *testing.T) {
	var a = assert.New(t)
	t.Run("find role by id", func(t *testing.T) {
		server := web.NewServer(false)
		var svc = new(MockService)
		var controller = NewController(server, svc, logger)
		controller.RegisterRoutes()
		var request = httptest.NewRequest(fiber.MethodGet, "/api/v1/roles/123", nil)
		request.Header.Add("Content-Type", "application/json")
		svc.On("FindById", mock.AnythingOfType("IdRequest")).Return(Response{Id: int64(123)}, nil)
		resp, err := server.App.Test(request)
		a.Nil(err)
		a.NotEmpty(resp)
		a.Equal(http.StatusOK, resp.StatusCode)
		bytesData, err := io.ReadAll(resp.Body)
		a.Nil(err)
		var responseBody common.Response[Response]
		err = json.Unmarshal(bytesData, &responseBody)
		a.Nil(err)
		a.Equal(int64(123), responseBody.Data.Id)
	})
	t.Run("find role - incorrect id", func(t *testing.T) {
		server := web.NewServer(false)
		var svc = new(MockService)
		var controller = NewController(server, svc, logger)
		controller.RegisterRoutes()
		var request = httptest.NewRequest(fiber.MethodGet, "/api/v1/roles/ffff", nil)
		request.Header.Add("Content-Type", "application/json")
		resp, err := server.App.Test(request)
		a.Nil(err)
		a.NotEmpty(resp)
		a.Equal(http.StatusBadRequest, resp.StatusCode)
	})
	t.Run("find role - validation error", func(t *testing.T) {
		server := web.NewServer(false)
		var svc = new(MockService)
		var controller = NewController(server, svc, logger)
		controller.RegisterRoutes()
		var request = httptest.NewRequest(fiber.MethodGet, "/api/v1/roles/0", nil)
		request.Header.Add("Content-Type", "application/json")
		message := "Field validation for 'Name' failed on the 'min' tag"
		svc.On("FindById", mock.AnythingOfType("IdRequest")).Return(Response{}, common.RequestValidationError{
			Message: message,
		})
		resp, err := server.App.Test(request)
		a.Nil(err)
		a.NotEmpty(resp)
		a.Equal(http.StatusBadRequest, resp.StatusCode)
		bytesData, err := io.ReadAll(resp.Body)
		a.Nil(err)
		var responseBody common.Response[Response]
		err = json.Unmarshal(bytesData, &responseBody)
		a.Nil(err)
		a.Equal(message, responseBody.Message)
	})
	t.Run("find role - not found error", func(t *testing.T) {
		server := web.NewServer(false)
		var svc = new(MockService)
		var controller = NewController(server, svc, logger)
		controller.RegisterRoutes()
		var request = httptest.NewRequest(fiber.MethodGet, "/api/v1/roles/123", nil)
		request.Header.Add("Content-Type", "application/json")
		message := "error finding role with id 123"
		svc.On("FindById", mock.AnythingOfType("IdRequest")).Return(Response{}, common.NotFoundError{
			Message: message,
		})
		resp, err := server.App.Test(request)
		a.Nil(err)
		a.NotEmpty(resp)
		a.Equal(http.StatusOK, resp.StatusCode)
		bytesData, err := io.ReadAll(resp.Body)
		a.Nil(err)
		var responseBody common.Response[Response]
		err = json.Unmarshal(bytesData, &responseBody)
		a.Nil(err)
		a.Equal(message, responseBody.Message)
	})
}

func TestFindAllRoles(t *testing.T) {
	var a = assert.New(t)
	t.Run("find all roles", func(t *testing.T) {
		server := web.NewServer(false)
		var svc = new(MockService)
		var controller = NewController(server, svc, logger)
		controller.RegisterRoutes()
		var request = httptest.NewRequest(fiber.MethodGet, "/api/v1/roles", nil)
		request.Header.Add("Content-Type", "application/json")
		responses := []Response{
			{Id: int64(123)},
			{Id: int64(124)},
			{Id: int64(125)},
		}
		svc.On("FindAll").Return(
			responses, nil)
		resp, err := server.App.Test(request)
		a.Nil(err)
		a.NotEmpty(resp)
		a.Equal(http.StatusOK, resp.StatusCode)
		bytesData, err := io.ReadAll(resp.Body)
		a.Nil(err)
		var responseBody common.Response[[]Response]
		err = json.Unmarshal(bytesData, &responseBody)
		a.Nil(err)
		a.Equal(responses, responseBody.Data)
	})
	t.Run("find all with error", func(t *testing.T) {
		server := web.NewServer(false)
		var svc = new(MockService)
		var controller = NewController(server, svc, logger)
		controller.RegisterRoutes()
		var request = httptest.NewRequest(fiber.MethodGet, "/api/v1/roles", nil)
		request.Header.Add("Content-Type", "application/json")
		message := "error finding all roles"
		svc.On("FindAll").Return([]Response{}, common.NotFoundError{
			Message: message,
		})
		resp, err := server.App.Test(request)
		a.Nil(err)
		a.NotEmpty(resp)
		a.Equal(http.StatusOK, resp.StatusCode)
		bytesData, err := io.ReadAll(resp.Body)
		a.Nil(err)
		var responseBody common.Response[[]Response]
		err = json.Unmarshal(bytesData, &responseBody)
		a.Nil(err)
		a.Equal(message, responseBody.Message)
	})
}

func TestFindRolesByIds(t *testing.T) {
	var a = assert.New(t)
	t.Run("find roles by ids", func(t *testing.T) {
		server := web.NewServer(false)
		var svc = new(MockService)
		var controller = NewController(server, svc, logger)
		controller.RegisterRoutes()
		var request = httptest.NewRequest(fiber.MethodGet, "/api/v1/roles/find?ids=123,124,125", nil)
		request.Header.Add("Content-Type", "application/json")
		responses := []Response{
			{Id: int64(123)},
			{Id: int64(124)},
			{Id: int64(125)},
		}
		svc.On("FindByIds", mock.AnythingOfType("IdsRequest")).Return(responses, nil)
		resp, err := server.App.Test(request)
		a.Nil(err)
		a.NotEmpty(resp)
		a.Equal(http.StatusOK, resp.StatusCode)
		bytesData, err := io.ReadAll(resp.Body)
		a.Nil(err)
		var responseBody common.Response[[]Response]
		err = json.Unmarshal(bytesData, &responseBody)
		a.Nil(err)
		a.Equal(responses, responseBody.Data)
	})
	t.Run("find roles by ids - error parsing", func(t *testing.T) {
		server := web.NewServer(false)
		var svc = new(MockService)
		var controller = NewController(server, svc, logger)
		controller.RegisterRoutes()
		var request = httptest.NewRequest(fiber.MethodGet, "/api/v1/roles/find?ids=fff,124,125", nil)
		request.Header.Add("Content-Type", "application/json")
		responses := []Response{
			{Id: int64(123)},
			{Id: int64(124)},
			{Id: int64(125)},
		}
		message := "strconv.ParseInt: parsing \"fff\": invalid syntax"
		svc.On("FindByIds", mock.AnythingOfType("IdsRequest")).Return(responses, nil)
		resp, err := server.App.Test(request)
		a.Nil(err)
		a.NotEmpty(resp)
		a.Equal(http.StatusBadRequest, resp.StatusCode)
		bytesData, err := io.ReadAll(resp.Body)
		a.Nil(err)
		var responseBody common.Response[[]Response]
		err = json.Unmarshal(bytesData, &responseBody)
		a.Nil(err)
		a.Equal(message, responseBody.Message)
	})
	t.Run("find roles by ids - validation error", func(t *testing.T) {
		server := web.NewServer(false)
		var svc = new(MockService)
		var controller = NewController(server, svc, logger)
		controller.RegisterRoutes()
		var request = httptest.NewRequest(fiber.MethodGet, "/api/v1/roles/find?ids=123,124,125", nil)
		request.Header.Add("Content-Type", "application/json")
		message := "role not found"
		svc.On("FindByIds", mock.AnythingOfType("IdsRequest")).Return([]Response{}, common.RequestValidationError{
			Message: message,
		})
		resp, err := server.App.Test(request)
		a.Nil(err)
		a.NotEmpty(resp)
		a.Equal(http.StatusBadRequest, resp.StatusCode)
		bytesData, err := io.ReadAll(resp.Body)
		a.Nil(err)
		var responseBody common.Response[[]Response]
		err = json.Unmarshal(bytesData, &responseBody)
		a.Nil(err)
		a.Equal(message, responseBody.Message)
	})
	t.Run("find roles by ids - not found error", func(t *testing.T) {
		server := web.NewServer(false)
		var svc = new(MockService)
		var controller = NewController(server, svc, logger)
		controller.RegisterRoutes()
		var request = httptest.NewRequest(fiber.MethodGet, "/api/v1/roles/find?ids=123,124,125", nil)
		request.Header.Add("Content-Type", "application/json")
		message := "error finding roles by ids"
		svc.On("FindByIds", mock.AnythingOfType("IdsRequest")).Return([]Response{}, common.NotFoundError{
			Message: message,
		})
		resp, err := server.App.Test(request)
		a.Nil(err)
		a.NotEmpty(resp)
		a.Equal(http.StatusOK, resp.StatusCode)
		bytesData, err := io.ReadAll(resp.Body)
		a.Nil(err)
		var responseBody common.Response[[]Response]
		err = json.Unmarshal(bytesData, &responseBody)
		a.Nil(err)
		a.Equal(message, responseBody.Message)
	})
}

func TestDeleteRoleById(t *testing.T) {
	var a = assert.New(t)
	t.Run("find role by id", func(t *testing.T) {
		server := web.NewServer(false)
		var svc = new(MockService)
		var controller = NewController(server, svc, logger)
		controller.RegisterRoutes()
		var request = httptest.NewRequest(fiber.MethodDelete, "/api/v1/roles/123", nil)
		request.Header.Add("Content-Type", "application/json")
		svc.On("DeleteById", mock.AnythingOfType("IdRequest")).Return(nil)
		resp, err := server.App.Test(request)
		a.Nil(err)
		a.NotEmpty(resp)
		a.Equal(http.StatusOK, resp.StatusCode)
		bytesData, err := io.ReadAll(resp.Body)
		a.Nil(err)
		var responseBody common.Response[Response]
		err = json.Unmarshal(bytesData, &responseBody)
		a.Nil(err)
		a.Equal(int64(123), responseBody.Data.Id)
	})
	t.Run("find role - incorrect id", func(t *testing.T) {
		server := web.NewServer(false)
		var svc = new(MockService)
		var controller = NewController(server, svc, logger)
		controller.RegisterRoutes()
		var request = httptest.NewRequest(fiber.MethodDelete, "/api/v1/roles/ffff", nil)
		request.Header.Add("Content-Type", "application/json")
		resp, err := server.App.Test(request)
		a.Nil(err)
		a.NotEmpty(resp)
		a.Equal(http.StatusBadRequest, resp.StatusCode)
	})
	t.Run("find role - validation error", func(t *testing.T) {
		server := web.NewServer(false)
		var svc = new(MockService)
		var controller = NewController(server, svc, logger)
		controller.RegisterRoutes()
		var request = httptest.NewRequest(fiber.MethodDelete, "/api/v1/roles/0", nil)
		request.Header.Add("Content-Type", "application/json")
		message := "Field validation for 'Name' failed on the 'min' tag"
		svc.On("DeleteById", mock.AnythingOfType("IdRequest")).Return(common.RequestValidationError{
			Message: message,
		})
		resp, err := server.App.Test(request)
		a.Nil(err)
		a.NotEmpty(resp)
		a.Equal(http.StatusBadRequest, resp.StatusCode)
		bytesData, err := io.ReadAll(resp.Body)
		a.Nil(err)
		var responseBody common.Response[Response]
		err = json.Unmarshal(bytesData, &responseBody)
		a.Nil(err)
		a.Equal(message, responseBody.Message)
	})
	t.Run("find role - not found error", func(t *testing.T) {
		server := web.NewServer(false)
		var svc = new(MockService)
		var controller = NewController(server, svc, logger)
		controller.RegisterRoutes()
		var request = httptest.NewRequest(fiber.MethodDelete, "/api/v1/roles/123", nil)
		request.Header.Add("Content-Type", "application/json")
		message := "error finding role with id 123"
		svc.On("DeleteById", mock.AnythingOfType("IdRequest")).Return(common.NotFoundError{
			Message: message,
		})
		resp, err := server.App.Test(request)
		a.Nil(err)
		a.NotEmpty(resp)
		a.Equal(http.StatusOK, resp.StatusCode)
		bytesData, err := io.ReadAll(resp.Body)
		a.Nil(err)
		var responseBody common.Response[Response]
		err = json.Unmarshal(bytesData, &responseBody)
		a.Nil(err)
		a.Equal(message, responseBody.Message)
	})
}

func TestDeleteRolesByIds(t *testing.T) {
	var a = assert.New(t)
	t.Run("delete roles by ids", func(t *testing.T) {
		server := web.NewServer(false)
		var svc = new(MockService)
		var controller = NewController(server, svc, logger)
		controller.RegisterRoutes()
		var request = httptest.NewRequest(fiber.MethodDelete, "/api/v1/roles/delete?ids=123,124,125", nil)
		request.Header.Add("Content-Type", "application/json")
		responses := []Response{
			{Id: int64(123)},
			{Id: int64(124)},
			{Id: int64(125)},
		}
		svc.On("DeleteByIds", mock.AnythingOfType("IdsRequest")).Return(nil)
		resp, err := server.App.Test(request)
		a.Nil(err)
		a.NotEmpty(resp)
		a.Equal(http.StatusOK, resp.StatusCode)
		bytesData, err := io.ReadAll(resp.Body)
		a.Nil(err)
		var responseBody common.Response[[]Response]
		err = json.Unmarshal(bytesData, &responseBody)
		a.Nil(err)
		a.Equal(responses, responseBody.Data)
	})
	t.Run("delete roles by ids - error parsing", func(t *testing.T) {
		server := web.NewServer(false)
		var svc = new(MockService)
		var controller = NewController(server, svc, logger)
		controller.RegisterRoutes()
		var request = httptest.NewRequest(fiber.MethodDelete, "/api/v1/roles/delete?ids=fff,124,125", nil)
		request.Header.Add("Content-Type", "application/json")
		message := "strconv.ParseInt: parsing \"fff\": invalid syntax"
		svc.On("DeleteByIds", mock.AnythingOfType("IdsRequest")).Return(nil)
		resp, err := server.App.Test(request)
		a.Nil(err)
		a.NotEmpty(resp)
		a.Equal(http.StatusBadRequest, resp.StatusCode)
		bytesData, err := io.ReadAll(resp.Body)
		a.Nil(err)
		var responseBody common.Response[[]Response]
		err = json.Unmarshal(bytesData, &responseBody)
		a.Nil(err)
		a.Equal(message, responseBody.Message)
	})
	t.Run("delete roles by ids - validation error", func(t *testing.T) {
		server := web.NewServer(false)
		var svc = new(MockService)
		var controller = NewController(server, svc, logger)
		controller.RegisterRoutes()
		var request = httptest.NewRequest(fiber.MethodDelete, "/api/v1/roles/delete?ids=123,124,125", nil)
		request.Header.Add("Content-Type", "application/json")
		message := "role not found"
		svc.On("DeleteByIds", mock.AnythingOfType("IdsRequest")).Return(common.RequestValidationError{
			Message: message,
		})
		resp, err := server.App.Test(request)
		a.Nil(err)
		a.NotEmpty(resp)
		a.Equal(http.StatusBadRequest, resp.StatusCode)
		bytesData, err := io.ReadAll(resp.Body)
		a.Nil(err)
		var responseBody common.Response[[]Response]
		err = json.Unmarshal(bytesData, &responseBody)
		a.Nil(err)
		a.Equal(message, responseBody.Message)
	})
	t.Run("delete roles by ids - not found error", func(t *testing.T) {
		server := web.NewServer(false)
		var svc = new(MockService)
		var controller = NewController(server, svc, logger)
		controller.RegisterRoutes()
		var request = httptest.NewRequest(fiber.MethodDelete, "/api/v1/roles/delete?ids=123,124,125", nil)
		request.Header.Add("Content-Type", "application/json")
		message := "error finding roles by ids"
		svc.On("DeleteByIds", mock.AnythingOfType("IdsRequest")).Return(common.NotFoundError{
			Message: message,
		})
		resp, err := server.App.Test(request)
		a.Nil(err)
		a.NotEmpty(resp)
		a.Equal(http.StatusOK, resp.StatusCode)
		bytesData, err := io.ReadAll(resp.Body)
		a.Nil(err)
		var responseBody common.Response[[]Response]
		err = json.Unmarshal(bytesData, &responseBody)
		a.Nil(err)
		a.Equal(message, responseBody.Message)
	})
}
