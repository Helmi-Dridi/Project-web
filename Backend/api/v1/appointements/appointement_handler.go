package bookingappointment

import (
	"labs/constants"
	"labs/domains"
	"labs/utils"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

// CreateBooking handles POST /appointments
func (db Database) CreateBooking(ctx *gin.Context) {
	session := utils.ExtractJWTValues(ctx)
	input := new(AppointmentBookingInput)

	if err := ctx.ShouldBindJSON(input); err != nil {
		logrus.Warn("Invalid JSON input: ", err)
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	logrus.Infof("Received booking request: user=%s admin=%s date=%s slot=%s",
		session.UserID, input.AdminID, input.Date, input.TimeSlot)

	// Parse date
	parsedDate, err := time.Parse("2006-01-02", input.Date)
	if err != nil {
		logrus.Error("Invalid date format: ", input.Date)
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, "Invalid date format (expected YYYY-MM-DD)", utils.Null())
		return
	}

	// Validate AdminID
	if input.AdminID == uuid.Nil {
		logrus.Warn("Missing AdminID in request")
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, "Receiver (admin) ID is required", utils.Null())
		return
	}

	// Check if slot is already booked
	exists, err := domains.CheckSlotTaken(db.DB, input.AdminID, parsedDate, input.TimeSlot)
	if err != nil {
		logrus.Error("Database error while checking availability: ", err)
		utils.BuildErrorResponse(ctx, http.StatusInternalServerError, constants.UNKNOWN_ERROR, utils.Null())
		return
	}
	if exists {
		logrus.Warnf("Slot already taken: admin=%s date=%s time=%s", input.AdminID, parsedDate, input.TimeSlot)
		utils.BuildErrorResponse(ctx, http.StatusConflict, "This time slot is already booked for the selected admin", gin.H{
			"date":     input.Date,
			"timeSlot": input.TimeSlot,
		})
		return
	}

	// Proceed with booking
	booking := &domains.BookingAppointment{
		ID:         uuid.New(),
		ClientID:   session.UserID,
		ReceiverID: input.AdminID,
		Date:       parsedDate,
		TimeSlot:   input.TimeSlot,
		Status:     "booked",
		CreatedAt:  time.Now(),
	}

	if err := CreateBooking(db.DB, booking); err != nil {
		logrus.Error("Failed to create booking: ", err)
		utils.BuildErrorResponse(ctx, http.StatusInternalServerError, constants.UNKNOWN_ERROR, utils.Null())
		return
	}

	logrus.Infof("Booking created: ID=%s client=%s admin=%s", booking.ID, booking.ClientID, booking.ReceiverID)

	utils.BuildResponse(ctx, http.StatusCreated, constants.CREATED, gin.H{
		"id":        booking.ID,
		"date":      input.Date,
		"timeSlot":  input.TimeSlot,
		"receiver":  input.AdminID,
		"client_id": session.UserID,
	})
}

// GetUserAppointments handles GET /appointments
func (db Database) GetUserAppointments(ctx *gin.Context) {
	session := utils.ExtractJWTValues(ctx)
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))

	if page <= 0 || limit <= 0 {
		logrus.Warn("Invalid pagination parameters")
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	offset := (page - 1) * limit

	bookings, err := ReadUserBookingsPaginated(db.DB, session.UserID, limit, offset)
	if err != nil {
		logrus.Error("Error reading user bookings: ", err)
		utils.BuildErrorResponse(ctx, http.StatusInternalServerError, constants.UNKNOWN_ERROR, utils.Null())
		return
	}

	var response []AppointmentView
	for _, b := range bookings {
		response = append(response, AppointmentView{
			ID:       b.ID,
			Date:     b.Date.Format("2006-01-02"),
			TimeSlot: b.TimeSlot,
			Status:   b.Status,
		})
	}

	count, err := domains.ReadTotalCount(db.DB, &domains.BookingAppointment{}, "client_id", session.UserID)
	if err != nil {
		logrus.Error("Error counting user bookings: ", err)
		utils.BuildErrorResponse(ctx, http.StatusInternalServerError, constants.UNKNOWN_ERROR, utils.Null())
		return
	}

	utils.BuildResponse(ctx, http.StatusOK, constants.SUCCESS, gin.H{
		"items":      response,
		"page":       page,
		"limit":      limit,
		"totalCount": count,
	})
}

