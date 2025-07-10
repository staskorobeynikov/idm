package employee

import (
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"idm/inner/common"
	"idm/inner/web"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
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

func (svc *MockService) FindWithOffset(request PageRequest) (PageResponse, error) {
	args := svc.Called(request)
	return args.Get(0).(PageResponse), args.Error(1)
}

func (svc *MockService) DeleteById(request IdRequest) error {
	args := svc.Called(request)
	return args.Error(0)
}

func (svc *MockService) DeleteByIds(request IdsRequest) error {
	args := svc.Called(request)
	return args.Error(0)
}

func TestCreateEmployee(t *testing.T) {
	var a = assert.New(t)
	t.Run("create employee without error", func(t *testing.T) {
		var claims = &web.IdmClaims{
			RealmAccess: web.RealmAccessClaims{Roles: []string{web.IdmAdmin}},
		}
		var auth = func(c *fiber.Ctx) error {
			c.Locals(web.JwtKey, &jwt.Token{Claims: claims})
			return c.Next()
		}
		server := web.NewServer()
		server.GroupApiV1.Use(auth)
		var svc = new(MockService)
		var controller = NewController(server, svc)
		controller.RegisterRoutes()
		var body = strings.NewReader("{\"name\": \"john doe\", \"role_id\": 1}")
		var request = httptest.NewRequest(fiber.MethodPost, "/api/v1/employees", body)
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
	t.Run("create employee validation error - name required", func(t *testing.T) {
		var claims = &web.IdmClaims{
			RealmAccess: web.RealmAccessClaims{Roles: []string{web.IdmAdmin}},
		}
		var auth = func(c *fiber.Ctx) error {
			c.Locals(web.JwtKey, &jwt.Token{Claims: claims})
			return c.Next()
		}
		server := web.NewServer()
		server.GroupApiV1.Use(auth)
		var svc = new(MockService)
		var controller = NewController(server, svc)
		controller.RegisterRoutes()
		var body = strings.NewReader("{\"name\": \"\", \"role_id\": 1}")
		var request = httptest.NewRequest(fiber.MethodPost, "/api/v1/employees", body)
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
	t.Run("create employee validation error - short name", func(t *testing.T) {
		var claims = &web.IdmClaims{
			RealmAccess: web.RealmAccessClaims{Roles: []string{web.IdmAdmin}},
		}
		var auth = func(c *fiber.Ctx) error {
			c.Locals(web.JwtKey, &jwt.Token{Claims: claims})
			return c.Next()
		}
		server := web.NewServer()
		server.GroupApiV1.Use(auth)
		var svc = new(MockService)
		var controller = NewController(server, svc)
		controller.RegisterRoutes()
		var body = strings.NewReader("{\"name\": \"john doe\", \"role_id\": 1}")
		var request = httptest.NewRequest(fiber.MethodPost, "/api/v1/employees", body)
		request.Header.Add("Content-Type", "application/json")
		message := "employee already exists: john doe"
		svc.On("Save", mock.AnythingOfType("CreateRequest")).Return(Response{}, common.AlreadyExistsError{
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
	t.Run("create employee without token", func(t *testing.T) {
		fakeAuth := func(c *fiber.Ctx) error {
			authHeader := c.Get("Authorization")
			if authHeader == "" {
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
					"success": false,
					"error":   "missing or malformed JWT",
				})
			}
			c.Locals(web.JwtKey, &jwt.Token{})
			return c.Next()
		}
		server := web.NewServer()
		server.GroupApiV1.Use(fakeAuth)
		var svc = new(MockService)
		var controller = NewController(server, svc)
		controller.RegisterRoutes()
		var body = strings.NewReader("{\"name\": \"john doe\", \"role_id\": 1}")
		var request = httptest.NewRequest(fiber.MethodPost, "/api/v1/employees", body)
		request.Header.Add("Content-Type", "application/json")
		svc.On("Save", mock.AnythingOfType("CreateRequest")).Return(Response{Id: int64(123)}, nil)
		message := "missing or malformed JWT"
		resp, err := server.App.Test(request)
		a.Nil(err)
		a.NotEmpty(resp)
		a.Equal(http.StatusUnauthorized, resp.StatusCode)
		bytesData, err := io.ReadAll(resp.Body)
		a.Nil(err)
		var responseBody common.Response[Response]
		err = json.Unmarshal(bytesData, &responseBody)
		a.Nil(err)
		a.Equal(message, responseBody.Message)
	})
	t.Run("create employee invalid token", func(t *testing.T) {
		fakeAuth := func(c *fiber.Ctx) error {
			authHeader := c.Get("Authorization")
			tokenString := strings.TrimPrefix(authHeader, "Bearer ")
			tokenString = strings.TrimSpace(tokenString)
			if len(strings.Split(tokenString, ".")) != 3 {
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
					"success": false,
					"error":   "token is malformed: token contains an invalid number of segments",
				})
			}
			c.Locals(web.JwtKey, &jwt.Token{})
			return c.Next()
		}
		server := web.NewServer()
		server.GroupApiV1.Use(fakeAuth)
		var svc = new(MockService)
		var controller = NewController(server, svc)
		controller.RegisterRoutes()
		var body = strings.NewReader("{\"name\": \"john doe\", \"role_id\": 1}")
		var request = httptest.NewRequest(fiber.MethodPost, "/api/v1/employees", body)
		request.Header.Set("Authorization", "Bearer this.is.not.a.jwt")
		request.Header.Add("Content-Type", "application/json")
		svc.On("Save", mock.AnythingOfType("CreateRequest")).Return(Response{Id: int64(123)}, nil)
		message := "token is malformed: token contains an invalid number of segments"
		resp, err := server.App.Test(request)
		a.Nil(err)
		a.NotEmpty(resp)
		a.Equal(http.StatusUnauthorized, resp.StatusCode)
		bytesData, err := io.ReadAll(resp.Body)
		a.Nil(err)
		var responseBody common.Response[Response]
		err = json.Unmarshal(bytesData, &responseBody)
		a.Nil(err)
		a.Equal(message, responseBody.Message)
	})
	t.Run("create employee expired token", func(t *testing.T) {
		var claims = &web.IdmClaims{
			RealmAccess: web.RealmAccessClaims{
				Roles: []string{web.IdmAdmin},
			},
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(-1 * time.Hour)),
			},
		}
		var auth = func(c *fiber.Ctx) error {
			if claims.ExpiresAt != nil && claims.ExpiresAt.Time.Before(time.Now()) {
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
					"success": false,
					"error":   "token has invalid claims: token is expired",
				})
			}
			c.Locals(web.JwtKey, &jwt.Token{Claims: claims})
			return c.Next()
		}
		server := web.NewServer()
		server.GroupApiV1.Use(auth)
		var svc = new(MockService)
		var controller = NewController(server, svc)
		controller.RegisterRoutes()
		var body = strings.NewReader("{\"name\": \"john doe\", \"role_id\": 1}")
		var request = httptest.NewRequest(fiber.MethodPost, "/api/v1/employees", body)
		request.Header.Add("Content-Type", "application/json")
		svc.On("Save", mock.AnythingOfType("CreateRequest")).Return(Response{Id: int64(123)}, nil)
		message := "token has invalid claims: token is expired"
		resp, err := server.App.Test(request)
		a.Nil(err)
		a.NotEmpty(resp)
		a.Equal(http.StatusUnauthorized, resp.StatusCode)
		bytesData, err := io.ReadAll(resp.Body)
		a.Nil(err)
		var responseBody common.Response[Response]
		err = json.Unmarshal(bytesData, &responseBody)
		a.Nil(err)
		a.Equal(message, responseBody.Message)
	})
	t.Run("create employee without role admin", func(t *testing.T) {
		var claims = &web.IdmClaims{
			RealmAccess: web.RealmAccessClaims{Roles: []string{web.IdmUser}},
		}
		var auth = func(c *fiber.Ctx) error {
			c.Locals(web.JwtKey, &jwt.Token{Claims: claims})
			return c.Next()
		}
		server := web.NewServer()
		server.GroupApiV1.Use(auth)
		var svc = new(MockService)
		var controller = NewController(server, svc)
		controller.RegisterRoutes()
		var body = strings.NewReader("{\"name\": \"john doe\", \"role_id\": 1}")
		var request = httptest.NewRequest(fiber.MethodPost, "/api/v1/employees", body)
		request.Header.Add("Content-Type", "application/json")
		svc.On("Save", mock.AnythingOfType("CreateRequest")).Return(Response{Id: int64(123)}, nil)
		resp, err := server.App.Test(request)
		message := "Permission denied"
		a.Nil(err)
		a.NotEmpty(resp)
		a.Equal(http.StatusForbidden, resp.StatusCode)
		bytesData, err := io.ReadAll(resp.Body)
		a.Nil(err)
		var responseBody common.Response[Response]
		err = json.Unmarshal(bytesData, &responseBody)
		a.Nil(err)
		a.Equal(message, responseBody.Message)
	})
}

func TestFindEmployeeById(t *testing.T) {
	var a = assert.New(t)
	t.Run("find employee by id", func(t *testing.T) {
		var claims = &web.IdmClaims{
			RealmAccess: web.RealmAccessClaims{Roles: []string{web.IdmAdmin}},
		}
		var auth = func(c *fiber.Ctx) error {
			c.Locals(web.JwtKey, &jwt.Token{Claims: claims})
			return c.Next()
		}
		server := web.NewServer()
		server.GroupApiV1.Use(auth)
		var svc = new(MockService)
		var controller = NewController(server, svc)
		controller.RegisterRoutes()
		var request = httptest.NewRequest(fiber.MethodGet, "/api/v1/employees/123", nil)
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
	t.Run("find employee - incorrect id", func(t *testing.T) {
		var claims = &web.IdmClaims{
			RealmAccess: web.RealmAccessClaims{Roles: []string{web.IdmAdmin}},
		}
		var auth = func(c *fiber.Ctx) error {
			c.Locals(web.JwtKey, &jwt.Token{Claims: claims})
			return c.Next()
		}
		server := web.NewServer()
		server.GroupApiV1.Use(auth)
		var svc = new(MockService)
		var controller = NewController(server, svc)
		controller.RegisterRoutes()
		var request = httptest.NewRequest(fiber.MethodGet, "/api/v1/employees/ffff", nil)
		request.Header.Add("Content-Type", "application/json")
		resp, err := server.App.Test(request)
		a.Nil(err)
		a.NotEmpty(resp)
		a.Equal(http.StatusBadRequest, resp.StatusCode)
	})
	t.Run("find employee - validation error", func(t *testing.T) {
		var claims = &web.IdmClaims{
			RealmAccess: web.RealmAccessClaims{Roles: []string{web.IdmAdmin}},
		}
		var auth = func(c *fiber.Ctx) error {
			c.Locals(web.JwtKey, &jwt.Token{Claims: claims})
			return c.Next()
		}
		server := web.NewServer()
		server.GroupApiV1.Use(auth)
		var svc = new(MockService)
		var controller = NewController(server, svc)
		controller.RegisterRoutes()
		var request = httptest.NewRequest(fiber.MethodGet, "/api/v1/employees/0", nil)
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
	t.Run("find employee - not found error", func(t *testing.T) {
		var claims = &web.IdmClaims{
			RealmAccess: web.RealmAccessClaims{Roles: []string{web.IdmAdmin}},
		}
		var auth = func(c *fiber.Ctx) error {
			c.Locals(web.JwtKey, &jwt.Token{Claims: claims})
			return c.Next()
		}
		server := web.NewServer()
		server.GroupApiV1.Use(auth)
		var svc = new(MockService)
		var controller = NewController(server, svc)
		controller.RegisterRoutes()
		var request = httptest.NewRequest(fiber.MethodGet, "/api/v1/employees/123", nil)
		request.Header.Add("Content-Type", "application/json")
		message := "error finding employee with id 123"
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

func TestFindAllEmployees(t *testing.T) {
	var a = assert.New(t)
	t.Run("find all employees", func(t *testing.T) {
		var claims = &web.IdmClaims{
			RealmAccess: web.RealmAccessClaims{Roles: []string{web.IdmAdmin}},
		}
		var auth = func(c *fiber.Ctx) error {
			c.Locals(web.JwtKey, &jwt.Token{Claims: claims})
			return c.Next()
		}
		server := web.NewServer()
		server.GroupApiV1.Use(auth)
		var svc = new(MockService)
		var controller = NewController(server, svc)
		controller.RegisterRoutes()
		var request = httptest.NewRequest(fiber.MethodGet, "/api/v1/employees", nil)
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
		var claims = &web.IdmClaims{
			RealmAccess: web.RealmAccessClaims{Roles: []string{web.IdmAdmin}},
		}
		var auth = func(c *fiber.Ctx) error {
			c.Locals(web.JwtKey, &jwt.Token{Claims: claims})
			return c.Next()
		}
		server := web.NewServer()
		server.GroupApiV1.Use(auth)
		var svc = new(MockService)
		var controller = NewController(server, svc)
		controller.RegisterRoutes()
		var request = httptest.NewRequest(fiber.MethodGet, "/api/v1/employees", nil)
		request.Header.Add("Content-Type", "application/json")
		message := "error finding all employees"
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

func TestFindEmployeesByIds(t *testing.T) {
	var a = assert.New(t)
	t.Run("find employees by ids", func(t *testing.T) {
		var claims = &web.IdmClaims{
			RealmAccess: web.RealmAccessClaims{Roles: []string{web.IdmAdmin}},
		}
		var auth = func(c *fiber.Ctx) error {
			c.Locals(web.JwtKey, &jwt.Token{Claims: claims})
			return c.Next()
		}
		server := web.NewServer()
		server.GroupApiV1.Use(auth)
		var svc = new(MockService)
		var controller = NewController(server, svc)
		controller.RegisterRoutes()
		var request = httptest.NewRequest(fiber.MethodGet, "/api/v1/employees/find?ids=123,124,125", nil)
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
	t.Run("find employees by ids - error parsing", func(t *testing.T) {
		var claims = &web.IdmClaims{
			RealmAccess: web.RealmAccessClaims{Roles: []string{web.IdmAdmin}},
		}
		var auth = func(c *fiber.Ctx) error {
			c.Locals(web.JwtKey, &jwt.Token{Claims: claims})
			return c.Next()
		}
		server := web.NewServer()
		server.GroupApiV1.Use(auth)
		var svc = new(MockService)
		var controller = NewController(server, svc)
		controller.RegisterRoutes()
		var request = httptest.NewRequest(fiber.MethodGet, "/api/v1/employees/find?ids=fff,124,125", nil)
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
	t.Run("find employees by ids - validation error", func(t *testing.T) {
		var claims = &web.IdmClaims{
			RealmAccess: web.RealmAccessClaims{Roles: []string{web.IdmAdmin}},
		}
		var auth = func(c *fiber.Ctx) error {
			c.Locals(web.JwtKey, &jwt.Token{Claims: claims})
			return c.Next()
		}
		server := web.NewServer()
		server.GroupApiV1.Use(auth)
		var svc = new(MockService)
		var controller = NewController(server, svc)
		controller.RegisterRoutes()
		var request = httptest.NewRequest(fiber.MethodGet, "/api/v1/employees/find?ids=123,124,125", nil)
		request.Header.Add("Content-Type", "application/json")
		message := "employee not found"
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
	t.Run("find employees by ids - not found error", func(t *testing.T) {
		var claims = &web.IdmClaims{
			RealmAccess: web.RealmAccessClaims{Roles: []string{web.IdmAdmin}},
		}
		var auth = func(c *fiber.Ctx) error {
			c.Locals(web.JwtKey, &jwt.Token{Claims: claims})
			return c.Next()
		}
		server := web.NewServer()
		server.GroupApiV1.Use(auth)
		var svc = new(MockService)
		var controller = NewController(server, svc)
		controller.RegisterRoutes()
		var request = httptest.NewRequest(fiber.MethodGet, "/api/v1/employees/find?ids=123,124,125", nil)
		request.Header.Add("Content-Type", "application/json")
		message := "error finding employees by ids"
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

func TestDeleteEmployeeById(t *testing.T) {
	var a = assert.New(t)
	t.Run("find employee by id", func(t *testing.T) {
		var claims = &web.IdmClaims{
			RealmAccess: web.RealmAccessClaims{Roles: []string{web.IdmAdmin}},
		}
		var auth = func(c *fiber.Ctx) error {
			c.Locals(web.JwtKey, &jwt.Token{Claims: claims})
			return c.Next()
		}
		server := web.NewServer()
		server.GroupApiV1.Use(auth)
		var svc = new(MockService)
		var controller = NewController(server, svc)
		controller.RegisterRoutes()
		var request = httptest.NewRequest(fiber.MethodDelete, "/api/v1/employees/123", nil)
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
	t.Run("find employee - incorrect id", func(t *testing.T) {
		var claims = &web.IdmClaims{
			RealmAccess: web.RealmAccessClaims{Roles: []string{web.IdmAdmin}},
		}
		var auth = func(c *fiber.Ctx) error {
			c.Locals(web.JwtKey, &jwt.Token{Claims: claims})
			return c.Next()
		}
		server := web.NewServer()
		server.GroupApiV1.Use(auth)
		var svc = new(MockService)
		var controller = NewController(server, svc)
		controller.RegisterRoutes()
		var request = httptest.NewRequest(fiber.MethodDelete, "/api/v1/employees/ffff", nil)
		request.Header.Add("Content-Type", "application/json")
		resp, err := server.App.Test(request)
		a.Nil(err)
		a.NotEmpty(resp)
		a.Equal(http.StatusBadRequest, resp.StatusCode)
	})
	t.Run("find employee - validation error", func(t *testing.T) {
		var claims = &web.IdmClaims{
			RealmAccess: web.RealmAccessClaims{Roles: []string{web.IdmAdmin}},
		}
		var auth = func(c *fiber.Ctx) error {
			c.Locals(web.JwtKey, &jwt.Token{Claims: claims})
			return c.Next()
		}
		server := web.NewServer()
		server.GroupApiV1.Use(auth)
		var svc = new(MockService)
		var controller = NewController(server, svc)
		controller.RegisterRoutes()
		var request = httptest.NewRequest(fiber.MethodDelete, "/api/v1/employees/0", nil)
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
	t.Run("find employee - not found error", func(t *testing.T) {
		var claims = &web.IdmClaims{
			RealmAccess: web.RealmAccessClaims{Roles: []string{web.IdmAdmin}},
		}
		var auth = func(c *fiber.Ctx) error {
			c.Locals(web.JwtKey, &jwt.Token{Claims: claims})
			return c.Next()
		}
		server := web.NewServer()
		server.GroupApiV1.Use(auth)
		var svc = new(MockService)
		var controller = NewController(server, svc)
		controller.RegisterRoutes()
		var request = httptest.NewRequest(fiber.MethodDelete, "/api/v1/employees/123", nil)
		request.Header.Add("Content-Type", "application/json")
		message := "error finding employee with id 123"
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

func TestDeleteEmployeesByIds(t *testing.T) {
	var a = assert.New(t)
	t.Run("delete employees by ids", func(t *testing.T) {
		var claims = &web.IdmClaims{
			RealmAccess: web.RealmAccessClaims{Roles: []string{web.IdmAdmin}},
		}
		var auth = func(c *fiber.Ctx) error {
			c.Locals(web.JwtKey, &jwt.Token{Claims: claims})
			return c.Next()
		}
		server := web.NewServer()
		server.GroupApiV1.Use(auth)
		var svc = new(MockService)
		var controller = NewController(server, svc)
		controller.RegisterRoutes()
		var request = httptest.NewRequest(fiber.MethodDelete, "/api/v1/employees/delete?ids=123,124,125", nil)
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
	t.Run("delete employees by ids - error parsing", func(t *testing.T) {
		var claims = &web.IdmClaims{
			RealmAccess: web.RealmAccessClaims{Roles: []string{web.IdmAdmin}},
		}
		var auth = func(c *fiber.Ctx) error {
			c.Locals(web.JwtKey, &jwt.Token{Claims: claims})
			return c.Next()
		}
		server := web.NewServer()
		server.GroupApiV1.Use(auth)
		var svc = new(MockService)
		var controller = NewController(server, svc)
		controller.RegisterRoutes()
		var request = httptest.NewRequest(fiber.MethodDelete, "/api/v1/employees/delete?ids=fff,124,125", nil)
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
	t.Run("delete employees by ids - validation error", func(t *testing.T) {
		var claims = &web.IdmClaims{
			RealmAccess: web.RealmAccessClaims{Roles: []string{web.IdmAdmin}},
		}
		var auth = func(c *fiber.Ctx) error {
			c.Locals(web.JwtKey, &jwt.Token{Claims: claims})
			return c.Next()
		}
		server := web.NewServer()
		server.GroupApiV1.Use(auth)
		var svc = new(MockService)
		var controller = NewController(server, svc)
		controller.RegisterRoutes()
		var request = httptest.NewRequest(fiber.MethodDelete, "/api/v1/employees/delete?ids=123,124,125", nil)
		request.Header.Add("Content-Type", "application/json")
		message := "employee not found"
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
	t.Run("delete employees by ids - not found error", func(t *testing.T) {
		var claims = &web.IdmClaims{
			RealmAccess: web.RealmAccessClaims{Roles: []string{web.IdmAdmin}},
		}
		var auth = func(c *fiber.Ctx) error {
			c.Locals(web.JwtKey, &jwt.Token{Claims: claims})
			return c.Next()
		}
		server := web.NewServer()
		server.GroupApiV1.Use(auth)
		var svc = new(MockService)
		var controller = NewController(server, svc)
		controller.RegisterRoutes()
		var request = httptest.NewRequest(fiber.MethodDelete, "/api/v1/employees/delete?ids=123,124,125", nil)
		request.Header.Add("Content-Type", "application/json")
		message := "error finding employees by ids"
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
