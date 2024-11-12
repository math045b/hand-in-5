package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	pb "hand-in-5/proto"
	"os"
	"strconv"
	"time"
)

func main() {
	clientId := os.Args[1]
	stringPort := os.Args[2]
	port, _ := strconv.Atoi(stringPort)

	node := &Client{
		ID:        clientId,
		nextPort:  port,
		timestamp: 0,
	}

	node.StartClient()

	for {
		time.Sleep(1 * time.Second)
		bidReg := &pb.BidRequest{
			Amount: 5,
		}
		bidResponse, _ := node.client.PlaceBid(context.Background(), bidReg)
		println(bidResponse)
	}
}

type Client struct {
	client    pb.ServiceClient
	ID        string
	nextPort  int
	timestamp int64
}

func (n *Client) StartClient() {
	address := "localhost:" + strconv.Itoa(n.nextPort)
	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Errorf("failed to connect to node %s: %v", n.ID, err)
	}

	n.client = pb.NewServiceClient(conn)
	n.nextPort++
}
