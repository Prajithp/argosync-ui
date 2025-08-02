// Legacy sample data - kept for backward compatibility
// For optimized data model, see optimizedSampleData.ts

import { Deployment } from "@/components/DeploymentTree";
import { transformToLegacyFormat } from './optimizedSampleData';

// Use optimized data transformed to legacy format
export const sampleDeployments: Deployment[] = transformToLegacyFormat();