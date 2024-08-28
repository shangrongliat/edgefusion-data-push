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

func (a *API) GetHlsPullPath(ctx *config.Context) (any, error) {
	node := ctx.Param("node")
	app := ctx.Param("app")
	hls := "172.16.100.14"
	// 示例： http://172.16.100.14:8020/node01/live.m3u8?vhost=edgefusion
	hlsPath := fmt.Sprintf("http://%s:8020/%s/%s.m3u8?vhost=edgefusion", hls, node, app)
	return hlsPath, nil
}
