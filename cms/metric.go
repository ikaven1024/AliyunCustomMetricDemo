package cms

// @see https://github.com/aliyun/aliyun-cms-java-sdk/blob/master/src/main/java/com/aliyun/openservices/cms/model/CustomMetric.java

const (
	typeValue = 0
	typeAgg = 1
)

type CustomMetric struct {
	// app groupId
	// you can find it on aliyun monitor console
	GroupId int `json:"groupId"`
	// key value pair
	// which defined one time series
	// the key must [_0-9a-zA-Z]
	Dimensions map[string] string `json:"dimensions"`
	// metric name
	// the value must be [_0-9a-zA-Z]
	MetricName string `json:"metricName"`
	// the time of the value set
	Time Time `json:"time"`
	// 0 original value
	// 1 pre agg value
	Type int `json:"type"`
	// only when type = 1
	// only set the values in
	//     * 15
	//     * 60
	//     * 300
	//     * 900
	//     * 1800
	Period int `json:"period"`

    // if type is 0
    // key: 'value'
    // value: original value
    //
    // if type is 1
    // key set is
    // Average,Minimum,Maximum,Sum,SampleCount
    // LastValue: the last of the original value
    // SumPerSecond: Sum/period or the ewma(or other algorithm) value of the original value
    // CountPerSecond: SampleCount/period or the ewma(or other algorithm) value of the original value
    //
    // PXX value: xx% is less then the PXX
    // P10: 10% of the  original value is less then the P10
    // P20
    // P30
    // P40
    // P50
    // P60
    // P70
    // P75
    // P80
    // P90
    // P95
    // P98
    // P99
    // the value is the pre agg value in the period
	Values map[string] interface{} `json:"values"`
}


