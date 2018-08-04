package mysql

import (
	"database/sql"

	"github.com/madskrogh/finisafricae"
	"github.com/madskrogh/finisafricae/util"
)

//BookService represents a MySQL implementation of the finisafricae.BookService interface.
type BookService struct {
	DB *sql.DB
}

//Book returns a book for a given id.
func (s *BookService) Book(id string) (*finisafricae.Book, error) {
	var b finisafricae.Book
	row := s.DB.QueryRow(`SELECT * FROM book WHERE id = ?`, id)
	if err := row.Scan(&b.ID, &b.UserID, &b.Title, &b.Author, &b.Year, &b.Genre, &b.Notes); err != nil {
		return nil, err
	}
	return &b, nil
}

//Books returns all book
func (s *BookService) Books(userID string) ([]*finisafricae.Book, error) {
	bs := make([]*finisafricae.Book, 0)
	rows, err := s.DB.Query(`SELECT * FROM book WHERE userid = ?`, userID)
	util.HandleError(err)
	for rows.Next() {
		b := finisafricae.Book{}
		err := rows.Scan(&b.ID, &b.UserID, &b.Title, &b.Author, &b.Year, &b.Genre, &b.Notes)
		if err != nil {
			return nil, err
		}
		bs = append(bs, &b)
	}
	return bs, nil
}

//CreateBook inserts new book into table
func (s *BookService) CreateBook(b *finisafricae.Book) error {
	sqlStatement := `INSERT INTO book (id,userid,title,author,year,genre,notes) VALUES (?, ?, ?, ?, ?, ?, ?)`
	_, err := s.DB.Exec(sqlStatement, &b.ID, &b.UserID, &b.Title, &b.Author, &b.Year, &b.Genre, &b.Notes)
	return err
}

//UpdateBook updates a book in the table
func (s *BookService) UpdateBook(b *finisafricae.Book) error {
	sqlStatement := `UPDATE book SET userid=?, title=?, author=?, year=?, genre=?, notes=? WHERE id=?`
	_, err := s.DB.Exec(sqlStatement, &b.ID, &b.UserID, &b.Title, &b.Author, &b.Year, &b.Genre, &b.Notes, &b.ID)
	return err
}

//DeleteBook deletes record with matching id
func (s *BookService) DeleteBook(id string) error {
	sqlStatement := `DELETE FROM book WHERE id=?`
	_, err := s.DB.Exec(sqlStatement, id)
	return err
}
