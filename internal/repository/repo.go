package repository

type UserInfo struct {
	ID        int
	Name      string
	Username  string
	Email     string
	UserPhone string
	Password  string
}

type ReservationInfo struct {
	ID        int
	RoomID    int
	HotelID   int
	UserID    int
	UserPhone string
	DateFrom  string
	DateTo    string
	Status    string
}

type Store interface {
	CreateUser(user UserInfo) error
	GetUserByUsername(username string) (*UserInfo, error)
	CreateReservation(r ReservationInfo) error
}
