package main

import (
	"database/sql"
)

type user struct {
	ID 		int		`json:"id`
	UserName 	string		`json:"userName"`
	Password	string		`json:"password"`
}

func (u *user) getUser(db *sql.DB) error {
	return db.QueryRow("SELECT id, userName, password FROM users WHERE userName=$1 and password=$2",
		u.ID).Scan(&u.UserName, &u.Password)
}

func getUsers(db *sql.DB) ([]user, error) {
	rows, err := db.Query(
		"SELECT id, userName, password FROM users")

	if (err != nil) {
		return nil, err
	}

	defer rows.Close()

	users := []user{}

	for rows.Next() {
		var u user
		if err := rows.Scan(&u.ID, &u.UserName, &u.Password); err != nil {
			return nil, err
		}

		users = append(users, u)
	}

	return users, nil
}

func (u *user) updatePassword(db *sql.DB) error {
	_, err :=
		db.Exec("UPDATE users SET password=$2 WHERE userName=$1", u.UserName, u.Password)

	return err
}

func (u *user) deleteUser(db *sql.DB) error {
	_, err :=
		db.Exec("DELETE FROM users WHERE id=$1", u.ID)

	return err
}

func (u *user) createUser(db *sql.DB) error {
	_, err :=
		db.Exec("INSERT INTO users(userName, password) VALUES ($1, $2)", u.UserName, u.Password)

	if err != nil {
		return err
	}

	return nil
}