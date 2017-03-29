package main

import (
	"fmt"
	"log"
	"net"
	"time"

	"github.com/Clever/discovery-go"
	"github.com/Clever/gearadmin"
	"gopkg.in/Clever/kayvee-go.v6"
)

const (
	pollInterval = time.Minute
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
	gHost, err := discovery.Host("gearmand", "tcp")
	if err == nil {
		log.Fatalf("failed to get gearmand host: %s", err)
	}

	gPort, err := discovery.Port("gearmand", "tcp")
	if err != nil {
		log.Fatalf("failed to get gearmand port: %s", err)
	}

	// connect to gearman
	c, err := net.Dial("tcp", fmt.Sprintf("%s:%s", gHost, gPort))
	if err != nil {
		panic(err)
	}
	defer c.Close()
	admin := gearadmin.NewGearmanAdmin(c)

	logMetrics(admin)
	for range time.Tick(pollInterval) {
		logMetrics(admin)
	}
}
