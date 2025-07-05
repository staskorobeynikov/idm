package employee

import (
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v3"
	"go.uber.org/zap"
	"idm/inner/common"
	"idm/inner/middleware"
	"idm/inner/web"
	"strconv"
	"strings"
)

type Controller struct {
	server          *web.Server
	employeeService Svc
}

type Svc interface {
	Save(request CreateRequest) (Response, error)
	FindById(request IdRequest) (Response, error)
	FindAll() ([]Response, error)
	FindByIds(request IdsRequest) ([]Response, error)
	FindWithOffset(request PageRequest) (PageResponse, error)
	DeleteById(request IdRequest) error
	DeleteByIds(request IdsRequest) error
}

func NewController(
	server *web.Server,
	employeeService Svc,
) *Controller {
	return &Controller{
		server:          server,
		employeeService: employeeService,
	}
}

func (c *Controller) RegisterRoutes() {
	c.server.GroupApiV1.Post("/employees", c.CreateEmployee)
	c.server.GroupApiV1.Get("/employees/find", c.FindByIds)
	c.server.GroupApiV1.Get("/employees/page", c.FindWithOffset)
	c.server.GroupApiV1.Get("/employees/:id", c.FindById)
	c.server.GroupApiV1.Get("/employees", c.FindAll)
	c.server.GroupApiV1.Delete("/employees/delete", c.DeleteByIds)
	c.server.GroupApiV1.Delete("/employees/:id", c.DeleteById)
}

// Функция-хендлер, которая будет вызываться при POST запросе по маршруту "/api/v1/employees"
// @Summary create a new employee
// @Description Create a new employee.
// @Tags employee
// @Accept json
// @Produce json
// @Param request body employee.CreateRequest true "create employee request"
// @Success 200 {object} common.Response[int64]
// @Failure 400 {object} common.Response[string]
// @Router /employees [post]
func (c *Controller) CreateEmployee(ctx fiber.Ctx) error {
	requestId := string(ctx.Response().Header.Peek(fiber.HeaderXRequestID))
	logger := middleware.GetLogger(ctx.Context())
	logger.Info(fmt.Sprint("Handling request with ID:", requestId))
	var request CreateRequest
	if err := ctx.Bind().Body(&request); err != nil {
		logger.Error("create employee", zap.Error(err))
		return common.ErrResponse(ctx, fiber.StatusBadRequest, err.Error())
	}
	logger.Info("create employee: received request", zap.Any("request", request))
	var response, err = c.employeeService.Save(request)
	if err != nil {
		switch {
		case errors.As(err, &common.RequestValidationError{}) || errors.As(err, &common.AlreadyExistsError{}):
			logger.Error("create employee", zap.Error(err))
			return common.ErrResponse(ctx, fiber.StatusBadRequest, err.Error())
		default:
			logger.Error("create employee", zap.Error(err))
			return common.ErrResponse(ctx, fiber.StatusInternalServerError, err.Error())
		}
	}
	return common.OkResponse(ctx, response.Id)
}

// Функция-хендлер, которая будет вызываться при GET запросе по маршруту "/api/v1/employees/find"
// @Summary Get employees with dynamic filter(optional) and pagination.
// @Description get employees with dynamic filter(optional) and pagination
// @Tags employee
// @Accept json
// @Produce json
// @Param pageNumber  query int true "Page number (0 is first page)"
// @Param pageSize    query int true "Page size (number of employee on the page)"
// @Param textFilter  query string false "Filter name of employees"
// @Success 200 {object} common.PageResponse[[]employee.Response]
// @Failure 400 {object} common.Response[string]
// @Failure 500 {object} common.Response[string]
// @Router /employees/page [get]
func (c *Controller) FindWithOffset(ctx fiber.Ctx) error {
	requestId := string(ctx.Response().Header.Peek(fiber.HeaderXRequestID))
	logger := middleware.GetLogger(ctx.Context())
	logger.Info(fmt.Sprint("Handling request with ID:", requestId))
	pageSize, err := strconv.Atoi(ctx.Query("pageSize", "100"))
	if err != nil {
		logger.Error("parse page size with error: ", zap.Error(err))
		return common.ErrResponse(ctx, fiber.StatusBadRequest, err.Error())
	}
	pageNumber, err := strconv.Atoi(ctx.Query("pageNumber", "0"))
	if err != nil {
		logger.Error("parse page number with error: ", zap.Error(err))
		return common.ErrResponse(ctx, fiber.StatusBadRequest, err.Error())
	}
	textFilter := ctx.Query("textFilter", "")
	request := PageRequest{
		PageSize:   pageSize,
		PageNumber: pageNumber,
		TextFilter: textFilter,
	}
	logger.Info("find employees with offset", zap.Any("request", request))
	response, err := c.employeeService.FindWithOffset(request)
	if err != nil {
		switch {
		case errors.As(err, &common.RequestValidationError{}):
			logger.Error("find employee with offset", zap.Error(err))
			return common.ErrResponse(ctx, fiber.StatusBadRequest, err.Error())
		case errors.As(err, &common.NotFoundError{}):
			logger.Error("find employee with offset", zap.Error(err))
			return common.ErrResponse(ctx, fiber.StatusOK, err.Error())
		default:
			logger.Error("find employee with offset", zap.Error(err))
			return common.ErrResponse(ctx, fiber.StatusInternalServerError, err.Error())
		}
	}
	return common.OkResponse(ctx, response)
}

