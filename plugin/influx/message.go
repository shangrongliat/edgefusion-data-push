package influx

type Detection struct {
	// 目标类别
	Class string `json:"class"`
	// 目标名称
	Name string ` json:"name,omitempty"`
	// 得分/概率
	Score float32 `json:"score,omitempty"`
	// 目标坐标，格式为 (x,y,w,h) x,y图片中心坐标，w宽 h高
	Box string `json:"box,omitempty"`
	// 目标切片/对目标进行标注后的图片
	// 注意，所有图片格式统一为png/jpg，不会再单独加一个字段表示图片格式
	Image []byte `json:"image,omitempty"`
	// 目标地理位置，格式为(lon,lat,height) 经度、纬度和高度，有些场景下可以从图片中解算出地理位置
	Location string `json:"location,omitempty"`
	// 格式化时间
	Time string `json:"time,omitempty"`
	// 时间戳
	Timestamp uint64 `json:"timestamp,omitempty"`
}
