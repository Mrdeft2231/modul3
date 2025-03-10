package service

import (
	"fmt"
	"log"
	"rest/internal/http/repository"
	"rest/internal/model"
	"rest/pkg/auth"
)

type UserServiceInterface interface {
	ServiceCreateUsers(users string, email string, password string) (*model.User, error)
	UserAuthService(users, email, password string) (*model.User, error)
}

type UserService struct {
	repo repository.RepoInterface
}

func NewService(repo repository.RepoInterface) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) ServiceCreateUsers(users string, email string, password string) (*model.User, error) {

	passwordHah, err := auth.HashPassword(password)
	if err != nil {
		log.Fatalf("Не удалось зашифровать пароль, %v", err)
	}

	user := &model.User{
		Login:    users,
		Email:    email,
		Password: passwordHah,
	}

	s.repo.RepoCreateUser(user.Login, user.Email, user.Password)

	return user, nil

}

func (s *UserService) UserAuthService(users, email, password string) (*model.User, error) {

	user, err := s.repo.RepoAuthUser(users, email)
	if err != nil {
		fmt.Errorf("не удалось отправить данные, %v", err)
		return nil, err
	}

	bol := auth.CheckPassword(user.Password, password)
	if bol {
		return user, nil
	}

	return nil, err
}
