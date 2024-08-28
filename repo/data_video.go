package repo

import (
	"edgefusion-data-push/plugin/context"
	"edgefusion-data-push/repo/model"
	"gorm.io/gorm"
)

type IDataVideo interface {
	Create(dataInfo *model.DataVideo) error
	GetByNodeIdAndAppName(nodeId, appName string) (*model.DataVideo, error)
	UpdateStateByNodeIdAndAppName(nodeId, appName string, state int8) error
}

type DataVideo struct {
	db *gorm.DB
}

func NewDataVideo() IDataVideo {
	return &DataVideo{
		db: context.DatabaseMapHandle(),
	}
}

func (d *DataVideo) Create(dataInfo *model.DataVideo) error {
	tx := d.db.
		Debug().
		Model(&model.DataVideo{}).
		Create(dataInfo)
	return tx.Error
}

func (d *DataVideo) GetByNodeIdAndAppName(nodeId, appName string) (*model.DataVideo, error) {
	var dataVideo model.DataVideo
	tx := d.db.
		Debug().
		Model(&model.DataVideo{}).
		Where("node_id = ?", nodeId).
		Where("app_name = ?", appName).
		First(&dataVideo)
	return &dataVideo, tx.Error
}

func (d *DataVideo) UpdateStateByNodeIdAndAppName(nodeId, appName string, state int8) error {
	tx := d.db.
		Debug().
		Model(&model.DataVideo{}).
		Where("node_id = ?", nodeId).
		Where("app_name = ?", appName).
		Update("state", state)
	return tx.Error
}
