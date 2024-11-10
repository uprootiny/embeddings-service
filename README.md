# Embeddings Service

Welcome to the `embeddings-service` — a Go-based server designed for scanning, indexing, and mapping intents to project entry points. This service can seamlessly integrate into various projects, supporting use cases in software gardening, financial intelligence, chaotic service discovery, and LLM experimentation.

## Overview

The `embeddings-service` is built to aid in the discovery and navigation of project repositories and their associated entry points, connecting intents and user requests to the most relevant codebases. It scans local working directories, identifies key scripts and configurations, and maps them for rapid and meaningful interaction.

## Key Features

- **Repository Scanning**: Discover and list project repositories from specified local directories.
- **Intent Mapping**: Use embeddings to map user intents to relevant projects and execution points.
- **Detailed Project Insights**: Retrieve information on project entry points, configurations, and modification times.
- **RESTful API Endpoints**: A structured API to access all the service capabilities programmatically.
- **Service Logging**: Comprehensive logs to track API usage and diagnostics.

## Workflows and Use Cases

### 1. **Software Gardening**
In large codebases, maintaining awareness of what's present and its state can be challenging. `embeddings-service` can help by regularly scanning and exposing the current projects and their health.

**Example API Calls:**

- **List Recent Projects**:
    ```bash
    curl http://localhost:8085/list-repos
    ```
    **Response**:
    ```json
    [
      {
        "name": "diagnostics",
        "path": "/home/uprootiny/Projects/diagnostics",
        "entryPoints": ["main.go", "web.clj"],
        "lastModified": "2024-11-10T13:15:56Z"
      },
      ...
    ]
    ```

- **Get Project Details**:
    ```bash
    curl "http://localhost:8085/repo-details?path=/home/uprootiny/Projects/diagnostics"
    ```
    **Response**:
    ```json
    [
      "-rw-r--r-- 1 user group 2034 Nov 10 2024 main.go",
      "drwxr-xr-x 2 user group 4096 Nov 10 2024 src/",
      ...
    ]
    ```

### 2. **Financial Intelligence**
For financial data analysts and traders, integrating different data sources, scraping scripts, and analysis tools is crucial. Use `embeddings-service` to manage and link these disparate components effectively.

**Example Workflow**:
- **Identify Available Scrapers and Run Them**:
    ```bash
    curl http://localhost:8085/list-repos | jq '.[] | select(.name | contains("scraper"))'
    ```

- **Map Intent to Execute Financial Data Analysis**:
    ```bash
    curl -X POST -H "Content-Type: application/json" -d '{"intent": "Analyze tech news impact"}' http://localhost:8085/map-intent
    ```
    **Response**:
    ```json
    {
      "mappedProject": "news_scraper",
      "entryPoints": ["fetch_news.go"],
      "params": "news_params.json"
    }
    ```

### 3. **Service Discovery in a Chaotic Environment**
In environments with many interdependent services, keeping track of what's running and what isn't can be a headache. `embeddings-service` offers APIs to track projects, discover entry points, and analyze logs.

**Example API Calls**:
- **Check for Running Services**:
    ```bash
    curl http://localhost:8085/network-services
    ```
    **Response**:
    ```json
    [
      "Controller: *:8080",
      "Main: *:5000",
      "Postgres: 127.0.0.1:5432",
      ...
    ]
    ```

### 4. **LLM Experimentation and Rapid Prototyping**
Experiment with embeddings and language models to map user intents to code execution paths, allowing for quick integration and prototyping of LLM-powered features.

**Example Workflow**:
- **Map Intent Using LLMs**:
    ```bash
    curl -X POST -H "Content-Type: application/json" -d '{"intent": "Generate market trend forecast"}' http://localhost:8085/map-intent
    ```
    **Response**:
    ```json
    {
      "mappedProject": "market_forecaster",
      "entryPoints": ["forecast.py"],
      "params": "trend_params.json"
    }
    ```

- **Execute CLI Commands Directly**:
    ```bash
    curl http://localhost:8085/execute?cmd="python3 /path/to/market_forecaster/forecast.py"
    ```
    **Response**:
    ```json
    {
      "status": "success",
      "output": "Forecast generated successfully. Output file: forecast_results.csv"
    }
    ```

## How to Use

### Running the Service
1. **Clone the repository**:
    ```bash
    git clone git@github.com:uprootiny/embeddings-service.git
    cd embeddings-service
    ```

2. **Build and run the service**:
    ```bash
    go build -o embeddings-service .
    ./embeddings-service
    ```

3. **Access the API**:
    Visit `http://localhost:8085` or use `curl` commands to interact with the API endpoints.

### API Endpoints

- **`/list-repos`**: Lists recently modified repositories.
- **`/repo-details?path=/uprootiny/embeddings-service`**: Provides details for a specific repository.
- **`/map-intent`**: Maps a user-provided intent to relevant projects and entry points.
  `curl -X GET http://localhost:8085/map-intent`
- **`/network-services`**: Lists active network services on the server.
- **`/execute?cmd=your-command`**: Executes a command on the server (use with caution).

## Advanced Features

- **Embeddings Integration**: The service is built to leverage embeddings to map high-level user intents to specific projects and code paths.
- **Modular Design**: The project can be extended with additional handlers, models, and scanning capabilities.
- **Detailed Logs**: Keep track of all interactions and API calls through robust logging.

## Roadmap and Future Work

- **Full Embeddings Implementation**: Integrate with embeddings libraries for more advanced intent mapping.
- **PWA Interface**: Develop a Progressive Web App for interactive use of the service.
- **Security Enhancements**: Add authentication and authorization layers for executing sensitive commands.
- **Monitoring Dashboards**: Build a UI to track running services, resource usage, and project statuses in real time.

## Contributing

Contributions are welcome! Feel free to submit pull requests or open issues for features, bug fixes, or enhancements.

---

**Explore, connect, and manage your projects efficiently with `embeddings-service`.** Your bridge to better code discovery and integration in a complex ecosystem.
