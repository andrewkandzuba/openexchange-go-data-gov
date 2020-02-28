package ingress

import (
	"bufio"
	"github.com/andrewkandzuba/openexchange-go-data-gov/pkg/config"
	db2 "github.com/andrewkandzuba/openexchange-go-data-gov/pkg/db"
	"github.com/buger/jsonparser"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
	"log"
	"os"
	"testing"
)

func Test_NewInstance_Success(t *testing.T) {
	repo := db2.NewsFeedRepository{
		&gorm.DB{},
	}

	sink, err := NewNewsFeedConsumer(&repo)

	assert.NotNil(t, sink)
	assert.Nil(t, err)
}

func Test_NewInstance_Failure(t *testing.T) {
	repo := db2.NewsFeedRepository{
		nil,
	}

	sink, err := NewNewsFeedConsumer(&repo)

	assert.Nil(t, sink)
	assert.NotNil(t, err)
	assert.Equal(t, "Repo.Db: zero value", err.Error())
}

func Test_ConsumeFromStream_Success(t *testing.T) {
	ch := make(chan string)
	go func() {
		defer close(ch)

		handle, err := os.Open("testdata/news.json")
		if err != nil {
			log.Fatal(err)
		}
		defer handle.Close()

		var json string
		scanner := bufio.NewScanner(handle)
		for scanner.Scan() {
			json += scanner.Text()
		}

		_, err = jsonparser.ArrayEach([]byte(json), func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
			ch <- string(value)
		}, "data")
		if err != nil {
			log.Fatal(err)
		}

	}()

	cfg := config.NewConfig("testdata/application-test.yaml")

	db, err := gorm.Open(cfg.Db.Dialect, cfg.Db.Host)
	if err != nil {
		panic(err)
	}
	defer func() {
		db.DropTableIfExists(&db2.Article{}, &db2.AdminOfficial{})
	}()
	db.AutoMigrate(&db2.Article{}, &db2.AdminOfficial{})

	repo, _ := db2.NewNewsFeedRepository(db)
	consumer, _ := NewNewsFeedConsumer(repo)

	err = consumer.From(ch)

	assert.Nil(t, err)

	var articles []db2.Article
	db.Find(&articles)
	assert.Equal(t, 3, len(articles))

	var adminOfficials []db2.AdminOfficial
	db.Find(&adminOfficials)
	assert.Equal(t, 2, len(adminOfficials))
}
