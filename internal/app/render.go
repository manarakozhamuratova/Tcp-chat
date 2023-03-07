package app

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"forum/internal/model"
)

func (s *ServiceServer) render(w http.ResponseWriter, pageName string, data interface{}) {
	path := fmt.Sprintf("./templates/html/%s.html", pageName)
	t, err := template.ParseFiles(path)
	if err != nil {
		log.Println("ERROR:\nrender:", err)

		s.ErrorHandler(w, model.NewErrorWeb(http.StatusInternalServerError))
		return
	}

	if t.ExecuteTemplate(w, pageName, data) != nil {
		log.Println("ERROR:\nrender:", err)

		s.ErrorHandler(w, model.NewErrorWeb(http.StatusInternalServerError))
	}
}
