package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"git.amocrm.ru/study_group/in_memory_database/internal/controller/http/v1"
	"git.amocrm.ru/study_group/in_memory_database/internal/repository/mysql/account_integrations"
	"git.amocrm.ru/study_group/in_memory_database/internal/repository/mysql/accounts"
	"git.amocrm.ru/study_group/in_memory_database/internal/repository/mysql/contacts"
	"git.amocrm.ru/study_group/in_memory_database/internal/usecase/account"
	"git.amocrm.ru/study_group/in_memory_database/internal/usecase/account_integration"
	"git.amocrm.ru/study_group/in_memory_database/internal/usecase/amo_client"
	"git.amocrm.ru/study_group/in_memory_database/internal/usecase/contact"
	"git.amocrm.ru/study_group/in_memory_database/pkg/amocrm"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	port := os.Getenv("DB_PORT")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		user, password, host, port, dbName)

	var db *gorm.DB
	var err error

	maxAttempts := 10
	delay := 3 * time.Second

	for i := 1; i <= maxAttempts; i++ {
		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err == nil {
			fmt.Println("Connected to MySQL")
			break
		}
		time.Sleep(delay)
	}

	if err != nil {
		log.Fatalf("Could not connect to DB after %d attempts: %v", maxAttempts, err)
	}

	accountsRepo := accounts.NewAccountRepoMySQL(db)
	integrationsRepo := account_integrations.NewIntegrationRepoMySQL(db)
	contactsRepo := contacts.NewContactRepoMySQL(db)

	httpClient := &http.Client{Timeout: 20 * time.Second}
	amoClient := amocrm.NewAMOClient(httpClient)

	accountService := account.NewAccountUsecase(accountsRepo)
	integrationService := account_integration.NewAccountInegrationUsecase(integrationsRepo)
	contactService := contact.NewContactUsecase(contactsRepo)
	amoClientService := amo_client.NewAmoClientServiceService(amoClient, accountsRepo)

	accountHandler := v1.NewAccountHandler(accountService)
	integrationHandler := v1.NewAccountIntegrationHandler(integrationService)
	contactHandler := v1.NewContactHandler(contactService)
	amoClientHandler := v1.NewAmoClientHandler(amoClientService)

	handler := v1.NewHandler(accountHandler, integrationHandler, contactHandler, amoClientHandler)
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
