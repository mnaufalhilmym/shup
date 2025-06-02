package main

import (
	"net/http"

	"shup.hilmy.dev/src/internal/app"
	"shup.hilmy.dev/src/internal/config"
)

func main() {
	conf := config.NewConfiguration()

	config.ConfigureLogger(conf.LogLevel())

	mux := http.NewServeMux()

	db := config.NewDatabase(conf.DBDriverName(), conf.DBConnectionURL())
	defer db.Close()

	app.Bootstrap(app.BootstrapParams{
		Mux:        mux,
		DB:         db,
		WebAddress: conf.WebAddress(),
	})
}
