package employee

import (
	"errors"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
)

type MockRepo struct {
	mock.Mock
}

func (r *MockRepo) BeginTransaction() (*sqlx.Tx, error) {
	args := r.Called()
	return args.Get(0).(*sqlx.Tx), args.Error(1)
}

func (r *MockRepo) Save(tx *sqlx.Tx, e Entity) (int64, error) {
	args := r.Called(tx, e)
	return args.Get(0).(int64), args.Error(1)
}

func (r *MockRepo) FindById(id int64) (employee Entity, err error) {
	args := r.Called(id)
	return args.Get(0).(Entity), args.Error(1)
}

func (r *MockRepo) FindByName(tx *sqlx.Tx, name string) (bool, error) {
	args := r.Called(tx, name)
	return args.Bool(0), args.Error(1)
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
	t.Run("should return wrapped error because begin transaction was failed", func(t *testing.T) {
		a := assert.New(t)
		db, mck, err := sqlmock.New()
		mck.ExpectBegin().WillReturnError(fmt.Errorf("error creating transaction"))
		sqlxDb := sqlx.NewDb(db, "sqlmock")
		repo := NewRepository(sqlxDb)
		_, err = repo.BeginTransaction()
		a.Error(err)
		a.Equal("error creating transaction", err.Error())
	})
	t.Run("should return wrapped error because findByName was failed", func(t *testing.T) {
		a := assert.New(t)
		db, mck, err := sqlmock.New()
		if err != nil {
			t.Fatal(err)
		}
		sqlxDb := sqlx.NewDb(db, "sqlmock")
		mck.ExpectBegin()
		tx, err := sqlxDb.Beginx()
		if err != nil {
			t.Fatal(err)
		}
		var repo = new(MockRepo)
		var svc = NewService(repo)
		repo.On("BeginTransaction").Return(tx, nil)
		var entity = Entity{
			Id:        1,
			Name:      "test",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			RoleId:    1,
		}
		err = errors.New("database error")
		repo.On("FindByName", tx, entity.Name).Return(false, err)
		var want = fmt.Errorf("error finding employee: %w", err)
		response, got := svc.Save(entity)
		a.Empty(response)
		a.NotNil(got)
		a.Equal(want, got)
		a.True(repo.AssertNumberOfCalls(t, "BeginTransaction", 1))
		a.True(repo.AssertNumberOfCalls(t, "FindByName", 1))
		a.True(repo.AssertNumberOfCalls(t, "Save", 0))
	})
	t.Run("should return response with id of employee because findByName return true", func(t *testing.T) {
		a := assert.New(t)
		db, mck, err := sqlmock.New()
		if err != nil {
			t.Fatal(err)
		}
		sqlxDb := sqlx.NewDb(db, "sqlmock")
		mck.ExpectBegin()
		tx, err := sqlxDb.Beginx()
		if err != nil {
			t.Fatal(err)
		}
		var repo = new(MockRepo)
		var svc = NewService(repo)
		repo.On("BeginTransaction").Return(tx, nil)
		var entity = Entity{
			Id:        1,
			Name:      "test",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			RoleId:    1,
		}
		err = errors.New("database error")
		repo.On("FindByName", tx, entity.Name).Return(true, nil)
		var want = Response{
			Id: 1,
		}
		got, err := svc.Save(entity)
		a.NotEmpty(got)
		a.Nil(err)
		a.Equal(want, got)
		a.True(repo.AssertNumberOfCalls(t, "BeginTransaction", 1))
		a.True(repo.AssertNumberOfCalls(t, "FindByName", 1))
		a.True(repo.AssertNumberOfCalls(t, "Save", 0))
	})
	t.Run("should return response with id of employee because findByName return true", func(t *testing.T) {
		a := assert.New(t)
		db, mck, err := sqlmock.New()
		if err != nil {
			t.Fatal(err)
		}
		sqlxDb := sqlx.NewDb(db, "sqlmock")
		mck.ExpectBegin()
		tx, err := sqlxDb.Beginx()
		if err != nil {
			t.Fatal(err)
		}
		var repo = new(MockRepo)
		var svc = NewService(repo)
		repo.On("BeginTransaction").Return(tx, nil)
		var entity = Entity{
			Id:        1,
			Name:      "test",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			RoleId:    1,
		}
		err = errors.New("database error")
		repo.On("FindByName", tx, entity.Name).Return(false, nil)
		repo.On("Save", tx, entity).Return(entity.Id, err)
		var want = fmt.Errorf("error saving employee with: %w", err)
		response, got := svc.Save(entity)
		a.Empty(response)
		a.NotNil(err)
		a.Equal(want, got)
		a.True(repo.AssertNumberOfCalls(t, "BeginTransaction", 1))
		a.True(repo.AssertNumberOfCalls(t, "FindByName", 1))
		a.True(repo.AssertNumberOfCalls(t, "Save", 1))
	})
	t.Run("should return id new employee because findByName return false", func(t *testing.T) {
		a := assert.New(t)
		db, mck, err := sqlmock.New()
		if err != nil {
			t.Fatal(err)
		}
		sqlxDb := sqlx.NewDb(db, "sqlmock")
		mck.ExpectBegin()
		tx, err := sqlxDb.Beginx()
		if err != nil {
			t.Fatal(err)
		}
		var repo = new(MockRepo)
		var svc = NewService(repo)
		repo.On("BeginTransaction").Return(tx, nil)
		var entity = Entity{
			Id:        1,
			Name:      "test",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			RoleId:    1,
		}
		err = errors.New("database error")
		repo.On("FindByName", tx, entity.Name).Return(false, nil)
		repo.On("Save", tx, entity).Return(entity.Id, nil)
		var want = Response{
			Id: 1,
		}
		got, err := svc.Save(entity)
		a.NotEmpty(got)
		a.Nil(err)
		a.Equal(want, got)
		a.True(repo.AssertNumberOfCalls(t, "BeginTransaction", 1))
		a.True(repo.AssertNumberOfCalls(t, "FindByName", 1))
		a.True(repo.AssertNumberOfCalls(t, "Save", 1))
	})
}

