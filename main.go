package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"sync"
)

type templateHandler struct {
	once     sync.Once
	template *template.Template
	filename string
}

type authHandler struct {
	next http.Handler
}

var (
	host = flag.String("host", "0.0.0.0", "Hostname")
	port = flag.String("port", ":80", "Port")
)

func (h *authHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	_, err := req.Cookie("auth")
	if err == http.ErrNoCookie {
		w.Header().Set("Location", "/login")
		w.WriteHeader(http.StatusTemporaryRedirect)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	h.next.ServeHTTP(w, req)
}

// MustAuth create handler with auth
func MustAuth(handler http.Handler) http.Handler {
	return &authHandler{next: handler}
}

func (t *templateHandler) ServeHTTP(writer http.ResponseWriter, req *http.Request) {
	t.once.Do(func() {
		t.template = template.Must(template.ParseFiles(filepath.Join("public", t.filename)))
	})
	t.template.Execute(writer, req)
}

func main() {
	flag.Parse()
	// gomniauth.SetSecurityKey("lequangphu")
	// gomniauth.WithProviders(
	// 	facebook.New(
	// 		"1737355836575055",
	// 		"10af2ef771ecd35a309794eaeb40009b",
	// 		"http://localhost:2711/auth/callback/facebook"
	// 	),
	// 	github.New(
	// 		"d2fbdfc461257bd1a8de",
	// 		"7ed62279e9e0a94c0b73718bf5b5ab850e5716cb",
	// 		"http://localhost:2711/auth/callback/github"
	// 	),
	// 	google.New(
	// 		"866874494879-35970trihs9fsuu3c3lsj0tqarm6ntdo.apps.googleusercontent.com",
	// 		"wTJ_oQSUN-znU2JBiQ6d2iae",
	// 		"http://localhost:2711/auth/callback/google"
	// 	),
	// )
	r := newRoom()
	// http.Handle("/assets/", http.StripPrefix("/assets", http.FileServer(http.Dir("public/assets"))))
	http.Handle("/chat", MustAuth(&templateHandler{filename: "chat.html"}))
	http.Handle("/login", &templateHandler{filename: "login.html"})
	http.HandleFunc("/auth/", loginHandler)
	http.Handle("/room", r)
	go r.run()

	log.Printf("Starting web server on http://%s:%s\n", *host, *port)
	if err := http.ListenAndServe(fmt.Sprintf("%s:%s", *host, *port), nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
