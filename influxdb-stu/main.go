package main

import (
	"context"
	"fmt"
	"log"
	"time"

	influxdb "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api/write"
)

func connInflux() influxdb.Client {
	return influxdb.NewClient("http://localhost:8086", "YuE8PZ5uxwXFa4rlmXZWFKOPdiPnmYaE_G3sm9Y5_iajb9oCTQWc8d7c34TorpA1oVuvuMWPyxn7Nk3ijJvIlw==")
}

func main() {
	client := connInflux()
	org := "lbw"
	bucket := "test"
	writeAPI := client.WriteAPIBlocking(org, bucket)
	for value := 0; value < 5; value++ {
		tags := map[string]string{
			"tagname1": "tagvalue1",
		}
		fields := map[string]interface{}{
			"field1": value,
		}
		point := write.NewPoint("measurement1", tags, fields, time.Now())
		time.Sleep(1 * time.Second) // separate points by 1 second

		if err := writeAPI.WritePoint(context.Background(), point); err != nil {
			log.Fatal(err)
		}
	}

	queryAPI := client.QueryAPI(org)
	query := `from(bucket: "<BUCKET>")
            |> range(start: -10m)
            |> filter(fn: (r) => r._measurement == "measurement1")`
	results, err := queryAPI.Query(context.Background(), query)
	if err != nil {
		log.Fatal(err)
	}
	for results.Next() {
		fmt.Println(results.Record())
	}
	if err := results.Err(); err != nil {
		log.Fatal(err)
	}
}
