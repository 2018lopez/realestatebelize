<div align="center">

  <h1>RealEstate Belize Listing API</h1>
  
  <p>
    Advance Web Development Project
  </p>
  
<h4> Developed By: Imer Lopez </h4>
</div>
<!-- Table of Contents -->
# :notebook_with_decorative_cover: Table of Contents

- [About the Project](#star2-about-the-project)
  * [Tech Stack](#space_invader-tech-stack)
  * [Features](#dart-features)
- [Getting Started](#toolbox-getting-started)
  * [HealthCheck Endpoint](#gear-health-check-endpoint)
  * [Users Endpoints](#gear-users-endpoints)
  * [Listings Endpoints](#gear-listings-endpoints)
  * [Report Endpoints](#gear-reports-endpoints)
  * [Currency Rate Endpoint](#gear-currency-rate-endpoint)
  * [Server File Endpoint](#gear-server-file-endpoint)


<!-- About the Project -->
## :star2: About the Project
<p> My final project proposal is to develop a Real Estate API. Over the years, the real
estate industry has been growing fast, from international investors and locals finding new
investment opportunities. The API will have provides properties listing within the country. </p>

<!-- Features -->
### :dart: Features

- Creating, Updating Users accounts (agents/admin) - users will upload a profile image as part of account creation
- Assigning agent to listings/properties
- User's Account Activation - a token is provide via email when the account is created
- User Authentication to access certain part of API that required access permission
- Creating, Updating listing - Images upload will be part of listing creation
- Reports - Top Agents - Listing Status -> Available vs Sold/leased
- Currency Rate Third Party API


<!-- TechStack -->
### :space_invader: Tech Stack


<details>
  <summary>Server</summary>
  <ul>
    <li><a href="https://go.dev/">Golang</a></li>
   </ul>
</details>

<details>
<summary>Database</summary>
  <ul>
    <li><a href="https://www.postgresql.org/">PostgreSQL</a></li>
   
  </ul>
</details>



<!-- Getting Started -->
## 	:toolbox: Getting Started

<!-- Healthcheck Endpoint -->
### :gear: Health Check Endpoint

Health Check End point
```bash
 GET : /v1/healthcheck
```

<!-- Users Endpoint -->
### :gear: Users Endpoints

User End points
```bash
 POST : /v1/users
```
```bash
 GET : /v1/users/:id
```
```bash
 POST : /v1/users/image
```
```bash
 PUTH : /v1/users/updated/:id
```
```bash
 PUTH : /v1/users/activated
```
```bash
 POST: /v1/tokens/authentication
```

<!-- Listings -->
### :gear: Listings Endpoints

Listings Endpoints
```bash
 POST: /v1/listings
```
```bash
 GET: /v1/listings
```
```bash
 GET: /v1/listings/:id
```
```bash
 POST: /v1/listings/images
```
```bash
 POST: /v1/users/listings
```
```bash
 PUTH: /v1/listings/update/:id
```
```bash
 GET: /v1/agent/listings/:id
```

<!-- REPORTS -->
### :gear: Reports Endpoints

Reports Endpoints

```bash
 GET: /v1/report/agents
```
```bash
 GET: /v1/report/listings
```
```bash
 GET: /v1/report/total-sales
```



<!-- CURRENCY RATE -->
### :gear: Currency Rate Endpoint

Currency Endpoint 

```bash
 GET: /v1/currencyrate/:id
```

<!-- Server File Endpoint -->
### :gear: Server File Endpoint

Server File Endpoint
```bash
 GET: /uploads/*filepath
```

