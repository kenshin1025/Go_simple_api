package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"simpleAPI/apierr"
)

type userTabel struct {
	id       int
	name     string
	email    string
	password string
}

type db struct {
	user []*userTabel
}

func newDB() *db {
	return &db{}
}

func main() {
	//apiの起動確認
	fmt.Printf("Starting server at 'http://localhost:8080'\n")

	// dbのインスタンスを作る
	db := newDB()

	http.HandleFunc("/", showUsers(db))
	http.HandleFunc("/user/create", createUser(db))
	http.ListenAndServe(":8080", nil)
}

//ちゃんとユーザーが入っているか確認する用
func showUsers(db *db) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		for _, user := range db.user {
			fmt.Fprintf(w, "%+v\n", user)
		}
	}
}

//createUserのリクエストのjsonに合わせた構造体の定義
type ReqCreateUserJSON struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

//createUserのレスポンスのjsonに合わせた構造体の定義
type ResCreateUserJSON struct {
	ID int `json:"id"`
}

func createUser(db *db) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//POSTリクエストの時だけ処理をする
		if r.Method == "POST" {
			//jsonからgoの構造体にデコードする
			var req ReqCreateUserJSON
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				log.Fatal(err)
				return
			}

			user := &userTabel{
				name:     req.Name,
				email:    req.Email,
				password: req.Password,
			}

			//送られてきたユーザーのEmailがすでに登録されていないかチェックする
			err := findByEmil(db, user.email)
			//すでに登録されていたらエラーを返す
			if errors.Is(err, apierr.ErrEmailAlreadyExists) {
				w.Header().Set("Content-Type", "application/json;charset=utf-8")
				w.WriteHeader(http.StatusBadRequest)
				return
			} else if err != nil {
				log.Fatal(err)
				return
			}

			err = create(db, user)
			if err != nil {
				log.Fatal(err)
				return
			}

			//レスポンス用にヘッダーをセットする
			w.Header().Set("Content-Type", "application/json;charset=utf-8")
			w.WriteHeader(http.StatusCreated)

			if err := json.NewEncoder(w).Encode(&ResCreateUserJSON{
				ID: user.id,
			}); err != nil {
				log.Fatal(err)
			}
		}
	}
}

//送られてきたユーザーのEmailがすでに登録されていないかチェックする
func findByEmil(db *db, email string) error {
	for _, user := range db.user {
		if user.email == email {
			return apierr.ErrEmailAlreadyExists
		}
	}
	return nil
}

//送られてきたユーザーの情報にidを割り振り、DBに登録する
func create(db *db, user *userTabel) error {
	user.id = len(db.user) + 1
	db.user = append(db.user, user)
	return nil
}
