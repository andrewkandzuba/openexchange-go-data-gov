package ingress

import (
	"errors"
	"github.com/andrewkandzuba/openexchange-go-data-gov/pkg/db"
	"github.com/buger/jsonparser"
	"gopkg.in/validator.v2"
	"time"
)

type NewsFeedConsumer struct {
	Repo *db.NewsFeedRepository `validate:"nonzero"`
}

func NewNewsFeedConsumer(repo *db.NewsFeedRepository) (*NewsFeedConsumer, error) {
	instance := &NewsFeedConsumer{repo}

	if errs := validator.Validate(instance); errs != nil {
		// ToDo: Create a test to handle log.Fatal(...)
		return nil, errors.New(errs.Error())
	}

	return instance, nil
}

func (c *NewsFeedConsumer) From(ch chan string) error {
	var ex error
RangeLoop:
	for {
		select {
		case v, err := <-ch:
			if !err {
				break RangeLoop
			}
			ex := c.save(v)
			if ex != nil {
				break RangeLoop
			}
		case <-time.After(3 * time.Second):
			break RangeLoop
		}
	}
	return ex
}

func (c *NewsFeedConsumer) save(json string) error {
	var err error

	a := &db.Article{}

	a.Type, err = jsonparser.GetString([]byte(json), "type")
	a.ExternalId, err = jsonparser.GetString([]byte(json),"id")
	a.UUID, err = jsonparser.GetString([]byte(json), "uuid")
	a.Label, err = jsonparser.GetString([]byte(json),"label")
	a.Body, err = jsonparser.GetString([]byte(json), "body")
	a.Created, err = jsonparser.GetInt([]byte(json),"created")
	a.Updated, err = jsonparser.GetInt([]byte(json), "updated")
	a.Href, err = jsonparser.GetString([]byte(json),"href")
	a.Status, err = jsonparser.GetString([]byte(json), "release_status")

	a.AdminOfficials = make([]db.AdminOfficial, 0)

	_, err = jsonparser.ArrayEach([]byte(json), func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		ao := db.AdminOfficial{}
		ao.ExternalId, err = jsonparser.GetString(value, "id")
		ao.Label, err = jsonparser.GetString(value, "label")
		ao.Href, err = jsonparser.GetString(value,"href")
		a.AdminOfficials = append(a.AdminOfficials, ao)
	}, "admin_officials")

	err = c.Repo.Insert(a)

	return err
}
