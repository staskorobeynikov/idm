package tests

import (
	"github.com/stretchr/testify/assert"
	"idm/inner/database"
	"idm/inner/employee"
	"idm/inner/role"
	"testing"
	"time"
)

func TestEmployeeRepository(t *testing.T) {
	a := assert.New(t)
	var db = database.ConnectDb()
	var clearDatabase = func() {
		db.MustExec("DELETE FROM employee")
	}
	defer func() {
		if r := recover(); r != nil {
			clearDatabase()
		}
	}()
	defer func() {
		db.MustExec("DELETE FROM role")
	}()
	var employeeRepository = employee.NewRepository(db)
	var emplFixture = Fixture{
		employees: employeeRepository,
		db:        db,
	}
	var roleRepository = role.NewRepository(db)
	var roleFixture = NewRoleFixture(roleRepository)
	var newRoleId = roleFixture.Role("Test Name")
	_ = emplFixture.CreateDatabase(db)
	t.Run("find an employee by id", func(t *testing.T) {
		var newEmployeeId = emplFixture.Employee("Test Name", newRoleId)
		got, err := employeeRepository.FindById(newEmployeeId)
		a.Nil(err)
		a.NotEmpty(got)
		a.NotEmpty(got.Id)
		a.NotEmpty(got.CreatedAt)
		a.NotEmpty(got.UpdatedAt)
		a.Equal("Test Name", got.Name)
		clearDatabase()
	})
	t.Run("find all employees", func(t *testing.T) {
		_ = emplFixture.Employee("Test Name", newRoleId)
		_ = emplFixture.Employee("Test Name 1", newRoleId)
		_ = emplFixture.Employee("Test Name 2", newRoleId)
		_ = emplFixture.Employee("Test Name 3", newRoleId)
		got, err := employeeRepository.FindAll()
		a.Nil(err)
		a.NotEmpty(got)
		a.Equal(len(got), 4)
		a.Equal("Test Name", got[0].Name)
		a.Equal("Test Name 1", got[1].Name)
		a.Equal("Test Name 2", got[2].Name)
		a.Equal("Test Name 3", got[3].Name)
		clearDatabase()
	})
	t.Run("find employees by ids", func(t *testing.T) {
		var newEmployeeId = emplFixture.Employee("Test Name", newRoleId)
		_ = emplFixture.Employee("Test Name 1", newRoleId)
		_ = emplFixture.Employee("Test Name 2", newRoleId)
		_ = emplFixture.Employee("Test Name 3", newRoleId)
		var newEmployeeId1 = emplFixture.Employee("Test Name 4", newRoleId)
		var newEmployeeId2 = emplFixture.Employee("Test Name 5", newRoleId)
		ids := []int64{
			newEmployeeId,
			newEmployeeId1,
			newEmployeeId2,
		}
		got, err := employeeRepository.FindByIds(ids)
		a.Nil(err)
		a.NotEmpty(got)
		a.Equal(len(got), 3)
		a.Equal("Test Name", got[0].Name)
		a.Equal("Test Name 4", got[1].Name)
		a.Equal("Test Name 5", got[2].Name)
		clearDatabase()
	})
	t.Run("delete employee by id", func(t *testing.T) {
		_ = emplFixture.Employee("Test Name", newRoleId)
		_ = emplFixture.Employee("Test Name 1", newRoleId)
		var newEmployeeId = emplFixture.Employee("Test Name 2", newRoleId)
		_ = emplFixture.Employee("Test Name 3", newRoleId)
		err := employeeRepository.DeleteById(newEmployeeId)
		got, _ := employeeRepository.FindAll()
		a.Nil(err)
		a.NotEmpty(got)
		a.Equal(len(got), 3)
		a.Equal("Test Name", got[0].Name)
		a.Equal("Test Name 1", got[1].Name)
		a.Equal("Test Name 3", got[2].Name)
		clearDatabase()
	})
	t.Run("delete employees by ids", func(t *testing.T) {
		_ = emplFixture.Employee("Test Name", newRoleId)
		var newEmployeeId = emplFixture.Employee("Test Name 1", newRoleId)
		var newEmployeeId1 = emplFixture.Employee("Test Name 2", newRoleId)
		var newEmployeeId2 = emplFixture.Employee("Test Name 3", newRoleId)
		_ = emplFixture.Employee("Test Name 4", newRoleId)
		var newEmployeeId3 = emplFixture.Employee("Test Name 5", newRoleId)
		ids := []int64{
			newEmployeeId,
			newEmployeeId1,
			newEmployeeId2,
			newEmployeeId3,
		}
		err := employeeRepository.DeleteByIds(ids)
		got, _ := employeeRepository.FindAll()
		a.Nil(err)
		a.NotEmpty(got)
		a.Equal(len(got), 2)
		a.Equal("Test Name", got[0].Name)
		a.Equal("Test Name 4", got[1].Name)
		clearDatabase()
	})
	t.Run("find by name and save employee in one tx", func(t *testing.T) {
		tx, err := employeeRepository.BeginTransaction()
		a.NoError(err)
		defer func() {
			err := tx.Rollback()
			a.NoError(err)
		}()
		empl := employee.Entity{
			Id:        1,
			Name:      "Test Name",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			RoleId:    newRoleId,
		}
		isExist, _ := employeeRepository.FindByName(tx, empl.Name)
		a.False(isExist)
		got, err := employeeRepository.Save(tx, empl)
		a.NoError(err)
		a.NotEmpty(got)
		found, err := employeeRepository.FindByName(tx, empl.Name)
		a.NoError(err)
		a.NotEmpty(found)
		a.True(found)
	})
}
