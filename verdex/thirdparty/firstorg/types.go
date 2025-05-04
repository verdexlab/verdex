package firstorg

type EpssResponse struct {
	Data []*EpssData `json:"data"`
}

type EpssData struct {
	CveID string `json:"cve"`
	EPSS  string `json:"epss"`
}
