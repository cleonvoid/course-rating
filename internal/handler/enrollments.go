package handler

import (
	"net/http"
	"strconv"

	"github.com/cleonvoid/course-rating/internal/db"
	"github.com/cleonvoid/course-rating/internal/session"
)

type EnrollmentHandler struct {
	queries       *db.Queries
	courses       *CourseHandler
	sessionSecret string
}

func NewEnrollmentHandler(queries *db.Queries, courses *CourseHandler, sessionSecret string) *EnrollmentHandler {
	return &EnrollmentHandler{queries: queries, courses: courses, sessionSecret: sessionSecret}
}

func (h *EnrollmentHandler) Enroll(w http.ResponseWriter, r *http.Request) {
	courseID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "invalid course id", http.StatusBadRequest)
		return
	}

	userID, ok := session.Get(r, h.sessionSecret)
	if !ok {
		http.Error(w, "not signed in", http.StatusUnauthorized)
		return
	}

	if err := h.queries.EnrollUser(r.Context(), db.EnrollUserParams{
		UserID:   userID,
		CourseID: int32(courseID),
	}); err != nil {
		http.Error(w, "failed to enroll", http.StatusInternalServerError)
		return
	}

	course, err := h.queries.GetCourse(r.Context(), int32(courseID))
	if err != nil {
		http.Error(w, "failed to fetch course", http.StatusInternalServerError)
		return
	}

	currentUser, err := h.queries.GetUserByID(r.Context(), userID)
	if err != nil {
		http.Error(w, "failed to fetch user", http.StatusInternalServerError)
		return
	}

	data := map[string]any{
		"CourseID":      course.ID,
		"AvgRating":     course.AvgRating,
		"RatingCount":   course.RatingCount,
		"CurrentUser":   &currentUser,
		"IsEnrolled":    true,
		"StudentRating": 0,
	}

	if err := h.courses.tmpls.ExecuteTemplate(w, "enroll_section.html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
