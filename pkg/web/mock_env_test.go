package web

import (
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	code := m.Run()
	os.Exit(code)
}
