package tests

import (
	"context"
	"github.com/jmoiron/sqlx"
	"idm/inner/employee"
	"os"
)

type Fixture struct {
	db        *sqlx.DB
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
	tx, err := f.db.BeginTxx(context.Background(), nil)
	if err != nil {
		return -1
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		} else {
			_ = tx.Commit()
		}
	}()
	var newId, _ = f.employees.Save(tx, entity)
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
