package openmeteo

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

type ForecastParams struct {
	Longitude float32
	Latitude  float32
	Timezone  string
}

type Client struct {
	host   string
	client *http.Client
}

const defaultApiHost = "https://api.open-meteo.com"

var ErrOpenMeteoService = errors.New("open-meteo service error")

func New(host string) *Client {
	apiHost := defaultApiHost
	if apiHost != "" {
		apiHost = host
	}

	return &Client{
		host: apiHost,
		client: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

type ForecastResponse struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Current   struct {
		Time        string  `json:"time"`
		Temperature float64 `json:"temperature_2m"`
	} `json:"current"`
	Error  bool   `json:"error"`
	Reason string `json:"reason"`
}

const forecastURL = "/v1/forecast?latitude=%f&longitude=%f&current=temperature_2m&hourly=temperature_2m&timezone=%s"

// Forecast returns current temperature for the provided location
func (cl *Client) Forecast(params ForecastParams) (float64, error) {
	addr := fmt.Sprintf(cl.host+forecastURL, params.Latitude, params.Longitude, params.Timezone)

	req, err := http.NewRequest(http.MethodGet, addr, nil)

	resp, err := cl.client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return 0, ErrOpenMeteoService
	}

	var response ForecastResponse

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	if err := json.Unmarshal(responseBody, &response); err != nil {
		return 0, err
	}

	if response.Error {
		return 0, errors.New(response.Reason)
	}

	return response.Current.Temperature, nil
}
