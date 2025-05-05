package server

import (
	"log/slog"
	"net/http"
	"os"

	"byteXlearn/cmd/web"

	"github.com/a-h/templ"
	"github.com/google/uuid"
)

// createSessionCookie function is a helper function that creates a session cookie
func createSessionCookie(sessionID uuid.UUID) http.Cookie {
	// Determine the 'Secure' flag based on the environment.
	secureFlag := os.Getenv("ENV") == "production"
	cookie := http.Cookie{
		Name:     "sessionid",
		Value:    sessionID.String(),
		Path:     "/",
		MaxAge:   3600 * 24 * 7 * 2, // 2 weeks
		HttpOnly: true,
		Secure:   secureFlag,
		SameSite: http.SameSiteStrictMode,
	}

	return cookie
}

// renderDashboardComponent renders a component either as a full dashboard page
// (when not an HTMX request) or just the component (when it's an HTMX request).
func (s *Server) renderComponent(w http.ResponseWriter, r *http.Request, children templ.Component) {
	if r.Header.Get("HX-Request") == "true" {
		if err := children.Render(r.Context(), w); err != nil {
			writeError(w, http.StatusBadRequest, err.Error())
			slog.Error("Failed to render dashboard component", "error", err)
		}
	} else {
		userRole, ok := r.Context().Value(userContextKey).(User)
		if !ok {
			writeError(w, http.StatusUnauthorized, "unauthorised")
			return
		}
		user := web.DashboardUserRole{
			Role: userRole.Role,
		}

		rooms, err := s.queries.ListPublicRooms(r.Context())
		if err != nil {
			writeError(w, http.StatusInternalServerError, "no rooms found")
		}

		ctx := templ.WithChildren(r.Context(), children)
		if err := web.Dashboard(user, rooms).Render(ctx, w); err != nil {
			writeError(w, http.StatusBadRequest, err.Error())
			slog.Error("Failed to render dashboard layout", "error", err)
		}
	}
}
