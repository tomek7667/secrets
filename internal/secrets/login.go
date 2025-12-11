package secrets

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/tomek7667/go-http-helpers/h"
)

type GetLoginDto struct {
	Username     string `json:"username"`
	Password     string `json:"password"`
	CaptchaToken string `json:"captchaToken"`
}

func (s *Server) PostLogin() {
	s.Router.With(withRateLimit(s.loginLimiter)).Post("/login", func(w http.ResponseWriter, r *http.Request) {
		dto, err := h.GetDto[GetLoginDto](r)
		if err != nil {
			h.ResBadRequest(w, err)
			return
		}
		err = s.verifyCaptcha(r.Context(), dto.CaptchaToken)
		if err != nil {
			h.ResBadRequest(w, err)
			return
		}

		user, err := s.Db.Queries.GetUserByUsername(r.Context(), dto.Username)
		if err != nil || user.Password != dto.Password {
			slog.Warn(
				"s.Dber.GetUserByUsername(r.Context(), dto.Username)",
				"err", err,
				"given pass", dto.Password,
			)
			s.Log(LoginFailedEvent, fmt.Sprintf("GetUserByUsername for %s with password %s", dto.Username, dto.Password), r)
			h.ResErr(w, fmt.Errorf("invalid username or password"))
			return
		}
		token, err := s.auther.GetToken(&user)
		if err != nil {
			s.Log(LoginFailedEvent, fmt.Sprintf("GetToken for %s failed: %s", dto.Username, err.Error()), r)
			h.ResErr(w, err)
			return
		}
		s.Log(LoginSuccessEvent, fmt.Sprintf("%s logged in", user.Username), r)
		h.ResSuccess(w, map[string]string{
			"token": token,
		})
	})
}
