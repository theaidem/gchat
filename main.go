package main

import (
	"html/template"
	"log"
	"net/http"
	"net/url"
	"path"

	"github.com/gorilla/sessions"
	"github.com/trevex/golem"
)

const (
	secret      = "secret_key"
	sessionName = "sid"
)

var store = sessions.NewCookieStore([]byte(secret))
var messages []*Message
var users map[string]string
var room = golem.NewRoom()

type User struct {
	Name string
	Conn *golem.Connection
}

func init() {
	users = make(map[string]string, 0)
}

func main() {
	wsrouter := golem.NewRouter()
	wsrouter.SetConnectionExtension(NewUser)

	wsrouter.On("msg", msg)
	wsrouter.OnClose(onClose)
	wsrouter.OnConnect(onConnect)

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	http.HandleFunc("/", index)
	http.HandleFunc("/login", login)
	http.HandleFunc("/logout", logout)

	http.HandleFunc("/ws", wsrouter.Handler())

	// Listen Server
	if err := http.ListenAndServe(":3333", nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}

func NewUser(conn *golem.Connection) *User {
	return &User{Conn: conn}
}

func index(w http.ResponseWriter, r *http.Request) {

	if session, err := store.Get(r, sessionName); err == nil {
		if session.Values["isAuth"] != nil {
			uName := session.Values["uName"]

			fp := path.Join("templates", "index.tpl")
			tmpl, err := template.ParseFiles(fp)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			if err := tmpl.Execute(w, uName); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}

		} else {
			http.Redirect(w, r, "/login", http.StatusFound)
		}
	}
}

func login(w http.ResponseWriter, r *http.Request) {

	session, _ := store.Get(r, sessionName)

	switch r.Method {
	case "GET":

		if session.Values["isAuth"] != nil {
			//TODO: More check
			http.Redirect(w, r, "/", http.StatusFound)
		}

		fp := path.Join("templates", "login.tpl")
		tmpl, err := template.ParseFiles(fp)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if err := tmpl.Execute(w, nil); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

	case "POST":
		uName := r.PostFormValue("uName")
		if len(uName) == 0 {
			http.Redirect(w, r, "/login", http.StatusFound)
		}
		session.Values["isAuth"] = true
		session.Values["uName"] = uName
		session.Save(r, w)
		http.Redirect(w, r, "/", http.StatusFound)
	}
}

func logout(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		if session, err := store.Get(r, sessionName); err == nil {
			delete(session.Values, "isAuth")
			session.Save(r, w)
		}
		http.Redirect(w, r, "/", http.StatusFound)
	}
}

type ConnectData struct {
	Name     string     `json:"name"`
	Messages []*Message `json:"messages"`
	Users    []string   `json:"users"`
}

func onConnect(conn *User, r *http.Request) {
	m, err := url.ParseQuery(r.RequestURI)
	if err != nil {
		log.Fatal(err)
	}
	conn.Name = m["/ws?uname"][0]
	users[conn.Name] = conn.Name

	var room_users []string
	for name, _ := range users {
		room_users = append(room_users, name)
	}

	var room_messages []*Message
	if len(messages) > 50 {
		room_messages = messages[:50]
	} else {
		room_messages = messages
	}

	room.Join(conn.Conn)
	room.Emit("join", &ConnectData{conn.Name, room_messages, room_users})
}

func onClose(conn *User) {
	// Potential panic
	delete(users, conn.Name)
	room.Leave(conn.Conn)
	room.Emit("leave", conn.Name)
}

type Message struct {
	UName string `json:"uname"`
	Msg   string `json:"msg"`
}

func msg(conn *User, data *Message) {
	messages = append(messages, data)
	room.Emit("message", data)
	log.Println(messages)
}
