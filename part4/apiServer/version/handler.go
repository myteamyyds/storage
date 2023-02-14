package version

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"storage/lib/es"
	"strings"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	method := r.Method
	fmt.Println("进入versions + ", method)
	if method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	from := 0
	size := 1000
	name := strings.Split(r.URL.EscapedPath(), "/")[2]
	for {
		metas, err := es.SearchAllVersions(name, from, size)
		fmt.Println(metas)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		for i := range metas {
			body, _ := json.Marshal(metas[i])
			w.Write(body)
			w.Write([]byte("\n"))
		}
		if len(metas) != size {
			return
		}
		from += size
	}
}
