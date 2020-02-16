package main

import (
	"context"
	"fmt"
	"job/pkg/jobpb"
	"log"
	"net"
	"os"
	"sync"

	"google.golang.org/grpc"
)

const (
	portKey = "PORT"
	echoKey = "ECHO"
)

var (
	address = fmt.Sprintf(":%s", os.Getenv(portKey))
	echo    = os.Getenv(echoKey)
)

func init() {
	for _, k := range []string{portKey, echoKey} {
		if os.Getenv(k) == "" {
			log.Fatalf("%s empty and expected from env variables", k)
		}
	}
}

func main() {
	s := &server{}
	s.start()
}

type server struct {
	wg *sync.WaitGroup
}

func (s *server) start() {
	if s.wg == nil {
		s.wg = &sync.WaitGroup{}
	}
	log.Printf("Listening on %s", address)
	s.wg.Add(1)
	go func(s *server) {
		lis, err := net.Listen("tcp", address)
		if err != nil {
			log.Fatal(err)
		}
		gs := grpc.NewServer()
		jobpb.RegisterWorkerServer(gs, s)
		err = gs.Serve(lis)
		if err != nil {
			log.Fatal(err)
		}

	}(s)
	s.wg.Wait()
}

func (s *server) Echo(ctx context.Context, in *jobpb.Job) (result *jobpb.Job, err error) {
	result = &jobpb.Job{
		Id: echo,
	}
	return
}

func (s *server) Stop(ctx context.Context, in *jobpb.Job) (result *jobpb.Job, err error) {
	result = &jobpb.Job{
		Id: echo,
	}
	s.wg.Done()
	return
}
