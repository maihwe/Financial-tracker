package storage

import (
	"context"
	"strings"

	"financial-tracker/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

func CreateUserInDB(pool *pgxpool.Pool, user models.User) (models.User, error) {

	email := strings.ToLower(strings.TrimSpace(user.Email))
	name := strings.TrimSpace(user.Name)

	err := pool.QueryRow(
		context.Background(),
		`
		INSERT INTO users (name, email, password_hash)
		VALUES ($1, $2, $3)
		RETURNING id, name, email, password_hash, role, created_at
		`,
		name,
		email,
		user.PasswordHash,
	).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.PasswordHash,
		&user.Role,
		&user.CreatedAt,
	)

	if err != nil {
		return models.User{}, err
	}

	return user, nil
}

func GetUserByEmailFromDB(pool *pgxpool.Pool, email string) (models.User, error) {

	var user models.User

	email = strings.ToLower(strings.TrimSpace(email))

	err := pool.QueryRow(
		context.Background(),
		`
		SELECT id, name, email, password_hash, role, created_at
		FROM users
		WHERE email = $1
		`,
		email,
	).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.PasswordHash,
		&user.Role,
		&user.CreatedAt,
	)

	if err != nil {
		return models.User{}, err
	}

	return user, nil
}

func GetUserByIDFromDB(pool *pgxpool.Pool, id int) (models.User, error) {

	var user models.User

	err := pool.QueryRow(
		context.Background(),
		`
		SELECT id, name, email, password_hash, role, created_at
		FROM users
		WHERE id = $1
		`,
		id,
	).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.PasswordHash,
		&user.Role,
		&user.CreatedAt,
	)

	if err != nil {
		return models.User{}, err
	}

	return user, nil
}

func ListUsersInDB(pool *pgxpool.Pool) ([]models.User, error) {

	rows, err := pool.Query(
		context.Background(),
		`
		SELECT id, name, email, role, created_at
		FROM users
		ORDER BY id
		`,
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var users []models.User

	for rows.Next() {

		var user models.User

		err := rows.Scan(
			&user.ID,
			&user.Name,
			&user.Email,
			&user.Role,
			&user.CreatedAt,
		)

		if err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}
