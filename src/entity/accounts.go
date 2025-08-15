package entity

import (
	"github.com/shopspring/decimal"
	"time"
)

type Accounts struct {
	AccountId     string `json:"account_id" gorm:"column:account_id; type:VARCHAR(50); primaryKey"`
	UserId        string `json:"user_id" gorm:"column:user_id; type:VARCHAR(50)"`
	Type          string `json:"type" gorm:"column:type; type:VARCHAR(50)"`
	Currency      string `json:"currency" gorm:"column:currency; type:VARCHAR(10)"`
	AccountNumber string `json:"account_number" gorm:"column:account_number; type:VARCHAR(20)"`
	Issuer        string `json:"issuer" gorm:"column:issuer; type:VARCHAR(100)"`
	DummyCol3     string `json:"dummy_col_3" gorm:"column:dummy_col_3; type:VARCHAR(255)"`
}

func (Accounts) TableName() string { return "accounts" }

type AccountBalances struct {
	AccountId string          `json:"account_id" gorm:"column:account_id; type:VARCHAR(50); primaryKey"`
	UserId    string          `json:"user_id" gorm:"column:user_id; type:VARCHAR(50)"`
	Amount    decimal.Decimal `json:"amount" gorm:"column:amount; type:DECIMAL(15,2)"`
	DummyCol4 string          `json:"dummy_col_4" gorm:"column:dummy_col_4; type:VARCHAR(255)"`
}

func (AccountBalances) TableName() string { return "account_balances" }

type AccountDetails struct {
	AccountId     string `json:"account_id" gorm:"column:account_id; type:VARCHAR(50); primaryKey"`
	UserId        string `json:"user_id" gorm:"column:user_id; type:VARCHAR(50)"`
	Color         string `json:"color" gorm:"column:color; type:VARCHAR(10)"`
	IsMainAccount bool   `json:"is_main_account" gorm:"column:is_main_account; type:TINYINT(1)"`
	Progress      int    `json:"progress" gorm:"column:progress; type:INT"`
	DummyCol5     string `json:"dummy_col_5" gorm:"column:dummy_col_5; type:VARCHAR(255)"`
}

func (AccountDetails) TableName() string { return "account_details" }

type AccountFlags struct {
	FlagId    int       `json:"flag_id" gorm:"column:flag_id; type:INT; primaryKey; autoIncrement"`
	AccountId string    `json:"account_id" gorm:"column:account_id; type:VARCHAR(50); not null"`
	UserId    string    `json:"user_id" gorm:"column:user_id; type:VARCHAR(50); not null"`
	FlagType  string    `json:"flag_type" gorm:"column:flag_type; type:VARCHAR(50); not null"`
	FlagValue string    `json:"flag_value" gorm:"column:flag_value; type:VARCHAR(30); not null"`
	CreatedAt time.Time `json:"created_at" gorm:"<-:create; column:created_at; autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"column:updated_at; autoUpdateTime"`
}

func (AccountFlags) TableName() string { return "account_flags" }

type SavedAccounts struct {
	UserId        string    `json:"user_id" gorm:"column:user_id; type:VARCHAR(50); primaryKey"`
	AccountName   string    `json:"account_name" gorm:"column:account_name; type:VARCHAR(100)"`
	AccountNumber string    `json:"account_number" gorm:"column:account_number; type:VARCHAR(20)"`
	CreatedAt     time.Time `json:"created_at" gorm:"<-:create; column:created_at; autoCreateTime"`
	UpdatedAt     time.Time `json:"updated_at" gorm:"column:updated_at; autoUpdateTime"`
}

func (SavedAccounts) TableName() string { return "saved_accounts" }
