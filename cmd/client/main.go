package main

import (
	"context"
	"fmt"
	api "github.com/jamalkaksouri/Book-Service-grpc-golang/api/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
	"log"
)

var serverAddress = "localhost:8080"

func main() {
	conn, err := grpc.Dial(serverAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	defer func(conn *grpc.ClientConn) {
		err := conn.Close()
		if err != nil {

		}
	}(conn)

	client := api.NewBookServiceClient(conn)

	bookDto := &api.Book{
		Title:       "Go Language Programming",
		Author:      "John Doe",
		Description: "Go is a robust programming language",
		Language:    "English",
		FinishTime:  timestamppb.Now(),
	}
	resCreated, err := client.CreateBook(context.Background(), &api.CreateBookRequest{Book: bookDto})
	if err != nil {
		fmt.Println(status.FromError(err))
	}
	log.Printf("Book created with id: %d\n", resCreated.Id)

	var id int64 = 1
	resRetrieved, err := client.RetrieveBook(context.Background(), &api.RetrieveBookRequest{Id: id})
	if err != nil {
		fmt.Println(status.FromError(err))
	} else {
		log.Printf("Book retrieved: %v\n", resRetrieved.Book.String())
	}

	var bidUpdate int64 = 3
	bookUpdate := &api.Book{
		Id:          bidUpdate,
		Title:       "Go Programming-updated",
		Author:      "John Doe",
		Description: "Go is a programming language",
		Language:    "English",
		FinishTime:  timestamppb.Now(),
	}
	_, err = client.UpdateBook(context.Background(), &api.UpdateBookRequest{Book: bookUpdate})
	if err != nil {
		errStatus, _ := status.FromError(err)
		fmt.Println(errStatus.Message())
		fmt.Println(errStatus.Code())
	} else {
		log.Printf("Book updated: %v\n", bookUpdate.String())
	}

	var bidDelete int64 = 2
	_, err = client.DeleteBook(context.Background(), &api.DeleteBookRequest{Id: bidDelete})
	if err != nil {
		errStatus, _ := status.FromError(err)
		fmt.Println(errStatus.Message())
		fmt.Println(errStatus.Code())
	} else {
		log.Printf("Book deleted bid: %v\n", bidDelete)
	}

	resList, err := client.ListBook(context.Background(), &api.ListBookRequest{})
	if err != nil {
		errStatus, _ := status.FromError(err)
		fmt.Println(errStatus.Message())
		fmt.Println(errStatus.Code())
	} else {
		log.Printf("Book list: %v\n", resList.Books)
	}
}
