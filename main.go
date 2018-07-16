package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

const Address string = "https://checkpoint-api.hashicorp.com/v1"

var product string

var products []string = []string{
	"consul",
	"nomad",
	"packer",
	"terraform",
	"vagrant",
}

type Check struct {
	Alerts              []string `json:"alerts"`
	CurrentChangelogURL string   `json:"current_changelog_url"`
	CurrentDownloadURL  string   `json:"current_download_url"`
	CurrentRelease      uint     `json:"current_release"`
	CurrentVersion      string   `json:"current_version"`
	Product             string   `json:"product"`
	ProjectWebsite      string   `json:"project_website"`
}

func init() {
	flag.StringVar(&product, "product", "", fmt.Sprintf("Must be one of %s", strings.Join(products, ", ")))
}

func main() {
	flag.Parse()
	if product == "" {
		fmt.Println("checkpoint: product cannot be empty")
		os.Exit(1)
	}
	var found bool
	for i := range products {
		if product == products[i] {
			found = true
			break
		}
	}
	if !found {
		fmt.Printf("checkpoint: product must be one of %s\n", strings.Join(products, ", "))
		os.Exit(1)
	}
	response, err := http.Get(fmt.Sprintf("%s/check/%s", Address, product))
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()
	b, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	var check Check
	if err := json.Unmarshal(b, &check); err != nil {
		log.Fatal(err)
	}
	fmt.Println(check.CurrentVersion)
	os.Exit(0)
}
