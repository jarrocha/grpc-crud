package main

import (
	"log"

	"gitlab.com/codelittinc/golang-interview-project-jaime/grpc-crud/controller"
	"gitlab.com/codelittinc/golang-interview-project-jaime/grpc-crud/hirepb"
	"google.golang.org/grpc"
)

type Client struct{}

func NewClient() *Client {
	return &Client{}
}

var server_conn *grpc.ClientConn

/*
 Run handles the main routine for the client
 The client provides a text based menu to interact with the database server.
 Features:
 - Create: will ask for the name, role, duration, type, tags for the employee
 - List: will request the database to show all current hires
 - Find: will ask the user to find by name, type, role, duration, or tags
 - Delete: will ask the user which hire to delete by name
 - Update: will ask the user which hire to update by name
*/
func (c *Client) Run() {
	StartServerConnection()

	controller.DisplayMainMenu()

	StopServerConnection()

}

func StartServerConnection() {
	var err error
	server_conn, err = grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalln("gRPC Dial error. ", err)
	}
	controller.ServiceClient = hirepb.NewHireDataServiceClient(server_conn)
}

func StopServerConnection() {
	server_conn.Close()
}
