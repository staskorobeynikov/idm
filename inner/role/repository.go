package role

import (
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

func (r *Repository) Save(e Entity) (int64, error) {
	var id int64
	err := r.db.QueryRow(
		"INSERT INTO role (name) VALUES ($1) RETURNING id",
		e.Name).Scan(&id)
	if err != nil {
		return -1, err
	}
	return id, nil
}

func (r *Repository) FindById(id int64) (res Entity, err error) {
	err = r.db.Get(&res, "SELECT * FROM role WHERE id = $1", id)
	return res, err
}

func (r *Repository) FindAll() ([]Entity, error) {
	var roles []Entity
	rows, err := r.db.Queryx("SELECT * FROM role")
	if err != nil {
		return roles, err
	}
	for rows.Next() {
		var role Entity
		if err := rows.StructScan(&role); err != nil {
			return roles, err
		}
		roles = append(roles, role)
	}
	return roles, nil
}

func (r *Repository) FindByIds(ids []int64) ([]Entity, error) {
	var roles []Entity
	query, args, err := sqlx.In("SELECT * FROM role WHERE id IN (?)", ids)
	if err != nil {
		return roles, err
	}
	query = r.db.Rebind(query)
	rows, err := r.db.Queryx(query, args...)
	if err != nil {
		return roles, err
	}
	for rows.Next() {
		var role Entity
		if err := rows.StructScan(&role); err != nil {
			return roles, err
		}
		roles = append(roles, role)
	}
	return roles, nil
}

func (r *Repository) DeleteById(id int64) error {
	_, err := r.db.Exec("DELETE FROM role WHERE id = $1", id)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) DeleteByIds(ids []int64) error {
	query, args, err := sqlx.In("DELETE FROM role WHERE id IN (?)", ids)
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
