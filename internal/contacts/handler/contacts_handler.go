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
