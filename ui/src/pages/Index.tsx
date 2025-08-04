import { useState } from "react";
import { Deployment } from "@/components/DeploymentTree";
import { HierarchicalDeploymentTree } from "@/components/HierarchicalDeploymentTree";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { Settings, RefreshCw, Plus } from "lucide-react";
import { useToast } from "@/hooks/use-toast";
import { rollbackDeployment } from "@/api/deploymentApi";
import { DeploymentForm } from "@/components/DeploymentForm";
import { Dialog, DialogContent, DialogTrigger } from "@/components/ui/dialog";

const Index = () => {
  const [isLoading, setIsLoading] = useState(false);
  const { toast } = useToast();

  const handleRollback = async (deployment: Deployment) => {
    try {
      setIsLoading(true);
      await rollbackDeployment(
        deployment.applicationName,
        deployment.environment,
        deployment.region,
        deployment.version
      );
      
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
              <Button variant="ghost" size="sm">
                <Settings className="h-4 w-4" />
              </Button>
            </div>
          </div>
        </div>
      </header>

      {/* Main Content */}
      <main className="container mx-auto px-4 py-8">
        <div className="flex justify-end items-center mb-4">
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
        ) : (
          <HierarchicalDeploymentTree
            onRollback={handleRollback}
          />
        )}
      </main>
    </div>
  );
};

export default Index;
