package db

import (
	"fmt"

	"github.com/dwburke/go-tools"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/spf13/viper"
)

var gorm_db *gorm.DB

func init() {
	viper.SetDefault("db.postgres.maxidleconnections", 5)
	viper.SetDefault("db.postgres.maxopenconnections", 50)

	viper.SetDefault("db.postgres.host", "localhost")
	viper.SetDefault("db.postgres.port", 5432)
	viper.SetDefault("db.postgres.name", "app")
	viper.SetDefault("db.postgres.user", "user")
	viper.SetDefault("db.postgres.pass", "")
}

func OpenPostgres() *gorm.DB {
	var err error

	connStr := PgConnectString()

	gorm_db, err = gorm.Open("postgres", connStr)
	tools.FatalError(err)

	gorm_db = gorm_db.Set("gorm:auto_preload", true)

	gorm_db.DB().SetMaxIdleConns(viper.GetInt("db.postgres.maxidleconnections"))
	gorm_db.DB().SetMaxOpenConns(viper.GetInt("db.postgres.maxopenconnections"))

	return gorm_db
}

func PgConnectString() string {
	connStr := fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s sslmode=disable",
		viper.GetString("db.postgres.host"),
		viper.GetInt("db.postgres.port"),
		viper.GetString("db.postgres.user"),
		viper.GetString("db.postgres.name"),
		viper.GetString("db.postgres.pass"),
	)

	return connStr
}
