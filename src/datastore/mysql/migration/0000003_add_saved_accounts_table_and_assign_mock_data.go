package migration

import (
	"fmt"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

var addSavedAccountsTableAndAssignMockDataMigration = &Migration{
	Number: 3,
	Name:   "create saved accounts table and assign mock data",
	Forwards: func(db *gorm.DB) error {
		const sql = `
			CREATE TABLE IF NOT EXISTS saved_accounts (
				user_id VARCHAR(50) NOT NULL,
			    account_name VARCHAR(100) NOT NULL,
			    account_number VARCHAR(20) NOT NULL,
				image VARCHAR(255),
				created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
				updated_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
				PRIMARY KEY (user_id)
			) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
		`

		err := db.Exec(sql).Error
		if err != nil {
			return errors.Wrap(err, "unable to create saved accounts table")
		}

		initUserSavedAccountsSql := fmt.Sprintf(`
			INSERT INTO saved_accounts (user_id, account_name, account_number, image) 
			SELECT u.user_id, '%s', '%s', '%s' FROM users u 
			LEFT JOIN saved_accounts s ON s.user_id = u.user_id 
			WHERE s.user_id IS NULL;`,
			"Dummy Name", "1234567890", "https://dummyimage.com/54x54/999/fff")

		err = db.Exec(initUserSavedAccountsSql).Error
		if err != nil {
			return errors.Wrap(err, "unable to init default saved account for users")
		}
		return nil
	},
}

func init() {
	Migrations = append(Migrations, addSavedAccountsTableAndAssignMockDataMigration)
}
