package bean

type ConnectInfo struct {
	TcUrl   string `json:"tc_url"`
	PageUrl string `json:"page_url"`
}

type PublishInfo struct {
	Action   string `json:"action"`
	ClientId int    `json:"client_id"`
	IP       string `json:"ip"`
	Vhost    string `json:"vhost"`
	App      string `json:"app"`
	Stream   string `json:"stream"`
}

type PlayInfo struct {
	Action   string `json:"action"`
	ClientId int    `json:"client_id"`
	IP       string `json:"ip"`
	Vhost    string `json:"vhost"`
	App      string `json:"app"`
	Stream   string `json:"stream"`
	Param    string `json:"param"`
	PageUrl  string `json:"page_url"`
}
