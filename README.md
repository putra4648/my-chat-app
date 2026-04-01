# Real-Time Chat Application

A high-performance, real-time chat application built with **Go** and the **Fiber v3** framework. This application features secure user authentication, private messaging, a real-time friend list, and persistent message history.

## 🚀 Key Features

- **Real-Time Messaging**: Bi-directional communication using WebSockets (Fiber v3 extension).
- **Private Conversations**: Instant 1-on-1 chats between friends.
- **Friend List**: Dynamic sidebar to find and start chats with other users.
- **Message History**: Persistent storage of all conversations using PostgreSQL.
- **Secure by Design**: 
  - CSRF Protection (Token-based).
  - Session Management with Redis.
  - Secure Password Hashing (Bcrypt).
  - Auth Middleware for route protection.
- **Modern UI**: Clean, responsive dashboard designed with Bootstrap 5 and Bootstrap Icons.
- **Responsive Layout**: Mobile-first design that adapts seamlessly to all screen sizes.

## 🛠️ Technology Stack

- **Backend**: [Go](https://go.dev/) (Golang)
- **Framework**: [Fiber v3](https://docs.gofiber.io/)
- **WebSocket**: [Fiber v3 WebSocket Extension](https://github.com/gofiber/contrib/tree/main/websocket)
- **Database**: [PostgreSQL](https://www.postgresql.org/) (via `pgxpool`)
- **Session Cache**: [Redis](https://redis.io/)
- **Dependency Injection**: [Parsley](https://github.com/matzefriedrich/parsley)
- **Frontend**: HTML5, Vanilla JavaScript, CSS3, Bootstrap 5

## 📋 Prerequisites

- **Go 1.21+** installed.
- **PostgreSQL** instance running.
- **Redis** instance running.

## 📦 Getting Started

1. **Clone the repository**:
   ```bash
   git clone https://github.com/putra4648/my-chat-app.git
   cd my-chat-app
   ```

2. **Environment Variables**:
   Create a `.env` file in the root directory and configure your connection strings:
   ```env
   DATABASE_URL=postgres://user:password@localhost:5432/chat_db
   REDIS_URL=localhost:6379
   PORT=3000
   ```

3. **Install Dependencies**:
   ```bash
   go mod download
   ```

4. **Initialize Database**:
   Run the migration SQL located in `migrations/001_create_table.sql`.

5. **Run the Application**:
   ```bash
   go run cmd/main.go
   ```

6. **Access the App**:
   Open `http://localhost:3000` in your browser.

## 🏗️ Project Architecture

The project follows a clean, modular structure:
- **/cmd**: Application entry point.
- **/internal/app**: Main application setup and dependency orchestration.
- **/internal/configs**: Modular middleware configurations (CORS, CSRF, Session, etc.).
- **/internal/handlers**: Route handlers organized by domain (Auth, Dashboard, WS).
- **/internal/repositories**: Data access layer using PostgreSQL.
- **/internal/services**: Business logic layer.
- **/internal/middleware**: Custom security and system middlewares.
- **/views**: HTML templates and frontend assets.

---
Built with ❤️ by [Antigravity](https://github.com/google/advanced-agentic-coding)
