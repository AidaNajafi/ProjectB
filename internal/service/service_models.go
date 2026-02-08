package service

type SignUpRequest struct {
	ID       int
	Name     string
	Username string
	Email    string
	Password string
}

type LoginRequest struct {
	Username string
	Password string
}
