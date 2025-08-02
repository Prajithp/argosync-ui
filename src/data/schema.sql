-- Optimized CI/CD Platform Database Schema
-- Designed for 1000+ services with efficient querying

-- Applications/Services table
CREATE TABLE applications (
  id SERIAL PRIMARY KEY,
  name VARCHAR(255) NOT NULL UNIQUE,
  description TEXT,
  team VARCHAR(100),
  repository_url VARCHAR(500),
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Environments table (normalized)
CREATE TABLE environments (
  id SERIAL PRIMARY KEY,
  name VARCHAR(50) NOT NULL UNIQUE, -- dev, staging, production
  description TEXT,
  priority INTEGER DEFAULT 0
);

-- Regions table (normalized)
CREATE TABLE regions (
  id SERIAL PRIMARY KEY,
  code VARCHAR(20) NOT NULL UNIQUE, -- us-east-1, eu-west-1, etc.
  name VARCHAR(100) NOT NULL, -- US East (N. Virginia)
  continent VARCHAR(50),
  active BOOLEAN DEFAULT true
);

-- Main deployments table
CREATE TABLE deployments (
  id SERIAL PRIMARY KEY,
  application_id INTEGER NOT NULL REFERENCES applications(id),
  environment_id INTEGER NOT NULL REFERENCES environments(id),
  region_id INTEGER NOT NULL REFERENCES regions(id),
  version VARCHAR(100) NOT NULL,
  status VARCHAR(20) NOT NULL DEFAULT 'active', -- active, inactive, failed
  deployed_by VARCHAR(100),
  deployed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  rollback_target_id INTEGER REFERENCES deployments(id),
  metadata JSONB, -- Additional deployment info (build_id, commit_hash, etc.)
  
  -- Ensure only one active deployment per app/env/region combination
  CONSTRAINT unique_active_deployment 
    EXCLUDE (application_id WITH =, environment_id WITH =, region_id WITH =) 
    WHERE (status = 'active')
);

-- Deployment history for audit trail
CREATE TABLE deployment_history (
  id SERIAL PRIMARY KEY,
  deployment_id INTEGER NOT NULL REFERENCES deployments(id),
  action VARCHAR(50) NOT NULL, -- deploy, rollback, deactivate
  performed_by VARCHAR(100),
  performed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  details JSONB
);

-- Indexes for optimal performance with 1000+ services
CREATE INDEX idx_deployments_app_env_region ON deployments(application_id, environment_id, region_id);
CREATE INDEX idx_deployments_status ON deployments(status);
CREATE INDEX idx_deployments_app_status ON deployments(application_id, status);
CREATE INDEX idx_deployments_deployed_at ON deployments(deployed_at DESC);
CREATE INDEX idx_applications_name ON applications(name);
CREATE INDEX idx_applications_team ON applications(team);
CREATE INDEX idx_deployment_history_deployment_id ON deployment_history(deployment_id);
CREATE INDEX idx_deployment_history_performed_at ON deployment_history(performed_at DESC);

-- Insert default environments
INSERT INTO environments (name, description, priority) VALUES
('development', 'Development environment', 1),
('staging', 'Staging environment', 2),
('production', 'Production environment', 3);

-- Insert regions
INSERT INTO regions (code, name, continent) VALUES
('us-east-1', 'US East (N. Virginia)', 'North America'),
('us-west-2', 'US West (Oregon)', 'North America'),
('eu-west-1', 'Europe (Ireland)', 'Europe'),
('eu-central-1', 'Europe (Frankfurt)', 'Europe'),
('ap-southeast-1', 'Asia Pacific (Singapore)', 'Asia'),
('ap-northeast-1', 'Asia Pacific (Tokyo)', 'Asia');

-- Views for common queries
CREATE VIEW active_deployments AS
SELECT 
  a.name as application_name,
  a.team,
  e.name as environment,
  r.code as region_code,
  r.name as region_name,
  d.version,
  d.deployed_at,
  d.deployed_by
FROM deployments d
JOIN applications a ON d.application_id = a.id
JOIN environments e ON d.environment_id = e.id
JOIN regions r ON d.region_id = r.id
WHERE d.status = 'active'
ORDER BY a.name, e.priority, r.code;

CREATE VIEW deployment_summary AS
SELECT 
  a.name as application_name,
  a.team,
  COUNT(DISTINCT d.environment_id) as env_count,
  COUNT(DISTINCT d.region_id) as region_count,
  COUNT(*) as total_deployments,
  MAX(d.deployed_at) as last_deployment
FROM applications a
LEFT JOIN deployments d ON a.id = d.application_id
GROUP BY a.id, a.name, a.team
ORDER BY last_deployment DESC NULLS LAST;