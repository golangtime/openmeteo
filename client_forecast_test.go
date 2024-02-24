package openmeteo

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

const testResponse = `{
	"latitude": 0,
	"longitude": 0,
	"utc_offset_seconds": 0,
	"timezone": "UTC",
	"timezone_abbreviation": "UTC",
	"elevation": 141,
	"current_units": {
	  "time": "iso8601",
	  "interval": "seconds",
	  "temperature_2m": "°C"
	},
	"current": {
	  "time": "2000-01-01T01:02",
	  "interval": 900,
	  "temperature_2m": -0.9
	}
}`

const testResponseFailed = `{
	""
}`

func TestForecastOk(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, testResponse)
	}))
	defer svr.Close()

	cl := New(svr.URL)

	resp, err := cl.Forecast(ForecastParams{
		Latitude:  0,
		Longitude: 0,
		Timezone:  "UTC",
	})

	if err != nil {
		t.Errorf("getting unexpected error: %v", err)
	}

	require.Equal(t, -0.9, resp)
}

func TestForecastJsonUnmarshalFailed(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, testResponseFailed)
	}))
	defer svr.Close()

	cl := New(svr.URL)

	_, err := cl.Forecast(ForecastParams{
		Latitude:  0,
		Longitude: 0,
		Timezone:  "UTC",
	})

	if err == nil {
		t.Errorf("error was expected, got nil")
		return
	}

	require.Equal(t, "invalid character '}' after object key", err.Error())
}

func TestForecastServerFailed(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer svr.Close()

	cl := New(svr.URL)

	_, err := cl.Forecast(ForecastParams{
		Latitude:  0,
		Longitude: 0,
		Timezone:  "UTC",
	})

	if err == nil {
		t.Errorf("error was expected, got nil")
		return
	}

	require.Equal(t, "open-meteo service error", err.Error())
}

func TestForecastServerPanic(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("panic"))
		panic("service error")
	}))
	defer svr.Close()

	cl := New(svr.URL)

	_, err := cl.Forecast(ForecastParams{
		Latitude:  0,
		Longitude: 0,
		Timezone:  "UTC",
	})

	if err == nil {
		t.Errorf("error was expected, got nil")
		return
	}

	errText := strings.TrimPrefix(err.Error(), "Get \""+svr.URL)

	fmt.Println(svr.URL, errText)

	expectedError := `/v1/forecast?latitude=0.000000&longitude=0.000000&current=temperature_2m&hourly=temperature_2m&timezone=UTC": EOF`

	require.Equal(t, expectedError, errText)
}
