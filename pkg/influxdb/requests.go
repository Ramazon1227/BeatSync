package influxdb

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Ramazon1227/BeatSync/config"
)

func DeleteUser(ctx context.Context,email string) error{
    // Define your InfluxDB parameters
	cfg:= config.Load()
	influxURL := cfg.InfluxURL // Replace with your InfluxDB URL
    org := cfg.InfluxOrg                                // Replace with your organization name
    bucket := cfg.InfluxBucket                          // Replace with your bucket name
    token := cfg.InfluxToken                            // Replace with your API token

    // Define the time range for deletion
    start := "2025-03-01T00:00:00Z"
    stop := time.Now().Format(time.RFC3339)

    // Define the predicate for deletion
    predicate := fmt.Sprintf(`_measurement=\"user_info\" AND email=\"%s\"`, email)

    // Construct the delete request body
    deleteBody := fmt.Sprintf(`{
        "start": "%s",
        "stop": "%s",
        "predicate": "%s"
    }`, start, stop, predicate)

    // Create the HTTP request
    req, err := http.NewRequestWithContext(context.Background(), http.MethodPost, fmt.Sprintf("%s/api/v2/delete?org=%s&bucket=%s", influxURL, org, bucket), bytes.NewBuffer([]byte(deleteBody)))
    if err != nil {
        panic(err)
    }

    // Set the appropriate headers
    req.Header.Set("Authorization", "Token "+token)
    req.Header.Set("Content-Type", "application/json")
    fmt.Println("Request:", req)
    // Execute the request
    client := &http.Client{Timeout: 10 * time.Second}
    resp, err := client.Do(req)
    if err != nil {

        log.Fatal(err)
		return err
    }
    defer resp.Body.Close()

    // Check the response status
    if resp.StatusCode != http.StatusNoContent {
        fmt.Printf("Failed to delete data: %s\n", resp.Status)
		return fmt.Errorf("failed to delete data: %s", resp.Status)
    } else {
        fmt.Println("Data deletion successful.")
    }
	return nil
}
