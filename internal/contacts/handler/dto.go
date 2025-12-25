package http

// SyncContactsRequest Exported request struct for Swaggo
type SyncContactsRequest struct {
	Phones []string `json:"phones" binding:"required"`
}
