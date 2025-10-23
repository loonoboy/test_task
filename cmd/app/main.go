package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"git.amocrm.ru/study_group/in_memory_database/internal/controller/http/v1"
	"git.amocrm.ru/study_group/in_memory_database/internal/repository/account_integrations"
	"git.amocrm.ru/study_group/in_memory_database/internal/repository/accounts"
	"git.amocrm.ru/study_group/in_memory_database/internal/usecase/account"
	"git.amocrm.ru/study_group/in_memory_database/internal/usecase/account_integration"
	"git.amocrm.ru/study_group/in_memory_database/internal/usecase/contact"
	"git.amocrm.ru/study_group/in_memory_database/pkg/amocrm"
)

func main() {
	accountsRepo := accounts.NewAccountsRepository()
	integrationsRepo := account_integrations.NewIntegrationsRepository()

	httpClient := &http.Client{Timeout: 20 * time.Second}
	amoClient := amocrm.NewAMOClient(httpClient)

	accountService := account.NewAccountUsecase(accountsRepo, integrationsRepo, amoClient)
	integrationService := account_integration.NewAccountInegrationUsecase(integrationsRepo)
	contactServic := contact.NewContactsService(amoClient, accountsRepo)

	accountHandler := v1.NewAccountHandler(accountService)
	integrationHandler := v1.NewAccountIntegrationHandler(integrationService)
	contactHandler := v1.NewContactHandler(contactServic)

	handler := v1.NewHandler(accountHandler, integrationHandler, contactHandler)
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

	<-stop
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server stopped gracefully")
}
