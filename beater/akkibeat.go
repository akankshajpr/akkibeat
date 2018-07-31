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

var Loc = "159205000000025039"

type  Mydata []struct {
	Message	string	`json:"message"`
	Data struct {
		ChartData struct {
			ResponseTimeReportChart struct {
				Location struct {
					Max float64 `json:"max"`
					Label string `json:"label"`
					Min float64 `json:"min"`
					Avg float64 `json:"average"`
				} `json:Loc`
			} `json:"ResponseTimeReportChart"`
		}`json:"chart_data"`
	}`json:"data"`
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

	ticker := time.NewTicker(bt.config.Period)
	counter := 1
	for {
		select {
		case <-bt.done:
			return nil
		case <-ticker.C:
		}

		client := &http.Client{}
		req, err := http.NewRequest("GET", bt.config.Url, nil)
		req.Header.Add("Authorization", bt.config.Authorization)

		resp, err := client.Do(req)
		if err != nil {
				return  err
		}
		defer resp.Body.Close()

		var akkidata Mydata
		if resp.StatusCode == http.StatusOK {
				bodyBytes, err2 := ioutil.ReadAll(resp.Body)
				if err2 != nil {
						return  err
				}
				 json.Unmarshal(bodyBytes, &akkidata)
		}
		
		fmt.Println(akkidata)
		for d := range akkidata {

			event := beat.Event{
				Timestamp: time.Now(),
				Fields: common.MapStr{
				"@timestamp":  common.Time(time.Now()),	
				"type":	"akkibeat",
				"counter": counter,
				"Message": akkidata[d].Message,
				"MaxResponseTime": akkidata[d].Data.ChartData.ResponseTimeReportChart.Location.Max,
				"Label": akkidata[d].Data.ChartData.ResponseTimeReportChart.Location.Label,
				"MinResponseTime": akkidata[d].Data.ChartData.ResponseTimeReportChart.Location.Min,
				"AvgResponseTime": akkidata[d].Data.ChartData.ResponseTimeReportChart.Location.Avg,
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
