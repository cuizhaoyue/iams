package server

import (
	"net/http"

	"github.com/cuizhaoyue/iams/pkg/log"
)

// ServeHealthCheck 运行一个http服务用来提供检查pump服务健康状态的api.
func ServeHealthCheck(healthPath string, healthAddress string) {
	http.HandleFunc("/"+healthPath, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status": "ok"}`))
	})

	if err := http.ListenAndServe(healthAddress, nil); err != nil {
		log.Fatalf("Error serving health check endpoint: %s", err.Error())
	}
}
