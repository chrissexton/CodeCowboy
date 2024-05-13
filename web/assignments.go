package web

import (
	"cso/codecowboy/classroom"
	"fmt"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func (w *Web) setupAssignmentHandlers() chi.Router {
	router := chi.NewRouter()

	router.Get("/newAssignment", w.handleNewAssignmentForm)
	router.Get("/{assignment}", w.handleAssignmentDetails)
	router.Delete("/{assignment}", w.handleRmAssignment)
	router.Post("/", w.handleNewAssignment)

	return router
}

func (w *Web) handleAssignmentDetails(wr http.ResponseWriter, r *http.Request) {
	courseName := chi.URLParam(r, "course")
	assignmentName := chi.URLParam(r, "assignment")
	course, err := classroom.New(w.db, courseName)
	if err != nil {
		w.renderErr(r.Context(), wr, err)
	}
	for _, a := range course.Assignments {
		if a.Name == assignmentName {
			if a.Expr == "" {
				a.Expr = classroom.DEFAULT_EXPR
			}
			w.Index(assignmentName, w.assignmentDetails(a)).Render(r.Context(), wr)
			return
		}
	}
	w.renderErr(r.Context(), wr, fmt.Errorf("could not find assignment"))
}

func (w *Web) handleRmAssignment(wr http.ResponseWriter, r *http.Request) {
	assignmentName := chi.URLParam(r, "assignment")
	courseName := chi.URLParam(r, "course")
	cls, err := classroom.New(w.db, courseName)
	if err != nil {
		w.renderErr(r.Context(), wr, err)
		return
	}
	assignments := classroom.Assignments{}
	for _, a := range cls.Assignments {
		if a.Name != assignmentName {
			assignments = append(assignments, a)
		}
	}
	cls.Assignments = assignments
	err = cls.Save()
	if err != nil {
		w.renderErr(r.Context(), wr, err)
		return
	}
	wr.Header().Set("HX-Redirect", "/courses/"+cls.Name)
}

func (w *Web) handleNewAssignmentForm(wr http.ResponseWriter, r *http.Request) {
	w.Index("New Assignment", nil).Render(r.Context(), wr)
}

func (w *Web) handleNewAssignment(wr http.ResponseWriter, r *http.Request) {
	cls, err := classroom.New(w.db, r.FormValue("course"))
	if err != nil {
		w.renderErr(r.Context(), wr, err)
		return
	}
	assign := classroom.AssignmentSpec{
		Name:      r.FormValue("name"),
		Path:      r.FormValue("path"),
		Course:    cls.Name,
		ExtrasSrc: r.FormValue("extrasSrc"),
		ExtrasDst: r.FormValue("extrasDst"),
		Expr:      r.FormValue("expr"),
	}

	cls.Assignments = append(cls.Assignments, assign)

	err = cls.Save()
	if err != nil {
		w.renderErr(r.Context(), wr, err)
		return
	}

	wr.Header().Set("HX-Redirect", "/courses/"+cls.Name)
	w.courseDetails(cls).Render(r.Context(), wr)
}