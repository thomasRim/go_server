package login

import "net/http"

func checkBaseAuth(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)

	username, password, authOK := r.BasicAuth()
	if authOK == false {
		http.Error(w, "Not authorized", 401)
		return
	}

	if username != "username" || password != "password" {
		http.Error(w, "Not authorized", 401)
		return
	}
}
