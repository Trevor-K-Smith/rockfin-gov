# rockfin-gov

## Overview

**rockfin-gov** is a GoLang-based project aimed at:
1. Fetching federal opportunities from **[sam.gov](https://sam.gov/)** using their public API.
2. Storing that data into a **PostgreSQL** database.
3. Leveraging AI (summaries and action recommendations) on top of the stored data.
4. Integrating new opportunity sources in the future, such as state-specific sites.

This README describes the project structure, key components, and how everything ties together.

---

## File Structure

rockfin-gov
├── README.md                 # Project documentation
├── config.yaml               # Global configuration file
├── go.mod                    # Go module definition
├── go.sum                    # Go module checksums
├── cmd                       # Executable entry points
│   └── rockfin-gov
│       └── main.go           # Main application entry point
├── internal                  # Private packages (internal to this repo)
│   ├── aggregator            # Module for fetching and aggregating data
│   │   ├── aggregator.go
│   │   ├── federal           # Federal-level data fetchers
│   │   │   └── samgovclient.go
│   │   └── states            # State-specific data fetchers
│   │       └── examplestateclient.go
│   ├── db                    # Database-related logic
│   │   └── dbmanager.go
│   ├── ai                    # AI or NLP functionalities
│   │   └── summarizer.go
│   └── routes                # Web API routes / endpoints
│       └── opportunities.go
└── test                      # Test files (unit/integration)
    └── aggregator_test.go

