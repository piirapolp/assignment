package model_mysql

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

func TestGetUserAccounts_Success_WithFlags(t *testing.T) {
	mock, teardown := setupMockDB(t)
	defer teardown()

	repo := &ModelMysqlRepository{}
	ctx := context.Background()
	userID := "user-1"

	mainQuery := `
		SELECT 
			a.account_id,
			a.type,
			a.currency,
			a.account_number,
			a.issuer,
			ab.amount,
			ad.color,
			ad.is_main_account,
			ad.progress
		FROM accounts AS a
		JOIN account_balances AS ab ON ab.account_id = a.account_id AND ab.user_id = a.user_id
		JOIN account_details  AS ad ON ad.account_id = a.account_id AND ad.user_id = a.user_id
		WHERE a.user_id = ?
		ORDER BY ad.is_main_account DESC, a.account_id
	`

	mainRows := sqlmock.NewRows([]string{"account_id", "type", "currency", "account_number", "issuer", "amount", "color", "is_main_account", "progress"}).
		AddRow("acc-1", "savings", "USD", "111111", "BANKX", "123.45", "blue", 1, 80).
		AddRow("acc-2", "checking", "EUR", "222222", "BANKY", "0", "red", 0, 10)

	mock.ExpectQuery(mainQuery).WithArgs(userID).WillReturnRows(mainRows)

	flagsQuery := "SELECT * FROM `account_flags` WHERE user_id = ? AND account_id IN (?,?) ORDER BY account_id, flag_type, flag_value"
	flagRows := sqlmock.NewRows([]string{"user_id", "account_id", "flag_type", "flag_value"}).
		AddRow(userID, "acc-1", "restricted", "true").
		AddRow(userID, "acc-1", "vip", "gold").
		AddRow(userID, "acc-2", "promo", "summer-2025")

	mock.ExpectQuery(flagsQuery).WithArgs(userID, "acc-1", "acc-2").WillReturnRows(flagRows)

	accounts, err := repo.GetUserAccounts(ctx, userID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(accounts) != 2 {
		t.Fatalf("expected 2 accounts, got %d", len(accounts))
	}

	// Order: is_main_account DESC first, then by account_id
	if accounts[0].AccountID != "acc-1" || !accounts[0].IsMainAccount {
		t.Fatalf("expected first account to be acc-1 (main), got %+v", accounts[0])
	}
	if !accounts[0].Amount.Equal(decimal.RequireFromString("123.45")) {
		t.Fatalf("expected amount 123.45, got %s", accounts[0].Amount.String())
	}
	if len(accounts[0].Flags) != 2 {
		t.Fatalf("expected 2 flags for acc-1, got %d", len(accounts[0].Flags))
	}

	if accounts[1].AccountID != "acc-2" || accounts[1].IsMainAccount {
		t.Fatalf("expected second account to be acc-2 (not main), got %+v", accounts[1])
	}
	if !accounts[1].Amount.Equal(decimal.RequireFromString("0")) {
		t.Fatalf("expected amount 0, got %s", accounts[1].Amount.String())
	}
	if len(accounts[1].Flags) != 1 {
		t.Fatalf("expected 1 flag for acc-2, got %d", len(accounts[1].Flags))
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestGetUserAccounts_Empty_NoFlagsQuery(t *testing.T) {
	mock, teardown := setupMockDB(t)
	defer teardown()

	repo := &ModelMysqlRepository{}
	ctx := context.Background()
	userID := "user-2"

	mainQuery := `
		SELECT 
			a.account_id,
			a.type,
			a.currency,
			a.account_number,
			a.issuer,
			ab.amount,
			ad.color,
			ad.is_main_account,
			ad.progress
		FROM accounts AS a
		JOIN account_balances AS ab ON ab.account_id = a.account_id AND ab.user_id = a.user_id
		JOIN account_details  AS ad ON ad.account_id = a.account_id AND ad.user_id = a.user_id
		WHERE a.user_id = ?
		ORDER BY ad.is_main_account DESC, a.account_id
	`

	mock.ExpectQuery(mainQuery).WithArgs(userID).WillReturnRows(
		sqlmock.NewRows([]string{
			"account_id", "type", "currency", "account_number", "issuer", "amount", "color", "is_main_account", "progress",
		}),
	)

	accounts, err := repo.GetUserAccounts(ctx, userID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(accounts) != 0 {
		t.Fatalf("expected 0 accounts, got %d", len(accounts))
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestGetUserAccounts_QueryError(t *testing.T) {
	mock, teardown := setupMockDB(t)
	defer teardown()

	repo := &ModelMysqlRepository{}
	ctx := context.Background()
	userID := "user-err"

	mainQuery := `
		SELECT 
			a.account_id,
			a.type,
			a.currency,
			a.account_number,
			a.issuer,
			ab.amount,
			ad.color,
			ad.is_main_account,
			ad.progress
		FROM accounts AS a
		JOIN account_balances AS ab ON ab.account_id = a.account_id AND ab.user_id = a.user_id
		JOIN account_details  AS ad ON ad.account_id = a.account_id AND ad.user_id = a.user_id
		WHERE a.user_id = ?
		ORDER BY ad.is_main_account DESC, a.account_id
	`

	mock.ExpectQuery(mainQuery).WithArgs(userID).WillReturnError(gorm.ErrInvalidDB)

	_, err := repo.GetUserAccounts(ctx, userID)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestGetUserCards_Success(t *testing.T) {
	mock, teardown := setupMockDB(t)
	defer teardown()

	repo := &ModelMysqlRepository{}
	ctx := context.Background()
	userID := "user-cards"

	cardsQuery := `
		SELECT 
			dc.card_id,
			dc.name,
			dc_design.color,
			dc_design.border_color,
			dc_details.number,		
			dc_details.issuer,
			dc_s.status
		FROM debit_cards AS dc
		JOIN debit_card_design AS dc_design ON dc_design.card_id = dc.card_id AND dc_design.user_id = dc.user_id
		JOIN debit_card_details AS dc_details ON dc_details.card_id = dc.card_id AND dc_details.user_id = dc.user_id
		JOIN debit_card_status AS dc_s ON dc_s.card_id = dc.card_id AND dc_s.user_id = dc.user_id
		WHERE dc.user_id = ?
		ORDER BY dc.card_id
	`

	rows := sqlmock.NewRows([]string{
		"card_id", "name", "color", "border_color", "number", "issuer", "status",
	}).
		AddRow("card-1", "My Card 1", "green", "darkgreen", "4444333322221111", "VISA", "active").
		AddRow("card-2", "My Card 2", "purple", "violet", "5555444433332222", "MASTERCARD", "inactive")

	mock.ExpectQuery(cardsQuery).WithArgs(userID).WillReturnRows(rows)

	cards, err := repo.GetUserCards(ctx, userID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cards) != 2 {
		t.Fatalf("expected 2 cards, got %d", len(cards))
	}
	if cards[0].CardId != "card-1" || cards[0].Status != "active" {
		t.Fatalf("unexpected first card: %+v", cards[0])
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestGetUserCards_QueryError(t *testing.T) {
	mock, teardown := setupMockDB(t)
	defer teardown()

	repo := &ModelMysqlRepository{}
	ctx := context.Background()
	userID := "user-cards-err"

	cardsQuery := `
		SELECT 
			dc.card_id,
			dc.name,
			dc_design.color,
			dc_design.border_color,
			dc_details.number,		
			dc_details.issuer,
			dc_s.status
		FROM debit_cards AS dc
		JOIN debit_card_design AS dc_design ON dc_design.card_id = dc.card_id AND dc_design.user_id = dc.user_id
		JOIN debit_card_details AS dc_details ON dc_details.card_id = dc.card_id AND dc_details.user_id = dc.user_id
		JOIN debit_card_status AS dc_s ON dc_s.card_id = dc.card_id AND dc_s.user_id = dc.user_id
		WHERE dc.user_id = ?
		ORDER BY dc.card_id
	`

	mock.ExpectQuery(cardsQuery).WithArgs(userID).WillReturnError(gorm.ErrInvalidDB)

	_, err := repo.GetUserCards(ctx, userID)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestGetUserSavedAccounts_Success(t *testing.T) {
	mock, teardown := setupMockDB(t)
	defer teardown()

	repo := &ModelMysqlRepository{}
	ctx := context.Background()
	userID := "user-saved"

	// We don't rely on the exact table name; match the core of the query.
	savedQuery := "SELECT account_name, account_number, image FROM `saved_accounts` WHERE user_id = ?"

	rows := sqlmock.NewRows([]string{"account_name", "account_number", "image"}).
		AddRow("Alice", "1234567890", "https://cdn/img1.png").
		AddRow("Bob", "9876543210", "https://cdn/img2.png")

	mock.ExpectQuery(savedQuery).WithArgs(userID).WillReturnRows(rows)

	accounts, err := repo.GetUserSavedAccounts(ctx, userID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(accounts) != 2 {
		t.Fatalf("expected 2 saved accounts, got %d", len(accounts))
	}
	if accounts[0].AccountName != "Alice" || accounts[0].AccountNumber != "1234567890" {
		t.Fatalf("unexpected first saved account: %+v", accounts[0])
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestGetUserSavedAccounts_QueryError(t *testing.T) {
	mock, teardown := setupMockDB(t)
	defer teardown()

	repo := &ModelMysqlRepository{}
	ctx := context.Background()
	userID := "user-saved-err"

	savedQuery := "SELECT account_name, account_number, image FROM `saved_accounts` WHERE user_id = ?"
	mock.ExpectQuery(savedQuery).WithArgs(userID).WillReturnError(gorm.ErrInvalidDB)

	_, err := repo.GetUserSavedAccounts(ctx, userID)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}
