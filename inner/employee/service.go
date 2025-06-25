package employee

import (
	"fmt"
	"github.com/jmoiron/sqlx"
)

type Service struct {
	repo Repo
}

type Repo interface {
	BeginTransaction() (*sqlx.Tx, error)
	Save(tx *sqlx.Tx, e Entity) (int64, error)
	FindById(id int64) (Entity, error)
	FindByName(tx *sqlx.Tx, name string) (bool, error)
	FindAll() ([]Entity, error)
	FindByIds(ids []int64) ([]Entity, error)
	DeleteById(id int64) error
	DeleteByIds(ids []int64) error
}

func NewService(repo Repo) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) Save(e Entity) (Response, error) {
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
	isExist, err := s.repo.FindByName(tx, e.Name)
	if err != nil {
		return Response{}, fmt.Errorf("error finding employee: %w", err)
	}
	if isExist {
		return Response{
			Id: e.Id,
		}, nil
	}
	var id, err1 = s.repo.Save(tx, e)
	if err1 != nil {
		return Response{}, fmt.Errorf("error saving employee with: %w", err1)
	}
	return Response{
		Id: id,
	}, nil
}

func (s *Service) FindById(id int64) (Response, error) {
	var entity, err = s.repo.FindById(id)
	if err != nil {
		return Response{}, fmt.Errorf("error finding employee with id %d: %w", id, err)
	}
	return entity.toResponse(), nil
}

func (s *Service) FindAll() ([]Response, error) {
	var employees, err = s.repo.FindAll()
	if err != nil {
		return nil, fmt.Errorf("error finding all employees: %w", err)
	}
	var response []Response
	for _, employee := range employees {
		response = append(response, employee.toResponse())
	}
	return response, nil
}

func (s *Service) FindByIds(ids []int64) ([]Response, error) {
	var employees, err = s.repo.FindByIds(ids)
	if err != nil {
		return nil, fmt.Errorf("error finding employees by ids: %w", err)
	}
	var response []Response
	for _, employee := range employees {
		response = append(response, employee.toResponse())
	}
	return response, nil
}

func (s *Service) DeleteById(id int64) error {
	var err = s.repo.DeleteById(id)
	if err != nil {
		return fmt.Errorf("error deleting employee with id %d: %w", id, err)
	}
	return nil
}

func (s *Service) DeleteByIds(ids []int64) error {
	var err = s.repo.DeleteByIds(ids)
	if err != nil {
		return fmt.Errorf("error deleting employee with ids %d: %w", ids, err)
	}
	return nil
}
