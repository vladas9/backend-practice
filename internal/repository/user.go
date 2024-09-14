package repository

import (
	"database/sql"
	"github.com/google/uuid"
	m "github.com/vladas9/backend-practice/internal/models"
)

type UserRepo struct {
	tx *sql.Tx
}

func NewUserRepo(tx *sql.Tx) *UserRepo {
	return &UserRepo{tx}
}

func (r *UserRepo) GetById(id uuid.UUID) (*m.UserModel, error) {
	item := &m.UserModel{}
	query := `
		SELECT 
			id,
			username,
			email,
			password,
			address,
			phone_number,
			user_type,
			registered_date
		FROM
			users
		WHERE
			id = $1
	`
	row := r.tx.QueryRow(query, id)
	if err := row.Scan(
		&item.ID,
		&item.Username,
		&item.Email,
		&item.Password,
		&item.Address,
		&item.PhoneNumber,
		&item.UserType,
		&item.RegisteredDate,
	); err != nil {
		return nil, err
	}
	return item, nil
}

func (r *UserRepo) GetByEmail(email string) (*m.UserModel, error) {
	item := &m.UserModel{}
	query := `
		SELECT 
			id,
			username,
			password,
			address,
			phone_number,
			user_type,
			registered_date
		FROM
			users
		WHERE
			email = $1
	`

	row := r.tx.QueryRow(query, email)
	if err := row.Scan(
		&item.ID,
		&item.Username,
		&item.Password,
		&item.Address,
		&item.PhoneNumber,
		&item.UserType,
		&item.RegisteredDate,
	); err != nil {
		return nil, err
	}
	return item, nil
}

func (r *UserRepo) GetAll() ([]*m.UserModel, error) {
	var users []*m.UserModel
	query := `
		SELECT 
			id,
			username,
			email,
			password,
			address,
			phone_number,
			user_type,
			registered_date
		FROM
			users
	`
	rows, err := r.tx.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		item := &m.UserModel{}
		if err := rows.Scan(
			&item.ID,
			&item.Username,
			&item.Email,
			&item.Password,
			&item.Address,
			&item.PhoneNumber,
			&item.UserType,
			&item.RegisteredDate,
		); err != nil {
			return nil, err
		}
		users = append(users, item)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return users, nil
}

func (r *UserRepo) Update(item *m.UserModel) error {
	query := `
		UPDATE 
			users
		SET
			username = $1,
			email = $2,
			address = $3,
			password = $4
			phone_number = $5,
			user_type = $6,
			registered_date = $7
		WHERE
			id = $8
	`
	_, err := r.tx.Exec(query,
		&item.Username,
		&item.Email,
		&item.Address,
		&item.Password,
		&item.PhoneNumber,
		&item.UserType,
		&item.RegisteredDate,
		&item.ID,
	)

	return err
}

func (r *UserRepo) Remove(id uuid.UUID) error {
	query := `
		DELETE FROM 
			users
		WHERE 
			id = $1
	`
	_, err := r.tx.Exec(query, id)

	return err
}

func (r *UserRepo) Insert(item *m.UserModel) (uuid.UUID, error) {
	query := `
        INSERT INTO users (
            username,
            email,
            address,
            password,
            phone_number,
            user_type
        ) VALUES (
            $1, $2, $3, $4, $5, $6
        ) RETURNING id
    `
	var userId uuid.UUID

	err := r.tx.QueryRow(query,
		item.Username,
		item.Email,
		item.Address,
		item.Password,
		item.PhoneNumber,
		item.UserType,
	).Scan(&userId)
	if err != nil {
		return uuid.Nil, err
	}

	return userId, nil
}
