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
	"time"

	p2p "github.com/frederikgantriis/P2P-DISYS/src"
	"google.golang.org/grpc"
)

func main() {
	var ownPort int32
	var f *os.File
	var err error

	// Check if ressources folder exists
	_, ferr := os.Stat("ressources")
	if os.IsNotExist(ferr) {
		os.Mkdir("ressources", os.ModePerm)
	}

	// Open or create ports.txt file
	f, err = os.OpenFile("ressources/ports.txt", os.O_RDWR, 0x0666)
	if err != nil {
		err = nil
		f, err = os.Create("ressources/ports.txt")
		if err != nil {
			log.Println(err)
			log.Fatalln("Could not read nor create port file")
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
				log.Fatalf("Could not convert last port to int: %s\n", portString)
			}
		}
	}
	if len(ports) == 0 {
		ownPort = int32(3000)
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
		id:      ownPort,
		lamport: 0,
		state:   State(RELEASED),
		clients: make(map[int32]p2p.ReqAccessToCSClient),
		ctx:     ctx,
	}

	// Create listener tcp on port ownPort
	list, err := net.Listen("tcp", fmt.Sprintf(":%v", ownPort))
	if err != nil {
		log.Fatalf("Failed to listen on port: %v\n", err)
	}
	grpcServer := grpc.NewServer()
	p2p.RegisterReqAccessToCSServer(grpcServer, p)

	go func() {
		if err := grpcServer.Serve(list); err != nil {
			log.Fatalf("failed to server %v\n", err)
		}
	}()

	connections := []*grpc.ClientConn{}
	go func() {
		for _, port := range ports {
			var conn *grpc.ClientConn
			fmt.Printf("Trying to dial: %v\n", port)
			conn, err := grpc.Dial(fmt.Sprintf(":%v", port), grpc.WithInsecure(), grpc.WithBlock())
			if err != nil {
				log.Fatalf("Could not connect: %s\n", err)
			}
			connections = append(connections, conn)
			c := p2p.NewReqAccessToCSClient(conn)
			p.clients[port] = c
			log.Printf("Connected to port: %d\n", port)
		}
	}()

	go func() {
		port := ownPort + 1
		for {
			var conn *grpc.ClientConn
			fmt.Printf("Trying to dial: %v\n", port)
			conn, err := grpc.Dial(fmt.Sprintf(":%v", port), grpc.WithInsecure(), grpc.WithBlock())
			if err != nil {
				log.Fatalf("Could not connect: %s\n", err)
			}
			connections = append(connections, conn)
			c := p2p.NewReqAccessToCSClient(conn)
			p.clients[port] = c
			log.Printf("Connected to port: %d\n", port)
			port++
		}
	}()

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		if scanner.Text() == "-q" {
			_ports := getPortsExcept(ownPort)
			deletePortFile()
			if len(_ports) > 0 {
				CreateAndWritePortsToFile(_ports)
			}
			break
		}
		p.sendReqAccessToCSToAll()
	}

	for _, conn := range connections {
		conn.Close()
	}
}

type peer struct {
	p2p.UnimplementedReqAccessToCSServer
	id      int32
	state   State
	lamport int32
	clients map[int32]p2p.ReqAccessToCSClient
	ctx     context.Context
}

type State int32

const (
	RELEASED State = 0
	WANTED   State = 1
	HELD     State = 2
)

func (p *peer) ReqAccessToCS(ctx context.Context, req *p2p.Request) (*p2p.Reply, error) {
	for p.state == State(HELD) || (p.state == State(WANTED) && (p.lamport < req.Lamport || (p.lamport == req.Lamport && p.id < req.Id))) {
		continue
	}
	rep := &p2p.Reply{Lamport: p.lamport}
	setLamportTimestamp(p, int(rep.Lamport))
	return rep, nil
}

func (p *peer) sendReqAccessToCSToAll() {
	p.lamport++
	p.state = State(WANTED)
	fmt.Printf("peer; %v, wants access to critical section \n", p.id)

	request := &p2p.Request{Id: p.id, Lamport: p.lamport}

	for id, client := range p.clients {
		reply, err := client.ReqAccessToCS(p.ctx, request)
		if err != nil {
			delete(p.clients, id)
			continue
		}
		fmt.Printf("Got reply from id %v: %v\n", id, reply.Lamport)
	}

	p.state = State(HELD)
	fmt.Printf("peer; %v, has gained access to the critical section\n", p.id)
	criticalSection(p)
	p.state = State(RELEASED)
	fmt.Printf("peer; %v, has released access to the critical section\n", p.id)
}

func setLamportTimestamp(p *peer, incoming int) {
	p.lamport = int32(math.Max(float64(p.lamport), float64(incoming)) + 1)
}

func criticalSection(p *peer) {
	fmt.Printf("A Critical Hello World from %v\n", p.id)
	time.Sleep(3 * time.Second)
}

func deletePortFile() {
	os.Remove("ressources/ports.txt")
}

func CreateAndWritePortsToFile(ports []string) {
	f, err := os.Create("ressources/ports.txt")

	if err != nil {
		log.Fatalf("Couldn't create file: ", err)
	}

	for _, port := range ports {
		f.WriteString(port + "\n")
	}
}

func getPortsExcept(port int32) []string {
	portString := strconv.Itoa(int(port))

	f, err := os.OpenFile("ressources/ports.txt", os.O_RDWR, 0x0666)
	if err != nil {
		log.Fatalln("Couldn't read file: ", err)
	}

	// Read text file
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
		if string(b[:readTotal-1]) != portString {
			portStrings = append(portStrings, string(b[:readTotal-1]))
		}
	}

	f.Close()

	return portStrings
}
