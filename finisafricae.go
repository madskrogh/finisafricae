//Package finisafricae defines the simple datatypes of the application
package finisafricae

type User struct {
	ID       string
	Uname    string
	Email    string
	Password string
}

type UserService interface {
	User(id string) (*User, error)
	Users() ([]*User, error)
	CreateUser(u *User) error
	UpdateUser(u *User) error
	DeleteUser(id string) error
	UserFromEmail(email string) (*User, error)
}

type Book struct {
	ID     string
	UserID string
	Title  string
	Author string
	Year   string
	Genre  string
	Notes  string
}

type BookService interface {
	Book(id string) (*Book, error)
	Books(userId string) ([]*Book, error)
	CreateBook(b *Book) error
	UpdateBook(b *Book) error
	DeleteBook(id string) error
}

type Session struct {
	ID     string
	UserID string
	Time   string
}

//shared type

type SessionService interface {
	Session(id string) (*Session, error)
	Sessions() ([]*Session, error)
	CreateSession(s *Session) error
	UpdateSession(s *Session) error
	DeleteSession(id string) error
}
