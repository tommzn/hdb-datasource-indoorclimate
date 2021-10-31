
[![Go Reference](https://pkg.go.dev/badge/github.com/tommzn/hdb-datasource-indoorclimate.svg)](https://pkg.go.dev/github.com/tommzn/hdb-datasource-indoorclimate)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/tommzn/hdb-datasource-indoorclimate)
![GitHub release (latest by date)](https://img.shields.io/github/v/release/tommzn/hdb-datasource-indoorclimate)
[![Go Report Card](https://goreportcard.com/badge/github.com/tommzn/hdb-datasource-indoorclimate)](https://goreportcard.com/report/github.com/tommzn/hdb-datasource-indoorclimate)
[![Actions Status](https://github.com/tommzn/hdb-datasource-indoorclimate/actions/workflows/go.pkg.auto-ci.yml/badge.svg)](https://github.com/tommzn/hdb-datasource-indoorclimate/actions)

# HomeDashboard Indoor Climate DataSource
Fetches indoor climate data from a MQTT broker and publishes to HomebDaskboard backend.

## Config
Config have to contain URL of an exchange rate API and a list of currencie pairs an exchange rate should be fetched.
More details about loading config at https://github.com/tommzn/go-config

### Config example
```yaml
exchangerate:
  url: "https://api.frankfurter.app/latest"
  date_format: "2006-01-02"
  conversions:
    - from: "EUR"
      to: "USD"
    - from: "USD"
      to: "EUR"
```

## Usage
After creating a new datasource, you can fetch specified exchange rates. If anything works well Fetch will return a [exchange rates struct](https://github.com/tommzn/hdb-events-go/blob/main/exchangerate.pb.go) or otherwise an error.
```golang

    import (
       exchangerate "github.com/tommzn/hdb-datasource-exchangerate"  
       events "github.com/tommzn/hdb-events-go"  
    )
    
    datasource, err := exchangerate.New(config)
    if err != nil {
        panic(err)
    }

    weatherData, err := datasource.Fetch()
    if err != nil {
        panic(err)
    }

    fmt.Printf("Exchange Rates: %d\n", len(weatherData.(events.ExchangeRates).Rates))
```
