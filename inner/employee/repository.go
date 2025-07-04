package employee

import (
	"fmt"
	"github.com/jmoiron/sqlx"
)

type Repository struct {
	db *sqlx.DB
}

func NewRepository(database *sqlx.DB) *Repository {
	return &Repository{
		db: database,
	}
}

func (r *Repository) BeginTransaction() (*sqlx.Tx, error) {
	return r.db.Beginx()
}

func (r *Repository) Save(tx *sqlx.Tx, e Entity) (int64, error) {
	var id int64
	err := tx.QueryRow(
		"INSERT INTO employee (name, role_id) VALUES ($1, $2) RETURNING id",
		e.Name, e.RoleId).Scan(&id)
	if err != nil {
		return -1, err
	}
	return id, nil
}

func (r *Repository) FindById(id int64) (res Entity, err error) {
	err = r.db.Get(&res, "SELECT * FROM employee WHERE id = $1", id)
	return res, err
}

func (r *Repository) FindByName(tx *sqlx.Tx, name string) (isExist bool, err error) {
	err = tx.Get(
		&isExist,
		"SELECT EXISTS(SELECT 1 FROM employee WHERE name = $1)",
		name,
	)
	if err != nil {
		return false, err
	}
	return isExist, nil
}

func (r *Repository) FindAll() ([]Entity, error) {
	var employees []Entity
	rows, err := r.db.Queryx("SELECT * FROM employee")
	if err != nil {
		return employees, err
	}
	for rows.Next() {
		var e Entity
		if err := rows.StructScan(&e); err != nil {
			return employees, err
		}
		employees = append(employees, e)
	}
	return employees, nil
}

func (r *Repository) FindByIds(ids []int64) ([]Entity, error) {
	var employees []Entity
	query, args, err := sqlx.In("SELECT * FROM employee WHERE id IN (?)", ids)
	if err != nil {
		return employees, err
	}
	query = r.db.Rebind(query)
	rows, err := r.db.Queryx(query, args...)
	if err != nil {
		return employees, err
	}
	for rows.Next() {
		var e Entity
		if err := rows.StructScan(&e); err != nil {
			return employees, err
		}
		employees = append(employees, e)
	}
	return employees, nil
}

func (r *Repository) FindWithOffset(offset int, limit int, filter string) ([]Entity, error) {
	var employees []Entity
	query := "SELECT * FROM employee WHERE 1 = 1"
	var args []interface{}
	paramIdx := 1
	if filter != "" {
		query += fmt.Sprintf(" AND name ILIKE $%d", paramIdx)
		args = append(args, "%"+filter+"%")
		paramIdx++
	}
	query += fmt.Sprintf(" ORDER BY id OFFSET $%d LIMIT $%d", paramIdx, paramIdx+1)
	args = append(args, offset, limit)
	err := r.db.Select(&employees, query, args...)
	if err != nil {
		return employees, err
	}
	return employees, nil
}

func (r *Repository) DeleteById(id int64) error {
	_, err := r.db.Exec("DELETE FROM employee WHERE id = $1", id)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) DeleteByIds(ids []int64) error {
	query, args, err := sqlx.In("DELETE FROM employee WHERE id IN (?)", ids)
	if err != nil {
		return err
	}
	query = r.db.Rebind(query)
	_, err = r.db.Exec(query, args...)
	if err != nil {
		return err
	}
	return nil
}
