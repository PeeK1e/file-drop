package routes

import "net/http"

func RedirectStatic(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/ul", 301)
}
