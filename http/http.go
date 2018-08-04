//Package http defines the handlers and handlerfunctions for handling http requests and isolates all net/http dependencies
package http

import (
	"html/template"
	"net/http"
	"time"

	"github.com/madskrogh/finisafricae"
	"github.com/madskrogh/finisafricae/util"

	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
)

type IndexHandler struct {
	UserService    finisafricae.UserService
	SessionService finisafricae.SessionService
	Templates      *template.Template
}

func (h *IndexHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if isLoggedIn(h.SessionService, h.UserService, r) {
		http.Redirect(w, r, "/home", http.StatusSeeOther)
		return
	}
	err := h.Templates.ExecuteTemplate(w, "index.gohtml", nil)
	util.HandleError(err)
	return
}

type LoginHandler struct {
	UserService    finisafricae.UserService
	SessionService finisafricae.SessionService
	Templates      *template.Template
}

func (h *LoginHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	//Redirects if user is loggedin or method is GET
	if isLoggedIn(h.SessionService, h.UserService, r) {
		http.Redirect(w, r, "/home", http.StatusSeeOther)
		return
	} else if r.Method == "GET" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	//Parses request form and gets user from db
	err := r.ParseForm()
	util.HandleError(err)
	u, _ := h.UserService.UserFromEmail(r.Form["email"][0])

	if u != nil {
		//Compares hashed password from form with stored password
		err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(r.Form["password"][0]))
		if err == nil {
			//Passwords match. Create new session and cookie for the user.
			sID, _ := uuid.NewV4()
			c := &http.Cookie{
				Name:  "session",
				Value: sID.String(),
			}
			http.SetCookie(w, c)
			t := time.Now().Format(time.RFC1123)
			s := finisafricae.Session{ID: sID.String(), UserID: u.ID, Time: t}
			err = h.SessionService.CreateSession(&s)
			util.HandleError(err)
			http.Redirect(w, r, "/home", http.StatusSeeOther)
			return
		}
	}
	err = h.Templates.ExecuteTemplate(w, "index.gohtml", "Wrong email or password.")
	util.HandleError(err)
	return
}

type LogoutHandler struct {
	UserService    finisafricae.UserService
	SessionService finisafricae.SessionService
}

func (h *LogoutHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !isLoggedIn(h.SessionService, h.UserService, r) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	//Gets cookie and deletes session
	c, err := r.Cookie("session")
	util.HandleError(err)
	c.MaxAge = -1
	h.SessionService.DeleteSession(c.Value)
	http.Redirect(w, r, "/", http.StatusSeeOther)
	return
}

type SignupHandler struct {
	UserService    finisafricae.UserService
	SessionService finisafricae.SessionService
	Templates      *template.Template
}

func (h *SignupHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if isLoggedIn(h.SessionService, h.UserService, r) {
		http.Redirect(w, r, "/home", http.StatusSeeOther)
		return
	} else if r.Method == "GET" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	err := r.ParseForm()
	if r.Form["password"][0] == "" || r.Form["password2"][0] == "" || r.Form["email"][0] == "" || r.Form["uname"][0] == "" {
		//One or more required form fields are empty. User is sent back.
		err := h.Templates.ExecuteTemplate(w, "index.gohtml", "One or more fields are empty. Try again.")
		util.HandleError(err)
		return
	} else if u, _ := h.UserService.UserFromEmail(r.Form["email"][0]); u != nil {
		//A user with the given email is already present in the db. User is sent back.
		err := h.Templates.ExecuteTemplate(w, "index.gohtml", "A user with this email already exist.")
		util.HandleError(err)
		return
	} else if r.Form["password"][0] != r.Form["password2"][0] {
		//The password and repeated password doesn't match. User is sent back.
		err := h.Templates.ExecuteTemplate(w, "index.gohtml", "Passwords does not match. Try again.")
		util.HandleError(err)
		return
	}
	//Form is correctly filled. User created and stored. Redirects to login page.
	uID, err := uuid.NewV4()
	uP, err := bcrypt.GenerateFromPassword([]byte(r.Form["password"][0]), bcrypt.DefaultCost)
	u := finisafricae.User{
		ID:       uID.String(),
		Uname:    r.Form["uname"][0],
		Email:    r.Form["email"][0],
		Password: string(uP),
	}
	err = h.UserService.CreateUser(&u)
	util.HandleError(err)
	err = h.Templates.ExecuteTemplate(w, "index.gohtml", "User was succesfully created. Login to continue.")
	util.HandleError(err)
	return
}

