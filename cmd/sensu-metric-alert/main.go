package main

import (
	"bufio"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/danielb42/whiteflag"
)

var (
	prmEndpoint string
	prmMetric   string
)

func main() {
	log.SetOutput(os.Stderr)

	whiteflag.Alias("e", "endpoint", "endpoint to scrape (http://url:port/path)")
	whiteflag.Alias("m", "metric", "the metric to watch")
	whiteflag.Alias("l", "lt", "value must be lower than this for crit status")
	whiteflag.Alias("g", "gt", "value must be greater than this for crit status")
	whiteflag.Alias("q", "eq", "value must be equal to this for crit status")
	whiteflag.Alias("n", "ne", "value must be different from this for crit status")
	whiteflag.ParseCommandLine()

	if !whiteflag.CheckString("e") || !whiteflag.CheckString("m") ||
		(!whiteflag.CheckInt("l")) && !whiteflag.CheckInt("g") && !whiteflag.CheckInt("q") && !whiteflag.CheckInt("n") {
		println("usage: sensu-metric-alert -e <endpoint> -m <metric> --lt|--gt|--eq|--ne <crit value>")
		os.Exit(2)
	}

	prmEndpoint = whiteflag.GetString("e")
	prmMetric = whiteflag.GetString("m")

	resp, err := http.Get(prmEndpoint)
	if err != nil {
		log.Println("could not scrape metrics from", prmEndpoint, err.Error())
		os.Exit(2)
	}
	defer resp.Body.Close()

	foundMetric := false
	s := bufio.NewScanner(resp.Body)
	for s.Scan() {
		line := strings.Split(s.Text(), " ")

		if strings.HasPrefix(line[0], prmMetric) {
			foundMetric = true
			val, err := strconv.ParseFloat(line[1], 64)
			if err != nil {
				log.Println(err)
				os.Exit(2)
			}

			evaluate(val)
		}
	}

	if !foundMetric {
		log.Println("metric", prmMetric, "not found in endpoint output")
		os.Exit(2)
	}
}

func evaluate(val float64) {
	log.Printf("%s = %f\n", prmMetric, val)

	if whiteflag.CheckInt("lt") && val < float64(whiteflag.GetInt("lt")) {
		os.Exit(2)
	}

	if whiteflag.CheckInt("gt") && val > float64(whiteflag.GetInt("gt")) {
		os.Exit(2)
	}

	if whiteflag.CheckInt("eq") && val == float64(whiteflag.GetInt("eq")) {
		os.Exit(2)
	}

	if whiteflag.CheckInt("ne") && val != float64(whiteflag.GetInt("ne")) {
		os.Exit(2)
	}

	os.Exit(0)
}
