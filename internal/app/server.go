package app

import "forum/internal/services"

type ServiceServer struct {
	authService    services.AuthService
	userService    services.UserService
	postService    services.PostService
	sessionService services.SessionService
}

func NewServiceServer(
	authService services.AuthService,
	userService services.UserService,
	postService services.PostService,
	sessionService services.SessionService,
) ServiceServer {
	return ServiceServer{
		authService:    authService,
		userService:    userService,
		postService:    postService,
		sessionService: sessionService,
	}
}
