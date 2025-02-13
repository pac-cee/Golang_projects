package health

import (
    "encoding/json"
    "net/http"
)

type Health struct {
    Status string `json:"status"`
}

func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
    health := Health{Status: "UP"}
    json.NewEncoder(w).Encode(health)
}