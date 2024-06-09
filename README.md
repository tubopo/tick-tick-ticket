# TickTickTicket ⌚️

[![CI [branch]](https://github.com/tubopo/tick-tick-ticket/actions/workflows/ci-branch.yml/badge.svg)](https://github.com/tubopo/tick-tick-ticket/actions/workflows/ci-branch.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/tubopo/tick-tick-ticket)](https://goreportcard.com/report/github.com/tubopo/tick-tick-ticket)
[![MIT License](https://img.shields.io/github/license/mashape/apistatus.svg?maxAge=2592000)](https://github.com/tubopo/tick-tick-ticket/blob/main/LICENSE)

TickTickTicket is a command-line tool that takes the hassle out of logging meeting hours. Say goodbye to tracking down each meeting in your calendar and typing it manually into Jira. TickTickTicket does the heavy lifting for you, automating the process with just a few simple commands.

Simplicity and efficiency are the cores of TickTickTicket. It's fast, it's fuss-free, and it's built to help you reclaim your time for what truly matters. Let's dive in!

## Features

+ Automatically pulls time spent in meetings from your calendar (Microsoft).
+ Logs those hours directly into your Jira tickets.
+ Supports authentication to keep your data secure.
+ Command-line interface for easy integration into your workflow.

## Getting Started

Before you begin, make sure you have Go installed on your system. You can download and install Go from the official website.

## Build

Clone the repository to your local machine:

```sh
git clone https://github.com/tubopo/ticktickticket.git
```

Navigate to the cloned directory and build the project:

```sh
cd cmd/tick-tick-ticket
go build -v
```

## Configuration

Create a config.json file with the following structure and provide your Microsoft calendar API keys and Jira settings:

```json
{
    "calendar": {
        "clientId": "",
        "clientSecret": "",
        "tenantId": ""
    },
    "jira": {
        "domain": "",
        "apiToken": ""
    }
}
```

> You can get your Microsoft calendar API keys from the Azure portal. Please see [documentation](https://learn.microsoft.com/en-us/entra/identity-platform/scenario-desktop-app-registration).
> To get jira personal token see details [here](https://confluence.atlassian.com/enterprise/using-personal-access-tokens-1026032365.html).

## Usage

Run TickTickTicket with the date and Jira ticket as arguments:

```sh
./tick-tick-ticket --config "/path/to/config.json" --verbose --date="2024-06-17" --ticket="JIRA-123"
```

Shorten version, using current date:

```sh
./tick-tick-ticket --ticket="JIRA-123"
```

This command will extract the time spent in meetings for the given date and log it to the specified Jira ticket.
Now, managing your meeting time is as simple as a tick, tick, and a ticket!
