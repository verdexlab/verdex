package core

type CVE struct {
	ID                    string   `json:"id"`
	GHSA                  *string  `json:"ghsa"`
	Summary               string   `json:"summary"`
	Description           string   `json:"description"`
	Severity              string   `json:"severity"`
	CvssScore             float32  `json:"cvss_score"`
	EpssScore             float32  `json:"epss_score"`
	KevSince              *string  `json:"kev_since"`
	POCs                  []string `json:"pocs"`
	References            []string `json:"references"`
	VendorAdvisory        *string  `json:"vendor_advisory"`
	NucleiTemplate        *string  `json:"nuclei_template"`
	VulnerableVersions    string   `json:"vulnerable_versions"`
	NearestPatchedVersion *string  `json:"nearest_patched_version"`
	PublishedAt           string   `json:"published_at"`
	UpdatedAt             string   `json:"updated_at"`
}
