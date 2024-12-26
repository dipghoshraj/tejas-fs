package apis

import (
	"net/http"
)

func (nh *NMHandler) SaveFile(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(10 << 20) // 10 MB limit
	if err != nil {
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}

	// filename := uuid.New().String()

}
