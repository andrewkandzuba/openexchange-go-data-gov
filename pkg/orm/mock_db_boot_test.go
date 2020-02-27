package orm

import (
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
	"os"
	"testing"
)

var db gorm.DB

func TestMain(m *testing.M) {
	db, err := gorm.Open("sqlite3", "test.db")
	if err != nil {
		panic(err)
	}

	code := m.Run()
	defer db.Close()
	os.Exit(code)
}

