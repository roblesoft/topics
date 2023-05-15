package usecase

import (
	"fmt"

	entity "github.com/roblesoft/topics/internal/entity"
	repo "github.com/roblesoft/topics/internal/usecase/repo"
	token "github.com/roblesoft/topics/pkg/token"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	UserRepo *repo.UserRepository
}

func NewService(UserRepo *repo.UserRepository) *Service {
	return &Service{
		UserRepo: UserRepo,
	}
}

func (s *Service) GetUser(username string) (*entity.User, error) {
	return s.UserRepo.Get(username)
}

func (s *Service) CreateUser(b *entity.User) error {
	return s.UserRepo.Create(b)
}

func (s *Service) verifyPassword(password, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func (s *Service) LoginCheck(username string, password string) (string, error) {
	user, err := s.GetUser(username)

	if err != nil {
		return "", err
	}

	err = s.verifyPassword(password, user.Password)

	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		fmt.Println("verify")
		fmt.Println(err)
		return "", err
	}

	token, err := token.GenerateToken(user.ID)

	if err != nil {
		fmt.Println("generate")
		fmt.Println(err)
		return "", err
	}

	return token, nil

}
