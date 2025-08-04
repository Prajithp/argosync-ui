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
  applications,
  regions,
  environments,
  optimizedDeployments,
  getApplicationById,
  getEnvironmentById,
  getRegionById,
  getDeploymentHistory,
  type Application,
  type Region,
  type Environment,
  type OptimizedDeployment
} from "@/data/optimizedSampleData";
import { Deployment } from "@/components/DeploymentTree";

interface HierarchicalDeploymentTreeProps {
  onRollback: (deployment: Deployment) => void;
}

interface Version {
  id: number;
  version: string;
  status: 'active' | 'inactive';
  deployedAt: string;
  deployedBy: string;
}

export const HierarchicalDeploymentTree = ({ onRollback }: HierarchicalDeploymentTreeProps) => {
  // State for hierarchical data
  const [appData] = useState<Application[]>(applications);
  const [regionsMap, setRegionsMap] = useState<Record<number, Region[]>>({});
  const [environmentsMap, setEnvironmentsMap] = useState<Record<string, Environment[]>>({});
  const [versionsMap, setVersionsMap] = useState<Record<string, Version[]>>({});
  
  // Loading states - simplified for sample data
  const [isLoadingApps] = useState(false);
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

  // Load data from sample data instead of API
  useEffect(() => {
    // Pre-load all data for better UX since it's sample data
    preloadSampleData();
  }, []);

  const preloadSampleData = () => {
    // Pre-populate regions for each application
    const newRegionsMap: Record<number, Region[]> = {};
    applications.forEach(app => {
      // Get regions that have deployments for this app
      const appRegionIds = optimizedDeployments
        .filter(d => d.applicationId === app.id)
        .map(d => d.regionId);
      const uniqueRegionIds = [...new Set(appRegionIds)];
      newRegionsMap[app.id] = regions.filter(r => uniqueRegionIds.includes(r.id));
    });
    setRegionsMap(newRegionsMap);
  };

  const toggleApplication = (appId: number) => {
    const newExpandedApps = new Set(expandedApps);
    if (newExpandedApps.has(appId)) {
      newExpandedApps.delete(appId);
    } else {
      newExpandedApps.add(appId);
    }
    setExpandedApps(newExpandedApps);
  };

  const toggleRegion = (appId: number, regionId: number) => {
    const key = `${appId}-${regionId}`;
    const newExpandedRegions = new Set(expandedRegions);
    if (newExpandedRegions.has(key)) {
      newExpandedRegions.delete(key);
    } else {
      newExpandedRegions.add(key);
      // Load environments for this app-region combination
      if (!environmentsMap[key]) {
        loadEnvironmentsForAppAndRegion(appId, regionId);
      }
    }
    setExpandedRegions(newExpandedRegions);
  };

  const loadEnvironmentsForAppAndRegion = (appId: number, regionId: number) => {
    const key = `${appId}-${regionId}`;
    setLoadingEnvironmentsFor(key);
    
    // Get environments that have deployments for this app-region combination
    const appEnvIds = optimizedDeployments
      .filter(d => d.applicationId === appId && d.regionId === regionId)
      .map(d => d.environmentId);
    const uniqueEnvIds = [...new Set(appEnvIds)];
    const envs = environments.filter(e => uniqueEnvIds.includes(e.id));
    
    setEnvironmentsMap(prev => ({
      ...prev,
      [key]: envs
    }));
    setLoadingEnvironmentsFor(null);
  };

  const toggleEnvironment = (appId: number, regionId: number, envId: number) => {
    const key = `${appId}-${regionId}-${envId}`;
    const newExpandedEnvironments = new Set(expandedEnvironments);
    if (newExpandedEnvironments.has(key)) {
      newExpandedEnvironments.delete(key);
    } else {
      newExpandedEnvironments.add(key);
      // Load versions if not already loaded
      if (!versionsMap[key]) {
        loadVersionsForAppEnvRegion(appId, envId, regionId);
      }
    }
    setExpandedEnvironments(newExpandedEnvironments);
  };

  const loadVersionsForAppEnvRegion = (appId: number, envId: number, regionId: number) => {
    const key = `${appId}-${regionId}-${envId}`;
    setLoadingVersionsFor(key);
    
    // Get deployment history for this specific combination
    const deploymentHistory = getDeploymentHistory(appId, envId, regionId);
    const versions: Version[] = deploymentHistory
      .filter(deployment => deployment.status !== 'failed') // Filter out failed deployments
      .map(deployment => ({
        id: deployment.id,
        version: deployment.version,
        status: deployment.status as 'active' | 'inactive',
        deployedAt: deployment.deployedAt,
        deployedBy: deployment.deployedBy
      }));
    
    setVersionsMap(prev => ({
      ...prev,
      [key]: versions
    }));
    setLoadingVersionsFor(null);
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
      timestamp: version.deployedAt,
      status: version.status,
      deployedBy: version.deployedBy
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
        const bDate = new Date(b.deployedAt).getTime();
        const aDate = new Date(a.deployedAt).getTime();
        return bDate - aDate;
      })
      .map((version) => {
        const isActive = version.status === 'active';
        return (
          <div
            key={`version-${version.id}`}
            className="ml-8 flex items-center justify-between p-3 border rounded-lg bg-card hover:bg-accent/50 transition-all duration-200 shadow-sm"
          >
            <div className="flex items-center gap-3">
              <Package className="h-4 w-4 text-muted-foreground" />
              <span className="font-mono text-sm font-medium">{version.version}</span>
              <Badge variant={isActive ? 'default' : 'secondary'} className="font-medium">
                {isActive ? 'active' : 'inactive'}
              </Badge>
              <div className="flex flex-col gap-1">
                <span className="text-xs text-muted-foreground">
                  {new Date(version.deployedAt).toLocaleString()}
                </span>
                <span className="text-xs text-muted-foreground">
                  by {version.deployedBy}
                </span>
              </div>
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
      const key = `${appId}-${regionId}-${env.id}`;
      const isExpanded = expandedEnvironments.has(key);
      const isLoading = loadingVersionsFor === key;
      const versions = versionsMap[key] || [];
      
      return (
        <div key={`env-${env.id}`} className="ml-6">
          <div
            className="flex items-center gap-3 p-3 cursor-pointer hover:bg-accent/50 rounded-lg transition-all duration-200"
            onClick={() => toggleEnvironment(appId, regionId, env.id)}
          >
            {isExpanded ? <ChevronDown className="h-4 w-4" /> : <ChevronRight className="h-4 w-4" />}
            <Layers className="h-4 w-4 text-accent" />
            <span className="font-medium capitalize">{env.name}</span>
            {isLoading ? (
              <Badge variant="outline" className="animate-pulse">Loading...</Badge>
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
      const key = `${appId}-${region.id}`;
      const isExpanded = expandedRegions.has(key);
      const isLoading = loadingEnvironmentsFor === key;
      const environments = environmentsMap[key] || [];
      
      return (
        <div key={`region-${region.id}`} className="ml-4">
          <div
            className="flex items-center gap-3 p-3 cursor-pointer hover:bg-accent/50 rounded-lg transition-all duration-200"
            onClick={() => toggleRegion(appId, region.id)}
          >
            {isExpanded ? <ChevronDown className="h-4 w-4" /> : <ChevronRight className="h-4 w-4" />}
            <Globe className="h-4 w-4 text-primary" />
            <span className="font-medium">{region.name}</span>
            <Badge variant="outline" className="text-xs">{region.code}</Badge>
            {isLoading ? (
              <Badge variant="secondary" className="animate-pulse">Loading...</Badge>
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
        ) : appData.length === 0 ? (
          <div className="text-center py-8 text-muted-foreground">
            No applications found. Deploy an application to see it here.
          </div>
        ) : (
          <div className="space-y-3">
            {appData.map(app => {
              const appId = app.id;
              const isExpanded = expandedApps.has(appId);
              const isLoading = loadingRegionsFor === appId;
              const appRegions = regionsMap[appId] || [];
              
              return (
                <div key={`app-${app.id}`} className="border rounded-xl p-4 bg-gradient-to-r from-card to-card/50 shadow-sm hover:shadow-md transition-all duration-200">
                  <div
                    className="flex items-center gap-3 p-3 cursor-pointer hover:bg-accent/30 rounded-lg transition-all duration-200"
                    onClick={() => toggleApplication(app.id)}
                  >
                    {isExpanded ? <ChevronDown className="h-5 w-5" /> : <ChevronRight className="h-5 w-5" />}
                    <Server className="h-6 w-6 text-primary" />
                    <div className="flex-1">
                      <div className="flex items-center gap-3">
                        <span className="font-semibold text-lg">{app.name}</span>
                        <Badge variant="outline" className="text-xs">{app.team}</Badge>
                      </div>
                      {app.description && (
                        <p className="text-sm text-muted-foreground mt-1">{app.description}</p>
                      )}
                    </div>
                    {isLoading ? (
                      <Badge className="animate-pulse">Loading...</Badge>
                    ) : (
                      appRegions.length > 0 && <Badge>{appRegions.length} regions</Badge>
                    )}
                  </div>
                  {isExpanded && appRegions.length > 0 && renderRegions(appRegions, app.id, app.name)}
                </div>
              );
            })}
          </div>
        )}
      </CardContent>
    </Card>
  );
};