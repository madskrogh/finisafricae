package mysql

import (
	"database/sql"

	"github.com/madskrogh/finisafricae"
	"github.com/madskrogh/finisafricae/util"
)

//UserService represents a MySQL implementation of the finisafricae.UserService interaface.
type UserService struct {
	DB *sql.DB
}

//User returns a user for a given id.
func (s *UserService) User(id string) (*finisafricae.User, error) {
	var u finisafricae.User
	row := s.DB.QueryRow(`SELECT * FROM user WHERE id = ?`, id)
	if err := row.Scan(&u.ID, &u.Uname, &u.Email, &u.Password); err != nil {
		return nil, err
	}
	return &u, nil
}

//UserFromEmail returns a user for given email (used for login)
func (s *UserService) UserFromEmail(email string) (*finisafricae.User, error) {
	var u finisafricae.User
	row := s.DB.QueryRow(`SELECT * FROM user WHERE email = ?`, email)
	if err := row.Scan(&u.ID, &u.Uname, &u.Email, &u.Password); err != nil {
		return nil, err
	}
	return &u, nil
}

//Users returns all user in the table
func (s *UserService) Users() ([]*finisafricae.User, error) {
	us := make([]*finisafricae.User, 0)
	rows, err := s.DB.Query(`SELECT * FROM user`)
	util.HandleError(err)
	for rows.Next() {
		u := finisafricae.User{}
		err := rows.Scan(&u.ID, &u.Uname, &u.Email, &u.Password) // order matters
		if err != nil {
			return nil, err
		}
		us = append(us, &u)
	}
	return us, nil
}

//CreateUser inserts new user into table
func (s *UserService) CreateUser(u *finisafricae.User) error {
	sqlStatement := `INSERT INTO user (id, uname, email, password) VALUES (?, ?, ?, ?)`
	_, err := s.DB.Exec(sqlStatement, &u.ID, &u.Uname, &u.Email, &u.Password)
	return err
}

//UpdateUser updates user in table
func (s *UserService) UpdateUser(u *finisafricae.User) error {
	sqlStatement := `UPDATE user SET uname=?, email=?, password=? WHERE id = ?`
	_, err := s.DB.Exec(sqlStatement, u.Uname, u.Email, u.Password, u.ID)
	return err
}

//DeleteUser deletes record with matching id from table
func (s *UserService) DeleteUser(id string) error {
	sqlStatement := `DELETE FROM user WHERE id=?`
	_, err := s.DB.Exec(sqlStatement, id)
	return err
}
