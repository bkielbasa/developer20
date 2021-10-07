package main

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

//go:embed email-template.html
var emailTemplate string

type config struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	RefreshToken string `json:"expiry"`
	Expiry       string `json:"expiry"`
}

type collection struct {
	ID    int `json:"_id"`
	Title string
}

type collectionResponse struct {
	Items []collection
}

type rainDrop struct {
	ID      int `json:"_id"`
	Title   string
	Excerpt string
	Link    string
	Tags    []string `json:"tags"`
	Cover   string   `json:"cover"`
}

type rainDropResponse struct {
	Items []rainDrop
}

func main() {
	b, err := os.ReadFile("token.json")
	if err != nil {
		log.Fatal(err)
	}

	conf := config{}
	if err = json.Unmarshal(b, &conf); err != nil {
		log.Fatal(err)
	}

	req, _ := http.NewRequest(http.MethodGet, "https://api.raindrop.io/rest/v1/collections", nil)
	req.Header["Authorization"] = []string{fmt.Sprintf("%s %s", conf.TokenType, conf.AccessToken)}

	client := http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	coll := collectionResponse{}
	if err = json.Unmarshal(body, &coll); err != nil {
		log.Fatal(err)
	}

	toSend := map[string][]rainDrop{}

	for _, c := range coll.Items {
		log.Print(c.Title)
		ids := []int{}
		req, _ = http.NewRequest(http.MethodGet, fmt.Sprintf("https://api.raindrop.io/rest/v1/raindrops/%d?search=notag:true", c.ID), nil)
		req.Header["Authorization"] = []string{fmt.Sprintf("%s %s", conf.TokenType, conf.AccessToken)}
		resp, err = client.Do(req)
		if err != nil {
			log.Fatal(err)
		}

		body, err = io.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}

		drops := rainDropResponse{}
		if err = json.Unmarshal(body, &drops); err != nil {
			log.Fatal(err)
		}

		for _, drop := range drops.Items {
			log.Printf(" * %s %s\n", drop.Title, drop.Cover)
			ids = append(ids, drop.ID)
		}

		if len(drops.Items) > 0 {
			toSend[c.Title] = drops.Items
			body, _ = json.Marshal(struct {
				CollectionID int      `json:"collectionId"`
				Tags         []string `json:"tags"`
				IDs          []int    `json:"ids,omitempty"`
			}{
				CollectionID: c.ID,
				Tags:         []string{},
				IDs:          ids,
			})

			log.Printf("updating raindrops: %s", string(body))
			req, _ = http.NewRequest(http.MethodPut, fmt.Sprintf("https://api.raindrop.io/rest/v1/raindrops/%d", c.ID), bytes.NewReader(body))
			req.Header["Authorization"] = []string{fmt.Sprintf("%s %s", conf.TokenType, conf.AccessToken)}
			req.Header["Content-Type"] = []string{"application/json"}

			resp, err = client.Do(req)
			if err != nil {
				log.Fatal(err)
			}

			body, _ = io.ReadAll(resp.Body)
			if resp.StatusCode != 200 {
				panic(fmt.Sprintf("cannot update tags: %d %s", resp.StatusCode, string(body)))
			}
		}
	}

	t, err := template.New("foo").Parse(emailTemplate)
	if err != nil {
		panic(err)
	}

	buf := new(bytes.Buffer)
	err = t.Execute(buf, toSend)
	if err != nil {
		panic(err)
	}

	if err = sendCampaign(buf.String()); err != nil {
		panic(err)
	}
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
