package model

type Comment struct {
	ID       int64
	PostID   int64
	UserID   int64
	Username string
	Message  string
	Like     int
	Dislike  int
}
