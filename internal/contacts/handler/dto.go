package http

type SyncContactsRequest struct {
	Phones []string `json:"phones" binding:"required"`
}
