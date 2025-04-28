package models

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
	Role     string `json:"role"`
}

type UserCreateModel struct {
	Name           string `json:"name" example:"John Doe"`
	Phone          string `json:"phone"`
	Email          string `json:"email"`
	DeviceID       string `json:"device_id"`
	Password       string `json:"password"`
	ConfirmPassword string `json:"confirm_password"`
}
