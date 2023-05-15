package repository

import (
	entity "github.com/roblesoft/topics/internal/entity"
	"gorm.io/gorm"
)

type UserRepository struct {
	Db *gorm.DB
}

func (r *UserRepository) Get(username string) (*entity.User, error) {
	var b *entity.User

	err := r.Db.Where("username = ?", username).First(&b).Error

	return b, err
}

func (r *UserRepository) Create(user *entity.User) error {
	err := r.Db.Create(&user).Error
	return err
}
