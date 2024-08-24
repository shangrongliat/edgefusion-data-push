package bean

type StreamInfo struct {
	ClientId     string `json:"client_id"`
	App          string `json:"app"`
	Vhost        string `json:"vhost"`
	GenerateTime int    `json:"generate_time"`
	Active       bool   `json:"active" default:"false"`
}

type ClientStreamInfo struct {
	ClientId     string `json:"client_id"`
	GenerateTime int    `json:"generate_time"`
	Active       bool   `json:"active" default:"false"`
	Type         string `json:"type"`
}
