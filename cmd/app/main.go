package main

import (
	"flag"
	"github.com/tubopo/tick-tick-ticket/internal/jira"
	"github.com/tubopo/tick-tick-ticket/internal/microsoft"
	"github.com/tubopo/tick-tick-ticket/pkg/config"
	"log"
)

func main() {
	configPath := flag.String("config", "./config.json", "path to the config file")
	date := flag.String("date", "", "The date for which to retrieve the schedule in YYYY-MM-DD format")
	ticket := flag.String("ticket", "", "The JIRA ticket to log work against")

	flag.Parse()

	if *date == "" || *ticket == "" {
		log.Fatal("You must specify both a date and a JIRA ticket")
	}

	cfg, err := config.Load(*configPath)
	if err != nil {
		log.Fatalf("Loading config failed: %v", err)
	}

	microsoftAuthenticator := &microsoft.Authenticator{Cfg: &cfg.Microsoft}
	jiraAuthenticator := &jira.Authenticator{Cfg: cfg.Jira}

	// Setup the service clients
	calendarService := microsoft.NewService(microsoftAuthenticator)
	jiraService := jira.NewService(jiraAuthenticator)

	// Perform the actions
	err = calendarService.GetCalendarEvents(*date, *ticket, jiraService)
	if err != nil {
		log.Fatalf("Failed to log time to JIRA: %v", err)
	}
}
