package main

import (
	"context"
	"fmt"
	"hand-in-5/proto"
	"log"
	"net"
	"os"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ServerNode struct {
	proto.UnimplementedServiceServer
	port     string
	nextNode proto.ServiceClient
	action   map[string]int64
}

var isRunning bool
var startTime time.Time
var node *ServerNode
var winner string

func main() {
	// starting node

	fmt.Printf("Starting server node with port %s and client port %s\n", os.Args[1], os.Args[2])
	clientPort := os.Args[1]
	isRunning = false
	winner = ""
	node = &ServerNode{
		port: clientPort,
	}
	if len(os.Args) > 2 {
		fmt.Printf("Server port: %s, clientport: %s \n", os.Args[1], os.Args[2])
		node.StartClient(os.Args[2])
	}
	go RunAuction()
	node.StartServer()
}

func RunAuction() {
	for {
		if isRunning {
			if time.Since(startTime) > (30 * time.Second) {
				fmt.Println("start: ", startTime.String())
				log.Println("The auction was ended")
				isRunning = false
				winner = AssertWinner()
				node.nextNode.UpdateNodes(
					context.Background(),
					&proto.NodeUpdate{
						Auctionstate:   node.action,
						AuctionOngoing: isRunning,
						Time:           startTime.Unix(),
					})
				break
			}
		}
	}
}

func AssertWinner() (max string) {
	maxVal := int64(0)
	maxKey := ""
	for key, value := range node.action {
		if value > maxVal {
			maxVal = value
			maxKey = key
		}
	}
	return maxKey
}

func (n *ServerNode) UpdateNodes(ctx context.Context, data *proto.NodeUpdate) (*proto.Empty, error) {
	fmt.Printf("Port: %s Received update\n", n.port)
	// updating node auction map
	n.action = data.Auctionstate
	isRunning = data.AuctionOngoing
	startTime = time.Unix(data.Time, 0)

	winner = AssertWinner()

	state := &proto.NodeUpdate{
		Auctionstate:   n.action,
		AuctionOngoing: isRunning,
		Time:           data.Time,
	}

	if n.nextNode != nil {
		n.nextNode.UpdateNodes(context.Background(), state)
	}

	return &proto.Empty{}, nil
}

func (n *ServerNode) startTimer(t time.Time) {
	log.Println("The auction was started")
	startTime = t
	update := &proto.NodeUpdate{
		Auctionstate:   n.action,
		AuctionOngoing: isRunning,
		Time:           t.Unix(),
	}
	isRunning = true
	n.nextNode.UpdateNodes(context.Background(), update)
}

func (n *ServerNode) PlaceBid(ctx context.Context, request *proto.BidRequest) (*proto.BidResponse, error) {
	var response *proto.BidResponse
	if n.action == nil {
		n.startTimer(time.Now())
		n.action = make(map[string]int64)
		n.action[request.Port] = request.Amount
		response = &proto.BidResponse{Response: fmt.Sprintf("%s: You have joined the auction with your bid of %d", request.Port, request.Amount)}
	} else if n.action != nil && !isRunning {
		response = &proto.BidResponse{Response: "The auction has ended"}
	} else {
		n.action[request.Port] = request.Amount
		response = &proto.BidResponse{Response: fmt.Sprintf("Bid Received! New Bid: %d", n.action[request.Port])}
		update := &proto.NodeUpdate{Auctionstate: n.action, AuctionOngoing: isRunning, Time: startTime.Unix()}
		if n.nextNode != nil {
			_, err := n.nextNode.UpdateNodes(context.Background(), update)
			if err != nil {
				fmt.Printf("Failed to forward update to next node: %v\n", err)
			}
		} else {
			fmt.Printf("Port: %s - No next node to forward update to\n", n.port)
		}
	}

	return response, nil
}

func (n *ServerNode) AuctionResult(context context.Context, request *proto.ResultRequest) (*proto.ResultResponse, error) {
	log.Println("Auction result was requested")
	winner = AssertWinner()
	var response *proto.ResultResponse
	if !isRunning && winner != "" {
		response = &proto.ResultResponse{Response: fmt.Sprintf("Auction is over the winner was %s", winner)}
	} else {
		response = &proto.ResultResponse{Response: fmt.Sprintf("Auction is ongoing the current max bid is: %d", node.action[winner])}
	}

	return response, nil
}

func (n *ServerNode) StartClient(port string) error {
	println("client is starting brrr")
	address := "localhost:" + port
	fmt.Printf("%s is trying to create a client on", n.port)
	fmt.Println(address)
	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return fmt.Errorf("failed to connect to next node on port %s: %w", port, err)
	}
	n.nextNode = proto.NewServiceClient(conn)
	fmt.Printf("Port: %s - Connected to next node on port %s\n", n.port, port)
	return nil
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