// GetAvailabilityByDate handles GET /appointments/availability?date=YYYY-MM-DD
func (db Database) GetAvailabilityByDate(ctx *gin.Context) {
	date := ctx.Query("date")
	if date == "" {
		logrus.Warn("Missing date parameter in availability request")
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	bookings, err := ReadAvailabilityByDate(db.DB, date)
	if err != nil {
		logrus.Error("Error reading availability: ", err)
		utils.BuildErrorResponse(ctx, http.StatusInternalServerError, constants.UNKNOWN_ERROR, utils.Null())
		return
	}

	var bookedSlots []string
	for _, b := range bookings {
		bookedSlots = append(bookedSlots, b.TimeSlot)
	}

	logrus.Infof("Fetched availability for date=%s bookedSlots=%v", date, bookedSlots)

	utils.BuildResponse(ctx, http.StatusOK, constants.SUCCESS, gin.H{
		"bookedSlots": bookedSlots,
	})
}

// CancelBooking handles PATCH /appointments/:id/cancel
func (db Database) CancelBooking(ctx *gin.Context) {
	session := utils.ExtractJWTValues(ctx)
	bookingID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		logrus.Warn("Invalid booking ID format")
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	booking, err := ReadBookingByID(db.DB, bookingID)
	if err != nil || booking.ClientID != session.UserID {
		logrus.Warnf("Booking not found or unauthorized access: user=%s booking=%s", session.UserID, bookingID)
		utils.BuildErrorResponse(ctx, http.StatusNotFound, constants.DATA_NOT_FOUND, utils.Null())
		return
	}

	if booking.Status != "booked" {
		logrus.Warnf("Attempt to cancel non-booked status: ID=%s status=%s", bookingID, booking.Status)
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, "Booking already canceled or invalid", utils.Null())
		return
	}

	if err := CancelBooking(db.DB, bookingID); err != nil {
		logrus.Error("Error canceling booking: ", err)
		utils.BuildErrorResponse(ctx, http.StatusInternalServerError, constants.UNKNOWN_ERROR, utils.Null())
		return
	}

	logrus.Infof("Booking canceled: ID=%s by user=%s", bookingID, session.UserID)

	utils.BuildResponse(ctx, http.StatusOK, constants.SUCCESS, gin.H{
		"message": "Booking canceled successfully",
	})
}
func (db Database) GetReceivedAppointments(ctx *gin.Context) {
	session := utils.ExtractJWTValues(ctx)
	dateParam := ctx.Query("date") // optional query param

	query := db.DB.Model(&domains.BookingAppointment{}).Where("receiver_id = ?", session.UserID)

	// ✅ Fix: Ensure date is compared as a pure DATE string to avoid time mismatch
	if dateParam != "" {
		_, err := time.Parse("2006-01-02", dateParam)
		if err != nil {
			logrus.Warnf("Invalid date param: %s", dateParam)
			utils.BuildErrorResponse(ctx, http.StatusBadRequest, "Invalid date format (expected YYYY-MM-DD)", utils.Null())
			return
		}
		// Use SQL string comparison rather than Go time.Time
		query = query.Where("DATE(date) = ?", dateParam)
	}

	var bookings []domains.BookingAppointment
	if err := query.Order("date DESC, time_slot ASC").Find(&bookings).Error; err != nil {
		logrus.Error("Error reading received appointments: ", err)
		utils.BuildErrorResponse(ctx, http.StatusInternalServerError, constants.UNKNOWN_ERROR, utils.Null())
		return
	}

	var response []AppointmentView
	for _, b := range bookings {
		response = append(response, AppointmentView{
			ID:       b.ID,
			Date:     b.Date.Format("2006-01-02"),
			TimeSlot: b.TimeSlot,
			Status:   b.Status,
		})
	}

	utils.BuildResponse(ctx, http.StatusOK, constants.SUCCESS, gin.H{
		"items": response,
	})
}


// UpdateAppointmentStatus handles PATCH /appointments/:id/status
func (db Database) UpdateAppointmentStatus(ctx *gin.Context) {
	session := utils.ExtractJWTValues(ctx)

	bookingID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		logrus.Warn("Invalid booking ID format")
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	var input UpdateAppointmentStatusInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		logrus.Warn("Invalid status payload: ", err)
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	// Lookup booking
	booking, err := ReadBookingByID(db.DB, bookingID)
	if err != nil {
		logrus.Error("Booking not found: ", err)
		utils.BuildErrorResponse(ctx, http.StatusNotFound, constants.DATA_NOT_FOUND, utils.Null())
		return
	}

	// Only allow the admin assigned to the appointment to update it
	if booking.ReceiverID != session.UserID {
		logrus.Warnf("Forbidden: Admin %s tried to update booking not assigned to them", session.UserID)
		utils.BuildErrorResponse(ctx, http.StatusForbidden, constants.UNAUTHORIZED, utils.Null())
		return
	}

	// Update the status
	err = db.DB.Model(&booking).Update("status", input.Status).Error
	if err != nil {
		logrus.Error("Failed to update booking status: ", err)
		utils.BuildErrorResponse(ctx, http.StatusInternalServerError, constants.UNKNOWN_ERROR, utils.Null())
		return
	}

	// 🔁 If status is confirmed, trigger Google Meet link creation in background
	if input.Status == "confirmed" {
		go func() {
			startTime, err := parseTime(booking.Date, booking.TimeSlot)
			if err != nil {
				logrus.Warnf("⛔ Failed to parse appointment time: %v", err)
				return
			}

			// ✅ Call without attendee email
			meetLink, err := CreateGoogleMeetEvent(
				"Rendez-vous confirmé",
				"Lien généré automatiquement",
				startTime,
				startTime.Add(30*time.Minute),
			)
			if err != nil {
				logrus.Error("❌ Failed to create Google Meet link: ", err)
				return
			}

			if err := db.DB.Model(&booking).Update("meet_link", meetLink).Error; err != nil {
				logrus.Error("❌ Failed to save Meet link to database: ", err)
			} else {
				logrus.Infof("✅ Google Meet link generated and saved for appointment %s", booking.ID)
			}
		}()
	}

	utils.BuildResponse(ctx, http.StatusOK, constants.SUCCESS, gin.H{
		"message": "Appointment status updated",
		"id":      booking.ID,
		"status":  input.Status,
	})
}

