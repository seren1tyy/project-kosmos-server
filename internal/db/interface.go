package db

type User struct {
	ID       int    `json:"id"`
	Login    string `json:"login"`
	Password string `json:"password"`
}

type DB interface {
	Connect(dsn string) error
	FindByLogin(login string) (*User, bool, error)
	HasCharacters(accountID int) (bool, error)
	Close() error
}
