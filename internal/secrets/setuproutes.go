package secrets

func (s *Server) SetupRoutes() {
	s.PostLogin()
	s.AddUsersRoutes()
	s.AddSecretsRoutes()
	s.AddTokensRoutes()
}
