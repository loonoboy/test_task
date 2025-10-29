package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"git.amocrm.ru/study_group/in_memory_database/internal/controller/grpc/proto"
	"git.amocrm.ru/study_group/in_memory_database/internal/controller/grpc/server"
	"git.amocrm.ru/study_group/in_memory_database/internal/controller/http/v1"
	"git.amocrm.ru/study_group/in_memory_database/internal/provider"
	"git.amocrm.ru/study_group/in_memory_database/internal/repository/mysql/account_integrations"
	"git.amocrm.ru/study_group/in_memory_database/internal/repository/mysql/accounts"
	"git.amocrm.ru/study_group/in_memory_database/internal/repository/mysql/contacts"
	"git.amocrm.ru/study_group/in_memory_database/internal/repository/mysql/init_mysql"
	"git.amocrm.ru/study_group/in_memory_database/internal/usecase/account"
	"git.amocrm.ru/study_group/in_memory_database/internal/usecase/account_integration"
	"git.amocrm.ru/study_group/in_memory_database/internal/usecase/amo_client"
	"git.amocrm.ru/study_group/in_memory_database/internal/usecase/contact"
	"git.amocrm.ru/study_group/in_memory_database/internal/usecase/unisender"
	"git.amocrm.ru/study_group/in_memory_database/pkg/amocrm"
	"google.golang.org/grpc"
)

func main() {

	db := init_mysql.NewConnectMySQL()
	accountsRepo := accounts.NewAccountRepoMySQL(db)
	integrationsRepo := account_integrations.NewIntegrationRepoMySQL(db)
	contactsRepo := contacts.NewContactRepoMySQL(db)

	httpClient := &http.Client{Timeout: 20 * time.Second}
	amoClient := amocrm.NewAMOClient(httpClient)
	unisenderProvider := provider.NewUnisenderProvider()

	accountService := account.NewAccountUsecase(accountsRepo)
	integrationService := account_integration.NewAccountInegrationUsecase(integrationsRepo)
	contactService := contact.NewContactUsecase(contactsRepo)
	amoClientService := amo_client.NewAmoClientServiceService(amoClient, accountsRepo)
	unisenderService := unisender.NewUnisenderService(accountsRepo, contactsRepo, unisenderProvider)

	gserver := grpc.NewServer()
	gRPCServerStruct := server.NewGRPCServer(accountService)
	proto.RegisterDeleteAccountServiceServer(gserver, gRPCServerStruct)

	accountHandler := v1.NewAccountHandler(accountService)
	integrationHandler := v1.NewAccountIntegrationHandler(integrationService)
	contactHandler := v1.NewContactHandler(contactService)
	amoClientHandler := v1.NewAmoClientHandler(amoClientService)
	unisenderHandler := v1.NewUnisenderHandler(unisenderService)

	handler := v1.NewHandler(accountHandler, integrationHandler, contactHandler, amoClientHandler, unisenderHandler)
	router := v1.NewRouter(handler)

	addr := "localhost:8080"

	srv := &http.Server{
		Addr:    addr,
		Handler: router,
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Printf("API service starting at %s\n", addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	go func() {
		l, err := net.Listen("tcp", ":50051")
		if err != nil {
			log.Fatal(err)
		}
		if err := gserver.Serve(l); err != nil {
			log.Fatal(err)
		}
	}()

	<-stop
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server stopped gracefully")
}
