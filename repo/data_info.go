package repo

import (
	"edgefusion-data-push/plugin/context"
	"edgefusion-data-push/repo/model"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type IDataInfo interface {
	Create(dataInfo model.DataInfo) error
	GetByNodeIdAndAppNameAndType(nodeId, appName string, dType int) ([]model.DataInfo, error)
}

type DataInfo struct {
	db *gorm.DB
}

func NewDataInfo() IDataInfo {
	return &DataInfo{
		db: context.DatabaseMapHandle(),
	}
}

func (d *DataInfo) Create(dataInfo model.DataInfo) error {
	tx := d.db.Debug().Model(&model.DataInfo{}).Create(dataInfo)
	if err := tx.Error; err != nil {
		return err
	}
	return nil
}

func (d *DataInfo) GetByNodeIdAndAppNameAndType(nodeId, appName string, dType int) ([]model.DataInfo, error) {
	var data []model.DataInfo
	tx := d.db.
		Debug().
		Model(&model.DataInfo{}).
		Where("node_id = ? and app_name = ? and type = ?", nodeId, appName, dType).
		Find(&data)
	if err := tx.Error; err != nil {
		return nil, errors.New("find data fail.")
	}
	return data, nil
}
