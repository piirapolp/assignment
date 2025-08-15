package model_mysql

import (
	"assignment/datastore/mysql"
	"assignment/entity"
	"context"
	"github.com/shopspring/decimal"
)

type AccountWithDetails struct {
	AccountID     string          `json:"account_id"`
	Type          string          `json:"type"`
	Currency      string          `json:"currency"`
	AccountNumber string          `json:"account_number"`
	Issuer        string          `json:"issuer"`
	Amount        decimal.Decimal `json:"amount"`
	Color         string          `json:"color"`
	IsMainAccount bool            `json:"is_main_account"`
	Progress      int             `json:"progress"`
	Flags         []AccountFlags  `json:"flags" gorm:"-"`
}

type AccountFlags struct {
	FlagType  string `json:"flag_type"`
	FlagValue string `json:"flag_value"`
}

func (repository *ModelMysqlRepository) GetUserAccounts(ctx context.Context, userId string) ([]AccountWithDetails, error) {
	var result []AccountWithDetails
	if err := mysql.DB.WithContext(ctx).
		Table("accounts AS a").
		Select(`
			a.account_id,
			a.type,
			a.currency,
			a.account_number,
			a.issuer,
			ab.amount,
			ad.color,
			ad.is_main_account,
			ad.progress
		`).
		Joins("JOIN account_balances AS ab ON ab.account_id = a.account_id AND ab.user_id = a.user_id").
		Joins("JOIN account_details  AS ad ON ad.account_id = a.account_id AND ad.user_id = a.user_id").
		Where("a.user_id = ?", userId).
		Order("ad.is_main_account DESC, a.account_id").
		Scan(&result).Error; err != nil {
		return nil, err
	}
	if len(result) == 0 {
		return result, nil
	}

	accountIndex := make(map[string]int, len(result))
	accountIds := make([]string, 0, len(result))
	for i, r := range result {
		accountIndex[r.AccountID] = i
		accountIds = append(accountIds, r.AccountID)
	}

	var flagList []entity.AccountFlags
	if err := mysql.DB.WithContext(ctx).
		Where("user_id = ? AND account_id IN ?", userId, accountIds).
		Order("account_id, flag_type, flag_value").
		Find(&flagList).Error; err != nil {
		return nil, err
	}

	for _, flag := range flagList {
		if idx, ok := accountIndex[flag.AccountId]; ok {
			result[idx].Flags = append(result[idx].Flags, AccountFlags{
				FlagType:  flag.FlagType,
				FlagValue: flag.FlagValue,
			})
		}
	}

	return result, nil
}

type CardsWithDetails struct {
	CardId      string `json:"card_id"`
	Name        string `json:"name"`
	Status      string `json:"status"`
	Number      string `json:"number"`
	Issuer      string `json:"issuer"`
	Color       string `json:"color"`
	BorderColor string `json:"border_color"`
}

func (repository *ModelMysqlRepository) GetUserCards(ctx context.Context, userId string) ([]CardsWithDetails, error) {
	var result []CardsWithDetails
	if err := mysql.DB.WithContext(ctx).
		Table("debit_cards AS dc").
		Select(`
			dc.card_id,
			dc.name,
			dc_design.color,
			dc_design.border_color,
			dc_details.number,		
			dc_details.issuer,
			dc_s.status
		`).
		Joins("JOIN debit_card_design AS dc_design ON dc_design.card_id = dc.card_id AND dc_design.user_id = dc.user_id").
		Joins("JOIN debit_card_details AS dc_details ON dc_details.card_id = dc.card_id AND dc_details.user_id = dc.user_id").
		Joins("JOIN debit_card_status AS dc_s ON dc_s.card_id = dc.card_id AND dc_s.user_id = dc.user_id").
		Where("dc.user_id = ?", userId).
		Order("dc.card_id").
		Scan(&result).Error; err != nil {
		return nil, err
	}
	return result, nil
}

type SavedAccounts struct {
	AccountName   string `json:"name"`
	AccountNumber string `json:"number"`
	Image         string `json:"image"`
}

func (repository *ModelMysqlRepository) GetUserSavedAccounts(ctx context.Context, userId string) ([]SavedAccounts, error) {
	var result []SavedAccounts
	if err := mysql.DB.WithContext(ctx).Table(entity.SavedAccounts{}.TableName()).
		Select(`account_name, account_number, image`).
		Where("user_id = ?", userId).Scan(&result).Error; err != nil {
		return nil, err
	}
	return result, nil
}
