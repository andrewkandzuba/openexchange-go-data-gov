package main

import (
	"fmt"
	"github.com/andrewkandzuba/openexchange-go-data-gov/pkg/config"
	db2 "github.com/andrewkandzuba/openexchange-go-data-gov/pkg/db"
	"github.com/andrewkandzuba/openexchange-go-data-gov/pkg/stream"
	"github.com/andrewkandzuba/openexchange-go-data-gov/pkg/ingress"
	"github.com/andrewkandzuba/openexchange-go-data-gov/pkg/connector"
	"github.com/andrewkandzuba/openexchange-go-data-gov/pkg/web"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
	"os"
)

func main() {
	log.Println("Data.gov Commerce News Web Service")

	// @ToDo: Switch to MySQL and move the environment's initialization into Terraform or Google Cloud Deployment Manager.
	// 1. Spin up environment
	dialect := "sqlite3"
	dbHost := os.TempDir() + "/test.db"
	db, err := gorm.Open(dialect, dbHost)
	if err != nil {
		panic(err)
	}
	defer func() {
		db.DropTableIfExists(&db2.Article{}, &db2.AdminOfficial{})
		db.Close()
		os.Remove(dbHost)
	}()
	db.AutoMigrate(&db2.Article{}, &db2.AdminOfficial{})

	_ = os.Setenv("DB_HOST", dbHost)
	_ = os.Setenv("DB_DIALECT", dialect)

	// @ToDo: Externalise configuration into some configuration manager.
	// 2. Load configuration
	workdir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(workdir)
	cfg := config.NewConfig(workdir + "/cmd/application/application.yaml")

	// 3. Start stream and ingress services
	go func() {
		api, err := connector.NewConnector(cfg.Api.Endpoint, cfg.Api.Key)
		if err != nil {
			log.Fatal(err)
		}

		cs, err := stream.NewCommerceStream(api)
		if err != nil {
			log.Fatal(err)
		}

		repo, err := db2.NewArticleRepository(db)
		if err != nil {
			log.Fatal(err)
		}

		c, err := ingress.NewArticleFeedConsumer(repo)
		if err != nil {
			log.Fatal(err)
		}

		err = c.From(cs.Stream())
		if err != nil {
			log.Fatal(err)
		}
	}()

	// 4. Initialize and Start webservice
	repo, err := db2.NewArticleRepository(db)
	if err != nil {
		log.Fatal(err)
	}
	ws := web.NewService(repo)

	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	e.GET("/articles", ws.GetArticles)
	e.GET("/articles/:id", ws.GetArticleById)

	e.Logger.Fatal(e.Start(cfg.Web.Address))
}
