package api

import (
	"fmt"
	"strings"
	"time"

	"edgefusion-data-push/bean"
	"edgefusion-data-push/cache"
	"edgefusion-data-push/plugin/config"
	log "edgefusion-data-push/plugin/logs"
)

func (a *API) Connect(ctx *config.Context) any {
	var connect bean.ConnectInfo
	if err := ctx.ShouldBindJSON(&connect); err != nil {
		log.L().Error("connect info 解析失败.", log.Error(err))
		return 500
	}
	return 0
}

func (a *API) Publish(ctx *config.Context) any {
	var publish bean.PublishInfo
	if err := ctx.ShouldBindJSON(&publish); err != nil {
		log.L().Error("publish info 解析失败.", log.Error(err))
		return 500
	}
	authStreamInfo, err := cache.Cache.PullCache(publish.Stream)
	if err != nil || authStreamInfo == nil {
		// 回调 关闭此次请求
		return 500
	}
	if shell, ok := authStreamInfo.(bean.StreamInfo); ok {
		shell.Active = true
		shell.GenerateTime = time.Now().Second()
		return 0
	}
	// 不存在则返回失败
	return 500
}

func (a *API) UnPublish(ctx *config.Context) any {
	var publish bean.PublishInfo
	if err := ctx.ShouldBindJSON(&publish); err != nil {
		log.L().Error("un_publish info 解析失败.", log.Error(err))
		return 500
	}
	log.L().Debug("un_publish", log.Any("data", publish))
	return 0
}

func (a *API) DvrFinish(ctx *config.Context) (any, error) {
	var dvr bean.DvrCallBackInfo
	if err := ctx.ShouldBindJSON(&dvr); err != nil {
		log.L().Error("dvr info 解析失败.", log.Error(err))
	}
	// 录播结束执行 文件存储，将文件存储到minio中，并建立节点-应用-文件关系
	// 假设直播路径为 rtmp://172.16.100.14:1935/{app}/{stream}?vhost=edgefusion
	// app = 节点ID-应用名称
	// stream = 随机串
	a.storage.VideoStorage(dvr)
	return 0, nil
}

func (a *API) Play(ctx *config.Context) any {
	//var play bean.PlayInfo
	//if err := ctx.ShouldBindJSON(&play); err != nil {
	//	log.L().Error("play info 解析失败.", log.Error(err))
	//	return 500
	//} 54.01
	//paramMap := a.paramMap(play.Param)
	//m := paramMap["key"]
	return 0
}

func (a *API) getDeleteUrl(clientId int) string {
	return fmt.Sprintf("%s%d", a.api, clientId)
}

func (a *API) paramMap(param string) map[string]string {
	paramMap := make(map[string]string)
	if len(param) > 0 {
		replace := strings.Replace(param, "?", "", -1)
		params := strings.Split(replace, "&")
		for _, s := range params {
			split1 := strings.Split(s, "=")
			paramMap[split1[0]] = split1[1]
		}
	}
	return map[string]string{}
}

func (a *API) getClientStreamInfo(app, vhost, stream, key string, clientId int) bool {
	check := checkAppAndStreamAndVhost(app, vhost, stream)
	if check {

	}
	return false
}

func checkAppAndStreamAndVhost(app, stream, vhost string) bool {
	info := getStream(stream)
	if info == nil {
		return false
	}
	return info.App == app && info.Vhost == vhost
}

func getStream(stream string) *bean.StreamInfo {
	param, err := cache.Cache.GetCache(stream)
	if err != nil || param == nil {
		return nil
	}
	if pa, ok := param.(bean.StreamInfo); ok {
		return &pa
	}
	return nil
}
