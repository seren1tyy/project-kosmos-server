package db

import (
	"database/sql"
	"fmt"
)

type MariaDB struct {
	conn *sql.DB
}

func NewMariaDB() DB { return &MariaDB{} }

func (d *MariaDB) Connect(dsn string) error {
	conn, err := sql.Open("mysql", dsn)
	if err != nil {
		return fmt.Errorf("sql.Open: %w", err)
	}
	if err := conn.Ping(); err != nil {
		return fmt.Errorf("ping db: %w", err)
	}
	conn.SetMaxOpenConns(25)
	conn.SetMaxIdleConns(5)
	conn.SetConnMaxLifetime(5 * 60) // 5 мин
	d.conn = conn
	return nil
}

func (d *MariaDB) FindByLogin(login string) (*User, bool, error) {
	var u User
	err := d.conn.QueryRow("SELECT id, login, password FROM accounts WHERE login = ?", login).
		Scan(&u.ID, &u.Login, &u.Password)
	if err == sql.ErrNoRows {
		return nil, false, nil
	}
	if err != nil {
		return nil, false, err
	}
	return &u, true, nil
}

func (d *MariaDB) HasCharacters(accountID int) (bool, error) {
	var count int
	err := d.conn.QueryRow("SELECT COUNT(*) FROM characters WHERE account_id = ?", accountID).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (d *MariaDB) Close() error {
	if d.conn != nil {
		return d.conn.Close()
	}
	return nil
}

func (d *MariaDB) GetDB() *sql.DB {
	return d.conn
}
