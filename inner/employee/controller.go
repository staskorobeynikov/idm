package employee

import (
	"errors"
	"github.com/gofiber/fiber/v3"
	"go.uber.org/zap"
	"idm/inner/common"
	"idm/inner/web"
	"strconv"
	"strings"
)

type Controller struct {
	server          *web.Server
	employeeService Svc
	logger          *common.Logger
}

type Svc interface {
	Save(request CreateRequest) (Response, error)
	FindById(request IdRequest) (Response, error)
	FindAll() ([]Response, error)
	FindByIds(request IdsRequest) ([]Response, error)
	DeleteById(request IdRequest) error
	DeleteByIds(request IdsRequest) error
}

func NewController(
	server *web.Server,
	employeeService Svc,
	logger *common.Logger,
) *Controller {
	return &Controller{
		server:          server,
		employeeService: employeeService,
		logger:          logger,
	}
}

func (c *Controller) RegisterRoutes() {
	c.server.GroupApiV1.Post("/employees", c.CreateEmployee)
	c.server.GroupApiV1.Get("/employees/find", c.FindByIds)
	c.server.GroupApiV1.Get("/employees/:id", c.FindById)
	c.server.GroupApiV1.Get("/employees", c.FindAll)
	c.server.GroupApiV1.Delete("/employees/delete", c.DeleteByIds)
	c.server.GroupApiV1.Delete("/employees/:id", c.DeleteById)
}

func (c *Controller) CreateEmployee(ctx fiber.Ctx) error {
	var request CreateRequest
	if err := ctx.Bind().Body(&request); err != nil {
		c.logger.Error("create employee", zap.Error(err))
		return common.ErrResponse(ctx, fiber.StatusBadRequest, err.Error())
	}
	c.logger.Info("create employee: received request", zap.Any("request", request))
	var response, err = c.employeeService.Save(request)
	if err != nil {
		switch {
		case errors.As(err, &common.RequestValidationError{}) || errors.As(err, &common.AlreadyExistsError{}):
			c.logger.Error("create employee", zap.Error(err))
			return common.ErrResponse(ctx, fiber.StatusBadRequest, err.Error())
		default:
			c.logger.Error("create employee", zap.Error(err))
			return common.ErrResponse(ctx, fiber.StatusInternalServerError, err.Error())
		}
	}
	return common.OkResponse(ctx, response.Id)
}

func (c *Controller) FindById(ctx fiber.Ctx) error {
	var param = ctx.Params("id")
	id, err := strconv.Atoi(param)
	if err != nil {
		c.logger.Error("find by id employee", zap.Error(err))
		return common.ErrResponse(ctx, fiber.StatusBadRequest, err.Error())
	}
	request := IdRequest{Id: int64(id)}
	c.logger.Info("find by id employee: received request", zap.Any("request", request))
	response, err := c.employeeService.FindById(request)
	if err != nil {
		switch {
		case errors.As(err, &common.RequestValidationError{}):
			c.logger.Error("find by id employee", zap.Error(err))
			return common.ErrResponse(ctx, fiber.StatusBadRequest, err.Error())
		case errors.As(err, &common.NotFoundError{}):
			c.logger.Error("find by id employee", zap.Error(err))
			return common.ErrResponse(ctx, fiber.StatusOK, err.Error())
		default:
			c.logger.Error("find by id employee", zap.Error(err))
			return common.ErrResponse(ctx, fiber.StatusInternalServerError, err.Error())
		}
	}
	return common.OkResponse(ctx, response)
}

func (c *Controller) FindAll(ctx fiber.Ctx) error {
	response, err := c.employeeService.FindAll()
	if err != nil {
		c.logger.Error("find all employees", zap.Error(err))
		return common.ErrResponse(ctx, fiber.StatusOK, err.Error())
	}
	return common.OkResponse(ctx, response)
}

func (c *Controller) FindByIds(ctx fiber.Ctx) error {
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
	c.logger.Info("find by ids employees: received request", zap.Any("request", request))
	var response, err = c.employeeService.FindByIds(request)
	if err != nil {
		switch {
		case errors.As(err, &common.RequestValidationError{}):
			c.logger.Error("find by ids employees", zap.Error(err))
			return common.ErrResponse(ctx, fiber.StatusBadRequest, err.Error())
		case errors.As(err, &common.NotFoundError{}):
			c.logger.Error("find by ids employees", zap.Error(err))
			return common.ErrResponse(ctx, fiber.StatusOK, err.Error())
		default:
			c.logger.Error("find by ids employees", zap.Error(err))
			return common.ErrResponse(ctx, fiber.StatusInternalServerError, err.Error())
		}
	}
	return common.OkResponse(ctx, response)
}

func (c *Controller) DeleteById(ctx fiber.Ctx) error {
	var param = ctx.Params("id")
	id, err := strconv.Atoi(param)
	if err != nil {
		return common.ErrResponse(ctx, fiber.StatusBadRequest, err.Error())
	}
	request := IdRequest{Id: int64(id)}
	c.logger.Info("delete by id employee: received request", zap.Any("request", request))
	err = c.employeeService.DeleteById(request)
	if err != nil {
		switch {
		case errors.As(err, &common.RequestValidationError{}):
			c.logger.Error("delete by id employee", zap.Error(err))
			return common.ErrResponse(ctx, fiber.StatusBadRequest, err.Error())
		case errors.As(err, &common.NotFoundError{}):
			c.logger.Error("delete by id employee", zap.Error(err))
			return common.ErrResponse(ctx, fiber.StatusOK, err.Error())
		default:
			c.logger.Error("delete by id employee", zap.Error(err))
			return common.ErrResponse(ctx, fiber.StatusInternalServerError, err.Error())
		}
	}
	return common.OkResponse(ctx, Response{Id: int64(id)})
}

func (c *Controller) DeleteByIds(ctx fiber.Ctx) error {
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
	c.logger.Info("delete by ids employees: received request", zap.Any("request", request))
	err := c.employeeService.DeleteByIds(request)
	if err != nil {
		switch {
		case errors.As(err, &common.RequestValidationError{}):
			c.logger.Error("delete by ids employees", zap.Error(err))
			return common.ErrResponse(ctx, fiber.StatusBadRequest, err.Error())
		case errors.As(err, &common.NotFoundError{}):
			c.logger.Error("delete by ids employees", zap.Error(err))
			return common.ErrResponse(ctx, fiber.StatusOK, err.Error())
		default:
			c.logger.Error("delete by ids employees", zap.Error(err))
			return common.ErrResponse(ctx, fiber.StatusInternalServerError, err.Error())
		}
	}
	var responses []Response
	for _, id := range ids {
		responses = append(responses, Response{Id: int64(id)})
	}
	return common.OkResponse(ctx, responses)
}
