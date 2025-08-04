import { useState, useEffect } from "react";
import { ChevronDown, ChevronRight, Server, Globe, Layers, Package, RotateCcw, AlertTriangle } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { cn } from "@/lib/utils";
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
  Application, 
  Region, 
  Environment, 
  Version,
  fetchAllApplications,
  fetchRegionsForApplication,
  fetchEnvironmentsForApplicationAndRegion,
  fetchVersionsForApplicationEnvironmentRegion
} from "@/api/deploymentApi";
import { Deployment } from "@/components/DeploymentTree";

interface HierarchicalDeploymentTreeProps {
  onRollback: (deployment: Deployment) => void;
}

export const HierarchicalDeploymentTree = ({ onRollback }: HierarchicalDeploymentTreeProps) => {
  // State for hierarchical data
  const [applications, setApplications] = useState<Application[]>([]);
  const [regionsMap, setRegionsMap] = useState<Record<number, Region[]>>({});
  const [environmentsMap, setEnvironmentsMap] = useState<Record<string, Environment[]>>({});
  const [versionsMap, setVersionsMap] = useState<Record<string, Version[]>>({});
  
  // Loading states
  const [isLoadingApps, setIsLoadingApps] = useState(true);
  const [loadingRegionsFor, setLoadingRegionsFor] = useState<number | null>(null);
  const [loadingEnvironmentsFor, setLoadingEnvironmentsFor] = useState<string | null>(null);
  const [loadingVersionsFor, setLoadingVersionsFor] = useState<string | null>(null);
  
  // Expanded nodes tracking
  const [expandedApps, setExpandedApps] = useState<Set<number>>(new Set());
  const [expandedRegions, setExpandedRegions] = useState<Set<string>>(new Set());
  const [expandedEnvironments, setExpandedEnvironments] = useState<Set<string>>(new Set());
  
  // Rollback confirmation
  const [deploymentToRollback, setDeploymentToRollback] = useState<Deployment | null>(null);
  const [showRollbackConfirmation, setShowRollbackConfirmation] = useState(false);

  // Load applications on component mount
  useEffect(() => {
    loadApplications();
  }, []);

  const loadApplications = async () => {
    try {
      setIsLoadingApps(true);
      const apps = await fetchAllApplications();
      setApplications(apps);
    } catch (error) {
      console.error("Error loading applications:", error);
    } finally {
      setIsLoadingApps(false);
    }
  };

  const toggleApplication = async (appId: number) => {
    // Validate appId
    if (appId === undefined || appId === null || isNaN(appId)) {
      console.error("Invalid application ID:", appId);
      return;
    }
    
    // Toggle expanded state
    const newExpandedApps = new Set(expandedApps);
    if (newExpandedApps.has(appId)) {
      newExpandedApps.delete(appId);
    } else {
      newExpandedApps.add(appId);
      // Load regions if not already loaded
      if (!regionsMap[appId]) {
        try {
          await loadRegionsForApplication(appId);
        } catch (error) {
          console.error(`Error loading regions for application ${appId}:`, error);
        }
      }
    }
    setExpandedApps(newExpandedApps);
  };

  const loadRegionsForApplication = async (appId: number) => {
    try {
      setLoadingRegionsFor(appId);
      const regions = await fetchRegionsForApplication(appId);
      setRegionsMap(prev => ({
        ...prev,
        [appId]: regions
      }));
    } catch (error) {
      console.error(`Error loading regions for application ${appId}:`, error);
    } finally {
      setLoadingRegionsFor(null);
    }
  };

  const toggleRegion = async (appId: number, regionId: number) => {
    // Validate IDs
    if (appId === undefined || appId === null || isNaN(appId)) {
      console.error("Invalid application ID:", appId);
      return;
    }
    if (regionId === undefined || regionId === null || isNaN(regionId)) {
      console.error("Invalid region ID:", regionId);
      return;
    }
    
    const key = `${appId}-${regionId}`;
    const newExpandedRegions = new Set(expandedRegions);
    if (newExpandedRegions.has(key)) {
      newExpandedRegions.delete(key);
    } else {
      newExpandedRegions.add(key);
      // Load environments if not already loaded
      if (!environmentsMap[key]) {
        try {
          await loadEnvironmentsForApplicationAndRegion(appId, regionId);
        } catch (error) {
          console.error(`Error loading environments for app ${appId} and region ${regionId}:`, error);
        }
      }
    }
    setExpandedRegions(newExpandedRegions);
  };

  const loadEnvironmentsForApplicationAndRegion = async (appId: number, regionId: number) => {
    const key = `${appId}-${regionId}`;
    try {
      setLoadingEnvironmentsFor(key);
      const environments = await fetchEnvironmentsForApplicationAndRegion(appId, regionId);
      setEnvironmentsMap(prev => ({
        ...prev,
        [key]: environments
      }));
    } catch (error) {
      console.error(`Error loading environments for app ${appId} and region ${regionId}:`, error);
    } finally {
      setLoadingEnvironmentsFor(null);
    }
  };

  const toggleEnvironment = async (appId: number, regionId: number, envId: number) => {
    // Validate IDs
    if (appId === undefined || appId === null || isNaN(appId)) {
      console.error("Invalid application ID:", appId);
      return;
    }
    if (regionId === undefined || regionId === null || isNaN(regionId)) {
      console.error("Invalid region ID:", regionId);
      return;
    }
    if (envId === undefined || envId === null || isNaN(envId)) {
      console.error("Invalid environment ID:", envId);
      return;
    }
    
    const key = `${appId}-${regionId}-${envId}`;
    const newExpandedEnvironments = new Set(expandedEnvironments);
    if (newExpandedEnvironments.has(key)) {
      newExpandedEnvironments.delete(key);
    } else {
      newExpandedEnvironments.add(key);
      // Load versions if not already loaded
      if (!versionsMap[key]) {
        try {
          await loadVersionsForApplicationEnvironmentRegion(appId, envId, regionId);
        } catch (error) {
          console.error(`Error loading versions for app ${appId}, env ${envId}, and region ${regionId}:`, error);
        }
      }
    }
    setExpandedEnvironments(newExpandedEnvironments);
  };

  const loadVersionsForApplicationEnvironmentRegion = async (appId: number, envId: number, regionId: number) => {
    const key = `${appId}-${regionId}-${envId}`;
    try {
      setLoadingVersionsFor(key);
      const versions = await fetchVersionsForApplicationEnvironmentRegion(appId, envId, regionId);
      setVersionsMap(prev => ({
        ...prev,
        [key]: versions
      }));
    } catch (error) {
      console.error(`Error loading versions for app ${appId}, env ${envId}, and region ${regionId}:`, error);
    } finally {
      setLoadingVersionsFor(null);
    }
  };

  const handleConfirmRollback = () => {
    if (deploymentToRollback) {
      // Close the confirmation dialog
      setShowRollbackConfirmation(false);
      
      // Call the actual rollback function
      onRollback(deploymentToRollback);
      
      // Reset the state
      setDeploymentToRollback(null);
    }
  };

  // Convert Version to Deployment for rollback
  const versionToDeployment = (
    version: Version, 
    appName: string, 
    envName: string, 
    regionName: string
  ): Deployment => {
    return {
      applicationName: appName,
      environment: envName,
      region: regionName,
      version: version.version,
      timestamp: version.created_at || version.CreatedAt,
      status: 'inactive', // Assuming all versions in history are inactive except the current one
      deployedBy: 'unknown' // This information might not be available
    };
  };

  const renderVersions = (
    versions: Version[], 
    appId: number, 
    regionId: number, 
    envId: number,
    appName: string,
    envName: string,
    regionName: string
  ) => {
    return versions
      .sort((a, b) => {
        const bDate = new Date(b.created_at || b.CreatedAt).getTime();
        const aDate = new Date(a.created_at || a.CreatedAt).getTime();
        return bDate - aDate;
      })
      .map((version) => {
        const isActive = (version.status || version.Status) === 'active';
        return (
          <div
            key={`version-${version.id || version.ID}`}
            className="ml-8 flex items-center justify-between p-2 border rounded-md bg-card hover:bg-accent/50 transition-colors"
          >
            <div className="flex items-center gap-2">
              <Package className="h-4 w-4 text-muted-foreground" />
              <span className="font-mono text-sm">{version.version}</span>
              <Badge variant={isActive ? 'default' : 'secondary'}>
                {isActive ? 'active' : 'inactive'}
              </Badge>
              <span className="text-xs text-muted-foreground">
                {new Date(version.created_at || version.CreatedAt).toLocaleString()}
              </span>
            </div>
            {!isActive && (
              <Button
                size="sm"
                variant="outline"
                onClick={() => {
                  const deployment = versionToDeployment(version, appName, envName, regionName);
                  setDeploymentToRollback(deployment);
                  setShowRollbackConfirmation(true);
                }}
                className="h-8"
              >
                <RotateCcw className="h-3 w-3 mr-1" />
                Rollback
              </Button>
            )}
          </div>
        );
      });
  };

  const renderEnvironments = (
    environments: Environment[], 
    appId: number, 
    regionId: number,
    appName: string,
    regionName: string
  ) => {
    return environments.map(env => {
      const key = `${appId}-${regionId}-${env.id || env.ID}`;
      const isExpanded = expandedEnvironments.has(key);
      const isLoading = loadingVersionsFor === key;
      const versions = versionsMap[key] || [];
      
      return (
        <div key={`env-${env.id || env.ID}`} className="ml-6">
          <div
            className="flex items-center gap-2 p-2 cursor-pointer hover:bg-accent/50 rounded-md transition-colors"
            onClick={() => toggleEnvironment(appId, regionId, env.id || env.ID)}
          >
            {isExpanded ? <ChevronDown className="h-4 w-4" /> : <ChevronRight className="h-4 w-4" />}
            <Layers className="h-4 w-4 text-accent" />
            <span className="font-medium capitalize">{env.name}</span>
            {isLoading ? (
              <Badge variant="outline">Loading...</Badge>
            ) : (
              versions.length > 0 && <Badge variant="outline">{versions.length} versions</Badge>
            )}
          </div>
          {isExpanded && versions.length > 0 && 
            renderVersions(versions, appId, regionId, env.id, appName, env.name, regionName)}
        </div>
      );
    });
  };

  const renderRegions = (regions: Region[], appId: number, appName: string) => {
    return regions.map(region => {
      const key = `${appId}-${region.id || region.ID}`;
      const isExpanded = expandedRegions.has(key);
      const isLoading = loadingEnvironmentsFor === key;
      const environments = environmentsMap[key] || [];
      
      return (
        <div key={`region-${region.id || region.ID}`} className="ml-4">
          <div
            className="flex items-center gap-2 p-2 cursor-pointer hover:bg-accent/50 rounded-md transition-colors"
            onClick={() => toggleRegion(appId, region.id || region.ID)}
          >
            {isExpanded ? <ChevronDown className="h-4 w-4" /> : <ChevronRight className="h-4 w-4" />}
            <Globe className="h-4 w-4 text-primary" />
            <span className="font-medium">{region.name} ({region.code})</span>
            {isLoading ? (
              <Badge variant="secondary">Loading...</Badge>
            ) : (
              environments.length > 0 && 
              <Badge variant="secondary">{environments.length} environments</Badge>
            )}
          </div>
          {isExpanded && environments.length > 0 && 
            renderEnvironments(environments, appId, region.id, appName, region.code)}
        </div>
      );
    });
  };

  return (
    <Card className="w-full">
      <AlertDialog open={showRollbackConfirmation} onOpenChange={setShowRollbackConfirmation}>
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle className="flex items-center gap-2">
              <AlertTriangle className="h-5 w-5 text-amber-500" />
              Confirm Rollback
            </AlertDialogTitle>
            <AlertDialogDescription>
              {deploymentToRollback && (
                <>
                  Are you sure you want to rollback <strong>{deploymentToRollback.applicationName}</strong> in <strong>{deploymentToRollback.environment}</strong> ({deploymentToRollback.region}) to version <strong>{deploymentToRollback.version}</strong>?
                  <p className="mt-2 text-amber-500">This action cannot be undone and may cause service disruption.</p>
                </>
              )}
            </AlertDialogDescription>
          </AlertDialogHeader>
          <AlertDialogFooter>
            <AlertDialogCancel>Cancel</AlertDialogCancel>
            <AlertDialogAction onClick={handleConfirmRollback} className="bg-amber-500 hover:bg-amber-600">
              Yes, Rollback
            </AlertDialogAction>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog>

      <CardHeader>
        <CardTitle className="flex items-center gap-2">
          <Server className="h-5 w-5" />
          Services
        </CardTitle>
      </CardHeader>
      <CardContent>
        {isLoadingApps ? (
          <div className="text-center py-8 text-muted-foreground">
            Loading applications...
          </div>
        ) : applications.length === 0 ? (
          <div className="text-center py-8 text-muted-foreground">
            No applications found. Deploy an application to see it here.
          </div>
        ) : (
          <div className="space-y-2">
            {applications.map(app => {
              const appId = app.id || app.ID;
              const isExpanded = expandedApps.has(appId);
              const isLoading = loadingRegionsFor === appId;
              const regions = regionsMap[appId] || [];
              
              return (
                <div key={`app-${app.id || app.ID}`} className="border rounded-lg p-2">
                  <div
                    className="flex items-center gap-2 p-2 cursor-pointer hover:bg-accent/50 rounded-md transition-colors"
                    onClick={() => toggleApplication(app.id || app.ID)}
                  >
                    {isExpanded ? <ChevronDown className="h-4 w-4" /> : <ChevronRight className="h-4 w-4" />}
                    <Server className="h-5 w-5 text-success" />
                    <span className="font-semibold text-lg">{app.name}</span>
                    {isLoading ? (
                      <Badge>Loading...</Badge>
                    ) : (
                      regions.length > 0 && <Badge>{regions.length} regions</Badge>
                    )}
                  </div>
                  {isExpanded && regions.length > 0 && renderRegions(regions, app.id || app.ID, app.name)}
                </div>
              );
            })}
          </div>
        )}
      </CardContent>
    </Card>
  );
};