package model

import (
	"time"
)

//PretratamientoDisp struct for company table at the database
type PretratamientoDisp struct {
	ForDate  time.Time `gorm:"type:datetime;default:getdate();"`
	OnHour   int       `gorm:"type:int;default:0;"`
	Total    float64   `gorm:"type:int;default:0;"`
	TagIndex int       `gorm:"type:int;default:0;"`
}

//PretratamientoDispWeek struct for company table at the database
type PretratamientoDispWk struct {
	ForDate time.Time `gorm:"type:datetime;default:getdate();"`
	WeekDay int       `gorm:"type:int;"`
	Total   float64   `gorm:"type:int;"`
	Tag     int       `gorm:"type:int;"`
}

//PretratamientoDispMonth struct for company table at the database
type PretratamientoDispMonth struct {
	OnYear   int     `gorm:"type:int;"`
	OnWeek   int     `gorm:"type:int;"`
	Total    float64 `gorm:"type:int;default:0;"`
	TagIndex int     `gorm:"type:int;default:0;"`
}

//RegisterFilter struct for company table at the database
type RegisterFilter struct {
	SelectedDate time.Time
	TagIndex     int
	TableName    string
	TypeReport   int
}

//Company struct for company table at the database
type Company struct {
	CompanyID           string    `gorm:"primaryKey;type:varchar(50);"`
	CompanyName         string    `gorm:"type:varchar(100);"`
	CompanyAddress      string    `gorm:"type:varchar(100);"`
	CompanyTel          string    `gorm:"type:varchar(50);"`
	companyStatus       int       `gorm:"type:bit;default:1;"`
	companyCreationDate time.Time `gorm:"type:datetime;default:getdate();"`
}

//AppUser struct for app_user table at the database
type AppUser struct {
	AppUserID           string    `gorm:"primaryKey;type:varchar(50);"`
	AppUserName         string    `gorm:"type:varchar(50);"`
	AppUserLastName     string    `gorm:"type:varchar(100);default:'';"`
	AppUserPassword     string    `gorm:"type:text;"`
	AppUserURLPhoto     string    `gorm:"type:text;"`
	AppUserEmail        string    `gorm:"type:varchar(150);default:'';"`
	AppUserStatus       bool      `gorm:"type:bit;default:1;"`
	AppUserCreationDate time.Time `gorm:"type:datetime;default:getdate();"`
	CompanyID           string    `gorm:"foreignKey:CompanyID;type:varchar(50);"`
}

//AppProfile struct for profile table at the database
type AppProfile struct {
	ProfileID           string    `gorm:"primaryKey;type:varchar(60);"`
	CompanyID           string    `gorm:"primaryKey;type:varchar(50);"`
	ProfileDescription  string    `gorm:"type:text;"`
	ProfileStatus       bool      `gorm:"type:bit;default:1;"`
	ProfileCreationDate time.Time `gorm:"type:datetime;default:getdate();"`
}

//UserProfile struct for user_profile table at the database
type UserProfile struct {
	AppUserID               string    `gorm:"primaryKey;type:varchar(50);"`
	ProfileID               string    `gorm:"primaryKey;type:varchar(60);"`
	CompanyID               string    `gorm:"primaryKey;type:varchar(50);"`
	UserProfileStatus       bool      `gorm:"type:bit;default:1;"`
	UserProfileCreationDate time.Time `gorm:"type:datetime;default:getdate();"`
}

//RecordTable struct for record_table table at the database
type RecordTable struct {
	RecordTableName     string    `gorm:"type:varchar(50);"`
	RecordTableAction   string    `gorm:"type:varchar(50);"`
	AppUserID           string    `gorm:"type:varchar(50);"`
	RecordTableActionID string    `gorm:"type:varchar(50);"`
	RecordTableDate     time.Time `gorm:"type:datetime;default:getdate();"`
}
