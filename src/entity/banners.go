package entity

type Banners struct {
	BannerId    string `json:"banner_id" gorm:"column:banner_id; type:VARCHAR(50); primaryKey"`
	UserId      string `json:"-" gorm:"column:user_id; type:VARCHAR(50)"`
	Title       string `json:"title" gorm:"column:title; type:VARCHAR(255)"`
	Description string `json:"description" gorm:"column:description; type:text"`
	Image       string `json:"image" gorm:"column:image; type:VARCHAR(255)"`
	DummyCol11  string `json:"-" gorm:"column:dummy_col_11; type:VARCHAR(255)"`
}

func (Banners) TableName() string { return "banners" }
