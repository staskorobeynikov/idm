package employee

import (
	"context"
	"errors"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"idm/inner/common"
	"idm/inner/validator"
	"strings"
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

func (r *MockRepo) FindWithOffset(offset int, limit int, filter string) ([]Entity, error) {
	args := r.Called(offset, limit, filter)
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
		db, mck, _ := sqlmock.New()
		mck.ExpectBegin().WillReturnError(fmt.Errorf("error creating transaction"))
		sqlxDb := sqlx.NewDb(db, "sqlmock")
		repo := NewRepository(sqlxDb)
		_, err := repo.BeginTransaction()
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
		var svc = NewService(repo, validator.New())
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
		response, got := svc.Save(
			context.Background(),
			CreateRequest{
				Name:   "test",
				RoleId: 1,
			})
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
		var svc = NewService(repo, validator.New())
		repo.On("BeginTransaction").Return(tx, nil)
		var entity = Entity{
			Id:        1,
			Name:      "test",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			RoleId:    1,
		}
		repo.On("FindByName", tx, entity.Name).Return(true, nil)
		got, err := svc.Save(
			context.Background(),
			CreateRequest{
				Name:   "test",
				RoleId: 1,
			})
		var want = common.AlreadyExistsError{Message: fmt.Sprintf("employee already exists: %v", entity.Name)}
		a.Empty(got)
		a.NotNil(err)
		a.Equal(want, err)
		a.True(repo.AssertNumberOfCalls(t, "BeginTransaction", 1))
		a.True(repo.AssertNumberOfCalls(t, "FindByName", 1))
		a.True(repo.AssertNumberOfCalls(t, "Save", 0))
	})
	t.Run("should return error because save fails", func(t *testing.T) {
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
		var svc = NewService(repo, validator.New())
		repo.On("BeginTransaction").Return(tx, nil)
		var entity = Entity{
			Name:      "test",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			RoleId:    int64(1),
		}
		err = errors.New("database error")
		repo.On("FindByName", tx, entity.Name).Return(false, nil)
		repo.On("Save", mock.Anything, mock.MatchedBy(func(e Entity) bool {
			return e.Name == "test" && e.RoleId == 1
		})).Return(int64(-1), err)
		var want = fmt.Errorf("error saving employee with: %w", err)

		response, got := svc.Save(
			context.Background(),
			CreateRequest{
				Name:   "test",
				RoleId: 1,
			})
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
		var svc = NewService(repo, validator.New())
		repo.On("BeginTransaction").Return(tx, nil)
		var entity = Entity{
			Id:        1,
			Name:      "test",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			RoleId:    1,
		}
		repo.On("FindByName", tx, entity.Name).Return(false, nil)
		repo.On("Save", mock.Anything, mock.MatchedBy(func(e Entity) bool {
			return e.Name == "test" && e.RoleId == 1
		})).Return(entity.Id, nil)
		var want = Response{
			Id: 1,
		}
		got, err := svc.Save(
			context.Background(),
			CreateRequest{
				Name:   "test",
				RoleId: 1,
			})
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
		var want = common.NotFoundError{Message: fmt.Sprintf("error finding employee with id %d: %v", id, err)}
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
		var svc = NewService(repo, validator.New())
		var entities []Entity
		var err = errors.New("database error")
		var want = common.NotFoundError{Message: fmt.Sprintf("error finding all employees: %v", err)}
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
		var want = common.NotFoundError{Message: fmt.Sprintf("error finding employees by ids: %v", err)}
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
		var want = common.NotFoundError{Message: fmt.Sprintf("error deleting employee with id %d: %v", id, err)}
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
		var want = common.NotFoundError{Message: fmt.Sprintf("error deleting employee with ids %d: %v", ids, err)}
		repo.On("DeleteByIds", ids).Return(err)
		var got = svc.DeleteByIds(IdsRequest{Ids: ids})
		a.NotNil(got)
		a.Equal(want, got)
		a.True(repo.AssertNumberOfCalls(t, "DeleteByIds", 1))
	})
}

func TestCreateRequest(t *testing.T) {
	a := assert.New(t)
	v := validator.New()
	t.Run("correct create request", func(t *testing.T) {
		err := v.Validate(CreateRequest{
			Name:   "test",
			RoleId: 1,
		})
		a.Nil(err)
	})
	t.Run("incorrect create request - no name", func(t *testing.T) {
		err := v.Validate(CreateRequest{
			RoleId: 1,
		})
		want := "Key: 'CreateRequest.Name' Error:Field validation for 'Name' failed on the 'required' tag"
		a.Error(err)
		a.Equal(want, err.Error())
	})
	t.Run("incorrect create request - short name", func(t *testing.T) {
		err := v.Validate(CreateRequest{
			Name:   "t",
			RoleId: 1,
		})
		want := "Key: 'CreateRequest.Name' Error:Field validation for 'Name' failed on the 'min' tag"
		a.Error(err)
		a.Equal(want, err.Error())
	})
	t.Run("incorrect create request - long name", func(t *testing.T) {
		err := v.Validate(CreateRequest{
			Name:   strings.Repeat("abcde", 32),
			RoleId: 1,
		})
		want := "Key: 'CreateRequest.Name' Error:Field validation for 'Name' failed on the 'max' tag"
		a.Error(err)
		a.Equal(want, err.Error())
	})
	t.Run("incorrect create request - no roleId", func(t *testing.T) {
		err := v.Validate(CreateRequest{
			Name: "test",
		})
		want := "Key: 'CreateRequest.RoleId' Error:Field validation for 'RoleId' failed on the 'required' tag"
		a.Error(err)
		a.Equal(want, err.Error())
	})
	t.Run("incorrect create request - roleId less than min", func(t *testing.T) {
		err := v.Validate(CreateRequest{
			Name:   "test",
			RoleId: -1,
		})
		want := "Key: 'CreateRequest.RoleId' Error:Field validation for 'RoleId' failed on the 'min' tag"
		a.Error(err)
		a.Equal(want, err.Error())
	})
}

