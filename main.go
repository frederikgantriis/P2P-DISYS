package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"math"
	"net"
	"os"
	"strconv"

	p2p "github.com/frederikgantriis/P2P-DISYS/src"
	"google.golang.org/grpc"
)

func main() {
	arg1, _ := strconv.ParseInt(os.Args[1], 10, 32)
	ownPort := int32(arg1) + 3000

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	p := &peer{
		id:      ownPort,
		Lamport: 0,
		clients: make(map[int32]p2p.ReqAccessToCSClient),
		ctx:     ctx,
	}

	// Create listener tcp on port ownPort
	list, err := net.Listen("tcp", fmt.Sprintf(":%v", ownPort))
	if err != nil {
		log.Fatalf("Failed to listen on port: %v", err)
	}
	grpcServer := grpc.NewServer()
	p2p.RegisterReqAccessToCSServer(grpcServer, p)

	go func() {
		if err := grpcServer.Serve(list); err != nil {
			log.Fatalf("failed to server %v", err)
		}
	}()

	for i := 0; i < 3; i++ {
		port := int32(3000) + int32(i)

		if port == ownPort {
			continue
		}

		var conn *grpc.ClientConn
		fmt.Printf("Trying to dial: %v\n", port)
		conn, err := grpc.Dial(fmt.Sprintf(":%v", port), grpc.WithInsecure(), grpc.WithBlock())
		if err != nil {
			log.Fatalf("Could not connect: %s", err)
		}
		defer conn.Close()
		c := p2p.NewReqAccessToCSClient(conn)
		p.clients[port] = c
	}

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		p.sendReqAccessToCSToAll()
	}
}

type peer struct {
	p2p.UnimplementedReqAccessToCSServer
	id      int32
	Lamport int32
	clients map[int32]p2p.ReqAccessToCSClient
	ctx     context.Context
}

func (p *peer) ReqAccessToCS(ctx context.Context, req *p2p.Request) (*p2p.Reply, error) {
	rep := &p2p.Reply{Lamport: p.Lamport}
	setLamportTimestamp(p, int(rep.Lamport))
	return rep, nil
}

func (p *peer) sendReqAccessToCSToAll() {
	p.Lamport++
	request := &p2p.Request{Id: p.id, Lamport: p.Lamport}

	for id, client := range p.clients {
		reply, err := client.ReqAccessToCS(p.ctx, request)
		if err != nil {
			fmt.Println("something went wrong: ", err)
		}
		fmt.Printf("Got reply from id %v: %v\n", id, reply.Lamport)
	}
}

func setLamportTimestamp(p *peer, incoming int) {
	p.Lamport = int32(math.Max(float64(p.Lamport), float64(incoming)) + 1)
}

func criticalSection() {
	fmt.Println("A Critical Hello World")
}
