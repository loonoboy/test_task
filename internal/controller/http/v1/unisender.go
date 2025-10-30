package v1

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"git.amocrm.ru/study_group/in_memory_database/internal/controller"
	"git.amocrm.ru/study_group/in_memory_database/internal/usecase/dto"
)

const webHookURL = "https://ab5319a02227.ngrok-free.app/webhook"

type UnisenderHandler struct {
	usecase   controller.UnisenderUsecaseInterface
	queue     controller.QueueInterface
	amoClient controller.AmoClientUsecaseInterface
}

func NewUnisenderHandler(usecase controller.UnisenderUsecaseInterface, queue controller.QueueInterface,
	amoClient controller.AmoClientUsecaseInterface) *UnisenderHandler {
	return &UnisenderHandler{
		usecase:   usecase,
		queue:     queue,
		amoClient: amoClient,
	}
}

func (h *UnisenderHandler) Initialization(w http.ResponseWriter, r *http.Request) {
	var update dto.UpdateAccount

	if err := r.ParseForm(); err != nil {
		fmt.Println("not parsed")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var AccountID int
	update.UnisenderKey = r.FormValue("unisender_key")
	idStr := r.FormValue("account_id")
	if idStr != "" {
		AccountID, _ = strconv.Atoi(idStr)
	}

	err := h.usecase.SaveUnisenderKey(AccountID, update)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	body, _ := json.Marshal(map[string]int{"account_id": AccountID})
	log.Println("adding a job for first sync for account ", AccountID)
	if h.queue.AddJob(body, "first_sync") != nil {
		log.Fatal("error adding job")
	}

	err = h.usecase.SaveExistingContacts(AccountID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	log.Println("First sync done")

	domain := r.URL.Query().Get("referer")

	err = h.amoClient.RegisterWebHook(AccountID, webHookURL, domain)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}
