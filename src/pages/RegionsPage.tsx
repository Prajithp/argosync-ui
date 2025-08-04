import { useState, useEffect } from "react";
import { useParams, useNavigate, Link } from "react-router-dom";
import { Globe, RefreshCw, ChevronRight, Home } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { useToast } from "@/hooks/use-toast";
import { Region, Application, fetchRegionsForApplication, fetchAllApplications } from "@/api/deploymentApi";

const RegionsPage = () => {
  const { appId } = useParams<{ appId: string }>();
  const [regions, setRegions] = useState<Region[]>([]);
  const [application, setApplication] = useState<Application | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const { toast } = useToast();
  const navigate = useNavigate();

  useEffect(() => {
    if (!appId || isNaN(Number(appId))) {
      setError("Invalid application ID");
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
        
        // Load regions for this application
        const appRegions = await fetchRegionsForApplication(Number(appId));
        console.log("Loaded regions:", appRegions);
        setRegions(appRegions);
      } catch (err) {
        setError("Failed to load regions");
        toast({
          variant: "destructive",
          title: "Error",
          description: "Failed to load regions from the backend",
        });
        console.error(err);
      } finally {
        setIsLoading(false);
      }
    };

    loadData();
  }, [appId, toast]);

  const handleRegionClick = (regionId: number) => {
    if (regionId && !isNaN(regionId)) {
      navigate(`/applications/${appId}/regions/${regionId}/environments`);
    } else {
      toast({
        variant: "destructive",
        title: "Navigation Error",
        description: "Invalid region ID",
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
                <span className="ml-1 text-sm font-medium text-primary">
                  {application?.name || "Loading..."} Regions
                </span>
              </div>
            </li>
          </ol>
        </nav>
      </div>

      {/* Main Content */}
      <main className="container mx-auto px-4 py-8">
        <h2 className="text-xl font-semibold mb-4">
          Regions for {application?.name || "Loading..."}
        </h2>
        
        {isLoading ? (
          <div className="flex justify-center items-center h-64">
            <RefreshCw className="h-8 w-8 animate-spin text-primary" />
          </div>
        ) : error ? (
          <div className="flex justify-center items-center h-64 text-destructive">
            <p>{error}</p>
          </div>
        ) : regions.length === 0 ? (
          <div className="text-center py-8 text-muted-foreground">
            No regions found for this application.
          </div>
        ) : (
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
            {regions.map(region => (
              <Card
                key={`region-${region.id || region.ID}`}
                className="cursor-pointer hover:bg-accent/10 transition-colors"
                onClick={() => handleRegionClick(region.id || region.ID)}
              >
                <CardHeader className="pb-2">
                  <CardTitle className="flex items-center gap-2">
                    <Globe className="h-5 w-5 text-primary" />
                    <span>{region.name}</span>
                    <span className="text-sm text-muted-foreground">({region.code})</span>
                  </CardTitle>
                </CardHeader>
                <CardContent>
                  <p className="text-sm text-muted-foreground">
                    Click to view environments in this region
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

export default RegionsPage;