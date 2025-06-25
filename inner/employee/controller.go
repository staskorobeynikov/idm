package employee

import (
	"errors"
	"github.com/gofiber/fiber/v3"
	"idm/inner/common"
	"idm/inner/web"
	"strconv"
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
	DeleteById(request IdRequest) error
	DeleteByIds(request IdsRequest) error
}

func NewController(server *web.Server, employeeService Svc) *Controller {
	return &Controller{
		server:          server,
		employeeService: employeeService,
	}
}

func (c *Controller) RegisterRoutes() {
	c.server.GroupApiV1.Post("/employees", c.CreateEmployee)
	c.server.GroupApiV1.Get("/employees/:id", c.FindById)
	c.server.GroupApiV1.Get("/employees", c.FindAll)
	c.server.GroupApiV1.Get("/employees/find", c.FindByIds)
	c.server.GroupApiV1.Delete("/employees/:id", c.DeleteById)
	c.server.GroupApiV1.Delete("/employees/:ids", c.DeleteByIds)
}

func (c *Controller) CreateEmployee(ctx fiber.Ctx) error {
	var request CreateRequest
	if err := ctx.Bind().Body(&request); err != nil {
		return common.ErrResponse(ctx, fiber.StatusBadRequest, err.Error())
	}
	var response, err = c.employeeService.Save(request)
	if err != nil {
		switch {
		case errors.As(err, &common.RequestValidationError{}) || errors.As(err, &common.AlreadyExistsError{}):
			return common.ErrResponse(ctx, fiber.StatusBadRequest, err.Error())
		default:
			return common.ErrResponse(ctx, fiber.StatusInternalServerError, err.Error())
		}
	}
	return common.OkResponse(ctx, response.Id)
}

func (c *Controller) FindById(ctx fiber.Ctx) error {
	var param = ctx.Params("id")
	id, err := strconv.Atoi(param)
	if err != nil {
		return common.ErrResponse(ctx, fiber.StatusBadRequest, err.Error())
	}
	request := IdRequest{Id: int64(id)}
	response, err := c.employeeService.FindById(request)
	if err != nil {
		switch {
		case errors.As(err, &common.RequestValidationError{}):
			return common.ErrResponse(ctx, fiber.StatusBadRequest, err.Error())
		case errors.As(err, &common.NotFoundError{}):
			return common.ErrResponse(ctx, fiber.StatusOK, err.Error())
		default:
			return common.ErrResponse(ctx, fiber.StatusInternalServerError, err.Error())
		}
	}
	return common.OkResponse(ctx, response)
}

func (c *Controller) FindAll(ctx fiber.Ctx) error {
	response, err := c.employeeService.FindAll()
	if err != nil {
		return common.ErrResponse(ctx, fiber.StatusOK, err.Error())
	}
	return common.OkResponse(ctx, response)
}

func (c *Controller) FindByIds(ctx fiber.Ctx) error {
	var request IdsRequest
	if err := ctx.Bind().Body(&request); err != nil {
		return common.ErrResponse(ctx, fiber.StatusBadRequest, err.Error())
	}
	var response, err = c.employeeService.FindByIds(request)
	if err != nil {
		switch {
		case errors.As(err, &common.RequestValidationError{}):
			return common.ErrResponse(ctx, fiber.StatusBadRequest, err.Error())
		case errors.As(err, &common.NotFoundError{}):
			return common.ErrResponse(ctx, fiber.StatusOK, err.Error())
		default:
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
	err = c.employeeService.DeleteById(request)
	if err != nil {
		switch {
		case errors.As(err, &common.RequestValidationError{}):
			return common.ErrResponse(ctx, fiber.StatusBadRequest, err.Error())
		case errors.As(err, &common.NotFoundError{}):
			return common.ErrResponse(ctx, fiber.StatusOK, err.Error())
		default:
			return common.ErrResponse(ctx, fiber.StatusInternalServerError, err.Error())
		}
	}
	return common.OkResponse(ctx, Response{Id: int64(id)})
}

func (c *Controller) DeleteByIds(ctx fiber.Ctx) error {
	var request IdsRequest
	if err := ctx.Bind().Body(&request); err != nil {
		return common.ErrResponse(ctx, fiber.StatusBadRequest, err.Error())
	}
	err := c.employeeService.DeleteByIds(request)
	if err != nil {
		switch {
		case errors.As(err, &common.RequestValidationError{}):
			return common.ErrResponse(ctx, fiber.StatusBadRequest, err.Error())
		case errors.As(err, &common.NotFoundError{}):
			return common.ErrResponse(ctx, fiber.StatusOK, err.Error())
		default:
			return common.ErrResponse(ctx, fiber.StatusInternalServerError, err.Error())
		}
	}
	return common.OkResponse(ctx, Response{Id: int64(0)})
}
