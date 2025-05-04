package core

type CVE struct {
	ID                 string   `json:"id"`
	Description        string   `json:"description"`
	CvssScore          float32  `json:"cvss_score"`
	EpssScore          float32  `json:"epss_score"`
	IsKEV              bool     `json:"is_kev"`
	VulnerableVersions []string `json:"vulnerable_versions"`
	PublishedAt        string   `json:"published_at"`
}
