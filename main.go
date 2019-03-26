package main

import (
	"aliyun-metric/cms"
	"log"
)

func main() {
	client := cms.NewClient("xxx", "xxxx")
	resp, err := client.PutCustomMetrics([]cms.CustomMetric{
		{GroupId: 102, MetricName: "testMetric", Period: 15, Time: cms.Now(), Type: 1,
			Dimensions: map[string]string{"ip": "127.0.0.1", "key": "value"},
			Values: map[string]interface{} {"LastValue": 100}},
	})

	if err != nil {
		log.Panic(err)
	}
	log.Println(resp.Code, resp.Message, resp.RequestId)
}
