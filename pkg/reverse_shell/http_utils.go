package reverse_shell

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
)

func ParseJSONBody(r *http.Request, v any) error {
	if r.Header.Get("Content-Type") != "application/json" {
		return errors.New("content type mismatch")
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}
	err = json.Unmarshal(body, v)
	return err
}

func ErrorBadRequest(w http.ResponseWriter) {
	http.Error(w, "bad request", http.StatusBadRequest)
}

func ErrorInternal(w http.ResponseWriter) {
	http.Error(w, "internal server error", http.StatusInternalServerError)

}

func HandleRestFunc(method string, path string, handler HttpHandleFunc) {
	http.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != method {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		handler(w, r)
	})
}

func WriteJSON(w http.ResponseWriter, v interface{}) error {
	jbytes, err := json.Marshal(v)
	if err != nil {
		return err
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jbytes)
	return nil
}
