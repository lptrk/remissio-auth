package auth

import (
	"log"

	"gorm.io/gorm"
)

type Repository struct {
	DB *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{DB: db}
}

func (r *Repository) Save(u *User) error {
	return r.DB.Create(u).Error
}

func (r *Repository) GetByUsername(username string) (*User, error) {
	var user User
	result := r.DB.Where("username = ?", username).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func (r *Repository) UserAlreadyExists(username, email string) (bool, error) {
	var user User
	result := r.DB.Where("username = ? OR email = ?", username, email).First(&user)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return false, nil
		}
		return false, result.Error
	}

	return true, nil
}

func (r *Repository) SetSessionToken(t, u string) error {
	user, err := r.GetByUsername(u)
	if err != nil {
		log.Fatal("[Error] Error while setting session token")
		return err
	}
	user.SessionToken = t
	return r.DB.Save(&user).Error
}

func (r *Repository) SetCSRFToken(t, u string) error {
	user, err := r.GetByUsername(u)
	if err != nil {
		log.Fatal("[Error] Error while setting session token")
		return err
	}
	user.CSRFToken = t
	return r.DB.Save(&user).Error
}

func (r *Repository) ResetTokensForUser(username string) error {
	user, err := r.GetByUsername(username)
	if err != nil {
		log.Fatalf("[Error] Unable to fetch user: %s", username)
		return err
	}
	user.CSRFToken = ""
	user.SessionToken = ""
	return r.DB.Save(&user).Error
}
