package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"time"

	"github.com/mongodb/mongo-go-driver/mongo"
	"gitlab.com/codelittinc/golang-interview-project-jaime/grpc-crud/controller"
	"gitlab.com/codelittinc/golang-interview-project-jaime/grpc-crud/hirepb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type Server struct{}

func NewServer() *Server {
	return &Server{}
}

// Run handles the main routine for the server
func (s *Server) Run() {

	getMongoConnection()

	li, err := net.Listen("tcp", "localhost:"+grpcPort)
	if err != nil {
		log.Fatalln("Error on port listen. ", err)
	}

	// creating grpc server and registering
	opts := []grpc.ServerOption{}
	srv := grpc.NewServer(opts...)
	hirepb.RegisterHireDataServiceServer(srv, &controller.HireDataController{})

	// Register reflection service on gRPC server.
	reflection.Register(srv)

	go func() {
		log.Println("Hire Server started")
		if err := srv.Serve(li); err != nil {
			log.Fatalln("grpc Server error. ", err)
		}
	}()

	// wait for interrupt signal
	ch := make(chan os.Signal)
	signal.Notify(ch, os.Interrupt)

	<-ch
	shutdownServer(li, srv)

}

// shutdown server gracefully on interruption
func shutdownServer(li net.Listener, s *grpc.Server) {
	log.Println("Stopping Server")
	s.GracefulStop()
	li.Close()
	controller.SrvCancelFunc()
	controller.SrvClient.Disconnect(controller.SrcCtx)
}

// getMongoConnection is used to establish the connection to the database
func getMongoConnection() {
	log.Println("Starting database")

	client, err := mongo.NewClient("mongodb://localhost:" + mongoPort)
	if err != nil {
		log.Fatalln(err)
	}

	controller.SrcCtx, controller.SrvCancelFunc =
		context.WithTimeout(context.Background(), 4000*time.Second)

	err = client.Connect(controller.SrcCtx)
	if err != nil {
		log.Fatalln(err)
	}

	// opens MongDB database and its collection (creates them if they don't exitst)
	controller.SrvClient = client
	controller.SrvCol = client.Database("hireDB").Collection("hires")
}
