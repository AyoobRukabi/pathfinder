<p align="center">
<img src="https://img.shields.io/badge/kood/Sisu-brightgreen?logo=gitea&logoColor=white&labelColor=8A2BE2">
<img src="https://img.shields.io/badge/1.26.0-00ADD8?logo=go&logoColor=white&labelColor=gray">
<!-- <img src="https://img.shields.io/badge/platforms-Linux%20|%20macOS-green.svg"> -->
<img src="https://img.shields.io/badge/license-MIT-green.svg">
</p>

# Pathfinder

> [!NOTE]
> Project Context: Developed in 2026 as part of the Go Module within a project-based coding program.

---

> [!WARNING]
> The rest is a boilerplate.


## Task learning objectives:
- bla bla bla

## Key learnings and Results:
- some learning

## Key Features

### 1. Some feature
The application

## Installation

### Prerequisites

* **Go:** Version 1.22 or higher.

## Quick Start

### 1. Clone the repo

1. Navigate to the project repo.
```bash
git clone <GITEA_URL>
cd pathfinder
```
2. Build the application
```bash
go build -o pathfinder cmd/pathfinder/main.go
```
3. Run the application:
```bash
./pathfinder
```

> [!WARNING]
> Boilerplate Ends.

## Structure

```text
pathfinder/
├── cmd/
│   └── pathfinder/
│       └── main.go             # Application entry point (Wires dependencies & starts App)
├── config/
│   └── local/
│       └── local.json          # Configuration for local environment
├── internal/
│   ├── app/                    # Composition Root (Initializes Core, Repository, Web)
│   ├── config/                 # Configuration structs and parsing logic
│   ├── domain/                 # Core Business Entities (structs that we need, e.g. stations, graph map, trains, etc.)
│   ├── lib/                    # Any useful functions or helpers not related to some of the layers.
│   │   └── e/                  # Error wrapping utilities
│   ├── storage/                # Data Access Layer (e.g. we store maps in files and memory)
│   │   └── memory/
│   └── usecase/
│       └── dispetcher/         # Business Logic (algorithm to find the shortest root, optimisation, and so on)
├── pkg/                        # Reusable Library Code (No domain dependencies)
│   ├── cache/                  # Thread-safe Cache with Janitor
│   ├── httpclient/             # Resilient HTTP Client wrapper (not needed, just as an example)
│   └── logger/                 # Structured Logger setup
├── go.mod                      # Go Module definitions
├── CONTRIBUTING.md             # Team Collaboration & Git Workflow Guide
├── LICENSE.txt                 # License
├── README.md                   # Project Documentation
└── TODO.md                     # Todo list to mark some ideas, tasks and so on. Not required since we can mark items as issues in Git.
```

---

> [!CAUTION]
> ## ⚖️ Academic Integrity Disclaimer
> This project was submitted as part of the Go Module within a project-based coding program. It is intended for portfolio purposes only.
>
> If you are a student currently enrolled in a similar course, please be aware that using any part of this code may violate your institution's academic honesty policy.
>
> Use this as a reference, but write your own code!

## License

This project is licensed under the [MIT License](https://github.com/KorbenSweetheart/car-catalog/blob/main/LICENSE.txt).