package main

import (
	"fmt"
	"hand-in-5/proto"
	"net"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ServerNode struct {
	proto.UnimplementedServiceServer
	port     string
	nextNode proto.ServiceClient
}

func main() {
	// starting node
	clientPort := os.Args[1]

	node := &ServerNode{
		port: clientPort,
	}

	node.StartServer()

	if len(os.Args) > 3 {

	}

}

func (n *ServerNode) StartClient(port string) {
	conn, err := grpc.NewClient("localhost:"+port, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Println("could not make client")
	}
	n.nextNode = proto.NewServiceClient(conn)
}

func (n *ServerNode) StartServer() {
	server := grpc.NewServer()
	listener, err := net.Listen("tcp", ":"+n.port)

	if err != nil {
		fmt.Println("Listener did not start")
	}

	proto.RegisterServiceServer(server, n)

	err = server.Serve(listener)
	if err != nil {
		fmt.Println("Server could not be served")
	}
}
