package domains

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BookingAppointment struct {
    ID 			uuid.UUID `gorm:"column:id;primaryKey;type:uuid;not null" json:"id"`

    ClientID     uuid.UUID      `gorm:"type:uuid;not null" json:"client_id"`     // user with role 'Student' or 'User'
    ReceiverID   uuid.UUID      `gorm:"type:uuid;not null" json:"receiver_id"`   // user with role 'CEO' or 'Manager'
    Date         time.Time      `gorm:"type:date;not null" json:"date"`          // YYYY-MM-DD
    TimeSlot     string         `gorm:"type:varchar(20);not null" json:"timeSlot"` // "10:00 AM"
    Status       string         `gorm:"type:varchar(20);default:'booked'" json:"status"` // booked | cancelled | completed
    Notes        string         `gorm:"type:text" json:"notes,omitempty"`        // optional message
    CreatedAt    time.Time      `json:"created_at"`
    UpdatedAt    time.Time      `json:"updated_at"`
    DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
	MeetLink string `json:"meetLink"`
}
type AppointmentStats struct {
	Total     int64 `json:"total"`
	Confirmed int64 `json:"confirmed"`
	Canceled  int64 `json:"canceled"`
	Completed int64 `json:"completed"`
}

func ReadAppointmentStats(db *gorm.DB, adminID uuid.UUID) (*AppointmentStats, error) {
	var stats AppointmentStats
	if err := db.Model(&BookingAppointment{}).
			Where("receiver_id = ? AND deleted_at IS NULL", adminID).
			Count(&stats.Total).Error; err != nil {
		return nil, err
	}
	db.Model(&BookingAppointment{}).
		Where("receiver_id = ? AND status = ?", adminID, "confirmed").
		Count(&stats.Confirmed)
	db.Model(&BookingAppointment{}).
		Where("receiver_id = ? AND status = ?", adminID, "canceled").
		Count(&stats.Canceled)
	db.Model(&BookingAppointment{}).
		Where("receiver_id = ? AND status = ?", adminID, "completed").
		Count(&stats.Completed)

	return &stats, nil
}

func ReadUpcomingAppointments(db *gorm.DB, userID uuid.UUID) ([]BookingAppointment, error) {
	var upcoming []BookingAppointment
	today := time.Now().Truncate(24 * time.Hour)
	err := db.Where("receiver_id = ? AND date >= ?", userID, today).
		Order("date ASC, time_slot ASC").
		Find(&upcoming).Error
	return upcoming, err
}

func ReadPastAppointments(db *gorm.DB, userID uuid.UUID) ([]BookingAppointment, error) {
	var past []BookingAppointment
	today := time.Now().Truncate(24 * time.Hour)
	err := db.Where("receiver_id = ? AND date < ?", userID, today).
		Order("date DESC, time_slot ASC").
		Find(&past).Error
	return past, err
}

func SendReminderToUser(db *gorm.DB, bookingID uuid.UUID) error {
	// For now this is a placeholder
	return nil
}

func RescheduleBooking(db *gorm.DB, bookingID uuid.UUID, newDate time.Time, newSlot string) error {
	return db.Model(&BookingAppointment{}).
		Where("id = ?", bookingID).
		Updates(map[string]interface{}{
			"date":      newDate,
			"time_slot": newSlot,
		}).Error
}

func ReadStudentAppointments(db *gorm.DB, userID uuid.UUID) ([]BookingAppointment, error) {
	var appointments []BookingAppointment
	err := db.Where("client_id = ?", userID).Find(&appointments).Error
	if err != nil {
		return nil, err
	}
	return appointments, nil
}
func CheckSlotTaken(db *gorm.DB, receiverID uuid.UUID, date time.Time, timeSlot string) (bool, error) {
	var count int64
	err := db.Model(&BookingAppointment{}).
		Where("receiver_id = ? AND date = ? AND time_slot = ? AND deleted_at IS NULL", receiverID, date, timeSlot).
		Count(&count).Error
	return count > 0, err
}


