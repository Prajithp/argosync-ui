# Heirloom Deployment Platform

A full-stack application for managing deployments across different environments and regions.

## Project Structure

- `backend/`: Go backend API using Echo framework
- `src/`: Frontend React application

## Backend

The backend is a Go API built with the Echo framework that provides endpoints for managing deployments:

- `POST /release`: Deploy a new version
- `POST /rollback`: Rollback to a previous version
- `GET /deployments`: Get active deployments
- `GET /history`: Get deployment history
- `GET /all-deployments`: Get all deployments in a format compatible with the frontend

### Setting Up the Backend

The backend supports both PostgreSQL and SQLite databases:

#### Using SQLite (Recommended for Testing)

SQLite is the default database and is recommended for testing:

1. Navigate to the backend directory:
   ```bash
   cd backend
   ```
2. Run the server with SQLite:
   ```bash
   go run main.go --db-type=sqlite --sqlite-path=heirloom.db --init-db
   ```

The `--init-db` flag will initialize the database with the schema and sample data.

#### Using PostgreSQL (For Production)

For production environments, you can use PostgreSQL:

1. Create a PostgreSQL database
2. Navigate to the backend directory:
   ```bash
   cd backend
   ```
3. Run the server with PostgreSQL:
   ```bash
   go run main.go --db-type=postgres --db-host=localhost --db-port=5432 --db-user=postgres --db-pass=postgres --db-name=heirloom --init-db
   ```

## Frontend

The frontend is a React application that uses the backend API to manage deployments.

### Integrating the Backend with the Frontend

1. The `src/api/deploymentApi.ts` file provides functions for interacting with the backend API:
   - `fetchDeployments()`: Fetch all deployments
   - `releaseVersion()`: Deploy a new version
   - `rollbackDeployment()`: Rollback to a previous version
   - `fetchDeploymentHistory()`: Get deployment history

2. The `src/components/DeploymentManager.tsx` component demonstrates how to use these functions with the existing UI components.

### Example Usage

```tsx
import { DeploymentManager } from "@/components/DeploymentManager";

function App() {
  return (
    <div className="App">
      <DeploymentManager />
    </div>
  );
}
```

## Data Format

The backend API returns data in the same format as the sample data used by the frontend:

```typescript
interface Deployment {
  applicationName: string;
  environment: string;
  region: string;
  version: string;
  timestamp: string;
  status: 'active' | 'inactive';
}
```

## Testing

1. Start the backend server:
   ```bash
   cd backend
   go run main.go
   ```

2. Run the test script to verify the API endpoints:
   ```bash
   cd backend
   ./test.sh
   ```

3. Start the frontend development server:
   ```bash
   npm run dev
   ```

4. Open the application in your browser and use the DeploymentForm to create new deployments.

## Notes

- The backend supports both PostgreSQL and SQLite databases:
  - SQLite is recommended for testing and development
  - PostgreSQL is recommended for production environments
- The frontend expects the backend to be running on `http://localhost:8080`. Update the `API_BASE_URL` in `src/api/deploymentApi.ts` if your backend is running on a different URL.
- Sample data is automatically loaded when using the `--init-db` flag
