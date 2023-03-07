package app

import (
	"html/template"
	"log"
	"net/http"

	"forum/internal/model"
)

func (s *ServiceServer) ErrorHandler(w http.ResponseWriter, errorStatus *model.ErrorWeb) {
	t, err := template.ParseFiles("./templates/html/error.html")
	if err != nil {
		http.Error(w, errorStatus.StatusText, errorStatus.StatusCode)
		return
	}

	w.WriteHeader(errorStatus.StatusCode)
	if t.ExecuteTemplate(w, "error", errorStatus) != nil {
		log.Println("ERROR:\nErrorHandler:", err)

		http.Error(w, errorStatus.StatusText, errorStatus.StatusCode)
	}
}
