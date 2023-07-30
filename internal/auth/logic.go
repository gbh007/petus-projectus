package auth

import (
	"app/internal/auth/storage"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
)

// Неверный логин/пароль
var LoginOrPasswordIncorrectErr = errors.New("login or password incorrect")

// createUser - создает нового пользователя
func (s *authServer) createUser(ctx context.Context, login, password string) (int64, error) {
	salt := randomSHA256String()
	login = strings.ToLower(login)

	id, err := s.db.CreateUser(ctx, &storage.User{
		Login:    login,
		Password: saltPassword(password, salt),
		Salt:     salt,
	})
	if err != nil {
		return 0, fmt.Errorf("create user: %w", err)
	}

	return id, nil
}

// createSession - создает новую сессию пользователя
func (domain *authServer) createSession(ctx context.Context, login, password string) (string, error) {
	user, err := domain.checkUser(ctx, login, password)
	if err != nil {
		return "", fmt.Errorf("create session: %w", err)
	}

	token := randomSHA256String()

	err = domain.db.CreateSession(ctx, &storage.Session{
		Token:  token,
		UserID: user.ID,
	})
	if err != nil {
		return "", fmt.Errorf("create session: %w", err)
	}

	return token, nil
}

// deleteSession - удаляет сессию пользователя
func (domain *authServer) deleteSession(ctx context.Context, token string) error {
	err := domain.db.DeleteSessionByToken(ctx, token)
	if err != nil {
		return fmt.Errorf("delete session: %w", err)
	}

	return nil
}

// checkUser - проверяет данные пользователя
func (domain *authServer) checkUser(ctx context.Context, login, password string) (*storage.User, error) {
	login = strings.ToLower(login)

	user, err := domain.db.GetUserByLogin(ctx, login)

	// Такого пользователя не существует
	if errors.Is(err, sql.ErrNoRows) {
		return nil, LoginOrPasswordIncorrectErr
	}

	if err != nil {
		return nil, err
	}

	// Проверка пароля
	if saltPassword(password, user.Salt) != user.Password {
		return nil, LoginOrPasswordIncorrectErr
	}

	return user, nil
}
