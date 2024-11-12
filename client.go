package main

import (
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	pb "hand-in-5/proto"
)

type Client struct {
	clients   []pb.ServiceClient
	ID        string
	timestamp int64
}

func (n *Client) StartClient(addresses []string) {
	for _, address := range addresses {
		conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			fmt.Errorf("failed to connect to node %s: %v", n.ID, err)
		}

		client := pb.NewServiceClient(conn)
		n.clients = append(n.clients, client)
	}
}
