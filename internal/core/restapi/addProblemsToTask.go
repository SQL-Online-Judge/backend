package restapi

import (
	"net/http"
)

func addProblemsToTask(w http.ResponseWriter, r *http.Request) {
	updateTaskProblem(w, r, "add")
}
