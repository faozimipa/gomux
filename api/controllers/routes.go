package controllers

import "github.com/faozimipa/gomux/api/middlewares"

func (s *Server) initializeRoutes() {

	// Home Route
	s.Router.HandleFunc("/", middlewares.SetMiddlewareJSON(s.Home)).Methods("GET")

	// Login Route
	s.Router.HandleFunc("/auth/login", middlewares.SetMiddlewareJSON(s.Login)).Methods("POST")
	s.Router.HandleFunc("/auth/logout", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(s.Logout))).Methods("POST")

	//Users routes
	s.Router.HandleFunc("/users", middlewares.SetMiddlewareJSON(s.CreateUser)).Methods("POST")
	s.Router.HandleFunc("/users", middlewares.SetMiddlewareJSON(s.GetUsers)).Methods("GET")
	s.Router.HandleFunc("/users/{id}", middlewares.SetMiddlewareJSON(s.GetUser)).Methods("GET")
	s.Router.HandleFunc("/users/{id}", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(s.UpdateUser))).Methods("PUT")
	s.Router.HandleFunc("/users/{id}", middlewares.SetMiddlewareAuthentication(s.DeleteUser)).Methods("DELETE")

	//Posts routes
	s.Router.HandleFunc("/posts", middlewares.SetMiddlewareJSON(s.CreatePost)).Methods("POST")
	s.Router.HandleFunc("/posts", middlewares.SetMiddlewareJSON(s.GetPosts)).Methods("GET")
	s.Router.HandleFunc("/posts/{id}", middlewares.SetMiddlewareJSON(s.GetPost)).Methods("GET")
	s.Router.HandleFunc("/posts/{id}", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(s.UpdatePost))).Methods("PUT")
	s.Router.HandleFunc("/posts/{id}", middlewares.SetMiddlewareAuthentication(s.DeletePost)).Methods("DELETE")

	//test redis
	s.Router.HandleFunc("/test", middlewares.SetMiddlewareJSON(s.TestRedis)).Methods("GET")
	s.Router.HandleFunc("/test-set", middlewares.SetMiddlewareJSON(s.SetData)).Methods("GET")
	s.Router.HandleFunc("/test-get", middlewares.SetMiddlewareJSON(s.GetData)).Methods("GET")
}
