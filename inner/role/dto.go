package role

import "time"

type Entity struct {
	Id        int64     `db:"id"`
	Name      string    `db:"name"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type Response struct {
	Id        int64     `db:"id"`
	Name      string    `db:"name"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

func (e *Entity) toResponse() Response {
	return Response{
		Id:        e.Id,
		Name:      e.Name,
		CreatedAt: e.CreatedAt,
		UpdatedAt: e.UpdatedAt,
	}
}

type CreateRequest struct {
	Name string `json:"name" validate:"required,min=2,max=155"`
}

func (req *CreateRequest) ToEntity() Entity {
	return Entity{
		Name: req.Name,
	}
}

type IdRequest struct {
	Id int64 `json:"id" validate:"required,min=1"`
}

type IdsRequest struct {
	Ids []int64 `json:"ids" validate:"required,min=1,dive"`
}
