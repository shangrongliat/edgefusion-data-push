package model

type VideoInfo struct {
	ID      string `json:"id"`
	NodeId  string `json:"node_id"`
	AppName string `json:"app_name"`
	Status  string `json:"status"`
	With    string `json:"with"`
	Height  string `json:"height"`
	Fps     string `json:"fps"`
}
