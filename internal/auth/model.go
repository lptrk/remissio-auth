package auth

import "gorm.io/gorm"

type User struct {
	gorm.Model
	ID           string  `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	Firstname    *string `gorm:"type:varchar(100)" json:"firstname,omitempty"`
	Lastname     *string `gorm:"type:varchar(100)" json:"lastname,omitempty"`
	Username     string  `gorm:"type:varchar(50);unique;not null" json:"username"`
	Email        string  `gorm:"type:varchar(100);unique;not null" json:"email"`
	Password     string  `gorm:"type:varchar(100);unique;not null" json:"password"`
	Age          *int    `gorm:"null" json:"age,omitempty"`
	Gender       *string `gorm:"type:varchar(10)" json:"gender,omitempty"`
	DateOfBirth  *string `gorm:"type:date" json:"date_of_birth,omitempty"`
	SessionToken string  `gorm:"type:varchar(250)" json:"sessionToken,omitempty"`
	CSRFToken    string  `gorm:"type:varchar(250)" json:"csrfToken,omitempty"`
}
