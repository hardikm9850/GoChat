package http

import (
	"github.com/hardikm9850/GoChat/internal/contacts/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ContactsHandler struct {
	contactService service.ContactService
}

func NewContactsHandler(contactService service.ContactService) *ContactsHandler {
	return &ContactsHandler{
		contactService: contactService,
	}
}

// @Summary Sync contacts
// @Description Syncs a list of phone numbers with the server
// @Tags Sync contacts
// @Accept json
// @Produce json
// @Param request body SyncContactsRequest true "Sync Contacts Request"
// @Success 200 {array} string "List of synced contacts"
// @Failure 400 {object} map[string]string "Invalid request body"
// @Failure 500 {object} map[string]string "Failed to sync contacts"
// @Router /contacts/sync [post]
func (h *ContactsHandler) SyncContacts(c *gin.Context) {
	var req SyncContactsRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body",
		})
		return
	}

	contacts, err := h.contactService.SyncContacts(req.Phones)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to sync contacts",
		})
		return
	}

	c.JSON(http.StatusOK, contacts)
}
