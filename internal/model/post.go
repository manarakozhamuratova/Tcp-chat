package model

type Post struct {
	ID         int64
	Title      string
	Content    string
	User       User
	Comments   []Comment
	Categories []Category
	Like       int
	Dislike    int
}
