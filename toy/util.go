package toy

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

func MarshalStringAndHttpWriteStatus(w http.ResponseWriter, msg string, statusCode int) {
	w.Header().Set("Content-Type", "text/html; charset=UTF-8")
	w.WriteHeader(statusCode)
	w.Write([]byte(msg))

}
func MarshalJsonAndHttpWriteStatus(w http.ResponseWriter, obj interface{}, statusCode int) error {
	if obj != nil {
		b, err := json.MarshalIndent(obj, "", "\t")
		if err == nil {
			w.WriteHeader(statusCode)
			w.Write(b)
			w.Write([]byte("\n")) //Send \n so that when using CURL it is easier to cut and paste output
		} else {
			return fmt.Errorf("MarshalJsonAndHttpWriteStatus: %s\n", err)
		}
		return nil
	} else {
		return errors.New("Object for marshaling was nil")
	}
}
