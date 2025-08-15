package entity

import "time"

type UserPin struct {
	UserId    string    `json:"user_id" gorm:"column:user_id; type:VARCHAR(50); primaryKey"`
	Pin       string    `json:"pin" gorm:"column:pin; type:VARCHAR(255)"`
	CreatedAt time.Time `json:"created_at" gorm:"<-:create; column:created_at; not null; autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"column:updated_at; not null; autoUpdateTime"`
}

func (UserPin) TableName() string { return "user_pin" }

type Users struct {
	UserId    string `json:"user_id" gorm:"column:user_id; type:VARCHAR(50); primaryKey"`
	Name      string `json:"name" gorm:"column:name; type:VARCHAR(100)"`
	DummyCol1 string `json:"dummy_col_1" gorm:"column:dummy_col_1; type:VARCHAR(255)"`
}

func (Users) TableName() string { return "users" }

type UserGreetings struct {
	UserId    string `json:"user_id" gorm:"column:user_id; type:VARCHAR(50); primaryKey"`
	Greeting  string `json:"greeting" gorm:"column:greeting; type:text"`
	DummyCol2 string `json:"dummy_col_2" gorm:"column:dummy_col_2; type:VARCHAR(255)"`
}

func (UserGreetings) TableName() string { return "user_greetings" }
