import { useState, useEffect } from "react";
import { DeploymentTree, Deployment } from "@/components/DeploymentTree";
import { DeploymentForm } from "@/components/DeploymentForm";
import { fetchDeployments, releaseVersion, rollbackDeployment } from "@/api/deploymentApi";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Alert, AlertDescription, AlertTitle } from "@/components/ui/alert";
import { AlertCircle, CheckCircle2 } from "lucide-react";

export const DeploymentManager = () => {
  const [deployments, setDeployments] = useState<Deployment[]>([]);
  const [loading, setLoading] = useState<boolean>(true);
  const [error, setError] = useState<string | null>(null);
  const [success, setSuccess] = useState<string | null>(null);

  // Fetch deployments on component mount
  useEffect(() => {
    loadDeployments();
  }, []);

  // Load deployments from the API
  const loadDeployments = async () => {
    try {
      setLoading(true);
      setError(null);
      const data = await fetchDeployments();
      setDeployments(data);
    } catch (err) {
      setError("Failed to load deployments. Please try again later.");
      console.error(err);
    } finally {
      setLoading(false);
    }
  };

  // Handle deployment form submission
  const handleDeploy = async (formData: {
    applicationName: string;
    environment: string;
    region: string;
    version: string;
    deployedBy: string;
  }) => {
    try {
      setLoading(true);
      setError(null);
      setSuccess(null);
      
      await releaseVersion(
        formData.applicationName,
        formData.environment,
        formData.region,
        formData.version,
        formData.deployedBy
      );
      
      setSuccess(`Successfully deployed ${formData.applicationName} version ${formData.version} to ${formData.environment} in ${formData.region}`);
      
      // Reload deployments to show the new deployment
      await loadDeployments();
    } catch (err) {
      setError("Failed to deploy. Please check your inputs and try again.");
      console.error(err);
    } finally {
      setLoading(false);
    }
  };

  // Handle rollback
  const handleRollback = async (deployment: Deployment) => {
    try {
      setLoading(true);
      setError(null);
      setSuccess(null);
      
      // For rollbacks, we'll use the current user or system as the deployedBy value
      // You might want to add a confirmation dialog that captures this information
      await rollbackDeployment(
        deployment.applicationName,
        deployment.environment,
        deployment.region,
        deployment.version, // version is undefined for standard rollback
        "system" // Using default value for rollbacks, could be replaced with a user input
      );
      
      setSuccess(`Successfully rolled back ${deployment.applicationName} in ${deployment.environment} (${deployment.region}) to version ${deployment.version}`);
      
      // Reload deployments to show the rollback
      await loadDeployments();
    } catch (err) {
      setError("Failed to rollback. Please try again later.");
      console.error(err);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="container mx-auto py-8 space-y-8">
      <h1 className="text-3xl font-bold">Deployment Manager</h1>
      
      {error && (
        <Alert variant="destructive">
          <AlertCircle className="h-4 w-4" />
          <AlertTitle>Error</AlertTitle>
          <AlertDescription>{error}</AlertDescription>
        </Alert>
      )}
      
      {success && (
        <Alert className="bg-green-50 border-green-200">
          <CheckCircle2 className="h-4 w-4 text-green-600" />
          <AlertTitle className="text-green-800">Success</AlertTitle>
          <AlertDescription className="text-green-700">{success}</AlertDescription>
        </Alert>
      )}
      
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-8">
        <DeploymentForm onDeploy={handleDeploy} />
        
        <Card>
          <CardHeader>
            <CardTitle>API Integration</CardTitle>
            <CardDescription>
              This component demonstrates how to integrate the backend API with the existing UI components.
            </CardDescription>
          </CardHeader>
          <CardContent>
            <ul className="list-disc pl-5 space-y-2">
              <li>Uses <code>fetchDeployments()</code> to load data from the backend</li>
              <li>Calls <code>releaseVersion()</code> when the form is submitted</li>
              <li>Calls <code>rollbackDeployment()</code> when rollback is requested</li>
              <li>Handles loading states and error messages</li>
            </ul>
          </CardContent>
        </Card>
      </div>
      
      {loading ? (
        <Card>
          <CardContent className="py-8">
            <p className="text-center text-muted-foreground">Loading deployments...</p>
          </CardContent>
        </Card>
      ) : (
        <DeploymentTree deployments={deployments} onRollback={handleRollback} />
      )}
    </div>
  );
};
