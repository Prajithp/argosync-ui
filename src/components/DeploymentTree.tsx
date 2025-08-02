import { useState } from "react";
import { ChevronDown, ChevronRight, Server, Globe, Layers, Package, RotateCcw } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { cn } from "@/lib/utils";

export interface Deployment {
  applicationName: string;
  environment: string;
  region: string;
  version: string;
  timestamp: string;
  status: 'active' | 'inactive';
}

interface DeploymentTreeProps {
  deployments: Deployment[];
  onRollback: (deployment: Deployment) => void;
}

interface TreeNode {
  [key: string]: any;
}

export const DeploymentTree = ({ deployments, onRollback }: DeploymentTreeProps) => {
  const [expandedNodes, setExpandedNodes] = useState<Set<string>>(new Set());

  const toggleNode = (nodeId: string) => {
    const newExpanded = new Set(expandedNodes);
    if (newExpanded.has(nodeId)) {
      newExpanded.delete(nodeId);
    } else {
      newExpanded.add(nodeId);
    }
    setExpandedNodes(newExpanded);
  };

  // Organize deployments hierarchically
  const organizeDeployments = () => {
    const tree: TreeNode = {};
    
    deployments.forEach(deployment => {
      if (!tree[deployment.applicationName]) {
        tree[deployment.applicationName] = {};
      }
      if (!tree[deployment.applicationName][deployment.region]) {
        tree[deployment.applicationName][deployment.region] = {};
      }
      if (!tree[deployment.applicationName][deployment.region][deployment.environment]) {
        tree[deployment.applicationName][deployment.region][deployment.environment] = [];
      }
      
      tree[deployment.applicationName][deployment.region][deployment.environment].push(deployment);
    });

    return tree;
  };

  const tree = organizeDeployments();

  const renderVersions = (versions: Deployment[], path: string) => {
    return versions
      .sort((a, b) => new Date(b.timestamp).getTime() - new Date(a.timestamp).getTime())
      .map((deployment, index) => (
        <div key={`${path}-${deployment.version}-${index}`} className="ml-8 flex items-center justify-between p-2 border rounded-md bg-card hover:bg-accent/50 transition-colors">
          <div className="flex items-center gap-2">
            <Package className="h-4 w-4 text-muted-foreground" />
            <span className="font-mono text-sm">{deployment.version}</span>
            <Badge variant={deployment.status === 'active' ? 'default' : 'secondary'}>
              {deployment.status}
            </Badge>
            <span className="text-xs text-muted-foreground">
              {new Date(deployment.timestamp).toLocaleString()}
            </span>
          </div>
          {deployment.status === 'inactive' && (
            <Button
              size="sm"
              variant="outline"
              onClick={() => onRollback(deployment)}
              className="h-8"
            >
              <RotateCcw className="h-3 w-3 mr-1" />
              Rollback
            </Button>
          )}
        </div>
      ));
  };

  const renderEnvironments = (environments: TreeNode, appName: string, region: string) => {
    return Object.entries(environments).map(([envName, versions]) => {
      const envPath = `${appName}-${region}-${envName}`;
      const isExpanded = expandedNodes.has(envPath);
      
      return (
        <div key={envPath} className="ml-6">
          <div
            className="flex items-center gap-2 p-2 cursor-pointer hover:bg-accent/50 rounded-md transition-colors"
            onClick={() => toggleNode(envPath)}
          >
            {isExpanded ? <ChevronDown className="h-4 w-4" /> : <ChevronRight className="h-4 w-4" />}
            <Layers className="h-4 w-4 text-accent" />
            <span className="font-medium capitalize">{envName}</span>
            <Badge variant="outline">{(versions as Deployment[]).length} versions</Badge>
          </div>
          {isExpanded && renderVersions(versions as Deployment[], envPath)}
        </div>
      );
    });
  };

  const renderRegions = (regions: TreeNode, appName: string) => {
    return Object.entries(regions).map(([regionName, environments]) => {
      const regionPath = `${appName}-${regionName}`;
      const isExpanded = expandedNodes.has(regionPath);
      
      return (
        <div key={regionPath} className="ml-4">
          <div
            className="flex items-center gap-2 p-2 cursor-pointer hover:bg-accent/50 rounded-md transition-colors"
            onClick={() => toggleNode(regionPath)}
          >
            {isExpanded ? <ChevronDown className="h-4 w-4" /> : <ChevronRight className="h-4 w-4" />}
            <Globe className="h-4 w-4 text-primary" />
            <span className="font-medium">{regionName}</span>
            <Badge variant="secondary">
              {Object.values(environments).flat().length} deployments
            </Badge>
          </div>
          {isExpanded && renderEnvironments(environments, appName, regionName)}
        </div>
      );
    });
  };

  return (
    <Card className="w-full">
      <CardHeader>
        <CardTitle className="flex items-center gap-2">
          <Server className="h-5 w-5" />
          Deployment Hierarchy
        </CardTitle>
      </CardHeader>
      <CardContent>
        {Object.keys(tree).length === 0 ? (
          <div className="text-center py-8 text-muted-foreground">
            No deployments found. Deploy an application to see it here.
          </div>
        ) : (
          <div className="space-y-2">
            {Object.entries(tree).map(([appName, regions]) => {
              const appPath = appName;
              const isExpanded = expandedNodes.has(appPath);
              
              return (
                <div key={appPath} className="border rounded-lg p-2">
                  <div
                    className="flex items-center gap-2 p-2 cursor-pointer hover:bg-accent/50 rounded-md transition-colors"
                    onClick={() => toggleNode(appPath)}
                  >
                    {isExpanded ? <ChevronDown className="h-4 w-4" /> : <ChevronRight className="h-4 w-4" />}
                    <Server className="h-5 w-5 text-success" />
                    <span className="font-semibold text-lg">{appName}</span>
                    <Badge>
                      {Object.values(regions).flatMap(r => Object.values(r)).flat().length} total deployments
                    </Badge>
                  </div>
                  {isExpanded && renderRegions(regions, appName)}
                </div>
              );
            })}
          </div>
        )}
      </CardContent>
    </Card>
  );
};