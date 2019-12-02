package views

type HealthCheckResponse struct {
	Name     string `json:"name"`
	Database bool   `json:"database"`
}
