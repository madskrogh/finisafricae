package mysql

import (
	"database/sql"

	"github.com/madskrogh/finisafricae"
	"github.com/madskrogh/finisafricae/util"
)

//SessionService represents a MySQL implementation of the finisafricae.SessionService interface.
type SessionService struct {
	DB *sql.DB
}

//Session returns a Session for a given id.
func (s *SessionService) Session(id string) (*finisafricae.Session, error) {
	var se finisafricae.Session
	row := s.DB.QueryRow(`SELECT * FROM session WHERE id = ?`, id)
	if err := row.Scan(&se.ID, &se.UserID, &se.Time); err != nil {
		return nil, err
	}
	return &se, nil
}

//Sessions returns all Sessions in the database
func (s *SessionService) Sessions() ([]*finisafricae.Session, error) {
	ses := make([]*finisafricae.Session, 0)
	rows, err := s.DB.Query(`SELECT * FROM session`)
	util.HandleError(err)
	for rows.Next() {
		se := finisafricae.Session{}
		err := rows.Scan(&se.ID, &se.UserID, &se.Time) // order matters
		if err != nil {
			return nil, err
		}
		ses = append(ses, &se)
	}
	return ses, nil
}

//CreateSession inserts new Session into table
func (s *SessionService) CreateSession(se *finisafricae.Session) error {
	sqlStatement := `INSERT INTO session (id,userid,time) VALUES (?, ?, ?)`
	_, err := s.DB.Exec(sqlStatement, &se.ID, &se.UserID, &se.Time)
	return err
}

//UpdateSession updates a Session in the table
func (s *SessionService) UpdateSession(se *finisafricae.Session) error {
	sqlStatement := `UPDATE session SET id=?, userid=?, time=? WHERE id = ?`
	_, err := s.DB.Exec(sqlStatement, se.ID, se.UserID, se.Time, se.ID)
	return err
}

//DeleteSession deletes record with matching id from table
func (s *SessionService) DeleteSession(id string) error {
	sqlStatement := `DELETE FROM session WHERE id=?`
	_, err := s.DB.Exec(sqlStatement, id)
	return err
}
