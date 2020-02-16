package main

import (
	"context"
	"fmt"
	"job/pkg/jobpb"
	"log"
	"net/http"
	"os"
	"strconv"

	"google.golang.org/grpc"
)

const (
	portKey       = "PORT"
	nKey          = "N"
	jobChannelKey = "JOB_CHANNEL"
)

var (
	n          int
	jobChannel = os.Getenv(jobChannelKey)
	address    = fmt.Sprintf(":%s", os.Getenv(portKey))
)

func init() {
	var err error
	for _, k := range []string{nKey, jobChannelKey, portKey} {
		if os.Getenv(k) == "" {
			log.Fatalf("%s empty but expected from env variables", k)
		}
	}
	n, err = strconv.Atoi(os.Getenv(nKey))
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	http.HandleFunc("/", alive)
	http.HandleFunc("/job", runJob)
	log.Printf("Listening on %s\n", address)
	err := http.ListenAndServe(address, nil)
	if err != nil {
		log.Fatal(err)
	}
}

func alive(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "I'm alive")
}

func runJob(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL.Path)
	ctx := context.Background()
	conn, err := grpc.Dial(jobChannel, grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	client := jobpb.NewWorkerClient(conn)
	i := -1
	for {
		var err error
		var result *jobpb.Job
		i++
		if i >= n {
			return
		}
		result, err = client.Echo(ctx, &jobpb.Job{})
		if err != nil {
			log.Println(err)
			return
		}
		log.Println(result.Id)
	}
}
