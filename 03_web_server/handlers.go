package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type Response struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func jsonResponse(w http.ResponseWriter, status int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshalling JSON: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(response)
}

func aboutHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		jsonResponse(w, http.StatusMethodNotAllowed, Response{
			Status:  "error",
			Message: "Method not allowed",
		})
		return
	}

	jsonResponse(w, http.StatusOK, Response{
		Status:  "success",
		Message: "About page",
		Data: map[string]string{
			"description": "This is a professional Go web server",
			"version":     "1.0.0",
		},
	})
}

func contactHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		jsonResponse(w, http.StatusMethodNotAllowed, Response{
			Status:  "error",
			Message: "Method not allowed",
		})
		return
	}

	jsonResponse(w, http.StatusOK, Response{
		Status:  "success",
		Message: "Contact information",
		Data: map[string]string{
			"email":    "contact@example.com",
			"phone":    "+1-555-555-5555",
			"location": "San Francisco, CA",
		},
	})
}

func servicesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		jsonResponse(w, http.StatusMethodNotAllowed, Response{
			Status:  "error",
			Message: "Method not allowed",
		})
		return
	}

	jsonResponse(w, http.StatusOK, Response{
		Status:  "success",
		Message: "Our services",
		Data: []map[string]string{
			{
				"name":        "Web Development",
				"description": "Full-stack web development services",
			},
			{
				"name":        "Cloud Hosting",
				"description": "Secure and scalable cloud hosting solutions",
			},
			{
				"name":        "API Integration",
				"description": "Custom API development and integration",
			},
		},
	})
}
