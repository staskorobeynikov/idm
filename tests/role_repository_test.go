package tests

import (
	"github.com/stretchr/testify/assert"
	"idm/inner/database"
	"idm/inner/role"
	"testing"
)

func TestRoleRepository(t *testing.T) {
	a := assert.New(t)
	var db = database.ConnectDb()
	var clearDatabase = func() {
		db.MustExec("DELETE FROM role")
	}
	defer func() {
		if r := recover(); r != nil {
			clearDatabase()
		}
	}()
	var roleRepository = role.NewRoleRepository(db)
	var roleFixture = NewRoleFixture(roleRepository)
	_ = roleFixture.CreateDatabase(db)
	t.Run("find an role by id", func(t *testing.T) {
		var newRoleId = roleFixture.Role("Test Name")
		got, err := roleRepository.GetById(newRoleId)
		a.Nil(err)
		a.NotEmpty(got)
		a.NotEmpty(got.Id)
		a.NotEmpty(got.CreatedAt)
		a.NotEmpty(got.UpdatedAt)
		a.Equal("Test Name", got.Name)
		clearDatabase()
	})
	t.Run("find all roles", func(t *testing.T) {
		_ = roleFixture.Role("Test Name")
		_ = roleFixture.Role("Test Name 1")
		_ = roleFixture.Role("Test Name 2")
		got, err := roleRepository.FindAll()
		a.Nil(err)
		a.NotEmpty(got)
		a.Equal("Test Name", got[0].Name)
		a.Equal("Test Name 1", got[1].Name)
		a.Equal("Test Name 2", got[2].Name)
		clearDatabase()
	})
	t.Run("find roles by ids", func(t *testing.T) {
		_ = roleFixture.Role("Test Name")
		var newRoleId1 = roleFixture.Role("Test Name 1")
		var newRoleId2 = roleFixture.Role("Test Name 2")
		_ = roleFixture.Role("Test Name 3")
		var newRoleId4 = roleFixture.Role("Test Name 4")
		ids := []int64{
			newRoleId1,
			newRoleId2,
			newRoleId4,
		}
		got, err := roleRepository.FindByIds(ids)
		a.Nil(err)
		a.NotEmpty(got)
		a.Equal("Test Name 1", got[0].Name)
		a.Equal("Test Name 2", got[1].Name)
		a.Equal("Test Name 4", got[2].Name)
		clearDatabase()
	})
	t.Run("delete role by id", func(t *testing.T) {
		_ = roleFixture.Role("Test Name")
		var newRoleId = roleFixture.Role("Test Name 1")
		_ = roleFixture.Role("Test Name 2")
		err := roleRepository.DeleteById(newRoleId)
		got, _ := roleRepository.FindAll()
		a.Nil(err)
		a.NotEmpty(got)
		a.Equal(len(got), 2)
		a.Equal("Test Name", got[0].Name)
		a.Equal("Test Name 2", got[1].Name)
		clearDatabase()
	})
	t.Run("delete roles by ids", func(t *testing.T) {
		var newRoleId = roleFixture.Role("Test Name")
		_ = roleFixture.Role("Test Name 1")
		var newRoleId2 = roleFixture.Role("Test Name 2")
		_ = roleFixture.Role("Test Name 3")
		var newRoleId4 = roleFixture.Role("Test Name 4")
		ids := []int64{
			newRoleId,
			newRoleId2,
			newRoleId4,
		}
		err := roleRepository.DeleteByIds(ids)
		got, _ := roleRepository.FindAll()
		a.Nil(err)
		a.NotEmpty(got)
		a.Equal(len(got), 2)
		a.Equal("Test Name 1", got[0].Name)
		a.Equal("Test Name 3", got[1].Name)
		clearDatabase()
	})
}
