package db

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"os"
	"testing"
)

func TestMain(m *testing.M){
	dialect := "sqlite3"
	dbHost := os.TempDir() + "/test.db"
	db, err := gorm.Open(dialect, dbHost)
	if err != nil {
		panic(err)
	}
	defer func() {
		db.Close()
		os.Remove(dbHost)
	}()

	_ = os.Setenv("DB_HOST", dbHost)
	_ = os.Setenv("DB_DIALECT", dialect)

	code := m.Run()
	defer db.Close()
	os.Exit(code)
}
