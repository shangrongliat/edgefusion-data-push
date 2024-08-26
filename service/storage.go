package service

import (
	"edgefusion-data-push/plugin/json"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"sync"
	"time"

	"edgefusion-data-push/bean"
	"edgefusion-data-push/common"
	"edgefusion-data-push/message"
	"edgefusion-data-push/plugin/config"
	"edgefusion-data-push/plugin/influx"
	log "edgefusion-data-push/plugin/logs"
	plugin "edgefusion-data-push/plugin/minio"
	"edgefusion-data-push/plugin/utils"
	"edgefusion-data-push/repo"
	"edgefusion-data-push/repo/model"
)

var TimeSeriesDataInfo = sync.Map{}

type StorageService interface {
	VideoStorage(dvr bean.DvrCallBackInfo)
	ImageStorage(nodeId, appName string, data []byte)
	TimeSeriesStorage(nodeId, appId string, time uint64, target *message.Target)
	Test()
	Testw()
}

type StorageServiceImpl struct {
	client plugin.Minio
	influx influx.InfluxRepo
}

func NewStorageService(cfg *config.Config) (StorageService, error) {
	influxClient, err := influx.NewInflux()
	if err != nil {
		return nil, err
	}
	minio, err := plugin.NewMinioService(cfg)
	if err != nil {
		return nil, err
	}
	return &StorageServiceImpl{
		client: minio,
		influx: influxClient,
	}, nil
}

func (s *StorageServiceImpl) VideoStorage(dvr bean.DvrCallBackInfo) {
	apps := strings.Split(dvr.App, "-")
	nodeId := apps[0]
	appName := apps[1]
	fileName := strings.Split(dvr.File, "/")
	bucketName := fmt.Sprintf("%s/%s/%s", nodeId, appName, fileName[len(fileName)-1])
	uploadIndo, err := s.client.PutFileObject(dvr.File, common.MinioVideoBucket, bucketName)
	if err != nil {
		log.L().Error("录播视频上传失败", log.Error(err))
		return
	}
	minioFileName := fmt.Sprintf("%s/%s", common.MinioVideoBucket, bucketName)
	createDataInfo(nodeId, appName, fileName[len(fileName)-1], minioFileName, common.ImageType, common.ImagePng, float64(uploadIndo.Size))
}

func (s *StorageServiceImpl) ImageStorage(nodeId, appName string, data []byte) {
	fileName := fmt.Sprintf("%v%v.jpg", time.Now().UnixMilli(), rand.Intn(900)+100)
	bucketName := fmt.Sprintf("%s/%s/%s", nodeId, appName, fileName)
	if err := s.client.PutStreamObject(common.MinioImageBucket, bucketName, data); err != nil {
		log.L().Error("图片上传失败",
			log.Error(err),
			log.Any("nodeId", nodeId),
			log.Any("appName", appName))
	}
	path := fmt.Sprintf("%s/%s", common.MinioImageBucket, bucketName)
	createDataInfo(nodeId, appName, fileName, path, common.ImageType, common.ImagePng, float64(len(data)))
}

func (s *StorageServiceImpl) TimeSeriesStorage(nodeId, appName string, time uint64, target *message.Target) {
	if info := getTimeSeriesInfo(nodeId, appName); info == nil {
		// 判断内存中没有对应的time series 信息
		// 查询数据中是否存在
		infos, err := repo.NewDataInfo().GetByNodeIdAndAppNameAndType(nodeId, appName, common.TimeSeries)
		if err != nil {
			log.L().Debug("未查询到有效信息", log.Error(err))
			return
		}
		if len(infos) != 0 {
			// 不等于0 说明有对应的数据，将数据添加到内存中
			setTimeSeriesInfo(nodeId, appName, &infos[0])
		} else {
			// 不存在则创建对应数据
			createDataInfo(nodeId, appName, fmt.Sprintf("%s-%s", nodeId, appName), "", common.TimeSeries, common.NilType, 0)
		}
	}
	var inTarget influx.Detection
	inTarget.Score = target.Score
	inTarget.Box = target.Box
	inTarget.Location = target.Location
	inTarget.Class = target.Class
	inTarget.Name = target.Name
	inTarget.Time = common.Time2String(time)
	marshal, err := json.Marshal(target)
	if err != nil {
		log.L().Error("", log.Error(err))
	}
	fields := map[string]interface{}{
		// 目标类别
		"target": marshal,
	}
	if err := s.influx.Save(nodeId, appName, fields); err != nil {
		log.L().Error("时序数据存储失败",
			log.Error(err),
			log.Any("nodeId", nodeId),
			log.Any("appName", appName),
			log.Any("fields", fields),
		)
		return
	}
}

func (s *StorageServiceImpl) Test() {
	if err := s.influx.Get("IGA0LwM2w1WGVmXw", "ef-msg-distributor"); err != nil {
		log.L().Error("时序数据查询失败", log.Error(err))
	}
}

func (s *StorageServiceImpl) Testw() {
	for i := 0; i < 50; i++ {
		var target influx.Detection
		target.Score = 12.21312312
		target.Box = "[1,2,3,4]"
		target.Location = "[1.23,12.321.321]"
		target.Class = fmt.Sprintf("12345%v", i)
		target.Name = fmt.Sprintf("car%v", i)
		marshal, err := json.Marshal(target)
		if err != nil {
			log.L().Error("", log.Error(err))
		}
		fields := map[string]interface{}{
			// 目标类别
			"target": marshal,
		}
		if err := s.influx.Save("FDt4zjxNrTnohMt3", "skills-test-004", fields); err != nil {
			log.L().Error("时序数据查询失败", log.Error(err))
		}
		time.Sleep(1 * time.Second)
	}
}

func createDataInfo(nodeId, appName, fileName, path string, dType, dataType int, dataSize float64) {
	value, err := strconv.ParseFloat(fmt.Sprintf("%.2f", dataSize/(1024*1024)), 64)
	if err != nil {
		return
	}
	var info model.DataInfo
	info.ID = utils.ToStringUuid()
	info.AppName = appName
	info.NodeID = nodeId
	info.DataName = fileName
	info.DataPath = path
	info.Size = value
	info.Type = dType
	info.DataType = dataType
	info.CreateTime = utils.Time2String(time.Now().Unix())
	if err := repo.NewDataInfo().Create(info); err != nil {
		log.L().Error("数据信息写入失败.", log.Error(err))
	}
	if dType == common.TimeSeries {
		setTimeSeriesInfo(nodeId, appName, &info)
	}
}

func getTimeSeriesInfo(nodeId, appName string) *model.DataInfo {
	value, ok := TimeSeriesDataInfo.Load(nodeId + "_" + appName)
	if ok {
		return value.(*model.DataInfo)
	} else {
		return nil
	}
}

func setTimeSeriesInfo(nodeId, appName string, data *model.DataInfo) {
	TimeSeriesDataInfo.Store(nodeId+"_"+appName, data)
}
