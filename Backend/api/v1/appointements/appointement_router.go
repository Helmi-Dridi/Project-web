package bookingappointment

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// AppointmentRouterInit initializes the routes for booking appointments.
func AppointmentRouterInit(router *gin.RouterGroup, db *gorm.DB) {
	// Initialize database wrapper
	baseInstance := Database{DB: db}

	// Auto-migrate Appointment schema
	NewBookingAppointmentRepository(db)

	// Routes group: /appointments
	appointments := router.Group("/appointments")
	{
		appointments.POST("", baseInstance.CreateBooking)                     // POST /appointments
		appointments.GET("", baseInstance.GetUserAppointments)               // GET /appointments
		appointments.GET("/availability", baseInstance.GetAvailabilityByDate) // GET /appointments/availability
		appointments.POST("/:id/cancel", baseInstance.CancelBooking)        // PATCH /appointments/:id/cancel
		appointments.GET("/received", baseInstance.GetReceivedAppointments)
		appointments.POST("/:id/status", baseInstance.UpdateAppointmentStatus)
		// 🔥 NEW ROUTES
		appointments.GET("/:id", baseInstance.GetAppointmentByID)             // Get single appointment details
		appointments.GET("/stats", baseInstance.GetAppointmentStats)          // Admin stats
		appointments.GET("/upcoming", baseInstance.GetUpcomingAppointments)   // Upcoming appointments
		appointments.GET("/history", baseInstance.GetAppointmentHistory)      // Past appointments
		appointments.POST("/:id/remind", baseInstance.SendReminder)           // Manual reminder
		appointments.POST("/:id/reschedule", baseInstance.RescheduleBooking)  // Reschedule
	}
}
