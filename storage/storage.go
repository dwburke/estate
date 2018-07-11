package storage

import (
	"database/sql"
	"errors"

	_ "github.com/go-sql-driver/mysql"
	"github.com/spf13/viper"
)

var (
	ErrNoRows = errors.New("storage: no rows returned")
)

type Storage struct {
	Handle *sql.DB
}

var err error

// create table prefs (var varchar(255), value text, primary key(var));

func (st *Storage) Close() {
	if st.Handle != nil {
		st.Handle.Close()
	}
}

func New() (*Storage, error) {
	db, err := sql.Open(viper.GetString("prefs.storage.type"), viper.GetString("prefs.storage.dsn"))
	if err != nil {
		return nil, err
	}

	st := Storage{Handle: db}

	return &st, nil
}

func (st *Storage) Set(key string, value string) error {
	db := st.Handle

	stmt, err := db.Prepare("replace into prefs (var, value) values (?, ?)")
	if err != nil {
		return err
	}

	_, err = stmt.Exec(key, value)
	if err != nil {
		return err
	}

	return nil
}

func (st *Storage) Get(key string) (string, error) {
	db := st.Handle

	var value string

	row := db.QueryRow("select value from prefs where var = ?", key)

	switch err := row.Scan(&value); err {
	case sql.ErrNoRows:
		return "", ErrNoRows
	case nil:
		return value, nil
	default:
		return "", err
	}
}

func (st *Storage) Delete(key string) error {
	db := st.Handle

	stmt, err := db.Prepare("delete from prefs where key = ?")
	if err != nil {
		return err
	}

	_, err = stmt.Exec(key)
	if err != nil {
		return err
	}

	return nil
}
