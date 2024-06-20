//go:build !solution

package jsonrpc

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"reflect"
)

func MakeHandler(service interface{}) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		method := reflect.ValueOf(service).MethodByName(r.URL.Path[1:])
		if !method.IsValid() {
			http.Error(w, "Method not found", http.StatusMethodNotAllowed)
			return
		}

		reqType := method.Type().In(1)
		reqValue := reflect.New(reqType.Elem())

		if err := json.NewDecoder(r.Body).Decode(reqValue.Interface()); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		result := method.Call([]reflect.Value{reflect.ValueOf(ctx), reqValue})

		if len(result) != 2 {
			http.Error(w, "method error", http.StatusInternalServerError)
			return
		}

		if !result[1].IsNil() {
			err := result[1].Interface().(error)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		rsp := result[0]

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(rsp.Interface()); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})
}

func Call(ctx context.Context, endpoint string, method string, req, rsp interface{}) error {
	url := fmt.Sprintf("%s/%s", endpoint, method)

	reqBody, err := json.Marshal(req)
	if err != nil {
		return err
	}

	httpClient := http.DefaultClient
	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(reqBody))
	if err != nil {
		return err
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := httpClient.Do(httpReq)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		return errors.New(string(bodyBytes))
	}

	return json.NewDecoder(resp.Body).Decode(rsp)
}
