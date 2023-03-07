package app

import (
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"

	"forum/internal/model"
)

func (s *ServiceServer) CreateComment(w http.ResponseWriter, r *http.Request, session model.Session) {
	if r.Method != http.MethodPost {
		s.ErrorHandler(w, model.NewErrorWeb(http.StatusMethodNotAllowed))
		return
	}

	postID, err := s.getID(r)
	if err != nil {
		log.Println("ERROR:CreateComment:", err)

		if errors.Is(err, model.ErrValueNotSet) {
			s.ErrorHandler(w, model.NewErrorWeb(http.StatusBadRequest))
			return
		}

		s.ErrorHandler(w, model.NewErrorWeb(http.StatusInternalServerError))
		return
	}

	if r.ParseForm() != nil {
		log.Println("ERROR:CreateComment:", err)

		s.ErrorHandler(w, model.NewErrorWeb(http.StatusBadGateway))
		return
	}

	comment := model.Comment{PostID: int64(postID), UserID: session.User.ID, Username: session.User.Username, Message: r.PostFormValue("comment")}
	temp := strings.Trim(comment.Message, " ")
	if temp == "" {
		log.Println("ERROR:\nCreateComment:", model.ErrMessageInvalid)

		s.ErrorHandler(w, model.NewErrorWeb(http.StatusBadRequest))
		return
	}

	err = s.postService.CreateComment(&comment)
	if err != nil {
		log.Println("ERROR:\nCreateComment:", err)
		if errors.Is(err, model.ErrPostNotFound) {
			s.ErrorHandler(w, model.NewErrorWeb(http.StatusBadRequest))
			return
		}

		s.ErrorHandler(w, model.NewErrorWeb(http.StatusInternalServerError))
		return
	}

	http.Redirect(w, r, "/post?ID="+strconv.Itoa(postID), http.StatusFound)
}