func TestFindById(t *testing.T) {
	var a = assert.New(t)
	t.Run("should return found employee", func(t *testing.T) {
		var repo = new(MockRepo)
		var svc = NewService(repo)
		var entity = Entity{
			Id:        1,
			Name:      "test",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		var want = entity.toResponse()
		repo.On("FindById", int64(1)).Return(entity, nil)
		var got, err = svc.FindById(1)
		a.Nil(err)
		a.Equal(want, got)
		a.True(repo.AssertNumberOfCalls(t, "FindById", 1))
	})
	t.Run("should return wrapped error", func(t *testing.T) {
		var repo = new(MockRepo)
		var svc = NewService(repo)
		var entity = Entity{}
		var err = errors.New("database error")
		var id = int64(1)
		var want = fmt.Errorf("error finding employee with id %d: %w", id, err)
		repo.On("FindById", id).Return(entity, err)
		var response, got = svc.FindById(id)
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
		var svc = NewService(repo)
		var entities = []Entity{
			{Id: 1, Name: "test1", CreatedAt: time.Now(), UpdatedAt: time.Now(), RoleId: 1},
			{Id: 1, Name: "test2", CreatedAt: time.Now(), UpdatedAt: time.Now(), RoleId: 1},
			{Id: 1, Name: "test3", CreatedAt: time.Now(), UpdatedAt: time.Now(), RoleId: 1},
			{Id: 1, Name: "test4", CreatedAt: time.Now(), UpdatedAt: time.Now(), RoleId: 1},
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
		var svc = NewService(repo)
		var entities []Entity
		var err = errors.New("database error")
		var want = fmt.Errorf("error finding all employees: %w", err)
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
		var svc = NewService(repo)
		var entities = []Entity{
			{Id: 2, Name: "test2", CreatedAt: time.Now(), UpdatedAt: time.Now(), RoleId: 1},
			{Id: 4, Name: "test4", CreatedAt: time.Now(), UpdatedAt: time.Now(), RoleId: 1},
		}
		var want []Response
		for _, entity := range entities {
			if entity.Id%2 == 0 {
				want = append(want, entity.toResponse())
			}
		}
		var ids = []int64{2, 4}
		repo.On("FindByIds", ids).Return(entities, nil)
		var got, err = svc.FindByIds(ids)
		a.Nil(err)
		a.Equal(len(got), 2)
		a.Equal(want, got)
		a.True(repo.AssertNumberOfCalls(t, "FindByIds", 1))
	})
	t.Run("should return wrapped error", func(t *testing.T) {
		var repo = new(MockRepo)
		var svc = NewService(repo)
		var entities []Entity
		var err = errors.New("database error")
		var want = fmt.Errorf("error finding employees by ids: %w", err)
		var ids = []int64{2, 4}
		repo.On("FindByIds", ids).Return(entities, err)
		var response, got = svc.FindByIds(ids)
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
		var svc = NewService(repo)
		repo.On("DeleteById", int64(1)).Return(nil)
		var got = svc.DeleteById(1)
		a.Nil(got)
		a.True(repo.AssertNumberOfCalls(t, "DeleteById", 1))
	})
	t.Run("should return wrapped error", func(t *testing.T) {
		var repo = new(MockRepo)
		var svc = NewService(repo)
		var err = errors.New("database error")
		var id = int64(1)
		var want = fmt.Errorf("error deleting employee with id %d: %w", id, err)
		repo.On("DeleteById", id).Return(err)
		var got = svc.DeleteById(id)
		a.NotNil(got)
		a.Equal(want, got)
		a.True(repo.AssertNumberOfCalls(t, "DeleteById", 1))
	})
}

func TestDeleteByIds(t *testing.T) {
	var a = assert.New(t)
	t.Run("should delete employee by ids", func(t *testing.T) {
		var repo = new(MockRepo)
		var svc = NewService(repo)
		repo.On("DeleteByIds", []int64{2, 4}).Return(nil)
		var got = svc.DeleteByIds([]int64{2, 4})
		a.Nil(got)
		a.True(repo.AssertNumberOfCalls(t, "DeleteByIds", 1))
	})
	t.Run("should return wrapped error", func(t *testing.T) {
		var repo = new(MockRepo)
		var svc = NewService(repo)
		var err = errors.New("database error")
		var ids = []int64{2, 4}
		var want = fmt.Errorf("error deleting employee with ids %d: %w", ids, err)
		repo.On("DeleteByIds", ids).Return(err)
		var got = svc.DeleteByIds(ids)
		a.NotNil(got)
		a.Equal(want, got)
		a.True(repo.AssertNumberOfCalls(t, "DeleteByIds", 1))
	})
}
