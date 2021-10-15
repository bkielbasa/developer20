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

//go:embed email-template.xml
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
				Tags:         []string{"newsletter", time.Now().Format("2006-01-02")},
				IDs:          ids,
			})

			log.Print(string(body))
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
			log.Print(string(body))
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

	tmpFileName := "generated.mjml"
	os.WriteFile(tmpFileName, []byte(buf.String()), 777)
}
