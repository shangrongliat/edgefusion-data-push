package api

import (
	"edgefusion-data-push/plugin/config"
)

func (a *API) GetInfluxData(ctx *config.Context) (any, error) {
	a.storage.Test()
	return nil, nil
}

func (a *API) WInfluxData(ctx *config.Context) (any, error) {
	a.storage.Testw()
	return nil, nil
}
