package restapi

import (
	"net/http"
)

func removeStudentsFromClass(w http.ResponseWriter, r *http.Request) {
	updateClassStudents(w, r, "remove")
}
