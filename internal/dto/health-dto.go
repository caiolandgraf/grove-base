package dto

// HealthCheckResponse is the health check response
type HealthCheckResponse struct {
	Status string `json:"status"`
}

// HealthCheckDetailedResponse is the detailed health check response
type HealthCheckDetailedResponse struct {
	Status   string            `json:"status"`
	Services map[string]string `json:"services"`
}
