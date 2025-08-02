import { useState, useEffect } from "react";
import { DeploymentTree, Deployment } from "@/components/DeploymentTree";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { Settings, Download, RefreshCw, Plus } from "lucide-react";
import { useToast } from "@/hooks/use-toast";
import { fetchDeployments, rollbackDeployment } from "@/api/deploymentApi";
import { DeploymentForm } from "@/components/DeploymentForm";
import { Dialog, DialogContent, DialogTrigger } from "@/components/ui/dialog";

const Index = () => {
  const [deployments, setDeployments] = useState<Deployment[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [versionLimit, setVersionLimit] = useState<number>(10);
  const { toast } = useToast();

  // Load deployments on component mount
  useEffect(() => {
    loadDeployments();
  }, []);

  const loadDeployments = async () => {
    try {
      setIsLoading(true);
      setError(null);
      const data = await fetchDeployments(versionLimit);
      setDeployments(data);
    } catch (err) {
      setError("Failed to load deployments");
      toast({
        variant: "destructive",
        title: "Error",
        description: "Failed to load deployments from the backend",
      });
      console.error(err);
    } finally {
      setIsLoading(false);
    }
  };

  const handleRollback = async (deployment: Deployment) => {
    try {
      setIsLoading(true);
      await rollbackDeployment(
        deployment.applicationName,
        deployment.environment,
        deployment.region
      );
      
      // Reload deployments to show the rollback
      await loadDeployments();
      
      toast({
        title: "Rollback Successful",
        description: `Rolled back ${deployment.applicationName} to previous version in ${deployment.environment}`,
      });
    } catch (err) {
      toast({
        variant: "destructive",
        title: "Rollback Failed",
        description: "Failed to rollback deployment",
      });
      console.error(err);
    } finally {
      setIsLoading(false);
    }
  };

  const handleLoadFromBackend = async () => {
    await loadDeployments();
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
      await fetch(`http://localhost:8080/release`, {
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
      
      // Reload deployments to show the new deployment
      await loadDeployments();
      
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

  const totalDeployments = deployments.length;
  const activeDeployments = deployments.filter(d => d.status === 'active').length;
  const uniqueApps = new Set(deployments.map(d => d.applicationName)).size;

  return (
    <div className="min-h-screen bg-background">
      {/* Header */}
      <header className="border-b bg-card/50 backdrop-blur">
        <div className="container mx-auto px-4 py-4">
          <div className="flex items-center justify-between">
            <div>
              <h1 className="text-2xl font-bold bg-gradient-to-r from-primary to-accent bg-clip-text text-transparent">
                CI/CD Platform
              </h1>
              <p className="text-muted-foreground">Deployment Management Dashboard</p>
            </div>
            <div className="flex items-center gap-4">
              <div className="flex gap-2 items-center">
                <Badge variant="outline" className="gap-1">
                  <span className="w-2 h-2 rounded-full bg-success"></span>
                  {activeDeployments} Active
                </Badge>
                <Badge variant="outline">{uniqueApps} Apps</Badge>
                <Badge variant="outline">{totalDeployments} Total</Badge>
                <div className="flex items-center gap-2 ml-4">
                  <span className="text-xs text-muted-foreground">Versions:</span>
                  <select
                    className="text-xs border rounded px-1 py-0.5"
                    value={versionLimit}
                    onChange={(e) => {
                      const newLimit = parseInt(e.target.value);
                      setVersionLimit(newLimit);
                      loadDeployments();
                    }}
                  >
                    <option value="5">5</option>
                    <option value="10">10</option>
                    <option value="20">20</option>
                    <option value="50">50</option>
                    <option value="100">100</option>
                  </select>
                </div>
              </div>
              <Button 
                variant="outline" 
                size="sm"
                onClick={handleLoadFromBackend}
                disabled={isLoading}
                className="gap-2"
              >
                {isLoading ? (
                  <RefreshCw className="h-4 w-4 animate-spin" />
                ) : (
                  <Download className="h-4 w-4" />
                )}
                Load from Backend
              </Button>
              <Button variant="ghost" size="sm">
                <Settings className="h-4 w-4" />
              </Button>
            </div>
          </div>
        </div>
      </header>

      {/* Main Content */}
      <main className="container mx-auto px-4 py-8">
        <div className="flex justify-end mb-4">
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
        
        {isLoading ? (
          <div className="flex justify-center items-center h-64">
            <RefreshCw className="h-8 w-8 animate-spin text-primary" />
          </div>
        ) : error ? (
          <div className="flex justify-center items-center h-64 text-destructive">
            <p>{error}</p>
          </div>
        ) : (
          <DeploymentTree
            deployments={deployments}
            onRollback={handleRollback}
          />
        )}
      </main>
    </div>
  );
};

export default Index;
