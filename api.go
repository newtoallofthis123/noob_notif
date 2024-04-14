package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/websocket"
)

type Server struct {
	subs map[*websocket.Conn]Subscriber
	db   *sql.DB
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func NewServer() *Server {
	return &Server{
		subs: make(map[*websocket.Conn]Subscriber),
		db:   NewDbInstance(),
	}
}

func (s *Server) StartServer(port string) {
	http.HandleFunc("/subscriber", s.handleSubscribe)
	http.HandleFunc("/notif", s.handleNotification)
	http.HandleFunc("/broadcast", s.handleBroadCast)
	http.HandleFunc("/", handleFrontend)

	err := http.ListenAndServe(port, nil)
	if err != nil {
		panic(err)
	}
}

func (s *Server) handleSubscribe(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		err := r.ParseForm()
		if err != nil {
			fmt.Fprintln(w, "form parse err: ", err)
			return
		}

		name := r.Form.Get("name")
		if name == "" {
			fmt.Fprintln(w, "no attr name found")
			return
		}

		sub, err := s.InsertSub(name)
		if err != nil {
			fmt.Fprintln(w, "Error inserting: ", err)
		}

		fmt.Fprintf(w, "Subscribed %s: Password - %s", name, sub.password)
	} else {
		fmt.Fprintf(w, "Only POST route :)")
	}
}

func (s *Server) handleNotification(w http.ResponseWriter, r *http.Request) {
	auth_header := r.Header.Get("Authorization")
	if auth_header == "" {
		auth_cookie, err := r.Cookie("Authorization")
		if err != nil {
			fmt.Fprintln(w, "Unable to parse header nor cookie")
			return
		}

		auth_header = auth_cookie.Value
	}
	parsed := strings.Split(auth_header, "|")
	id := parsed[0]
	password := parsed[1]
	if id == "" || password == "" {
		fmt.Fprintln(w, "No id and password found")
		fmt.Println("No id and password found")
		return
	}

	sub, err := s.GetSub(id)
	if err != nil {
		fmt.Fprintln(w, "Getting sub of id failed: ", err)
		fmt.Println(err)
		return
	}
	if sub.password != password {
		fmt.Println(err)
		fmt.Fprintln(w, "Password are not matching!")
		return
	}

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Fprintln(w, "Unable to upgrade the connection: ", err)
		fmt.Println(err)
		return
	}

	s.subs[ws] = sub
	ws.WriteMessage(websocket.TextMessage, []byte("Connected to the notification server"))
	fmt.Printf("%s with addr %s connected!\n", sub.name, ws.RemoteAddr())
}

func (s *Server) handleBroadCast(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		err := r.ParseForm()
		if err != nil {
			fmt.Fprintln(w, "Unable to parse form: ", err)
			return
		}

		msg := r.Form.Get("msg")
		if msg == "" {
			fmt.Fprintln(w, "msg attr not found")
			return
		}
		count := 0

		for conn, sub := range s.subs {
			if sub.valid {
				conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("%s! %s", sub.name, msg)))
				count++
			}
		}
		fmt.Fprintf(w, "Broadcasted Notification to %d clients!\n", count)
	} else {
		fmt.Fprintln(w, "Post only :)")
	}
}

func handleFrontend(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		fmt.Fprintln(w, "Get only :)")
	}
	content, err := os.ReadFile("frontend/index.html")
	if err != nil {
		fmt.Fprintln(w, "Error reading index.html to send to browser")
	}

	fmt.Fprintln(w, string(content))
}
