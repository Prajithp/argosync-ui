import { useState } from "react";
import { DeploymentForm } from "@/components/DeploymentForm";
import { DeploymentTree, Deployment } from "@/components/DeploymentTree";
import { sampleDeployments } from "@/data/sampleData";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { Settings, Download, RefreshCw } from "lucide-react";
import { useToast } from "@/hooks/use-toast";

const Index = () => {
  const [deployments, setDeployments] = useState<Deployment[]>(sampleDeployments);
  const [isLoading, setIsLoading] = useState(false);
  const { toast } = useToast();

  const handleDeploy = (newDeployment: Omit<Deployment, 'timestamp' | 'status'>) => {
    // Mark previous versions of the same app/env/region as inactive
    const updatedDeployments = deployments.map(deployment => {
      if (
        deployment.applicationName === newDeployment.applicationName &&
        deployment.environment === newDeployment.environment &&
        deployment.region === newDeployment.region
      ) {
        return { ...deployment, status: 'inactive' as const };
      }
      return deployment;
    });

    // Add new deployment
    const deployment: Deployment = {
      ...newDeployment,
      timestamp: new Date().toISOString(),
      status: 'active'
    };

    setDeployments([...updatedDeployments, deployment]);
    
    toast({
      title: "Deployment Successful",
      description: `${newDeployment.applicationName} ${newDeployment.version} deployed to ${newDeployment.environment} in ${newDeployment.region}`,
    });
  };

  const handleRollback = (deployment: Deployment) => {
    // Mark current active deployment as inactive and this one as active
    const updatedDeployments = deployments.map(d => {
      if (
        d.applicationName === deployment.applicationName &&
        d.environment === deployment.environment &&
        d.region === deployment.region
      ) {
        if (d.version === deployment.version) {
          return { ...d, status: 'active' as const, timestamp: new Date().toISOString() };
        } else {
          return { ...d, status: 'inactive' as const };
        }
      }
      return d;
    });

    setDeployments(updatedDeployments);
    
    toast({
      title: "Rollback Successful",
      description: `Rolled back ${deployment.applicationName} to ${deployment.version} in ${deployment.environment}`,
    });
  };

  const handleLoadFromBackend = async () => {
    setIsLoading(true);
    // Simulate API call
    await new Promise(resolve => setTimeout(resolve, 1500));
    setIsLoading(false);
    
    toast({
      title: "Backend Integration Ready",
      description: "Connect to your CI/CD API to load real deployment data",
    });
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
              <div className="flex gap-2">
                <Badge variant="outline" className="gap-1">
                  <span className="w-2 h-2 rounded-full bg-success"></span>
                  {activeDeployments} Active
                </Badge>
                <Badge variant="outline">{uniqueApps} Apps</Badge>
                <Badge variant="outline">{totalDeployments} Total</Badge>
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
        <div className="grid grid-cols-1 xl:grid-cols-3 gap-8">
          {/* Deployment Form */}
          <div className="xl:col-span-1">
            <DeploymentForm onDeploy={handleDeploy} />
          </div>
          
          {/* Deployment Tree */}
          <div className="xl:col-span-2">
            <DeploymentTree 
              deployments={deployments} 
              onRollback={handleRollback}
            />
          </div>
        </div>
      </main>
    </div>
  );
};

export default Index;
