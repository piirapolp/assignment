package entity

import "time"

type Tokens struct {
	SessionId string    `json:"session_id" gorm:"column:session_id; type:VARCHAR(255); primaryKey"`
	UserId    string    `json:"user_id" gorm:"column:user_id; type:VARCHAR(50)"`
	IssuedAt  time.Time `json:"issued_at" gorm:"<-:create; column:issued_at; not null; autoCreateTime"`
	ExpiredAt time.Time `json:"expired_at"`
}

func (Tokens) TableName() string { return "tokens" }
