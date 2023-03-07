package app

import (
	"errors"
	"fmt"
	"log"
	"net/http"
)

var ErrServer = errors.New("server fallen")

func (s *ServiceServer) Run() error {
	mux := http.NewServeMux()
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./templates/static"))))
	mux.HandleFunc("/", s.IndexWithSession)
	mux.HandleFunc("/signIn", s.SignIn)
	mux.HandleFunc("/signUp", s.SignUp)
	mux.HandleFunc("/signOut", s.authMiddleware(s.SignOut))
	mux.HandleFunc("/post", s.Post)
	mux.HandleFunc("/postLike", s.authMiddleware(s.PostLike))
	mux.HandleFunc("/postDislike", s.authMiddleware(s.PostDislike))
	mux.HandleFunc("/newPost", s.authMiddleware(s.NewPost))
	mux.HandleFunc("/createComment", s.authMiddleware(s.CreateComment))
	mux.HandleFunc("/commentLike", s.authMiddleware(s.CommentLike))
	mux.HandleFunc("/commentDislike", s.authMiddleware(s.CommentDislike))
	mux.HandleFunc("/category", s.Category)
	mux.HandleFunc("/createdPosts", s.authMiddleware(s.CreatedPosts))
	mux.HandleFunc("/ratedPosts", s.authMiddleware(s.RatedPosts))

	server := http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
	log.Printf("server started at http://localhost%s", server.Addr)

	err := server.ListenAndServe()
	if err != nil {
		return fmt.Errorf("server.listenAndServe: %w", err)
	}
	return nil
}
