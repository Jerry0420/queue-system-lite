# Queue System

# Table Of Contents
- [Queue System](#queue-system)
- [Table Of Contents](#table-of-contents)
- [Overview](#overview)
- [Screenshots](#screenshots)
- [Architecture](#architecture)
  - [Services](#services)
    - [Backend](#backend)
    - [Frontend](#frontend)
    - [PostgreSQL](#postgresql)
- [Backend Project Layout](#backend-project-layout)
- [Testing](#testing)

# Overview
Queue System is an online scheduling service for small and medium-sized businesses. It provides services for stores to manage their queues and their customers' states.

Once a store owner signup an account and sign in to Queue System, it will display a QRcode and the current state of all queues. When a customer enters a store, the customer will scan the QRcode and send out information to join a line. Once a customer scans a QRcode, the QRcode will change automatically. This mechanism can prevent unfair queuing and protect the other customers that haven't been in the queue yet. A store owner can change every single customer's status at any time. When a new customer joins the line or the state of a customer is updated, the website will refresh immediately and automatically.

Every account will persist for 24 hrs. After 24 hrs, Queue System will suspend the account (a.k.a. close store). Once the store is closed, Queue System will send a CSV report to the store owner and delete its related data in the database. This CSV file contains all customers' detailed information to let the store owner implement further business analysis. Next time, the store owner can signup an account again (a.k.a. open store) to use Queue System.

# Screenshots

![](https://raw.githubusercontent.com/Jerry0420/queue-system-lite/main/images/open-store.png)
* Open Store (signup an account). 
* Every account will persist for 24 hrs. After 24 hrs, Queue System will suspend the account (a.k.a. close store). Next time, the store owner can signup an account again (a.k.a. open store) to use Queue System.
<br /><br />

![](https://raw.githubusercontent.com/Jerry0420/queue-system-lite/main/images/signin-store.png)
* Sign in to the account
<br /><br />

![](https://raw.githubusercontent.com/Jerry0420/queue-system-lite/main/images/close-store.png)
* The store can be actively closed by the store owner or passively closed by Queue System. 24 hrs after opening, Queue System will suspend the account (a.k.a. close store). 
<br /><br />

![](https://raw.githubusercontent.com/Jerry0420/queue-system-lite/main/images/csv-content.png)
* Once the store is closed, Queue System will send a CSV report to the store owner. This CSV file contains all customers' detailed information to let the store owner implement further business analysis.   
<br /><br />

![](https://raw.githubusercontent.com/Jerry0420/queue-system-lite/main/images/store-summary.png)
* The store owner can see the summary of all queues and the QRcode block after signing into their account.
<br /><br />

![](https://raw.githubusercontent.com/Jerry0420/queue-system-lite/main/images/customer-scan.png)
* Once a customer scans a QRcode, the QRcode will change automatically.    
<br /><br />

![](https://raw.githubusercontent.com/Jerry0420/queue-system-lite/main/images/update-customer.png)
* The store owner can change every single customer's status at any time.
<br /><br />

![](https://raw.githubusercontent.com/Jerry0420/queue-system-lite/main/images/update-customer-refresh.png)
* When a new customer joins the line or the state of a customer is updated, the website will refresh immediately and automatically.   

# Architecture
* Services
  * Backend
  * Frontend
  * CronJob
    * Send a REST API to the backend server every minute to check if there is any store that has already opened 24 hrs.
    * Queue System will suspend these selected accounts and send CSV reports to these store owners that contain all customers' detailed information.
* Storage
  * PostgreSQL
* Reverse Proxy
  * Nginx

## Services

### Backend
Developed by `Golang`.   
Provide REST API and [Server-Sent Event API](https://developer.mozilla.org/en-US/docs/Web/API/Server-sent_events/Using_server-sent_events) interfaces.

### Frontend
Developed by `ReactJS` with `Typescript`.   
Interact with backend by REST API and [Server-Sent Event API](https://developer.mozilla.org/en-US/docs/Web/API/Server-sent_events/Using_server-sent_events) interfaces.

### PostgreSQL
The PostgreSQL server stores all data.

# [Backend Project Layout](https://github.com/Jerry0420/queue-system-lite/tree/main/backend)
*For simplicity, I put the backend and the frontend in one repo.*

* All file names that postfix with `_test.go` are golang test files.
* Reference the following two websites for implementation of clean architecture in this Go project.
  * https://github.com/bxcodec/go-clean-arch
  * https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html
* Separate this Golang project into four parts according to clean architecture:
  * [Delivery](https://github.com/Jerry0420/queue-system-lite/tree/main/backend/delivery/httpAPI)
    * Define all HTTP routes.
    * Validate parameters of HTTP requests.
    * Build HTTP responses.
  * [Domain](https://github.com/Jerry0420/queue-system-lite/tree/main/backend/domain)
    * Define all data models.
  * [Repository](https://github.com/Jerry0420/queue-system-lite/tree/main/backend/repository)
    * Interact with the PostgreSQL server.
  * [Usecase](https://github.com/Jerry0420/queue-system-lite/tree/main/backend/usecase)
    * Clean all data from the PostgreSQL server.
    * Implement all business logic.

# Testing
* Unit Tests
  * All file names that postfix with `_test.go`
  * Separately unit test `Delivery`, `Usecase` and `Repository` with Mock and Dependency Injection.
* [Integration Tests](https://github.com/Jerry0420/queue-system-lite/tree/main/backend/integration_tests)
  * Test the whole workflow by user stories.
