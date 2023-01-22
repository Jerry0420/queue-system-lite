# Queue System

# Table Of Contents
- [Queue System](#queue-system)
- [Table Of Contents](#table-of-contents)
- [Overview](#overview)
- [Screenshots](#screenshots)
- [Architecture](#architecture)
  - [Services Inside The K8s Cluster](#services-inside-the-k8s-cluster)
    - [Backend](#backend)
    - [Frontend](#frontend)
    - [gRPC](#grpc)
  - [Services Outside The K8s Cluster](#services-outside-the-k8s-cluster)
    - [PostgreSQL](#postgresql)
    - [Vault](#vault)
  - [Authorization For Services](#authorization-for-services)
    - [Container User](#container-user)
    - [PostgreSQL Database User](#postgresql-database-user)
    - [Vault Server](#vault-server)
- [Backend Project Layout](#backend-project-layout)
- [Testing](#testing)

# Overview
Queue System is an online scheduling service for small and medium-sized businesses. It provides services for stores to manage their queues and their customers' states.

Once a store owner signup an account and sign in to Queue System, it will display a QRcode and the current state of all queues. When a customer enters a store, the customer will scan the QRcode and send out information to join a line. Once a customer scans a QRcode, the QRcode will change automatically. This mechanism can prevent unfair queuing and protect the other customers that haven't been in the queue yet. A store owner can change every single customer's status at any time. When a new customer joins the line or the state of a customer is updated, the website will refresh immediately and automatically.

Every account will persist for 24 hrs. After 24 hrs, Queue System will suspend the account (a.k.a. close store). Once the store is closed, Queue System will send a CSV report to the store owner and delete its related data in the database. This CSV file contains all customers' detailed information to let the store owner implement further business analysis. Next time, the store owner can signup an account again (a.k.a. open store) to use Queue System.

# Screenshots

![](https://raw.githubusercontent.com/Jerry0420/queue-system/main/images/open-store.png)
* Open Store (signup an account). 
* Every account will persist for 24 hrs. After 24 hrs, Queue System will suspend the account (a.k.a. close store). Next time, the store owner can signup an account again (a.k.a. open store) to use Queue System.
<br /><br />

![](https://raw.githubusercontent.com/Jerry0420/queue-system/main/images/signin-store.png)
* Sign in to the account
<br /><br />

![](https://raw.githubusercontent.com/Jerry0420/queue-system/main/images/close-store.png)
* The store can be actively closed by the store owner or passively closed by Queue System. 24 hrs after opening, Queue System will suspend the account (a.k.a. close store). 
<br /><br />

![](https://raw.githubusercontent.com/Jerry0420/queue-system/main/images/csv-content.png)
* Once the store is closed, Queue System will send a CSV report to the store owner. This CSV file contains all customers' detailed information to let the store owner implement further business analysis.   
<br /><br />

![](https://raw.githubusercontent.com/Jerry0420/queue-system/main/images/store-summary.png)
* The store owner can see the summary of all queues and the QRcode block after signing into their account.
<br /><br />

![](https://raw.githubusercontent.com/Jerry0420/queue-system/main/images/customer-scan.png)
* Once a customer scans a QRcode, the QRcode will change automatically.    
<br /><br />

![](https://raw.githubusercontent.com/Jerry0420/queue-system/main/images/update-customer.png)
* The store owner can change every single customer's status at any time.
<br /><br />

![](https://raw.githubusercontent.com/Jerry0420/queue-system/main/images/update-customer-refresh.png)
* When a new customer joins the line or the state of a customer is updated, the website will refresh immediately and automatically.   

# Architecture
![](https://raw.githubusercontent.com/Jerry0420/queue-system/main/images/architecture.png)
Queue System setups a Kubernetes cluster inside an AWS EC2 instance with MicroK8s and uses `Nginx Ingress Controller` as a Load-Balancer to expose Queue System service to the Content Delivery Network (CDN).

*p.s. Considering the cost, I used MicroK8s to create the k8s cluster and Nginx Ingress Controller as load-balancer instead of using AWS EKS and AWS ELB.*

Deploy the following resources inside the K8s cluster:
* Deployments:
  * Backend
  * Frontend
  * gRPC
* Services
  * Backend
  * Frontend
  * gRPC
* CronJob
  * Send a REST API to the backend server every minute to check if there is any store that has already opened 24 hrs.
  * Queue System will suspend these selected accounts and send CSV reports to these store owners that contain all customers' detailed information.
* Ingress
  * Nginx Ingress Controller
  * Ingress

## Services Inside The K8s Cluster

### Backend
Developed by `Golang`.   
Provide REST API and [Server-Sent Event API](https://developer.mozilla.org/en-US/docs/Web/API/Server-sent_events/Using_server-sent_events) interfaces.

### Frontend
Developed by `ReactJS` with `Typescript`.   
Interact with backend by REST API and [Server-Sent Event API](https://developer.mozilla.org/en-US/docs/Web/API/Server-sent_events/Using_server-sent_events) interfaces.

### gRPC
Developed by `Golang`.   
The backend server passes tasks of sending emails and generating CSV files to the gRPC server.

## Services Outside The K8s Cluster

Deploy the following two services (PostgreSQL, Vault) outside the K8s cluster:

*p.s. Considering the cost, I deployed PostgreSQL and Vault by myself rather than using AWS RDS and AWS Secrets Manager.*

### PostgreSQL
The PostgreSQL server stores all data.

### Vault
The Vault server connects to the PostgreSQL server and rotates the user every hour. Then the backend server connects to the Vault server to retrieve the account and password of the PostgreSQL server.

## Authorization For Services

### Container User
Create a new user (`appuser`, UID=1000) and a new group (`appgroup`, GID=1001) inside Dockerfiles of Backend, Frontend, and gRPC. Then, change permissions of all data and Volumes that are inside these Docker images to be read, edited, and executed by `appuser`.

### PostgreSQL Database User
When setting the PostgreSQL Database, use the user root to create the user `migration` and the user `vault`.
* user migration:
  * Use this user to connect to the PostgreSQL Server only when migrating the database.
  * This user can create, read, update and delete all data inside all tables.
  * This user can create, alter and delete all tables' schemas.   
* user vault: 
  * Use this user to connect to the PostgreSQL Server only when using the Vault Server.
  * This user's role is superuser.
  * This user can create other database users.  

### Vault Server
The Vault Server uses the user vault to connect to the PostgreSQL Server and create new database users.
  * All newly created users can create, read, update and delete all data inside all tables.
  * Every newly created user will auto-expire by the PostgreSQL Server after 1 hr.
  * The Vault Server will create a new database user every 1 hr. 

# [Backend Project Layout](https://github.com/Jerry0420/queue-system/tree/main/backend)
*For simplicity, I put the backend, the frontend, and the gRPC server in one repo.*

* All file names that postfix with `_test.go` are golang test files.
* Reference the following two websites for implementation of clean architecture in this Go project.
  * https://github.com/bxcodec/go-clean-arch
  * https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html
* Separate this Golang project into four parts according to clean architecture:
  * [Delivery](https://github.com/Jerry0420/queue-system/tree/main/backend/delivery/httpAPI)
    * Define all HTTP routes.
    * Validate parameters of HTTP requests.
    * Build HTTP responses.
  * [Domain](https://github.com/Jerry0420/queue-system/tree/main/backend/domain)
    * Define all data models.
  * [Repository](https://github.com/Jerry0420/queue-system/tree/main/backend/repository)
    * Interact with the PostgreSQL server and the gRPC server.
  * [Usecase](https://github.com/Jerry0420/queue-system/tree/main/backend/usecase)
    * Clean all data from the PostgreSQL server.
    * Implement all business logic.

# Testing
* Unit Tests
  * All file names that postfix with `_test.go`
  * Separately unit test `Delivery`, `Usecase` and `Repository` with Mock and Dependency Injection.
* [Integration Tests](https://github.com/Jerry0420/queue-system/tree/main/backend/integration_tests)
  * Write a [mock grpc server](https://github.com/Jerry0420/queue-system/tree/main/backend/integration_tests/mock_grpc).
  * Test the whole workflow by user stories.
