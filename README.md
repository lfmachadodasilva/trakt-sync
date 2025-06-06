# Trakt Sync

Welcome to the Trakt Sync! This application allows you to sync your media center with Trakt, providing a seamless experience for tracking your watched items.

## Media Center Integration

- Emby ✅
- Jellyfin 🚧
- Plex 🚧

## Getting Started

To get started with the Trakt Sync application, you can use Docker to run the service. Below is a sample `docker-compose.yml` file that you can use to set up the application.

### Docker Compose Configuration

```yaml
services:
  api:
    image: lfmachadodasilva/trakt-sync:latest
    container_name: trakt-sync
    ports:
      - "3000:3000"
    volumes:
      - ./data:/app/data
    restart: unless-stopped
```

### Running the Application

1. Ensure you have Docker and Docker Compose installed on your machine.
2. Create a directory named `data`. This directory will be used to store the application's data, e.g., configuration files and database.
3. Run the following command to start the application:
   ```bash
   docker-compose up -d
   ```
4. Access the application at `http://localhost:3000`.

### Running Locally

If you prefer to run the application locally without Docker, you can follow these steps:

### Backend Setup

1. Install Go 1.24.3 on your machine.
2. Run the following commands to set up the backend:
   ```bash
   go mod tidy
   go run cmd/api/
   ```
3. Check folder http/ for `.http` files to test the API endpoints.
4. The backend will start on port `4000` by default when you run the application locally (outside of Docker).

### Frontend Setup

1. Install bun https://bun.sh/
2. Run the following commands to set up the frontend:

   ```bash
   cd cmd/web
   bun install
   bun dev
   ```

3. Access the application at `http://localhost:3000`.

## How to use

To use the Trakt Sync application, follow these steps:

1. **Configuration**:
   - Navigate to `http://localhost:3000` to set up your Trakt and media center credentials.
   - Ensure you have the necessary API keys and tokens for Trakt.
2. **Syncing**:
   - Use the sync button to start syncing your watched items between your media center and Trakt.
   - The application will fetch watched items from your media center and update them on Trakt.
3. **Webhooks**:
   - Set up webhooks to receive real-time updates from your media center.
   - The application will listen for events and sync them with Trakt automatically.

## Tech Stack

- **Backend**: Go
- **Frontend**: React with Bun
- **Styling**: Shadcn, Tailwind CSS
- **Database**: SQLite

## Contributing

We welcome contributions to the Trakt Sync application! If you have ideas for improvements or new features. please feel free to open an issue or submit a pull request.
