package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	httpCon "git.amocrm.ru/study_group/in_memory_database/internal/controller/http"
	"git.amocrm.ru/study_group/in_memory_database/internal/repository/account_integrations"
	"git.amocrm.ru/study_group/in_memory_database/internal/repository/accounts"
	"git.amocrm.ru/study_group/in_memory_database/internal/usecase/account"
	"git.amocrm.ru/study_group/in_memory_database/internal/usecase/account_integration"
)

func main() {
	accountsRepo := accounts.NewAccountsRepository()
	integrationsRepo := account_integrations.NewIntegrationsRepository()

	accountService := account.NewAccountUsecase(accountsRepo)
	integrationService := account_integration.NewAccountInegrationUsecase(integrationsRepo)

	accountHandler := httpCon.NewAccountHandler(accountService)
	integrationHandler := httpCon.NewAccountIntegrationHandler(integrationService)

	h := httpCon.NewHandler(accountHandler, integrationHandler)
	mux := http.NewServeMux()
	h.RegisterRoutes(mux)

	addr := "localhost:8080"

	srv := &http.Server{
		Addr:    addr,
		Handler: mux,
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