// Функция-хендлер, которая будет вызываться при GET запросе по маршруту "/api/v1/employees/:id"
// @Summary Get employee by ID
// @Description returns details of a single employee by their unique ID
// @Tags employee
// @Accept json
// @Produce json
// @Param id path int true "Employee ID"
// @Success 200 {object} common.Response[employee.Response]
// @Failure 400 {object} common.Response[string]
// @Failure 500 {object} common.Response[string]
// @Router /employees/{id} [get]
func (c *Controller) FindById(ctx fiber.Ctx) error {
	requestId := string(ctx.Response().Header.Peek(fiber.HeaderXRequestID))
	logger := middleware.GetLogger(ctx.Context())
	logger.Info(fmt.Sprint("Handling request with ID:", requestId))
	var param = ctx.Params("id")
	id, err := strconv.Atoi(param)
	if err != nil {
		logger.Error("find by id employee", zap.Error(err))
		return common.ErrResponse(ctx, fiber.StatusBadRequest, err.Error())
	}
	request := IdRequest{Id: int64(id)}
	logger.Info("find by id employee: received request", zap.Any("request", request))
	response, err := c.employeeService.FindById(request)
	if err != nil {
		switch {
		case errors.As(err, &common.RequestValidationError{}):
			logger.Error("find by id employee", zap.Error(err))
			return common.ErrResponse(ctx, fiber.StatusBadRequest, err.Error())
		case errors.As(err, &common.NotFoundError{}):
			logger.Error("find by id employee", zap.Error(err))
			return common.ErrResponse(ctx, fiber.StatusOK, err.Error())
		default:
			logger.Error("find by id employee", zap.Error(err))
			return common.ErrResponse(ctx, fiber.StatusInternalServerError, err.Error())
		}
	}
	return common.OkResponse(ctx, response)
}

// Функция-хендлер, которая будет вызываться при GET запросе по маршруту "/api/v1/employees"
// @Summary Get all employees
// @Description returns a list of all employees
// @Tags employee
// @Accept json
// @Produce json
// @Success 200 {object} common.Response[[]employee.Response]
// @Failure 400 {object} common.Response[string]
// @Failure 500 {object} common.Response[string]
// @Router /employees [get]
func (c *Controller) FindAll(ctx fiber.Ctx) error {
	requestId := string(ctx.Response().Header.Peek(fiber.HeaderXRequestID))
	logger := middleware.GetLogger(ctx.Context())
	logger.Info(fmt.Sprint("Handling request with ID:", requestId))
	response, err := c.employeeService.FindAll()
	if err != nil {
		logger.Error("find all employees", zap.Error(err))
		return common.ErrResponse(ctx, fiber.StatusOK, err.Error())
	}
	return common.OkResponse(ctx, response)
}

// Функция-хендлер, которая будет вызываться при GET запросе по маршруту "/api/v1/employees/find?ids=1,2,3"
// @Summary Get employees by multiple IDs
// @Description Returns a list of employees matching the provided IDs
// @Tags employee
// @Accept json
// @Produce json
// @Param ids query []int true "Comma-separated list of employee IDs (e.g., 1,2,3)"
// @Success 200 {object} common.Response[[]employee.Response]
// @Failure 400 {object} common.Response[string]
// @Failure 500 {object} common.Response[string]
// @Router /employees/find [get]
func (c *Controller) FindByIds(ctx fiber.Ctx) error {
	requestId := string(ctx.Response().Header.Peek(fiber.HeaderXRequestID))
	logger := middleware.GetLogger(ctx.Context())
	logger.Info(fmt.Sprint("Handling request with ID:", requestId))
	idsParam := ctx.Query("ids")
	stringIds := strings.Split(idsParam, ",")
	var ids []int64
	for _, id := range stringIds {
		id, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			return common.ErrResponse(ctx, fiber.StatusBadRequest, err.Error())
		}
		ids = append(ids, id)
	}
	var request = IdsRequest{Ids: ids}
	logger.Info("find by ids employees: received request", zap.Any("request", request))
	var response, err = c.employeeService.FindByIds(request)
	if err != nil {
		switch {
		case errors.As(err, &common.RequestValidationError{}):
			logger.Error("find by ids employees", zap.Error(err))
			return common.ErrResponse(ctx, fiber.StatusBadRequest, err.Error())
		case errors.As(err, &common.NotFoundError{}):
			logger.Error("find by ids employees", zap.Error(err))
			return common.ErrResponse(ctx, fiber.StatusOK, err.Error())
		default:
			logger.Error("find by ids employees", zap.Error(err))
			return common.ErrResponse(ctx, fiber.StatusInternalServerError, err.Error())
		}
	}
	return common.OkResponse(ctx, response)
}

