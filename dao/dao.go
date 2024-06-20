//go:build !solution

package dao

import (
	"context"
	"database/sql"
	"log"

	"github.com/jackc/pgx/v5"
)

type UsersDao struct {
	conn *pgx.Conn
}

func (u UsersDao) Close() error {
	if u.conn != nil {
		ctx := context.Background()
		return u.conn.Close(ctx)
	}
	return nil
}

func (u UsersDao) Create(ctx context.Context, user *User) (UserID, error) {
	var userID UserID

	err := u.conn.QueryRow(
		ctx,
		`
		INSERT INTO users(name) VALUES ($1) RETURNING id;
		`,
		user.Name,
	).Scan(&userID)

	if err != nil {
		return -1, err
	}

	return userID, nil
}

func (u UsersDao) Delete(ctx context.Context, id UserID) error {
	_, err := u.conn.Exec(
		ctx,
		"DELETE FROM users WHERE id = $1",
		id,
	)
	return err
}

func (u UsersDao) List(ctx context.Context) ([]User, error) {
	var users []User

	rows, err := u.conn.Query(ctx, "SELECT id, name FROM users;")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var id UserID
		var name string

		if err = rows.Scan(&id, &name); err != nil {
			return nil, err
		}
		users = append(users, User{
			ID:   id,
			Name: name,
		})
	}

	err = rows.Err()

	return users, err
}

func (u UsersDao) Lookup(ctx context.Context, id UserID) (User, error) {
	user := User{
		ID: id,
	}

	err := u.conn.QueryRow(ctx, "SELECT name FROM users WHERE id = $1", id).Scan(&user.Name)
	if err != nil && err.Error() == "no rows in result set" {
		return user, sql.ErrNoRows
	}

	return user, nil
}

func (u UsersDao) Update(ctx context.Context, user *User) error {
	_, err := u.Lookup(ctx, user.ID)
	if err != nil {
		return err
	}

	_, err = u.conn.Exec(
		ctx,
		"UPDATE users SET name = $1 WHERE id = $2",
		user.Name, user.ID,
	)
	return err
}

func CreateDao(ctx context.Context, dsn string) (Dao, error) {
	conn, err := pgx.Connect(ctx, dsn)
	if err != nil {
		log.Fatal(err)
	}

	dao := UsersDao{
		conn: conn,
	}

	_, err = conn.Exec(
		ctx,
		`
		CREATE TABLE users(
			id BIGSERIAL PRIMARY KEY,
			name varchar(50)
		);
		`,
	)

	if err != nil {
		conn.Close(ctx)
		return nil, err
	}

	return dao, nil
}