type UserHandler struct {
	UserService    finisafricae.UserService
	SessionService finisafricae.SessionService
	Templates      *template.Template
}

func (h *UserHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !isLoggedIn(h.SessionService, h.UserService, r) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	err := h.Templates.ExecuteTemplate(w, "user.gohtml", nil)
	util.HandleError(err)
	return
}

type UpdatePasswordHandler struct {
	UserService    finisafricae.UserService
	SessionService finisafricae.SessionService
	Templates      *template.Template
}

func (h *UpdatePasswordHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !isLoggedIn(h.SessionService, h.UserService, r) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	err := r.ParseForm()
	util.HandleError(err)
	if r.Form["npassword"][0] != "" || r.Form["npassword2"][0] != "" {
		//Requried fields are filled out. Retrieve cookie, session and current user.
		c, err := r.Cookie("session")
		util.HandleError(err)
		s, err := h.SessionService.Session(c.Value)
		util.HandleError(err)
		u, err := h.UserService.User(s.UserID)
		util.HandleError(err)
		if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(r.Form["password"][0])); err == nil {
			//The given password matches the users password.
			if r.Form["npassword"][0] == r.Form["npassword2"][0] {
				//New password matches repeat password. Hash new password and update current user.
				p, err := bcrypt.GenerateFromPassword([]byte(r.Form["npassword"][0]), bcrypt.MinCost)
				util.HandleError(err)
				u.Password = string(p)
				err = h.UserService.UpdateUser(u)
				util.HandleError(err)
				err = h.Templates.ExecuteTemplate(w, "user.gohtml", "Your password was updated")
				util.HandleError(err)
				return
			}
			//New password and repeat password doesn't match.
			err := h.Templates.ExecuteTemplate(w, "user.gohtml", "The passwords doesn't match. Try again.")
			util.HandleError(err)
			return

		}
		//The give password doesn't match the user
		err = h.Templates.ExecuteTemplate(w, "user.gohtml", "Wrong password.")
		util.HandleError(err)
		return
	}
	http.Redirect(w, r, "/user", http.StatusSeeOther)
	return
}

type UpdateEmailHandler struct {
	UserService    finisafricae.UserService
	SessionService finisafricae.SessionService
	Templates      *template.Template
}

func (h *UpdateEmailHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !isLoggedIn(h.SessionService, h.UserService, r) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	//Retrieve cookie, session and user.
	err := r.ParseForm()
	util.HandleError(err)
	c, err := r.Cookie("session")
	util.HandleError(err)
	s, err := h.SessionService.Session(c.Value)
	util.HandleError(err)
	u, err := h.UserService.User(s.UserID)
	util.HandleError(err)

	if r.Form["email"][0] != "" {
		//Email field is not empty
		if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(r.Form["password"][0])); err == nil {
			//Given password matches the users password
			if u1, _ := h.UserService.UserFromEmail(r.Form["email"][0]); u1 == nil {
				//New email is not taken
				u.Email = r.Form["email"][0]
				err := h.UserService.UpdateUser(u)
				util.HandleError(err)
				m := "Your email was updated to " + u.Email
				err = h.Templates.ExecuteTemplate(w, "user.gohtml", m)
				util.HandleError(err)
				return
			}
			//Email is already taken
			err = h.Templates.ExecuteTemplate(w, "user.gohtml", "A user with this email already exist. Try again.")
			util.HandleError(err)
			return
		}
		//Given password doesn't match the user password
		err = h.Templates.ExecuteTemplate(w, "user.gohtml", "Wrong password.")
		util.HandleError(err)
		return
	}
	http.Redirect(w, r, "/user", http.StatusSeeOther)

}

type HomeHandler struct {
	UserService    finisafricae.UserService
	BookService    finisafricae.BookService
	SessionService finisafricae.SessionService

	Templates *template.Template
}

