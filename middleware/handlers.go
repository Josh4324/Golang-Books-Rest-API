package middleware

import (
	"database/sql"
	"encoding/json" // package to encode and decode the json into struct and vice versa
	"fmt"
	"go-rest-api/models" // models package where Book schema is defined
	"log"
	"net/http" // used to access the request and response object of the api
	"os"       // used to read the environment variable
	"strconv"  // package used to covert string into int type

	"github.com/gorilla/mux" // used to get the params from the route

	"github.com/joho/godotenv" // package used to read the .env file
	_ "github.com/lib/pq"      // postgres golang driver
)

// response format
type response struct {
	ID      int64  `json:"id,omitempty"`
	Message string `json:"message,omitempty"`
}

// create connection with postgres db
func createConnection() *sql.DB {
	// load .env file
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	// Open the connection
	username := os.Getenv("db_user")
	password := os.Getenv("db_pass")
	dbName := os.Getenv("db_name")
	dbHost := os.Getenv("db_host")
	dbURI := fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable password=%s", dbHost, username, dbName, password) //Build connection string

	db, err := sql.Open("postgres", dbURI)

	if err != nil {
		panic(err)
	}

	// check the connection
	err = db.Ping()

	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully connected!")
	// return the connection
	return db
}

// CreateBook create a book in the postgres db
func CreateBook(w http.ResponseWriter, r *http.Request) {
	// set the header to content type x-www-form-urlencoded
	// Allow all origin to handle cors issue
	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	// create an empty book of type models.Book
	var book models.Book

	// decode the json request to book
	err := json.NewDecoder(r.Body).Decode(&book)

	if err != nil {
		log.Fatalf("Unable to decode the request body.  %v", err)
	}

	// call insert book function and pass the book
	insertID := insertBook(book)

	// format a response object
	res := response{
		ID:      insertID,
		Message: "Book created successfully",
	}

	// send the response
	json.NewEncoder(w).Encode(res)
}

// GetBook will return a single book by its id
func GetBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	// get the bookid from the request params, key is "id"
	params := mux.Vars(r)

	// convert the id type from string to int
	id, err := strconv.Atoi(params["id"])

	if err != nil {
		log.Fatalf("Unable to convert the string into int.  %v", err)
	}

	// call the getBook function with book id to retrieve a single book
	book, err := getBook(int64(id))

	if err != nil {
		log.Fatalf("Unable to get book. %v", err)
	}

	// send the response
	json.NewEncoder(w).Encode(book)
}

// GetAllBook will return all the books
func GetAllBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	// get all the books in the db
	books, err := getAllBooks()

	if err != nil {
		log.Fatalf("Unable to get all book. %v", err)
	}

	// send all the books as response
	json.NewEncoder(w).Encode(books)
}

// UpdateBook update book's detail in the postgres db
func UpdateBook(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "PUT")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	// get the bookid from the request params, key is "id"
	params := mux.Vars(r)

	// convert the id type from string to int
	id, err := strconv.Atoi(params["id"])

	if err != nil {
		log.Fatalf("Unable to convert the string into int.  %v", err)
	}

	// create an empty book of type models.Book
	var book models.Book

	// decode the json request to book
	err = json.NewDecoder(r.Body).Decode(&book)

	if err != nil {
		log.Fatalf("Unable to decode the request body.  %v", err)
	}

	// call update book to update the book
	updatedRows := updateBook(int64(id), book)

	// format the message string
	msg := fmt.Sprintf("Book updated successfully. Total rows/record affected %v", updatedRows)

	// format the response message
	res := response{
		ID:      int64(id),
		Message: msg,
	}

	// send the response
	json.NewEncoder(w).Encode(res)
}

// DeleteBook delete book's detail in the postgres db
func DeleteBook(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	// get the bookid from the request params, key is "id"
	params := mux.Vars(r)

	// convert the id in string to int
	id, err := strconv.Atoi(params["id"])

	if err != nil {
		log.Fatalf("Unable to convert the string into int.  %v", err)
	}

	// call the deleteBook, convert the int to int64
	deletedRows := deleteBook(int64(id))

	// format the message string
	msg := fmt.Sprintf("Book deleted successfully. Total rows/record affected %v", deletedRows)

	// format the reponse message
	res := response{
		ID:      int64(id),
		Message: msg,
	}

	// send the response
	json.NewEncoder(w).Encode(res)
}

