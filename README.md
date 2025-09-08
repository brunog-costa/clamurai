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

* Implement tests and documentations for the code 
	* Create Readme file
	* DrawIO C4 + DA/DDI 
* Improve detections and clamav customization 
	* Custom clam signature lib 
	* Create flag for alert and block mode 
* Tackle logging and outputs 
	* Generate an EFK stack in order to monitor the plugin
	* Create flag for stdout and file logging output 
	* Check if clamav logs can be displayed as json messing with the config file 
* Test the package 
	* Multi-cloud url re-write test logic 
	* Create pre-signed url infraestructure example 
	* Create unit tests for pkgs 
	* Create a simple application using javascript for testing the solution 

## 🧠 Want to contribute?

PRs and issues are always welcome, also, make sure to check the roadmap section in this file. 
