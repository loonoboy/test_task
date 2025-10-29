package grpc_client

import (
	"context"
	"log"

	"git.amocrm.ru/study_group/in_memory_database/internal/controller/grpc/proto"
	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	client := proto.NewDeleteAccountServiceClient(conn)
	response, err := client.DeleteAccount(context.Background(), &proto.DeleteAccountRequest{
		AccountId: 32726534,
	})
	if err != nil {
		log.Fatalf("could not delete: %v", err)
	}
	log.Println(response)
}
