import { useState, useEffect } from "react";
import { useParams, Link } from "react-router-dom";
import { Package, RefreshCw, ChevronRight, Home, Globe, Layers, RotateCcw, AlertTriangle } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { useToast } from "@/hooks/use-toast";
import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
} from "@/components/ui/alert-dialog";
import { 
  Version, 
  Application, 
  Region,
  Environment,
  fetchVersionsForApplicationEnvironmentRegion, 
  fetchAllApplications,
  fetchRegionsForApplication,
  fetchEnvironmentsForApplicationAndRegion,
  rollbackDeployment
} from "@/api/deploymentApi";
import { Deployment } from "@/components/DeploymentTree";

const VersionsPage = () => {
  const { appId, regionId, envId } = useParams<{ appId: string; regionId: string; envId: string }>();
  const [versions, setVersions] = useState<Version[]>([]);
  const [application, setApplication] = useState<Application | null>(null);
  const [region, setRegion] = useState<Region | null>(null);
  const [environment, setEnvironment] = useState<Environment | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const { toast } = useToast();
  
  // Rollback confirmation
  const [versionToRollback, setVersionToRollback] = useState<Version | null>(null);
  const [showRollbackConfirmation, setShowRollbackConfirmation] = useState(false);

  useEffect(() => {
    if (!appId || isNaN(Number(appId)) || !regionId || isNaN(Number(regionId)) || !envId || isNaN(Number(envId))) {
      setError("Invalid parameters");
      setIsLoading(false);
      return;
    }

    const loadData = async () => {
      try {
        setIsLoading(true);
        setError(null);
        
        // Load application details
        const apps = await fetchAllApplications();
        const app = apps.find(a => (a.id || a.ID) === Number(appId));
        if (!app) {
          setError("Application not found");
          setIsLoading(false);
          return;
        }
        setApplication(app);
        
        // Load region details
        const regions = await fetchRegionsForApplication(Number(appId));
        const reg = regions.find(r => (r.id || r.ID) === Number(regionId));
        if (!reg) {
          setError("Region not found");
          setIsLoading(false);
          return;
        }
        setRegion(reg);
        
        // Load environment details
        const envs = await fetchEnvironmentsForApplicationAndRegion(Number(appId), Number(regionId));
        const env = envs.find(e => (e.id || e.ID) === Number(envId));
        if (!env) {
          setError("Environment not found");
          setIsLoading(false);
          return;
        }
        setEnvironment(env);
        
        // Load versions for this application, environment, and region
        const vers = await fetchVersionsForApplicationEnvironmentRegion(
          Number(appId), 
          Number(envId), 
          Number(regionId)
        );
        console.log("Loaded versions:", vers);
        setVersions(vers);
      } catch (err) {
        setError("Failed to load versions");
        toast({
          variant: "destructive",
          title: "Error",
          description: "Failed to load versions from the backend",
        });
        console.error(err);
      } finally {
        setIsLoading(false);
      }
    };

    loadData();
  }, [appId, regionId, envId, toast]);

  const handleRollback = async (version: Version) => {
    setVersionToRollback(version);
    setShowRollbackConfirmation(true);
  };

  const confirmRollback = async () => {
    if (!versionToRollback || !application || !environment || !region) return;
    
    try {
      setIsLoading(true);
      await rollbackDeployment(
        application.name,
        environment.name,
        region.code,
        versionToRollback.version
      );
      
      // Reload versions to show the rollback
      const vers = await fetchVersionsForApplicationEnvironmentRegion(
        Number(appId), 
        Number(envId), 
        Number(regionId)
      );
      setVersions(vers);
      
      toast({
        title: "Rollback Successful",
        description: `Rolled back ${application.name} to version ${versionToRollback.version} in ${environment.name}`,
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
      setShowRollbackConfirmation(false);
      setVersionToRollback(null);
    }
  };

  return (
    <div className="min-h-screen bg-background">
      <AlertDialog open={showRollbackConfirmation} onOpenChange={setShowRollbackConfirmation}>
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle className="flex items-center gap-2">
              <AlertTriangle className="h-5 w-5 text-amber-500" />
              Confirm Rollback
            </AlertDialogTitle>
            <AlertDialogDescription>
              {versionToRollback && application && environment && region && (
                <>
                  Are you sure you want to rollback <strong>{application.name}</strong> in <strong>{environment.name}</strong> ({region.code}) to version <strong>{versionToRollback.version}</strong>?
                  <p className="mt-2 text-amber-500">This action cannot be undone and may cause service disruption.</p>
                </>
              )}
            </AlertDialogDescription>
          </AlertDialogHeader>
          <AlertDialogFooter>
            <AlertDialogCancel>Cancel</AlertDialogCancel>
            <AlertDialogAction onClick={confirmRollback} className="bg-amber-500 hover:bg-amber-600">
              Yes, Rollback
            </AlertDialogAction>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog>

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
          </div>
        </div>
      </header>

      {/* Breadcrumbs */}
      <div className="container mx-auto px-4 py-2 border-b">
        <nav className="flex" aria-label="Breadcrumb">
          <ol className="inline-flex items-center space-x-1 md:space-x-3">
            <li className="inline-flex items-center">
              <Link to="/" className="inline-flex items-center text-sm font-medium text-primary hover:text-primary/80">
                <Home className="w-4 h-4 mr-2" />
                Applications
              </Link>
            </li>
            <li>
              <div className="flex items-center">
                <ChevronRight className="w-4 h-4 text-muted-foreground" />
                <Link 
                  to={`/applications/${appId}/regions`} 
                  className="ml-1 text-sm font-medium text-primary hover:text-primary/80"
                >
                  <Globe className="w-4 h-4 mr-1 inline" />
                  {application?.name || "Loading..."} Regions
                </Link>
              </div>
            </li>
            <li>
              <div className="flex items-center">
                <ChevronRight className="w-4 h-4 text-muted-foreground" />
                <Link 
                  to={`/applications/${appId}/regions/${regionId}/environments`} 
                  className="ml-1 text-sm font-medium text-primary hover:text-primary/80"
                >
                  <Layers className="w-4 h-4 mr-1 inline" />
                  {region?.name || "Loading..."} Environments
                </Link>
              </div>
            </li>
            <li>
              <div className="flex items-center">
                <ChevronRight className="w-4 h-4 text-muted-foreground" />
                <span className="ml-1 text-sm font-medium text-primary">
                  {environment?.name || "Loading..."} Versions
                </span>
              </div>
            </li>
          </ol>
        </nav>
      </div>

      {/* Main Content */}
      <main className="container mx-auto px-4 py-8">
        <h2 className="text-xl font-semibold mb-4">
          Versions for {application?.name || "Loading..."} in {environment?.name || "Loading..."} ({region?.code || ""})
        </h2>
        
        {isLoading ? (
          <div className="flex justify-center items-center h-64">
            <RefreshCw className="h-8 w-8 animate-spin text-primary" />
          </div>
        ) : error ? (
          <div className="flex justify-center items-center h-64 text-destructive">
            <p>{error}</p>
          </div>
        ) : versions.length === 0 ? (
          <div className="text-center py-8 text-muted-foreground">
            No versions found for this application, environment, and region.
          </div>
        ) : (
          <div className="space-y-4">
            {versions
              .sort((a, b) => {
                const bDate = new Date(b.created_at || b.CreatedAt).getTime();
                const aDate = new Date(a.created_at || a.CreatedAt).getTime();
                return bDate - aDate;
              })
              .map((version) => {
                // Check the status field to determine if the version is active
                const isActive = (version.status || version.Status) === 'active';
                return (
                  <Card key={`version-${version.id || version.ID}`} className="overflow-hidden">
                    <div className={`${isActive ? 'bg-primary/10' : ''}`}>
                      <CardHeader className="pb-2">
                        <CardTitle className="flex items-center justify-between">
                          <div className="flex items-center gap-2">
                            <Package className="h-5 w-5 text-primary" />
                            <span className="font-mono">{version.version}</span>
                            <Badge variant={isActive ? 'default' : 'secondary'}>
                              {isActive ? 'active' : 'inactive'}
                            </Badge>
                          </div>
                          {!isActive && (
                            <Button
                              size="sm"
                              variant="outline"
                              onClick={() => handleRollback(version)}
                              className="h-8"
                            >
                              <RotateCcw className="h-3 w-3 mr-1" />
                              Rollback
                            </Button>
                          )}
                        </CardTitle>
                      </CardHeader>
                      <CardContent>
                        <div className="text-sm text-muted-foreground">
                          Deployed: {new Date(version.created_at || version.CreatedAt).toLocaleString()}
                        </div>
                      </CardContent>
                    </div>
                  </Card>
                );
              })}
          </div>
        )}
      </main>
    </div>
  );
};

export default VersionsPage;