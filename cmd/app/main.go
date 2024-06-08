package main

import (
	"context"
	"flag"
	"time"

	"github.com/tubopo/tick-tick-ticket/internal/jira"
	"github.com/tubopo/tick-tick-ticket/internal/microsoft"
	"github.com/tubopo/tick-tick-ticket/pkg/config"
	"github.com/tubopo/tick-tick-ticket/pkg/logger"
)

func main() {
	configPath := flag.String("config", "./config.json", "path to the config file")
	date := flag.String("date", "", "The date for which to retrieve the schedule in YYYY-MM-DD format")
	ticket := flag.String("ticket", "", "The JIRA ticket to log work against")

	flag.Parse()

	log := logger.New()

	if *ticket == "" {
		log.Fatalf("You must specify a JIRA ticket")
	}

	// If no date is specified, use today's date
	if *date == "" {
		*date = time.Now().Format("2006-01-02")
	}

	start, end, err := parseDateToStartEnd(*date)
	if err != nil {
		log.Fatal("Failed to parse date: %v", err)
	}

	cfg, err := config.Load(*configPath)
	if err != nil {
		log.Fatal("Loading config failed: %v", err)
	}

	calendarService := microsoft.NewService(cfg.Microsoft, log)

	ctx := context.Background()
	events, err := calendarService.GetCalendarEvents(start, end, ctx)
	if err != nil {
		log.Error("Failed to retrieve calendar events: %v", err)
	}

	timeSpent := calendarService.CalcTotalDuration(events)

	workLogger := jira.NewService(cfg.Jira, *ticket, log)

	workLogger.LogTime(timeSpent, ctx)
}

func parseDateToStartEnd(dateStr string) (start, end time.Time, err error) {
	location, _ := time.LoadLocation("Local") // Local timezone, or provide a specific one

	start, err = time.ParseInLocation("2006-01-02", dateStr, location)
	if err != nil {
		return time.Time{}, time.Time{}, err
	}

	// Set start time at beginning of the day
	start = time.Date(start.Year(), start.Month(), start.Day(), 0, 0, 0, 0, start.Location())

	// End time is just before the start of the next day
	end = time.Date(start.Year(), start.Month(), start.Day(), 23, 59, 59, 0, start.Location())

	return start, end, nil
}
