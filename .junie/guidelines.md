# UMD Developer Guidelines

This document provides essential information for developers working on the Universal Media Downloader (UMD) project.

## Project Overview

UMD (Universal Media Downloader) is an application for downloading media files from popular websites like Coomer, Fapello, Imaglr, Kemono, Reddit, and RedGifs. The application has two interfaces:

- **CLI**: Command-line interface for automated downloads
- **GUI**: Graphical user interface for interactive downloads

## Tech Stack

### Backend
- **Go**: Primary language for both CLI and GUI backends
- **Wails**: Framework for building desktop applications with Go and web technologies

### Frontend
- **TypeScript**: Programming language for the frontend
- **React**: UI library for building the interface
- **Vite**: Build tool for the frontend

### Build Tools
- **Task**: Task runner used for build automation
- **PNPM**: Package manager for Node.js dependencies

## Project Structure

The project is organized into the following key directories:

```
umd/
├── cli/            # CLI implementation
├── gui/            # GUI implementation
│   └── frontend/   # Web-based frontend (React/TypeScript)
├── shared/         # Common code used by both CLI and GUI
├── scripts/        # Utility scripts
├── assets/         # Static assets like images
└── build/          # Build output (created during build)
```

### Key Components

- **cli/**: Contains the command-line interface implementation
  - `cli.go`: Main CLI interface
  - `umd.go`: Core functionality for the CLI
  - `utils.go`: Utility functions

- **gui/**: Contains the graphical user interface implementation
  - `app.go`: Application structure
  - `binding.go`: Bindings between Go and the frontend
  - `events.go`: Event management
  - `main.go`: Entry point for the GUI application

- **shared/**: Contains code shared between CLI and GUI
  - `cache.go`: Caching functionality
  - `download.go`: Core download functionality
  - `models.go`: Data structures
  - `queue.go`: Download queue management
  - `telemetry.go`: Telemetry functionality
  - `utils.go`: Common utility functions

### Build Commands

Build the CLI version:
```bash
task cli os=<operating-system> arch=<architecture>
```

Build the GUI version:
```bash
task gui os=<operating-system> arch=<architecture>
```

Where:
- `<operating-system>`: `windows`, `darwin` (macOS), or `linux`
- `<architecture>`: `amd64` or `arm64`

Example:
```bash
task gui os=darwin arch=arm64
```

Build and package all versions:
```bash
task package oses="windows darwin linux" archs="amd64 arm64"
```

Clean build artifacts:
```bash
task clean
```

## Running the Application

### CLI

Run the CLI version:
```bash
./build/umd-dl -d <download-directory> <url>
```

Options:
- `-d`: Download directory (optional, default is current directory)
- `-l`: Maximum number of files to download (optional, default is 99,999)
- `--no-telemetry`: Disable telemetry (optional)

### GUI

Run the GUI version:
- Windows: Double-click `build/umd.exe`
- macOS: Double-click `build/UMD.app`
- Linux: Run `./build/umd`

## Development Workflow

### Making Changes to the CLI

1. Edit files in the `cli/` directory
2. Build using `task cli os=<os> arch=<arch>`
3. Test your changes by running the built binary

### Making Changes to the GUI

1. For backend changes, edit files in the `gui/` directory
2. For frontend changes, edit files in the `gui/frontend/src/` directory
3. Build using `task gui os=<os> arch=<arch>`
4. Test your changes by running the built application

### Making Changes to Shared Code

1. Edit files in the `shared/` directory
2. Build both CLI and GUI to test your changes
3. Ensure your changes work correctly in both interfaces

## Best Practices

### Code Organization

- Keep interface-specific code in the appropriate directory (`cli/` or `gui/`)
- Place common functionality in the `shared/` directory
- Use Go modules for dependency management

### Coding Standards

- Follow Go's [Effective Go](https://golang.org/doc/effective_go) guidelines
- Use TypeScript's strict mode for frontend code
- Format code using appropriate tools (Go fmt, Prettier)

### Feature Development

- Implement features in the shared module when possible
- Ensure features work consistently between CLI and GUI
- Consider both desktop use cases for GUI features