package model

type PostReaction struct {
	ID      int64
	Post    Post
	User    User
	Like    int
	Dislike int
}

type CommentReaction struct {
	ID      int64
	Comment Comment
	User    User
	Like    int
	Dislike int
}
