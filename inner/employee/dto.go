package employee

import "time"

type Entity struct {
	Id        int64     `db:"id"`
	Name      string    `db:"name"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
	RoleId    int64     `db:"role_id"`
}

type Response struct {
	Id        int64     `db:"id"`
	Name      string    `db:"name"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type CreateRequest struct {
	Name   string `json:"name" validate:"required,min=2,max=155"`
	RoleId int64  `json:"role_id" validate:"required,min=1"`
}

type IdRequest struct {
	Id int64 `json:"id" validate:"required,min=1"`
}

type IdsRequest struct {
	Ids []int64 `json:"ids" validate:"required,min=1,dive"`
}

type PageRequest struct {
	PageSize   int `validate:"min=1,max=100"`
	PageNumber int `validate:"min=0"`
}

type PageResponse struct {
	Result     []Response
	PageSize   int
	PageNumber int
	Total      int64
}

func (e *Entity) toResponse() Response {
	return Response{
		Id:        e.Id,
		Name:      e.Name,
		CreatedAt: e.CreatedAt,
		UpdatedAt: e.UpdatedAt,
	}
}

func (req *CreateRequest) ToEntity() Entity {
	return Entity{
		Name:   req.Name,
		RoleId: req.RoleId,
	}
}
