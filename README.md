# Clamurai Traefik Plugin

Wake up, clamurai...<br>
    &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;We have some malware to burn.

## Overview

Clamurai is a Traefik middleware plugin that provides real-time malware scanning for HTTP request bodies using ClamAV. Unlike traditional async scanning solutions (GuardDuty Malware Protection, Microsoft Defender for Storage), Clamurai inspects files in-flight, blocking malicious content before it reaches your backend infrastructure.

**Key Benefits:**
- 🛡️ Inline protection - scan before storage
- ⚡ Real-time blocking of malicious uploads
- 🔧 Easy integration with existing Traefik deployments
- 📊 Structured JSON logging with scan results and file hashes 

## Configuration

### Plugin Options

| Option                      | Type   | Description                                                                 | Default            |
| --------------------------- | ------ | --------------------------------------------------------------------------- | ------------------ |
| `clamdAddress`              | string | Host:port of the ClamAV daemon (clamd)                                      | `"localhost:3310"` |
| `clamavReadTimeout`         | uint64 | Timeout in seconds for waiting on ClamAV scan results                       | `3600`             |
| `clamavConnectionTimeout`   | uint64 | Timeout in seconds for establishing connection with ClamAV daemon           | `90`               |
| `alertMode`                 | bool   | When `true`, logs threats without blocking requests (monitoring mode)       | `true`             |


## Installation

1. Add the plugin to your Traefik static configuration:

```yaml
experimental:
  plugins:
    clamurai:
      moduleName: github.com/brunog-costa/clamurai
      version: v0.1.0  # Use latest version
```

2. Deploy a ClamAV daemon (clamd) accessible to your Traefik instance

3. Configure the middleware in your dynamic configuration (see examples below)

## Usage Examples

### Prerequisites

- Traefik instance with the clamurai plugin installed
- ClamAV daemon (clamd) running and accessible
- Basic knowledge of Traefik dynamic configuration

### Protecting S3 Uploads

Scan file uploads before they reach your S3 bucket:

```yaml
http:
  routers:
    s3-upload:
      rule: "PathPrefix(`/`) && Method(`PUT`)"
      service: s3-backend
      entryPoints:
        - web
      middlewares:
        - clamurai-scan
    
    s3-other:
      rule: "PathPrefix(`/`)"
      service: s3-backend
      entryPoints:
        - web

  services:
    s3-backend:
      loadBalancer:
        servers:
          - url: "https://s3.amazonaws.com"
        
  middlewares:
    clamurai-scan:
      plugin:
        clamurai:
          clamdAddress: "localhost:3310"
          clamavReadTimeout: 3600
          clamavConnectionTimeout: 90
          alertMode: false  # Block malicious uploads
```

### Monitoring Mode

Log threats without blocking (useful for testing):

```yaml
middlewares:
  clamurai-monitor:
    plugin:
      clamurai:
        clamdAddress: "clamav-service:3310"
        clamavReadTimeout: 3600
        clamavConnectionTimeout: 90
        alertMode: true  # Log only, don't block
```

## Features

- ✅ Real-time malware scanning using ClamAV
- ✅ SHA-256 hash calculation for all scanned files
- ✅ Structured JSON logging with scan results
- ✅ Alert mode for monitoring without blocking
- ✅ Configurable timeouts for scan operations
- ✅ Request body restoration for upstream services

## Roadmap

- [ ] **Extended Use Cases**
  - Azure Blob Storage examples
  - Google Cloud Storage examples
  - IaC templates (Terraform, CloudFormation)
  - Performance testing scripts

- [ ] **Enhanced Configuration**
  - Configurable logging destinations (file, stdout, remote endpoint)
  - Custom response messages
  - Whitelist/blacklist by file type or size

- [ ] **ClamAV Management**
  - Health check endpoints
  - Signature update monitoring
  - Connection pool management

- [ ] **Advanced Security**
  - YARA rules integration
  - Custom signature support
  - Multi-engine scanning support 


## 🧠 Want to contribute?

PRs and issues are always welcome, also, make sure to check the roadmap section above. 
