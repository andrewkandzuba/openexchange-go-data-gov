package db

import (
	"errors"
	"github.com/jinzhu/gorm"
	"gopkg.in/validator.v2"
)

type (
	ArticleRepositoryImpl interface {
		Insert(article *Article) error
		FindAll() ([]Article, error)
		Find(id int) (*Article, error)
	}

	Article struct {
		gorm.Model
		ExternalId     string          `json:"externalId"`
		UUID           string          `json:"uuid" gorm:"unique_index:idx_article_uuid"`
		Type           string          `json:"type"`
		Label          string          `json:"label"`
		Created        int64           `json:"created"`
		Updated        int64           `json:"updated"`
		Href           string          `json:"href"`
		Body           string          `json:"body"`
		Status         string          `json:"status"`
		AdminOfficials []AdminOfficial `json:"adminOfficials" gorm:"many2many:article_admin_officials;"`
	}

	AdminOfficial struct {
		gorm.Model
		ExternalId string    `json:"externalId" gorm:"unique_index:idx_admin_official_uuid"`
		Label      string    `json:"label"`
		Href       string    `json:"href"`
		Articles   []Article `json:"-" gorm:"many2many:article_admin_officials;"`
	}

	ArticleRepository struct {
		Db *gorm.DB `validate:"nonzero"`
	}
)

func NewArticleRepository(db *gorm.DB) (*ArticleRepository, error) {
	instance := &ArticleRepository{Db: db}

	if errs := validator.Validate(instance); errs != nil {
		// ToDo: Create a test to handle log.Fatal(...)
		return nil, errors.New(errs.Error())
	}

	return instance, nil
}

func (r *ArticleRepository) Insert(article *Article) error {
	for i := 0; i < len(article.AdminOfficials); i++ {
		_ = r.Db.Where("external_id=?", article.AdminOfficials[i].ExternalId).Find(&article.AdminOfficials[i])
	}
	res := r.Db.Create(&article)
	return res.Error
}

func (r *ArticleRepository) FindAll() ([]Article, error) {
	var articles []Article
	res := r.Db.Find(&articles)
	return articles, res.Error
}

func (r *ArticleRepository) Find(id int) (*Article, error) {
	var article Article
	res := r.Db.First(&article, id)
	if res.Error != nil {
		return nil, res.Error
	}
	return &article, nil
}
