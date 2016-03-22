package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"time"

	"github.com/Clever/discovery-go"
	"github.com/Clever/gearadmin"
	"gopkg.in/Clever/kayvee-go.v3"
)

func logMetrics(g gearadmin.GearmanAdmin) {
	statuses, err := g.Status()
	if err != nil {
		log.Fatalf("error retrieving gearman status: %s", err)
	}
	for _, status := range statuses {
		fmt.Println(kayvee.FormatLog("gearlogger", "info", "status", map[string]interface{}{
			"type":     "gauge",
			"function": status.Function,
			"running":  status.Running,
			"total":    status.Total,
			"workers":  status.AvailableWorkers,
		}))
	}
}

// This script polls gearman and outputs metrics to syslog in
// a standard format for later parsing
func main() {
	host := flag.String("host", "", "target gearman host")
	port := flag.Int("port", 4730, "target gearman port")
	interval := flag.Duration("interval", time.Minute, "interval to log at")
	flag.Parse()

	h := *host
	if discoveredHost, err := discovery.Host("gearmand", "tcp"); err == nil {
		h = discoveredHost
	}

	if discoveredPort, err := discovery.Port("gearmand", "tcp"); err == nil {
		if portInt, err := strconv.Atoi(discoveredPort); err == nil {
			port = &portInt
		}
	}

	if h == "" {
		log.Println("must specify host")
		flag.PrintDefaults()
		os.Exit(1)
	}

	// connect to gearman
	c, err := net.Dial("tcp", fmt.Sprintf("%s:%d", h, *port))
	if err != nil {
		panic(err)
	}
	defer c.Close()
	admin := gearadmin.NewGearmanAdmin(c)

	logMetrics(admin)
	for _ = range time.Tick(*interval) {
		logMetrics(admin)
	}
}
