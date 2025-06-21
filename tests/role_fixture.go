package tests

import (
	"github.com/jmoiron/sqlx"
	"idm/inner/role"
	"os"
)

type RoleFixture struct {
	roles *role.RoleRepository
}

func NewRoleFixture(roles *role.RoleRepository) *RoleFixture {
	return &RoleFixture{
		roles: roles,
	}
}

func (f *RoleFixture) Role(name string) int64 {
	var entity = role.RoleEntity{
		Name: name,
	}
	var newId, err = f.roles.Add(entity)
	if err != nil {
		panic(err)
	}
	return newId
}

func (f *RoleFixture) CreateDatabase(db *sqlx.DB) error {
	data, _ := os.ReadFile("./scripts/role.sql")
	_, err := db.Exec(string(data))
	if err != nil {
		return err
	}
	return nil
}
