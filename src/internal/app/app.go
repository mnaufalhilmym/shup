package app

import (
	"database/sql"
	"net/http"

	"github.com/rs/zerolog/log"
	"shup.hilmy.dev/src/internal/controller"
	"shup.hilmy.dev/src/internal/repository"
	"shup.hilmy.dev/src/internal/router"
	"shup.hilmy.dev/src/internal/scheduler"
	"shup.hilmy.dev/src/internal/service"
)

type BootstrapParams struct {
	Mux        *http.ServeMux
	DB         *sql.DB
	WebAddress string
}

func Bootstrap(params BootstrapParams) {
	fileRepository := repository.NewFile(params.DB)

	fileService := service.NewFile(params.DB, fileRepository)

	fileController := controller.NewFile(fileService)

	router := router.New(params.Mux, fileController)
	handler := router.Configure()

	go scheduler.FileDeletionScheduler(fileService)

	log.Debug().Msg("Listening and serving HTTP on " + params.WebAddress)
	if err := http.ListenAndServe(params.WebAddress, handler); err != nil {
		log.Error().Err(err).Msg("Failed to start web server")
	}
}
