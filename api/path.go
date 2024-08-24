package api

import (
	"fmt"
	"strings"

	"edgefusion-data-push/plugin/config"
	"github.com/google/uuid"
)

func (a *API) GetRtmpPutPath(ctx *config.Context) (any, error) {
	node := ctx.Param("node")
	app := ctx.Param("app")
	vhost := ctx.Param("vhost")
	stream := strings.ReplaceAll(uuid.New().String(), "-", "")
	localImageName := fmt.Sprintf("%s%s-%s/%s?vhost=%s", a.rtmp, node, app, stream, vhost)
	return localImageName, nil
}

func (a *API) GetRtmpPullPath(ctx *config.Context) (any, error) {
	node := ctx.Param("node")
	app := ctx.Param("app")
	stream := ctx.Param("stream")
	vhost := ctx.Param("vhost")
	key := strings.ReplaceAll(uuid.New().String(), "-", "")
	rtmpPushPath := fmt.Sprintf("%s%s-%s/%s?vhost=%s&key=%s", a.rtmp, node, app, stream, vhost, key)
	return rtmpPushPath, nil
}
