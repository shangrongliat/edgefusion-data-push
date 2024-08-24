package api

import (
	"edgefusion-data-push/plugin/config"
	log "edgefusion-data-push/plugin/logs"
	"edgefusion-data-push/service"
)

type API struct {
	api, rtmp, hls string
	log            *log.Logger
	storage        service.StorageService
}

// NewAPI new api
func NewAPI(config *config.Config) (*API, error) {
	storageService, err := service.NewStorageService(config)
	if err != nil {
		return nil, err
	}
	return &API{
		api:     config.Hook.API,
		rtmp:    config.Hook.RTMP,
		hls:     config.Hook.HLS,
		storage: storageService,
		log:     log.L().With(log.Any("api", "admin")),
	}, nil
}
