package main

import (
	"io"
	"net/http"
)

const form = `<html>
    <head>
    <title></title>
    </head>
    <body>
        <form action="/" method="post">
            <label>Логин <input type="text" name="login"></label>
            <label>Пароль <input type="password" name="password"></label>
            <input type="submit" value="Login">
        </form>
    </body>
</html>`

func Auth(login, password string) bool {
	return login == `guest` && password == `demo`
}

func mainPage(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		login := r.FormValue("login")
		password := r.FormValue("password")
		if Auth(login, password) {
			_, err := io.WriteString(w, "Добро пожаловать!")
			if err != nil {
				return
			}
		} else {
			http.Error(w, "Неверный логин или пароль", http.StatusUnauthorized)
		}
		return
	} else {
		_, err := w.Write([]byte(form))
		if err != nil {
			return
		}
	}
}

func main() {
	err := http.ListenAndServe(`:8080`, http.HandlerFunc(mainPage))
	if err != nil {
		panic(err)
	}
}
