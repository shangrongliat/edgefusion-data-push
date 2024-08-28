package influx

import (
	"context"
	"fmt"
	"time"

	"github.com/influxdata/influxdb-client-go/v2/api"

	log "edgefusion-data-push/plugin/logs"
	influxdb "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api/write"
)

// OZKfwLM43c_115CVH4BWm4vMwVcsPDrDz6hNlaHRK5wUjre7xz40k1sOJ0b4E4cu76awoHytAXMjo99R9djjFQ==

const (
	InfluxDB_org    = "edgefusion"
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

	url := "http://172.16.100.14:8086"
	client := influxdb.NewClient(url, "MDOZ9TWdnMlqUTKwvCg_IhjAtR3x1m8Ssjjqm7mLkjAJFNF_nSM0m52VN5oLmTL3TQN7yYV8o_Gf5FF6v9jKkw==")
	//_, err := client.Setup(context.Background(), "admin", "influxadmin", InfluxDB_org, InfluxDB_bucket, 0)
	//if err != nil {
	//	log.L().Error("初始化失败", log.Error(err))
	//}
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
	           |> range(start: -10h)
				|> filter(fn: (r) => r._measurement == "%s")`, InfluxDB_bucket, measurement)
	//qStr := fmt.Sprintf(`SELECT * FROM "%s"`, measurement)

	res, err := i.queryAPI.Query(context.Background(), query)
	if err != nil {
		log.L().Error("查询异常.", log.Error(err))
		return err
	}
	for res.Next() {
		if res.TableChanged() {
			fmt.Printf("表：%s\n", res.TableMetadata().String())
		}
		value := res.Record().Value()

		fmt.Println("----------", value)
		start := res.Record().Start()
		stop := res.Record().Stop()
		fmt.Println(start, "---------", stop)
	}

	return nil
}
