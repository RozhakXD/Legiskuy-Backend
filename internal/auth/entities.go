package auth

type User struct {
	ID       int
	Name     string
	Username string
	Password string
	Role     string
	HasVoted bool
}

type RegisterInput struct {
	Name     string `json:"name"`
	Username string `json:"username"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

type LoginInput struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserResponse struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Username string `json:"username"`
	Role     string `json:"role"`
}
