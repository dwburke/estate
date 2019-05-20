package mysql

import (
	"database/sql"

	"github.com/dwburke/estate/storage/common"
	"github.com/dwburke/estate/storage/meta"

	_ "github.com/go-sql-driver/mysql"
	"github.com/spf13/viper"
)

type Storage struct {
	meta.Storage
	Handle    *sql.DB
	tableName string
}

// create table estate (var varchar(255), value text, primary key(var));

func New() (*Storage, error) {
	db, err := sql.Open("mysql", viper.GetString("estate.storage.dsn"))
	if err != nil {
		return nil, err
	}

	st := Storage{
		Handle:    db,
		tableName: viper.GetString("estate.storage.table"),
	}

	return &st, nil
}

func (st *Storage) Set(key string, value []byte) error {
	db := st.Handle

	stmt, err := db.Prepare("replace into " + st.tableName + " (var, value) values (?, ?)")
	if err != nil {
		return err
	}

	_, err = stmt.Exec(key, value)
	if err != nil {
		return err
	}

	return nil
}

func (st *Storage) Get(key string) ([]byte, error) {
	db := st.Handle

	var value string

	row := db.QueryRow("select value from "+st.tableName+" where var = ?", key)

	switch err := row.Scan(&value); err {
	case sql.ErrNoRows:
		return nil, common.ErrNotFound
	case nil:
		return []byte(value), nil
	default:
		return nil, err
	}
}

func (st *Storage) Delete(key string) error {
	db := st.Handle

	stmt, err := db.Prepare("delete from " + st.tableName + " where key = ?")
	if err != nil {
		return err
	}

	_, err = stmt.Exec(key)
	if err != nil {
		return err
	}

	return nil
}

func (st *Storage) Close() error {
	if st.Handle != nil {
		st.Handle.Close()
	}
	return nil
}
