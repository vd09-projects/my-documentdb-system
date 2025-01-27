# My DocumentDB System

A Go + MongoDB application that demonstrates:

1. **User Registration and Authentication** (JWT-based).
2. **CSV/JSON Upload and Parsing**, with support for record types and field tracking.
3. **Data Quarantine** for invalid records, plus optional date-range filtering.
4. **Aggregation Operations** (e.g., SUM and AVERAGE) on numeric fields.

---

## Table of Contents

1. [Project Overview](#project-overview)
2. [Features](#features)
3. [Setup & Installation](#setup--installation)
4. [Running the Project (Local & Docker)](#running-the-project-local--docker)
5. [Usage Flow](#usage-flow)
6. [Future Improvements](#future-improvements)
7. [Flow Diagram](#flow-diagram)

---

## Project Overview

This repository builds a cloud-ready proof-of-concept system for ingesting data from multiple sources (CSV/JSON), storing them with minimal schema constraints, and offering basic data transformations and aggregations.

**Core technologies**:
- **Go (Golang)** for the backend and services.
- **MongoDB** as our NoSQL document store.
- **Docker & Docker Compose** to run and link the services.
- **JWT** for authentication.
- **HTML/Vanilla JS** for a minimal user interface, enabling:
  - **Login/Register**
  - **File Upload** (CSV/JSON + record type)
  - **Dashboard** to view filtered user data
  - **Aggregation** (Sum/Average) on numeric fields

---

## Features

- **User Management**:
  - Register with username/password (**bcrypt-hashed**)
  - Login returning a **JWT token**
  - Protected Endpoints require `Authorization: Bearer <token>`

- **Upload/Quarantine**:
  - CSV or JSON upload
  - Valid Records stored in `valid_records` (with a user-specified record type)
  - Invalid Records stored in `quarantine_records` with a “reason”

- **Record Type & Field Tracking**:
  - On each upload, the system captures field names and upserts them into a `record_fields` collection so we know which fields exist for each (user, recordType) pair.

- **Date Filtering**:
  - Each record is timestamped.
  - Users can specify `?from=YYYY-MM-DD&to=YYYY-MM-DD` to filter data in the Dashboard view.

- **Aggregation**:
  - A separate **Aggregator** page lets users pick:
    - A record type (e.g., “sales,” “inventory”)
    - A field (e.g., “price,” “quantity”)
    - An operation (SUM/AVERAGE)
  - Result is computed either via a small in-memory approach or custom pipeline, returning the final numeric aggregate.

---

## Setup & Installation

1. **Install Go** (>= 1.23).

2. **Install Docker** and (optionally) Docker Compose.

3. **Clone this repository**:
    ```bash
    git clone https://github.com/YOUR_ORG/my-documentdb-system.git
    cd my-documentdb-system
    ```

4. **Check `go.mod` and run `go mod tidy` if needed**:
    ```bash
    go mod tidy
    ```

---

## Running the Project (Local & Docker)

### 1. Local Build & Run
```bash
make build       # Compiles the Go binary into bin/my-documentdb-system
make run         # Runs the Go server locally on :8080
make docker-run  # Starts server and docker
```

### 2. Accessing via AI

To access the AI system, follow the steps below:

1. Open your browser.
2. Navigate to the following URL: [http://localhost:8080/](http://localhost:8080/)
3. The link will direct you to the user interface (UI) where you can interact with the system.

Ensure that the application is running locally before accessing the URL.

### 3. MongoDB Shell (Optional)
If you need to interact with the MongoDB instance directly, you can open a mongosh session inside the my-mongo-db container:
```bash
 make mongo-shell
```
Once inside the shell, you can run commands like:
```bash
- use myDB
- db.valid_records.find().pretty()
- db.quarantine_records.find().pretty()
```

---

## Usage Flow

1. **Register a new user**:
   - Send a `POST /register` with JSON body:
     ```json
     {
       "username": "...",
       "password": "..."
     }
     ```
   - **Or** open the UI, click **Register**, and submit the form.

2. **Login**:
   - Send a `POST /login` with JSON body:
     ```json
     {
       "username": "...",
       "password": "..."
     }
     ```
   - **Or** in the UI, fill in the Login form.
   - The server returns a JWT token on success.

3. **Upload (CSV/JSON) with Record Type**:
   - Navigate to **Upload** (once logged in).
   - Choose a record type (e.g., “sales” or “inventory”).
   - Upload the CSV/JSON file.
   - Valid data goes to `valid_records`; invalid lines go to `quarantine_records`.

4. **Dashboard**:
   - Click **Dashboard** to see your data.
   - Optionally specify `from` / `to` date filters (e.g., `2023-01-01`) to limit results.
   - The UI calls `GET /userData?from=...&to=...` behind the scenes.

5. **Aggregator**:
   - Choose a record type from the dropdown.
   - A list of fields (collected from your data) will appear.
   - Pick an operation (SUM/AVERAGE).
   - Click **Calculate** to see a numeric result.

---

## Future Improvements

1. **Dynamic Record Types**:
   - Currently, record types are hardcoded (e.g., “sales,” “inventory”).
   - We can let users create/edit custom record types from the UI, storing them in the DB.

2. **Indexing**:
   - For large data sets, we’d want to index fields like `userID`, `recordType`, `timestamp`, or frequently aggregated fields to speed up queries.

3. **Roles & Permissions**:
   - Introduce user roles (e.g., admin, viewer, analyst) to restrict who can upload data vs. who can only view data.

4. **More Advanced Aggregations**:
   - The aggregator currently supports SUM and AVERAGE. We could add MIN, MAX, or a full MongoDB pipeline based on user-defined expressions.

5. **Schema Validation / AI**:
   - Expand the system to do dynamic schema validation or use an AI approach to interpret new data fields.
   - Possibly integrate more advanced anomaly detection or data transformations.

6. **Pagination & Sorting**:
   - For large result sets, add queries for pagination (e.g., `limit`, `offset`, or `page`) and sorting by date or numeric fields.

7. **Cloud Deployment**:
   - Migrate from local Docker Compose to a full cloud environment (AWS ECS, GCP, Azure) with managed MongoDB (Atlas) or fully serverless solutions.

8. **Production Hardening**:
   - Load secrets from environment variables or a secrets manager (Vault, AWS Secret Manager, etc.).
   - Configure HTTPS certificates.
   - Add logging (structured logs) and monitoring (Prometheus/Grafana).

## Flow Diagram
![ezcv logo](https://img.plantuml.biz/plantuml/png/fPLHRzis4CVVzIbkVZ1iB6asO4_6t3QEcY35qcHst3LOXvgbpZ9HcjH8oeQ7VVWTpveh4qeG87wnHDxzoV_lZjHR7uGBzLfd5VqhgIfX0imzkgE1tiJPBGt2_Be7mjFVyVILXHcw3JgUvST4uCQQqkOJdiydMJOUmnHBvGE9NcgxXV4uYyl2wMjI7y5jXGREtFyXbIWncr_JRAC-WhlRsNqhOb1JjX5hFA5WxxVM5SECOrapdupWsb1808DRYC6W32pYWr-0j5gZ3CgQB9_0QMkPMIqbUdzsAPG-w3MRxF6EfKCHi02_ZrpMMkz-w85rWCvPnrv_iwKoJnZLFF-eIvv-ZiMFGgj214mDfguu3cmCsV0ZcIZG12MqJrrECTJEmFi_xY7ORqZxK4lWh25xEw_3AgPAfd3E57rgyyxPZIgmhUA3fKNZ9hMLGi_ebVmF4m23A6-T-aT4tH5CK3YI_Vmhgtn-lT_3lD9M592B8DALFuYDAYEIO9kmJiJrNc5mCVjuzdvs-m5-18UI4D_lAjgKqA61nkjcGB_ExYBwWi4peIzx3R-8h1T_ry898cNmCDuXDA-uWejzo9UbGXq5jYgie51UsWr6Rljnpa-As-p4eQzru82qRsUjLWN5uQJugN7iUnbgZwIpIHodEh_yz3FCZMWgX0Kbce8BaWk_-cg58w9N6E1cnxKPFbMSGqfIdwjtiJoK51NUB9rTZn_eTtfpSvhd_M1RRyyeq-yrkXsarVcwunDpdKCVcBuJfiGkkDDzIXSERXkk_zAwbfqNmpmlIuU4oRed-BMYkW1gvSFqqD3-6QIcA8gq2zbPGLBAoD0w89rzYA_AL_1dpmkDldS2FMKsQ89rkzvVNG3kR5LhPd-GXSBPgh1RyqA6RMGpr7YqzyTeNRBbDzy3nv-Is_OIWrmr1n64SsB7VFZkxRTjnAPDKMRTj2gd94PQT5cOe4A4zHYelT81qh7FYvMr7c70DaxYTyRn5Jey0QcL_0ONGalMwG8WccR28_zzNd7yqBVLUE3hzRzYWtnfTu8nzHh93L3HQlGDXEJW3LWH7bA5EtI4SWsqo5cp22W3qXPaoJgeL8LQ6HnsVwsQom-4vrF9f-YTAVZuiDb8QWIvyfMcgeLLX6c4RlByuUljixTeYdgj_mC0)