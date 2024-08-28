package service

import (
	"encoding/json"
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
	ImageStorage(nodeId, appName, time string, data []byte)
	TimeSeriesStorage(nodeId, appId string, time uint64, target *message.Target)
	TargetStorage(nodeId, appName string, time uint64, targets []*message.Target) error

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

func (s *StorageServiceImpl) ImageStorage(nodeId, appName, time string, data []byte) {
	fileName := fmt.Sprintf("%v%v.jpg", time, rand.Intn(900)+100)
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
	inTarget.Class = target.Class                // 目标类别
	inTarget.Name = target.Name                  // 目标名称
	inTarget.Score = target.Score                // 得分/概率
	inTarget.Box = target.Box                    // 目标坐标，格式为 (x,y,w,h) x,y图片中心坐标，w宽 h高
	inTarget.Location = target.Location          // 目标地理位置，格式为(lon,lat,height) 经度、纬度和高度，有些场景下可以从图片中解算出地理位置
	inTarget.Time = common.TimeUint2String(time) // 格式化时间
	inTarget.Timestamp = time                    // 时间戳
	marshal, err := json.Marshal(inTarget)
	if err != nil {
		log.L().Error("时序数据Json化失败", log.Error(err))
	}
	fields := map[string]interface{}{
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

func (s *StorageServiceImpl) TargetStorage(nodeId, appName string, time uint64, targets []*message.Target) error {
	if info := getTimeSeriesInfo(nodeId, appName); info == nil {
		// 判断内存中没有对应的time series 信息
		// 查询数据中是否存在
		infos, err := repo.NewDataInfo().GetByNodeIdAndAppNameAndType(nodeId, appName, common.TimeSeries)
		if err != nil {
			log.L().Debug("未查询到有效信息", log.Error(err))
			return err
		}
		if len(infos) != 0 {
			// 不等于0 说明有对应的数据，将数据添加到内存中
			setTimeSeriesInfo(nodeId, appName, &infos[0])
		} else {
			// 不存在则创建对应数据
			createDataInfo(nodeId, appName, fmt.Sprintf("%s-%s", nodeId, appName), "", common.TimeSeries, common.NilType, 0)
		}
	}
	var targetList []*influx.Detection
	for i, target := range targets {
		var inTarget influx.Detection
		inTarget.Class = target.Class                // 目标类别
		inTarget.Name = target.Name                  // 目标名称
		inTarget.Score = target.Score                // 得分/概率
		inTarget.Box = target.Box                    // 目标坐标，格式为 (x,y,w,h) x,y图片中心坐标，w宽 h高
		inTarget.Location = target.Location          // 目标地理位置，格式为(lon,lat,height) 经度、纬度和高度，有些场景下可以从图片中解算出地理位置
		inTarget.Time = common.TimeUint2String(time) // 格式化时间
		inTarget.Timestamp = time
		targetList = append(targetList, &inTarget)
		// 时间戳
		if len(target.Image) > 0 {
			// 如果长度大于0，则说明有图片
			// 进行图片存储
			s.ImageStorage(nodeId, appName, fmt.Sprintf("%v-%v", time, i), target.Image)
		}
	}
	marshal, err := json.Marshal(targetList)
	if err != nil {
		log.L().Error("时序数据Json化失败", log.Error(err))
		return err
	}
	fields := map[string]interface{}{
		"target": marshal,
	}
	if err = s.influx.Save(nodeId, appName, fields); err != nil {
		log.L().Error("时序数据存储失败",
			log.Error(err),
			log.Any("nodeId", nodeId),
			log.Any("appName", appName),
			log.Any("fields", fields),
		)
		return err
	}
	return nil
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
	info.DateTime = time.Now().Unix()
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
