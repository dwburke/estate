package db

import (
	"strconv"

	"github.com/dwburke/go-tools"
	"github.com/spf13/viper"
	"github.com/syndtr/goleveldb/leveldb"
	leveldb_errors "github.com/syndtr/goleveldb/leveldb/errors"
)

var conn_leveldb *leveldb.DB

func init() {
	viper.SetDefault("db.leveldb.datadir", "state")
}

func OpenLevelDB() *leveldb.DB {
	if conn_leveldb != nil {
		return conn_leveldb
	}

	var err error
	conn_leveldb, err = leveldb.OpenFile(viper.GetString("db.leveldb.datadir"), nil)
	tools.FatalError(err)

	return conn_leveldb
}

func LevelDBClose() {
	if conn_leveldb != nil {
		conn_leveldb.Close()
	}
}

func LevelDBGet(key string) (value string) {
	ldb := OpenLevelDB()

	data, err := ldb.Get([]byte(key), nil)

	if err != nil {
		if err != leveldb_errors.ErrNotFound {
			tools.FatalError(err)
		}
	} else if data != nil {
		value = string(data)
	}

	return
}

func LevelDBGetInt(key string) (value int64) {
	value, _ = strconv.ParseInt(LevelDBGet(key), 10, 64)
	return
}

func LevelDBSet(key string, value string) {
	ldb := OpenLevelDB()

	err := ldb.Put([]byte(key), []byte(value), nil)
	tools.FatalError(err)
}

func LevelDBSetInt(key string, value int64) {
	LevelDBSet(key, strconv.FormatInt(value, 10))
}
