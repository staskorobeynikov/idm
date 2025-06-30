package role

import (
	"fmt"
	"idm/inner/common"
)

type Service struct {
	repo      Repo
	validator Validator
}

type Repo interface {
	Save(entity Entity) (int64, error)
	FindById(id int64) (entity Entity, err error)
	FindAll() ([]Entity, error)
	FindByIds(ids []int64) ([]Entity, error)
	DeleteById(id int64) error
	DeleteByIds(ids []int64) error
}

type Validator interface {
	Validate(request any) error
}

func NewService(repo Repo, validator Validator) *Service {
	return &Service{
		repo:      repo,
		validator: validator,
	}
}

func (s *Service) Save(request CreateRequest) (Response, error) {
	err := s.validator.Validate(request)
	if err != nil {
		return Response{}, common.RequestValidationError{
			Message: err.Error(),
		}
	}
	var id int64
	id, err = s.repo.Save(request.ToEntity())
	if err != nil {
		return Response{}, fmt.Errorf("error saving role with: %w", err)
	}
	return Response{
		Id: id,
	}, nil
}

func (s *Service) FindById(request IdRequest) (Response, error) {
	var err = s.validator.Validate(request)
	if err != nil {
		return Response{}, common.RequestValidationError{Message: err.Error()}
	}
	entity, err := s.repo.FindById(request.Id)
	if err != nil {
		return Response{}, common.NotFoundError{Message: fmt.Sprintf("error finding role with id %d: %v", request.Id, err)}
	}
	return entity.toResponse(), nil
}

func (s *Service) FindAll() ([]Response, error) {
	var employees, err = s.repo.FindAll()
	if err != nil {
		return nil, common.NotFoundError{Message: fmt.Sprintf("error finding all roles: %v", err)}
	}
	var response []Response
	for _, employee := range employees {
		response = append(response, employee.toResponse())
	}
	return response, nil
}

func (s *Service) FindByIds(request IdsRequest) ([]Response, error) {
	if err := s.validator.Validate(request); err != nil {
		return []Response{}, common.RequestValidationError{Message: err.Error()}
	}
	var employees, err = s.repo.FindByIds(request.Ids)
	if err != nil {
		return nil, common.NotFoundError{Message: fmt.Sprintf("error finding roles by ids: %v", err)}
	}
	var response []Response
	for _, employee := range employees {
		response = append(response, employee.toResponse())
	}
	return response, nil
}

func (s *Service) DeleteById(request IdRequest) error {
	var err = s.validator.Validate(request)
	if err != nil {
		return common.RequestValidationError{Message: err.Error()}
	}
	err = s.repo.DeleteById(request.Id)
	if err != nil {
		return common.NotFoundError{Message: fmt.Sprintf("error deleting role with id %d: %v", request.Id, err)}
	}
	return nil
}

func (s *Service) DeleteByIds(request IdsRequest) error {
	if err := s.validator.Validate(request); err != nil {
		return common.RequestValidationError{Message: err.Error()}
	}
	var err = s.repo.DeleteByIds(request.Ids)
	if err != nil {
		return common.NotFoundError{Message: fmt.Sprintf("error deleting roles with ids %d: %v", request.Ids, err)}
	}
	return nil
}
