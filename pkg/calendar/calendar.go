package calendar

import (
	"github.com/tubopo/pkg/auth"
)

// CalendarService is an abstraction around a calendar service.
type CalendarService interface {
	GetTimeSpent(date string) (float64, error)
	LogTimeToJira(date, ticket string, logger JiraService) error
}

// JiraService is the interface needed to log work in JIRA.
type JiraService interface {
	LogWork(ticket string, timeSpent float64) error
}

// Service is a calendar service that implements CalendarService.
type Service struct {
	Authenticator auth.Authenticator
	// TODO: Other fields like an HTTP client, if needed.
}

// NewService creates a new calendar service instance.
func NewService(authenticator auth.Authenticator) CalendarService {
	return &Service{
		Authenticator: authenticator,
		// TODO: Initialize other fields.
	}
}

// GetTimeSpent gets the amount of time spent on events for a specific date.
func (s *Service) GetTimeSpent(date string) (float64, error) {
	// TODO: Authenticate if needed and fetch calendar events.
	// TODO: Calculate the total time spent on calendar events.
	return 0, nil // replace with actual time spent
}

// LogTimeToJira logs the time spent on a particular date to a JIRA ticket.
func (s *Service) LogTimeToJira(date, ticket string, jiraService JiraService) error {
	timeSpent, err := s.GetTimeSpent(date)
	if err != nil {
		return err
	}
	return jiraService.LogWork(ticket, timeSpent)
}
