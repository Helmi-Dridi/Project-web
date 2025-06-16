package bookingappointment

import (
	"time"

	"github.com/google/uuid"
)

// @Description AppointmentBookingInput represents booking request data.
// Used in: POST /appointments
// @name AppointmentBookingInput
type AppointmentBookingInput struct {
	Date     string    `json:"date" binding:"required"`
	TimeSlot string    `json:"timeSlot" binding:"required"`
	AdminID  uuid.UUID `json:"receiverId" binding:"required"` // 👈 must match frontend field
}

type UpdateAppointmentStatusInput struct {
	Status string `json:"status" binding:"required,oneof=booked confirmed completed canceled"`
}

// @Description AppointmentView is a simplified view for listing user bookings.
// @name AppointmentView
type AppointmentView struct {
	ID       uuid.UUID `json:"id"`
	Date     string    `json:"date"`
	TimeSlot string    `json:"timeSlot"`
	Status   string    `json:"status"` // e.g. "booked", "cancelled"
}

// @Description AppointmentPagination wraps paginated bookings for frontend.
// @name AppointmentPagination
type AppointmentPagination struct {
	Items      []AppointmentView `json:"items"`
	Page       int               `json:"page"`
	Limit      int               `json:"limit"`
	TotalCount int               `json:"totalCount"`
}

// @Description BookingAppointmentTable is used for listing booked appointments.
// @name BookingAppointmentTable
type BookingAppointmentTable struct {
	ID       uuid.UUID `json:"id"`
	Date     string    `json:"date"`
	TimeSlot string    `json:"timeSlot"`
	Status   string    `json:"status"` // e.g. "pending", "confirmed", "canceled"
}

// @Description BookingAppointmentPagination holds paginated list of bookings.
// @name BookingAppointmentPagination
type BookingAppointmentPagination struct {
	Items      []BookingAppointmentTable `json:"items"`
	Page       uint                      `json:"page"`
	Limit      uint                      `json:"limit"`
	TotalCount uint                      `json:"totalCount"`
}

// @Description BookingAppointmentDetails contains full details for a booking.
// @name BookingAppointmentDetails
type BookingAppointmentDetails struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"userID"`
	Date      string    `json:"date"`
	TimeSlot  string    `json:"timeSlot"`
	Status    string    `json:"status"` // e.g. "pending", "confirmed", "canceled"
	CreatedAt time.Time `json:"createdAt"`
}

// @Description BookingAvailabilitySlot shows availability for a time slot.
// @name BookingAvailabilitySlot
type BookingAvailabilitySlot struct {
	TimeSlot   string `json:"timeSlot"`   // e.g. "10:00 AM"
	Available  bool   `json:"available"`  // Whether this slot is still available
	BookedBy   uint   `json:"bookedBy"`   // How many have booked this slot (for group support)
	MaxAllowed uint   `json:"maxAllowed"` // Max allowed bookings per slot
}

// @Description BookingAvailabilityDay shows available slots for a given date.
// @name BookingAvailabilityDay
type BookingAvailabilityDay struct {
	Date  string                     `json:"date"`  // e.g. "2025-06-06"
	Slots []BookingAvailabilitySlot `json:"slots"` // List of time slots and availability
}
