package restapi

import (
	"net/http"
)

func removeProblemsFromTask(w http.ResponseWriter, r *http.Request) {
	updateTaskProblem(w, r, "remove")
}
