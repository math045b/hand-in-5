package main

import (
	"bufio"
	"context"
	pb "hand-in-5/proto"
	"log"
	"os"
	"strconv"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	clientId := os.Args[1]

	node := &Client{
		ID:        clientId,
		timestamp: 0,
	}

	node.StartClient()

	reader := bufio.NewReader(os.Stdin)
	println("You can input the following: ")
	println("Bid {amount}")
	println("Result")

	for {

		text, _ := reader.ReadString('\n')
		splittext := strings.Split(text, " ")

		if len(splittext) < 2 {
			if strings.TrimSpace(splittext[0]) == "Result" {
				log.Println("Getting auction result")
				err := Request(node)
				if err != nil {
					log.Println("Error getting auction result")
					Request(node)
				}

			}
		} else {
			if splittext[0] == "Bid" {
				num, _ := strconv.Atoi(strings.TrimSpace(splittext[1]))
				Bid(node, num)
			}
		}
	}
}

func Bid(node *Client, num int) {
	bidReg := &pb.BidRequest{
		Amount: int64(num),
		Port:   node.ID,
	}
	bidResponse, _ := node.client.PlaceBid(context.Background(), bidReg)
	log.Println(bidResponse.Response)

}

func Request(node *Client) error {
	req := &pb.ResultRequest{}
	res, err := node.client.AuctionResult(context.Background(), req)
	if err != nil {
		return err
	}
	log.Println(res.Response)
	return err
}

type Client struct {
	client    pb.ServiceClient
	ID        string
	timestamp int64
}

func (n *Client) StartClient() {
	address := "localhost: 5050"
	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect to node %s: %v", n.ID, err)
	}
	n.client = pb.NewServiceClient(conn)
}
