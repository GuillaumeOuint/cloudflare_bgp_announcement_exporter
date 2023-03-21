package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type JSONResponse struct {
	Success bool            `json:"success"`
	Errors  []interface{}   `json:"errors"`
	Result  ResultStructure `json:"result"`
}

type ResultStructure struct {
	Serie0 Serie0Structure `json:"serie_0"`
	Meta   MetaStructure   `json:"meta"`
}

type Serie0Structure struct {
	Timestamps []time.Time `json:"timestamps"`
	Values     []string    `json:"values"`
}

type MetaStructure struct {
	DateRange      DateRangeStructure      `json:"dateRange"`
	ConfidenceInfo ConfidenceInfoStructure `json:"confidenceInfo"`
	Normalization  string                  `json:"normalization"`
	AggInterval    string                  `json:"aggInterval"`
	LastUpdated    time.Time               `json:"lastUpdated"`
}

type DateRangeStructure struct {
	StartTime time.Time `json:"startTime"`
	EndTime   time.Time `json:"endTime"`
}

type ConfidenceInfoStructure struct {
	Level       interface{} `json:"level"`
	Annotations []string    `json:"annotations"`
}

var (
	apiResponseGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "bgp_dfz_announcements",
			Help: "API response values.",
		},
		[]string{"asn"},
	)
)

func init() {
	prometheus.MustRegister(apiResponseGauge)
}

func fetchApiData() {
	// make http client with headers for authorization bearer token
	client := &http.Client{}
	asns := os.Getenv("ASN")
	splitAsns := strings.Split(asns, ",")
	for {
		for _, asn := range splitAsns {
			req, err := http.NewRequest("GET", "https://api.cloudflare.com/client/v4/radar/bgp/timeseries?asn="+asn+"&dateRange=1d", nil)
			if err != nil {
				fmt.Println("Error creating request:", err)
			}
			req.Header.Add("Authorization", "Bearer "+os.Getenv("CLOUDFLARE_API_TOKEN"))
			// get response from api
			resp, err := client.Do(req)
			if err != nil {
				fmt.Println("Error fetching data:", err)
				continue
			}

			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				fmt.Println("Error reading response body:", err)
				continue
			}
			var apiResponse JSONResponse
			err = json.Unmarshal(body, &apiResponse)
			if err != nil {
				fmt.Println("Error unmarshalling response:", err)
				continue
			}
			apiVal := apiResponse.Result.Serie0.Values
			value := apiVal[len(apiVal)-1]
			valueFloat, err := strconv.ParseFloat(value, 64)
			if err != nil {
				fmt.Println("Error converting value to float:", err)
				time.Sleep(1 * time.Minute)
				continue
			}
			apiResponseGauge.WithLabelValues(asn).Set(valueFloat)

			resp.Body.Close()
		}
		time.Sleep(1 * time.Minute)
	}
}

func main() {
	go fetchApiData()
	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":8080", nil)
}
