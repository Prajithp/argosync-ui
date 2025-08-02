// Optimized sample data for 1000+ services CI/CD platform
// This represents a realistic enterprise setup with multiple teams and services

export interface Application {
  id: number;
  name: string;
  team: string;
  description?: string;
}

export interface Environment {
  id: number;
  name: string;
  priority: number;
}

export interface Region {
  id: number;
  code: string;
  name: string;
  continent: string;
}

export interface OptimizedDeployment {
  id: number;
  applicationId: number;
  environmentId: number;
  regionId: number;
  version: string;
  status: 'active' | 'inactive' | 'failed';
  deployedBy: string;
  deployedAt: string;
  metadata?: {
    buildId?: string;
    commitHash?: string;
    duration?: number;
  };
}

// Reference data (normalized)
export const environments: Environment[] = [
  { id: 1, name: 'development', priority: 1 },
  { id: 2, name: 'staging', priority: 2 },
  { id: 3, name: 'production', priority: 3 },
];

export const regions: Region[] = [
  { id: 1, code: 'us-east-1', name: 'US East (N. Virginia)', continent: 'North America' },
  { id: 2, code: 'us-west-2', name: 'US West (Oregon)', continent: 'North America' },
  { id: 3, code: 'eu-west-1', name: 'Europe (Ireland)', continent: 'Europe' },
  { id: 4, code: 'eu-central-1', name: 'Europe (Frankfurt)', continent: 'Europe' },
  { id: 5, code: 'ap-southeast-1', name: 'Asia Pacific (Singapore)', continent: 'Asia' },
  { id: 6, code: 'ap-northeast-1', name: 'Asia Pacific (Tokyo)', continent: 'Asia' },
];

// Sample applications representing different teams and service types
export const applications: Application[] = [
  // Auth Team
  { id: 1, name: 'auth-service', team: 'auth', description: 'User authentication and authorization' },
  { id: 2, name: 'user-management', team: 'auth', description: 'User profile and settings management' },
  { id: 3, name: 'identity-provider', team: 'auth', description: 'SSO and identity management' },
  
  // Payment Team
  { id: 4, name: 'payment-gateway', team: 'payment', description: 'Payment processing service' },
  { id: 5, name: 'billing-service', team: 'payment', description: 'Subscription and billing management' },
  { id: 6, name: 'fraud-detection', team: 'payment', description: 'Real-time fraud detection' },
  
  // Platform Team
  { id: 7, name: 'api-gateway', team: 'platform', description: 'Main API gateway and routing' },
  { id: 8, name: 'notification-service', team: 'platform', description: 'Email, SMS, and push notifications' },
  { id: 9, name: 'logging-service', team: 'platform', description: 'Centralized logging and monitoring' },
  { id: 10, name: 'config-service', team: 'platform', description: 'Configuration management' },
  
  // Product Team
  { id: 11, name: 'product-catalog', team: 'product', description: 'Product information and catalog' },
  { id: 12, name: 'inventory-service', team: 'product', description: 'Inventory tracking and management' },
  { id: 13, name: 'recommendation-engine', team: 'product', description: 'AI-powered product recommendations' },
  
  // Analytics Team
  { id: 14, name: 'analytics-collector', team: 'analytics', description: 'Event collection and processing' },
  { id: 15, name: 'reporting-service', team: 'analytics', description: 'Business intelligence and reports' },
  { id: 16, name: 'data-warehouse', team: 'analytics', description: 'Data warehousing and ETL' },
  
  // Mobile Team
  { id: 17, name: 'mobile-api', team: 'mobile', description: 'Mobile-specific API endpoints' },
  { id: 18, name: 'push-service', team: 'mobile', description: 'Mobile push notification service' },
  
  // DevOps Team
  { id: 19, name: 'deployment-service', team: 'devops', description: 'CI/CD deployment orchestration' },
  { id: 20, name: 'monitoring-service', team: 'devops', description: 'Infrastructure monitoring' },
];