func TestIdRequest(t *testing.T) {
	a := assert.New(t)
	v := validator.New()
	t.Run("correct id request", func(t *testing.T) {
		err := v.Validate(IdRequest{Id: int64(1)})
		a.Nil(err)
	})
	t.Run("incorrect id request - error required", func(t *testing.T) {
		err := v.Validate(IdRequest{})
		want := "Key: 'IdRequest.Id' Error:Field validation for 'Id' failed on the 'required' tag"
		a.Error(err)
		a.Equal(want, err.Error())
	})
	t.Run("incorrect id request - error required", func(t *testing.T) {
		err := v.Validate(IdRequest{Id: int64(-1)})
		want := "Key: 'IdRequest.Id' Error:Field validation for 'Id' failed on the 'min' tag"
		a.Error(err)
		a.Equal(want, err.Error())
	})
}

func TestIdsRequest(t *testing.T) {
	a := assert.New(t)
	v := validator.New()
	t.Run("correct ids request", func(t *testing.T) {
		err := v.Validate(IdsRequest{
			Ids: []int64{1, 2, 3},
		})
		a.Nil(err)
	})
	t.Run("incorrect ids request - error required", func(t *testing.T) {
		err := v.Validate(IdsRequest{})
		want := "Key: 'IdsRequest.Ids' Error:Field validation for 'Ids' failed on the 'required' tag"
		a.Error(err)
		a.Equal(want, err.Error())
	})
	t.Run("incorrect ids request - error length", func(t *testing.T) {
		err := v.Validate(IdsRequest{
			Ids: []int64{},
		})
		want := "Key: 'IdsRequest.Ids' Error:Field validation for 'Ids' failed on the 'min' tag"
		a.Error(err)
		a.Equal(want, err.Error())
	})
}

func TestPageRequest(t *testing.T) {
	a := assert.New(t)
	v := validator.New()
	t.Run("correct page", func(t *testing.T) {
		err := v.Validate(PageRequest{
			PageSize:   5,
			PageNumber: 5,
		})
		a.Nil(err)
	})
	t.Run("incorrect page - PageSize < 1", func(t *testing.T) {
		err := v.Validate(PageRequest{
			PageSize:   0,
			PageNumber: 5,
		})
		want := "Key: 'PageRequest.PageSize' Error:Field validation for 'PageSize' failed on the 'min' tag"
		a.Error(err)
		a.Equal(want, err.Error())
	})
	t.Run("incorrect page - PageSize > 105", func(t *testing.T) {
		err := v.Validate(PageRequest{
			PageSize:   105,
			PageNumber: 5,
		})
		want := "Key: 'PageRequest.PageSize' Error:Field validation for 'PageSize' failed on the 'max' tag"
		a.Error(err)
		a.Equal(want, err.Error())
	})
	t.Run("incorrect page - PageNumber < 0", func(t *testing.T) {
		err := v.Validate(PageRequest{
			PageSize:   50,
			PageNumber: -1,
		})
		want := "Key: 'PageRequest.PageNumber' Error:Field validation for 'PageNumber' failed on the 'min' tag"
		a.Error(err)
		a.Equal(want, err.Error())
	})
	t.Run("incorrect page - TextFilter is empty", func(t *testing.T) {
		err := v.Validate(PageRequest{
			PageSize:   50,
			PageNumber: -1,
			TextFilter: "",
		})
		want := "Key: 'PageRequest.PageNumber' Error:Field validation for 'PageNumber' failed on the 'min' tag"
		a.Error(err)
		a.Equal(want, err.Error())
	})
	t.Run("incorrect page - TextFilter is space symbols", func(t *testing.T) {
		err := v.Validate(PageRequest{
			PageSize:   50,
			PageNumber: 1,
			TextFilter: "    ",
		})
		want := "Key: 'PageRequest.TextFilter' Error:Field validation for 'TextFilter' failed on the 'minnows3' tag"
		a.Error(err)
		a.Equal(want, err.Error())
	})
	t.Run("incorrect page - TextFilter is \n symbols", func(t *testing.T) {
		err := v.Validate(PageRequest{
			PageSize:   50,
			PageNumber: 1,
			TextFilter: "\n\n\n",
		})
		want := "Key: 'PageRequest.TextFilter' Error:Field validation for 'TextFilter' failed on the 'minnows3' tag"
		a.Error(err)
		a.Equal(want, err.Error())
	})
	t.Run("incorrect page - TextFilter has not enough text symbols", func(t *testing.T) {
		err := v.Validate(PageRequest{
			PageSize:   50,
			PageNumber: 1,
			TextFilter: "a  b ",
		})
		want := "Key: 'PageRequest.TextFilter' Error:Field validation for 'TextFilter' failed on the 'minnows3' tag"
		a.Error(err)
		a.Equal(want, err.Error())
	})
	t.Run("correct page", func(t *testing.T) {
		err := v.Validate(PageRequest{
			PageSize:   50,
			PageNumber: 1,
			TextFilter: "a  b c d",
		})
		a.NoError(err)
	})
}
