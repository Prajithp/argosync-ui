import { Deployment } from "@/components/DeploymentTree";

export const sampleDeployments: Deployment[] = [
  {
    applicationName: "user-service",
    environment: "production",
    region: "us-east-1",
    version: "v2.1.0",
    timestamp: "2024-01-15T10:30:00Z",
    status: "active"
  },
  {
    applicationName: "user-service",
    environment: "production",
    region: "us-east-1",
    version: "v2.0.5",
    timestamp: "2024-01-14T15:20:00Z",
    status: "inactive"
  },
  {
    applicationName: "user-service",
    environment: "staging",
    region: "us-east-1",
    version: "v2.2.0-beta",
    timestamp: "2024-01-16T09:15:00Z",
    status: "active"
  },
  {
    applicationName: "user-service",
    environment: "development",
    region: "us-west-2",
    version: "v2.3.0-alpha",
    timestamp: "2024-01-17T11:45:00Z",
    status: "active"
  },
  {
    applicationName: "api-gateway",
    environment: "production",
    region: "us-east-1",
    version: "v1.5.2",
    timestamp: "2024-01-13T14:00:00Z",
    status: "active"
  },
  {
    applicationName: "api-gateway",
    environment: "production",
    region: "eu-west-1",
    version: "v1.5.1",
    timestamp: "2024-01-12T08:30:00Z",
    status: "active"
  },
  {
    applicationName: "api-gateway",
    environment: "staging",
    region: "us-east-1",
    version: "v1.6.0-rc1",
    timestamp: "2024-01-15T16:20:00Z",
    status: "active"
  },
  {
    applicationName: "notification-service",
    environment: "production",
    region: "ap-southeast-1",
    version: "v3.0.1",
    timestamp: "2024-01-16T12:10:00Z",
    status: "active"
  },
  {
    applicationName: "notification-service",
    environment: "production",
    region: "ap-southeast-1",
    version: "v3.0.0",
    timestamp: "2024-01-15T18:45:00Z",
    status: "inactive"
  },
  {
    applicationName: "notification-service",
    environment: "staging",
    region: "us-west-2",
    version: "v3.1.0-beta",
    timestamp: "2024-01-17T13:25:00Z",
    status: "active"
  }
];