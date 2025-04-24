# Project Overview

This project provides a microservices-based platform for **crowdsourced navigation** and **data aggregation**. It comprises the following services:

- **Navigation Service**: Routes users to their destinations using OSRM.
- **Crowdsourcing Service**: Collects real-time user feedback and location data.
- **Aggregation Service**: Processes and clusters the collected crowdsourced data using the Louvain community detection algorithm.
- **OSRM Service**: Provides high-performance routing via the `hamirho/osrm:trafficless` Docker image.
- **Authentication Service**: Secures API endpoints and manages user identities (ensure this is running or reachable).

---

## 🚀 Features

- **Real-time route calculation** with OSRM.
- **Crowdsourced data collection** for traffic, POIs, and user feedback.
- **Graph-based clustering** of crowdsourced events using Louvain.
- **Modular microservices** architecture for scalability and maintainability.
- **Out-of-the-box Docker support** using Docker Compose.

---

## 📦 Technologies

- **Docker & Docker Compose** for containerization.
- **Go  / Python** (service-specific stacks).
- **OSRM** for routing engine.
- **Louvain algorithm** for community detection.
- **RESTful APIs** for inter-service communication.

---


## 🛠 Prerequisites

- Docker Engine (>= 20.10)
- Docker Compose (>= 1.29)
- (Optional) A running Authentication Service exposing a REST API

---

## 🔧 Installation & Setup

1. **Clone the repository**


2. **Build and start services**
   ```bash
   docker-compose up --build -d
   ```
   This command will:
   - Build `navigation`, `crowdsourcing`, and `aggregation` from local Dockerfiles
   - Pull and run the OSRM image (`hamirho/osrm:trafficless`)
   - Start all containers in detached mode

3. **Verify**
   ```bash
   docker-compose ps
   ```
   You should see:
   - `navigation` listening on port **8080**
   - `crowdsourcing` listening on port **8081** (container port 8080)
   - `osrm` listening on port **5000**
   - `aggregation` restarting on exit to process new data



---

## 🧮 Aggregation Service & Louvain Algorithm

The **Aggregation Service**:

- Reads crowdsourced event data (from a database ).
- Constructs a graph where nodes represent events/locations and edges represent proximity or similarity.
- Applies the **Louvain community detection** algorithm to identify clusters (communities) of related events.
- Outputs cluster assignments or summary statistics for downstream consumption.

### Running Aggregation Manually

If you need to trigger aggregation on-demand:
```bash
docker-compose run --rm aggregation
```


---

## 📈 Monitoring & Logs

```bash
# View live logs
docker-compose logs -f

# View logs for a single service
docker logs -f navigation
```



