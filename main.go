package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"math"
	"net"
	"os"
	"strconv"

	ping "github.com/frederikgantriis/P2P-DISYS/src"
	"google.golang.org/grpc"
)

var lamport int32

func main() {
	lamport = 0
	var ownPort int32

	// Creating/reading text file
	var f *os.File
	var err error

	// Open or create ports.txt file
	f, err = os.OpenFile("ressources/ports.txt", os.O_RDWR, 0x0666)
	if err != nil {
		err = nil
		f, err = os.Create("ressources/ports.txt")
		defer deletePortFile()
		if err != nil {
			log.Println(err)
			log.Fatalf("Could not read nor create port file")
		}
	}

	// Read all ports from text file and add to slice
	const maxSz = 5
	b := make([]byte, maxSz)
	portStrings := []string{}
	for {
		// read content to buffer
		readTotal, err := f.Read(b)
		if err != nil {
			if err != io.EOF {
				fmt.Println(err)
			}
			break
		}
		fmt.Println("I read this port: " + string(b[:readTotal-1])) // print content from buffer
		portStrings = append(portStrings, string(b[:readTotal-1]))
	}

	// Print the ports read
	for _, str := range portStrings {
		fmt.Println("Saved port: " + str)
	}

	// Convert port strings to int32's
	ports := []int32{}
	for _, portString := range portStrings {
		if len(portStrings) > 0 {
			port, err := strconv.Atoi(portString)
			ports = append(ports, int32(port))
			if err != nil {
				log.Fatalf("Could not convert last port to int: %s", portString)
			}
		}
	}
	if len(ports) == 0 {
		ownPort = int32(5000)
	} else {
		ownPort = int32(ports[len(ports)-1] + 1)
	}
	log.Printf("My port is: %d\n", ownPort)
	f.WriteString(fmt.Sprint(ownPort) + "\n")
	f.Close()
	log.Println("Wrote to and closed file")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	p := &peer{
		id:            ownPort,
		lamport:       lamport,
		amountOfPings: make(map[int32]int32),
		clients:       make(map[int32]ping.PingClient),
		ctx:           ctx,
	}

	// Create listener tcp on port ownPort
	list, err := net.Listen("tcp", fmt.Sprintf(":%v", ownPort))
	if err != nil {
		log.Fatalf("Failed to listen on port: %v", err)
	}
	grpcServer := grpc.NewServer()
	ping.RegisterPingServer(grpcServer, p)

	go func() {
		if err := grpcServer.Serve(list); err != nil {
			log.Fatalf("failed to server %v", err)
		}
	}()

	go func() {
		for _, port := range ports {
			var conn *grpc.ClientConn
			fmt.Printf("Trying to dial: %v\n", port)
			conn, err := grpc.Dial(fmt.Sprintf(":%v", port), grpc.WithInsecure(), grpc.WithBlock())
			if err != nil {
				log.Fatalf("Could not connect: %s", err)
			}
			defer conn.Close()
			c := ping.NewPingClient(conn)
			p.clients[port] = c
			log.Printf("Connected to port: %d\n", port)
		}
	}()

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		if scanner.Text() == "-q" {
			break
		}
		p.sendPingToAll()
	}
}

type peer struct {
	ping.UnimplementedPingServer
	id            int32
	lamport       int32
	amountOfPings map[int32]int32
	clients       map[int32]ping.PingClient
	ctx           context.Context
}

func (p *peer) Ping(ctx context.Context, req *ping.Request) (*ping.Reply, error) {
	id := req.Id
	p.amountOfPings[id] += 1

	rep := &ping.Reply{Amount: p.amountOfPings[id]}
	return rep, nil
}

func (p *peer) sendPingToAll() {
	lamport++
	request := &ping.Request{Id: p.id, Lamport: p.lamport}
	for id, client := range p.clients {
		reply, err := client.Ping(p.ctx, request)
		if err != nil {
			fmt.Println("something went wrong")
		}
		fmt.Printf("Got reply from id %v: %v\n", id, reply.Amount)
	}
}

func setLamportTimestamp(incoming int) {
	lamport = int32(math.Max(float64(lamport), float64(incoming)) + 1)
}

func criticalSection() {
	fmt.Println("A Critical Hello World")
}

func deletePortFile() {
	os.Remove("ressources/ports.txt")
}
