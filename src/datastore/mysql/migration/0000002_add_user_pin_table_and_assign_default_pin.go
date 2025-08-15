package migration

import (
	"assignment/util"
	"fmt"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

var addUserPinTableMigration = &Migration{
	Number: 2,
	Name:   "create user pin table",
	Forwards: func(db *gorm.DB) error {
		const sql = `
			CREATE TABLE IF NOT EXISTS user_pin (
				user_id VARCHAR(50) NOT NULL,
				pin VARCHAR(255) NOT NULL,
				created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
				updated_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
				PRIMARY KEY (user_id)
			) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
		`

		err := db.Exec(sql).Error
		if err != nil {
			return errors.Wrap(err, "unable to create user pin table")
		}

		defaultPin, err := util.HashPassword(viper.GetString("DefaultPin"))
		if err != nil {
			return errors.Wrap(err, "unable to hash default pin")
		}
		initUsersPinSql := fmt.Sprintf(`
		INSERT INTO user_pin (user_id, pin)
		SELECT u.user_id, '%s'
		FROM users u
		LEFT JOIN user_pin p ON p.user_id = u.user_id
		WHERE p.user_id IS NULL;
		`, defaultPin)

		err = db.Exec(initUsersPinSql).Error
		if err != nil {
			return errors.Wrap(err, "unable to init default user pin table")
		}
		return nil
	},
}

func init() {
	Migrations = append(Migrations, addUserPinTableMigration)
}
