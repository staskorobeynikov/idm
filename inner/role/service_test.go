package role

import (
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"idm/inner/common"
	"idm/inner/validator"
	"testing"
	"time"
)

type MockRepo struct {
	mock.Mock
}

func (r *MockRepo) Save(e Entity) (int64, error) {
	args := r.Called(e)
	return args.Get(0).(int64), args.Error(1)
}

func (r *MockRepo) FindById(id int64) (employee Entity, err error) {
	args := r.Called(id)
	return args.Get(0).(Entity), args.Error(1)
}

func (r *MockRepo) FindAll() ([]Entity, error) {
	args := r.Called()
	return args.Get(0).([]Entity), args.Error(1)
}

func (r *MockRepo) FindByIds(ids []int64) ([]Entity, error) {
	args := r.Called(ids)
	return args.Get(0).([]Entity), args.Error(1)
}

func (r *MockRepo) DeleteById(id int64) error {
	args := r.Called(id)
	return args.Error(0)
}

func (r *MockRepo) DeleteByIds(ids []int64) error {
	args := r.Called(ids)
	return args.Error(0)
}

func TestSave(t *testing.T) {
	var a = assert.New(t)
	t.Run("should return id new employee", func(t *testing.T) {
		var repo = new(MockRepo)
		var svc = NewService(repo, validator.New())
		var entity = Entity{
			Name: "test",
		}
		repo.On("Save", entity).Return(entity.Id, nil)
		_, err := svc.Save(CreateRequest{
			Name: entity.Name,
		})
		a.Nil(err)
		a.True(repo.AssertNumberOfCalls(t, "Save", 1))
	})
	t.Run("should return wrapped error", func(t *testing.T) {
		var repo = new(MockRepo)
		var svc = NewService(repo, validator.New())
		var entity = Entity{
			Name: "test",
		}
		var err = errors.New("database error")
		var want = fmt.Errorf("error saving role with: %w", err)
		repo.On("Save", entity).Return(entity.Id, err)
		var response, got = svc.Save(CreateRequest{
			Name: entity.Name,
		})
		a.Empty(response)
		a.NotNil(got)
		a.Equal(want, got)
		a.True(repo.AssertNumberOfCalls(t, "Save", 1))
	})
}

