package main

import (
	"bufio"
	"fmt"
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
	prmPrintAll bool
)

func main() {
	whiteflag.Alias("e", "endpoint", "endpoint to scrape (http://url:port/path)")
	whiteflag.Alias("m", "metric", "the metric to watch")
	whiteflag.Alias("l", "lt", "value must be lower than this for crit status")
	whiteflag.Alias("g", "gt", "value must be greater than this for crit status")
	whiteflag.Alias("q", "eq", "value must be equal to this for crit status")
	whiteflag.Alias("n", "ne", "value must be different from this for crit status")
	whiteflag.Alias("p", "printall", "print all scraped metrics (without threshold evaluation)")

	if !whiteflag.FlagPresent("e") || !whiteflag.FlagPresent("m") ||
		(!whiteflag.FlagPresent("l")) && !whiteflag.FlagPresent("g") && !whiteflag.FlagPresent("q") && !whiteflag.FlagPresent("n") {
		println("usage: sensu-metric-alert -e <endpoint> -m <metric> --lt|--gt|--eq|--ne <crit value>")
		os.Exit(2)
	}

	prmEndpoint = whiteflag.GetString("e")
	prmMetric = whiteflag.GetString("m")
	prmPrintAll = whiteflag.FlagPresent("p")

	resp, err := http.Get(prmEndpoint)
	if err != nil {
		log.Println("could not scrape metrics from", prmEndpoint, err.Error())
		os.Exit(2)
	}
	defer resp.Body.Close()

	s := bufio.NewScanner(resp.Body)
	for s.Scan() {
		line := strings.Split(s.Text(), " ")

		scrapedMetric := line[0]
		scrapedVal := line[1]

		if prmPrintAll {
			fmt.Println(scrapedMetric, scrapedVal)
		} else if strings.HasPrefix(scrapedMetric, prmMetric) {
			val, err := strconv.ParseFloat(scrapedVal, 64)
			if err != nil {
				log.Println(err)
				os.Exit(2)
			}

			evaluate(val)
		}
	}

	if !prmPrintAll {
		log.Println("metric", prmMetric, "not found in endpoint output")
		os.Exit(2)
	}
}

func evaluate(val float64) {
	log.Printf("%s = %f\n", prmMetric, val)

	if whiteflag.FlagPresent("lt") && val < float64(whiteflag.GetInt("lt")) {
		log.Println("Value should be >=", float64(whiteflag.GetInt("lt")))
		os.Exit(2)
	}

	if whiteflag.FlagPresent("gt") && val > float64(whiteflag.GetInt("gt")) {
		log.Println("Value should be <=", float64(whiteflag.GetInt("gt")))
		os.Exit(2)
	}

	if whiteflag.FlagPresent("eq") && val == float64(whiteflag.GetInt("eq")) {
		log.Println("Value should be !=", float64(whiteflag.GetInt("eq")))
		os.Exit(2)
	}

	if whiteflag.FlagPresent("ne") && val != float64(whiteflag.GetInt("ne")) {
		log.Println("Value should be =", float64(whiteflag.GetInt("ne")))
		os.Exit(2)
	}

	log.Println("OK: Value in expected range")
	os.Exit(0)
}
