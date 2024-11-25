package main

import (
	"bufio"
	"context"
	"fmt"
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
				log.Println("Getting acution result")
				err := Request(node)
				if err != nil {
					log.Println("Error getting auction result")
					node.IncrementPort()
					Request(node)
				}

			}
		} else {
			if splittext[0] == "Bid" {
				num, _ := strconv.Atoi(strings.TrimSpace(splittext[1]))
				if num <= lastbet {
					log.Println("bet must be larger than last bet!")
					continue
				}
				lastbet = num
				err := Bid(node, num)
				if err != nil {
					log.Println("port was unresponsive incrementing port num")
					node.IncrementPort()
					Bid(node, num)
					continue
				}

			}
		}
	}
}

func (n *Client) IncrementPort() {
	n.StartClient()
}

func Bid(node *Client, num int) error {
	bidReg := &pb.BidRequest{
		Amount: int64(num),
		Port:   node.ID,
	}
	bidResponse, err := node.client.PlaceBid(context.Background(), bidReg)
	if err != nil {
		return err
	}
	log.Println(bidResponse.Response)
	return err

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