func TestFindById(t *testing.T) {
	var a = assert.New(t)
	t.Run("should return found employee", func(t *testing.T) {
		var repo = new(MockRepo)
		var svc = NewService(repo, validator.New())
		var entity = Entity{
			Id:        1,
			Name:      "test",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		var want = entity.toResponse()
		repo.On("FindById", int64(1)).Return(entity, nil)
		var got, err = svc.FindById(IdRequest{Id: int64(1)})
		a.Nil(err)
		a.Equal(want, got)
		a.True(repo.AssertNumberOfCalls(t, "FindById", 1))
	})
	t.Run("should return wrapped error", func(t *testing.T) {
		var repo = new(MockRepo)
		var svc = NewService(repo, validator.New())
		var entity = Entity{}
		var err = errors.New("database error")
		var id = int64(1)
		var want = common.NotFoundError{Message: fmt.Sprintf("error finding role with id %d: %v", id, err)}
		repo.On("FindById", id).Return(entity, err)
		var response, got = svc.FindById(IdRequest{Id: id})
		a.Empty(response)
		a.NotNil(got)
		a.Equal(want, got)
		a.True(repo.AssertNumberOfCalls(t, "FindById", 1))
	})
}

func TestFindAll(t *testing.T) {
	var a = assert.New(t)
	t.Run("should return all employees", func(t *testing.T) {
		var repo = new(MockRepo)
		var svc = NewService(repo, validator.New())
		var entities = []Entity{
			{Id: 1, Name: "test1", CreatedAt: time.Now(), UpdatedAt: time.Now()},
			{Id: 1, Name: "test2", CreatedAt: time.Now(), UpdatedAt: time.Now()},
			{Id: 1, Name: "test3", CreatedAt: time.Now(), UpdatedAt: time.Now()},
			{Id: 1, Name: "test4", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		}
		var want []Response
		for _, entity := range entities {
			want = append(want, entity.toResponse())
		}
		repo.On("FindAll").Return(entities, nil)
		var got, err = svc.FindAll()
		a.Nil(err)
		a.Equal(want, got)
		a.True(repo.AssertNumberOfCalls(t, "FindAll", 1))
	})
	t.Run("should return wrapped error", func(t *testing.T) {
		var repo = new(MockRepo)
		var svc = NewService(repo, validator.New())
		var entities []Entity
		var err = errors.New("database error")
		var want = common.NotFoundError{Message: fmt.Sprintf("error finding all roles: %v", err)}
		repo.On("FindAll").Return(entities, err)
		var response, got = svc.FindAll()
		a.Empty(response)
		a.NotNil(got)
		a.Equal(want, got)
		a.True(repo.AssertNumberOfCalls(t, "FindAll", 1))
	})
}

func TestFindByIds(t *testing.T) {
	var a = assert.New(t)
	t.Run("should return employees by ids", func(t *testing.T) {
		var repo = new(MockRepo)
		var svc = NewService(repo, validator.New())
		var entities = []Entity{
			{Id: 2, Name: "test2", CreatedAt: time.Now(), UpdatedAt: time.Now()},
			{Id: 4, Name: "test4", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		}
		var want []Response
		for _, entity := range entities {
			if entity.Id%2 == 0 {
				want = append(want, entity.toResponse())
			}
		}
		var ids = []int64{2, 4}
		repo.On("FindByIds", ids).Return(entities, nil)
		var got, err = svc.FindByIds(IdsRequest{Ids: ids})
		a.Nil(err)
		a.Equal(len(got), 2)
		a.Equal(want, got)
		a.True(repo.AssertNumberOfCalls(t, "FindByIds", 1))
	})
	t.Run("should return wrapped error", func(t *testing.T) {
		var repo = new(MockRepo)
		var svc = NewService(repo, validator.New())
		var entities []Entity
		var err = errors.New("database error")
		var want = common.NotFoundError{Message: fmt.Sprintf("error finding roles by ids: %v", err)}
		var ids = []int64{2, 4}
		repo.On("FindByIds", ids).Return(entities, err)
		var response, got = svc.FindByIds(IdsRequest{Ids: ids})
		a.Empty(response)
		a.NotNil(got)
		a.Equal(want, got)
		a.True(repo.AssertNumberOfCalls(t, "FindByIds", 1))
	})
}

func TestDeleteById(t *testing.T) {
	var a = assert.New(t)
	t.Run("should delete employee by id", func(t *testing.T) {
		var repo = new(MockRepo)
		var svc = NewService(repo, validator.New())
		repo.On("DeleteById", int64(1)).Return(nil)
		var got = svc.DeleteById(IdRequest{Id: int64(1)})
		a.Nil(got)
		a.True(repo.AssertNumberOfCalls(t, "DeleteById", 1))
	})
	t.Run("should return wrapped error", func(t *testing.T) {
		var repo = new(MockRepo)
		var svc = NewService(repo, validator.New())
		var err = errors.New("database error")
		var id = int64(1)
		var want = common.NotFoundError{Message: fmt.Sprintf("error deleting role with id %d: %v", id, err)}
		repo.On("DeleteById", id).Return(err)
		var got = svc.DeleteById(IdRequest{Id: id})
		a.NotNil(got)
		a.Equal(want, got)
		a.True(repo.AssertNumberOfCalls(t, "DeleteById", 1))
	})
}

func TestDeleteByIds(t *testing.T) {
	var a = assert.New(t)
	t.Run("should delete employee by ids", func(t *testing.T) {
		var repo = new(MockRepo)
		var svc = NewService(repo, validator.New())
		repo.On("DeleteByIds", []int64{2, 4}).Return(nil)
		var got = svc.DeleteByIds(IdsRequest{Ids: []int64{2, 4}})
		a.Nil(got)
		a.True(repo.AssertNumberOfCalls(t, "DeleteByIds", 1))
	})
	t.Run("should return wrapped error", func(t *testing.T) {
		var repo = new(MockRepo)
		var svc = NewService(repo, validator.New())
		var err = errors.New("database error")
		var ids = []int64{2, 4}
		var want = common.NotFoundError{Message: fmt.Sprintf("error deleting roles with ids %d: %v", ids, err)}
		repo.On("DeleteByIds", ids).Return(err)
		var got = svc.DeleteByIds(IdsRequest{Ids: ids})
		a.NotNil(got)
		a.Equal(want, got)
		a.True(repo.AssertNumberOfCalls(t, "DeleteByIds", 1))
	})
}