// ----------------------
// 🆕 NEW API HANDLERS
// ----------------------

// GET /appointments/:id
func (db Database) GetAppointmentByID(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, "Invalid appointment ID", utils.Null())
		return
	}

	appt, err := ReadBookingByID(db.DB, id)
	if err != nil {
		utils.BuildErrorResponse(ctx, http.StatusNotFound, constants.DATA_NOT_FOUND, utils.Null())
		return
	}

	utils.BuildResponse(ctx, http.StatusOK, constants.SUCCESS, appt)
}

// GET /appointments/stats
func (db Database) GetAppointmentStats(ctx *gin.Context) {
	session := utils.ExtractJWTValues(ctx)
	counts, err := domains.ReadAppointmentStats(db.DB, session.UserID)
	if err != nil {
		utils.BuildErrorResponse(ctx, http.StatusInternalServerError, constants.UNKNOWN_ERROR, utils.Null())
		return
	}

	utils.BuildResponse(ctx, http.StatusOK, constants.SUCCESS, counts)
}

// GET /appointments/upcoming
func (db Database) GetUpcomingAppointments(ctx *gin.Context) {
	session := utils.ExtractJWTValues(ctx)
	upcoming, err := domains.ReadUpcomingAppointments(db.DB, session.UserID)
	if err != nil {
		utils.BuildErrorResponse(ctx, http.StatusInternalServerError, constants.UNKNOWN_ERROR, utils.Null())
		return
	}

	utils.BuildResponse(ctx, http.StatusOK, constants.SUCCESS, upcoming)
}

// GET /appointments/history
func (db Database) GetPastAppointments(ctx *gin.Context) {
	session := utils.ExtractJWTValues(ctx)
	history, err := domains.ReadPastAppointments(db.DB, session.UserID)
	if err != nil {
		utils.BuildErrorResponse(ctx, http.StatusInternalServerError, constants.UNKNOWN_ERROR, utils.Null())
		return
	}

	utils.BuildResponse(ctx, http.StatusOK, constants.SUCCESS, history)
}

// POST /appointments/:id/remind
func (db Database) SendReminder(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, "Invalid ID", utils.Null())
		return
	}

	err = domains.SendReminderToUser(db.DB, id)
	if err != nil {
		utils.BuildErrorResponse(ctx, http.StatusInternalServerError, "Failed to send reminder", utils.Null())
		return
	}

	utils.BuildResponse(ctx, http.StatusOK, constants.SUCCESS, gin.H{"message": "Reminder sent"})
}

// POST /appointments/:id/reschedule
func (db Database) RescheduleBooking(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, "Invalid ID", utils.Null())
		return
	}

	type Req struct {
		NewDate    string `json:"date" binding:"required"`
		NewTime    string `json:"timeSlot" binding:"required"`
	}

	var body Req
	if err := ctx.ShouldBindJSON(&body); err != nil {
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, "Missing fields", utils.Null())
		return
	}

	parsedDate, err := time.Parse("2006-01-02", body.NewDate)
	if err != nil {
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, "Invalid date format", utils.Null())
		return
	}

	if err := domains.RescheduleBooking(db.DB, id, parsedDate, body.NewTime); err != nil {
		utils.BuildErrorResponse(ctx, http.StatusInternalServerError, "Failed to reschedule", utils.Null())
		return
	}

	utils.BuildResponse(ctx, http.StatusOK, constants.SUCCESS, gin.H{"message": "Booking rescheduled"})
}
func (db Database) GetAppointmentHistory(ctx *gin.Context) {
	session := utils.ExtractJWTValues(ctx)

	bookings, err := domains.ReadPastAppointments(db.DB, session.UserID)
	if err != nil {
		utils.BuildErrorResponse(ctx, http.StatusInternalServerError, constants.UNKNOWN_ERROR, utils.Null())
		return
	}

	var response []AppointmentView
	for _, b := range bookings {
		response = append(response, AppointmentView{
			ID:       b.ID,
			Date:     b.Date.Format("2006-01-02"),
			TimeSlot: b.TimeSlot,
			Status:   b.Status,
		})
	}

	utils.BuildResponse(ctx, http.StatusOK, constants.SUCCESS, gin.H{
		"items": response,
	})
}
