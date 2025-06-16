package bookingappointment

import (
	"labs/domains"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// Database holds the DB instance for use in handlers.
type Database struct {
	DB *gorm.DB
}

// NewBookingAppointmentRepository migrates the BookingAppointment model.
func NewBookingAppointmentRepository(db *gorm.DB) {
	if err := db.AutoMigrate(&domains.BookingAppointment{}); err != nil {
		logrus.Fatal("Error migrating BookingAppointment structure: ", err)
	}
}

// CreateBooking creates a new appointment booking.
func CreateBooking(db *gorm.DB, model *domains.BookingAppointment) error {
	return db.Create(model).Error
}

// ReadUserBookingsPaginated retrieves paginated appointments for a user.
func ReadUserBookingsPaginated(db *gorm.DB, userID uuid.UUID, limit, offset int) ([]domains.BookingAppointment, error) {
	var bookings []domains.BookingAppointment
	err := db.
		Where("client_id = ?", userID).
		Limit(limit).
		Offset(offset).
		Order("date asc, time_slot asc").
		Find(&bookings).Error
	return bookings, err
}

// ReadUserBookingsAll retrieves all appointments for a user (non-paginated).
func ReadUserBookingsAll(db *gorm.DB, userID uuid.UUID) ([]domains.BookingAppointment, error) {
	var bookings []domains.BookingAppointment
	err := db.
		Where("client_id = ?", userID).
		Order("date asc, time_slot asc").
		Find(&bookings).Error
	return bookings, err
}

// ReadBookingByID retrieves a specific booking by ID.
func ReadBookingByID(db *gorm.DB, id uuid.UUID) (*domains.BookingAppointment, error) {
	var booking domains.BookingAppointment
	err := db.First(&booking, id).Error
	return &booking, err
}

// ReadAvailabilityByDate returns all appointments booked for a given date.
func ReadAvailabilityByDate(db *gorm.DB, date string) ([]domains.BookingAppointment, error) {
	var bookings []domains.BookingAppointment
	err := db.Where("date = ?", date).Find(&bookings).Error
	return bookings, err
}

// CancelBooking updates a booking's status to canceled.
func CancelBooking(db *gorm.DB, id uuid.UUID) error {
	return db.Model(&domains.BookingAppointment{}).Where("id = ?", id).Update("status", "canceled").Error
}
func ReadReceivedBookings(db *gorm.DB, receiverID uuid.UUID) ([]domains.BookingAppointment, error) {
	var bookings []domains.BookingAppointment
	err := db.Where("receiver_id = ?", receiverID).
		Order("date DESC").
		Find(&bookings).Error
	return bookings, err
}
func parseTime(date time.Time, timeSlot string) (time.Time, error) {
	layout := "2006-01-02 03:04 PM"
	combined := date.Format("2006-01-02") + " " + timeSlot
	return time.ParseInLocation(layout, combined, time.Local)
}
func ReadConfirmedAppointmentsForUser(db *gorm.DB, userID uuid.UUID) ([]domains.BookingAppointment, error) {
	var confirmed []domains.BookingAppointment
	err := db.Where("client_id = ? AND status = ?", userID, "confirmed").
		Order("date ASC, time_slot ASC").
		Find(&confirmed).Error
	return confirmed, err
}
