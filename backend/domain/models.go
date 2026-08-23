package domain

type Alert struct {
	ID           string `json:"id"`
	Region       string `json:"region"`
	Severity     string `json:"severity"`
	Status       string `json:"status"`
	DetectedAt   string `json:"detected_at"`
	CrewAssigned int    `json:"crew_assigned"`
}
type StatusChange struct {
	Status string `json:"status"`
}
