# finis Africae 
A web backend for simple management of books. After signing up, users can login and view their collection of books, create new books and do simple account management, updating email and password. Finis Africae uses the [bcrypt package](https://godoc.org/golang.org/x/crypto/bcrypt) to store hashed versions of passwords and uses cookies to keep track of sessions. It includes a basic html user interface for the purpose of demonstration. 

The project structure is modeled after the [Standard Package Layout](https://medium.com/@benbjohnson/standard-package-layout-7cdbc8391fc1) which allows for isolation of dependencies and easy implementation of different database solutions. For instance, the MySQL implementation of finisAfricae.UserService (mysql.UserService) used here could relatively easily be substituted by a PostgreSQL or MongoDB solution as long as said solutions satisfy the defined interface. 

The project is a work in progress and feedback/review is highly appreciated. 

Future features to add include:
* Book/library sharing between users
* Updating and deletion of existing books
* Tags and lists 

The name finis Africae refers to [The Name of the Rose](https://en.wikipedia.org/wiki/The_Name_of_the_Rose)