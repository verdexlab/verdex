package core

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

var reportUrl = "https://report.verdexlab.workers.dev"

type ReportData struct {
	TargetURL     string `json:"target_url"`
	Product       string `json:"product"`
	VerdexVersion string `json:"verdex_version"`
}

func ReportTarget(detection *Detection) error {
	payload := ReportData{
		TargetURL:     detection.Target,
		Product:       detection.Product,
		VerdexVersion: GetVerdexVersion(),
	}

	payloadJson, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	method := "POST"
	client := &http.Client{}
	req, err := http.NewRequest(method, reportUrl, bytes.NewReader(payloadJson))
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		body, err := io.ReadAll(res.Body)
		if err != nil {
			return err
		}

		return fmt.Errorf("failed to report target, invalid response (status %d): %s", res.StatusCode, string(body))
	}

	return nil
}
