package service

import (
	"context"
	"fmt"
	"log"
	"rest/internal/http/repository"
	"rest/internal/model"
	"rest/pkg/auth"
	"time"
)

type UserServiceInterface interface {
	ServiceCreateUsers(ctx context.Context, users string, email string, password, role string, status int) (*model.User, error)
	UserAuthService(ctx context.Context, users, email, password string) (*model.User, error)
	GetUsers() ([]model.User, error)
	DeleteUser(ctx context.Context, id int) error
	ChangePassword(ctx context.Context, id int, password, oldPassword string) error
	BlocketUser(ctx context.Context, id int) error
	GetUser(ctx context.Context, id int) (*model.User, error)
}

type UserService struct {
	repo repository.RepoInterface
}

func NewService(repo repository.RepoInterface) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) ServiceCreateUsers(ctx context.Context, users string, email string, password, role string, status int) (*model.User, error) {
	fmt.Println("сервис", users, email, password, role, status)
	passwordHah, err := auth.HashPassword(password)
	if err != nil {
		log.Printf("Не удалось зашифровать пароль, %v", err)
	}

	user := &model.User{
		Login:    users,
		Email:    email,
		Password: passwordHah,
		Status:   status,
		Role:     role,
	}

	if user.Status == 0 {
		user.Status = 1
	}

	err = s.repo.RepoCreateUser(ctx, user.Login, user.Email, user.Password, user.Role, user.Status)
	if err != nil {
		return nil, err
	}

	return user, nil

}

func (s *UserService) UserAuthService(ctx context.Context, users, email, password string) (*model.User, error) {

	user, err := s.repo.RepoAuthUser(ctx, users, email)
	if err != nil {
		return nil, fmt.Errorf("не удалось отправить данные, %v", err)
	}
	fmt.Println("В сервисе", user.CreateUser)
	if user.Block == 5 {
		return nil, auth.ErrBlockStatus
	}

	fmt.Println("юзер в сервисе", user)

	if time.Since(user.CreateUser) > 30*24*time.Hour {
		user.Block = 4
		err := s.repo.RepoBlocketUser(ctx, user.Id, user.Block)
		if err != nil {
			log.Printf("не удалось обновить блокировку польователю, %v", err)
		}
		return nil, auth.ErrBlockDate
	}

	if user.Block == 3 {
		user.Status = 0
		fmt.Println("Сработал блок", user.Id, user.Status, user.Block)
		err := s.repo.RepoUpdateStatus(ctx, user.Id, user.Status, user.Block)
		if err != nil {
			log.Printf("не удалось обновить статус, %v", err)
		}
		return nil,
			auth.ErrBlock
	}
	fmt.Println("id в сервисе", user.Id)
	bol := auth.CheckPassword(user.Password, password)

	if bol {
		return user, nil
	}

	if user.Password != "" {
		user.Block++
		err := s.repo.RepoBlocketUser(ctx, user.Id, user.Block)
		if err != nil {
			log.Printf("не удалось обновить блокировку польователю, %v", err)
		}
		return nil, auth.ErrPassword
	}

	return nil, err
}

func (s *UserService) GetUsers() ([]model.User, error) {
	users, err := s.repo.RepoGetUserAll(context.Background())
	if err != nil {
		log.Printf("не удалось получить данные из репы %v", err)
		return nil, err
	}

	return users, nil
}

func (s *UserService) DeleteUser(ctx context.Context, id int) error {
	err := s.repo.RepoDeleteUser(ctx, id)
	if err != nil {
		log.Printf("Не удалось удалить пользователя %v", err)
		return err
	}
	return nil
}

func (s *UserService) ChangePassword(ctx context.Context, id int, password, oldPassword string) error {
	users, err := s.repo.RepoGetPassword(ctx, id)
	if err != nil {
		log.Printf("Ошибка при получении пароля: %v", err)
		return err
	}

	// Проверяем старый пароль
	if !auth.CheckPassword(users.Password, oldPassword) {
		return auth.ErrIncorrectPassword
	}

	// Проверяем, что новый пароль не совпадает со старым
	if auth.CheckPassword(users.Password, password) {
		return auth.ErrSamePassword
	}

	// Хешируем новый пароль
	newPassword, err := auth.HashPassword(password)
	if err != nil {
		log.Printf("Ошибка хеширования пароля: %v", err)
		return err
	}

	// Обновляем пароль в базе
	users.PasswordChanged = true
	err = s.repo.RepoChangePassword(ctx, id, newPassword, users.PasswordChanged)
	if err != nil {
		return err
	}

	return nil
}

func (s *UserService) BlocketUser(ctx context.Context, id int) error {
	user, err := s.repo.RepoGetUser(ctx, id)
	if err != nil {
		return err
	}

	fmt.Println("в сервисе", user.Block, user.Status)

	if user.Status == 1 {
		user.Status = 0
		user.Block = 5
		err := s.repo.RepoUpdateStatus(ctx, id, user.Status, user.Block)
		if err != nil {
			return err
		}
		return nil
	}

	if user.Status == 0 {
		user.Status = 1
		user.Block = 0
		err := s.repo.RepoUpdateStatus(ctx, id, user.Status, user.Block)
		if err != nil {
			return err
		}
		return nil
	}

	return nil
}

func (s *UserService) GetUser(ctx context.Context, id int) (*model.User, error) {
	user, err := s.repo.RepoGetUser(ctx, id)
	if err != nil {
		return nil, err
	}

	return user, nil
}
