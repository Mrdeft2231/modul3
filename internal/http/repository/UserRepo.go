package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq"
	"log"
	"rest/internal/model"
	"rest/pkg/auth"
	"strings"
)

type RepoInterface interface {
	RepoCreateUser(ctx context.Context, user string, email string, password, role string, status int) error
	RepoAuthUser(ctx context.Context, username, email string) (*model.User, error)
	RepoGetPassword(ctx context.Context, id int) (*model.User, error)
	RepoGetUserAll(ctx context.Context) ([]model.User, error)
	RepoDeleteUser(ctx context.Context, id int) error
	RepoChangePassword(ctx context.Context, id int, password string, passwordChanged bool) error
	RepoBlocketUser(ctx context.Context, id, block int) error
	RepoGetUser(ctx context.Context, id int) (*model.User, error)
	RepoUpdateStatus(ctx context.Context, id, status, block int) error
}

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

func (r *Repository) RepoCreateUser(ctx context.Context, user, email, password, role string, status int) error {
	fmt.Println("Репа", user, email, password, role, status)
	query := "INSERT INTO Users (username, email, password_hash, status, role) VALUES ($1, $2, $3, $4, $5)"

	_, err := r.db.Exec(ctx, query, user, email, password, status, role)

	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value") {
			fmt.Println("попал")
			return auth.ErrUniqName
		}
		return fmt.Errorf("не удалось добавить данные в БД %v", err)
	}
	return nil
}

func (r *Repository) RepoAuthUser(ctx context.Context, username, email string) (*model.User, error) {
	query := "SELECT id, username, email, password_hash, block, create_at, password_changed, role FROM users WHERE username = $1 AND email = $2"

	user := &model.User{}

	err := r.db.QueryRow(ctx, query, username, email).Scan(&user.Id, &user.Login, &user.Email, &user.Password, &user.Block, &user.CreateUser, &user.PasswordChanged, &user.Role)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("пользователь не найден")
		}
		return nil, fmt.Errorf("ошибка запроса: %v", err)
	}

	fmt.Println("юзер в контроллере", user)

	return user, nil
}

func (r *Repository) RepoGetPassword(ctx context.Context, id int) (*model.User, error) {
	fmt.Println("id в репе", id)
	query := "SELECT password_hash FROM users where id = $1"
	user := &model.User{}

	err := r.db.QueryRow(ctx, query, id).Scan(&user.Password)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Printf("Ошибка получения пользователя")
		}
	}

	return user, nil
}

func (r *Repository) RepoGetUserAll(ctx context.Context) ([]model.User, error) {
	query := "SELECT id, username, email, password_hash, role, status, create_at FROM users"

	// Выполняем SQL-запрос (без передачи параметров)
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		log.Printf("не удалось выполнить запрос: %v", err)
		return nil, err
	}
	defer rows.Close() // Закрываем rows, чтобы избежать утечки соединений

	var users []model.User
	for rows.Next() {
		var user model.User // Создаём НОВЫЙ экземпляр структуры в каждой итерации
		if err := rows.Scan(&user.Id, &user.Login, &user.Email, &user.Password, &user.Role, &user.Status, &user.CreateUser); err != nil {
			log.Printf("Ошибка сканирования строки: %v", err)
			return nil, err
		}
		users = append(users, user) // Добавляем копию структуры в срез
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	fmt.Println("Полученные пользователи:", users)
	return users, nil
}

func (r *Repository) RepoDeleteUser(ctx context.Context, id int) error {
	query := "DELETE FROM users WHERE id = $1"
	_, err := r.db.Exec(ctx, query, id)
	if err != nil {
		log.Printf("Не удалось удалить %v", err)
		return err
	}
	return nil
}

func (r *Repository) RepoChangePassword(ctx context.Context, id int, password string, passwordChanged bool) error {
	query := "UPDATE users SET password_hash = $1, password_changed = $2 WHERE id = $3"

	_, err := r.db.Exec(ctx, query, password, passwordChanged, id)
	if err != nil {
		log.Printf("Не удалось выпольнить запрос на обновения пароля %v", err)
		return err
	}

	return nil
}

func (r *Repository) RepoBlocketUser(ctx context.Context, id, block int) error {
	query := "UPDATE users SET block = $1 WHERE id = $2"

	_, err := r.db.Exec(ctx, query, block, id)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) RepoGetUser(ctx context.Context, id int) (*model.User, error) {
	query := "SELECT id, username, email, password_hash, create_at, role, status, password_changed, block FROM users WHERE id = $1"

	user := &model.User{}
	err := r.db.QueryRow(ctx, query, id).Scan(
		&user.Id, &user.Login, &user.Email, &user.Password,
		&user.CreateUser, &user.Role, &user.Status,
		&user.PasswordChanged, &user.Block,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Printf("Пользователь с id %d не найден: %v", id, err)
			return nil, err // Возвращаем ошибку, чтобы обработать на верхнем уровне
		}
		log.Printf("Ошибка при получении пользователя: %v", err)
		return nil, err // Любая другая ошибка должна быть обработана выше
	}

	fmt.Println("данные в репозитории:", user)
	return user, nil
}

func (r *Repository) RepoUpdateStatus(ctx context.Context, id, status, block int) error {
	query := "UPDATE users SET status = $1, block = $2 WHERE id = $3"

	log.Printf("Обновление статуса пользователя: id=%d, status=%d, block=%d", id, status, block)

	_, err := r.db.Exec(ctx, query, status, block, id)
	if err != nil {
		log.Printf("Ошибка обновления пользователя id=%d: %v", id, err)
	}
	return err
}
