package db

// Note: upsert is not supported by GORM: https://github.com/jinzhu/gorm/issues/1188

import (
	"github.com/andrewkandzuba/openexchange-go-data-gov/pkg/config"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewNewsFeedRepository_NewInstance_Success(t *testing.T) {
	cfg := config.NewConfig("testdata/application-test.yaml")

	db, err := gorm.Open(cfg.Db.Dialect, cfg.Db.Host)
	if err != nil {
		panic(err)
	}

	repo, err := NewNewsFeedRepository(db)
	assert.Nil(t, err)
	assert.NotNil(t, repo)
}

func TestNewsFeedRepository_NewInstance_InvalidArgument_Failure(t *testing.T) {
	repo, err := NewNewsFeedRepository(nil)

	assert.NotNil(t, err)
	assert.Equal(t, "Db: zero value", err.Error())
	assert.Nil(t, repo)
}

func TestNewsFeedRepository_InsertArticle_Success(t *testing.T) {
	cfg := config.NewConfig("testdata/application-test.yaml")

	db, err := gorm.Open(cfg.Db.Dialect, cfg.Db.Host)
	if err != nil {
		panic(err)
	}
	defer func() {
		db.DropTableIfExists(&Article{}, &AdminOfficial{})
	}()
	db.AutoMigrate(&Article{}, &AdminOfficial{})

	repo, _ := NewNewsFeedRepository(db)

    err = repo.Insert(&Article{
		ExternalId: "12",
		UUID:       "595b7871-b652-4c06-80b9-6260db2cd4e6",
		Type:       "",
		Label:      "",
		Created:    0,
		Updated:    0,
		Href:       "",
		Body:       "",
		Status:     "",
		AdminOfficials: []AdminOfficial{
			{
				ExternalId: "9",
				Label:      "James Bond",
				Href:       "",
			},
			{
				ExternalId: "10",
				Label:      "Eric Clapton",
				Href:       "",
			},
		},
	})

    assert.Nil(t, err)

	var article Article
	db.First(&article, 1)

	assert.NotNil(t, article)
	assert.Equal(t, "595b7871-b652-4c06-80b9-6260db2cd4e6", article.UUID)
	assert.Equal(t, "12", article.ExternalId)
	assert.Equal(t, 0, len(article.AdminOfficials))

	var adminOfficials []AdminOfficial
	db.Find(&adminOfficials)

	assert.NotNil(t, adminOfficials)
	assert.Equal(t, 2, len(adminOfficials))
	assert.Equal(t, "9", adminOfficials[0].ExternalId)
	assert.Equal(t, "James Bond", adminOfficials[0].Label)
	assert.Equal(t, "10", adminOfficials[1].ExternalId)
	assert.Equal(t, "Eric Clapton", adminOfficials[1].Label)

	db.Preload("AdminOfficials").First(&article, 1)

	assert.NotNil(t, article)
	assert.Equal(t, "595b7871-b652-4c06-80b9-6260db2cd4e6", article.UUID)
	assert.Equal(t, "12", article.ExternalId)
	assert.Equal(t, 2, len(article.AdminOfficials))
	assert.Equal(t, "9", article.AdminOfficials[0].ExternalId)
	assert.Equal(t, "James Bond", article.AdminOfficials[0].Label)
	assert.Equal(t, "10", article.AdminOfficials[1].ExternalId)
	assert.Equal(t, "Eric Clapton", article.AdminOfficials[1].Label)

	var adminOfficial AdminOfficial
	db.Preload("Articles").Where("external_id=?", "9").Find(&adminOfficial)

	assert.NotNil(t, adminOfficial)
	assert.Equal(t, "James Bond", adminOfficial.Label)
	assert.Equal(t, 1, len(adminOfficial.Articles))
	assert.Equal(t, "595b7871-b652-4c06-80b9-6260db2cd4e6", adminOfficial.Articles[0].UUID)
}

func TestNewsFeedRepository_InsertDuplicateArticle_Failure(t *testing.T) {
	cfg := config.NewConfig("testdata/application-test.yaml")

	db, err := gorm.Open(cfg.Db.Dialect, cfg.Db.Host)
	if err != nil {
		panic(err)
	}
	defer func() {
		db.DropTableIfExists(&Article{}, &AdminOfficial{})
	}()
	db.AutoMigrate(&Article{}, &AdminOfficial{})

	repo, _ := NewNewsFeedRepository(db)

	article := &Article{
		ExternalId: "12",
		UUID:       "595b7871-b652-4c06-80b9-6260db2cd4e6",
		Type:       "",
		Label:      "",
		Created:    0,
		Updated:    0,
		Href:       "",
		Body:       "",
		Status:     "",
		AdminOfficials: []AdminOfficial{
			{
				ExternalId: "9",
				Label:      "James Bond",
				Href:       "",
			},
			{
				ExternalId: "10",
				Label:      "Eric Clapton",
				Href:       "",
			},
		},
	}

	var articles []Article

	err = repo.Insert(article)

	assert.Nil(t, err)
	db.Find(&articles)
	assert.Equal(t, 1, len(articles))

	err = repo.Insert(article)

	assert.NotNil(t, err)
	assert.Equal(t, "UNIQUE constraint failed: articles.id", err.Error())
	db.Find(&articles)
	assert.Equal(t, 1, len(articles))
}

func TestNewsFeedRepository_PersistAssociations_Success(t *testing.T) {
	cfg := config.NewConfig("testdata/application-test.yaml")

	db, err := gorm.Open(cfg.Db.Dialect, cfg.Db.Host)
	if err != nil {
		panic(err)
	}
	defer func() {
		db.DropTableIfExists(&Article{}, &AdminOfficial{})
	}()
	db.AutoMigrate(&Article{}, &AdminOfficial{})

	repo, _ := NewNewsFeedRepository(db)

	err = repo.Insert(&Article{
		ExternalId: "12",
		UUID:       "595b7871-b652-4c06-80b9-6260db2cd4e6",
		Type:       "",
		Label:      "",
		Created:    0,
		Updated:    0,
		Href:       "",
		Body:       "",
		Status:     "",
		AdminOfficials: []AdminOfficial{
			{
				ExternalId: "9",
				Label:      "James Bond",
				Href:       "",
			},
			{
				ExternalId: "10",
				Label:      "Eric Clapton",
				Href:       "",
			},
		},
	})

	assert.Nil(t, err)

	var articles []Article
	db.Find(&articles)
	assert.Equal(t, 1, len(articles))

	var adminOfficials []AdminOfficial
	db.Find(&adminOfficials)
	assert.Equal(t, 2, len(adminOfficials))

	err = repo.Insert(&Article{
		ExternalId: "13",
		UUID:       "5228c838-5cda-4c6d-8aed-2c0c55a51b31",
		Type:       "",
		Label:      "",
		Created:    0,
		Updated:    0,
		Href:       "",
		Body:       "",
		Status:     "",
		AdminOfficials: []AdminOfficial{
			{
				ExternalId: "9",
				Label:      "James Bond",
				Href:       "",
			},
			{
				ExternalId: "13",
				Label:      "Jon Bonjovi",
				Href:       "",
			},
		},
	})

	db.Find(&articles)
	assert.Equal(t, 2, len(articles))

	db.Find(&adminOfficials)
	assert.Equal(t, 3, len(adminOfficials))

	var adminOfficial AdminOfficial
	db.Preload("Articles").Where("external_id=?", "9").Find(&adminOfficial)

	assert.Equal(t, 2, len(adminOfficial.Articles))
}




