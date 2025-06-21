package employee

import (
	"github.com/jmoiron/sqlx"
	"time"
)

type EmployeeRepository struct {
	db *sqlx.DB
}

func NewEmployeeRepository(database *sqlx.DB) *EmployeeRepository {
	return &EmployeeRepository{db: database}
}

type EmployeeEntity struct {
	Id        int       `db:"id"`
	Name      string    `db:"name"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
	roleId    int       `db:"role_id"`
}

func (r *EmployeeRepository) Add(e EmployeeEntity) (int64, error) {
	var id int64
	err := r.db.QueryRow(
		"INSERT INTO employee (name, role_id) VALUES ($1, $2) RETURNING id",
		e.Name, e.roleId).Scan(&id)
	if err != nil {
		return -1, err
	}
	return id, nil
}

func (r *EmployeeRepository) GetById(id int64) (res EmployeeEntity, err error) {
	err = r.db.Get(&res, "SELECT * FROM employee WHERE id = $1", id)
	return res, err
}

func (r *EmployeeRepository) FindAll() ([]EmployeeEntity, error) {
	var employees []EmployeeEntity
	rows, err := r.db.Queryx("SELECT * FROM employee")
	if err != nil {
		return employees, err
	}
	for rows.Next() {
		var e EmployeeEntity
		if err := rows.StructScan(&e); err != nil {
			return employees, err
		}
		employees = append(employees, e)
	}
	return employees, nil
}

func (r *EmployeeRepository) FindByIds(ids []int64) ([]EmployeeEntity, error) {
	var employees []EmployeeEntity
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
		var e EmployeeEntity
		if err := rows.StructScan(&e); err != nil {
			return employees, err
		}
		employees = append(employees, e)
	}
	return employees, nil
}

func (r *EmployeeRepository) DeleteByUId(id int64) error {
	_, err := r.db.Exec("DELETE FROM employee WHERE id = $1", id)
	if err != nil {
		return err
	}
	return nil
}

func (r *EmployeeRepository) DeleteByIds(ids []int64) error {
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
