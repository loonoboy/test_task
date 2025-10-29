package v1

import (
	"fmt"
	"net/http"
	"strconv"

	"git.amocrm.ru/study_group/in_memory_database/internal/usecase/dto"
)

type UnisenderHandler struct {
	usecase UnisenderUsecaseInterface
}

func NewUnisenderHandler(usecase UnisenderUsecaseInterface) *UnisenderHandler {
	return &UnisenderHandler{usecase: usecase}
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
	err := h.usecase.SaveUnisenderKey(AccountID, update)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	h.usecase.SaveExistingContacts(AccountID)
	w.WriteHeader(http.StatusOK)
}
