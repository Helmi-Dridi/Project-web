package messages

import (
	"labs/domains"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type Database struct {
	DB *gorm.DB
}

func NewMessageRepository(db *gorm.DB) {
	if err := db.AutoMigrate(&domains.Message{}); err != nil {
		logrus.Fatal("Failed to migrate message model: ", err)
	}
}

func CreateMessage(db *gorm.DB, message *domains.Message) error {
	return db.Create(message).Error
}

func GetMessagesBetween(db *gorm.DB, user1, user2 uuid.UUID) ([]domains.Message, error) {
	var messages []domains.Message
	err := db.
		Where("(sender_id = ? AND receiver_id = ?) OR (sender_id = ? AND receiver_id = ?)", user1, user2, user2, user1).
		Order("created_at").
		Find(&messages).Error
	return messages, err
}

func MarkMessageRead(db *gorm.DB, messageID uuid.UUID) error {
	return db.Model(&domains.Message{}).Where("id = ?", messageID).Update("read", true).Error
}
func GetConversationPartners(db *gorm.DB, userID uuid.UUID) ([]uuid.UUID, error) {
	var partnerIDs []uuid.UUID
	err := db.Model(&domains.Message{}).
		Select("DISTINCT CASE WHEN sender_id = ? THEN receiver_id ELSE sender_id END", userID).
		Where("sender_id = ? OR receiver_id = ?", userID, userID).
		Pluck("receiver_id", &partnerIDs).Error
	return partnerIDs, err
}
func CountUnreadMessages(db *gorm.DB, receiverID uuid.UUID) (int64, error) {
	var count int64
	err := db.Model(&domains.Message{}).
		Where("receiver_id = ? AND read = false", receiverID).
		Count(&count).Error
	return count, err
}
func GetPaginatedMessages(db *gorm.DB, user1, user2 uuid.UUID, limit, offset int) ([]domains.Message, error) {
	var messages []domains.Message
	err := db.
		Where("(sender_id = ? AND receiver_id = ?) OR (sender_id = ? AND receiver_id = ?)", user1, user2, user2, user1).
		Order("created_at desc").
		Limit(limit).
		Offset(offset).
		Find(&messages).Error
	return messages, err
}
func DeleteMessage(db *gorm.DB, messageID, userID uuid.UUID) error {
	return db.Where("id = ? AND sender_id = ?", messageID, userID).Delete(&domains.Message{}).Error
}
func SearchMessages(db *gorm.DB, companyID uuid.UUID, userID uuid.UUID, keyword string) ([]domains.Message, error) {
	var messages []domains.Message
	err := db.
		Where("company_id = ? AND (sender_id = ? OR receiver_id = ?) AND content ILIKE ?", companyID, userID, userID, "%"+keyword+"%").
		Order("created_at desc").
		Find(&messages).Error
	return messages, err
}
func GetAdminInbox(db *gorm.DB, adminID uuid.UUID) (map[uuid.UUID][]domains.Message, error) {
	var messages []domains.Message
	err := db.Where("receiver_id = ?", adminID).Order("created_at asc").Find(&messages).Error
	if err != nil {
		return nil, err
	}

	grouped := make(map[uuid.UUID][]domains.Message)
	for _, msg := range messages {
		grouped[msg.SenderID] = append(grouped[msg.SenderID], msg)
	}
	return grouped, nil
}
