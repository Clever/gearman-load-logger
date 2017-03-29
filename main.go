package main

import (
	"fmt"
	"log"
	"net"
	"time"

	"github.com/Clever/discovery-go"
	"github.com/Clever/gearadmin"
	"gopkg.in/Clever/kayvee-go.v6/logger"
)

const (
	pollInterval = time.Minute
)

var (
	lg = logger.New("gearman-load-logger")
)

func logMetrics(g gearadmin.GearmanAdmin) {
	statuses, err := g.Status()
	if err != nil {
		log.Fatalf("error retrieving gearman status: %s", err)
	}
	for _, status := range statuses {
		lg.GaugeIntD("total_workers", status.Total, logger.M{
			"function":          status.Function,
			"running_workers":   status.Running,
			"available_workers": status.AvailableWorkers,
		})
	}
}

// This script polls gearman and outputs metrics to syslog in
// a standard format for later parsing
func main() {
	if err := logger.SetGlobalRouting("./kvconfig.yml"); err != nil {
		log.Fatalf("failed to find kayvee config: %s", err)
	}

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
