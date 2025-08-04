import { useState } from "react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Upload, AlertTriangle } from "lucide-react";
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

interface DeploymentFormProps {
  onDeploy: (deployment: {
    applicationName: string;
    environment: string;
    region: string;
    version: string;
    deployedBy: string;
  }) => void;
}

export const DeploymentForm = ({ onDeploy }: DeploymentFormProps) => {
  const [formData, setFormData] = useState({
    applicationName: "",
    environment: "",
    region: "",
    version: "",
    deployedBy: "",
  });
  const [showConfirmation, setShowConfirmation] = useState(false);

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (formData.applicationName && formData.environment && formData.region && formData.version) {
      setShowConfirmation(true);
    }
  };

  const handleConfirmDeploy = () => {
    // If deployedBy is empty, use the current user or a default value
    const deployData = {
      ...formData,
      deployedBy: formData.deployedBy || "system"
    };
    
    // Close the confirmation dialog
    setShowConfirmation(false);
    
    // Call the actual deploy function
    onDeploy(deployData);
    
    // Reset the form
    setFormData({ applicationName: "", environment: "", region: "", version: "", deployedBy: "" });
  };

  return (
    <Card className="w-full max-w-2xl">
      <AlertDialog open={showConfirmation} onOpenChange={setShowConfirmation}>
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle className="flex items-center gap-2">
              <AlertTriangle className="h-5 w-5 text-amber-500" />
              Confirm Deployment
            </AlertDialogTitle>
            <AlertDialogDescription>
              Are you sure you want to deploy <strong>{formData.applicationName}</strong> version <strong>{formData.version}</strong> to <strong>{formData.environment}</strong> in <strong>{formData.region}</strong>?
              
              {formData.environment === "production" && (
                <p className="mt-2 text-amber-500 font-semibold">
                  This is a PRODUCTION deployment and will affect live users!
                </p>
              )}
            </AlertDialogDescription>
          </AlertDialogHeader>
          <AlertDialogFooter>
            <AlertDialogCancel>Cancel</AlertDialogCancel>
            <AlertDialogAction
              onClick={handleConfirmDeploy}
              className={formData.environment === "production" ? "bg-amber-500 hover:bg-amber-600" : ""}
            >
              Yes, Deploy
            </AlertDialogAction>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog>


      <CardHeader>
        <CardTitle className="flex items-center gap-2">
          <Upload className="h-5 w-5" />
          Deploy Application
        </CardTitle>
        <CardDescription>
          Configure and deploy your application to the specified environment
        </CardDescription>
      </CardHeader>
      <CardContent>
        <form onSubmit={handleSubmit} className="space-y-4">
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div className="space-y-2">
              <Label htmlFor="applicationName">Application Name</Label>
              <Input
                id="applicationName"
                placeholder="e.g., user-service"
                value={formData.applicationName}
                onChange={(e) => setFormData({ ...formData, applicationName: e.target.value })}
                required
              />
            </div>
            <div className="space-y-2">
              <Label htmlFor="version">Version</Label>
              <Input
                id="version"
                placeholder="e.g., v1.2.3"
                value={formData.version}
                onChange={(e) => setFormData({ ...formData, version: e.target.value })}
                required
              />
            </div>
          </div>
          
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div className="space-y-2">
              <Label htmlFor="environment">Environment</Label>
              <Select value={formData.environment} onValueChange={(value) => setFormData({ ...formData, environment: value })}>
                <SelectTrigger>
                  <SelectValue placeholder="Select environment" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="development">Development</SelectItem>
                  <SelectItem value="staging">Staging</SelectItem>
                  <SelectItem value="production">Production</SelectItem>
                </SelectContent>
              </Select>
            </div>
            <div className="space-y-2">
              <Label htmlFor="region">Region</Label>
              <Select value={formData.region} onValueChange={(value) => setFormData({ ...formData, region: value })}>
                <SelectTrigger>
                  <SelectValue placeholder="Select region" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="us-east-1">US East (N. Virginia)</SelectItem>
                  <SelectItem value="us-west-2">US West (Oregon)</SelectItem>
                  <SelectItem value="eu-west-1">Europe (Ireland)</SelectItem>
                  <SelectItem value="ap-southeast-1">Asia Pacific (Singapore)</SelectItem>
                </SelectContent>
              </Select>
            </div>
          </div>
          
          <div className="space-y-2">
            <Label htmlFor="deployedBy">Deployed By</Label>
            <Input
              id="deployedBy"
              placeholder="e.g., john.doe"
              value={formData.deployedBy}
              onChange={(e) => setFormData({ ...formData, deployedBy: e.target.value })}
            />
          </div>
          
          <Button type="submit" className="w-full">
            Deploy Application
          </Button>
        </form>
      </CardContent>
    </Card>
  );
};