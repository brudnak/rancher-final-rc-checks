package main

import (
	"crypto/tls"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
)

func main() {
	// Handling command-line arguments
	fqdn := flag.String("fqdn", "", "Rancher FQDN, e.g. https://<Rancher FQDN>")
	apiBearerToken := flag.String("api", "", "API Bearer token")
	flag.Parse()

	// List of settings to check - URL tags
	listOfUrls := []string{"rke-version", "ui-index", "ui-dashboard-index", "cli-url-linux",
		"cli-url-darwin", "cli-url-windows", "system-catalog", "kdm-branch",
		"ui-k8s-supported-versions-range"}
	urlLabel := []string{"Released RKE version", "UI Tag", "UI Dashboard Index", "CLI URL (Linux)",
		"CLI URL (Darwin)", "CLI URL (Windows)", "System Chart Catalog", "KDM branch",
		"UI k8s supported versions range"}

	// Disable certificate verification
	httpClient := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	// Add content type and api_bearer_token to the header
	headers := make(map[string]string)
	headers["Content-Type"] = "application/json"
	headers["Authorization"] = "Bearer " + *apiBearerToken

	// Iterate through listOfUrls to get the values
	for i, url := range listOfUrls {
		webUrl := *fqdn + "/v3/settings/" + url
		request, err := http.NewRequest("GET", webUrl, nil)
		if err != nil {
			fmt.Println("Error creating request:", err)
			return
		}

		for key, value := range headers {
			request.Header.Add(key, value)
		}

		response, err := httpClient.Do(request)
		if err != nil {
			fmt.Println("Error on request:", err)
			return
		}
		defer response.Body.Close()

		if response.StatusCode == 200 {
			var result map[string]interface{}
			body, err := ioutil.ReadAll(response.Body)
			if err != nil {
				fmt.Println("Error reading response body:", err)
				return
			}
			err = json.Unmarshal(body, &result)
			if err != nil {
				fmt.Println("Error unmarshalling response:", err)
				return
			}
			fmt.Printf("\t%s: %s\n", urlLabel[i], result["value"])
		} else {
			fmt.Printf("\t%s is not set\n", urlLabel[i])
		}
	}
}
