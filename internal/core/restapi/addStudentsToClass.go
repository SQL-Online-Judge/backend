package restapi

import (
	"net/http"
)

func addStudentsToClass(w http.ResponseWriter, r *http.Request) {
	updateClassStudents(w, r, "add")
}
