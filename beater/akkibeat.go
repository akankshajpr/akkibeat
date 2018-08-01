package beater

import (
	"fmt"
	"time"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"github.com/elastic/beats/libbeat/beat"
	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/logp"

	"github.com/akankshajpr/akkibeat/config"
)

//var Loc = "6"

type  Mydata []struct {	
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    struct {
		TableData struct {
		} `json:"table_data"`
		Info struct {
			FormattedEndTime       string `json:"formatted_end_time"`
			MonitorType            string `json:"monitor_type"`
			ResourceID             string `json:"resource_id"`
			ResourceTypeName       string `json:"resource_type_name"`
			PeriodName             string `json:"period_name"`
			GeneratedTime          string `json:"generated_time"`
			MetricAggregationName  string `json:"metric_aggregation_name"`
			ReportName             string `json:"report_name"`
			EndTime                string `json:"end_time"`
			MetricAggregation      int    `json:"metric_aggregation"`
			StartTime              string `json:"start_time"`
			SegmentType            int    `json:"segment_type"`
			ReportType             int    `json:"report_type"`
			Period                 int    `json:"period"`
			ResourceName           string `json:"resource_name"`
			FormattedStartTime     string `json:"formatted_start_time"`
			FormattedGeneratedTime string `json:"formatted_generated_time"`
			ResourceType           int    `json:"resource_type"`
		} `json:"info"`
		ChartData []struct {
			ResponseTimeReportChart   []interface{} `json:"ResponseTimeReportChart,omitempty"`
			LocationResponseTimeChart []struct {
				Num1 struct {
					Max             []float64       `json:"max"`
					Label           string          `json:"label"`
					Min             []float64       `json:"min"`
					Nine5Percentile []float64       `json:"95_percentile"`
					Average         []float64       `json:"average"`
					ChartData       [][]interface{} `json:"chart_data"`
				} `json:"1,omitempty"`
				Num6 struct {
					Max             []float64       `json:"max"`
					Label           string          `json:"label"`
					Min             []float64       `json:"min"`
					Nine5Percentile []float64       `json:"95_percentile"`
					Average         []float64       `json:"average"`
					ChartData       [][]interface{} `json:"chart_data"`
				} `json:"6,omitempty"`
				Num15 struct {
					Max             []float64       `json:"max"`
					Label           string          `json:"label"`
					Min             []float64       `json:"min"`
					Nine5Percentile []float64       `json:"95_percentile"`
					Average         []float64       `json:"average"`
					ChartData       [][]interface{} `json:"chart_data"`
				} `json:"15,omitempty"`
			} `json:"LocationResponseTimeChart,omitempty"`
		} `json:"chart_data"`
	} `json:"data"`
}
}

// Akkibeat configuration.
type Akkibeat struct {
	done   chan struct{}
	config config.Config
	client beat.Client
}

// New creates an instance of akkibeat.
func New(b *beat.Beat, cfg *common.Config) (beat.Beater, error) {
	c := config.DefaultConfig
	if err := cfg.Unpack(&c); err != nil {
		return nil, fmt.Errorf("Error reading config file: %v", err)
	}

	bt := &Akkibeat{
		done:   make(chan struct{}),
		config: c,
	}
	return bt, nil
}

// Run starts akkibeat.
func (bt *Akkibeat) Run(b *beat.Beat) error {
	logp.Info("akkibeat is running! Hit CTRL-C to stop it.")

	var err error
	bt.client, err = b.Publisher.Connect()
	if err != nil {
		return err
	}
        fmt.Println("bt.config.Url", bt.config.Url)

	ticker := time.NewTicker(bt.config.Period)
	counter := 1
	for {
		select {
		case <-bt.done:
			return nil
		case <-ticker.C:
		}

		client := &http.Client{}
		req, err := http.NewRequest("GET",bt.config.Url,nil)
		req.Header.Add("Authorization", bt.config.Authorization)

		resp, err := client.Do(req)
		if err != nil {
				return  err
		}
		defer resp.Body.Close()
		fmt.Println(resp.Body)
		var akkidata Mydata
		if resp.StatusCode == http.StatusOK {
				bodyBytes, err2 := ioutil.ReadAll(resp.Body)
				//abc := string(bodyBytes)
				
					if err2 != nil {
						return  err
				}
				 json.Unmarshal(bodyBytes, &akkidata)
				fmt.Println(string(bodyBytes))		

}
		fmt.Println("new")
fmt.Println(akkidata[:])		
		fmt.Println(akkidata)
		for d := range akkidata {

			event := beat.Event{
				Timestamp: time.Now(),
				Fields: common.MapStr{
				"@timestamp":  common.Time(time.Now()),	
				"type":	"akkibeat",
				"counter": counter,
				"Message": akkidata[d].Message,
				"MaxResponseTime": akkidata[d].Data.ChartData[0].LocationResponseTimeChart[0].Num6.Max[0],
				"Label": akkidata[d].Data.ChartData[0].LocationResponseTimeChart[0].Num6.Label,
				"MinResponseTime": akkidata[d].Data.ChartData[0].LocationResponseTimeChart[0].Num6.Min[0],
				"AvgResponseTime": akkidata[d].Data.ChartData[0].LocationResponseTimeChart[0].Num6.Average[0],
			},
		}
		bt.client.Publish(event)
		logp.Info("Event sent")
	}
	counter++
	}
}
// Stop stops akkibeat.
func (bt *Akkibeat) Stop() {
	bt.client.Close()
	close(bt.done)
}
