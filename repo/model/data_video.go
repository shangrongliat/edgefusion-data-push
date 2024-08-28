package model

type DataVideo struct {
	ID         string  `json:"id" gorm:"id"`                   // 主键id
	NodeId     string  `json:"node_id" gorm:"node_id"`         // 节点id
	AppName    string  `json:"app_name" gorm:"app_name"`       // 应用名称
	Stream     string  `json:"stream" gorm:"stream"`           // 流
	Vhost      string  `json:"vhost" gorm:"vhost"`             // 直播后缀
	State      int8    `json:"state" gorm:"state"`             // 直播状态 1 直播中 2 直播结束
	DataName   string  `json:"data_name" gorm:"data_name"`     // 数据名称(图片与视频为文件名称，时序数据为measurement名称)
	LivePath   string  `json:"live_path" gorm:"live_path"`     // 文件名称(数据类型为时序信息时为空)
	Size       float64 `json:"size" gorm:"size"`               // 数据大小
	Width      int     `json:"width" gorm:"width"`             // 视频宽度
	Height     int     `json:"height" gorm:"height"`           // 视频高度
	Fps        int     `json:"fps" gorm:"fps"`                 // 视频帧率
	DFlag      int64   `json:"d_flag" gorm:"d_flag"`           // 是否删除 1 删除 0 未删除
	DateTime   int64   `json:"date_time" gorm:"date_time"`     // 创建时间戳
	CreateTime string  `json:"create_time" gorm:"create_time"` // 创建时间
}

// TableName 表名称
func (*DataVideo) TableName() string {
	return "ef_data_video"
}
