/*package main

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

type Weather struct {
	City string `json:"city"`
	Forecast string `json:"forecast"`
}

type Tests struct {
	name string
	server *httptest.Server
	response *Weather
	expectedError error
}

func TestGetWeather(t *testing.T) {
	tests := []Tests{
		{
			name: "basic-request",
			server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"city": "Denver, CO", "forecast": "sunny"}`))//creating mock response data
			})),
			response: &Weather{
				City: "Denver, CO",
				Forecast: "sunny",
			},
			expectedError: nil,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T){
			defer test.server.Close()

			resp, err :=GetWeather(test.server.URL)

			if !reflect.DeepEqual(resp, test.response){
				t.Errorf("FAILED: expected %v, got %v\n", test.response, resp)
			}
			if !errors.Is(err, test.expectedError){
				t.Errorf("Expected error FAILED: expected %v got %v",test.expectedError, err)
			}
		} )
	}
}*/

package main

import (
	"Gorilla_Websocket_Exercise/database"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

//Test 1.: correct method and route
func TestServeHomeGET(t *testing.T) {
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ServeHome)

	handler.ServeHTTP(rr, req)

	// Assert the response status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("ServeHome FAIL: unexpected status code: got %v, want %v", status, http.StatusOK)
	}

	// Assert the response body or any other expected behavior
	// ...
}

//Test 2.: incorrect method, correct route
func TestServeHomePOST(t *testing.T) {
	req, err := http.NewRequest("POST", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ServeHome)

	//here we apply method and route from line 87 to ServeHome
	handler.ServeHTTP(rr, req)

	// Assert the response status code
	if status := rr.Code; status == http.StatusOK {
		t.Errorf("ServeHome FAIL: unexpected status code: got %v, want %v", status, http.StatusOK)
	}

	// Assert the response body or any other expected behavior
	// ...
}

//Test 3.: correct method, incorrect route
func TestServeHomeROUTE(t *testing.T) {
	req, err := http.NewRequest("GET", "/Helena", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ServeHome)

	//here we apply method and route from line 87 to ServeHome
	handler.ServeHTTP(rr, req)

	// Assert the response status code
	if status := rr.Code; status == http.StatusOK {
		t.Errorf("ServeHome FAIL: unexpected status code: got %v, want %v", status, http.StatusOK)
	}

	// Assert the response body or any other expected behavior
	// ...
}

func TestMain(t *testing.T) {
	// Write test cases to test specific parts of the main function
	// For example, check if the database connection is successful
	//connect to db
	db, err := database.ConnectDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	//retrieve message to db
	type dbCheck struct {
		id      string
		user    string
		message string
	}

	res := []dbCheck{}

	//row, err := db.Exec("SELECT sender, content FROM messages WHERE message=Lovely day")
	/*err = db.QueryRow("SELECT sender, content FROM messages WHERE message=Lovely day").Scan(&res[0].user, &res[0].message)
	if err != nil {
		fmt.Println("error retrieving user and message from db", err)
	}*/

	allMsgs := dbCheck{}

	rows, err := db.Query("SELECT * FROM messages")
	if err != nil {
		log.Fatalln("unable to retrieve db data", err)
	}

	// Loop through 'message' table rows using only one struct
	for rows.Next() {
		err := rows.Scan(&allMsgs.id, &allMsgs.user, &allMsgs.message)
		if err != nil {
			log.Fatalln(err)
		}
		res = append(res, allMsgs)
	}

	compRES := dbCheck{
		id:      "2",
		user:    "Helena",
		message: "Lovely day",
	}

	if !reflect.DeepEqual(compRES, res[1]) {
		t.Errorf("FAILED: expected %v, got %v\n", compRES, res[1])
	}

	// You can also use the httptest package to test the server's behavior
	// For example, you can create an HTTP request to test the server's response
	// ...
}
