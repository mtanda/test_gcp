package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/golang/protobuf/ptypes"
	"golang.org/x/net/context"
	"google.golang.org/api/iterator"

	monitoring "cloud.google.com/go/monitoring/apiv3"
	monitoringpb "google.golang.org/genproto/googleapis/monitoring/v3"
)

func main() {
	projectID := flag.String("project", "", "project id")
	fmt.Printf("project = %s\n", *projectID)
	ctx := context.Background()

	client, err := monitoring.NewMetricClient(ctx)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	/*
		req := &monitoringpb.ListMetricDescriptorsRequest{
			Name: "projects/" + projectID,
		}
		it := client.ListMetricDescriptors(ctx, req)
		for {
			resp, err := it.Next()
			if err == iterator.Done {
				break
			}
			if err != nil {
				fmt.Printf("%+v\n", err)
				break
			}
			fmt.Printf("%+v\n", resp)
		}
	*/
	startTime, _ := ptypes.TimestampProto(time.Now().Add(-1 * time.Hour))
	req := &monitoringpb.ListTimeSeriesRequest{
		Name:   "projects/" + *projectID,
		Filter: "metric.type = \"compute.googleapis.com/instance/cpu/utilization\"",
		Interval: &monitoringpb.TimeInterval{
			StartTime: startTime,
			EndTime:   ptypes.TimestampNow(),
		},
	}
	it := client.ListTimeSeries(ctx, req)
	for {
		resp, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			fmt.Printf("%+v\n", err)
			break
		}
		fmt.Printf("%+v\n", resp.Metric)
		fmt.Printf("%+v\n", resp.Resource)
		fmt.Printf("%+v\n", resp.MetricKind)
		fmt.Printf("%+v\n", resp.ValueType)
		for _, p := range resp.Points {
			fmt.Printf("%+v %+v\n", p.Interval, p.Value)
		}
	}
}
