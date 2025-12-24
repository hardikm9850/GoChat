package handler

import (
    "log"
    "net/http"
    "time"

    "github.com/gin-gonic/gin"
    "github.com/gorilla/websocket"
    "github.com/hardikm9850/GoChat/internal/chat/domain"
    "github.com/hardikm9850/GoChat/internal/chat/usecase"
    "github.com/hardikm9850/GoChat/internal/hub"
)

var upgrader = websocket.Upgrader{
    CheckOrigin: func(r *http.Request) bool { return true },
}

type WSHandler struct {
    Hub                *hub.Hub
    SendMessageUseCase *usecase.SendMessageUseCase
}

type IncomingMessage struct {
    ConversationID string `json:"conversation_id"`
    Content        string `json:"content"`
}

type WSError struct {
    Type    string `json:"type"`
    Message string `json:"message"`
}

func NewWSHandler(h *hub.Hub, uc *usecase.SendMessageUseCase) *WSHandler {
    return &WSHandler{
        Hub:                h,
        SendMessageUseCase: uc,
    }
}

func (h *WSHandler) HandleWebSocket(c *gin.Context) {
    userID := c.GetString("userID")
    log.Println("Control received in HandleWebSocket")
    if userID == "" {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
        return
    }
    log.Println("WS connected user:", userID)

    conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
    if err != nil {
        return
    }
    // register connection in Hub via channel
    h.Hub.Register <- hub.ConnEvent{
        UserID: userID,
        Conn:   conn,
    } // we pass the connection reference received after upgrading http connection to long-lived ws

    // defer means run below func Run after the surrounding function finishes, no matter how it finishes.
    // In this case, readloop returns or HandleWebSocket returns then the deferred func executes
    // This is stack-based cleanup or try .. finally
    defer func() {
        h.Hub.Unregister <- hub.ConnEvent{
            UserID: userID,
            Conn:   conn,
        }
        _ = conn.Close()
    }()

    h.readLoop(conn, userID)
}

// readLoop Performs the followings
// Deserialize JSON,
// Identify sender (userID)
// Call application logic
func (h *WSHandler) readLoop(conn *websocket.Conn, userID string) {
    conn.SetReadDeadline(time.Now().Add(60 * time.Second))

    for {
        var payload IncomingMessage

        if err := conn.ReadJSON(&payload); err != nil {
            log.Println("ReadJSON error:", err)
            return
        }

        conn.SetReadDeadline(time.Now().Add(60 * time.Second))
        h.handleIncomingMessage(conn, userID, payload)
    }
}

func (h *WSHandler) handleIncomingMessage(
    conn *websocket.Conn,
    userID string,
    payload IncomingMessage,
) {
    log.Println("handleIncomingMessage invoked to Execute SMUC")
    _, err := h.SendMessageUseCase.Execute(
        domain.UserID(userID),
        domain.ConversationID(payload.ConversationID),
        payload.Content,
    )

    if err != nil {
        _ = conn.WriteJSON(WSError{
            Type:    "SEND_MESSAGE_FAILED",
            Message: err.Error(),
        })
        return
    }

}
