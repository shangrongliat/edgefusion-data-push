package model

type DataInfo struct {
	ID         string  `json:"id" gorm:"id"`                   //数据ID
	NodeID     string  `json:"node_id" gorm:"node_id"`         //节点名称
	AppName    string  `json:"app_name" gorm:"app_name"`       //应用名称
	Type       int     `json:"type" gorm:"type"`               //类型(1 视频;2 图片;3时序数据)
	DataName   string  `json:"data_name" gorm:"data_name"`     //数据名称(图片与视频为文件名称，时序数据为measurement名称)
	DataPath   string  `json:"data_path" gorm:"data_path"`     //文件路径
	DataType   int     `json:"data_type" gorm:"data_type"`     //数据类型(视频:1mp4;2 flv，图片:1 png/jpg;2 image/bmp，时序数据时为空)
	Size       float64 `json:"size" gorm:"size"`               //数据大小
	CreateTime string  `json:"create_time" gorm:"create_time"` //创建时间
	DFlag      int8    `json:"d_flag" json:"d_flag"`
}

// TableName 表名称
func (*DataInfo) TableName() string {
	return "ef_data_info"
}
