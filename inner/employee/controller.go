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
	request := PageRequest{
		PageSize:   pageSize,
		PageNumber: pageNumber,
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
