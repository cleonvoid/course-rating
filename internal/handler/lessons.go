package handler

import (
	"html/template"
	"net/http"
	"strconv"
	"strings"

	"github.com/cleonvoid/course-rating/internal/db"
	"github.com/cleonvoid/course-rating/internal/session"
)

type LessonHandler struct {
	queries       *db.Queries
	tmpls         *template.Template
	sessionSecret string
}

func NewLessonHandler(queries *db.Queries, tmpls *template.Template, sessionSecret string) *LessonHandler {
	return &LessonHandler{queries: queries, tmpls: tmpls, sessionSecret: sessionSecret}
}

func (h *LessonHandler) Detail(w http.ResponseWriter, r *http.Request) {
	lessonID, err := strconv.Atoi(r.PathValue("lid"))
	if err != nil {
		http.Error(w, "invalid lesson id", http.StatusBadRequest)
		return
	}

	lesson, err := h.queries.GetLesson(r.Context(), int32(lessonID))
	if err != nil {
		http.Error(w, "lesson not found", http.StatusNotFound)
		return
	}

	comments, err := h.queries.ListCommentsByLesson(r.Context(), lesson.ID)
	if err != nil {
		http.Error(w, "failed to fetch comments", http.StatusInternalServerError)
		return
	}

	data := struct {
		Lesson      db.Lesson
		Comments    []db.ListCommentsByLessonRow
		CurrentUser *db.GetUserByIDRow
	}{
		Lesson:      lesson,
		Comments:    comments,
		CurrentUser: h.currentUser(r),
	}

	if err := h.tmpls.ExecuteTemplate(w, "lesson.html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *LessonHandler) CreateComment(w http.ResponseWriter, r *http.Request) {
	courseID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "invalid course id", http.StatusBadRequest)
		return
	}

	lessonID, err := strconv.Atoi(r.PathValue("lid"))
	if err != nil {
		http.Error(w, "invalid lesson id", http.StatusBadRequest)
		return
	}

	userID, ok := session.Get(r, h.sessionSecret)
	if !ok {
		http.Error(w, "not signed in", http.StatusUnauthorized)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "invalid form", http.StatusBadRequest)
		return
	}

	body := strings.TrimSpace(r.FormValue("body"))
	if body == "" {
		http.Error(w, "comment cannot be empty", http.StatusBadRequest)
		return
	}

	if err := h.queries.CreateComment(r.Context(), db.CreateCommentParams{
		LessonID: int32(lessonID),
		UserID:   userID,
		Body:     body,
	}); err != nil {
		http.Error(w, "failed to save comment", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/courses/"+strconv.Itoa(courseID)+"/lessons/"+strconv.Itoa(lessonID), http.StatusSeeOther)
}

func (h *LessonHandler) DeleteComment(w http.ResponseWriter, r *http.Request) {
	courseID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "invalid course id", http.StatusBadRequest)
		return
	}

	lessonID, err := strconv.Atoi(r.PathValue("lid"))
	if err != nil {
		http.Error(w, "invalid lesson id", http.StatusBadRequest)
		return
	}

	commentID, err := strconv.Atoi(r.PathValue("cid"))
	if err != nil {
		http.Error(w, "invalid comment id", http.StatusBadRequest)
		return
	}

	userID, ok := session.Get(r, h.sessionSecret)
	if !ok {
		http.Error(w, "not signed in", http.StatusUnauthorized)
		return
	}

	if err := h.queries.DeleteComment(r.Context(), db.DeleteCommentParams{
		ID:     int32(commentID),
		UserID: userID,
	}); err != nil {
		http.Error(w, "failed to delete comment", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/courses/"+strconv.Itoa(courseID)+"/lessons/"+strconv.Itoa(lessonID), http.StatusSeeOther)
}

func (h *LessonHandler) currentUser(r *http.Request) *db.GetUserByIDRow {
	userID, ok := session.Get(r, h.sessionSecret)
	if !ok {
		return nil
	}
	user, err := h.queries.GetUserByID(r.Context(), userID)
	if err != nil {
		return nil
	}
	return &user
}
