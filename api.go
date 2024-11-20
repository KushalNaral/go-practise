package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

type PaginatedObject struct {
	Data []User
	Meta Pagination
}

type Pagination struct {
	Total        int
	PerPage      int
	Page         int
	CurrentTotal int
}

func Routes() {
	http.Handle("/", http.FileServer(http.Dir("./")))
	http.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			w.Header().Set("Content-Type", "application/json")
			users, err := GetAllUsers(r)
			if err != nil {
				http.Error(w, fmt.Sprintf("Err : %v", err), 500)
				return
			}
			fmt.Fprintf(w, users)
			break

		case "POST":
			fmt.Fprintf(w, "post method for /users")
			break

		case "DELETE":
			fmt.Fprintf(w, "delete method for /users")
			break
		}
	})

	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Printf("HTTP server failed: %v\n", err)
	}
}

func GetAllUsers(r *http.Request) (string, error) {

	per_page, _ := strconv.Atoi(r.URL.Query().Get("per_page"))
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	search := r.URL.Query().Get("search")

	fmt.Println("\nperpage : \n", per_page)

	if per_page == 0 {
		per_page = 10
	}

	if page == 0 {
		page = 1
	}

	limit := per_page * page

	db, err := OpenDB()
	if err != nil {
		return "", err
	}

	rows, err := db.Query(
		"SELECT * FROM users WHERE name LIKE ? ORDER BY id DESC LIMIT ?",
		"%"+search+"%", limit,
	)
	defer rows.Close()

	if err != nil {
		return "", err
	}

	var users []User

	count := 0
	for rows.Next() {
		count++
		var user User
		if err := rows.Scan(&user.Id, &user.Name, &user.Email, &user.Contact, &user.Income, &user.IsPresent, &user.JoinedAt); err != nil {
			return "", err
		}
		if page == 1 {
		}

		//page 2
		//per_page 10 limit = 20
		if count > (limit - per_page) {
			users = append(users, user)
		}

	}

	if err := rows.Err(); err != nil {
		return "", err
	}

	pObject := PaginatedObject{
		users, Pagination{
			Total:        count,
			PerPage:      per_page,
			Page:         page,
			CurrentTotal: len(users),
		},
	}

	usersJson, err := json.Marshal(pObject)

	return string(usersJson), nil

}
