package user

type UserRequest struct {
	Name string `json:"name" validate:"required"`
}
