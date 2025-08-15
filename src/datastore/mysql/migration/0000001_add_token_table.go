package migration

import (
	"fmt"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

var addTokensTableMigration = &Migration{
	Number: 1,
	Name:   "create tokens table",
	Forwards: func(db *gorm.DB) error {
		sql := `
			CREATE TABLE IF NOT EXISTS tokens (
				session_id varchar(255) NOT NULL,
				user_id varchar(50) NOT NULL,
				issued_at timestamp NOT NULL,
				expired_at timestamp NOT NULL,
				PRIMARY KEY (session_id)
			) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
		`

		if err := db.Exec(sql).Error; err != nil {
			return errors.Wrap(err, "Unable to create tokens table")
		}

		sql = fmt.Sprintf(`ALTER TABLE tokens ADD CONSTRAINT fk_user_id FOREIGN KEY(user_id) REFERENCES users (user_id) ON DELETE RESTRICT ON UPDATE RESTRICT`)

		if err := db.Exec(sql).Error; err != nil {
			return errors.Wrap(err, "Unable to link foreign key to users table")
		}
		return nil
	},
}

func init() {
	Migrations = append(Migrations, addTokensTableMigration)
}
