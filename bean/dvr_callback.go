package bean

type DvrCallBackInfo struct {
	Action   string `json:"action"`
	ClientID int    `json:"client_id"`
	IP       string `json:"ip"`
	Vhost    string `json:"vhost"`
	App      string `json:"app"`
	Stream   string `json:"stream"`
	Cwd      string `json:"cwd"`
	File     string `json:"file"`
}
