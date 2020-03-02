package web

import (
	"errors"
	"github.com/andrewkandzuba/openexchange-go-data-gov/pkg/db"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

type NewsFeedRepositoryStub struct {}

func Test_GetArticles_Success(t *testing.T) {
	e := echo.New()

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	c := e.NewContext(req, rec)
	c.SetPath("/articles")

	web := NewService(&NewsFeedRepositoryStub{})

	var expectedJSON = `{"articles":[{"ID":0,"CreatedAt":"0001-01-01T00:00:00Z","UpdatedAt":"0001-01-01T00:00:00Z","DeletedAt":null,"externalId":"12","uuid":"595b7871-b652-4c06-80b9-6260db2cd4e6","type":"","label":"","created":0,"updated":0,"href":"","body":"","status":"","adminOfficials":[{"ID":0,"CreatedAt":"0001-01-01T00:00:00Z","UpdatedAt":"0001-01-01T00:00:00Z","DeletedAt":null,"externalId":"9","label":"James Bond","href":""},{"ID":0,"CreatedAt":"0001-01-01T00:00:00Z","UpdatedAt":"0001-01-01T00:00:00Z","DeletedAt":null,"externalId":"10","label":"Eric Clapton","href":""}]},{"ID":0,"CreatedAt":"0001-01-01T00:00:00Z","UpdatedAt":"0001-01-01T00:00:00Z","DeletedAt":null,"externalId":"13","uuid":"5228c838-5cda-4c6d-8aed-2c0c55a51b31","type":"","label":"","created":0,"updated":0,"href":"","body":"","status":"","adminOfficials":[{"ID":0,"CreatedAt":"0001-01-01T00:00:00Z","UpdatedAt":"0001-01-01T00:00:00Z","DeletedAt":null,"externalId":"9","label":"James Bond","href":""},{"ID":0,"CreatedAt":"0001-01-01T00:00:00Z","UpdatedAt":"0001-01-01T00:00:00Z","DeletedAt":null,"externalId":"13","label":"Jon Bonjovi","href":""}]}]}
`
	if assert.NoError(t, web.GetArticles(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		actualJSON := rec.Body.String()
		assert.Equal(t, expectedJSON, actualJSON)
	}
}

func Test_GetArticleById_Success(t *testing.T) {
	e := echo.New()

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	c := e.NewContext(req, rec)
	c.SetPath("/articles/:id")
	// @ToDo: Bug!!! https://github.com/labstack/echo/issues/1492
	c.SetParamNames("id")
	c.SetParamValues("1")

	web := NewService(&NewsFeedRepositoryStub{})

	var expectedJSON = `{"ID":0,"CreatedAt":"0001-01-01T00:00:00Z","UpdatedAt":"0001-01-01T00:00:00Z","DeletedAt":null,"externalId":"12","uuid":"595b7871-b652-4c06-80b9-6260db2cd4e6","type":"","label":"","created":0,"updated":0,"href":"","body":"","status":"","adminOfficials":[{"ID":0,"CreatedAt":"0001-01-01T00:00:00Z","UpdatedAt":"0001-01-01T00:00:00Z","DeletedAt":null,"externalId":"9","label":"James Bond","href":""},{"ID":0,"CreatedAt":"0001-01-01T00:00:00Z","UpdatedAt":"0001-01-01T00:00:00Z","DeletedAt":null,"externalId":"10","label":"Eric Clapton","href":""}]}
`
	if assert.NoError(t, web.GetArticleById(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		actualJSON := rec.Body.String()
		assert.Equal(t, expectedJSON, actualJSON)
	}
}

func Test_GetArticleById_Notfound_Failure(t *testing.T) {
	e := echo.New()

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	c := e.NewContext(req, rec)
	c.SetPath("/articles/:id")
	// @ToDo: Bug!!! https://github.com/labstack/echo/issues/1492
	c.SetParamNames("id")
	c.SetParamValues("2")

	web := NewService(&NewsFeedRepositoryStub{})

	if assert.NoError(t, web.GetArticleById(c)) {
		assert.Equal(t, http.StatusNotFound, rec.Code)
	}
}

func Test_GetArticleById_InternalError_Failure(t *testing.T) {
	e := echo.New()

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	c := e.NewContext(req, rec)
	c.SetPath("/articles/:id")
	// @ToDo: Bug!!! https://github.com/labstack/echo/issues/1492
	c.SetParamNames("id")
	c.SetParamValues("100")

	web := NewService(&NewsFeedRepositoryStub{})

	if assert.NoError(t, web.GetArticleById(c)) {
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	}
}

func Test_GetArticleById_BadRequest_Failure(t *testing.T) {
	e := echo.New()

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	c := e.NewContext(req, rec)
	c.SetPath("/articles/:id")
	c.SetParamNames("id")
	c.SetParamValues("notint")

	web := NewService(&NewsFeedRepositoryStub{})

	if assert.NoError(t, web.GetArticleById(c)) {
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	}
}

func (r *NewsFeedRepositoryStub) Insert(article *db.Article) error {
	return nil
}

func (r *NewsFeedRepositoryStub) FindAll() ([]db.Article, error) {
	articles := []db.Article{
		{
			ExternalId: "12",
			UUID:       "595b7871-b652-4c06-80b9-6260db2cd4e6",
			Type:       "",
			Label:      "",
			Created:    0,
			Updated:    0,
			Href:       "",
			Body:       "",
			Status:     "",
			AdminOfficials: []db.AdminOfficial{
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
		},
		{
			ExternalId: "13",
			UUID:       "5228c838-5cda-4c6d-8aed-2c0c55a51b31",
			Type:       "",
			Label:      "",
			Created:    0,
			Updated:    0,
			Href:       "",
			Body:       "",
			Status:     "",
			AdminOfficials: []db.AdminOfficial{
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
		},
	}
	return articles, nil
}


func (r *NewsFeedRepositoryStub) Find(id int) (*db.Article, error) {
	if id == 1 {
		return &db.Article {
		ExternalId: "12",
			UUID:       "595b7871-b652-4c06-80b9-6260db2cd4e6",
				Type:       "",
				Label:      "",
				Created:    0,
				Updated:    0,
				Href:       "",
				Body:       "",
				Status:     "",
				AdminOfficials: []db.AdminOfficial{
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
		}, nil
	} else if id == 2 {
		return nil, errors.New("record not found")
	}
	return nil, errors.New("unknown error")
}
