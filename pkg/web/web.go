package web

import (
	"github.com/andrewkandzuba/openexchange-go-data-gov/pkg/db"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

type (
	resultLists struct {
		Articles []db.Article `json:"articles"`
	}

	Service struct {
		nfr db.ArticleRepositoryImpl
	}
)

func NewService(nfr db.ArticleRepositoryImpl) *Service {
	return &Service{nfr: nfr}
}

func (w *Service) GetArticles(c echo.Context) error {
	lists, err := w.nfr.FindAll()
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	a := &resultLists{
		Articles: lists,
	}
	return c.JSON(http.StatusOK, a)
}


func (w *Service) GetArticleById(c echo.Context) error {
	id := c.Param("id")
	idx, err := strconv.Atoi(id)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	a, err := w.nfr.Find(idx)
	if err != nil {
		if err.Error() == "record not found" {
			return c.String(http.StatusNotFound, err.Error())
		} else {
			return c.String(http.StatusInternalServerError, err.Error())
		}
	}
	return c.JSON(http.StatusOK, a)
}
