package model

type Data struct {
	Status     bool
	Categories []Category
	Posts      []Post
}

type DataPost struct {
	Session bool
	Post
}