func (h *HomeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !isLoggedIn(h.SessionService, h.UserService, r) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	err := r.ParseForm()
	util.HandleError(err)
	c, err := r.Cookie("session")
	util.HandleError(err)
	s, err := h.SessionService.Session(c.Value)
	util.HandleError(err)
	books, err := h.BookService.Books(s.UserID)
	util.HandleError(err)
	err = h.Templates.ExecuteTemplate(w, "home.gohtml", books)
	util.HandleError(err)
}

type BookHandler struct {
	UserService    finisafricae.UserService
	BookService    finisafricae.BookService
	SessionService finisafricae.SessionService
	Templates      *template.Template
}

func (h *BookHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !isLoggedIn(h.SessionService, h.UserService, r) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	err := h.Templates.ExecuteTemplate(w, "book.gohtml", nil)
	util.HandleError(err)
}

type NewBookHandler struct {
	UserService    finisafricae.UserService
	BookService    finisafricae.BookService
	SessionService finisafricae.SessionService
	Templates      *template.Template
}

func (h *NewBookHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !isLoggedIn(h.SessionService, h.UserService, r) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	err := h.Templates.ExecuteTemplate(w, "newbook.gohtml", nil)
	util.HandleError(err)
}

type SaveBookHandler struct {
	UserService    finisafricae.UserService
	BookService    finisafricae.BookService
	SessionService finisafricae.SessionService
	Templates      *template.Template
}

func (h *SaveBookHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !isLoggedIn(h.SessionService, h.UserService, r) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	//Retrieve cookie, session and books of current user.
	err := r.ParseForm()
	util.HandleError(err)
	c, err := r.Cookie("session")
	util.HandleError(err)
	s, err := h.SessionService.Session(c.Value)
	util.HandleError(err)
	books, err := h.BookService.Books(s.UserID)
	util.HandleError(err)

	if r.Form["title"][0] == "" {
		//The title field of the form is not filled
		err = h.Templates.ExecuteTemplate(w, "newbook.gohtml", "The book must have a title.")
		util.HandleError(err)
		return
	}
	for i := range books {
		//Ranges through books to see if title already exists (distinct titles are allowed)
		if books[i].Title == r.Form["title"][0] {
			//Title match found. User sendt back to form page.
			err = h.Templates.ExecuteTemplate(w, "newbook.gohtml", "A book with this title already exists.")
			util.HandleError(err)
			return
		}
	}
	//No title found that matches new title. New book created and stored. User sent back to home.
	bID, _ := uuid.NewV4()
	b := finisafricae.Book{
		ID:     bID.String(),
		Title:  r.Form["title"][0],
		UserID: s.UserID,
		Author: r.Form["author"][0],
		Year:   r.Form["year"][0],
		Genre:  r.Form["genre"][0],
		Notes:  r.Form["notes"][0],
	}
	err = h.BookService.CreateBook(&b)
	util.HandleError(err)
	http.Redirect(w, r, "/home", http.StatusSeeOther)
	return
}

type ShareHandler struct {
	UserService    finisafricae.UserService
	BookService    finisafricae.BookService
	SessionService finisafricae.SessionService
	Templates      *template.Template
}

func (h *ShareHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !isLoggedIn(h.SessionService, h.UserService, r) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	err := h.Templates.ExecuteTemplate(w, "share.gohtml", nil)
	util.HandleError(err)
}

//Returns true if user is logged in
func isLoggedIn(SessionService finisafricae.SessionService, UserService finisafricae.UserService, r *http.Request) bool {
	//Parse form and get cookie
	err := r.ParseForm()
	util.HandleError(err)
	c, err := r.Cookie("session")
	if err != nil {
		return false
	}
	//Cookie exists. Get session.
	s, err := SessionService.Session(c.Value)
	if err != nil {
		return false
	}
	//Session exists. Tjek if it is expired.
	t, err := time.Parse(time.RFC1123, s.Time)
	if err != nil {
		return false
	}
	if time.Now().Sub(t) > (time.Second * 300) {
		//Session expired. Delete session and cookie.
		err := SessionService.DeleteSession(s.ID)
		c.MaxAge = -1
		util.HandleError(err)
		return false
	}
	//Session is valid. Update time and session.
	t = time.Now()
	s.Time = t.Format(time.RFC1123)
	err = SessionService.UpdateSession(s)
	util.HandleError(err)
	return true
}
