package repository

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/lib/pq"
	"rest/internal/model"
)

type RepoInterface interface {
	RepoCreateUser(user string, email string, password string)
	RepoAuthUser(username, email string) (*model.User, error)
}

type Repository struct {
	sql *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{sql: db}
}

func (r *Repository) RepoCreateUser(user string, email string, password string) {
	query := "INSERT INTO Users (username, email, password_hash) VALUES ($1, $2, $3)"

	_, err := r.sql.Exec(query, user, email, password)
	if err != nil {
		fmt.Errorf("Не удалось добавить данные в БД %v", err)
	}
}

func (r *Repository) RepoAuthUser(username, email string) (*model.User, error) {
	query := "SELECT username, email, password_hash FROM users WHERE username = $1 AND email = $2"

	user := &model.User{}

	err := r.sql.QueryRow(query, username, email).Scan(&user.Login, &user.Email, &user.Password)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("пользователь не найден")
		}
		return nil, fmt.Errorf("ошибка запроса: %v", err)
	}

	return user, nil
}

func (r *Repository) RepoChangePassword(id int, password string) {
	query := "SELECT password FROM users where id = $1"
	var oldPassword string

	err := r.sql.QueryRow(query, id).Scan(oldPassword)
	if err != nil {
		fmt.Errorf("не удалось получить старый пароль %v", err)
	}
}
