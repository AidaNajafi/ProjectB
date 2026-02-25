package service

type SignUpRequest struct {
	ID        int64
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


