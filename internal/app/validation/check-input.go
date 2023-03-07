package validation

import (
	"fmt"
	"net/http"
	"strings"

	"forum/internal/model"
)

func CheckInput(r *http.Request, post *model.Post, allCategories []model.Category) error {
	if err := r.ParseForm(); err != nil {
		return fmt.Errorf("checkInput: %w", err)
	}

	post.Title = r.PostFormValue("title")

	if temp := strings.Trim(post.Title, " "); temp == "" || len(post.Title) > 50 {
		return fmt.Errorf("checkInput: %w", model.ErrMessageInvalid)
	}

	post.Content = r.PostFormValue("content")
	if temp := strings.Trim(post.Content, " "); temp == "" {
		return fmt.Errorf("checkInput: %w", model.ErrMessageInvalid)
	}

	categories := r.Form["categories"]
	if len(categories) == 0 {
		return fmt.Errorf("checkInput: %w", model.ErrMessageInvalid)
	}

	for i := 0; i < len(categories); i++ {
		status := false
		for j := 0; j < len(allCategories); j++ {
			if categories[i] == allCategories[j].Category {
				post.Categories = append(post.Categories, allCategories[j])
				status = true
				break
			}
		}

		if !status {
			return fmt.Errorf("checkInput: %w", model.ErrMessageInvalid)
		}
	}

	return nil
}
