package repositories

import (
	"github.com/helenoktaa/dannisa_sweet_be/config"
	"github.com/helenoktaa/dannisa_sweet_be/models"
)

type UserRepository struct{}

func NewUserRepository() *UserRepository {
	return &UserRepository{}
}

// Find By ID

func (r *UserRepository) FindByID(id string) (*models.User, error) {
	var user models.User

	result := config.DB.
		Preload("Jabatan").
		Where("id_user = ?", id).
		First(&user)

	if result.Error != nil {
		return nil, result.Error
	}

	return &user, nil
}

// Find By Email

func (r *UserRepository) FindByEmail(email string) (*models.User, error) {
	var user models.User

	result := config.DB.
		Preload("Jabatan").
		Where("email = ?", email).
		First(&user)

	if result.Error != nil {
		return nil, result.Error
	}

	return &user, nil
}


// Create User

func (r *UserRepository) Create(user *models.User) error {
	return config.DB.Create(user).Error
}


// Update User


func (r *UserRepository) Update(user *models.User) error {
	return config.DB.Save(user).Error
}

// Find All Users
func (r *UserRepository) FindAll(users *[]models.User) error {
	return config.DB.
		Preload("Jabatan").
		Find(users).Error
}

// Delete User
func (r *UserRepository) Delete(id string) error {
	return config.DB.
		Where("id_user = ?", id).
		Delete(&models.User{}).Error
}