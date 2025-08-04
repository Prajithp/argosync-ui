import { Toaster } from "@/components/ui/toaster";
import { Toaster as Sonner } from "@/components/ui/sonner";
import { TooltipProvider } from "@/components/ui/tooltip";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { BrowserRouter, Routes, Route, Navigate } from "react-router-dom";
import ApplicationsPage from "./pages/ApplicationsPage";
import RegionsPage from "./pages/RegionsPage";
import EnvironmentsPage from "./pages/EnvironmentsPage";
import VersionsPage from "./pages/VersionsPage";
import NotFound from "./pages/NotFound";

const queryClient = new QueryClient();

const App = () => (
  <QueryClientProvider client={queryClient}>
    <TooltipProvider>
      <Toaster />
      <Sonner />
      <BrowserRouter>
        <Routes>
          {/* Redirect from old index to new applications page */}
          <Route path="/" element={<ApplicationsPage />} />
          
          {/* Hierarchical navigation routes */}
          <Route path="/applications/:appId/regions" element={<RegionsPage />} />
          <Route path="/applications/:appId/regions/:regionId/environments" element={<EnvironmentsPage />} />
          <Route path="/applications/:appId/regions/:regionId/environments/:envId/versions" element={<VersionsPage />} />
          
          {/* Catch-all route */}
          <Route path="*" element={<NotFound />} />
        </Routes>
      </BrowserRouter>
    </TooltipProvider>
  </QueryClientProvider>
);

export default App;
