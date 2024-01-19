package main

import (
	"context"
	"fmt"
	"net"
	"sync"

	pb "github.com/Manish-Saharan/train-ticket-management/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type trainServer struct {
	//this is the struct to be created, pb is assigned above
	mu    sync.Mutex
	users map[string]*pb.Receipt
	pb.UnimplementedTicketServiceServer
}

func newTrainServer() *trainServer {
	return &trainServer{
		users: make(map[string]*pb.Receipt),
	}
}

func (s *trainServer) PurchaseTicket(ctx context.Context, req *pb.Receipt) (*pb.Receipt, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if req == nil || req.User == nil {
		return nil, fmt.Errorf("invalid request: missing user information")
	}

	if req.Section == nil {
		req.Section = &pb.TSection{}
	}

	// Assign section based on available seats
	if len(s.getUsersInSection("A")) < len(s.getUsersInSection("B")) {
		req.Section.Sec = "A"
	} else {
		req.Section.Sec = "B"
	}

	s.users[req.User.Email] = req

	fmt.Printf("Ticket purchased: %+v\n", req)
	return req, nil
}

func (s *trainServer) GetReceipt(ctx context.Context, req *pb.User) (*pb.Receipt, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if receipt, exists := s.users[req.Email]; exists {
		fmt.Printf("Viewing receipt for user %s\n", req.Email)
		return receipt, nil
	}
	return nil, fmt.Errorf("user not found")
}

func (s *trainServer) GetUsersBySection(section *pb.TSection, stream pb.TicketService_GetUsersBySectionServer) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, receipt := range s.getUsersInSection(section.Sec) {
		if err := stream.Send(receipt); err != nil {
			return err
		}
	}

	return nil
}

func (s *trainServer) RemoveUser(ctx context.Context, req *pb.User) (*pb.Receipt, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if receipt, exists := s.users[req.Email]; exists {
		delete(s.users, req.Email)
		fmt.Printf("User %s removed from train\n", req.Email)
		return receipt, nil
	}
	return nil, fmt.Errorf("user not found")
}

func (s *trainServer) ModifySeat(ctx context.Context, req *pb.Receipt) (*pb.Receipt, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if existing, exists := s.users[req.User.Email]; exists {
		existing.Section.Sec = req.Section.Sec
		s.users[req.User.Email] = existing

		fmt.Printf("Seat modified for user %s\n", req.User.Email)
		return req, nil
	}
	return nil, fmt.Errorf("user not found")
}

func (s *trainServer) getUsersInSection(section string) []*pb.Receipt {
	var result []*pb.Receipt
	for _, receipt := range s.users {
		if receipt.Section.Sec == section {
			result = append(result, receipt)
		}
	}
	return result
}

func main() {
	// create a new gRPC server
	server := grpc.NewServer()

	// register the trainTicket service
	pb.RegisterTicketServiceServer(server, newTrainServer())
	reflection.Register(server)

	//listen on the port
	listener, err := net.Listen("tcp", ":9080")
	if err != nil {
		fmt.Printf("Failed to listen: %v", err)
		return
	}

	fmt.Println("Server listening on :9080")
	err = server.Serve(listener)
	if err != nil {
		fmt.Printf("Failed to serve: %v", err)
	}
}
