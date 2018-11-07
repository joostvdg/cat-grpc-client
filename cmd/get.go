package cmd

import (
	"context"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"time"

	"github.com/joostvdg/cat/pkg/api/v1"
)

func init() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetReportCaller(true)
	log.SetLevel(log.InfoLevel)
}

func GetApplication(address string, uuid string) {
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	c := v1.NewApplicationServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Read
	req2 := v1.ReadRequest{
		Api:  apiVersion,
		Uuid: uuid,
	}
	res2, err := c.Read(ctx, &req2)
	if err != nil {
		log.Fatalf("Read failed: %v", err)
	}
	log.Printf("Read result: id=%v, <%+v>\n", res2.Application.Uuid, res2)
}
