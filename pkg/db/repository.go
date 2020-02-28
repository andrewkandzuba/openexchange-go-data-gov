package db

import (
	"errors"
	"github.com/jinzhu/gorm"
	"gopkg.in/validator.v2"
)

type Article struct {
	gorm.Model
	ExternalId     string
	UUID           string `gorm:"unique_index:idx_article_uuid"`
	Type           string
	Label          string
	Created        int64
	Updated        int64
	Href           string
	Body           string
	Status         string
	AdminOfficials []AdminOfficial `gorm:"many2many:article_admin_officials;"`
}

type AdminOfficial struct {
	gorm.Model
	ExternalId string `gorm:"unique_index:idx_admin_official_uuid"`
	Label   string
	Href    string
	Articles []Article `gorm:"many2many:article_admin_officials;"`
}

type NewsFeedRepository struct {
	Db *gorm.DB `validate:"nonzero"`
}

func NewNewsFeedRepository(db *gorm.DB) (*NewsFeedRepository, error) {
	instance := &NewsFeedRepository{Db:db}

	if errs := validator.Validate(instance); errs != nil {
		// ToDo: Create a test to handle log.Fatal(...)
		return nil, errors.New(errs.Error())
	}

	return instance, nil
}

func (r *NewsFeedRepository) Insert(article *Article) error {
	for i := 0; i< len(article.AdminOfficials); i++ {
		_ = r.Db.Where("external_id=?", article.AdminOfficials[i] .ExternalId).Find(&article.AdminOfficials[i])
	}
	res := r.Db.Create(&article)
	return res.Error
}