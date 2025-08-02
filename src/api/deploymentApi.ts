// API client for the Heirloom Deployment Backend

import { Deployment } from "@/components/DeploymentTree";

// API base URL - update this to match your backend server
const API_BASE_URL = "http://localhost:8888";

// Fetch all deployments from the backend
export const fetchDeployments = async (limit: number = 10): Promise<Deployment[]> => {
  try {
    const response = await fetch(`${API_BASE_URL}/all-deployments?limit=${limit}`);
    
    if (!response.ok) {
      throw new Error(`API error: ${response.status}`);
    }
    
    const data = await response.json();
    return data as Deployment[];
  } catch (error) {
    console.error("Error fetching deployments:", error);
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