//------------------------- handler functions ----------------
// insert one book in the DB
func insertBook(book models.Book) int64 {

	// create the postgres db connection
	db := createConnection()

	// close the db connection
	defer db.Close()

	// create the insert sql query
	// returning bookid will return the id of the inserted book
	sqlStatement := `INSERT INTO books (name, author, published_at) VALUES ($1, $2, $3) RETURNING bookid`

	// the inserted id will store in this id
	var id int64

	// execute the sql statement
	// Scan function will save the insert id in the id
	err := db.QueryRow(sqlStatement, book.Name, book.Author, book.PublishedAt).Scan(&id)

	if err != nil {
		log.Fatalf("Unable to execute the query. %v", err)
	}

	fmt.Printf("Inserted a single record %v", id)

	// return the inserted id
	return id
}

// get one book from the DB by its bookid
func getBook(id int64) (models.Book, error) {
	// create the postgres db connection
	db := createConnection()

	// close the db connection
	defer db.Close()

	// create a book of models.Book type
	var book models.Book

	// create the select sql query
	sqlStatement := `SELECT * FROM books WHERE bookid=$1`

	// execute the sql statement
	row := db.QueryRow(sqlStatement, id)

	// unmarshal the row object to book
	err := row.Scan(&book.ID, &book.Name, &book.Author, &book.PublishedAt)

	switch err {
	case sql.ErrNoRows:
		fmt.Println("No rows were returned!")
		return book, nil
	case nil:
		return book, nil
	default:
		log.Fatalf("Unable to scan the row. %v", err)
	}

	// return empty book on error
	return book, err
}

// get one book from the DB by its bookid
func getAllBooks() ([]models.Book, error) {
	// create the postgres db connection
	db := createConnection()

	// close the db connection
	defer db.Close()

	var books []models.Book

	// create the select sql query
	sqlStatement := `SELECT * FROM books`

	// execute the sql statement
	rows, err := db.Query(sqlStatement)

	if err != nil {
		log.Fatalf("Unable to execute the query. %v", err)
	}

	// close the statement
	defer rows.Close()

	// iterate over the rows
	for rows.Next() {
		var book models.Book

		// unmarshal the row object to book
		err = rows.Scan(&book.ID, &book.Name, &book.Author, &book.PublishedAt)

		if err != nil {
			log.Fatalf("Unable to scan the row. %v", err)
		}

		// append the book in the books slice
		books = append(books, book)

	}

	// return empty book on error
	return books, err
}

// update book in the DB
func updateBook(id int64, book models.Book) int64 {

	// create the postgres db connection
	db := createConnection()

	// close the db connection
	defer db.Close()

	// create the update sql query
	sqlStatement := `UPDATE books SET name=$2, author=$3, published_at=$4 WHERE bookid=$1`

	// execute the sql statement
	res, err := db.Exec(sqlStatement, id, book.Name, book.Author, book.PublishedAt)

	if err != nil {
		log.Fatalf("Unable to execute the query. %v", err)
	}

	// check how many rows affected
	rowsAffected, err := res.RowsAffected()

	if err != nil {
		log.Fatalf("Error while checking the affected rows. %v", err)
	}

	fmt.Printf("Total rows/record affected %v", rowsAffected)

	return rowsAffected
}

// delete book in the DB
func deleteBook(id int64) int64 {

	// create the postgres db connection
	db := createConnection()

	// close the db connection
	defer db.Close()

	// create the delete sql query
	sqlStatement := `DELETE FROM books WHERE bookid=$1`

	// execute the sql statement
	res, err := db.Exec(sqlStatement, id)

	if err != nil {
		log.Fatalf("Unable to execute the query. %v", err)
	}

	// check how many rows affected
	rowsAffected, err := res.RowsAffected()

	if err != nil {
		log.Fatalf("Error while checking the affected rows. %v", err)
	}

	fmt.Printf("Total rows/record affected %v", rowsAffected)

	return rowsAffected
}