package tests

import (
	"github.com/jmoiron/sqlx"
	"idm/inner/employee"
	"os"
)

type Fixture struct {
	employees *employee.Repository
}

func NewFixture(employees *employee.Repository) *Fixture {
	return &Fixture{
		employees: employees,
	}
}

func (f *Fixture) Employee(name string, roleId int64) int64 {
	var entity = employee.Entity{
		Name:   name,
		RoleId: roleId,
	}
	var newId, err = f.employees.Save(entity)
	if err != nil {
		panic(err)
	}
	return newId
}

func (f *Fixture) CreateDatabase(db *sqlx.DB) error {
	data, _ := os.ReadFile("./scripts/employee.sql")
	_, err := db.Exec(string(data))
	if err != nil {
		return err
	}
	return nil
}
