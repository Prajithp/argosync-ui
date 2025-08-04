import { useState, useEffect } from "react";
import { useNavigate } from "react-router-dom";
import { Server, RefreshCw, Plus } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { useToast } from "@/hooks/use-toast";
import { Application, fetchAllApplications } from "@/api/deploymentApi";
import { Dialog, DialogContent, DialogTrigger } from "@/components/ui/dialog";
import { DeploymentForm } from "@/components/DeploymentForm";

const ApplicationsPage = () => {
  const [applications, setApplications] = useState<Application[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const { toast } = useToast();
  const navigate = useNavigate();

  useEffect(() => {
    loadApplications();
  }, []);

  const loadApplications = async () => {
    try {
      setIsLoading(true);
      setError(null);
      const apps = await fetchAllApplications();
      
      // Add detailed logging for debugging
      console.log("Loaded applications:", JSON.stringify(apps, null, 2));
      if (apps && apps.length > 0) {
        console.log("First application details:", {
          id: apps[0].id || apps[0].ID,
          name: apps[0].name,
          created_at: apps[0].created_at || apps[0].CreatedAt,
          updated_at: apps[0].updated_at || apps[0].UpdatedAt,
          type_of_id: typeof (apps[0].id || apps[0].ID),
          is_valid_number: !isNaN(Number(apps[0].id || apps[0].ID)),
          is_valid_date: (apps[0].created_at || apps[0].CreatedAt) &&
                         !isNaN(new Date(apps[0].created_at || apps[0].CreatedAt).getTime())
        });
      } else {
        console.log("No applications loaded");
      }
      
      // Ensure all applications have valid IDs and dates
      const validatedApps = apps.map(app => {
        // Get ID from either lowercase or uppercase field
        const rawId = app.id || app.ID;
        // Validate ID
        const id = typeof rawId === 'string' ? parseInt(rawId) : rawId;
        
        // Get dates from either lowercase or uppercase field
        const rawCreatedAt = app.created_at || app.CreatedAt;
        const rawUpdatedAt = app.updated_at || app.UpdatedAt;
        
        // Validate dates
        const created_at = rawCreatedAt && !isNaN(new Date(rawCreatedAt).getTime())
          ? rawCreatedAt
          : new Date().toISOString();
        
        const updated_at = rawUpdatedAt && !isNaN(new Date(rawUpdatedAt).getTime())
          ? rawUpdatedAt
          : new Date().toISOString();
        
        return {
          ...app,
          id,
          created_at,
          updated_at
        };
      });
      
      console.log("Validated applications:", JSON.stringify(validatedApps, null, 2));
      setApplications(validatedApps);
    } catch (err) {
      setError("Failed to load applications");
      toast({
        variant: "destructive",
        title: "Error",
        description: "Failed to load applications from the backend",
      });
      console.error(err);
    } finally {
      setIsLoading(false);
    }
  };

  const handleDeploy = async (formData: {
    applicationName: string;
    environment: string;
    region: string;
    version: string;
  }) => {
    try {
      setIsLoading(true);
      
      // Call the API to release a new version
      await fetch(`/api/v1/release`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({
          application: formData.applicationName,
          environment: formData.environment,
          region: formData.region,
          version: formData.version,
          deployed_by: "user@example.com",
        }),
      });
      
      // Reload applications to show the new deployment
      await loadApplications();
      
      toast({
        title: "Deployment Successful",
        description: `Deployed ${formData.applicationName} version ${formData.version} to ${formData.environment}`,
      });
    } catch (err) {
      toast({
        variant: "destructive",
        title: "Deployment Failed",
        description: "Failed to deploy application",
      });
      console.error(err);
    } finally {
      setIsLoading(false);
    }
  };

  const handleApplicationClick = (appId: number | string) => {
    console.log("Application clicked with ID:", appId, "Type:", typeof appId);
    
    // Ensure appId is a valid number
    const numericId = Number(appId);
    
    if (!isNaN(numericId) && numericId > 0) {
      console.log("Navigating to regions page with appId:", numericId);
      navigate(`/applications/${numericId}/regions`);
    } else {
      console.error("Invalid application ID:", appId);
      toast({
        variant: "destructive",
        title: "Navigation Error",
        description: `Invalid application ID: ${appId}`,
      });
    }
  };

  return (
    <div className="min-h-screen bg-background">
      {/* Header */}
      <header className="border-b bg-card/50 backdrop-blur">
        <div className="container mx-auto px-4 py-4">
          <div className="flex items-center justify-between">
            <div>
              <h1 className="text-2xl font-bold bg-gradient-to-r from-primary to-accent bg-clip-text text-transparent">
                ArgoSync
              </h1>
              <p className="text-muted-foreground">Kubernetes Deployment Management Dashboard</p>
            </div>
            <div className="flex items-center gap-4">
              <Dialog>
                <DialogTrigger asChild>
                  <Button className="gap-2">
                    <Plus className="h-4 w-4" />
                    Deploy New Version
                  </Button>
                </DialogTrigger>
                <DialogContent className="sm:max-w-[600px]">
                  <DeploymentForm onDeploy={handleDeploy} />
                </DialogContent>
              </Dialog>
            </div>
          </div>
        </div>
      </header>

      {/* Breadcrumbs */}
      <div className="container mx-auto px-4 py-2 border-b">
        <nav className="flex" aria-label="Breadcrumb">
          <ol className="inline-flex items-center space-x-1 md:space-x-3">
            <li className="inline-flex items-center">
              <span className="inline-flex items-center text-sm font-medium text-primary">
                Applications
              </span>
            </li>
          </ol>
        </nav>
      </div>

      {/* Main Content */}
      <main className="container mx-auto px-4 py-8">
        <h2 className="text-xl font-semibold mb-4">Applications</h2>
        
        {isLoading ? (
          <div className="flex justify-center items-center h-64">
            <RefreshCw className="h-8 w-8 animate-spin text-primary" />
          </div>
        ) : error ? (
          <div className="flex justify-center items-center h-64 text-destructive">
            <p>{error}</p>
          </div>
        ) : applications.length === 0 ? (
          <div className="text-center py-8 text-muted-foreground">
            No applications found. Deploy an application to see it here.
          </div>
        ) : (
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
            {applications.map(app => {
              // Get ID from either lowercase or uppercase field
              const rawId = app.id || app.ID;
              // Ensure we have a valid ID
              const appId = typeof rawId === 'string' ? parseInt(rawId) : rawId;
              const isValidId = !isNaN(Number(appId)) && Number(appId) > 0;
              
              return (
                <Card
                  key={`app-${app.id}`}
                  className={`${isValidId ? 'cursor-pointer hover:bg-accent/10' : 'opacity-70'} transition-colors`}
                  onClick={() => isValidId ? handleApplicationClick(appId) : null}
                >
                  <CardHeader className="pb-2">
                    <CardTitle className="flex items-center gap-2">
                      <Server className="h-5 w-5 text-success" />
                      <span>{app.name}</span>
                      {!isValidId && (
                        <Badge variant="destructive" className="ml-2">Invalid ID</Badge>
                      )}
                    </CardTitle>
                  </CardHeader>
                  <CardContent>
                    <p className="text-sm text-muted-foreground">
                      Created: {(app.created_at || app.CreatedAt) &&
                               !isNaN(new Date(app.created_at || app.CreatedAt).getTime())
                        ? new Date(app.created_at || app.CreatedAt).toLocaleDateString()
                        : "Unknown date"}
                    </p>
                    {isValidId ? (
                      <p className="text-xs text-muted-foreground mt-1">
                        Click to view regions
                      </p>
                    ) : (
                      <p className="text-xs text-destructive mt-1">
                        ID: {(app.id || app.ID || "undefined")} (Invalid)
                      </p>
                    )}
                  </CardContent>
                </Card>
              );
            })}
          </div>
        )}
      </main>
    </div>
  );
};

export default ApplicationsPage;