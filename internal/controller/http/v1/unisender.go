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

type UnisenderHandler struct {
	usecase controller.UnisenderUsecaseInterface
	queue   controller.QueueInterface
}

func NewUnisenderHandler(usecase controller.UnisenderUsecaseInterface, queue controller.QueueInterface) *UnisenderHandler {
	return &UnisenderHandler{
		usecase: usecase,
		queue:   queue,
	}
}

func (h *UnisenderHandler) SynchronizationContacts(w http.ResponseWriter, r *http.Request) {
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
	body, _ := json.Marshal(map[string]int{"account_id": AccountID})
	log.Println("adding a job for first sync for account ", AccountID)
	if h.queue.AddJob(body, "first_sync") != nil {
		log.Fatal("error adding job")
	}
	err := h.usecase.SaveUnisenderKey(AccountID, update)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	h.usecase.SaveExistingContacts(AccountID)
	w.WriteHeader(http.StatusOK)
}
