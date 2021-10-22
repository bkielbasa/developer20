package main

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

//go:embed generated.html
var generated string

func main() {
	sendCampaign(generated)
}

func sendCampaign(template string) error {
	d := time.Now().Format("2 Jan 2006")
	createCampaign := map[string]interface{}{
		"subject": "List of links for the day " + d,
		"name":    "List of links for the day " + d,
		"type":    "regular",
		"groups":  []int{9634420},
	}

	token := "968a96c12b482a75e436e8f7b9c4371b"

	body, _ := json.Marshal(createCampaign)
	req, _ := http.NewRequest(http.MethodPost, "https://api.mailerlite.com/api/v2/campaigns", bytes.NewReader(body))
	req.Header["X-MailerLite-ApiKey"] = []string{token}
	req.Header["Content-Type"] = []string{"application/json"}

	client := http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("cannot prepare a draft of the campaing: %d", resp.StatusCode)
	}

	body, _ = io.ReadAll(resp.Body)
	createResponse := struct{ ID int }{}
	_ = json.Unmarshal(body, &createResponse)

	body, _ = json.Marshal(struct {
		HTML  string `json:"html"`
		Plain string `json:"plain"`
	}{
		HTML:  template,
		Plain: "Your email client does not support HTML emails. Open newsletter here: {$url}. If you do not want to receive emails from us, click here: {$unsubscribe}",
	})
	req, _ = http.NewRequest(http.MethodPut, fmt.Sprintf("https://api.mailerlite.com/api/v2/campaigns/%d/content", createResponse.ID), bytes.NewReader(body))
	req.Header["X-MailerLite-ApiKey"] = []string{token}
	req.Header["Content-Type"] = []string{"application/json"}
	resp, err = client.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		body, _ = io.ReadAll(resp.Body)
		return fmt.Errorf("cannot add content to the campaing: %d %s", resp.StatusCode, string(body))
	}

	req, _ = http.NewRequest(http.MethodPost, fmt.Sprintf("https://api.mailerlite.com/api/v2/campaigns/%d/actions/send", createResponse.ID), nil)
	req.Header["X-MailerLite-ApiKey"] = []string{token}
	req.Header["Content-Type"] = []string{"application/json"}
	resp, err = client.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		body, _ = io.ReadAll(resp.Body)
		return fmt.Errorf("cannot send the campaing: %d %s", resp.StatusCode, string(body))
	}

	return nil
}
