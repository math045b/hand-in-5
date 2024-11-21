package main

import (
	"bufio"
	"context"
	"fmt"
	pb "hand-in-5/proto"
	"os"
	"strconv"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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

	reader := bufio.NewReader(os.Stdin)
	println("You can input the following: ")
	println("Bid {amount}")
	println("Result")

	lastbet := 0

	for {

		text, _ := reader.ReadString('\n')
		splittext := strings.Split(text, " ")

		if len(splittext) < 2 {
			if strings.TrimSpace(splittext[0]) == "Result" {
				println("Getting acution result")
				req := &pb.ResultRequest{}
				res, err := node.client.AuctionResult(context.Background(), req)
				if err != nil {
					println("Error getting auction result")
					node.IncrementPort()
				}
				println(res.Response)
			}
		} else {
			if splittext[0] == "Bid" {
				num, _ := strconv.Atoi(strings.TrimSpace(splittext[1]))
				if num < lastbet {
					println("bet must be larger than last bet!")
					continue
				}
				lastbet = num
				bidReg := &pb.BidRequest{
					Amount: int64(num),
					Port:   clientId,
				}
				bidResponse, err := node.client.PlaceBid(context.Background(), bidReg)
				if err != nil {
					println("port was unresponsive incrementing port num")
					node.IncrementPort()

					continue
				}
				println(bidResponse.Response)
			}
		}
	}
}

func (n *Client) IncrementPort() {
	n.StartClient()
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
