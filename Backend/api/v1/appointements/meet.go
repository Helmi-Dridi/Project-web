package bookingappointment

import (
	"context"
	"fmt"
	"os"
	"time"

	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
)

func CreateGoogleMeetEvent(summary, description string, start, end time.Time) (string, error) {
	ctx := context.Background()

	// Lire le fichier credentials.json
	creds, err := os.ReadFile("C:/Users/HELMI/go/src/github.com/App/api/v1/appointements/appointmentapp-462814-5b55756fb624.json")
	if err != nil {
		return "", err
	}

	conf, err := google.JWTConfigFromJSON(creds, calendar.CalendarScope)
	if err != nil {
		return "", err
	}

	client := conf.Client(ctx)
	service, err := calendar.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return "", err
	}

	event := &calendar.Event{
		Summary:     summary,
		Description: description,
		Start: &calendar.EventDateTime{
			DateTime: start.Format(time.RFC3339),
			TimeZone: "Europe/Paris", // Modifie selon ton fuseau horaire
		},
		End: &calendar.EventDateTime{
			DateTime: end.Format(time.RFC3339),
			TimeZone: "Europe/Paris",
		},
		ConferenceData: &calendar.ConferenceData{
			CreateRequest: &calendar.CreateConferenceRequest{
				RequestId: fmt.Sprintf("appt-%d", time.Now().UnixNano()),
				ConferenceSolutionKey: &calendar.ConferenceSolutionKey{
					Type: "hangoutsMeet",
				},
			},
		},
		// 🔥 Attendees removed to avoid 403
	}

	inserted, err := service.Events.Insert("primary", event).ConferenceDataVersion(1).Do()
	if err != nil {
		return "", err
	}

	for _, ep := range inserted.ConferenceData.EntryPoints {
		if ep.EntryPointType == "video" {
			return ep.Uri, nil
		}
	}

	return "", fmt.Errorf("aucun lien Meet généré")
}
