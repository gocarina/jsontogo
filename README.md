JsonToGo
========

The JsonToGo package aims to transform a JSON element to its Golang struct representation

[![GoDoc](https://godoc.org/github.com/gocarina/jsontogo?status.png)](https://godoc.org/github.com/gocarina/jsontogo)

Usage
=====

```go

package main

import (
	"bytes"
	"fmt"
	"jsontogo"
)

const (
	JSON_TEST = `
    {
        "coord": {
            "lon": -0.13,
            "lat": 51.51
        },
        "sys_data": {
            "message": 0.0052,
            "country": "GB",
            "sunrise": 1401335451,
            "sunset": 1401393904
        },
        "weather": [
            {
                "id": 803,
                "main": "Clouds",
                "description": "broken clouds",
                "icon": "04d"
            }
        ]
    }
    `
)

func main() {
	stringWriter := &bytes.Buffer{}
	enc := jsontogo.NewEncoderWithNameAndTags(stringWriter, "Weather", []string{"json"})
	if err := enc.Encode([]byte(JSON_TEST)); err != nil {
		panic(err)
	}
	fmt.Println(stringWriter)
}

```

Will output

```go

type Weather struct {
	Weather []*struct {
		Icon string `json:"icon"`
		Id float64 `json:"id"`
		Main string `json:"main"`
		Description string `json:"description"`
	} `json:"weather"`
	Coord struct {
		Lon float64 `json:"lon"`
		Lat float64 `json:"lat"`
	} `json:"coord"`
	SysData struct {
		Message float64 `json:"message"`
		Country string `json:"country"`
		Sunrise float64 `json:"sunrise"`
		Sunset float64 `json:"sunset"`
	} `json:"sys_data"`
}

```
