package main_test

import (
	"os"
	"testing"

	"."
	"log"
	"net/http"
	"net/http/httptest"
	"strconv"
)

var a main.App
// to login psql:
// psql -d testwechat -U iris -W
// password: wMMvbj35
// to start psql server:
// lunchy start postgres
// to stop psql server:
// lunchy stop postgres
// https://www.moncefbelyamani.com/how-to-install-postgresql-on-a-mac-with-homebrew-and-lunchy/

func TestMain(m *testing.M) {
	a = main.App{}

	os.Setenv("TEST_DB_USER_NAME", "iris")
	os.Setenv("TEST_DB_PASSWORD", "wMMvbj35")
	os.Setenv("TEST_DB_NAME", "testwechat")

	a.Initialize(
		os.Getenv("TEST_DB_USERNAME"),
		os.Getenv("TEST_DB_PASSWORD"),
		os.Getenv("TEST_DB_NAME"))

	ensureTableExists()

	code := m.Run()

	//clearTable()

	os.Exit(code)
}

func ensureTableExists() {
	if _, err := a.DB.Exec(tableCreationQuery); err != nil {
		log.Fatal(err)
	}
}

func TestGetProducts(t *testing.T) {
	clearTable()
	addUsers(5)

	req, _ := http.NewRequest("GET", "/users", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)
}

func TestEmptyUsers(t *testing.T) {
	clearTable()

	req, _ := http.NewRequest("GET", "/users", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	if body := response.Body.String(); body != "[]" {
		t.Errorf("Expected an empty array. Got %s", body)
	}
}

func addUsers(count int) {
	if count < 1 {
		count = 1
	}

	for i := 0; i < count; i++ {
		a.DB.Exec("INSERT INTO users(userName, password) VALUES($1, $2)", "User" + strconv.Itoa(i), i)
	}
}

// helpers
func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	a.Router.ServeHTTP(rr, req)

	return rr
}


func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d, Got %d\n", expected, actual)
	}
}

const tableCreationQuery = `CREATE TABLE IF NOT EXISTS users
(
id SERIAL,
userName TEXT NOT NULL,
password TEXT NOT NULL,
CONSTRAINT user_pkey PRIMARY KEY (id)
)`

func clearTable() {
	a.DB.Exec("DELETE FROM users")
	a.DB.Exec("ALTER SEQUENCE users_id_seq RESETART WITH 1")
}


