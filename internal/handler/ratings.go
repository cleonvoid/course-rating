package handler

import (
	"net/http"
	"strconv"

	"github.com/cleonvoid/course-rating/internal/db"
	"github.com/cleonvoid/course-rating/internal/session"
)

type RatingHandler struct {
	queries       *db.Queries
	courses       *CourseHandler
	sessionSecret string
}

func NewRatingHandler(queries *db.Queries, courses *CourseHandler, sessionSecret string) *RatingHandler {
	return &RatingHandler{queries: queries, courses: courses, sessionSecret: sessionSecret}
}

func (h *RatingHandler) Upsert(w http.ResponseWriter, r *http.Request) {
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

	enrolled, err := h.queries.IsEnrolled(r.Context(), db.IsEnrolledParams{
		UserID:   userID,
		CourseID: int32(courseID),
	})
	if err != nil || !enrolled {
		http.Error(w, "not enrolled in this course", http.StatusForbidden)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "invalid form", http.StatusBadRequest)
		return
	}

	stars, err := strconv.Atoi(r.FormValue("stars"))
	if err != nil || stars < 1 || stars > 5 {
		http.Error(w, "stars must be between 1 and 5", http.StatusBadRequest)
		return
	}

	_, err = h.queries.UpsertRating(r.Context(), db.UpsertRatingParams{
		CourseID: int32(courseID),
		UserID:   userID,
		Stars:    int16(stars),
	})
	if err != nil {
		http.Error(w, "failed to save rating", http.StatusInternalServerError)
		return
	}

	course, err := h.queries.GetCourse(r.Context(), int32(courseID))
	if err != nil {
		http.Error(w, "failed to fetch course", http.StatusInternalServerError)
		return
	}

	data := map[string]any{
		"CourseID":      course.ID,
		"AvgRating":     course.AvgRating,
		"RatingCount":   course.RatingCount,
		"StudentRating": stars,
	}

	if err := h.courses.tmpls.ExecuteTemplate(w, "star_rating.html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
