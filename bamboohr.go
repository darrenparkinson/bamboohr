// Package bamboohr provides a library for the Bamboo HR API
package bamboohr

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"gopkg.in/errgo.v2/errors"
)

// Client represents connectivity to the bamboo hr API
type Client struct {
	// Base URL for Bamboo HR API which is set to v1 using the provided company domain if initiated with `bamboohr.New()`
	BaseURL string

	// HTTP Client to use for making requests allowing the user to supply their own if required.
	HTTPClient *http.Client

	// Base64 Encoded string based on the APIKey, used for Basic Authorization
	Auth string
}

// New is a helper function that returns a new instance of the bamboo hr client given a company domain and api key.
// An http.Client will be created if nil is provided.
func New(apikey string, companyDomain string, client *http.Client) (*Client, error) {
	if apikey == "" {
		return nil, errors.New("apikey required")
	}
	if companyDomain == "" {
		return nil, errors.New("companyDomain required")
	}
	if client == nil {
		client = &http.Client{
			Timeout: 10 * time.Second,
		}
	}
	c := &Client{
		BaseURL:    fmt.Sprintf("https://api.bamboohr.com/api/gateway.php/%s/v1", companyDomain),
		HTTPClient: client,
		Auth:       fmt.Sprintf("Basic %s", base64.StdEncoding.EncodeToString([]byte(apikey+":x"))),
	}
	return c, nil
}

// makeRequest provides a single function to add common items to the request.
func (c *Client) makeRequest(req *http.Request, v interface{}) error {
	// Set standard headers
	req.Header.Set("Authorization", c.Auth)
	req.Header.Set("Accept", "application/json")
	// Make the request
	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	// Check we have a desired status code, e.g. between 200 and 400
	if res.StatusCode < http.StatusOK || res.StatusCode >= http.StatusBadRequest {
		return fmt.Errorf("error from bamboo, status code: %d", res.StatusCode)
	}
	// If we're just getting a created (201), then it's ok. We might want to return a struct at some point
	if res.StatusCode == http.StatusCreated {
		return nil
	}
	// Decode the body to the supplied interface
	if err = json.NewDecoder(res.Body).Decode(&v); err != nil {
		return err
	}
	return nil
}
