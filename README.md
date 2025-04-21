# API Arfevrier

## Project Overview
This project provides a RESTful API for various functionalities, including YouTube video management, WebSocket connections, and Bitcoin price retrieval. It is designed to serve as a backend for applications requiring these services.

## Run
To start the API server, use the following command:
```bash
$ source .env
$ pm2 start --name api-v3 run.sh
```
The server runs on `127.0.0.1:1239` and is accessed via a reverse proxy.

## Features
- **YouTube API**: Fetch subscription videos and download YouTube content (video/audio).
- **WebSocket API**: Establish and manage WebSocket connections for real-time communication.
- **Bitcoin API**: Retrieve the current Bitcoin price from blockchain.info.

## Requirements
- Go 1.18 or later
- Gin framework
- Swaggo for API documentation
- Additional dependencies listed in `go.mod`

## Usage
### Generate Documentation
To generate Swagger documentation, run:
```bash
$ swag init
```

### API Endpoints
- **YouTube API**
  - `GET /youtube/subscriptions/:token`: Fetch subscription videos.
  - `GET /youtube/download/:type/:id`: Download YouTube video or audio.
- **WebSocket API**
  - `GET /webconnect/new/:channel`: Create a new WebSocket connection.
  - `GET /webconnect/connect/:id`: Connect to an existing WebSocket.
- **Bitcoin API**
  - `GET /bitcoin/price`: Retrieve the current Bitcoin price.
