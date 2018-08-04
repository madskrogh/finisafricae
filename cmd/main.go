package main

import (
	"database/sql"
	"net/http"

	handler "github.com/madskrogh/finisafricae/http"
	"github.com/madskrogh/finisafricae/mysql"
	"github.com/madskrogh/finisafricae/util"

	_ "github.com/go-sql-driver/mysql"

	"html/template"
)

var Templates *template.Template

func init() {
	Templates = template.Must(template.ParseGlob("/path/to/html/templaters"))
}

func main() {

	//Start mysql db
	db, err := sql.Open("mysql", "user:password@/database")
	util.HandleError(err)
	defer db.Close()

	mysql.InitDB(db)

	//Initialize services and inject the db
	us := &mysql.UserService{DB: db}
	bs := &mysql.BookService{DB: db}
	ss := &mysql.SessionService{DB: db}

	//Http router
	http.Handle("/", &handler.IndexHandler{UserService: us, SessionService: ss, Templates: Templates})
	http.Handle("/home", &handler.HomeHandler{UserService: us, SessionService: ss, BookService: bs, Templates: Templates})
	http.Handle("/book", &handler.BookHandler{UserService: us, SessionService: ss, Templates: Templates})
	http.Handle("/newbook", &handler.NewBookHandler{UserService: us, SessionService: ss, BookService: bs, Templates: Templates})
	http.Handle("/savebook", &handler.SaveBookHandler{UserService: us, SessionService: ss, BookService: bs, Templates: Templates})
	http.Handle("/login", &handler.LoginHandler{UserService: us, SessionService: ss, Templates: Templates})
	http.Handle("/logout", &handler.LogoutHandler{UserService: us, SessionService: ss})
	http.Handle("/signup", &handler.SignupHandler{UserService: us, SessionService: ss, Templates: Templates})
	http.Handle("/user", &handler.UserHandler{UserService: us, SessionService: ss, Templates: Templates})
	http.Handle("/updatepassword", &handler.UpdatePasswordHandler{UserService: us, SessionService: ss, Templates: Templates})
	http.Handle("/updateemail", &handler.UpdateEmailHandler{UserService: us, SessionService: ss, Templates: Templates})
	http.Handle("/share", &handler.ShareHandler{UserService: us, SessionService: ss, Templates: Templates})
	http.Handle("/favicon.ico", http.NotFoundHandler())
	http.ListenAndServe(":8080", nil)
}
