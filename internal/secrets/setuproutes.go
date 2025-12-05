package secrets

func (s *Server) SetupRoutes() {
	s.PostLogin()
	s.AddSecretsRoutes()
}
