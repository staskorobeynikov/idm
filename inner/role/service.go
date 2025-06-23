package role

import "fmt"

type Service struct {
	repo Repo
}

type Repo interface {
	Save(entity Entity) (int64, error)
	FindById(id int64) (entity Entity, err error)
	FindAll() ([]Entity, error)
	FindByIds(ids []int64) ([]Entity, error)
	DeleteById(id int64) error
	DeleteByIds(ids []int64) error
}

func NewService(repo Repo) Service {
	return Service{
		repo: repo,
	}
}

func (s *Service) Save(e Entity) (Response, error) {
	var id, err = s.repo.Save(e)
	if err != nil {
		return Response{}, fmt.Errorf("error saving role with: %w", err)
	}
	return Response{
		Id: id,
	}, nil
}

func (s *Service) FindById(id int64) (Response, error) {
	var entity, err = s.repo.FindById(id)
	if err != nil {
		return Response{}, fmt.Errorf("error finding role with id %d: %w", id, err)
	}
	return entity.toResponse(), nil
}

func (s *Service) FindAll() ([]Response, error) {
	var employees, err = s.repo.FindAll()
	if err != nil {
		return nil, fmt.Errorf("error finding all roles: %w", err)
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
		return nil, fmt.Errorf("error finding roles by ids: %w", err)
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
		return fmt.Errorf("error deleting role with id %d: %w", id, err)
	}
	return nil
}

func (s *Service) DeleteByIds(ids []int64) error {
	var err = s.repo.DeleteByIds(ids)
	if err != nil {
		return fmt.Errorf("error deleting roles with ids %d: %w", ids, err)
	}
	return nil
}
