package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/sensu-community/sensu-plugin-sdk/sensu"
	"github.com/sensu/sensu-go/types"
)

// Config holds the configuration options for the Sensu plugin.
type Config struct {
	sensu.PluginConfig
	Debug              bool   // Enable debug mode
	InsecureSkipVerify bool   // Skip TLS certificate verification (not recommended!)
	Timeout            int    // Request timeout in seconds
	URL                string // URL to test
}

var (
	config = Config{
		PluginConfig: sensu.PluginConfig{
			Name:     "sensu-fastapi-appsumo-healthcheck",
			Short:    "Check FastAPI AppSumo health status",
			Keyspace: "sensu.io/plugins/sensu-django-healthcheck/config",
		},
	}

	options = []*sensu.PluginConfigOption{
		{
			Path:      "debug",
			Env:       "SENSU_CHECK_DEBUG",
			Argument:  "debug",
			Shorthand: "d",
			Default:   false,
			Usage:     "Enable debug mode",
			Value:     &config.Debug,
		},
		{
			Path:      "insecure-skip-verify",
			Env:       "SENSU_CHECK_INSECURE_SKIP_VERIFY",
			Argument:  "insecure-skip-verify",
			Shorthand: "i",
			Default:   false,
			Usage:     "Skip TLS certificate verification (not recommended!)",
			Value:     &config.InsecureSkipVerify,
		},
		{
			Path:      "timeout",
			Env:       "SENSU_CHECK_TIMEOUT",
			Argument:  "timeout",
			Shorthand: "T",
			Default:   15,
			Usage:     "Request timeout in seconds",
			Value:     &config.Timeout,
		},
		{
			Path:      "url",
			Env:       "SENSU_CHECK_URL",
			Argument:  "url",
			Shorthand: "u",
			Default:   "http://localhost:80/",
			Usage:     "URL to test",
			Value:     &config.URL,
		},
	}
)

func main() {
	check := sensu.NewGoCheck(&config.PluginConfig, options, checkArgs, executeCheck, false)
	// Call Execute without expecting a return value
	check.Execute()
	// After Execute, if you need to handle specific cases, you can do so here
	// For example, you might want to log a message or set an exit status based on internal logic
}

// checkArgs checks if essential configuration parameters are provided.
func checkArgs(event *types.Event) (int, error) {
	if config.URL == "" {
		return sensu.CheckStateWarning, fmt.Errorf("url is required")
	}
	return sensu.CheckStateOK, nil
}

func executeCheck(event *types.Event) (int, error) {
	// Make HTTP GET request to URL
	client := &http.Client{Timeout: time.Duration(config.Timeout) * time.Second}
	resp, err := client.Get(config.URL)
	if err != nil {
		return sensu.CheckStateCritical, fmt.Errorf("failed to fetch URL: %v", err)
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return sensu.CheckStateCritical, fmt.Errorf("failed to read response body: %v", err)
	}

	// this is an example of what the response body will look like, there may be more or less keys
	// we need to check each key and make sure it has a value of "working"
	// {
	// 	"status":"ok",
	// 	"backends":
	// 		{
	// 			"Cache backend: default":"working",
	// 			"DatabaseBackend":"working",
	// 			"DefaultFileStorageHealthCheck":"working",
	// 			"MigrationsHealthCheck":"working",
	// 			"ProductsHealthCheckBackend":"working"
	// 		}
	// }

	// Parse JSON response into a map[string]interface{}
	var data map[string]interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		return sensu.CheckStateCritical, fmt.Errorf("failed to unmarshal JSON response: %v", err)
	}

	// Debug log the parsed JSON data if debug mode is enabled
	if config.Debug {
		fmt.Println("Parsed JSON data:", data)
	}

	allWorking := true
	warnings := []string{}

	// check if the status key has a value of "ok"
	if data["status"] != "ok" {
		allWorking = false
		warnings = append(warnings, fmt.Sprintf("status is not ok, expected ok"))
	}

	// now we need to check each key in the backends map and make sure it has a value of "working"
	backends, ok := data["backends"].(map[string]interface{})
	if !ok {
		return sensu.CheckStateCritical, fmt.Errorf("backends is not a map")
	}

	for key, value := range backends {
		if value != "working" {
			allWorking = false
			warnings = append(warnings, fmt.Sprintf("key %s has a value of %s, expected working", key, value))
		}
	}

	if !allWorking {
		fmt.Printf("some checks are not working: %v", warnings)
		return sensu.CheckStateCritical, nil
	}

	fmt.Printf("all checks are working: %v", data)
	return sensu.CheckStateOK, nil
}
