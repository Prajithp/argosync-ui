// API client for the Heirloom Deployment Backend

import { Deployment } from "@/components/DeploymentTree";

// API base URL - update this to match your backend server
const API_BASE_URL = "/api/v1";

// Pagination response interface
export interface PaginatedResponse<T> {
  deployments: T[];
  pagination: {
    page: number;
    pageSize: number;
    totalCount: number;
    totalPages: number;
  };
}

// Fetch all deployments from the backend with pagination
export const fetchDeployments = async (
  page: number = 1,
  pageSize: number = 10
): Promise<PaginatedResponse<Deployment>> => {
  try {
    const response = await fetch(
      `${API_BASE_URL}/all-deployments?page=${page}&pageSize=${pageSize}`
    );
    
    if (!response.ok) {
      throw new Error(`API error: ${response.status}`);
    }
    
    const data = await response.json();
    return data as PaginatedResponse<Deployment>;
  } catch (error) {
    console.error("Error fetching deployments:", error);
    throw error;
  }
};

// Legacy function for backward compatibility
export const fetchAllDeployments = async (limit: number = 10): Promise<Deployment[]> => {
  try {
    const response = await fetchDeployments(1, limit);
    return response.deployments;
  } catch (error) {
    console.error("Error fetching all deployments:", error);
    throw error;
  }
};

// Release a new version
export const releaseVersion = async (
  applicationName: string,
  environment: string,
  region: string,
  version: string,
  deployedBy: string = "system"
): Promise<any> => {
  try {
    const response = await fetch(`${API_BASE_URL}/release`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({
        application: applicationName,
        environment,
        region,
        version,
        deployed_by: deployedBy,
      }),
    });
    
    if (!response.ok) {
      throw new Error(`API error: ${response.status}`);
    }
    
    return await response.json();
  } catch (error) {
    console.error("Error releasing version:", error);
    throw error;
  }
};

// Rollback to a previous version
export const rollbackDeployment = async (
  applicationName: string,
  environment: string,
  region: string,
  version?: string,
  deployedBy: string = "system"
): Promise<any> => {
  try {
    const response = await fetch(`${API_BASE_URL}/rollback`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({
        application: applicationName,
        environment,
        region,
        version,
        deployed_by: deployedBy,
      }),
    });
    
    if (!response.ok) {
      throw new Error(`API error: ${response.status}`);
    }
    
    return await response.json();
  } catch (error) {
    console.error("Error rolling back deployment:", error);
    throw error;
  }
};

// Get deployment history
export const fetchDeploymentHistory = async (
  applicationName: string,
  environment: string,
  region: string
): Promise<any[]> => {
  try {
    const response = await fetch(
      `${API_BASE_URL}/history?application=${applicationName}&environment=${environment}&region=${region}`
    );
    
    if (!response.ok) {
      throw new Error(`API error: ${response.status}`);
    }
    
    return await response.json();
  } catch (error) {
    console.error("Error fetching deployment history:", error);
    throw error;
  }
};

// New hierarchical data loading functions

// Interface for Application
export interface Application {
  id?: number;
  ID?: number;
  name: string;
  created_at?: string;
  updated_at?: string;
  CreatedAt?: string;
  UpdatedAt?: string;
}

// Interface for Region
export interface Region {
  id?: number;
  ID?: number;
  code: string;
  name: string;
  created_at?: string;
  updated_at?: string;
  CreatedAt?: string;
  UpdatedAt?: string;
}

// Interface for Environment
export interface Environment {
  id?: number;
  ID?: number;
  name: string;
  created_at?: string;
  updated_at?: string;
  CreatedAt?: string;
  UpdatedAt?: string;
}

// Interface for Version
export interface Version {
  id?: number;
  ID?: number;
  version: string;
  status?: string;
  Status?: string;
  created_at?: string;
  updated_at?: string;
  CreatedAt?: string;
  UpdatedAt?: string;
}

// Fetch all applications
export const fetchAllApplications = async (): Promise<Application[]> => {
  try {
    const response = await fetch(`${API_BASE_URL}/applications`);
    
    if (!response.ok) {
      throw new Error(`API error: ${response.status}`);
    }
    
    const data = await response.json();
    
    // Log the raw data to help debug
    console.log("Raw applications data:", data);
    
    // Convert IDs to numbers when possible, but don't filter out any applications
    return data.map((app: any) => ({
      ...app,
      // Try to parse the ID as a number, but keep the original value if parsing fails
      id: app.id !== undefined && app.id !== null && !isNaN(Number(app.id))
          ? Number(app.id)
          : app.id
    }));
  } catch (error) {
    console.error("Error fetching applications:", error);
    throw error;
  }
};

// Fetch regions for an application
export const fetchRegionsForApplication = async (appId: number): Promise<Region[]> => {
  try {
    // Validate appId before making the request
    if (!appId || isNaN(appId)) {
      throw new Error("Invalid application ID");
    }
    
    const response = await fetch(`${API_BASE_URL}/applications/${appId}/regions`);
    
    if (!response.ok) {
      throw new Error(`API error: ${response.status}`);
    }
    
    const data = await response.json();
    
    // Ensure region IDs are valid numbers
    return data.map((region: any) => ({
      ...region,
      id: region.id && !isNaN(parseInt(region.id)) ? parseInt(region.id) : region.id
    }));
  } catch (error) {
    console.error("Error fetching regions for application:", error);
    throw error;
  }
};

// Fetch environments for an application and region
export const fetchEnvironmentsForApplicationAndRegion = async (
  appId: number,
  regionId: number
): Promise<Environment[]> => {
  try {
    // Validate IDs before making the request
    if (!appId || isNaN(appId)) {
      throw new Error("Invalid application ID");
    }
    if (!regionId || isNaN(regionId)) {
      throw new Error("Invalid region ID");
    }
    
    const response = await fetch(
      `${API_BASE_URL}/applications/${appId}/regions/${regionId}/environments`
    );
    
    if (!response.ok) {
      throw new Error(`API error: ${response.status}`);
    }
    
    const data = await response.json();
    
    // Ensure environment IDs are valid numbers
    return data.map((env: any) => ({
      ...env,
      id: env.id && !isNaN(parseInt(env.id)) ? parseInt(env.id) : env.id
    }));
  } catch (error) {
    console.error("Error fetching environments for application and region:", error);
    throw error;
  }
};

// Fetch versions for an application, environment, and region
export const fetchVersionsForApplicationEnvironmentRegion = async (
  appId: number,
  envId: number,
  regionId: number
): Promise<Version[]> => {
  try {
    // Validate IDs before making the request
    if (!appId || isNaN(appId)) {
      throw new Error("Invalid application ID");
    }
    if (!envId || isNaN(envId)) {
      throw new Error("Invalid environment ID");
    }
    if (!regionId || isNaN(regionId)) {
      throw new Error("Invalid region ID");
    }
    
    const response = await fetch(
      `${API_BASE_URL}/applications/${appId}/environments/${envId}/regions/${regionId}/versions`
    );
    
    if (!response.ok) {
      throw new Error(`API error: ${response.status}`);
    }
    
    const data = await response.json();
    
    // Ensure version IDs are valid numbers
    return data.map((version: any) => ({
      ...version,
      id: version.id && !isNaN(parseInt(version.id)) ? parseInt(version.id) : version.id
    }));
  } catch (error) {
    console.error("Error fetching versions for application, environment, and region:", error);
    throw error;
  }
};
