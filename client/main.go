package main

import (
	"context"
	"fmt"
	"io"
	"log"

	pb "github.com/Manish-Saharan/train-ticket-management/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.Dial("localhost:9080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Could not connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewTicketServiceClient(conn)

	for {
		fmt.Println("\nSelect an option:")
		fmt.Println("1. Purchase Ticket")
		fmt.Println("2. View Receipt")
		fmt.Println("3. View Users By Section")
		fmt.Println("4. Remove User")
		fmt.Println("5. Modify Seat")
		fmt.Println("0. Exit")

		choice := promptInt("Enter the option number: ")

		switch choice {
		case 1:
			// Purchase Ticket
			handlePurchaseTicket(client)
		case 2:
			// View Receipt
			handleViewReceipt(client)
		case 3:
			// View Users By Section
			handleViewUsersBySection(client)
		case 4:
			// Remove User
			handleRemoveUser(client)
		case 5:
			// Modify Seat
			handleModifySeat(client)
		case 0:
			// Exit
			fmt.Println("Exiting...")
			return
		default:
			fmt.Println("Invalid option. Please choose a valid option.")
		}
	}
}

func handlePurchaseTicket(client pb.TicketServiceClient) {

	from, to, firstName, lastName, email, pricePaid := getUserInput()

	receipt, err := client.PurchaseTicket(context.Background(), &pb.Receipt{
		From:      from,
		To:        to,
		User:      &pb.User{FirstName: firstName, LastName: lastName, Email: email},
		PricePaid: pricePaid,
	})

	if err != nil {
		log.Fatalf("Error purchasing ticket: %v", err)
	}
	fmt.Printf("Purchase successful! Receipt: %+v\n", receipt)

}

func handleViewReceipt(client pb.TicketServiceClient) {

	email := prompt("Enter the email to view receipt: ")
	viewReceipt, err := client.GetReceipt(context.Background(), &pb.User{Email: email})
	if err != nil {
		log.Fatalf("Error viewing receipt: %v", err)
	}
	fmt.Printf("View Receipt: %+v\n", viewReceipt)
}

func handleViewUsersBySection(client pb.TicketServiceClient) {

	section := prompt("Enter section (A or B): ")
	stream, err := client.GetUsersBySection(context.Background(), &pb.TSection{Sec: section})
	if err != nil {
		log.Fatalf("Error viewing users by section: %v", err)
	}

	fmt.Printf("Users in Section %s:\n", section)
	for {
		receipt, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Error receiving user: %v", err)
		}
		fmt.Printf("%+v\n", receipt)
	}
}

func handleRemoveUser(client pb.TicketServiceClient) {

	email := prompt("Enter the email to remove user: ")
	removedReceipt, err := client.RemoveUser(context.Background(), &pb.User{Email: email})
	if err != nil {
		log.Fatalf("Error removing user: %v", err)
	}
	fmt.Printf("User removed! Receipt: %+v\n", removedReceipt)
}

func handleModifySeat(client pb.TicketServiceClient) {

	from, to, firstName, lastName, email, pricePaid := getUserInput()
	modifiedReceipt, err := client.ModifySeat(context.Background(), &pb.Receipt{
		From:      from,
		To:        to,
		User:      &pb.User{FirstName: firstName, LastName: lastName, Email: email},
		PricePaid: pricePaid,
		Section:   &pb.TSection{Sec: "B"}, // Change section
	})
	if err != nil {
		log.Fatalf("Error modifying seat: %v", err)
	}
	fmt.Printf("Seat modified! Receipt: %+v\n", modifiedReceipt)
}

func getUserInput() (from, to, firstName, lastName, email string, pricePaid int32) {
	from = prompt("From: ")
	to = prompt("To: ")
	firstName = prompt("First Name: ")
	lastName = prompt("Last Name: ")
	email = prompt("Email: ")
	pricePaid = promptInt("Price Paid: ")
	return
}

// prompt prompts the user for input and returns the entered string.
func prompt(promptText string) string {
	fmt.Print(promptText)
	var text string
	fmt.Scanln(&text)
	return text
}

// promptInt prompts the user for an integer input and returns the entered value.
func promptInt(promptText string) int32 {
	for {
		text := prompt(promptText)
		var value int32
		_, err := fmt.Sscanf(text, "%d", &value)
		if err == nil {
			return value
		}
		fmt.Println("Invalid input. Please enter a valid integer.")
	}
}
