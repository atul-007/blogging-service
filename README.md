﻿# Blogging Service

This is a microservices-based blogging service implemented using Go Fiber, Elasticsearch, Docker Compose, and Kubernetes. The service allows users to submit blog entries and search for blog entries based on their content.

## Table of Contents
1. [Architecture](#architecture)
2. [Services](#services)
   - [Blog Submission API](#blog-submission-api)
   - [Queue Service](#queue-service)
   - [Queue Consumer](#queue-consumer)
   - [Search API](#search-api)
3. [Running with Docker Compose](#running-with-docker-compose)
4. [Deploying to Kubernetes](#deploying-to-kubernetes)
5. [Testing](#testing)

## Architecture

![blog-submission-design](./blog-submission-design.png)

The service consists of four main components:
1. **Blog Submission API**: Accepts blog submissions and enqueues them for asynchronous processing.
2. **Queue Service**: Manages the queue of blog submissions.
3. **Queue Consumer**: Consumes the queued blog submissions and writes them to Elasticsearch.
4. **Search API**: Allows searching of blog entries in Elasticsearch.

## Services

### Blog Submission API

- **Endpoint**: `/submit`
- **Method**: POST
- **Payload**:
  ```json
  {
    "title": "My First Blog",
    "text": "This is the content of my first blog",
    "user_id": "user1"
  }
  ```
- **Response**:
  ```json
  {
    "status": "Blog submission queued"
  }
  ```

### Queue Service

- **Enqueue Endpoint**: `/enqueue`
  - **Method**: POST
  - **Payload**:
    ```json
    {
      "title": "My First Blog",
      "text": "This is the content of my first blog",
      "user_id": "user1"
    }
    ```
  - **Response**:
    ```json
    {
      "status": "Blog enqueued"
    }
    ```

- **Dequeue Endpoint**: `/dequeue`
  - **Method**: GET
  - **Response**:
    ```json
    {
      "title": "My First Blog",
      "text": "This is the content of my first blog",
      "user_id": "user1"
    }
    ```

### Queue Consumer

The queue consumer service reads from the queue and writes the entries to an Elasticsearch database. It continuously polls the queue for new blog submissions.

### Search API

- **Endpoint**: `/search`
- **Method**: GET
- **Query Parameters**:
  - `q`: The search query string.
- **Response**:
  ```json
  [
    {
      "title": "My First Blog",
      "text": "This is the content of my first blog"
    }
  ]
  ```

## Running with Docker Compose

To run the services locally using Docker Compose:

1. Clone the repository:

   ```sh
   git clone https://github.com/atul-007/blogging-service.git
   cd blogging-service
   ```

2. Build and start the services:

   ```sh
   docker-compose up --build
   ```

3. The services will be available at the following ports:
   - Blog Submission API: `http://localhost:3000`
   - Queue Service: `http://localhost:3001`
   - Search API: `http://localhost:3002`
   - Elasticsearch: `http://localhost:9200`

## Deploying to Kubernetes

To deploy the services to a Kubernetes cluster:

1. Ensure you have a Kubernetes cluster running and `kubectl` configured.

2. Apply the deployment files:

   ```sh
   kubectl apply -f kubernetes/elasticsearch-deployment.yaml
   kubectl apply -f kubernetes/queue-consumer-deployment.yaml
   kubectl apply -f kubernetes/blog-submission-deployment.yaml
   kubectl apply -f kubernetes/search-api-deployment.yaml
   ```

3. Expose the services using a load balancer or ingress as required.

## Testing


https://github.com/atul-007/blogging-service/assets/94912544/d9f200c1-ec60-4b42-a36e-d23560294c0d



### Blog Submission

Submit a blog entry:

```sh
curl -X POST http://localhost:3000/submit -H "Content-Type: application/json" -d '{"title": "My First Blog", "text": "This is the content of my first blog", "user_id": "user1"}'
```

### Search Blog Entries

Search for blog entries containing a specific term:

```sh
curl http://localhost:3002/search?q=First
```

### Elasticsearch

Verify the entries in Elasticsearch:

```sh
curl http://localhost:9200/blogs/_search?q=First
```

---

