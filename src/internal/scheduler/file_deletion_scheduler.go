package scheduler

import (
	"context"
	"time"

	"github.com/rs/zerolog/log"
	"shup.hilmy.dev/src/internal/service"
)

func FileDeletionScheduler(fileService *service.File) {
	ticker := time.NewTicker(10 * time.Minute)
	defer ticker.Stop()

	for {
		// Run task
		log.Info().Msg("Run file deletion scheduler")
		fileService.DeleteAllExpired(context.Background())

		// Wait for the next tick
		<-ticker.C
	}
}
