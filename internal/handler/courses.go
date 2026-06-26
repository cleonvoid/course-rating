package handler

import (
	"html/template"
	"net/http"
	"strconv"

	"github.com/cleonvoid/course-rating/internal/db"
	"github.com/cleonvoid/course-rating/internal/session"
)

type CourseHandler struct {
	queries       *db.Queries
	tmpls         *template.Template
	sessionSecret string
}

func NewCourseHandler(queries *db.Queries, tmpls *template.Template, sessionSecret string) *CourseHandler {
	return &CourseHandler{queries: queries, tmpls: tmpls, sessionSecret: sessionSecret}
}

func (h *CourseHandler) List(w http.ResponseWriter, r *http.Request) {
	courses, err := h.queries.ListCourses(r.Context())
	if err != nil {
		http.Error(w, "failed to fetch courses", http.StatusInternalServerError)
		return
	}

	data := struct {
		Courses     []db.ListCoursesRow
		CurrentUser *db.GetUserByIDRow
	}{
		Courses:     courses,
		CurrentUser: h.currentUser(r),
	}

	if err := h.tmpls.ExecuteTemplate(w, "courses.html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *CourseHandler) Detail(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "invalid course id", http.StatusBadRequest)
		return
	}

	course, err := h.queries.GetCourse(r.Context(), int32(id))
	if err != nil {
		http.Error(w, "course not found", http.StatusNotFound)
		return
	}

	currentUser := h.currentUser(r)

	var isEnrolled bool
	var studentRating int

	if currentUser != nil {
		enrolled, err := h.queries.IsEnrolled(r.Context(), db.IsEnrolledParams{
			UserID:   currentUser.ID,
			CourseID: int32(id),
		})
		if err == nil {
			isEnrolled = enrolled
		}

		if isEnrolled {
			rating, err := h.queries.GetRatingByUserAndCourse(r.Context(), db.GetRatingByUserAndCourseParams{
				CourseID: int32(id),
				UserID:   currentUser.ID,
			})
			if err == nil {
				studentRating = int(rating.Stars)
			}
		}
	}

	lessons, err := h.queries.ListLessonsByCourse(r.Context(), int32(id))
	if err != nil {
		http.Error(w, "failed to fetch lessons", http.StatusInternalServerError)
		return
	}

	data := struct {
		Course        db.GetCourseRow
		CurrentUser   *db.GetUserByIDRow
		IsEnrolled    bool
		StudentRating int
		Lessons       []db.Lesson
	}{
		Course:        course,
		CurrentUser:   currentUser,
		IsEnrolled:    isEnrolled,
		StudentRating: studentRating,
		Lessons:       lessons,
	}

	if err := h.tmpls.ExecuteTemplate(w, "course.html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *CourseHandler) currentUser(r *http.Request) *db.GetUserByIDRow {
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
