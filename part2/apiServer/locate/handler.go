package locate

import "net/http"

func Handler(w http.ResponseWriter, r *http.Request) {
	method := r.Method
	if method == r.Method {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

}
