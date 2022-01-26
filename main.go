package main

import (
	"bufio"
	"flag"
	"log"
	"net/http"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	disk_usage_remaining = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "disk_usage_remaining",
			Help: "Disk Usage Remaining",
		},
		[]string{"physicaldrive"},
	)
	disk_status = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "disk_status",
			Help: "Disk Status (OK = 1)",
		},
		[]string{"physicaldrive"},
	)
	disk_current_temperature = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "disk_current_temperature",
			Help: "Disk Current Temperature",
		},
		[]string{"physicaldrive"},
	)
	disk_maximum_temperature = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "disk_maximum_temperature",
			Help: "Disk Maximum Temperature",
		},
		[]string{"physicaldrive"},
	)
	disk_power_on_hours = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "disk_power_on_hours",
			Help: "Disk Power on Hours ",
		},
		[]string{"physicaldrive"},
	)
	disk_estimated_life_remaining = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "disk_estimated_life_remaining",
			Help: "SSD Disk Estimated Life Remaining",
		},
		[]string{"physicaldrive"},
	)
	usageRG  *regexp.Regexp = regexp.MustCompile("^.*: (.*)%")
	dayRG    *regexp.Regexp = regexp.MustCompile("^.*date: (.*) ")
	bayRG    *regexp.Regexp = regexp.MustCompile("^.*Bay: (.*)")
	boxRG    *regexp.Regexp = regexp.MustCompile("^.*Box: (.*)")
	statusRG *regexp.Regexp = regexp.MustCompile("         Status: (.*)")
	ctempRG  *regexp.Regexp = regexp.MustCompile("^.*Current Temperature.*: (.*)")
	mtempRG  *regexp.Regexp = regexp.MustCompile("^.*Maximum Temperature.*: (.*)")
	powerRG  *regexp.Regexp = regexp.MustCompile("^.*Power On Hours: (.*)")
	typeRG   *regexp.Regexp = regexp.MustCompile("^.*Interface Type: (.*)")
)

func init() {
	prometheus.MustRegister(disk_usage_remaining)
	prometheus.MustRegister(disk_status)
	prometheus.MustRegister(disk_current_temperature)
	prometheus.MustRegister(disk_maximum_temperature)
	prometheus.MustRegister(disk_power_on_hours)
	prometheus.MustRegister(disk_estimated_life_remaining)
}

func parse() {
	bay_current := "1"
	box_current := "1"
	disktype_current := "none"
	var MetricValue float64 = 0
	output := string(runcmd("ssacli ctrl slot=0 physicaldrive all show detail", true))

	scanner := bufio.NewScanner(strings.NewReader(output))
	for scanner.Scan() {

		bay := bayRG.FindStringSubmatch(scanner.Text())
		box := boxRG.FindStringSubmatch(scanner.Text())
		usage := usageRG.FindStringSubmatch(scanner.Text())
		day := dayRG.FindStringSubmatch(scanner.Text())
		status := statusRG.FindStringSubmatch(scanner.Text())
		disktype := typeRG.FindStringSubmatch(scanner.Text())
		ctmep := ctempRG.FindStringSubmatch(scanner.Text())
		mtemp := mtempRG.FindStringSubmatch(scanner.Text())
		power := powerRG.FindStringSubmatch(scanner.Text())

		if len(bay) != 0 {
			bay_current = bay[1]
		}
		if len(box) != 0 {
			box_current = box[1]
		}
		if len(disktype) != 0 {
			disktype_current = disktype[1]
		}
		if len(usage) != 0 {
			name := "box " + box_current + " bay " + bay_current + " type " + disktype_current
			MetricValue, _ := strconv.ParseFloat(strings.TrimSpace(usage[1]), 64)
			disk_usage_remaining.With(prometheus.Labels{"physicaldrive": name}).Set(MetricValue)
		}
		if len(day) != 0 {
			name := "box " + box_current + " bay " + bay_current + " type " + disktype_current
			MetricValue, _ := strconv.ParseFloat(strings.TrimSpace(day[1]), 64)
			disk_estimated_life_remaining.With(prometheus.Labels{"physicaldrive": name}).Set(MetricValue)
		}
		if len(status) != 0 {
			name := "box " + box_current + " bay " + bay_current + " type " + disktype_current
			if status[1] == "OK" {
				MetricValue = 1
			}
			disk_status.With(prometheus.Labels{"physicaldrive": name}).Set(MetricValue)
		}
		if len(ctmep) != 0 {
			name := "box " + box_current + " bay " + bay_current + " type " + disktype_current
			MetricValue, _ := strconv.ParseFloat(strings.TrimSpace(ctmep[1]), 64)
			disk_current_temperature.With(prometheus.Labels{"physicaldrive": name}).Set(MetricValue)
		}
		if len(mtemp) != 0 {
			name := "box " + box_current + " bay " + bay_current + " type " + disktype_current
			MetricValue, _ := strconv.ParseFloat(strings.TrimSpace(mtemp[1]), 64)
			disk_maximum_temperature.With(prometheus.Labels{"physicaldrive": name}).Set(MetricValue)
		}
		if len(power) != 0 {
			name := "box " + box_current + " bay " + bay_current + " type " + disktype_current
			MetricValue, _ := strconv.ParseFloat(strings.TrimSpace(power[1]), 64)
			disk_power_on_hours.With(prometheus.Labels{"physicaldrive": name}).Set(MetricValue)
		}
	}
}

func recordMetrics(interval time.Duration) {
	go func() {
		for {
		    log.Println("Reading new metrics..")
			parse()
			time.Sleep(interval)
		}
	}()
}

func runcmd(cmd string, shell bool) []byte {
    log.Printf("Executing command : %v", cmd)

	if shell {
		out, err := exec.Command(cmd).Output()
		if err != nil {
		    log.Println("Error while executing the command: ", err)
			panic("some error found")
		}
		return out
	}
	out, err := exec.Command(cmd).Output()
	if err != nil {
		log.Fatal(err)
	}
	return out
}

func main() {
	Port := flag.Int("Port", 9109, "Port Number to listen")
	ProbingRate := flag.String("ProbingRate", "1m", "The rate in which the ssacli tool is probed for new values")
    flag.Parse()

    interval, _ := time.ParseDuration(*ProbingRate)

    log.Printf("Starting exporter on port '%d' with probing interval '%s'", *Port, *ProbingRate)

	recordMetrics(interval)
	var port = ":" + strconv.Itoa(*Port)
	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(port, nil))
}
