package main

import (
	"log"

	app "github.com/AlyonaAg/bot-detector/internal/app/bot-detector"
	scriptsdb "github.com/AlyonaAg/bot-detector/internal/db/bot-detector"
	auth "github.com/AlyonaAg/bot-detector/internal/service"
)

func main() {
	repo, err := scriptsdb.NewRepository()
	if err != nil {
		log.Fatalf("repo error: %v", err)
	}

	authSrv, err := auth.NewAuthorizer()

	s := app.NewCheckerServer(repo, authSrv)
	s.Start()
}
