package unisender

import (
	"log"

	"git.amocrm.ru/study_group/in_memory_database/internal/usecase"
	"git.amocrm.ru/study_group/in_memory_database/internal/usecase/dto"
)

type UnisenderService struct {
	AccountRepo usecase.AccountRepository
	ContactRepo usecase.ContactRepository
}

func NewUnisenderService(accountRepo usecase.AccountRepository, contactRepo usecase.ContactRepository) *UnisenderService {
	return &UnisenderService{
		AccountRepo: accountRepo,
		ContactRepo: contactRepo,
	}
}

func (u *UnisenderService) SaveUnisenderKey(id int, update dto.UpdateAccount) error {
	err := u.AccountRepo.UpdateAccount(id, update)
	if err != nil {
		log.Println(err)
	}
	return nil
}
