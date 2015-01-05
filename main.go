package main

import (
	"flag"
	"github.com/Clever/gearadmin"
	"gopkg.in/Clever/kayvee-go.v1"
	"log"
	"net"
	"os"
	"strconv"
)

// This script polls gearman and outputs metrics to syslog in
// a standard format for later parsing
func main() {
	hostPtr := flag.String("host", "", "target gearman host")
	portPtr := flag.Int("port", 4730, "target gearman port")
	flag.Parse()

	if *hostPtr == "" {
		log.Println("must specify host")
		flag.PrintDefaults()
		os.Exit(1)
	}

	// connect to gearman
	c, err := net.Dial("tcp", *hostPtr+":"+strconv.Itoa(*portPtr))
	if err != nil {
		panic(err)
	}
	defer c.Close()
	admin := gearadmin.NewGearmanAdmin(c)

	var metrics map[string]interface{}

	// get status and output to syslog
	status, _ := admin.Status()
	for _, funcStatus := range status {
		metrics = make(map[string]interface{})

		metrics["function"] = funcStatus.Function
		metrics["total_jobs"] = float64(funcStatus.Total)
		metrics["running_jobs"] = float64(funcStatus.Running)
		metrics["available_workers"] = float64(funcStatus.AvailableWorkers)

		metrics["worker_load"] = 0.0
		if funcStatus.AvailableWorkers > 0 {
			metrics["worker_load"] = metrics["total_jobs"].(float64) / metrics["available_workers"].(float64)
		} else if funcStatus.Total > 0 {
			metrics["worker_load"] = 99.0 // some large number to stand in for infinity
		}

		log.Printf(kayvee.FormatLog("GEARMAN", "INFO", "QUEUE_INFO", metrics))
	}
}
