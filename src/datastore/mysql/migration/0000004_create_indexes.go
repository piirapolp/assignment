package migration

import (
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

var createIndexesMigration = &Migration{
	Number: 4,
	Name:   "create indexes",
	Forwards: func(db *gorm.DB) error {

		const sql = `
			CREATE INDEX idx_banners_user ON banners (user_id);
			CREATE INDEX idx_account_user ON accounts (user_id);
			CREATE INDEX idx_ab_user ON account_balances (user_id, account_id);
			CREATE INDEX idx_ad_user ON account_details (user_id, account_id);
			CREATE INDEX idx_af_user ON account_flags (user_id, account_id);
			CREATE INDEX idx_dc_user_card ON debit_cards (user_id);
			CREATE INDEX idx_dc_design_user_card ON debit_card_design (user_id, card_id);
			CREATE INDEX idx_dc_details_user_card ON debit_card_details (user_id, card_id);
			CREATE INDEX idx_dc_s_user_card ON debit_card_status (user_id, card_id);
			CREATE INDEX idx_sa_user ON saved_accounts (user_id);
		`

		err := db.Exec(sql).Error
		if err != nil {
			return errors.Wrap(err, "unable to create indexes")
		}
		return nil
	},
}

func init() {
	Migrations = append(Migrations, createIndexesMigration)
}
