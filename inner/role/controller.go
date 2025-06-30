package role

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
	server      *web.Server
	roleService Svc
	logger      *common.Logger
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
	roleService Svc,
	logger *common.Logger,
) *Controller {
	return &Controller{
		server:      server,
		roleService: roleService,
		logger:      logger,
	}
}

func (c *Controller) RegisterRoutes() {
	c.server.GroupApiV1.Post("/roles", c.CreateRole)
	c.server.GroupApiV1.Get("/roles/find", c.FindByIds)
	c.server.GroupApiV1.Get("/roles/:id", c.FindById)
	c.server.GroupApiV1.Get("/roles", c.FindAll)
	c.server.GroupApiV1.Delete("/roles/delete", c.DeleteByIds)
	c.server.GroupApiV1.Delete("/roles/:id", c.DeleteById)
}

func (c *Controller) CreateRole(ctx fiber.Ctx) error {
	var request CreateRequest
	if err := ctx.Bind().Body(&request); err != nil {
		c.logger.Error("create role", zap.Error(err))
		return common.ErrResponse(ctx, fiber.StatusBadRequest, err.Error())
	}
	c.logger.Info("create role: received request", zap.Any("request", request))
	var response, err = c.roleService.Save(request)
	if err != nil {
		switch {
		case errors.As(err, &common.RequestValidationError{}):
			c.logger.Error("create role", zap.Error(err))
			return common.ErrResponse(ctx, fiber.StatusBadRequest, err.Error())
		default:
			c.logger.Error("create role", zap.Error(err))
			return common.ErrResponse(ctx, fiber.StatusInternalServerError, err.Error())
		}
	}
	return common.OkResponse(ctx, response.Id)
}

func (c *Controller) FindById(ctx fiber.Ctx) error {
	var param = ctx.Params("id")
	id, err := strconv.Atoi(param)
	if err != nil {
		c.logger.Error("find role by id", zap.Error(err))
		return common.ErrResponse(ctx, fiber.StatusBadRequest, err.Error())
	}
	request := IdRequest{Id: int64(id)}
	c.logger.Info("find role by id: received request", zap.Any("request", request))
	response, err := c.roleService.FindById(request)
	if err != nil {
		switch {
		case errors.As(err, &common.RequestValidationError{}):
			c.logger.Error("find role by id", zap.Error(err))
			return common.ErrResponse(ctx, fiber.StatusBadRequest, err.Error())
		case errors.As(err, &common.NotFoundError{}):
			c.logger.Error("find role by id", zap.Error(err))
			return common.ErrResponse(ctx, fiber.StatusOK, err.Error())
		default:
			c.logger.Error("find role by id", zap.Error(err))
			return common.ErrResponse(ctx, fiber.StatusInternalServerError, err.Error())
		}
	}
	return common.OkResponse(ctx, response)
}

func (c *Controller) FindAll(ctx fiber.Ctx) error {
	response, err := c.roleService.FindAll()
	if err != nil {
		c.logger.Error("find all roles", zap.Error(err))
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
			c.logger.Error("find roles by ids", zap.Error(err))
			return common.ErrResponse(ctx, fiber.StatusBadRequest, err.Error())
		}
		ids = append(ids, id)
	}
	var request = IdsRequest{Ids: ids}
	c.logger.Info("find roles by ids: received request", zap.Any("request", request))
	var response, err = c.roleService.FindByIds(request)
	if err != nil {
		switch {
		case errors.As(err, &common.RequestValidationError{}):
			c.logger.Error("find roles by ids", zap.Error(err))
			return common.ErrResponse(ctx, fiber.StatusBadRequest, err.Error())
		case errors.As(err, &common.NotFoundError{}):
			c.logger.Error("find roles by ids", zap.Error(err))
			return common.ErrResponse(ctx, fiber.StatusOK, err.Error())
		default:
			c.logger.Error("find roles by ids", zap.Error(err))
			return common.ErrResponse(ctx, fiber.StatusInternalServerError, err.Error())
		}
	}
	return common.OkResponse(ctx, response)
}

func (c *Controller) DeleteById(ctx fiber.Ctx) error {
	var param = ctx.Params("id")
	id, err := strconv.Atoi(param)
	if err != nil {
		c.logger.Error("delete role by id", zap.Error(err))
		return common.ErrResponse(ctx, fiber.StatusBadRequest, err.Error())
	}
	request := IdRequest{Id: int64(id)}
	c.logger.Info("delete role by id: received request", zap.Any("request", request))
	err = c.roleService.DeleteById(request)
	if err != nil {
		switch {
		case errors.As(err, &common.RequestValidationError{}):
			c.logger.Error("delete role by id", zap.Error(err))
			return common.ErrResponse(ctx, fiber.StatusBadRequest, err.Error())
		case errors.As(err, &common.NotFoundError{}):
			c.logger.Error("delete role by id", zap.Error(err))
			return common.ErrResponse(ctx, fiber.StatusOK, err.Error())
		default:
			c.logger.Error("delete role by id", zap.Error(err))
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
			c.logger.Error("delete roles by ids", zap.Error(err))
			return common.ErrResponse(ctx, fiber.StatusBadRequest, err.Error())
		}
		ids = append(ids, id)
	}
	var request = IdsRequest{Ids: ids}
	c.logger.Info("delete roles by ids: received request", zap.Any("request", request))
	err := c.roleService.DeleteByIds(request)
	if err != nil {
		switch {
		case errors.As(err, &common.RequestValidationError{}):
			c.logger.Error("delete roles by ids", zap.Error(err))
			return common.ErrResponse(ctx, fiber.StatusBadRequest, err.Error())
		case errors.As(err, &common.NotFoundError{}):
			c.logger.Error("delete roles by ids", zap.Error(err))
			return common.ErrResponse(ctx, fiber.StatusOK, err.Error())
		default:
			c.logger.Error("delete roles by ids", zap.Error(err))
			return common.ErrResponse(ctx, fiber.StatusInternalServerError, err.Error())
		}
	}
	var responses []Response
	for _, id := range ids {
		responses = append(responses, Response{Id: int64(id)})
	}
	return common.OkResponse(ctx, responses)
}
