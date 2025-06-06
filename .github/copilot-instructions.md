# Copilot Instructions

## Overview

This application is a homelab project designed to synchronize Trakt with media centers like Emby, Plex, and Jellyfin. It consists of a Go-based backend and a React-based frontend, leveraging modern tools and libraries for efficient development and deployment.

---

## Backend (`cmd/api`)

- **Language & Version**: Go (v1.24.3)
- **Purpose**: Handles API endpoints for syncing Trakt with media centers.
- **Key Features**:
  - **Cron Jobs**: Uses `github.com/robfig/cron/v3` to schedule and manage synchronization tasks.
  - **CORS Middleware**: Ensures cross-origin requests are handled properly.
  - **Database**: Likely uses SQLite for storing configurations and sync data (`internal/database`).
  - **API Endpoints**:
    - `/config`: Manages configuration settings.
    - `/emby/`: Handles Emby-specific operations.
    - `/trakt/`: Handles Trakt-specific operations.
    - `/sync`: Triggers synchronization tasks.
  - **Configuration Management**: `internal/config` contains logic for managing Trakt and media center configurations.
  - **Static File Serving**: Serves React-built static files from the `static` directory.

---

## Frontend (`cmd/web`)

- **Frameworks & Tools**:
  - **React**: For building the user interface.
  - **Bun**: For faster dependency management and builds.
  - **shadcn**: For UI components and styling.
  - **framer-motion**: For animations (e.g., fade-in effects).
- **Key Features**:
  - **API Integration**: Fetches and updates configurations via API calls (`config/fetch.ts`).
  - **Dynamic UI**: Uses `shadcn` components for a modern and responsive interface.
  - **Animations**: Adds visual enhancements using `framer-motion`.
  - **Error Handling**: Ensures API calls return promises instead of throwing errors.

---

## Internal Logic (`internal/`)

- **Configuration (`internal/config`)**:
  - Manages Trakt and media center configurations.
  - Includes a `CronManager` for scheduling sync jobs dynamically.
- **Database (`internal/database`)**:
  - Handles SQLite database connections and operations.
- **Media Center Integration**:
  - **Emby (`internal/emby`)**: Functions for fetching users, items, and marking items as watched.
  - **Trakt (`internal/trakt`)**: Functions for authentication, fetching watched items, and marking items as watched.
- **Utilities (`internal/utils`)**:
  - Common helper functions for authentication, requests, and data processing.

---

## Deployment

- **Docker**:
  - Multi-stage Dockerfile for building and running the application.
  - Uses `bun` for building the React frontend and `golang` for building the backend.
  - Final runtime image is based on `alpine` for minimal size.
  - Exposes port `3000` for production and `4000` for local development.

---

## Copilot Instructions

1. **Backend Development**:

   - Use Go 1.24.3 for development.
   - Add new API endpoints in `cmd/api` and corresponding logic in `internal/`.
   - Use `internal/config` for managing configurations and `internal/database` for database interactions.
   - Test cron jobs using `github.com/robfig/cron/v3`.

2. **Frontend Development**:

   - Use React with Bun for dependency management and builds.
   - Add new UI components in `cmd/web/src/components`.
   - Use `shadcn` for consistent styling and `framer-motion` for animations.
   - Update API calls in `cmd/web/src/config/fetch.ts` to align with backend changes.

3. **Deployment**:

   - Build the Docker image using the provided Dockerfile.
   - Set `NODE_ENV` to `production` for production builds and `development` for local testing.
   - Use port `3000` for production and `4000` for local development.

4. **Testing**:

   - Test API endpoints using the `.http` files in the `http/` directory.
   - Ensure frontend and backend integration works seamlessly.

5. **Future Enhancements**:
   - Add support for Plex and Jellyfin in `internal/`.
   - Improve logging and monitoring for cron jobs and API calls.
   - Enhance the frontend UI/UX with additional animations and error handling.
