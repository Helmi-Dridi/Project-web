package domains

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Message represents a message sent between users.
type Message struct {
	ID         uuid.UUID  `gorm:"column:id;primaryKey;type:uuid;not null" json:"id"`
	SenderID   uuid.UUID  `gorm:"type:uuid;not null" json:"sender_id"`
	ReceiverID uuid.UUID  `gorm:"type:uuid;not null" json:"receiver_id"`
	Content    string     `gorm:"type:text" json:"content"`
	Attachment *string    `gorm:"type:text" json:"attachment,omitempty"` // omit if null
	CompanyID  uuid.UUID  `gorm:"type:uuid;not null" json:"company_id"`
	CreatedAt  time.Time  `json:"created_at"`
	Read       bool       `json:"read"`
}
func CreateMessage(db *gorm.DB, message *Message) error {
	return db.Create(message).Error
}
func GetConversationBetweenUsers(db *gorm.DB, user1, user2 uuid.UUID) ([]Message, error) {
	var messages []Message
	err := db.
		Where("(sender_id = ? AND receiver_id = ?) OR (sender_id = ? AND receiver_id = ?)", user1, user2, user2, user1).
		Order("created_at ASC").
		Find(&messages).Error
	return messages, err
}
func MarkMessageAsRead(db *gorm.DB, messageID uuid.UUID) error {
	return db.Model(&Message{}).Where("id = ?", messageID).Update("read", true).Error
}
func GetAllConversationPartners(db *gorm.DB, adminID uuid.UUID) ([]uuid.UUID, error) {
	var userIDs []uuid.UUID
	err := db.
		Model(&Message{}).
		Select("DISTINCT sender_id").
		Where("receiver_id = ?", adminID).
		Pluck("sender_id", &userIDs).Error
	return userIDs, err
}
