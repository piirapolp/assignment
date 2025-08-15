package controller

import (
	"assignment/global"
	model_mysql "assignment/model/mysql"
	"context"
	"strings"
	"unicode/utf8"
)

type GetAccountsOutput struct {
	Accounts []model_mysql.AccountWithDetails `json:"accounts"`
}

func (controller Controller) GetUserAccounts(ctx context.Context) (GetAccountsOutput, error) {
	controller.Logger.Info("getting user accounts")
	output := GetAccountsOutput{}

	var err error
	output.Accounts, err = controller.ModelRepository.GetUserAccounts(ctx, controller.UserId)
	if err != nil {
		controller.Logger.Errorf("get user accounts failed because: %s", err.Error())
		return output, global.SystemError{
			Code:    global.DatabaseError,
			Message: err.Error(),
		}
	}

	controller.Logger.Info("get user accounts completed")

	return output, nil
}

type GetDebitCardsOutput struct {
	DebitCards []model_mysql.CardsWithDetails `json:"debit_cards"`
}

func (controller Controller) GetUserDebitCards(ctx context.Context) (GetDebitCardsOutput, error) {
	controller.Logger.Info("getting user debit cards")
	output := GetDebitCardsOutput{}

	var err error
	output.DebitCards, err = controller.ModelRepository.GetUserCards(ctx, controller.UserId)
	if err != nil {
		controller.Logger.Errorf("get user debit cards failed because: %s", err.Error())
		return output, global.SystemError{
			Code:    global.DatabaseError,
			Message: err.Error(),
		}
	}

	// Masked middle half of card number
	for i := range output.DebitCards {
		parts := strings.Fields(output.DebitCards[i].Number)
		if len(parts) <= 2 {
			continue
		}
		for i := 1; i < len(parts)-1; i++ {
			parts[i] = strings.Repeat("*", utf8.RuneCountInString(parts[i]))
		}
		output.DebitCards[i].Number = strings.Join(parts, " ")
	}
	controller.Logger.Info("get user debit cards completed")

	return output, nil
}

type GetSavedAccountsOutput struct {
	SavedAccounts []model_mysql.SavedAccounts `json:"saved_accounts"`
}

func (controller Controller) GetUserSavedAccounts(ctx context.Context) (GetSavedAccountsOutput, error) {
	controller.Logger.Info("getting user saved accounts")
	output := GetSavedAccountsOutput{}

	var err error
	output.SavedAccounts, err = controller.ModelRepository.GetUserSavedAccounts(ctx, controller.UserId)
	if err != nil {
		controller.Logger.Errorf("get user saved accounts failed because: %s", err.Error())
		return output, global.SystemError{
			Code:    global.DatabaseError,
			Message: err.Error(),
		}
	}

	controller.Logger.Info("get user saved accounts completed")

	return output, nil
}