// Sample deployments with realistic patterns
export const optimizedDeployments: OptimizedDeployment[] = [
  // Auth services - deployed across all environments
  { id: 1, applicationId: 1, environmentId: 3, regionId: 1, version: 'v2.4.1', status: 'active', deployedBy: 'alex.smith', deployedAt: '2024-01-15T10:30:00Z', metadata: { buildId: 'build-1234', commitHash: 'abc123' } },
  { id: 2, applicationId: 1, environmentId: 3, regionId: 3, version: 'v2.4.1', status: 'active', deployedBy: 'alex.smith', deployedAt: '2024-01-15T10:35:00Z', metadata: { buildId: 'build-1234', commitHash: 'abc123' } },
  { id: 3, applicationId: 1, environmentId: 2, regionId: 1, version: 'v2.5.0-rc1', status: 'active', deployedBy: 'sarah.jones', deployedAt: '2024-01-20T14:20:00Z', metadata: { buildId: 'build-1267', commitHash: 'def456' } },
  { id: 4, applicationId: 1, environmentId: 1, regionId: 1, version: 'v2.5.0-beta2', status: 'active', deployedBy: 'mike.wilson', deployedAt: '2024-01-22T09:15:00Z', metadata: { buildId: 'build-1289', commitHash: 'ghi789' } },
  
  // Payment services - critical production deployments
  { id: 5, applicationId: 4, environmentId: 3, regionId: 1, version: 'v1.8.3', status: 'active', deployedBy: 'lisa.chen', deployedAt: '2024-01-10T08:00:00Z', metadata: { buildId: 'build-1156', commitHash: 'xyz987' } },
  { id: 6, applicationId: 4, environmentId: 3, regionId: 3, version: 'v1.8.3', status: 'active', deployedBy: 'lisa.chen', deployedAt: '2024-01-10T08:05:00Z', metadata: { buildId: 'build-1156', commitHash: 'xyz987' } },
  { id: 7, applicationId: 4, environmentId: 3, regionId: 5, version: 'v1.8.3', status: 'active', deployedBy: 'lisa.chen', deployedAt: '2024-01-10T08:10:00Z', metadata: { buildId: 'build-1156', commitHash: 'xyz987' } },
  
  // Previous versions for rollback scenarios
  { id: 8, applicationId: 4, environmentId: 3, regionId: 1, version: 'v1.8.2', status: 'inactive', deployedBy: 'lisa.chen', deployedAt: '2024-01-08T16:30:00Z', metadata: { buildId: 'build-1145', commitHash: 'prev123' } },
  { id: 9, applicationId: 4, environmentId: 3, regionId: 1, version: 'v1.8.1', status: 'inactive', deployedBy: 'tom.brown', deployedAt: '2024-01-05T11:45:00Z', metadata: { buildId: 'build-1134', commitHash: 'prev456' } },
  
  // Platform services
  { id: 10, applicationId: 7, environmentId: 3, regionId: 1, version: 'v3.1.0', status: 'active', deployedBy: 'david.kim', deployedAt: '2024-01-18T13:20:00Z', metadata: { buildId: 'build-1245', commitHash: 'api789' } },
  { id: 11, applicationId: 7, environmentId: 3, regionId: 2, version: 'v3.1.0', status: 'active', deployedBy: 'david.kim', deployedAt: '2024-01-18T13:25:00Z', metadata: { buildId: 'build-1245', commitHash: 'api789' } },
  { id: 12, applicationId: 8, environmentId: 3, regionId: 1, version: 'v2.2.1', status: 'active', deployedBy: 'emma.davis', deployedAt: '2024-01-16T15:10:00Z', metadata: { buildId: 'build-1223', commitHash: 'notif456' } },
  
  // Analytics services
  { id: 13, applicationId: 14, environmentId: 3, regionId: 1, version: 'v1.5.2', status: 'active', deployedBy: 'ryan.miller', deployedAt: '2024-01-19T11:30:00Z', metadata: { buildId: 'build-1256', commitHash: 'analytics123' } },
  { id: 14, applicationId: 15, environmentId: 2, regionId: 1, version: 'v1.6.0-beta1', status: 'active', deployedBy: 'sophia.garcia', deployedAt: '2024-01-21T10:00:00Z', metadata: { buildId: 'build-1278', commitHash: 'reports789' } },
  
  // Mobile services
  { id: 15, applicationId: 17, environmentId: 3, regionId: 1, version: 'v2.0.3', status: 'active', deployedBy: 'james.taylor', deployedAt: '2024-01-17T09:45:00Z', metadata: { buildId: 'build-1234', commitHash: 'mobile456' } },
  { id: 16, applicationId: 18, environmentId: 3, regionId: 1, version: 'v1.3.1', status: 'active', deployedBy: 'olivia.anderson', deployedAt: '2024-01-14T14:30:00Z', metadata: { buildId: 'build-1201', commitHash: 'push789' } },
];

// Helper functions for data transformation
export const getApplicationById = (id: number): Application | undefined => 
  applications.find(app => app.id === id);

export const getEnvironmentById = (id: number): Environment | undefined => 
  environments.find(env => env.id === id);

export const getRegionById = (id: number): Region | undefined => 
  regions.find(region => region.id === id);

// Transform optimized data to legacy format for UI compatibility
export const transformToLegacyFormat = () => {
  return optimizedDeployments
    .filter(deployment => deployment.status !== 'failed') // Filter out failed deployments for legacy compatibility
    .map(deployment => {
      const app = getApplicationById(deployment.applicationId);
      const env = getEnvironmentById(deployment.environmentId);
      const region = getRegionById(deployment.regionId);
      
      return {
        applicationName: app?.name || 'unknown',
        environment: env?.name || 'unknown',
        region: region?.code || 'unknown',
        version: deployment.version,
        status: deployment.status as 'active' | 'inactive', // Type assertion for legacy compatibility
        timestamp: deployment.deployedAt,
      };
    });
};

// Performance optimized queries (examples)
export const getActiveDeploymentsByTeam = (team: string) => {
  const teamApps = applications.filter(app => app.team === team);
  const teamAppIds = teamApps.map(app => app.id);
  
  return optimizedDeployments.filter(deployment => 
    teamAppIds.includes(deployment.applicationId) && deployment.status === 'active'
  );
};

export const getDeploymentHistory = (applicationId: number, environmentId: number, regionId: number) => {
  return optimizedDeployments
    .filter(deployment => 
      deployment.applicationId === applicationId &&
      deployment.environmentId === environmentId &&
      deployment.regionId === regionId
    )
    .sort((a, b) => new Date(b.deployedAt).getTime() - new Date(a.deployedAt).getTime());
};