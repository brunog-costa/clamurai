# Clamurai Traefik Plugin 

Wake up, Clamurai...<br>
    &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;We have some malware to burn.

## The Project

Clamurai is a Middleware built for protecting applications that need to store artifacts sent by the client to any sort of storage back-end. 

The inspiration from the project comes from tools such as GuardDuty Malware Inspection for S3 feature, which rely on receiving an object before inspecting it, requiring aditional logic and resources to protect the environment - with clamurai, we are able to stop malicious input before it touches our environment.


## Usage 

### Pre-requisites 

### Traefik Configuration  

| Option              | Type   | Description                              | Default            |
| ------------------- | ------ | ---------------------------------------- | ------------------ |
| `clamdHost`         | string | Host\:port of the ClamAV daemon (clamd)  | `"localhost:3310"` |
| `actionOnDetection` | string | `"block"`, `"allow"`, `"redirect"`       | `"block"`          |
| `scanTimeout`       | string | Timeout duration for scan                | `"5s"`             |
| `logLevel`          | string | `"info"`, `"debug"`, `"warn"`, `"error"` | `"info"`           |


## Examples 

### Protecting an AWS Application 
### Protecting an Azure Application 
### Protecting an GCP Application 


## 🧱 Roadmap

* Block/Allow list logic for signatures 
* Improved support for ZIP/GZIP decompression
* YARA Integration

## 🧠 Want to contribute?

PRs and issues are always welcome, also, make sure to check the roadmap section in this file. 
