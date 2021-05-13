package main

import (
	"crypto/tls"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type ReceiverFunc func(key string, value float64, indices []int, gaugeVecs map[string]*prometheus.GaugeVec)

func (receiver ReceiverFunc) Receive(key string, value float64, indices []int, gaugeVecs map[string]*prometheus.GaugeVec) {
	receiver(key, value, indices, gaugeVecs)
}

type Receiver interface {
	Receive(key string, value float64, indices []int, gaugeVecs map[string]*prometheus.GaugeVec)
}

func WalkJSON(path string, jsonData interface{}, indices []int, gaugeVecs map[string]*prometheus.GaugeVec, receiver Receiver) {
	switch v := jsonData.(type) {
	case int:
		receiver.Receive(path, float64(v), indices, gaugeVecs)
	case float64:
		receiver.Receive(path, v, indices, gaugeVecs)
	case bool:
		n := 0.0
		if v {
			n = 1.0
		}
		receiver.Receive(path, n, indices, gaugeVecs)
	case string:
		// ignore
	case nil:
		// ignore
	case []interface{}:
		prefix := ""
		if path != "" {
			prefix = path + "::"
		}
		indicesNext := make([]int, len(indices)+1)
		copy(indicesNext, indices)
		for i, x := range v {
			indicesNext[len(indices)] = i
			WalkJSON(fmt.Sprintf("%sarray_%d", prefix, len(indices)), x, indicesNext, gaugeVecs, receiver)
		}
	case map[string]interface{}:
		prefix := ""
		if path != "" {
			prefix = path + "::"
		}
		for k, x := range v {
			WalkJSON(fmt.Sprintf("%s%s", prefix, k), x, indices, gaugeVecs, receiver)
		}
	default:
		log.Printf("unkown type: %#v", v)
	}
}

func doProbe(client *http.Client, target string) (interface{}, error) {
	resp, err := client.Get(target)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var jsonData interface{}
	err = json.Unmarshal([]byte(bytes), &jsonData)
	if err != nil {
		return nil, err
	}

	return jsonData, nil
}

var httpClient *http.Client

func init() {
	httpClient = &http.Client{
		Transport: &http.Transport{
			MaxIdleConns: 100,
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}
}

func doWalkJSON(prefix string, jsonData interface{}, registry *prometheus.Registry) {
	WalkJSON(prefix, jsonData, []int{}, map[string]*prometheus.GaugeVec{}, ReceiverFunc(func(key string, value float64, indices []int, gaugeVecs map[string]*prometheus.GaugeVec) {
		g, ok := gaugeVecs[key]
		if !ok {
			labels := make([]string, len(indices))
			for array, _ := range indices {
				labels[array] = fmt.Sprintf("array_%d_index", array)
			}
			g = prometheus.NewGaugeVec(
				prometheus.GaugeOpts{
					Name: key,
					Help: "Retrieved value",
				},
				labels,
			)
			gaugeVecs[key] = g
			registry.MustRegister(g)
		}
		labelsWithValues := prometheus.Labels{}
		for array, index := range indices {
			labelsWithValues[fmt.Sprintf("array_%d_index", array)] = strconv.Itoa(index)
		}
		g.With(labelsWithValues).Set(value)
	}))
}

func probeHandler(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()

	target := params.Get("target")
	if target == "" {
		http.Error(w, "Target parameter is missing", http.StatusBadRequest)
		return
	}

	prefix := params.Get("prefix")

	jsonData, err := doProbe(httpClient, target)
	if err != nil {
		log.Print(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	// log.Printf("Retrieved value %v", jsonData)

	registry := prometheus.NewRegistry()

	doWalkJSON(prefix, jsonData, registry)

	h := promhttp.HandlerFor(registry, promhttp.HandlerOpts{})
	h.ServeHTTP(w, r)
}

var indexHTML = []byte(`<html>
<head><title>Json Exporter</title></head>
<body>
<h1>Json Exporter</h1>
<p><a href="/probe">Run a probe</a></p>
<p><a href="/metrics">Metrics</a></p>
</body>
</html>`)

func main() {
	addr := flag.String("listen-address", ":9116", "The address to listen on for HTTP requests.")
	flag.Parse()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write(indexHTML)
	})
	http.HandleFunc("/probe", probeHandler)
	http.Handle("/metrics", promhttp.Handler())

	log.Printf("listenning on %s", *addr)
	log.Fatal(http.ListenAndServe(*addr, nil))
}
