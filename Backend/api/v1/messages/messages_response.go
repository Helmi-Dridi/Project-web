package messages

import (
	"time"

	"github.com/google/uuid"
)

type MessageResponse struct {
	ID         uuid.UUID  `json:"id"`
	SenderID   uuid.UUID  `json:"senderId"`
	ReceiverID uuid.UUID  `json:"receiverId"`
	Content    string     `json:"content"`
	Attachment *string    `json:"attachment,omitempty"`
	CompanyID  uuid.UUID  `json:"companyId"`
	CreatedAt  time.Time  `json:"createdAt"`
	Read       bool       `json:"read"`
}
