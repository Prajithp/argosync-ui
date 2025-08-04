import { useState, useEffect } from "react";
import { useParams, useNavigate, Link } from "react-router-dom";
import { Layers, RefreshCw, ChevronRight, Home, Globe } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { useToast } from "@/hooks/use-toast";
import { 
  Environment, 
  Application, 
  Region,
  fetchEnvironmentsForApplicationAndRegion, 
  fetchAllApplications,
  fetchRegionsForApplication
} from "@/api/deploymentApi";

const EnvironmentsPage = () => {
  const { appId, regionId } = useParams<{ appId: string; regionId: string }>();
  const [environments, setEnvironments] = useState<Environment[]>([]);
  const [application, setApplication] = useState<Application | null>(null);
  const [region, setRegion] = useState<Region | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const { toast } = useToast();
  const navigate = useNavigate();

  useEffect(() => {
    if (!appId || isNaN(Number(appId)) || !regionId || isNaN(Number(regionId))) {
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
        
        // Load environments for this application and region
        const envs = await fetchEnvironmentsForApplicationAndRegion(Number(appId), Number(regionId));
        console.log("Loaded environments:", envs);
        setEnvironments(envs);
      } catch (err) {
        setError("Failed to load environments");
        toast({
          variant: "destructive",
          title: "Error",
          description: "Failed to load environments from the backend",
        });
        console.error(err);
      } finally {
        setIsLoading(false);
      }
    };

    loadData();
  }, [appId, regionId, toast]);

  const handleEnvironmentClick = (envId: number) => {
    if (envId && !isNaN(envId)) {
      navigate(`/applications/${appId}/regions/${regionId}/environments/${envId}/versions`);
    } else {
      toast({
        variant: "destructive",
        title: "Navigation Error",
        description: "Invalid environment ID",
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
                <span className="ml-1 text-sm font-medium text-primary">
                  {region?.name || "Loading..."} Environments
                </span>
              </div>
            </li>
          </ol>
        </nav>
      </div>

      {/* Main Content */}
      <main className="container mx-auto px-4 py-8">
        <h2 className="text-xl font-semibold mb-4">
          Environments for {application?.name || "Loading..."} in {region?.name || "Loading..."} ({region?.code || ""})
        </h2>
        
        {isLoading ? (
          <div className="flex justify-center items-center h-64">
            <RefreshCw className="h-8 w-8 animate-spin text-primary" />
          </div>
        ) : error ? (
          <div className="flex justify-center items-center h-64 text-destructive">
            <p>{error}</p>
          </div>
        ) : environments.length === 0 ? (
          <div className="text-center py-8 text-muted-foreground">
            No environments found for this application and region.
          </div>
        ) : (
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
            {environments.map(env => (
              <Card
                key={`env-${env.id || env.ID}`}
                className="cursor-pointer hover:bg-accent/10 transition-colors"
                onClick={() => handleEnvironmentClick(env.id || env.ID)}
              >
                <CardHeader className="pb-2">
                  <CardTitle className="flex items-center gap-2">
                    <Layers className="h-5 w-5 text-accent" />
                    <span>{env.name}</span>
                  </CardTitle>
                </CardHeader>
                <CardContent>
                  <p className="text-sm text-muted-foreground">
                    Click to view versions in this environment
                  </p>
                </CardContent>
              </Card>
            ))}
          </div>
        )}
      </main>
    </div>
  );
};

export default EnvironmentsPage;