// Функция-хендлер, которая будет вызываться при DELETE запросе по маршруту "/api/v1/employees/:id"
// @Summary Delete employee by ID
// @Description Deletes a single employee by their unique ID
// @Tags employee
// @Accept json
// @Produce json
// @Param id path int true "Employee ID"
// @Success 200 {object} common.Response[int64]
// @Failure 400 {object} common.Response[string]
// @Failure 500 {object} common.Response[string]
// @Router /employees/{id} [delete]
func (c *Controller) DeleteById(ctx fiber.Ctx) error {
	requestId := string(ctx.Response().Header.Peek(fiber.HeaderXRequestID))
	logger := middleware.GetLogger(ctx.Context())
	logger.Info(fmt.Sprint("Handling request with ID:", requestId))
	var param = ctx.Params("id")
	id, err := strconv.Atoi(param)
	if err != nil {
		return common.ErrResponse(ctx, fiber.StatusBadRequest, err.Error())
	}
	request := IdRequest{Id: int64(id)}
	logger.Info("delete by id employee: received request", zap.Any("request", request))
	err = c.employeeService.DeleteById(request)
	if err != nil {
		switch {
		case errors.As(err, &common.RequestValidationError{}):
			logger.Error("delete by id employee", zap.Error(err))
			return common.ErrResponse(ctx, fiber.StatusBadRequest, err.Error())
		case errors.As(err, &common.NotFoundError{}):
			logger.Error("delete by id employee", zap.Error(err))
			return common.ErrResponse(ctx, fiber.StatusOK, err.Error())
		default:
			logger.Error("delete by id employee", zap.Error(err))
			return common.ErrResponse(ctx, fiber.StatusInternalServerError, err.Error())
		}
	}
	return common.OkResponse(ctx, Response{Id: int64(id)})
}

// Функция-хендлер, которая будет вызываться при Delete запросе по маршруту "/api/v1/employees/delete?ids=1,2,3"
// @Summary Delete multiple employees by IDs
// @Description Deletes multiple employees matching the provided IDs
// @Tags employee
// @Accept json
// @Produce json
// @Param ids query []int true "Comma-separated list of employee IDs to delete (e.g., 1,2,3)"
// @Success 200 {object} common.Response[[]int64]
// @Failure 400 {object} common.Response[string]
// @Failure 500 {object} common.Response[string]
// @Router /employees/delete [delete]
func (c *Controller) DeleteByIds(ctx fiber.Ctx) error {
	requestId := string(ctx.Response().Header.Peek(fiber.HeaderXRequestID))
	logger := middleware.GetLogger(ctx.Context())
	logger.Info(fmt.Sprint("Handling request with ID:", requestId))
	idsParam := ctx.Query("ids")
	stringIds := strings.Split(idsParam, ",")
	var ids []int64
	for _, id := range stringIds {
		id, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			return common.ErrResponse(ctx, fiber.StatusBadRequest, err.Error())
		}
		ids = append(ids, id)
	}
	var request = IdsRequest{Ids: ids}
	logger.Info("delete by ids employees: received request", zap.Any("request", request))
	err := c.employeeService.DeleteByIds(request)
	if err != nil {
		switch {
		case errors.As(err, &common.RequestValidationError{}):
			logger.Error("delete by ids employees", zap.Error(err))
			return common.ErrResponse(ctx, fiber.StatusBadRequest, err.Error())
		case errors.As(err, &common.NotFoundError{}):
			logger.Error("delete by ids employees", zap.Error(err))
			return common.ErrResponse(ctx, fiber.StatusOK, err.Error())
		default:
			logger.Error("delete by ids employees", zap.Error(err))
			return common.ErrResponse(ctx, fiber.StatusInternalServerError, err.Error())
		}
	}
	var responses []Response
	for _, id := range ids {
		responses = append(responses, Response{Id: int64(id)})
	}
	return common.OkResponse(ctx, responses)
}
