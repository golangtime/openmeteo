# Openmeteo API client library

### Description

Non official client for Openmeteo API

### Example usage

```go
package main

import (
	"fmt"

	"github.com/golangtime/openmeteo"
)

func main() {
	// create new openmeteo client instance
	client := openmeteo.New("")

	latitude, longitude := float32(0), float32(0)

	// call the forecast method to get the current weather
	temperature, err := client.Forecast(openmeteo.ForecastParams{
		Latitude:  latitude,
		Longitude: longitude,
		Timezone:  timezone,
	})
	if err != nil {
		panic(err)
	}

	fmt.Println(temperature)
}
```

### Running example

Run the following command

```bash
go run .\examples\forecast\main.go --lat=55.7522 --lon=37.6156
```

### API Documentation

https://open-meteo.com/en/docs
