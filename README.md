# Clamurai Traefik Plugin 

Wake up, clamurai...<br>
    &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;We have some malware to burn.

## The Project

Clamurai is a middleware built for providing inline anti-malware features into a vast range of solutions such as object storage (S3 or Blob for instance) and MFT solutions. 

The inspiration from the project comes from tools such as GuardDuty Malware Protection and Microsoft Defender for Storage, which rely on receiving an object before inspecting it(async with queuing), forcing workflows to use batch format while inspecting files - with clamurai, we are able to stop malicious input before it touches our environment. 

### Traefik Configuration  

| Option              | Type   | Description                              | Default            |
| ------------------- | ------ | ---------------------------------------- | ------------------ |
| `clamdAddress`         | string | Host\:port of the ClamAV daemon (clamd)  | `"localhost:3310"` |
| `alertMode` | bool | When set to `true` this feature will not generate blocks on up/downstream traffic, only log events       | `false`          |
| `clamdReadTimeout`       | int | Ammount of time in seconds for the client to wait for scan on remove clamd                | `90`             |
| `clamavConnectionTimeout`          | int | Ammount of time in seconds for the client to attempt closing connection with remote clamd  | `90`           |


## Usage

In this section there are a few use cases where implementing the clamurai midleware as a standalone service might come in usefull.

### Pre-requisites 

In order to follow or implement the aftermentioned use cases, it's important that you consider provision: 

* Kubernetes or docker cluster
* Network access or permissions into desired environments such as a cloud providers or on-premises infrastructure
* A pot filled with fresh brewed coffee 

### Protecting file uploads in updog instance
TBD 

### Protecting Apache Airavata-mft
TBD 

### Protecting file uploads in S3 pre-signed urls 
TBD

### Protecting file uploads in azure with Blob SAS tokens 
TBD

## 🧱 Roadmap

* Implement tests and documentations for the code 
	* Document use cases
* Improve detections and clamav customization 
	* Custom clam signature lib 
	* Create flag for alert and block mode 
* Tackle logging and outputs 
	* Improve logger pkg
	* Check if clamav logs can be displayed as json messing with the config file 
* Test the package   
	* Create unit tests for pkgs 
	* Create a simple application using javascript for testing the solution 

## 🧠 Want to contribute?

PRs and issues are always welcome, also, make sure to check the roadmap section in this file. 
