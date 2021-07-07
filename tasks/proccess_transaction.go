package tasks

import (
	"github.com/luk3skyw4lker/codebank-go/DTO"
	"github.com/luk3skyw4lker/codebank-go/domain"
)

type TransactionTask struct {
	TransactionRepository domain.TransactionRepository
}

func NewTransactionTask(repository domain.TransactionRepository) TransactionTask {
	return TransactionTask{TransactionRepository: repository}
}

func (task *TransactionTask) ProcessTransaction(transactionDTO DTO.Transaction) (domain.Transaction, error) {
	creditCard := task.hydrateCreditCard(transactionDTO)
	ccInfo, err := task.TransactionRepository.GetCreditCard(*creditCard)

	if err != nil {
		return domain.Transaction{}, err
	}

	creditCard.ID = ccInfo.ID
	creditCard.Limit = ccInfo.Limit
	creditCard.Balance = ccInfo.Balance

	t := task.newTransaction(transactionDTO, ccInfo)

	t.Validate(creditCard)

	err = task.TransactionRepository.SaveTransaction(*t, *creditCard)

	if err != nil {
		return domain.Transaction{}, err
	}

	// transactionDTO.ID = t.ID
	// transactionDTO.CreatedAt = t.CreatedAt

	// transactionJson, err := json.Marshal(transactionDTO)

	// if err != nil {
	// 	return domain.Transaction{}, err
	// }

	return *t, nil
}

func (task *TransactionTask) hydrateCreditCard(transactionDTO DTO.Transaction) *domain.CreditCard {
	creditCard := domain.NewCreditCard()

	creditCard.Name = transactionDTO.Name
	creditCard.Number = transactionDTO.Number
	creditCard.ExpirationMonth = transactionDTO.ExpirationMonth
	creditCard.ExpirationYear = transactionDTO.ExpirationYear
	creditCard.CVV = transactionDTO.CVV

	return creditCard
}

func (task *TransactionTask) newTransaction(transactionDTO DTO.Transaction, cc domain.CreditCard) *domain.Transaction {
	t := domain.NewTransaction()

	t.CreditCardId = cc.ID
	t.Store = transactionDTO.Store
	t.Amount = transactionDTO.Amount
	t.CreatedAt = transactionDTO.CreatedAt
	t.Description = transactionDTO.Description

	return t
}
