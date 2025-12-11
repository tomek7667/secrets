package secrets

func (s *Server) SetupRoutes() {
	s.AddFrontendRoutes()
	s.PostLogin()
	s.AddUsersRoutes()
	s.AddSecretsRoutes()
	s.AddTokensRoutes()
	s.AddPermissionsRoutes()
	s.AddCertificatesRoutes()
}
