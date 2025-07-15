package employee

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	"idm/inner/common"
)

type Service struct {
	repo      Repo
	validator Validator
}

type Repo interface {
	BeginTransaction() (*sqlx.Tx, error)
	Save(tx *sqlx.Tx, e Entity) (int64, error)
	FindById(id int64) (Entity, error)
	FindByName(tx *sqlx.Tx, name string) (bool, error)
	FindAll() ([]Entity, error)
	FindByIds(ids []int64) ([]Entity, error)
	FindWithOffset(offset int, limit int, filter string) ([]Entity, error)
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

func (s *Service) Save(ctx context.Context, request CreateRequest) (Response, error) {
	err := s.validator.Validate(request)
	if err != nil {
		return Response{}, common.RequestValidationError{
			Message: err.Error(),
		}
	}
	tx, err := s.repo.BeginTransaction()
	if err != nil {
		return Response{}, fmt.Errorf("error creating transaction: %w", err)
	}
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("creating employee panic: %v", r)
			errTx := tx.Rollback()
			if errTx != nil {
				err = fmt.Errorf("creating employee: rolling back transaction errors: %w, %w", err, errTx)
			}
		} else if err != nil {
			errTx := tx.Rollback()
			if errTx != nil {
				err = fmt.Errorf("creating employee: rolling back transaction errors: %w, %w", err, errTx)
			}
		} else {
			errTx := tx.Commit()
			if errTx != nil {
				err = fmt.Errorf("creating employee: commiting transaction error: %w", errTx)
			}
		}
	}()
	isExist, err := s.repo.FindByName(tx, request.Name)
	if err != nil {
		return Response{}, fmt.Errorf("error finding employee: %w", err)
	}
	if isExist {
		return Response{}, common.AlreadyExistsError{Message: fmt.Sprintf("employee already exists: %v", request.Name)}
	}
	var id int64
	id, err = s.repo.Save(tx, request.ToEntity())
	if err != nil {
		return Response{}, fmt.Errorf("error saving employee with: %w", err)
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
		return Response{}, common.NotFoundError{Message: fmt.Sprintf("error finding employee with id %d: %v", request.Id, err)}
	}
	return entity.toResponse(), nil
}

func (s *Service) FindAll() ([]Response, error) {
	var employees, err = s.repo.FindAll()
	if err != nil {
		return nil, common.NotFoundError{Message: fmt.Sprintf("error finding all employees: %v", err)}
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
		return nil, common.NotFoundError{Message: fmt.Sprintf("error finding employees by ids: %v", err)}
	}
	var response []Response
	for _, employee := range employees {
		response = append(response, employee.toResponse())
	}
	return response, nil
}

func (s *Service) FindWithOffset(request PageRequest) (PageResponse, error) {
	if err := s.validator.Validate(request); err != nil {
		return PageResponse{}, common.RequestValidationError{Message: err.Error()}
	}
	employees, err := s.repo.FindWithOffset(request.PageSize*request.PageNumber, request.PageSize, request.TextFilter)
	if err != nil {
		return PageResponse{}, common.NotFoundError{Message: fmt.Sprintf("error finding employees with offset: %v", err)}
	}
	var response []Response
	for _, employee := range employees {
		response = append(response, employee.toResponse())
	}
	return PageResponse{
		Result:     response,
		PageSize:   request.PageSize,
		PageNumber: request.PageNumber,
		Total:      int64(len(response)),
	}, nil
}

func (s *Service) DeleteById(request IdRequest) error {
	var err = s.validator.Validate(request)
	if err != nil {
		return common.RequestValidationError{Message: err.Error()}
	}
	err = s.repo.DeleteById(request.Id)
	if err != nil {
		return common.NotFoundError{Message: fmt.Sprintf("error deleting employee with id %d: %v", request.Id, err)}
	}
	return nil
}

func (s *Service) DeleteByIds(request IdsRequest) error {
	if err := s.validator.Validate(request); err != nil {
		return common.RequestValidationError{Message: err.Error()}
	}
	var err = s.repo.DeleteByIds(request.Ids)
	if err != nil {
		return common.NotFoundError{Message: fmt.Sprintf("error deleting employee with ids %d: %v", request.Ids, err)}
	}
	return nil
}
