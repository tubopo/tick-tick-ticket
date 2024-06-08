package main

import (
	"context"
	"flag"

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

	if *date == "" || *ticket == "" {
		log.Fatal("You must specify both a date and a JIRA ticket")
	}

	cfg, err := config.Load(*configPath)
	if err != nil {
		log.Fatal("Loading config failed: %v", err)
	}

	calendarService := microsoft.NewService(cfg.Microsoft, log)

	ctx := context.Background()
	events, err := calendarService.GetCalendarEvents(*date, ctx)
	if err != nil {
		log.Error("Failed to retrieve calendar events: %v", err)
	}

	timeSpent := calendarService.CalcTotalDuration(events)

	workLogger := jira.NewService(cfg.Jira, log)

	workLogger.LogTime(*ticket, timeSpent, ctx)
}
