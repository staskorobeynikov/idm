package role

import (
	"github.com/jmoiron/sqlx"
	"time"
)

type RoleRepository struct {
	db *sqlx.DB
}

func NewRoleRepository(database *sqlx.DB) *RoleRepository {
	return &RoleRepository{db: database}
}

type RoleEntity struct {
	Id        int       `db:"id"`
	Name      string    `db:"name"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

func (r *RoleRepository) Add(e RoleEntity) (int64, error) {
	var id int64
	err := r.db.QueryRow(
		"INSERT INTO role (name) VALUES ($1) RETURNING id",
		e.Name).Scan(&id)
	if err != nil {
		return -1, err
	}
	return id, nil
}

func (r *RoleRepository) GetById(id int64) (res RoleEntity, err error) {
	err = r.db.Get(&res, "SELECT * FROM role WHERE id = $1", id)
	return res, err
}

func (r *RoleRepository) FindAll() ([]RoleEntity, error) {
	var roles []RoleEntity
	rows, err := r.db.Queryx("SELECT * FROM role")
	if err != nil {
		return roles, err
	}
	for rows.Next() {
		var role RoleEntity
		if err := rows.StructScan(&role); err != nil {
			return roles, err
		}
		roles = append(roles, role)
	}
	return roles, nil
}

func (r *RoleRepository) FindByIds(ids []int64) ([]RoleEntity, error) {
	var roles []RoleEntity
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
		var role RoleEntity
		if err := rows.StructScan(&role); err != nil {
			return roles, err
		}
		roles = append(roles, role)
	}
	return roles, nil
}

func (r *RoleRepository) DeleteByUId(id int64) error {
	_, err := r.db.Exec("DELETE FROM role WHERE id = $1", id)
	if err != nil {
		return err
	}
	return nil
}

func (r *RoleRepository) DeleteByIds(ids []int64) error {
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
