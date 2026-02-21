package service

type SignUpRequest struct {
	ID        int
	Name      string
	Username  string
	Email     string
	UserPhone string
	Password  string
}

type LoginRequest struct {
	Username  string
	UserPhone string
	Password  string
}


