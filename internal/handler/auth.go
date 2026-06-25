package handler

import (
	"html/template"
	"net/http"

	"golang.org/x/crypto/bcrypt"

	"github.com/cleonvoid/course-rating/internal/db"
	"github.com/cleonvoid/course-rating/internal/session"
)

type AuthHandler struct {
	queries       *db.Queries
	tmpls         *template.Template
	sessionSecret string
}

func NewAuthHandler(queries *db.Queries, tmpls *template.Template, sessionSecret string) *AuthHandler {
	return &AuthHandler{queries: queries, tmpls: tmpls, sessionSecret: sessionSecret}
}

func (h *AuthHandler) SignInPage(w http.ResponseWriter, r *http.Request) {
	if err := h.tmpls.ExecuteTemplate(w, "signin.html", map[string]any{"Error": ""}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *AuthHandler) SignIn(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "invalid form", http.StatusBadRequest)
		return
	}

	email := r.FormValue("email")
	password := r.FormValue("password")

	renderError := func(msg string) {
		h.tmpls.ExecuteTemplate(w, "signin.html", map[string]any{"Error": msg})
	}

	user, err := h.queries.GetUserByEmail(r.Context(), email)
	if err != nil {
		renderError("Invalid email or password.")
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		renderError("Invalid email or password.")
		return
	}

	session.Set(w, user.ID, h.sessionSecret)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h *AuthHandler) SignOut(w http.ResponseWriter, r *http.Request) {
	session.Clear(w)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
