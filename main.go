package main

import (
	"fmt"
	"log"
	"math"
	"net"
	"os"
	"time"

	"github.com/Clever/discovery-go"
	"github.com/Clever/gearadmin"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"gopkg.in/Clever/kayvee-go.v6/logger"
)

const (
	pollInterval = time.Minute
)

var (
	lg        = logger.New("gearman-load-logger")
	deployEnv = "unknown"
)

func logMetrics(g gearadmin.GearmanAdmin, cw *cloudwatch.CloudWatch) {
	statuses, err := g.Status()
	if err != nil {
		log.Fatalf("error retrieving gearman status: %s", err)
	}
	for _, status := range statuses {
		// Log data, routed via kvconfig
		lg.GaugeIntD("total_workers", status.Total, logger.M{
			"function":          status.Function,
			"running_workers":   status.Running,
			"available_workers": status.AvailableWorkers,
		})

		// Write data to CloudWatch
		jobToWorkerRatio := float64(status.Total) / float64(status.AvailableWorkers)
		if math.IsInf(jobToWorkerRatio, 1) || math.IsNaN(jobToWorkerRatio) {
			// handle divide-by-zero cases (x/0 or 0/0)
			continue
		}

		_, err := cw.PutMetricData(&cloudwatch.PutMetricDataInput{
			MetricData: []*cloudwatch.MetricDatum{
				&cloudwatch.MetricDatum{
					MetricName: aws.String("JobToWorkerRatio"),
					Unit:       aws.String(cloudwatch.StandardUnitNone),
					Value:      aws.Float64(jobToWorkerRatio),
					Dimensions: []*cloudwatch.Dimension{
						&cloudwatch.Dimension{
							Name:  aws.String("Environment"),
							Value: aws.String(deployEnv),
						},
						&cloudwatch.Dimension{
							Name:  aws.String("WorkerName"),
							Value: aws.String(status.Function),
						},
					},
					StorageResolution: aws.Int64(1), // 1 minute, for "coarse" resolution
				},
			},
			Namespace: aws.String("Gearman"),
		})
		if err != nil {
			lg.ErrorD("cloudwatch-write-error", logger.M{
				"worker-name": status.Function,
				"error":       err.Error(),
			})
		}
	}
}

// This script polls gearman and outputs metrics to syslog in
// a standard format for later parsing
func main() {
	if err := logger.SetGlobalRouting("./kvconfig.yml"); err != nil {
		log.Fatalf("failed to find kayvee config: %s", err)
	}

	gHost, err := discovery.Host("gearmand", "tcp")
	if err != nil {
		log.Fatalf("failed to get gearmand host: %s", err)
	}

	gPort, err := discovery.Port("gearmand", "tcp")
	if err != nil {
		log.Fatalf("failed to get gearmand port: %s", err)
	}

	// injected by Catapult
	if os.Getenv("_DEPLOY_ENV") != "" {
		deployEnv = os.Getenv("_DEPLOY_ENV")
	}

	// connect to AWS
	sess := session.Must(session.NewSession(aws.NewConfig().WithRegion("us-west-1").WithMaxRetries(4)))
	// create new cloudwatch client.
	cw := cloudwatch.New(sess)

	// connect to gearman
	c, err := net.Dial("tcp", fmt.Sprintf("%s:%s", gHost, gPort))
	if err != nil {
		panic(err)
	}
	defer c.Close()
	admin := gearadmin.NewGearmanAdmin(c)

	logMetrics(admin, cw)
	for range time.Tick(pollInterval) {
		logMetrics(admin, cw)
	}
}
