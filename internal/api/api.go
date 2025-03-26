package api

import (
	"github.com/mehmetalisavas/message-sender/config"
	"github.com/mehmetalisavas/message-sender/internal/service"
)

type Api struct {
	config         *config.Config
	storageService service.Storage
}

func New(cfg *config.Config, storageService service.Storage) *Api {
	return &Api{
		config:         cfg,
		storageService: storageService,
	}
}
