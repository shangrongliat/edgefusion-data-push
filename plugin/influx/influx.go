package influx

import (
	"context"
	"fmt"
	"time"

	log "edgefusion-data-push/plugin/logs"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
	"github.com/influxdata/influxdb-client-go/v2/api/write"
)

// OZKfwLM43c_115CVH4BWm4vMwVcsPDrDz6hNlaHRK5wUjre7xz40k1sOJ0b4E4cu76awoHytAXMjo99R9djjFQ==

const (
	InfluxDB_org    = "arboo"
	InfluxDB_bucket = "ef-data"
)

type InfluxRepo interface {
	Save(nodeId, appId string, fields map[string]interface{}) error
	Get(nodeId, appId string) error
}

type Influx struct {
	writeAPI api.WriteAPIBlocking
	queryAPI api.QueryAPI
}

func NewInflux() (InfluxRepo, error) {
	//token := os.Getenv("INFLUXDB_TOKEN")
	token := "OZKfwLM43c_115CVH4BWm4vMwVcsPDrDz6hNlaHRK5wUjre7xz40k1sOJ0b4E4cu76awoHytAXMjo99R9djjFQ=="
	url := "http://172.16.100.14:8086"
	client := influxdb2.NewClient(url, token)
	writeAPI := client.WriteAPIBlocking(InfluxDB_org, InfluxDB_bucket)
	queryAPI := client.QueryAPI(InfluxDB_org)
	return &Influx{
		writeAPI: writeAPI,
		queryAPI: queryAPI,
	}, nil
}

func (i *Influx) Save(nodeId, appId string, fields map[string]interface{}) error {
	// 索引列
	tags := map[string]string{
		"node_id": nodeId,
		"app_id":  appId,
	}
	measurement := fmt.Sprintf("%s-%s", nodeId, appId)
	point := write.NewPoint(measurement, tags, fields, time.Now())
	//time.Sleep(1 * time.Second) // separate points by 1 second
	if err := i.writeAPI.WritePoint(context.Background(), point); err != nil {
		log.L().Error("写入失败.", log.Error(err))
		return err
	}
	return nil
}

func (i *Influx) Get(nodeId, appId string) error {
	// 索引列
	//measurement := fmt.Sprintf("%s-%s", nodeId, appId)
	measurement := "FDt4zjxNrTnohMt3-skills-test-004"
	// Query with parameters
	query := fmt.Sprintf(`from(bucket:"%s")
                |> range(start: -1h)
				|> filter(fn: (r) => r._measurement == "%s")`, InfluxDB_bucket, measurement)

	result, err := i.queryAPI.Query(context.Background(), query)
	if err != nil {
		log.L().Error("查询异常.", log.Error(err))
		return err
	}
	if result.Err() != nil {
		log.L().Error("query parsing error", log.Error(result.Err()))
		return err
	}
	// 遍历查询结果
	for result.Next() {
		//mes := result.Record().Measurement()
		Class := result.Record().ValueByKey("Class")
		Name := result.Record().ValueByKey("Name")
		Box := result.Record().ValueByKey("Box")
		Score := result.Record().ValueByKey("Score")
		Location := result.Record().ValueByKey("Location")
		value := result.Record().Value()
		field := result.Record().Field()
		fmt.Printf("class: %v; name: %v; box: %v; score: %v; location: %v \n", Class, Name, Box, Score, Location)
		fmt.Println("value ===== ", value)
		fmt.Println("field ===== ", field)
	}
	return nil
}
