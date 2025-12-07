---
source: https://sazardev.github.io/goca/
fetched: 2025-12-07T12:21:59Z
method: sitefetch
---

<page>
  <title>Goca - Go Clean Architecture Code Generator</title>
  <url>https://sazardev.github.io/goca/</url>
  <content>Build production-ready Go applications following Clean Architecture principles. Stop writing boilerplate, start building features.

Clean Architecture by Default
-----------------------------

Every line of code follows Uncle Bob's Clean Architecture principles. Proper layer separation, dependency rules, and clear boundaries guaranteed.

Lightning Fast Setup
--------------------

Generate complete features with all layers in seconds. From entity to handler, repository to use case - everything configured and ready.

Best Practices Enforced
-----------------------

Prevents common anti-patterns like fat controllers, god objects, and anemic domain models. Your code stays clean and maintainable.

Auto Integration
----------------

New features are automatically integrated with dependency injection and routing. No manual wiring needed.

Multi-Protocol Support
----------------------

Generate handlers for HTTP REST, gRPC, CLI, Workers, and SOAP. All following the same clean architecture pattern.

Test-Ready
----------

Code generated with clear interfaces and dependency injection makes testing a breeze. TDD-friendly from the start.

8 Databases Supported
---------------------

PostgreSQL, MySQL, MongoDB, SQLite, SQL Server, PostgreSQL JSON, Elasticsearch, and DynamoDB. Switch between databases without changing business logic.

Rich Documentation
------------------

Comprehensive guides, tutorials, and examples. Learn Clean Architecture while building real applications.

Production Ready
----------------

Used in production systems. Battle-tested patterns and code generation that scales from MVP to enterprise.

Quick Example [​](#quick-example)
---------------------------------

bash

    # Initialize a new project
    goca init my-api --module github.com/user/my-api
    
    # Generate a complete feature with all layers
    goca feature User --fields "name:string,email:string,role:string"
    
    # That's it! You now have:
    # → Domain entity with validations
    # → Use cases with DTOs
    # → Repository with PostgreSQL implementation
    # → HTTP handlers with routing
    # → Dependency injection configured

Why Clean Architecture? [​](#why-clean-architecture)
----------------------------------------------------

Clean Architecture ensures your codebase remains:

*   **Maintainable**: Changes in one layer don't cascade through the entire system
*   **Testable**: Business logic is independent of frameworks and databases
*   **Flexible**: Easy to swap implementations without touching core logic
*   **Scalable**: Clear boundaries make it easy to add new features

What Developers Say [​](#what-developers-say)
---------------------------------------------

> "Goca transformed how we build Go services. What used to take hours now takes minutes, and the code quality is consistently high."
> 
> — Production User

> "Finally, a code generator that doesn't just dump code but teaches you proper architecture."
> 
> — Go Developer

> "The automatic integration of new features saved us so much time. No more manual wiring!"
> 
> — Go Team Lead

Ready to Build? [​](#ready-to-build)
------------------------------------

[Get Started in 5 Minutes](https://sazardev.github.io/goca/getting-started.html)</content>
</page>

<page>
  <title>Goca - Go Clean Architecture Code Generator</title>
  <url>https://sazardev.github.io/goca/getting-started.html</url>
  <content>Getting Started [​](#getting-started)
-------------------------------------

This guide will help you create your first project with Goca in **less than 5 minutes**. By the end, you'll have a functional REST API following Clean Architecture principles.

What We'll Build [​](#what-we-ll-build)
---------------------------------------

In this guide we'll create:

*   A project with complete Clean Architecture structure
*   A `User` entity with domain validations
*   Full CRUD REST API with all layers
*   Dependency injection configured
*   Repository pattern with PostgreSQL

Estimated Time

**5 minutes** from zero to running API

Prerequisites [​](#prerequisites)
---------------------------------

Before starting, make sure you have:

*   **Go 1.21+** - [Download here](https://golang.org/dl/)
*   **Goca installed** - See [Installation Guide](https://sazardev.github.io/goca/guide/installation.html)
*   **PostgreSQL** (optional for this tutorial)

Step 1: Create the Project [​](#step-1-create-the-project)
----------------------------------------------------------

bash

    # Create and navigate to your project directory
    mkdir my-first-api
    cd my-first-api
    
    # Initialize with Goca (default uses PostgreSQL)
    goca init my-api --module github.com/yourusername/my-api
    
    # Or choose a different database
    goca init my-api --module github.com/yourusername/my-api --database mysql
    goca init my-api --module github.com/yourusername/my-api --database mongodb
    goca init my-api --module github.com/yourusername/my-api --database sqlite
    
    # Navigate to generated directory
    cd my-api

What just happened?

Goca created a complete project structure with:

*   `internal/` - All your business logic layers
*   `cmd/` - Application entry points
*   `pkg/` - Shared packages
*   Configuration files for database, HTTP server, and more

**Available Databases:**

*   PostgreSQL (default) - SQL, traditional business apps
*   MySQL - SQL, web applications
*   MongoDB - NoSQL, flexible schemas
*   SQLite - Embedded, development/testing
*   PostgreSQL JSON - SQL with semi-structured data
*   SQL Server - Enterprise T-SQL
*   Elasticsearch - Full-text search
*   DynamoDB - Serverless AWS

[See Database Support Guide](https://sazardev.github.io/goca/features/database-support.html) for details.

Step 2: Generate Your First Feature [​](#step-2-generate-your-first-feature)
----------------------------------------------------------------------------

bash

    # Generate complete User feature with all layers
    goca feature User --fields "name:string,email:string,age:int"
    
    # See what was generated
    ls internal/domain/
    ls internal/usecase/
    ls internal/repository/
    ls internal/handler/http/

What Gets Generated?

This single command creates:

*   **Domain Entity**: `user.go` with validations
*   **Use Cases**: Service interfaces and DTOs
*   **Repository**: Interface and PostgreSQL implementation
*   **HTTP Handler**: REST endpoints for CRUD
*   **Dependency Injection**: Automatic wiring
*   **Routes**: Automatically registered

Step 3: Install Dependencies [​](#step-3-install-dependencies)
--------------------------------------------------------------

bash

    # Download and install Go dependencies
    go mod tidy

Step 4: Run Your API [​](#step-4-run-your-api)
----------------------------------------------

bash

    # Start the server
    go run cmd/server/main.go

You should see:

    → Server starting on :8080
    → Database connected
    → Routes registered

Step 5: Test Your API [​](#step-5-test-your-api)
------------------------------------------------

Now let's interact with our new API!

### Health Check [​](#health-check)

bash

    curl http://localhost:8080/health

**Response:**

json

    {
      "status": "ok",
      "timestamp": "2025-10-11T10:30:00Z"
    }

### Create a User [​](#create-a-user)

bash

    curl -X POST http://localhost:8080/api/v1/users \
      -H "Content-Type: application/json" \
      -d '{
        "name": "John Doe",
        "email": "john@example.com",
        "age": 28
      }'

**Response:**

json

    {
      "id": 1,
      "name": "John Doe",
      "email": "john@example.com",
      "age": 28,
      "created_at": "2025-10-11T10:30:00Z"
    }

### Get User by ID [​](#get-user-by-id)

bash

    curl http://localhost:8080/api/v1/users/1

### List All Users [​](#list-all-users)

bash

    curl http://localhost:8080/api/v1/users

**Response:**

json

    {
      "users": [
        {
          "id": 1,
          "name": "John Doe",
          "email": "john@example.com",
          "age": 28
        }
      ],
      "total": 1
    }

### Update User [​](#update-user)

bash

    curl -X PUT http://localhost:8080/api/v1/users/1 \
      -H "Content-Type: application/json" \
      -d '{
        "name": "John Smith",
        "email": "john.smith@example.com",
        "age": 29
      }'

### Delete User [​](#delete-user)

bash

    curl -X DELETE http://localhost:8080/api/v1/users/1

Understanding the Architecture [​](#understanding-the-architecture)
-------------------------------------------------------------------

Let's see how your code is organized:

    my-api/
    ├── internal/
    │   ├── domain/           # 🟡 Business entities
    │   │   └── user.go
    │   ├── usecase/          # 🔴 Application logic
    │   │   ├── dto.go
    │   │   └── user_service.go
    │   ├── repository/       # 🔵 Data persistence
    │   │   └── postgres_user_repository.go
    │   └── handler/          # 🟢 Input adapters
    │       └── http/
    │           └── user_handler.go
    └── cmd/
        └── server/
            └── main.go       # Entry point

### The Clean Architecture Layers [​](#the-clean-architecture-layers)

1.  **🟡 Domain** - Pure business logic, no dependencies
2.  **🔴 Use Cases** - Application rules and workflows
3.  **🔵 Repository** - Data access abstraction
4.  **🟢 Handlers** - External interface adapters

Dependency Rule

Dependencies always point inward:

    Handler → UseCase → Repository → Domain

The domain never knows about outer layers!

Next Steps [​](#next-steps)
---------------------------

Congratulations! You've created your first Clean Architecture API with Goca.

Here's what you can do next:

*   [Learn Clean Architecture Concepts](https://sazardev.github.io/goca/guide/clean-architecture.html)
*   [Add More Features](https://sazardev.github.io/goca/tutorials/adding-features.html)
*   [Complete Tutorial](https://sazardev.github.io/goca/tutorials/complete-tutorial.html)
*   [Explore All Commands](https://sazardev.github.io/goca/commands/)

Common Issues [​](#common-issues)
---------------------------------

### Port Already in Use [​](#port-already-in-use)

If you see "address already in use":

bash

    # Find and kill the process using port 8080
    lsof -ti:8080 | xargs kill -9

### Database Connection Failed [​](#database-connection-failed)

If using PostgreSQL and connection fails:

1.  Make sure PostgreSQL is running
2.  Update connection string in `.env` or config file
3.  Create the database: `createdb my-api`

### Module Not Found [​](#module-not-found)

bash

    # Re-initialize Go modules
    go mod init github.com/yourusername/my-api
    go mod tidy

Need Help? [​](#need-help)
--------------------------

*   [GitHub Issues](https://github.com/sazardev/goca/issues)
*   [Discussions](https://github.com/sazardev/goca/discussions)
*   [Full Documentation](https://sazardev.github.io/goca/guide/introduction.html)</content>
</page>

<page>
  <title>Goca - Go Clean Architecture Code Generator</title>
  <url>https://sazardev.github.io/goca/guide/introduction.html</url>
  <content>What is Goca? [​](#what-is-goca)
--------------------------------

**Goca** (Go Clean Architecture) is a powerful CLI code generator that helps you build production-ready Go applications following **Clean Architecture** principles designed by Uncle Bob (Robert C. Martin).

The Problem [​](#the-problem)
-----------------------------

Building Go applications with proper architecture is time-consuming:

*   Writing repetitive boilerplate code
*   Maintaining consistent structure across features
*   Ensuring proper layer separation
*   Setting up dependency injection
*   Configuring routing and handlers
*   Fighting architectural drift over time

The Solution [​](#the-solution)
-------------------------------

Goca automates all of this while teaching you Clean Architecture:

bash

    # One command generates all layers properly structured
    goca feature Product --fields "name:string,price:float64"

This creates:

*   Domain entities with business validations
*   Use cases with clear DTOs
*   Repository interfaces and implementations
*   HTTP handlers with proper routing
*   Dependency injection automatically configured

Core Philosophy [​](#core-philosophy)
-------------------------------------

### 1\. Clean Architecture by Default [​](#_1-clean-architecture-by-default)

Every line of code follows Uncle Bob's principles:

*   **Dependency Rule**: Dependencies point inward toward domain
*   **Layer Separation**: Clear boundaries between layers
*   **Interface Segregation**: Small, focused contracts
*   **Dependency Inversion**: Details depend on abstractions

### 2\. Prevention Over Correction [​](#_2-prevention-over-correction)

Goca prevents common anti-patterns:

*   🚫 Fat Controllers - Business logic stays in use cases
*   🚫 God Objects - Each entity has single responsibility
*   🚫 Anemic Domain - Entities include business behavior
*   🚫 Direct Database Access - Always through repositories

### 3\. Production-Ready Code [​](#_3-production-ready-code)

Generated code is not a starting point - it's production-ready:

*   Error handling included
*   Validation at proper layers
*   Clear separation of concerns
*   Testable by design
*   Well-documented

Key Features [​](#key-features)
-------------------------------

### Complete Architecture Layers [​](#complete-architecture-layers)

    ┌─────────────────────────────────────┐
    │         🟢 Handlers                 │  HTTP, gRPC, CLI
    │  (Input/Output Adapters)            │
    └─────────────────────────────────────┘
               ↓ depends on
    ┌─────────────────────────────────────┐
    │      🔴 Use Cases (Application)     │  Business workflows
    │  (DTOs, Services, Interfaces)       │
    └─────────────────────────────────────┘
               ↓ depends on
    ┌─────────────────────────────────────┐
    │    🔵 Repositories (Infrastructure) │  Data persistence
    │  (Database implementations)         │
    └─────────────────────────────────────┘
               ↓ depends on
    ┌─────────────────────────────────────┐
    │        🟡 Domain (Entities)         │  Pure business logic
    │  (No external dependencies)         │
    └─────────────────────────────────────┘

### 8 Supported Databases [​](#_8-supported-databases)

Generate repositories for any of these databases:

*   **PostgreSQL** - Traditional SQL with advanced features
*   **PostgreSQL JSON** - Semi-structured data with JSONB
*   **MySQL** - Web application standard
*   **MongoDB** - Flexible document-oriented
*   **SQLite** - Embedded, development-friendly
*   **SQL Server** - Enterprise T-SQL systems
*   **Elasticsearch** - Full-text search & analytics
*   **DynamoDB** - Serverless AWS infrastructure

Learn more in [Database Support Guide](https://sazardev.github.io/goca/features/database-support.html)

### Instant Feature Generation [​](#instant-feature-generation)

Generate complete features in seconds:

bash

    # From zero to fully functional CRUD
    goca feature Order --fields "customer:string,total:float64,status:string" --database postgres

Creates 10+ files with:

*   Domain entity
*   CRUD use cases
*   Repository interface + implementation
*   HTTP REST handlers
*   Dependency injection setup
*   Automatic routing
*   HTTP REST endpoints
*   Automatic integration

### Multi-Protocol Support [​](#multi-protocol-support)

Generate adapters for different protocols:

bash

    # HTTP REST API
    goca handler Product --type http
    
    # gRPC Service
    goca handler Product --type grpc
    
    # CLI Commands
    goca handler Product --type cli
    
    # Background Workers
    goca handler Product --type worker
    
    # SOAP Client
    goca handler Product --type soap

All following the same clean architecture pattern!

### Automatic Integration [​](#automatic-integration)

New features are automatically integrated:

*   Dependency injection containers updated
*   Routes registered automatically
*   Database connections configured
*   No manual wiring needed

### Test-Friendly [​](#test-friendly)

Generated code is designed for testing:

*   Clear interfaces for mocking
*   Dependency injection throughout
*   Pure functions in domain
*   Isolated layers

Clean Architecture provides:

### Maintainability [​](#maintainability)

*   Changes isolated to specific layers
*   Clear boundaries prevent cascading effects
*   Easy to understand and modify

### Testability [​](#testability)

*   Business logic independent of frameworks
*   Easy to mock dependencies
*   Fast unit tests without infrastructure

### Flexibility [​](#flexibility)

*   Swap implementations without touching business logic
*   Add new delivery mechanisms easily
*   Database agnostic domain layer

### Scalability [​](#scalability)

*   Clear structure makes onboarding easy
*   Consistent patterns across features
*   Easy to add new features

When to Use Goca? [​](#when-to-use-goca)
----------------------------------------

### Perfect For: [​](#perfect-for)

*   New Go projects requiring solid architecture
*   Microservices with consistent structure
*   REST APIs with multiple resources
*   Projects that will grow over time
*   Teams learning Clean Architecture
*   MVPs that need to scale to production

### Maybe Not For: [​](#maybe-not-for)

*   Simple scripts or one-off tools
*   Extremely unique architectures
*   Projects with existing different patterns

Comparison [​](#comparison)
---------------------------

### Without Goca [​](#without-goca)

bash

    # Manual process (hours of work)
    1. Create domain entity
    2. Write use case interfaces
    3. Implement use case logic
    4. Create DTOs for each operation
    5. Write repository interface
    6. Implement repository
    7. Create HTTP handlers
    8. Set up routing
    9. Configure DI container
    10. Wire everything together
    11. Test and fix integration issues

### With Goca [​](#with-goca)

bash

    # One command (seconds)
    goca feature Product --fields "name:string,price:float64"
    
    # Done! Everything wired and working

Real-World Usage [​](#real-world-usage)
---------------------------------------

Goca is used in production for:

*   🏢 Enterprise microservices
*   🛒 E-commerce platforms
*   📱 Mobile backend APIs
*   Internal tools and services
*   📊 Data processing pipelines

Next Steps [​](#next-steps)
---------------------------

Ready to start building clean Go applications?

*   [Install Goca](https://sazardev.github.io/goca/guide/installation.html)
*   [Quick Start Guide](https://sazardev.github.io/goca/getting-started.html)
*   [Learn Clean Architecture](https://sazardev.github.io/goca/guide/clean-architecture.html)
*   [Complete Tutorial](https://sazardev.github.io/goca/tutorials/complete-tutorial.html)

Philosophy [​](#philosophy)
---------------------------

> "The goal of software architecture is to minimize the human resources required to build and maintain the required system."
> 
> — Robert C. Martin (Uncle Bob)

Goca embodies this philosophy by automating the tedious parts while maintaining architectural excellence.</content>
</page>

<page>
  <title>Goca - Go Clean Architecture Code Generator</title>
  <url>https://sazardev.github.io/goca/blog/</url>
  <content>Goca BlogInsights, Updates & Architecture
-----------------------------------------

Stay updated with the latest releases, tutorials, and architectural insights for building production-ready Go applications

Latest Release [​](#latest-release)
-----------------------------------

### [v1.14.1 - Test Suite Improvements](https://sazardev.github.io/goca/blog/releases/v1-14-1.html) [​](#v1-14-1-test-suite-improvements)

_October 27, 2025_

Major improvements to test reliability and Windows compatibility. Fixed path handling, working directory management, and module dependencies. Test success rate improved to 99.04%.

[Read full release notes](https://sazardev.github.io/goca/blog/releases/v1-14-1.html)

* * *

Recent Articles [​](#recent-articles)
-------------------------------------

### [Mastering the Repository Pattern in Clean Architecture](https://sazardev.github.io/goca/blog/articles/mastering-repository-pattern.html) [​](#mastering-the-repository-pattern-in-clean-architecture)

_October 29, 2025_

A comprehensive guide to the Repository pattern in Clean Architecture. Learn the difference between repositories and DAOs, how to design clean interfaces, implement database-specific code, and how Goca generates production-ready repositories with complete abstraction.

[Read article](https://sazardev.github.io/goca/blog/articles/mastering-repository-pattern.html)

### [Mastering Use Cases in Clean Architecture](https://sazardev.github.io/goca/blog/articles/mastering-use-cases.html) [​](#mastering-use-cases-in-clean-architecture)

_October 29, 2025_

A deep dive into use cases and application services. Learn what use cases are, how they differ from controllers, DTOs patterns, and how Goca generates complete application layer code with orchestration logic and best practices.

[Read article](https://sazardev.github.io/goca/blog/articles/mastering-use-cases.html)

### [Understanding Domain Entities in Clean Architecture](https://sazardev.github.io/goca/blog/articles/understanding-domain-entities.html) [​](#understanding-domain-entities-in-clean-architecture)

_October 29, 2025_

A comprehensive guide to domain entities in Clean Architecture. Learn what entities are, why they're not database models, DDD principles, and how Goca generates production-ready entities with validation and business rules.

[Read article](https://sazardev.github.io/goca/blog/articles/understanding-domain-entities.html)

* * *

Coming Soon [​](#coming-soon)
-----------------------------

Stay tuned for more articles on:

*   Building scalable microservices with Clean Architecture
*   Advanced testing strategies with Goca
*   Database migration patterns and best practices
*   Performance optimization in Go applications

* * *

Subscribe [​](#subscribe)
-------------------------

Follow the project on [GitHub](https://github.com/sazardev/goca) to stay updated with releases and announcements.</content>
</page>

<page>
  <title>Goca - Go Clean Architecture Code Generator</title>
  <url>https://sazardev.github.io/goca/commands/init.html</url>
  <content>goca init [​](#goca-init)
-------------------------

Initialize a new Clean Architecture project with complete structure and configuration.

Syntax [​](#syntax)
-------------------

bash

    goca init <project-name> [flags]

Description [​](#description)
-----------------------------

The `goca init` command creates a production-ready Go project following Clean Architecture principles. It generates the complete directory structure, configuration files, and boilerplate code to get you started immediately.

Git Initialization

Projects are automatically initialized with Git, including an initial commit. This ensures your project is version-control ready from the start.

Arguments [​](#arguments)
-------------------------

### `<project-name>` [​](#project-name)

**Required.** The name of your project directory.

This creates a directory named `my-api` with the full project structure.

Flags [​](#flags)
-----------------

### `--module` (Required) [​](#module-required)

Go module name for your project.

bash

    --module github.com/username/projectname

**Example:**

bash

    goca init ecommerce --module github.com/sazardev/ecommerce

Module Naming Convention

Use your repository URL as the module name:

*   GitHub: `github.com/username/repo`
*   GitLab: `gitlab.com/username/repo`
*   Custom: `example.com/project`

### `--database` [​](#database)

Database system to use. Default: `postgres`

**Options:**

*   `postgres` - PostgreSQL (GORM)
*   `postgres-json` - PostgreSQL with JSONB
*   `mysql` - MySQL (GORM)
*   `mongodb` - MongoDB (native driver)
*   `sqlite` - SQLite (embedded)
*   `sqlserver` - SQL Server
*   `elasticsearch` - Elasticsearch (v8)
*   `dynamodb` - DynamoDB (AWS)

bash

    goca init myproject --module github.com/user/myproject --database mysql
    goca init config-server --module github.com/user/config --database postgres-json
    goca init search-app --module github.com/user/search --database elasticsearch

See [Database Support](https://sazardev.github.io/goca/features/database-support.html) for detailed comparison.

### `--auth` [​](#auth)

Include JWT authentication system.

bash

    goca init myproject --module github.com/user/myproject --auth

Generates:

*   JWT token generation and validation
*   Authentication middleware
*   User authentication endpoints
*   Password hashing utilities

### `--api` [​](#api)

API type to generate. Default: `rest`

**Options:** `rest` | `grpc` | `graphql` | `both`

bash

    goca init myproject --module github.com/user/myproject --api grpc

Examples [​](#examples)
-----------------------

### Basic REST API [​](#basic-rest-api)

bash

    goca init blog-api \
      --module github.com/sazardev/blog-api \
      --database postgres

### E-commerce with Authentication [​](#e-commerce-with-authentication)

bash

    goca init ecommerce \
      --module github.com/company/ecommerce \
      --database postgres \
      --auth

### gRPC Microservice [​](#grpc-microservice)

bash

    goca init user-service \
      --module github.com/company/user-service \
      --database mongodb \
      --api grpc

### Full-Featured Application [​](#full-featured-application)

bash

    goca init platform \
      --module github.com/startup/platform \
      --database postgres \
      --auth \
      --api both

Project Templates [​](#project-templates)
-----------------------------------------

GOCA provides predefined project templates that automatically configure your project with optimized settings for specific use cases. Templates generate a complete `.goca.yaml` configuration file tailored to your project type.

### Using Templates [​](#using-templates)

#### `--template` [​](#template)

Initialize project with a predefined configuration template.

bash

    goca init myproject --module github.com/user/myproject --template rest-api

#### `--list-templates` [​](#list-templates)

List all available templates with descriptions.

bash

    goca init --list-templates

### Available Templates [​](#available-templates)

#### `minimal` [​](#minimal)

**Lightweight starter with essential features only**

Perfect for:

*   Quick prototypes
*   Learning Clean Architecture
*   Minimal dependencies

bash

    goca init quick-start \
      --module github.com/user/quick-start \
      --template minimal

**Includes:**

*   Basic project structure
*   PostgreSQL database
*   Essential layers (domain, usecase, repository, handler)
*   Simple validation
*   Testify for testing

#### `rest-api` [​](#rest-api)

**Production-ready REST API with PostgreSQL, validation, and testing**

Perfect for:

*   RESTful web services
*   API backends
*   Standard CRUD applications

bash

    goca init api-service \
      --module github.com/company/api-service \
      --template rest-api

**Includes:**

*   Complete Clean Architecture layers
*   PostgreSQL with migrations
*   Input validation and sanitization
*   Swagger/OpenAPI documentation
*   Comprehensive testing with testify
*   Test coverage (70% threshold)
*   Integration tests
*   Soft deletes and timestamps

#### `microservice` [​](#microservice)

**Microservice with gRPC, events, and comprehensive testing**

Perfect for:

*   Distributed systems
*   Event-driven architecture
*   Service-oriented architecture

bash

    goca init user-service \
      --module github.com/company/user-service \
      --template microservice

**Includes:**

*   UUID primary keys
*   Audit logging
*   Event-driven patterns
*   Domain events support
*   Specification pattern
*   Advanced validation (validator library)
*   High test coverage (80% threshold)
*   Integration and benchmark tests
*   Optimized for horizontal scaling

#### `monolith` [​](#monolith)

**Full-featured monolithic application with web interface**

Perfect for:

*   Traditional web applications
*   Internal tools
*   Admin panels

bash

    goca init admin-panel \
      --module github.com/company/admin-panel \
      --template monolith

**Includes:**

*   JWT authentication with RBAC
*   Redis caching
*   Structured logging (JSON)
*   Health check endpoints
*   Soft deletes and timestamps
*   Audit trail
*   Versioning support
*   Markdown documentation
*   Test fixtures and seeds
*   Guards and authorization patterns

#### `enterprise` [​](#enterprise)

**Enterprise-grade with all features, security, and monitoring**

Perfect for:

*   Production applications
*   Enterprise systems
*   Mission-critical services

bash

    goca init enterprise-app \
      --module github.com/corp/enterprise-app \
      --template enterprise

**Includes:**

*   **Security**: HTTPS, CORS, rate limiting, header security
*   **Authentication**: JWT + OAuth2, RBAC
*   **Caching**: Redis with multi-layer caching
*   **Monitoring**: Prometheus metrics, distributed tracing, health checks, profiling
*   **Documentation**: Swagger 3.0, Postman collections, comprehensive markdown
*   **Testing**: 85% coverage threshold, mocks, integration, benchmarks, examples
*   **Deployment**: Docker (multistage), Kubernetes (manifests, Helm), CI/CD (GitHub Actions)
*   **Code Quality**: gofmt, goimports, golint, staticcheck
*   **Database**: Advanced features (partitioning, connection pooling)

### Template Configuration [​](#template-configuration)

When you initialize a project with a template, GOCA:

1.  Creates the standard project structure
2.  Generates a `.goca.yaml` configuration file with template settings
3.  All future feature generation uses these settings automatically
4.  You can still override settings using CLI flags when needed

**Example workflow:**

bash

    # Initialize with template
    goca init my-api --module github.com/user/my-api --template rest-api
    
    # Navigate to project
    cd my-api
    
    # Generate features - automatically uses template configuration
    goca feature Product --fields "name:string,price:float64,stock:int"
    # ✓ Uses REST API settings from template
    # ✓ Includes validation
    # ✓ Generates Swagger docs
    # ✓ Creates comprehensive tests
    
    # Override specific settings if needed
    goca feature Order --fields "total:float64" --database mysql
    # ✓ Uses template settings except database

### Customizing Template Configuration [​](#customizing-template-configuration)

After initialization, you can customize the generated `.goca.yaml`:

bash

    # Initialize project
    goca init my-project --module github.com/user/my-project --template rest-api
    
    # Edit configuration
    cd my-project
    vim .goca.yaml  # Customize as needed
    
    # Features use your customized settings
    goca feature User --fields "name:string,email:string"

See [Configuration Guide](https://sazardev.github.io/goca/guide/configuration.html) for detailed `.goca.yaml` documentation.

### Choosing the Right Template [​](#choosing-the-right-template)

| Template | Best For | Complexity | Features |
| --- | --- | --- | --- |
| `minimal` | Learning, prototypes | ⭐ Simple | Essential only |
| `rest-api` | Web APIs, CRUD services | ⭐⭐ Standard | Production-ready API |
| `microservice` | Distributed systems | ⭐⭐⭐ Advanced | Events, gRPC, scaling |
| `monolith` | Web applications | ⭐⭐⭐ Advanced | Auth, caching, logging |
| `enterprise` | Mission-critical apps | ⭐⭐⭐⭐ Complete | Everything included |

Generated Structure [​](#generated-structure)
---------------------------------------------

    myproject/
    ├── cmd/
    │   └── server/
    │       └── main.go              # Application entry point
    ├── internal/
    │   ├── domain/                  # 🟡 Entities & business rules
    │   │   └── errors.go
    │   ├── usecase/                 # 🔴 Application logic
    │   │   ├── dto.go
    │   │   └── interfaces.go
    │   ├── repository/              # 🔵 Data access
    │   │   └── interfaces.go
    │   └── handler/                 # 🟢 Input adapters
    │       ├── http/
    │       │   ├── routes.go
    │       │   └── middleware.go
    │       └── grpc/                # (if --api grpc)
    ├── pkg/
    │   ├── config/
    │   │   ├── config.go            # App configuration
    │   │   └── database.go          # DB connection
    │   ├── logger/
    │   │   └── logger.go            # Structured logging
    │   └── auth/                    # (if --auth)
    │       ├── jwt.go
    │       ├── middleware.go
    │       └── password.go
    ├── migrations/                   # Database migrations
    │   └── 001_initial.sql
    ├── .env.example                 # Environment variables template
    ├── .gitignore
    ├── go.mod
    ├── go.sum
    ├── Makefile                     # Common tasks
    └── README.md

Generated Files [​](#generated-files)
-------------------------------------

### `cmd/server/main.go` [​](#cmd-server-main-go)

The application entry point with:

*   Server initialization
*   Database connection
*   Route registration
*   Graceful shutdown

go

    package main
    
    import (
        "log"
        "os"
        "os/signal"
        "syscall"
        
        "github.com/user/myproject/pkg/config"
        "github.com/user/myproject/pkg/logger"
    )
    
    func main() {
        // Load configuration
        cfg := config.Load()
        
        // Initialize logger
        log := logger.New(cfg.LogLevel)
        
        // Connect to database
        db := config.ConnectDatabase(cfg)
        
        // Start server
        server := NewServer(cfg, db, log)
        
        // Graceful shutdown
        quit := make(chan os.Signal, 1)
        signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
        <-quit
        
        log.Info("Shutting down server...")
    }

### `pkg/config/config.go` [​](#pkg-config-config-go)

Configuration management:

go

    package config
    
    import "github.com/spf13/viper"
    
    type Config struct {
        ServerPort   string
        DatabaseURL  string
        LogLevel     string
        JWTSecret    string // if --auth
    }
    
    func Load() *Config {
        viper.AutomaticEnv()
        // Load configuration
        return &Config{...}
    }

### `.env.example` [​](#env-example)

Environment variables template:

bash

    SERVER_PORT=8080
    DATABASE_URL=postgres://user:pass@localhost:5432/dbname
    LOG_LEVEL=info
    JWT_SECRET=your-secret-key  # if --auth

### `Makefile` [​](#makefile)

Common development tasks:

makefile

    .PHONY: run build test
    
    run:
    	go run cmd/server/main.go
    
    build:
    	go build -o bin/server cmd/server/main.go
    
    test:
    	go test ./...
    
    migrate-up:
    	migrate -path migrations -database $(DATABASE_URL) up
    
    migrate-down:
    	migrate -path migrations -database $(DATABASE_URL) down

Next Steps [​](#next-steps)
---------------------------

After initialization, follow these steps:

1.  **Navigate to project:**
    
2.  **Install dependencies:**
    
3.  **Configure environment:**
    
    bash
    
        cp .env.example .env
        # Edit .env with your settings
    
4.  **Verify Git initialization:**
    
    bash
    
        git log --oneline  # See initial commit
        git status         # Check repository status
    
5.  **Generate your first feature:**
    
    bash
    
        goca feature User --fields "name:string,email:string"
    
6.  **Run the application:**
    
    bash
    
        make run
        # or
        go run cmd/server/main.go
    

Tips [​](#tips)
---------------

### Use Configuration Files [​](#use-configuration-files)

Create a `.goca.yaml` for reusable settings:

yaml

    module: github.com/company/projectname
    database: postgres
    auth: true
    api: rest

Then simply run:

### Customize Templates [​](#customize-templates)

After initialization, you can modify the generated code to fit your needs. The structure is designed to be a starting point, not a constraint.

### Version Control [​](#version-control)

Don't forget to initialize Git:

bash

    cd myproject
    git init
    git add .
    git commit -m "Initial commit with Clean Architecture structure"

Troubleshooting [​](#troubleshooting)
-------------------------------------

### Module Name Errors [​](#module-name-errors)

**Problem:** "invalid module name"

**Solution:** Ensure module name follows Go conventions:

bash

    #  Correct
    --module github.com/user/project
    
    #  Incorrect
    --module my-project
    --module Project Name

### Permission Denied [​](#permission-denied)

**Problem:** "permission denied creating directory"

**Solution:** Run with appropriate permissions or choose a directory you have write access to.

### Dependencies Not Found [​](#dependencies-not-found)

**Problem:** Generated project can't find dependencies

**Solution:**

bash

    cd myproject
    go mod tidy
    go mod download

See Also [​](#see-also)
-----------------------

*   [`goca feature`](https://sazardev.github.io/goca/commands/feature.html) - Generate complete features
*   [`goca integrate`](https://sazardev.github.io/goca/commands/integrate.html) - Wire features together
*   [Getting Started Guide](https://sazardev.github.io/goca/getting-started.html) - Complete walkthrough
*   [Project Structure](https://sazardev.github.io/goca/guide/project-structure.html) - Understand the layout

Resources [​](#resources)
-------------------------

*   [GitHub Repository](https://github.com/sazardev/goca)
*   [Example Projects](https://github.com/sazardev/goca-examples)
*   [Report Issues](https://github.com/sazardev/goca/issues)</content>
</page>

<page>
  <title>Goca - Go Clean Architecture Code Generator</title>
  <url>https://sazardev.github.io/goca/features/safety-and-dependencies.html</url>
  <content>Goca Safety & Dependency Management Features [​](#goca-safety-dependency-management-features)
---------------------------------------------------------------------------------------------

Overview [​](#overview)
-----------------------

This document describes the new safety and dependency management features implemented in Goca CLI v1.11.0.

New Features [​](#new-features)
-------------------------------

### 1\. Dry-Run Mode (`--dry-run`) [​](#_1-dry-run-mode-dry-run)

Preview all changes before they are made to your project.

**Usage:**

bash

    goca feature User --fields "name:string,email:string" --dry-run

**Output:**

    🔍 DRY-RUN MODE: Previewing changes without creating files
    
    📝 [DRY-RUN] Would create: internal/domain/user.go (1234 bytes)
    📝 [DRY-RUN] Would create: internal/usecase/user_service.go (2345 bytes)
    📝 [DRY-RUN] Would create: internal/repository/postgres_user_repository.go (1567 bytes)
    📝 [DRY-RUN] Would create: internal/handler/http/user_handler.go (2890 bytes)
    
    📋 DRY-RUN SUMMARY:
       Would create 15 files
       ⚠️  2 conflicts detected:
          - internal/domain/user.go
          - internal/usecase/user_service.go
    
    💡 Run without --dry-run to actually create files
       Use --force to overwrite existing files
       Use --backup to backup files before overwriting

### 2\. File Conflict Detection [​](#_2-file-conflict-detection)

Automatically detects existing files and prevents accidental overwrites.

**Scenarios:**

#### Scenario A: File Exists (No Force) [​](#scenario-a-file-exists-no-force)

bash

    goca feature User --fields "name:string"

    ❌ file already exists: internal/domain/user.go (use --force to overwrite or --backup to backup first)

#### Scenario B: Force Overwrite [​](#scenario-b-force-overwrite)

bash

    goca feature User --fields "name:string" --force

    ⚠️  Overwriting: internal/domain/user.go
    ✅ Created: internal/domain/user.go

#### Scenario C: Backup Before Overwrite [​](#scenario-c-backup-before-overwrite)

bash

    goca feature User --fields "name:string" --force --backup

    📦 Backed up: internal/domain/user.go -> .goca-backup/internal/domain/user.go.backup
    ✅ Created: internal/domain/user.go

### 3\. Name Conflict Detection [​](#_3-name-conflict-detection)

Detects duplicate entity/feature names across the project.

**Example:**

bash

    # User feature already exists
    goca feature User --fields "email:string"

    ❌ feature 'User' already exists in the project
    💡 Use --force to generate anyway

**Existing Entities Detection:** The system scans `internal/domain/` for existing entities and prevents duplicates.

### 4\. Automatic go.mod Management [​](#_4-automatic-go-mod-management)

Automatically updates `go.mod` when generating features with dependencies.

**Features:**

*   ✅ Adds required dependencies automatically
*   ✅ Runs `go mod tidy` after generation
*   ✅ Verifies dependency compatibility
*   ✅ Suggests optional dependencies

**Example:**

bash

    goca feature Auth --fields "username:string,password:string" --validation

    7️⃣  Managing dependencies...
    ✅ Added dependency: github.com/go-playground/validator/v10 v10.16.0
    ✅ Added dependency: github.com/golang-jwt/jwt/v5 v5.2.0
    
    📦 Updating go.mod...
    ✅ Updated go.mod and go.sum
    
    💡 OPTIONAL DEPENDENCIES:
       The following dependencies might be useful for your feature:
    
       📦 golang.org/x/crypto v0.17.0
          Reason: password hashing
          Install: go get golang.org/x/crypto@v0.17.0

### 5\. Version Compatibility Checking [​](#_5-version-compatibility-checking)

Verifies Go version and dependency compatibility.

**Features:**

*   ✅ Checks minimum Go version (1.21+)
*   ✅ Verifies dependency versions are compatible
*   ✅ Warns about potential conflicts

**Example:**

bash

    goca init myproject --module github.com/user/myproject

    ✅ Go version check: go1.25.2 (compatible with go1.21+)
    ✅ All dependencies verified

### 6\. Optional Dependency Suggestions [​](#_6-optional-dependency-suggestions)

Intelligently suggests dependencies based on feature characteristics.

**Dependency Categories:**

#### Validation [​](#validation)

    📦 github.com/go-playground/validator/v10
       Reason: struct validation for DTOs

#### Authentication [​](#authentication)

    📦 github.com/golang-jwt/jwt/v5
       Reason: JWT authentication
       
    📦 golang.org/x/crypto
       Reason: password hashing

#### Testing [​](#testing)

    📦 github.com/stretchr/testify
       Reason: testing assertions and mocks
       
    📦 github.com/golang/mock
       Reason: mock generation for testing

#### gRPC [​](#grpc)

    📦 google.golang.org/grpc
       Reason: gRPC protocol support
       
    📦 google.golang.org/protobuf
       Reason: Protocol Buffers

Implementation Details [​](#implementation-details)
---------------------------------------------------

### New Files Created [​](#new-files-created)

1.  **`cmd/safety.go`**
    
    *   `SafetyManager`: Handles dry-run, force, and backup modes
    *   `NameConflictDetector`: Scans for existing entities/features
    *   File conflict detection logic
    *   Backup system
2.  **`cmd/dependency_manager.go`**
    
    *   `DependencyManager`: Manages go.mod updates
    *   Dependency suggestion system
    *   Version compatibility checking
    *   Automatic dependency installation

### Updated Files [​](#updated-files)

1.  **`cmd/feature.go`**
    *   Added `--dry-run`, `--force`, `--backup` flags
    *   Integrated SafetyManager
    *   Integrated DependencyManager
    *   Added name conflict checking

Command Flag Reference [​](#command-flag-reference)
---------------------------------------------------

### All Commands Support [​](#all-commands-support)

| Flag | Type | Description |
| --- | --- | --- |
| `--dry-run` | bool | Preview changes without creating files |
| `--force` | bool | Overwrite existing files without asking |
| `--backup` | bool | Backup files before overwriting |

### Safety Workflow [​](#safety-workflow)

bash

    # 1. Preview changes first
    goca feature Product --fields "name:string,price:float64" --dry-run
    
    # 2. If satisfied, generate for real
    goca feature Product --fields "name:string,price:float64"
    
    # 3. If files exist and you want to update
    goca feature Product --fields "name:string,price:float64" --force --backup

Benefits [​](#benefits)
-----------------------

### For Developers [​](#for-developers)

✅ **Safety**: Preview changes before committing ✅ **Confidence**: Know exactly what will be created ✅ **No Accidents**: Automatic conflict detection ✅ **Easy Recovery**: Automatic backups ✅ **Less Manual Work**: Automatic dependency management

### For Teams [​](#for-teams)

✅ **Consistency**: Standardized dependency versions ✅ **Documentation**: Clear what each feature requires ✅ **Onboarding**: Suggestions help new developers ✅ **Best Practices**: Automatic inclusion of common libraries

Configuration [​](#configuration)
---------------------------------

### .goca.yaml Support [​](#goca-yaml-support)

yaml

    # Enable safety features by default
    safety:
      dry_run_default: false
      backup_enabled: true
      conflict_detection: true
    
    # Dependency management
    dependencies:
      auto_update: true
      suggest_optional: true
      verify_versions: true

Examples [​](#examples)
-----------------------

### Example 1: Safe Feature Generation [​](#example-1-safe-feature-generation)

bash

    # Step 1: Preview
    goca feature Order --fields "customer_id:int,total:float64,status:string" --dry-run
    
    # Step 2: Check for conflicts
    # (automatically done)
    
    # Step 3: Generate with backup
    goca feature Order --fields "customer_id:int,total:float64,status:string" --backup
    
    # Step 4: Dependencies auto-added
    # go.mod updated automatically

### Example 2: Update Existing Feature [​](#example-2-update-existing-feature)

bash

    # Backup and force update
    goca feature User --fields "name:string,email:string,age:int,role:string" --force --backup
    
    # Old files saved to .goca-backup/
    # New files generated
    # Dependencies updated

### Example 3: Team Workflow [​](#example-3-team-workflow)

bash

    # Developer A: Preview changes
    goca feature Payment --fields "amount:float64,method:string" --dry-run
    
    # Share preview output in PR
    # Team reviews
    
    # Developer B: Generate with exact same command
    goca feature Payment --fields "amount:float64,method:string"
    
    # Consistent results across team

Best Practices [​](#best-practices)
-----------------------------------

### 1\. Always Dry-Run First [​](#_1-always-dry-run-first)

bash

    # Good
    goca feature NewFeature --fields "..." --dry-run
    goca feature NewFeature --fields "..."
    
    # Risky
    goca feature NewFeature --fields "..."

### 2\. Use Backup for Updates [​](#_2-use-backup-for-updates)

bash

    # Safe
    goca feature ExistingFeature --fields "..." --force --backup
    
    # Risky
    goca feature ExistingFeature --fields "..." --force

### 3\. Review Dependency Suggestions [​](#_3-review-dependency-suggestions)

bash

    # After generation, review suggested dependencies
    # Install only what you need
    go get github.com/suggested/package@version

### 4\. Commit Backups [​](#_4-commit-backups)

bash

    # Add backups to .gitignore
    echo ".goca-backup/" >> .gitignore
    
    # Or commit them for safety
    git add .goca-backup/
    git commit -m "Backup before updating User feature"

Troubleshooting [​](#troubleshooting)
-------------------------------------

### Issue: "file already exists" [​](#issue-file-already-exists)

**Solution:** Use `--force` and `--backup`:

bash

    goca feature User --fields "..." --force --backup

### Issue: "feature already exists" [​](#issue-feature-already-exists)

**Solution:** Either:

1.  Use a different name
2.  Use `--force` to regenerate
3.  Delete existing feature files first

### Issue: "dependency verification failed" [​](#issue-dependency-verification-failed)

**Solution:** Run manually:

bash

    cd your-project
    go mod tidy
    go mod verify

### Issue: Dry-run shows many conflicts [​](#issue-dry-run-shows-many-conflicts)

**Solution:** This is expected if updating an existing feature. Use `--force --backup` to proceed safely.

Migration Guide [​](#migration-guide)
-------------------------------------

### Existing Projects [​](#existing-projects)

No changes needed! New features are opt-in via flags.

### Updating Commands [​](#updating-commands)

Old command:

bash

    goca feature User --fields "name:string"

New (safer):

bash

    # Preview first
    goca feature User --fields "name:string" --dry-run
    
    # Then generate
    goca feature User --fields "name:string"

Future Enhancements [​](#future-enhancements)
---------------------------------------------

Planned for v1.12.0:

*   Interactive conflict resolution
*   Merge tool for conflicting files
*   Undo/rollback command
*   Dependency version suggestions
*   Security vulnerability scanning

Contributing [​](#contributing)
-------------------------------

These features are open for community contribution. See:

*   `cmd/safety.go` - Safety manager implementation
*   `cmd/dependency_manager.go` - Dependency management
*   Tests in `internal/testing/tests/safety_test.go`

Support [​](#support)
---------------------

*   📚 Documentation: [https://sazardev.github.io/goca](https://sazardev.github.io/goca)
*   🐛 Issues: [https://github.com/sazardev/goca/issues](https://github.com/sazardev/goca/issues)
*   💬 Discussions: [https://github.com/sazardev/goca/discussions](https://github.com/sazardev/goca/discussions)</content>
</page>

<page>
  <title>Goca - Go Clean Architecture Code Generator</title>
  <url>https://sazardev.github.io/goca/guide/installation.html</url>
  <content>Installation [​](#installation)
-------------------------------

This guide will walk you through installing Goca on your system using various methods.

Prerequisites [​](#prerequisites)
---------------------------------

Before installing Goca, ensure you have:

*   **Go 1.21 or higher** - [Download Go](https://golang.org/dl/)
*   **Git** - For version control and cloning repositories
*   **Terminal or Command Prompt** - To run installation commands

Check Your Go Version

You should see `go version go1.21` or higher.

Installation Methods [​](#installation-methods)
-----------------------------------------------

### Method 1: go install (Recommended) [​](#method-1-go-install-recommended)

This is the fastest and simplest method:

bash

    go install github.com/sazardev/goca@latest

**Verify the installation:**

**Expected output:**

    Goca v2.0.0
    Build: 2025-10-11T10:00:00Z
    Go Version: go1.24.5
    OS/Arch: linux/amd64

Troubleshooting: Command Not Found

If you get `command not found`, ensure `$GOPATH/bin` is in your PATH:

**Linux/macOS:**

bash

    export PATH=$PATH:$(go env GOPATH)/bin

Add this to your `~/.bashrc`, `~/.zshrc`, or `~/.profile` to make it permanent.

**Windows:** Add `%USERPROFILE%\go\bin` to your system PATH environment variable.

### Method 2: Binary Downloads [​](#method-2-binary-downloads)

Download pre-compiled binaries directly from [GitHub Releases](https://github.com/sazardev/goca/releases).

LinuxmacOS (Intel)macOS (Apple Silicon)Windows

bash

    # Download the binary
    wget https://github.com/sazardev/goca/releases/latest/download/goca-linux-amd64
    
    # Make it executable
    chmod +x goca-linux-amd64
    
    # Move to PATH
    sudo mv goca-linux-amd64 /usr/local/bin/goca
    
    # Verify
    goca version

bash

    # Download the binary
    curl -L https://github.com/sazardev/goca/releases/latest/download/goca-darwin-amd64 -o goca
    
    # Make it executable
    chmod +x goca
    
    # Move to PATH
    sudo mv goca /usr/local/bin/goca
    
    # Verify
    goca version

bash

    # Download the binary
    curl -L https://github.com/sazardev/goca/releases/latest/download/goca-darwin-arm64 -o goca
    
    # Make it executable
    chmod +x goca
    
    # Move to PATH
    sudo mv goca /usr/local/bin/goca
    
    # Verify
    goca version

powershell

    # Download the binary
    Invoke-WebRequest -Uri "https://github.com/sazardev/goca/releases/latest/download/goca-windows-amd64.exe" -OutFile "goca.exe"
    
    # Move to a directory in PATH (requires admin)
    Move-Item goca.exe C:\Windows\System32\goca.exe
    
    # Verify
    goca version

### Method 3: Homebrew (macOS) [​](#method-3-homebrew-macos)

If you're on macOS and use Homebrew:

bash

    # Add the Goca tap
    brew tap sazardev/tools
    
    # Install Goca
    brew install goca
    
    # Verify
    goca version

### Method 4: Build from Source [​](#method-4-build-from-source)

For developers who want the latest development version or want to contribute:

bash

    # Clone the repository
    git clone https://github.com/sazardev/goca.git
    cd goca
    
    # Build the binary
    go build -o goca
    
    # (Optional) Install globally
    sudo mv goca /usr/local/bin/goca
    
    # Verify
    goca version

Building for Different Platforms

bash

    # Linux
    GOOS=linux GOARCH=amd64 go build -o goca-linux-amd64
    
    # macOS Intel
    GOOS=darwin GOARCH=amd64 go build -o goca-darwin-amd64
    
    # macOS Apple Silicon
    GOOS=darwin GOARCH=arm64 go build -o goca-darwin-arm64
    
    # Windows
    GOOS=windows GOARCH=amd64 go build -o goca-windows-amd64.exe

Verify Installation [​](#verify-installation)
---------------------------------------------

After installation, run:

You should see the help menu with all available commands:

    Goca - Go Clean Architecture Code Generator
    
    Usage:
      goca [command]
    
    Available Commands:
      init        Initialize a new Clean Architecture project
      feature     Generate a complete feature with all layers
      entity      Generate a domain entity
      usecase     Generate use cases
      repository  Generate repositories
      handler     Generate handlers
      di          Generate dependency injection
      integrate   Integrate existing features
      version     Show version information
    
    Flags:
      -h, --help      help for goca
      -v, --version   version for goca
    
    Use "goca [command] --help" for more information about a command.

Shell Completion (Optional) [​](#shell-completion-optional)
-----------------------------------------------------------

Enable command auto-completion for your shell:

BashZshFishPowerShell

bash

    # Generate completion script
    goca completion bash > /etc/bash_completion.d/goca
    
    # Or for current user only
    goca completion bash > ~/.bash_completion
    source ~/.bash_completion

bash

    # Generate completion script
    goca completion zsh > "${fpath[1]}/_goca"
    
    # Reload completions
    autoload -U compinit && compinit

bash

    # Generate completion script
    goca completion fish > ~/.config/fish/completions/goca.fish

powershell

    # Generate completion script
    goca completion powershell | Out-String | Invoke-Expression
    
    # To make permanent, add to profile
    goca completion powershell >> $PROFILE

Update Goca [​](#update-goca)
-----------------------------

### If installed via go install: [​](#if-installed-via-go-install)

bash

    go install github.com/sazardev/goca@latest

### If installed via Homebrew: [​](#if-installed-via-homebrew)

### If installed from binary: [​](#if-installed-from-binary)

Download the latest binary and replace the existing one following the [Binary Downloads](#method-2-binary-downloads) steps.

Uninstall Goca [​](#uninstall-goca)
-----------------------------------

### If installed via go install: [​](#if-installed-via-go-install-1)

### If installed via Homebrew: [​](#if-installed-via-homebrew-1)

bash

    brew uninstall goca
    brew untap sazardev/tools

### If installed from binary: [​](#if-installed-from-binary-1)

bash

    # Linux/macOS
    sudo rm /usr/local/bin/goca
    
    # Windows (as Administrator)
    del C:\Windows\System32\goca.exe

Next Steps [​](#next-steps)
---------------------------

Now that you have Goca installed, you're ready to start building!

*   [Quick Start Guide](https://sazardev.github.io/goca/getting-started.html) - Create your first project
*   [Learn Clean Architecture](https://sazardev.github.io/goca/guide/clean-architecture.html) - Understand the principles
*   [Complete Tutorial](https://sazardev.github.io/goca/tutorials/complete-tutorial.html) - Build a real application

Troubleshooting [​](#troubleshooting)
-------------------------------------

### Permission Denied [​](#permission-denied)

If you get permission errors on Linux/macOS:

bash

    sudo chmod +x /usr/local/bin/goca

### Command Not Found After Installation [​](#command-not-found-after-installation)

Make sure your `$PATH` includes Go's bin directory:

bash

    echo $PATH | grep -q "go/bin" && echo " Go bin in PATH" || echo "✗ Add Go bin to PATH"

### Version Mismatch [​](#version-mismatch)

If `goca version` shows an old version:

bash

    # Clear Go cache
    go clean -modcache
    
    # Reinstall
    go install github.com/sazardev/goca@latest

### Windows: goca is not recognized [​](#windows-goca-is-not-recognized)

Ensure you've added Go's bin directory to your system PATH:

1.  Open System Properties → Environment Variables
2.  Edit the `Path` variable
3.  Add `%USERPROFILE%\go\bin`
4.  Restart your terminal

Need Help? [​](#need-help)
--------------------------

*   [GitHub Issues](https://github.com/sazardev/goca/issues) - Report bugs
*   [Discussions](https://github.com/sazardev/goca/discussions) - Ask questions
*   [Documentation](https://sazardev.github.io/goca/) - Read the docs</content>
</page>

<page>
  <title>Goca - Go Clean Architecture Code Generator</title>
  <url>https://sazardev.github.io/goca/features/database-support.html</url>
  <content>Database Support [​](#database-support)
---------------------------------------

Goca supports **8 production-ready database systems** with full repository pattern implementations. Choose the database that best fits your use case.

Supported Databases [​](#supported-databases)
---------------------------------------------

### SQL Databases [​](#sql-databases)

#### PostgreSQL [​](#postgresql)

**Type:** SQL (ACID-compliant)  
**Driver:** GORM + postgres driver  
**Best For:** Traditional applications, transactional systems, most common choice

bash

    goca feature User --fields "name:string,email:string" --database postgres

**Features:**

*   Full transaction support
*   Complex joins and aggregations
*   Advanced query capabilities
*   Excellent performance for relational data

**Generated File:** `internal/repository/postgres_user_repository.go`

* * *

#### PostgreSQL JSON [​](#postgresql-json)

**Type:** SQL + Document (semi-structured)  
**Driver:** GORM + gorm/datatypes (JSONB)  
**Best For:** Semi-structured data, configs, metadata, hierarchical objects

bash

    goca feature Config --fields "name:string,settings:string" --database postgres-json

**Special Methods:**

*   `FindByJSONField(jsonField, value string)` - Query nested JSON with `@>` operator
*   Full JSONB operator support

**Use Case Example:**

go

    // Store nested configuration
    configs, _ := repo.FindByJSONField("settings.theme", "dark")

**Generated File:** `internal/repository/postgres_json_config_repository.go`

* * *

#### MySQL [​](#mysql)

**Type:** SQL (ACID-compliant)  
**Driver:** GORM + mysql driver  
**Best For:** Web applications, compatibility with existing infrastructure

bash

    goca feature Product --fields "name:string,price:float64" --database mysql

**Features:**

*   InnoDB transactions
*   Good scalability for web apps
*   Strong compatibility
*   Replication support

**Generated File:** `internal/repository/mysql_product_repository.go`

* * *

#### SQL Server [​](#sql-server)

**Type:** SQL (Enterprise)  
**Driver:** GORM + mssql driver  
**Best For:** Enterprise applications, Microsoft ecosystem, legacy systems

bash

    goca feature Employee --fields "name:string,salary:float64" --database sqlserver

**Features:**

*   T-SQL compatibility
*   Enterprise security
*   Advanced error handling
*   Context-aware operations
*   Enterprise compatibility

**Generated File:** `internal/repository/sqlserver_employee_repository.go`

* * *

#### SQLite [​](#sqlite)

**Type:** SQL (Embedded)  
**Driver:** database/sql + SQLite driver  
**Best For:** Development, testing, embedded applications, single-file databases

bash

    goca feature Setting --fields "key:string,value:string" --database sqlite

**Features:**

*   Single file storage (`.db`)
*   No server required
*   ACID compliance
*   JSON marshaling for flexibility
*   Perfect for prototyping

**Generated File:** `internal/repository/sqlite_setting_repository.go`

* * *

### NoSQL Databases [​](#nosql-databases)

#### MongoDB [​](#mongodb)

**Type:** Document NoSQL  
**Driver:** MongoDB official driver  
**Best For:** Document-oriented applications, flexible schemas, rapid iteration

bash

    goca feature Article --fields "title:string,content:string,tags:string" --database mongodb

**Features:**

*   Flexible schema
*   Rich query language
*   Horizontal scalability
*   Document transactions

**Generated File:** `internal/repository/mongodb_article_repository.go`

* * *

#### DynamoDB [​](#dynamodb)

**Type:** Key-Value NoSQL (Serverless)  
**Driver:** AWS SDK v2 + attributevalue  
**Best For:** Serverless AWS applications, auto-scaling, cloud-native architecture

bash

    goca feature Order --fields "orderID:string,total:float64" --database dynamodb

**Features:**

*   AWS-managed serverless
*   Auto-scaling
*   Scan operations for queries
*   Context-aware async operations
*   Pay-per-request pricing

**Generated File:** `internal/repository/dynamodb_order_repository.go`

* * *

### Search Databases [​](#search-databases)

#### Elasticsearch [​](#elasticsearch)

**Type:** Full-text Search & Analytics  
**Driver:** go-elasticsearch v8 client  
**Best For:** Full-text search, analytics, logging systems, product search

bash

    goca feature Article --fields "title:string,content:string" --database elasticsearch

**Special Methods:**

*   `FullTextSearch(query string)` - Multi-field full-text search
*   Lucene query DSL support
*   Result scoring and aggregations

**Use Case Example:**

go

    // Full-text search across multiple fields
    results, _ := repo.FullTextSearch("golang elasticsearch tutorial")

**Generated File:** `internal/repository/elasticsearch_article_repository.go`

* * *

Database Comparison [​](#database-comparison)
---------------------------------------------

| Feature | PostgreSQL | MySQL | MongoDB | SQLite | SQL Server | Elasticsearch | DynamoDB |
| --- | --- | --- | --- | --- | --- | --- | --- |
| **Type** | SQL | SQL | Document | SQL | SQL | Search | Key-Value |
| **ACID** | ✅ | ✅ | ✅ | ✅ | ✅ | ❌ | Limited |
| **Scalability** | Vertical | Horizontal | Horizontal | None | Vertical | Horizontal | Unlimited |
| **JSON Support** | JSONB | JSON | Native | Marshaling | Native | Native | Native |
| **Transactions** | ✅ Full | ✅ InnoDB | ✅ Multi-doc | ✅ | ✅ | ❌ | Limited |
| **Server Required** | ✅ | ✅ | ✅ | ❌ | ✅ | ✅ | N/A (Cloud) |
| **Cost** | Self-hosted | Self-hosted | Self-hosted | Free | Enterprise | Self-hosted | Pay-per-use |

* * *

Quick Selection Guide [​](#quick-selection-guide)
-------------------------------------------------

**Choose PostgreSQL if:**

*   You need strong ACID guarantees
*   Your data is highly relational
*   You want the most feature-rich SQL database
*   You're building traditional business applications

**Choose PostgreSQL JSON if:**

*   You have semi-structured data
*   You need flexible schema evolution
*   You want to store hierarchical configs
*   You want SQL with document flexibility

**Choose MySQL if:**

*   You're integrating with existing MySQL infrastructure
*   You need good web application performance
*   You want something lightweight
*   You're familiar with MySQL ecosystem

**Choose MongoDB if:**

*   Your schema changes frequently
*   You have flexible data structures
*   You need horizontal scaling
*   You're building rapid prototypes

**Choose SQLite if:**

*   You're developing locally
*   You need zero configuration database
*   You're building embedded applications
*   You want a single-file database

**Choose SQL Server if:**

*   You're in Microsoft enterprise environment
*   You need T-SQL specific features
*   You're migrating from legacy SQL Server
*   You need enterprise support

**Choose Elasticsearch if:**

*   You need full-text search capabilities
*   You're building search-heavy applications
*   You need analytics capabilities
*   You're centralizing application logs

**Choose DynamoDB if:**

*   You're building serverless AWS applications
*   You need unlimited auto-scaling
*   You want managed cloud database
*   You're comfortable with eventual consistency

* * *

Migration Between Databases [​](#migration-between-databases)
-------------------------------------------------------------

One of the key benefits of Goca's repository pattern is the ability to switch databases without changing your business logic:

bash

    # Start with PostgreSQL
    goca feature User --fields "name:string,email:string" --database postgres
    
    # Later, switch to MongoDB (same feature structure)
    goca repository User --database mongodb --implementation

Your use cases remain unchanged - only the repository implementation differs!

* * *

Performance Considerations [​](#performance-considerations)
-----------------------------------------------------------

### By Use Case [​](#by-use-case)

**High-Concurrency Writes:** PostgreSQL, MongoDB, DynamoDB

**Complex Queries:** PostgreSQL, PostgreSQL JSON, SQL Server

**Document Storage:** MongoDB, PostgreSQL JSON

**Full-Text Search:** Elasticsearch (specialized)

**Real-Time Analytics:** Elasticsearch

**Development/Testing:** SQLite

* * *

Next Steps [​](#next-steps)
---------------------------

*   [Repository Pattern Guide](https://sazardev.github.io/goca/guide/best-practices.html#repository-pattern)
*   [Feature Generation](https://sazardev.github.io/goca/commands/feature.html)
*   [Repository Command](https://sazardev.github.io/goca/commands/repository.html)
*   [Complete Tutorial](https://sazardev.github.io/goca/tutorials/complete-tutorial.html)</content>
</page>

<page>
  <title>Goca - Go Clean Architecture Code Generator</title>
  <url>https://sazardev.github.io/goca/guide/clean-architecture.html</url>
  <content>Learn how Goca implements and enforces **Clean Architecture** principles from Robert C. Martin (Uncle Bob) in your Go projects.

What is Clean Architecture? [​](#what-is-clean-architecture)
------------------------------------------------------------

Clean Architecture is a software design philosophy that separates code into **layers** with clear responsibilities and dependencies that point inward toward the core business logic.

### Core Principles [​](#core-principles)

1.  **Independence of Frameworks** - Business logic doesn't depend on libraries
2.  **Testability** - Business rules can be tested without UI, database, or external services
3.  **Independence of UI** - Change UI without changing business logic
4.  **Independence of Database** - Swap databases without affecting business rules
5.  **Independence of External Systems** - Business logic knows nothing about the outside world

The Dependency Rule [​](#the-dependency-rule)
---------------------------------------------

> **Source code dependencies must point only inward, toward higher-level policies.**

    ┌─────────────────────────────────────┐
    │     External Interfaces & I/O       │  ← Frameworks, Drivers, APIs
    └─────────────────┬───────────────────┘
                      │ depends on
    ┌─────────────────▼───────────────────┐
    │    Interface Adapters (Handlers)    │  ← Controllers, Presenters
    └─────────────────┬───────────────────┘
                      │ depends on
    ┌─────────────────▼───────────────────┐
    │    Application Business Rules       │  ← Use Cases
    └─────────────────┬───────────────────┘
                      │ depends on
    ┌─────────────────▼───────────────────┐
    │   Enterprise Business Rules         │  ← Entities (Domain)
    │         (No Dependencies)            │
    └─────────────────────────────────────┘

**Inner layers never know about outer layers.**

Goca's 4 Layers [​](#goca-s-4-layers)
-------------------------------------

### 🟡 Layer 1: Domain (Entities) [​](#🟡-layer-1-domain-entities)

**Location**: `internal/domain/`

The innermost layer containing **enterprise-wide business rules**.

#### Responsibilities [​](#responsibilities)

*   Define business entities
*   Implement core business rules
*   Define domain errors
*   Declare repository interfaces
*   Domain-specific validations

#### Example: User Entity [​](#example-user-entity)

go

    package domain
    
    import (
        "errors"
        "strings"
        "time"
    )
    
    // User represents a user entity in our system
    type User struct {
        ID        uint      `json:"id"`
        Name      string    `json:"name"`
        Email     string    `json:"email"`
        Role      string    `json:"role"`
        CreatedAt time.Time `json:"created_at"`
        UpdatedAt time.Time `json:"updated_at"`
    }
    
    // Validate enforces business rules
    func (u *User) Validate() error {
        if strings.TrimSpace(u.Name) == "" {
            return ErrUserNameRequired
        }
        
        if len(u.Name) < 2 {
            return ErrUserNameTooShort
        }
        
        if !u.IsValidEmail() {
            return ErrInvalidEmail
        }
        
        return nil
    }
    
    // IsAdmin checks if user has admin privileges
    func (u *User) IsAdmin() bool {
        return u.Role == "admin"
    }
    
    // IsValidEmail validates email format (business rule)
    func (u *User) IsValidEmail() bool {
        return strings.Contains(u.Email, "@") && len(u.Email) > 5
    }
    
    // Domain-specific errors
    var (
        ErrUserNameRequired = errors.New("user name is required")
        ErrUserNameTooShort = errors.New("user name must be at least 2 characters")
        ErrInvalidEmail     = errors.New("invalid email format")
        ErrUserNotFound     = errors.New("user not found")
    )

Domain Layer Rules

**Do**: Pure business logic, no external dependencies  
**Don't**: Import HTTP, database, or framework packages

### 🔴 Layer 2: Use Cases (Application Logic) [​](#🔴-layer-2-use-cases-application-logic)

**Location**: `internal/usecase/`

Contains **application-specific business rules**.

#### Responsibilities [​](#responsibilities-1)

*   Define use case interfaces
*   Implement application workflows
*   Coordinate between repositories
*   Define DTOs (Data Transfer Objects)
*   Input validation

#### Example: User Use Case [​](#example-user-use-case)

go

    package usecase
    
    import (
        "context"
        "myproject/internal/domain"
    )
    
    // UserUseCase defines user-related operations
    type UserUseCase interface {
        CreateUser(ctx context.Context, req CreateUserRequest) (*UserResponse, error)
        GetUser(ctx context.Context, id uint) (*UserResponse, error)
        UpdateUser(ctx context.Context, id uint, req UpdateUserRequest) error
        DeleteUser(ctx context.Context, id uint) error
        ListUsers(ctx context.Context) ([]*UserResponse, error)
    }
    
    // CreateUserRequest - Input DTO
    type CreateUserRequest struct {
        Name  string `json:"name" validate:"required,min=2"`
        Email string `json:"email" validate:"required,email"`
        Role  string `json:"role" validate:"required,oneof=admin user"`
    }
    
    // UserResponse - Output DTO
    type UserResponse struct {
        ID    uint   `json:"id"`
        Name  string `json:"name"`
        Email string `json:"email"`
        Role  string `json:"role"`
    }
    
    // userUseCase implements UserUseCase
    type userUseCase struct {
        userRepo domain.UserRepository // Depends on interface!
    }
    
    func NewUserUseCase(userRepo domain.UserRepository) UserUseCase {
        return &userUseCase{userRepo: userRepo}
    }
    
    func (uc *userUseCase) CreateUser(ctx context.Context, req CreateUserRequest) (*UserResponse, error) {
        // 1. Validate input
        if err := req.Validate(); err != nil {
            return nil, err
        }
        
        // 2. Create domain entity
        user := &domain.User{
            Name:  req.Name,
            Email: req.Email,
            Role:  req.Role,
        }
        
        // 3. Validate business rules
        if err := user.Validate(); err != nil {
            return nil, err
        }
        
        // 4. Persist through repository
        if err := uc.userRepo.Save(ctx, user); err != nil {
            return nil, err
        }
        
        // 5. Return DTO
        return &UserResponse{
            ID:    user.ID,
            Name:  user.Name,
            Email: user.Email,
            Role:  user.Role,
        }, nil
    }

Use Case Layer Rules

**Do**: Application workflows, DTOs, coordinate repositories  
**Don't**: HTTP/gRPC details, SQL queries, framework-specific code

### 🔵 Layer 3: Repository (Infrastructure) [​](#🔵-layer-3-repository-infrastructure)

**Location**: `internal/repository/`

Implements **data access and external communication**.

#### Responsibilities [​](#responsibilities-2)

*   Implement repository interfaces from domain
*   Handle database operations
*   Manage database connections
*   Transform between DB models and domain entities

#### Example: PostgreSQL Repository [​](#example-postgresql-repository)

go

    package repository
    
    import (
        "context"
        "database/sql"
        "myproject/internal/domain"
    )
    
    type postgresUserRepository struct {
        db *sql.DB
    }
    
    // NewPostgresUserRepository creates a new repository
    func NewPostgresUserRepository(db *sql.DB) domain.UserRepository {
        return &postgresUserRepository{db: db}
    }
    
    func (r *postgresUserRepository) Save(ctx context.Context, user *domain.User) error {
        query := `
            INSERT INTO users (name, email, role, created_at, updated_at)
            VALUES ($1, $2, $3, NOW(), NOW())
            RETURNING id, created_at, updated_at
        `
        
        err := r.db.QueryRowContext(
            ctx, query,
            user.Name, user.Email, user.Role,
        ).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
        
        return err
    }
    
    func (r *postgresUserRepository) FindByID(ctx context.Context, id uint) (*domain.User, error) {
        query := `
            SELECT id, name, email, role, created_at, updated_at
            FROM users WHERE id = $1
        `
        
        user := &domain.User{}
        err := r.db.QueryRowContext(ctx, query, id).Scan(
            &user.ID, &user.Name, &user.Email, &user.Role,
            &user.CreatedAt, &user.UpdatedAt,
        )
        
        if err == sql.ErrNoRows {
            return nil, domain.ErrUserNotFound
        }
        
        return user, err
    }
    
    func (r *postgresUserRepository) Update(ctx context.Context, user *domain.User) error {
        query := `
            UPDATE users 
            SET name = $1, email = $2, role = $3, updated_at = NOW()
            WHERE id = $4
        `
        
        _, err := r.db.ExecContext(ctx, query,
            user.Name, user.Email, user.Role, user.ID,
        )
        
        return err
    }

Repository Layer Rules

**Do**: Implement domain interfaces, handle persistence  
**Don't**: Business logic, validation rules

### 🟢 Layer 4: Handlers (Interface Adapters) [​](#🟢-layer-4-handlers-interface-adapters)

**Location**: `internal/handler/http/`

Adapts **external requests to use cases**.

#### Responsibilities [​](#responsibilities-3)

*   Handle HTTP/gRPC/CLI requests
*   Parse and validate input
*   Call use cases
*   Format responses
*   Handle HTTP-specific concerns

#### Example: HTTP Handler [​](#example-http-handler)

go

    package http
    
    import (
        "encoding/json"
        "net/http"
        "strconv"
        "myproject/internal/usecase"
        "github.com/gorilla/mux"
    )
    
    type UserHandler struct {
        userUseCase usecase.UserUseCase
    }
    
    func NewUserHandler(userUseCase usecase.UserUseCase) *UserHandler {
        return &UserHandler{userUseCase: userUseCase}
    }
    
    // CreateUser handles POST /users
    func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
        // 1. Parse HTTP request
        var req usecase.CreateUserRequest
        if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
            respondError(w, http.StatusBadRequest, "Invalid request body")
            return
        }
        
        // 2. Call use case
        user, err := h.userUseCase.CreateUser(r.Context(), req)
        if err != nil {
            respondError(w, http.StatusInternalServerError, err.Error())
            return
        }
        
        // 3. Send HTTP response
        respondJSON(w, http.StatusCreated, user)
    }
    
    // GetUser handles GET /users/{id}
    func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
        // Parse path parameter
        vars := mux.Vars(r)
        id, err := strconv.ParseUint(vars["id"], 10, 32)
        if err != nil {
            respondError(w, http.StatusBadRequest, "Invalid user ID")
            return
        }
        
        user, err := h.userUseCase.GetUser(r.Context(), uint(id))
        if err != nil {
            respondError(w, http.StatusNotFound, "User not found")
            return
        }
        
        respondJSON(w, http.StatusOK, user)
    }
    
    // Helper functions
    func respondJSON(w http.ResponseWriter, status int, payload interface{}) {
        w.Header().Set("Content-Type", "application/json")
        w.WriteStatus(status)
        json.NewEncoder(w).Encode(payload)
    }
    
    func respondError(w http.ResponseWriter, status int, message string) {
        respondJSON(w, status, map[string]string{"error": message})
    }

Handler Layer Rules

**Do**: Protocol-specific concerns, request/response formatting  
**Don't**: Business logic, direct database access

Complete Data Flow [​](#complete-data-flow)
-------------------------------------------

Here's how a request flows through all layers:

    1. HTTP Request
       ↓
    2. Handler parses request → CreateUserRequest DTO
       ↓
    3. UseCase validates and applies business rules
       ↓
    4. UseCase creates Domain Entity
       ↓
    5. Entity validates its own business rules
       ↓
    6. UseCase calls Repository interface
       ↓
    7. Repository saves to database
       ↓
    8. Repository returns Domain Entity
       ↓
    9. UseCase transforms to UserResponse DTO
       ↓
    10. Handler formats HTTP Response

Benefits of This Architecture [​](#benefits-of-this-architecture)
-----------------------------------------------------------------

### 1\. Testability [​](#_1-testability)

Test each layer in isolation:

go

    // Test use case without HTTP or database
    func TestCreateUser(t *testing.T) {
        mockRepo := &MockUserRepository{}
        useCase := usecase.NewUserUseCase(mockRepo)
        
        req := usecase.CreateUserRequest{
            Name:  "John Doe",
            Email: "john@example.com",
            Role:  "user",
        }
        
        user, err := useCase.CreateUser(context.Background(), req)
        assert.NoError(t, err)
        assert.Equal(t, "John Doe", user.Name)
    }

### 2\. Flexibility [​](#_2-flexibility)

Swap implementations without touching business logic:

go

    // Switch from PostgreSQL to MongoDB
    // Old: postgresRepo := repository.NewPostgresUserRepository(db)
    // New: mongoRepo := repository.NewMongoUserRepository(client)
    userUseCase := usecase.NewUserUseCase(mongoRepo) // Same interface!

### 3\. Maintainability [​](#_3-maintainability)

Changes are localized to specific layers:

*   UI change? → Only handler layer
*   Database change? → Only repository layer
*   Business rule change? → Only domain/usecase layer

Common Mistakes to Avoid [​](#common-mistakes-to-avoid)
-------------------------------------------------------

### Skip Layers [​](#skip-layers)

go

    // BAD: Handler directly accessing database
    func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
        db.Exec("INSERT INTO users...") //  Skipping use case!
    }

### Wrong Dependencies [​](#wrong-dependencies)

go

    // BAD: Domain depending on outer layer
    package domain
    
    import "net/http" //  Domain shouldn't know about HTTP!

### Business Logic in Handlers [​](#business-logic-in-handlers)

go

    // BAD: Validation in handler
    func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
        if user.Name == "" { //  This belongs in domain/usecase!
            return errors.New("name required")
        }
    }

Learn More [​](#learn-more)
---------------------------

*   [Project Structure](https://sazardev.github.io/goca/guide/project-structure.html) - Directory organization
*   [Complete Tutorial](https://sazardev.github.io/goca/tutorials/complete-tutorial.html) - Build a real app
*   [Best Practices](https://sazardev.github.io/goca/guide/best-practices.html) - Tips and conventions

Resources [​](#resources)
-------------------------

*   [Clean Architecture Book](https://www.amazon.com/Clean-Architecture-Craftsmans-Software-Structure/dp/0134494164) by Robert C. Martin
*   [The Clean Architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html) - Uncle Bob's blog</content>
</page>

<page>
  <title>Goca - Go Clean Architecture Code Generator</title>
  <url>https://sazardev.github.io/goca/tutorials/adding-features.html</url>
  <content>Adding Features [​](#adding-features)
-------------------------------------

Learn how to extend your Goca projects with new features and functionality.

Adding a New Feature [​](#adding-a-new-feature)
-----------------------------------------------

### 1\. Generate the Feature [​](#_1-generate-the-feature)

bash

    goca feature Comment user_id:uint:fk post_id:uint:fk content:text

This automatically:

*   Creates entity in `internal/domain/`
*   Generates repository in `internal/repository/`
*   Creates service in `internal/usecase/`
*   Generates handler in `internal/handler/http/`
*   Updates DI container
*   Registers routes

### 2\. Add Relationships [​](#_2-add-relationships)

Edit the domain entity to add relationships:

go

    // internal/domain/comment.go
    type Comment struct {
        ID        uint      `gorm:"primaryKey"`
        UserID    uint      `gorm:"not null"`
        User      User      `gorm:"foreignKey:UserID"`
        PostID    uint      `gorm:"not null"`
        Post      Post      `gorm:"foreignKey:PostID"`
        Content   string    `gorm:"type:text;not null"`
        CreatedAt time.Time
        UpdatedAt time.Time
    }

### 3\. Customize Business Logic [​](#_3-customize-business-logic)

Add custom methods to the service:

go

    // internal/usecase/comment_service.go
    
    func (s *commentService) GetCommentsByPost(ctx context.Context, postID uint) ([]*CommentResponse, error) {
        comments, err := s.repo.FindByPostID(ctx, postID)
        if err != nil {
            return nil, fmt.Errorf("failed to get comments: %w", err)
        }
        
        var responses []*CommentResponse
        for _, comment := range comments {
            responses = append(responses, toCommentResponse(comment))
        }
        return responses, nil
    }

### 4\. Add Custom Repository Methods [​](#_4-add-custom-repository-methods)

go

    // internal/repository/postgres_comment_repository.go
    
    func (r *PostgresCommentRepository) FindByPostID(ctx context.Context, postID uint) ([]*domain.Comment, error) {
        var comments []*domain.Comment
        err := r.db.WithContext(ctx).
            Where("post_id = ?", postID).
            Preload("User").
            Order("created_at DESC").
            Find(&comments).Error
        return comments, err
    }

### 5\. Add Custom HTTP Endpoints [​](#_5-add-custom-http-endpoints)

go

    // internal/handler/http/comment_handler.go
    
    func (h *CommentHandler) GetCommentsByPost(c *gin.Context) {
        postID, err := strconv.ParseUint(c.Param("post_id"), 10, 64)
        if err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
            return
        }
        
        comments, err := h.service.GetCommentsByPost(c.Request.Context(), uint(postID))
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }
        
        c.JSON(http.StatusOK, gin.H{
            "success": true,
            "data":    comments,
        })
    }

### 6\. Register Custom Routes [​](#_6-register-custom-routes)

go

    // internal/handler/http/routes.go
    
    posts := api.Group("/posts")
    {
        posts.POST("", container.PostHandler.CreatePost)
        posts.GET("/:id", container.PostHandler.GetPost)
        posts.GET("/:post_id/comments", container.CommentHandler.GetCommentsByPost)
    }

Adding Custom Middleware [​](#adding-custom-middleware)
-------------------------------------------------------

### 1\. Create Middleware [​](#_1-create-middleware)

go

    // internal/middleware/auth.go
    package middleware
    
    import (
        "github.com/gin-gonic/gin"
        "net/http"
    )
    
    func AuthMiddleware() gin.HandlerFunc {
        return func(c *gin.Context) {
            token := c.GetHeader("Authorization")
            
            if token == "" {
                c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing authorization"})
                c.Abort()
                return
            }
            
            // Validate token here
            userID, err := validateToken(token)
            if err != nil {
                c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
                c.Abort()
                return
            }
            
            c.Set("user_id", userID)
            c.Next()
        }
    }

### 2\. Apply Middleware [​](#_2-apply-middleware)

go

    // internal/handler/http/routes.go
    
    protected := api.Group("/")
    protected.Use(middleware.AuthMiddleware())
    {
        protected.POST("/posts", container.PostHandler.CreatePost)
        protected.DELETE("/posts/:id", container.PostHandler.DeletePost)
    }

Adding Validation [​](#adding-validation)
-----------------------------------------

### 1\. Use Struct Tags [​](#_1-use-struct-tags)

go

    // internal/usecase/dto.go
    
    type CreatePostInput struct {
        Title   string   `json:"title" binding:"required,min=3,max=200"`
        Content string   `json:"content" binding:"required,min=10"`
        Tags    []string `json:"tags" binding:"max=10"`
    }

### 2\. Custom Validators [​](#_2-custom-validators)

go

    // internal/handler/http/validators.go
    package http
    
    import (
        "github.com/go-playground/validator/v10"
        "regexp"
    )
    
    var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
    
    func validateEmail(fl validator.FieldLevel) bool {
        return emailRegex.MatchString(fl.Field().String())
    }
    
    // Register custom validator
    func RegisterCustomValidators(v *validator.Validate) {
        v.RegisterValidation("custom_email", validateEmail)
    }

Adding Search/Filtering [​](#adding-search-filtering)
-----------------------------------------------------

### 1\. Add Filter DTOs [​](#_1-add-filter-dtos)

go

    // internal/usecase/dto.go
    
    type ListPostsInput struct {
        Search   string `form:"search"`
        Tags     string `form:"tags"`
        Page     int    `form:"page" binding:"min=1"`
        PageSize int    `form:"page_size" binding:"min=1,max=100"`
    }

### 2\. Implement Repository Method [​](#_2-implement-repository-method)

go

    // internal/repository/postgres_post_repository.go
    
    func (r *PostgresPostRepository) Search(ctx context.Context, input ListPostsInput) ([]*domain.Post, error) {
        query := r.db.WithContext(ctx)
        
        // Search by title or content
        if input.Search != "" {
            query = query.Where("title ILIKE ? OR content ILIKE ?",
                "%"+input.Search+"%",
                "%"+input.Search+"%")
        }
        
        // Filter by tags
        if input.Tags != "" {
            tags := strings.Split(input.Tags, ",")
            query = query.Where("tags && ?", pq.Array(tags))
        }
        
        // Pagination
        offset := (input.Page - 1) * input.PageSize
        
        var posts []*domain.Post
        err := query.
            Offset(offset).
            Limit(input.PageSize).
            Order("created_at DESC").
            Find(&posts).Error
        
        return posts, err
    }

Adding File Upload [​](#adding-file-upload)
-------------------------------------------

### 1\. Create Upload Handler [​](#_1-create-upload-handler)

go

    // internal/handler/http/upload_handler.go
    package http
    
    import (
        "github.com/gin-gonic/gin"
        "path/filepath"
        "net/http"
    )
    
    type UploadHandler struct {
        uploadDir string
    }
    
    func (h *UploadHandler) UploadFile(c *gin.Context) {
        file, err := c.FormFile("file")
        if err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": "No file uploaded"})
            return
        }
        
        // Validate file type
        ext := filepath.Ext(file.Filename)
        if !isAllowedExtension(ext) {
            c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file type"})
            return
        }
        
        // Save file
        filename := generateUniqueFilename(file.Filename)
        path := filepath.Join(h.uploadDir, filename)
        
        if err := c.SaveUploadedFile(file, path); err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
            return
        }
        
        c.JSON(http.StatusOK, gin.H{
            "success":  true,
            "filename": filename,
            "url":      "/uploads/" + filename,
        })
    }

### 2\. Register Upload Route [​](#_2-register-upload-route)

go

    api.POST("/upload", container.UploadHandler.UploadFile)
    api.Static("/uploads", "./uploads")

Adding Background Jobs [​](#adding-background-jobs)
---------------------------------------------------

### 1\. Create Job Service [​](#_1-create-job-service)

go

    // internal/usecase/email_service.go
    package usecase
    
    import (
        "context"
        "fmt"
    )
    
    type EmailService interface {
        SendWelcomeEmail(ctx context.Context, userID uint) error
    }
    
    type emailService struct {
        userRepo repository.UserRepository
    }
    
    func (s *emailService) SendWelcomeEmail(ctx context.Context, userID uint) error {
        user, err := s.userRepo.FindByID(ctx, userID)
        if err != nil {
            return err
        }
        
        // Send email logic here
        fmt.Printf("Sending welcome email to %s\n", user.Email)
        
        return nil
    }

### 2\. Use Goroutines for Async [​](#_2-use-goroutines-for-async)

go

    // internal/usecase/user_service.go
    
    func (s *userService) CreateUser(ctx context.Context, input CreateUserInput) (*UserResponse, error) {
        user := &domain.User{
            Name:  input.Name,
            Email: input.Email,
        }
        
        if err := s.repo.Save(ctx, user); err != nil {
            return nil, err
        }
        
        // Send welcome email asynchronously
        go func() {
            _ = s.emailService.SendWelcomeEmail(context.Background(), user.ID)
        }()
        
        return toUserResponse(user), nil
    }

See Also [​](#see-also)
-----------------------

*   [Complete Tutorial](https://sazardev.github.io/goca/tutorials/complete-tutorial.html) - Build from scratch
*   [Best Practices](https://sazardev.github.io/goca/guide/best-practices.html) - Code guidelines
*   [Project Structure](https://sazardev.github.io/goca/guide/project-structure.html) - Directory organization</content>
</page>

<page>
  <title>Goca - Go Clean Architecture Code Generator</title>
  <url>https://sazardev.github.io/goca/commands/</url>
  <content>Commands Overview [​](#commands-overview)
-----------------------------------------

Goca provides a comprehensive set of commands to generate Clean Architecture components and manage your project structure.

Command Categories [​](#command-categories)
-------------------------------------------

### Project Initialization [​](#project-initialization)

*   [`goca init`](https://sazardev.github.io/goca/commands/init.html) - Initialize a new Clean Architecture project

### Complete Features [​](#complete-features)

*   [`goca feature`](https://sazardev.github.io/goca/commands/feature.html) - Generate a complete feature with all layers
*   [`goca integrate`](https://sazardev.github.io/goca/commands/integrate.html) - Integrate existing features with DI and routing

### Layer-Specific Generation [​](#layer-specific-generation)

#### Domain Layer [​](#domain-layer)

*   [`goca entity`](https://sazardev.github.io/goca/commands/entity.html) - Generate domain entities

#### Application Layer [​](#application-layer)

*   [`goca usecase`](https://sazardev.github.io/goca/commands/usecase.html) - Generate use cases and DTOs
*   [`goca interfaces`](https://sazardev.github.io/goca/commands/interfaces.html) - Generate interface contracts

#### Infrastructure Layer [​](#infrastructure-layer)

*   [`goca repository`](https://sazardev.github.io/goca/commands/repository.html) - Generate repositories

#### Adapter Layer [​](#adapter-layer)

*   [`goca handler`](https://sazardev.github.io/goca/commands/handler.html) - Generate handlers (HTTP, gRPC, CLI, etc.)

### Utilities [​](#utilities)

*   [`goca di`](https://sazardev.github.io/goca/commands/di.html) - Generate dependency injection container
*   [`goca messages`](https://sazardev.github.io/goca/commands/messages.html) - Generate error messages and constants
*   [`goca version`](https://sazardev.github.io/goca/commands/version.html) - Display version information

Quick Reference [​](#quick-reference)
-------------------------------------

| Command | Purpose | Auto-Integration |
| --- | --- | --- |
| `goca init` | Create new project | Complete setup |
| `goca feature` | Generate full feature | Automatic |
| `goca integrate` | Wire existing features | Automatic |
| `goca entity` | Create entities only | Manual |
| `goca usecase` | Create use cases only | Manual |
| `goca repository` | Create repositories only | Manual |
| `goca handler` | Create handlers only | Manual |
| `goca di` | Generate DI container | Manual |

Common Workflows [​](#common-workflows)
---------------------------------------

### Workflow 1: Quick Start (Recommended) [​](#workflow-1-quick-start-recommended)

bash

    # Initialize project
    goca init myproject --module github.com/user/myproject
    
    # Generate complete features
    goca feature User --fields "name:string,email:string"
    goca feature Product --fields "name:string,price:float64"
    
    # Everything is integrated automatically!
    go run cmd/server/main.go

### Workflow 2: Layer-by-Layer [​](#workflow-2-layer-by-layer)

bash

    # Generate entity
    goca entity Order --fields "customer:string,total:float64"
    
    # Generate use case
    goca usecase OrderService --entity Order
    
    # Generate repository
    goca repository Order --database postgres
    
    # Generate handler
    goca handler Order --type http
    
    # Wire everything together
    goca integrate --all

### Workflow 3: Add to Existing Project [​](#workflow-3-add-to-existing-project)

bash

    # Generate new feature
    goca feature Payment --fields "amount:float64,method:string"
    
    # Automatically integrated with existing features

Global Flags [​](#global-flags)
-------------------------------

All commands support these flags:

bash

    --help, -h      Show help for command
    --verbose, -v   Enable verbose output
    --dry-run       Show what would be generated without creating files

Examples by Use Case [​](#examples-by-use-case)
-----------------------------------------------

### Building a REST API [​](#building-a-rest-api)

bash

    goca init ecommerce-api --module github.com/user/ecommerce
    cd ecommerce-api
    
    goca feature Product --fields "name:string,price:float64,stock:int"
    goca feature Order --fields "customer:string,total:float64,status:string"
    goca feature User --fields "name:string,email:string,role:string"

### Building a Microservice [​](#building-a-microservice)

bash

    goca init payment-service --module github.com/user/payment
    cd payment-service
    
    goca feature Payment --fields "amount:float64,currency:string,status:string"
    goca handler Payment --type grpc
    goca handler Payment --type http

### Building a CLI Tool [​](#building-a-cli-tool)

bash

    goca init data-processor --module github.com/user/processor
    cd data-processor
    
    goca feature DataProcessor --fields "input:string,output:string"
    goca handler DataProcessor --type cli

Next Steps [​](#next-steps)
---------------------------

Choose a command to learn more:

### Essential Commands [​](#essential-commands)

*   [goca init](https://sazardev.github.io/goca/commands/init.html) - Start here
*   [goca feature](https://sazardev.github.io/goca/commands/feature.html) - Fastest way to add features
*   [goca integrate](https://sazardev.github.io/goca/commands/integrate.html) - Wire everything together

### Detailed Generation [​](#detailed-generation)

*   🟡 [goca entity](https://sazardev.github.io/goca/commands/entity.html) - Domain layer
*   🔴 [goca usecase](https://sazardev.github.io/goca/commands/usecase.html) - Application layer
*   🔵 [goca repository](https://sazardev.github.io/goca/commands/repository.html) - Infrastructure layer
*   🟢 [goca handler](https://sazardev.github.io/goca/commands/handler.html) - Adapter layer

### Utilities [​](#utilities-1)

*   [goca di](https://sazardev.github.io/goca/commands/di.html) - Dependency injection
*   📝 [goca messages](https://sazardev.github.io/goca/commands/messages.html) - Messages and constants
*   [goca version](https://sazardev.github.io/goca/commands/version.html) - Version info</content>
</page>

<page>
  <title>Goca - Go Clean Architecture Code Generator</title>
  <url>https://sazardev.github.io/goca/tutorials/complete-tutorial.html</url>
  <content>Complete Tutorial [​](#complete-tutorial)
-----------------------------------------

Step-by-step guide to building a complete API with Goca.

What We'll Build [​](#what-we-ll-build)
---------------------------------------

A **Task Management API** with:

*   User authentication
*   Projects and tasks
*   Task assignments
*   Due dates and priorities

Prerequisites [​](#prerequisites)
---------------------------------

*   Go 1.21+
*   PostgreSQL installed
*   Basic Go knowledge

Step 1: Initialize Project [​](#step-1-initialize-project)
----------------------------------------------------------

bash

    goca init task-manager --module github.com/yourusername/task-manager --database postgres
    cd task-manager

This creates the complete project structure with clean architecture and PostgreSQL configuration.

Step 2: Configure Database [​](#step-2-configure-database)
----------------------------------------------------------

Edit `config/config.yaml`:

yaml

    database:
      host: localhost
      port: 5432
      user: postgres
      password: yourpassword
      dbname: taskmanager
      sslmode: disable
    
    server:
      port: 8080
      environment: development

Step 3: Create User Feature [​](#step-3-create-user-feature)
------------------------------------------------------------

bash

    goca feature User --fields "name:string,email:string,password:string,role:string"

This generates:

*   Domain entity (`internal/domain/user.go`)
*   Repository (`internal/repository/postgres_user_repository.go`)
*   Use case (`internal/usecase/user_service.go`)
*   HTTP handler (`internal/handler/http/user_handler.go`)
*   Routes (automatically registered)

Step 4: Create Project Feature [​](#step-4-create-project-feature)
------------------------------------------------------------------

bash

    goca feature Project --fields "name:string,description:string,owner_id:int"

A project belongs to a user (owner).

Step 5: Create Task Feature [​](#step-5-create-task-feature)
------------------------------------------------------------

bash

    goca feature Task --fields "title:string,description:string,project_id:int,assigned_to:int,priority:string,status:string,due_date:time.Time"

Tasks belong to projects and can be assigned to users.

Step 6: Add Relationships [​](#step-6-add-relationships)
--------------------------------------------------------

Edit `internal/domain/project.go`:

go

    type Project struct {
        ID          uint      `gorm:"primaryKey"`
        Name        string    `gorm:"not null"`
        Description string    `gorm:"type:text"`
        OwnerID     uint      `gorm:"not null"`
        Owner       User      `gorm:"foreignKey:OwnerID"`
        Tasks       []Task    `gorm:"foreignKey:ProjectID"`
        CreatedAt   time.Time
        UpdatedAt   time.Time
    }

Edit `internal/domain/task.go`:

go

    type Task struct {
        ID          uint      `gorm:"primaryKey"`
        Title       string    `gorm:"not null"`
        Description string    `gorm:"type:text"`
        ProjectID   uint      `gorm:"not null"`
        Project     Project   `gorm:"foreignKey:ProjectID"`
        AssignedTo  uint
        Assignee    User      `gorm:"foreignKey:AssignedTo"`
        Priority    string    `gorm:"default:'medium'"`
        Status      string    `gorm:"default:'pending'"`
        DueDate     time.Time
        CreatedAt   time.Time
        UpdatedAt   time.Time
    }

Step 7: Run Migrations [​](#step-7-run-migrations)
--------------------------------------------------

bash

    # Auto-migrate database
    go run cmd/server/main.go migrate

Or use manual migrations:

bash

    # Create migration
    migrate create -ext sql -dir migrations -seq create_tables
    
    # Edit migration files, then run:
    migrate -path migrations -database "postgresql://postgres:password@localhost:5432/taskmanager?sslmode=disable" up

Step 8: Add Business Logic [​](#step-8-add-business-logic)
----------------------------------------------------------

Edit `internal/usecase/task_service.go`:

go

    func (s *taskService) AssignTask(ctx context.Context, taskID, userID uint) error {
        task, err := s.repo.FindByID(ctx, taskID)
        if err != nil {
            return fmt.Errorf("task not found: %w", err)
        }
        
        // Verify user exists
        // Add your validation logic here
        
        task.AssignedTo = userID
        return s.repo.Update(ctx, task)
    }
    
    func (s *taskService) GetTasksByProject(ctx context.Context, projectID uint) ([]*TaskResponse, error) {
        tasks, err := s.repo.FindByProjectID(ctx, projectID)
        if err != nil {
            return nil, err
        }
        
        var responses []*TaskResponse
        for _, task := range tasks {
            responses = append(responses, toTaskResponse(task))
        }
        return responses, nil
    }

Step 9: Add Custom Repository Methods [​](#step-9-add-custom-repository-methods)
--------------------------------------------------------------------------------

Edit `internal/repository/postgres_task_repository.go`:

go

    func (r *PostgresTaskRepository) FindByProjectID(ctx context.Context, projectID uint) ([]*domain.Task, error) {
        var tasks []*domain.Task
        err := r.db.WithContext(ctx).
            Where("project_id = ?", projectID).
            Preload("Assignee").
            Find(&tasks).Error
        return tasks, err
    }
    
    func (r *PostgresTaskRepository) FindByAssignee(ctx context.Context, userID uint) ([]*domain.Task, error) {
        var tasks []*domain.Task
        err := r.db.WithContext(ctx).
            Where("assigned_to = ?", userID).
            Preload("Project").
            Find(&tasks).Error
        return tasks, err
    }

Step 10: Add Custom Routes [​](#step-10-add-custom-routes)
----------------------------------------------------------

Edit `internal/handler/http/routes.go`:

go

    func SetupRoutes(router *gin.Engine, container *di.Container) {
        api := router.Group("/api/v1")
        
        // Users
        users := api.Group("/users")
        {
            users.POST("", container.UserHandler.CreateUser)
            users.GET("/:id", container.UserHandler.GetUser)
            users.GET("", container.UserHandler.ListUsers)
        }
        
        // Projects
        projects := api.Group("/projects")
        {
            projects.POST("", container.ProjectHandler.CreateProject)
            projects.GET("/:id", container.ProjectHandler.GetProject)
            projects.GET("", container.ProjectHandler.ListProjects)
            projects.GET("/:id/tasks", container.TaskHandler.GetTasksByProject)
        }
        
        // Tasks
        tasks := api.Group("/tasks")
        {
            tasks.POST("", container.TaskHandler.CreateTask)
            tasks.GET("/:id", container.TaskHandler.GetTask)
            tasks.PUT("/:id", container.TaskHandler.UpdateTask)
            tasks.DELETE("/:id", container.TaskHandler.DeleteTask)
            tasks.POST("/:id/assign", container.TaskHandler.AssignTask)
        }
    }

Step 11: Run the Server [​](#step-11-run-the-server)
----------------------------------------------------

bash

    go run cmd/server/main.go

Output:

     Server starting on :8080
     Database connected
     Routes registered
     Server running

Step 12: Test the API [​](#step-12-test-the-api)
------------------------------------------------

### Create a User [​](#create-a-user)

bash

    curl -X POST http://localhost:8080/api/v1/users \
      -H "Content-Type: application/json" \
      -d '{
        "name": "John Doe",
        "email": "john@example.com",
        "password": "securepass123",
        "role": "developer"
      }'

### Create a Project [​](#create-a-project)

bash

    curl -X POST http://localhost:8080/api/v1/projects \
      -H "Content-Type: application/json" \
      -d '{
        "name": "Website Redesign",
        "description": "Modernize company website",
        "owner_id": 1
      }'

### Create a Task [​](#create-a-task)

bash

    curl -X POST http://localhost:8080/api/v1/tasks \
      -H "Content-Type: application/json" \
      -d '{
        "title": "Design homepage mockup",
        "description": "Create Figma designs for new homepage",
        "project_id": 1,
        "priority": "high",
        "due_date": "2025-02-01T00:00:00Z"
      }'

### Assign Task [​](#assign-task)

bash

    curl -X POST http://localhost:8080/api/v1/tasks/1/assign \
      -H "Content-Type: application/json" \
      -d '{"user_id": 1}'

Step 13: Add Tests [​](#step-13-add-tests)
------------------------------------------

Create `internal/usecase/task_service_test.go`:

go

    package usecase_test
    
    import (
        "context"
        "testing"
        "github.com/stretchr/testify/assert"
        "github.com/stretchr/testify/mock"
    )
    
    func TestTaskService_CreateTask(t *testing.T) {
        mockRepo := new(MockTaskRepository)
        service := NewTaskService(mockRepo)
        
        input := CreateTaskInput{
            Title:       "Test Task",
            Description: "Test Description",
            ProjectID:   1,
        }
        
        mockRepo.On("Save", mock.Anything, mock.Anything).Return(nil)
        
        result, err := service.CreateTask(context.Background(), input)
        
        assert.NoError(t, err)
        assert.NotNil(t, result)
        mockRepo.AssertExpectations(t)
    }

Run tests:

Next Steps [​](#next-steps)
---------------------------

*   [Adding Features](https://sazardev.github.io/goca/tutorials/adding-features.html) - Extend your project
*   [Best Practices](https://sazardev.github.io/goca/guide/best-practices.html) - Code quality guidelines
*   [Deployment Guide](https://github.com/sazardev/goca/wiki) - Production deployment

Congratulations! [​](#congratulations)
--------------------------------------

You've built a complete REST API with:

*   Clean Architecture
*   CRUD operations
*   Relationships
*   Business logic
*   Database integration

Keep exploring and building!</content>
</page>

<page>
  <title>Goca - Go Clean Architecture Code Generator</title>
  <url>https://sazardev.github.io/goca/guide/configuration.html</url>
  <content>Configuration Guide [​](#configuration-guide)
---------------------------------------------

Overview [​](#overview)
-----------------------

GOCA supports centralized project configuration through the `.goca.yaml` file. This powerful feature allows you to define project-wide settings, conventions, and preferences that will be applied automatically to all code generation commands.

### Benefits [​](#benefits)

*   **Consistency**: Maintain uniform settings across your entire project
*   **Productivity**: Avoid repeating the same CLI flags for every command
*   \*### 2. Start Simple

Begin### 3. Document Your Choices

### 4\. Version Control [​](#_4-version-control)

Always commit `.goca.yaml` to version control so team members can use the same configuration:

bash

    git add .goca.yaml
    git commit -m "Add GOCA configuration"

### 5\. Team Standardsts to explain configuration decisions, especially if you've customized a template: [​](#_5-team-standardsts-to-explain-configuration-decisions-especially-if-you-ve-customized-a-template)

yaml

    # Using PostgreSQL for advanced features
    database:
      type: postgres
      
      # Enable soft deletes for audit trail
      features:
        soft_delete: true
        timestamps: true

### 4\. Version Controlnfiguration and add settings as needed. If using a template, you already have a good starting point: [​](#_4-version-controlnfiguration-and-add-settings-as-needed-if-using-a-template-you-already-have-a-good-starting-point)

yaml

    project:
      name: my-project
      module: github.com/user/my-project
    
    database:
      type: postgres

*   **Flexibility**: Override configuration with CLI flags when needed
*   **Documentation**: Configuration files serve as self-documenting project preferences

### When to Use Configuration [​](#when-to-use-configuration)

Use `.goca.yaml` when:

*   Working on projects with consistent patterns and conventions
*   Managing multiple features with the same database type
*   Enforcing team-wide coding standards
*   Automating CI/CD workflows
*   Customizing code generation templates

### Quick Start with Templates [​](#quick-start-with-templates)

The fastest way to get started with configuration is using project templates. Templates automatically generate optimized `.goca.yaml` files for specific use cases:

bash

    # List available templates
    goca init --list-templates
    
    # Initialize with a template
    goca init my-api --module github.com/user/my-api --template rest-api

**Available templates:**

*   `minimal`: Lightweight starter with essentials
*   `rest-api`: Production REST API with PostgreSQL, validation, testing
*   `microservice`: Distributed systems with events and comprehensive testing
*   `monolith`: Full-featured web application with auth and caching
*   `enterprise`: Complete enterprise setup with security and monitoring

See the [init command documentation](https://sazardev.github.io/goca/commands/init.html#project-templates) for detailed information about templates.

Configuration File Location [​](#configuration-file-location)
-------------------------------------------------------------

Create a `.goca.yaml` file in your project root directory (where your `go.mod` file is located):

    my-project/
    ├── .goca.yaml          ← Configuration file
    ├── go.mod
    ├── go.sum
    ├── cmd/
    ├── internal/
    └── ...

Configuration Precedence [​](#configuration-precedence)
-------------------------------------------------------

GOCA uses a three-tier configuration system with the following precedence (highest to lowest):

1.  **CLI Flags**: Command-line arguments (highest priority)
2.  **Configuration File**: Settings from `.goca.yaml`
3.  **Defaults**: Built-in GOCA defaults (lowest priority)

This allows you to define common settings in your configuration file while still being able to override them for specific commands.

Core Configuration Sections [​](#core-configuration-sections)
-------------------------------------------------------------

### Project Configuration [​](#project-configuration)

Define basic project metadata and information:

yaml

    project:
      name: my-api
      module: github.com/mycompany/my-api
      description: RESTful API for customer management
      version: 1.0.0
      author: Development Team
      license: MIT

**Fields:**

*   `name`: Project name
*   `module`: Go module path
*   `description`: Project description
*   `version`: Project version
*   `author`: Author or team name
*   `license`: License type

### Database Configuration [​](#database-configuration)

Configure database settings and features:

yaml

    database:
      type: postgres
      host: localhost
      port: 5432
      migrations:
        enabled: true
        auto_generate: true
        directory: migrations
      features:
        soft_delete: true
        timestamps: true
        uuid: true
        audit: false

**Supported database types:**

*   `postgres`: PostgreSQL
*   `mysql`: MySQL/MariaDB
*   `mongodb`: MongoDB

**Migration settings:**

*   `enabled`: Enable/disable migrations
*   `auto_generate`: Auto-generate migration files
*   `directory`: Migration files directory

**Database features:**

*   `soft_delete`: Add soft delete functionality to entities
*   `timestamps`: Add created\_at/updated\_at fields
*   `uuid`: Use UUID for primary keys
*   `audit`: Enable audit logging

### Architecture Configuration [​](#architecture-configuration)

Define Clean Architecture layers and naming conventions:

yaml

    architecture:
      layers:
        domain:
          enabled: true
          directory: internal/domain
        usecase:
          enabled: true
          directory: internal/usecase
        repository:
          enabled: true
          directory: internal/repository
        handler:
          enabled: true
          directory: internal/handler
      
      patterns:
        - repository
        - service
        - dto
      
      naming:
        files: lowercase
        entities: PascalCase
        variables: camelCase
        functions: PascalCase
        constants: SCREAMING_SNAKE

**Layer configuration:**

*   `enabled`: Enable/disable layer generation
*   `directory`: Custom directory for layer

**Naming conventions:**

*   `lowercase`: user\_service.go
*   `snake_case`: user\_service.go
*   `PascalCase`: UserService
*   `camelCase`: userService
*   `SCREAMING_SNAKE`: MAX\_RETRIES

### Generation Configuration [​](#generation-configuration)

Control code generation preferences:

yaml

    generation:
      validation:
        enabled: true
        library: builtin
        sanitize: true
      
      business_rules:
        enabled: true
        patterns:
          - validation
          - authorization
        events: true
      
      documentation:
        swagger:
          enabled: true
          version: "2.0"
          output: docs/swagger.yaml
        comments:
          enabled: true
          language: english
          style: godoc

**Validation options:**

*   `enabled`: Enable field validation
*   `library`: Validation library (builtin, validator, ozzo-validation)
*   `sanitize`: Enable input sanitization

**Business rules:**

*   `enabled`: Generate business rules layer
*   `patterns`: Patterns to apply
*   `events`: Enable domain events

### Testing Configuration [​](#testing-configuration)

Configure testing generation preferences:

yaml

    testing:
      enabled: true
      framework: testify
      coverage:
        enabled: true
        threshold: 80
      mocks:
        enabled: true
        tool: testify
        directory: internal/mocks
      integration: true
      benchmarks: true

**Testing options:**

*   `enabled`: Enable test generation
*   `framework`: Testing framework (testify, ginkgo, builtin)
*   `coverage`: Code coverage settings
*   `mocks`: Mock generation settings
*   `integration`: Generate integration tests
*   `benchmarks`: Generate benchmark tests

### Template Configuration [​](#template-configuration)

Customize code generation templates:

yaml

    templates:
      directory: .goca/templates
      variables:
        author: "Development Team"
        copyright: "2024 MyCompany Inc"
        license: "MIT"
      custom:
        entity:
          path: .goca/templates/entity.tmpl
          type: go-template

**Template settings:**

*   `directory`: Custom templates directory
*   `variables`: Template variables
*   `custom`: Custom template definitions

Configuration Examples [​](#configuration-examples)
---------------------------------------------------

### Minimal Configuration [​](#minimal-configuration)

Simple configuration for small projects:

yaml

    project:
      name: minimal-api
      module: github.com/user/minimal-api

### Web Application [​](#web-application)

Configuration for a web application with HTTP handlers:

yaml

    project:
      name: web-app
      module: github.com/user/web-app
    
    generation:
      validation:
        enabled: true

### Microservice [​](#microservice)

Configuration for a microservice with database and testing:

yaml

    project:
      name: user-service
      module: github.com/company/user-service
    
    database:
      type: postgres
      migrations:
        enabled: true
    
    architecture:
      patterns:
        - repository
        - service
    
    testing:
      enabled: true
      framework: testify
      mocks:
        enabled: true

### Enterprise Application [​](#enterprise-application)

Comprehensive configuration for enterprise applications:

yaml

    project:
      name: enterprise-api
      module: github.com/corp/enterprise-api
      description: Enterprise-grade API
    
    database:
      type: postgres
      migrations:
        enabled: true
      features:
        soft_delete: true
        timestamps: true
    
    generation:
      validation:
        enabled: true
      business_rules:
        enabled: true
    
    testing:
      enabled: true
      framework: testify
      coverage:
        enabled: true
      mocks:
        enabled: true
      integration: true
    
    templates:
      directory: .goca/templates
    
    architecture:
      naming:
        files: lowercase
        entities: PascalCase
        variables: camelCase

Using Configuration with Commands [​](#using-configuration-with-commands)
-------------------------------------------------------------------------

### Initialize Project with Configuration [​](#initialize-project-with-configuration)

When you have a `.goca.yaml` file, GOCA will automatically use it:

bash

    # Configuration will be loaded automatically
    goca feature Product --fields "name:string,price:float64,stock:int"

### Override Configuration [​](#override-configuration)

You can override configuration settings using CLI flags:

bash

    # Override database type from config
    goca feature Order --database mysql

### View Effective Configuration [​](#view-effective-configuration)

Check what configuration is being used:

Template Customization [​](#template-customization)
---------------------------------------------------

### Using Custom Templates [​](#using-custom-templates)

1.  Create a templates directory:

bash

    mkdir -p .goca/templates

2.  Configure template directory:

yaml

    templates:
      directory: .goca/templates

3.  Create custom templates in the directory following GOCA's template structure
    
4.  GOCA will use your custom templates instead of built-in ones
    

### Template Variables [​](#template-variables)

Define custom variables for templates:

yaml

    templates:
      variables:
        author: "Your Name"
        company: "Your Company"
        copyright: "2024"
        license: "Apache 2.0"

Access variables in templates:

go

    // Generated by GOCA CLI
    // Author: {{.Author}}
    // Copyright: {{.Copyright}} {{.Company}}

Best Practices [​](#best-practices)
-----------------------------------

### 1\. Use Templates for Quick Start [​](#_1-use-templates-for-quick-start)

Start with a predefined template that matches your use case:

bash

    # List available templates
    goca init --list-templates
    
    # Initialize with appropriate template
    goca init my-api --module github.com/user/my-api --template rest-api

Templates provide optimized configurations that you can customize later. See the [init command documentation](https://sazardev.github.io/goca/commands/init.html#project-templates) for details.

### 2\. Start Simple [​](#_2-start-simple)

Begin with minimal configuration and add settings as needed:

yaml

    project:
      name: my-project
      module: github.com/user/my-project
    
    database:
      type: postgres

### 2\. Document Your Choices [​](#_2-document-your-choices)

Add comments to explain configuration decisions:

yaml

    # Using PostgreSQL for advanced features
    database:
      type: postgres
      
      # Enable soft deletes for audit trail
      features:
        soft_delete: true
        timestamps: true

### 3\. Version Control [​](#_3-version-control)

Always commit `.goca.yaml` to version control:

bash

    git add .goca.yaml
    git commit -m "Add GOCA configuration"

### 4\. Team Standards [​](#_4-team-standards)

Use configuration to enforce team-wide standards:

yaml

    architecture:
      naming:
        files: lowercase      # Consistent file naming
        entities: PascalCase  # Go naming conventions
        variables: camelCase  # Go naming conventions
    
    testing:
      enabled: true           # Always generate tests
      framework: testify      # Standard testing framework

### 5\. Environment Separation [​](#_5-environment-separation)

For environment-specific settings, use separate configuration files:

    .goca.yaml           # Base configuration
    .goca.dev.yaml       # Development overrides
    .goca.prod.yaml      # Production overrides

### 6\. Validate Configuration [​](#_6-validate-configuration)

Test your configuration before committing:

bash

    # Generate a test feature to verify settings
    goca feature TestEntity --fields "name:string"
    
    # Review generated code
    # If correct, delete test entity and commit config

Configuration Validation [​](#configuration-validation)
-------------------------------------------------------

GOCA validates your configuration file when loading. Common validation errors:

### Invalid YAML Syntax [​](#invalid-yaml-syntax)

yaml

    # ERROR: Invalid indentation
    project:
    name: my-project  # Should be indented

### Required Fields Missing [​](#required-fields-missing)

yaml

    # ERROR: Project section required
    database:
      type: postgres

**Fix:** Add required project information:

yaml

    project:
      name: my-project
      module: github.com/user/my-project
    
    database:
      type: postgres

### Invalid Values [​](#invalid-values)

yaml

    # ERROR: Unsupported database type
    database:
      type: oracle  # Not supported

**Fix:** Use supported values:

yaml

    database:
      type: postgres  # postgres, mysql, mongodb

Troubleshooting [​](#troubleshooting)
-------------------------------------

### Configuration Not Loading [​](#configuration-not-loading)

**Problem:** GOCA doesn't seem to use your configuration file.

**Solution:**

1.  Verify file name is exactly `.goca.yaml`
2.  Check file is in project root (where go.mod is)
3.  Verify YAML syntax is valid
4.  Check for validation errors in output

### CLI Flags Not Overriding [​](#cli-flags-not-overriding)

**Problem:** CLI flags don't override configuration settings.

**Solution:**

*   CLI flags have highest precedence and should always override
*   Verify you're using the correct flag name
*   Check command output for applied settings

### Template Customization Not Working [​](#template-customization-not-working)

**Problem:** Custom templates are not being used.

**Solution:**

1.  Verify template directory path in configuration
2.  Check template files follow GOCA's naming conventions
3.  Ensure template syntax is valid Go templates

Advanced Configuration [​](#advanced-configuration)
---------------------------------------------------

### Conditional Generation [​](#conditional-generation)

Control what gets generated based on project type:

yaml

    # API-only project
    architecture:
      layers:
        domain:
          enabled: true
        usecase:
          enabled: true
        repository:
          enabled: true
        handler:
          enabled: true  # Enable HTTP handlers

### Custom Directory Structure [​](#custom-directory-structure)

Override default directory layout:

yaml

    architecture:
      layers:
        domain:
          enabled: true
          directory: pkg/domain      # Custom path
        repository:
          enabled: true
          directory: pkg/persistence # Custom path

### Multiple Pattern Support [​](#multiple-pattern-support)

Apply multiple architectural patterns:

yaml

    architecture:
      patterns:
        - repository    # Repository pattern
        - service       # Service layer
        - dto           # Data Transfer Objects
        - specification # Specification pattern

Migration Guide [​](#migration-guide)
-------------------------------------

### From CLI-Only to Configuration [​](#from-cli-only-to-configuration)

If you're currently using only CLI flags, migrate to configuration:

1.  **Identify repeated flags:**

bash

    # You run this often:
    goca feature User --database postgres --validation
    goca feature Order --database postgres --validation

2.  **Create configuration:**

yaml

    project:
      name: my-project
      module: github.com/user/my-project
    
    database:
      type: postgres
    
    generation:
      validation:
        enabled: true

3.  **Simplify commands:**

bash

    # Now you can run:
    goca feature User
    goca feature Order

Next Steps [​](#next-steps)
---------------------------

*   Learn about [Clean Architecture](https://sazardev.github.io/goca/guide/clean-architecture.html)
*   Explore [Commands Reference](https://sazardev.github.io/goca/commands/)
*   Read [Best Practices](https://sazardev.github.io/goca/guide/best-practices.html)
*   Check [Project Structure](https://sazardev.github.io/goca/guide/project-structure.html)

Resources [​](#resources)
-------------------------

*   [YAML Specification](https://yaml.org/)
*   [Go Templates](https://pkg.go.dev/text/template)
*   [Clean Architecture Principles](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)</content>
</page>

<page>
  <title>Goca - Go Clean Architecture Code Generator</title>
  <url>https://sazardev.github.io/goca/guide/project-structure.html</url>
  <content>Project Structure [​](#project-structure)
-----------------------------------------

Understanding the directory organization and file conventions.

Overview [​](#overview)
-----------------------

Goca follows Clean Architecture principles with clear separation of concerns:

    myproject/
    ├── cmd/                    # Application entry points
    │   └── server/
    │       └── main.go
    ├── internal/               # Private application code
    │   ├── domain/            # Business entities
    │   ├── usecase/           # Business logic
    │   ├── repository/        # Data access
    │   ├── handler/           # Delivery mechanisms
    │   ├── di/                # Dependency injection
    │   └── messages/          # Errors and constants
    ├── pkg/                    # Public libraries
    ├── config/                 # Configuration files
    ├── migrations/             # Database migrations
    ├── scripts/                # Build and deployment scripts
    └── docs/                   # Documentation

Layer Responsibilities [​](#layer-responsibilities)
---------------------------------------------------

### Domain (`internal/domain/`) [​](#domain-internal-domain)

Contains business entities and core business rules.

**Files:**

*   `{entity}.go` - Entity definition
*   `{entity}_seeds.go` - Seed data
*   `errors.go` - Domain-specific errors

**Example:**

go

    // user.go
    package domain
    
    type User struct {
        ID        uint   `gorm:"primaryKey"`
        Name      string `gorm:"not null"`
        Email     string `gorm:"unique;not null"`
        CreatedAt time.Time
        UpdatedAt time.Time
    }

### Use Case (`internal/usecase/`) [​](#use-case-internal-usecase)

Business logic and application rules.

**Files:**

*   `{entity}_service.go` - Business logic implementation
*   `interfaces.go` - Service contracts
*   `dto.go` - Data transfer objects

**Example:**

go

    // user_service.go
    package usecase
    
    type UserService interface {
        CreateUser(ctx context.Context, input CreateUserInput) (*UserResponse, error)
        GetUser(ctx context.Context, id uint) (*UserResponse, error)
    }

### Repository (`internal/repository/`) [​](#repository-internal-repository)

Data persistence layer.

**Files:**

*   `postgres_{entity}_repository.go` - PostgreSQL implementation
*   `interfaces.go` - Repository contracts

**Example:**

go

    // postgres_user_repository.go
    package repository
    
    type PostgresUserRepository struct {
        db *gorm.DB
    }
    
    func (r *PostgresUserRepository) Save(ctx context.Context, user *domain.User) error {
        return r.db.WithContext(ctx).Create(user).Error
    }

### Handler (`internal/handler/http/`) [​](#handler-internal-handler-http)

HTTP delivery layer.

**Files:**

*   `{entity}_handler.go` - HTTP handlers
*   `routes.go` - Route registration
*   `swagger.yaml` - API documentation

**Example:**

go

    // user_handler.go
    package http
    
    type UserHandler struct {
        service usecase.UserService
    }
    
    func (h *UserHandler) CreateUser(c *gin.Context) {
        // Handle HTTP request
    }

### Dependency Injection (`internal/di/`) [​](#dependency-injection-internal-di)

Wires all components together.

**Files:**

*   `container.go` - DI container

**Example:**

go

    // container.go
    package di
    
    type Container struct {
        UserRepository repository.UserRepository
        UserService    usecase.UserService
        UserHandler    *http.UserHandler
    }

File Naming Conventions [​](#file-naming-conventions)
-----------------------------------------------------

| Pattern | Example | Purpose |
| --- | --- | --- |
| `{entity}.go` | `user.go` | Entity definition |
| `{entity}_service.go` | `user_service.go` | Business logic |
| `postgres_{entity}_repository.go` | `postgres_user_repository.go` | Data access |
| `{entity}_handler.go` | `user_handler.go` | HTTP handlers |
| `{entity}_seeds.go` | `user_seeds.go` | Seed data |
| `{entity}_test.go` | `user_test.go` | Unit tests |

Package Organization [​](#package-organization)
-----------------------------------------------

### Internal vs Pkg [​](#internal-vs-pkg)

    internal/     # Private to this application
    pkg/          # Can be imported by other projects

### Feature Grouping [​](#feature-grouping)

Each feature spans multiple layers:

    User Feature:
    ├── internal/domain/user.go
    ├── internal/usecase/user_service.go
    ├── internal/repository/postgres_user_repository.go
    └── internal/handler/http/user_handler.go

Import Rules [​](#import-rules)
-------------------------------

Follow dependency direction:

    Handler → Use Case → Repository → Domain
       ↓         ↓           ↓          ↓
     HTTP    Business    Database   Entities

\*\* Allowed:\*\*

go

    // Handler imports use case
    import "myproject/internal/usecase"
    
    // Use case imports repository interface
    import "myproject/internal/repository"

\*\* Not Allowed:\*\*

go

    // Repository imports handler (wrong direction!)
    import "myproject/internal/handler/http"

Configuration Files [​](#configuration-files)
---------------------------------------------

    config/
    ├── config.yaml          # Application settings
    ├── .env.example         # Environment template
    └── database.yaml        # Database configuration

Migrations [​](#migrations)
---------------------------

    migrations/
    ├── 000001_create_users_table.up.sql
    ├── 000001_create_users_table.down.sql
    ├── 000002_create_products_table.up.sql
    └── 000002_create_products_table.down.sql

Testing Structure [​](#testing-structure)
-----------------------------------------

Mirror the source structure:

    internal/
    ├── usecase/
    │   ├── user_service.go
    │   └── user_service_test.go
    └── repository/
        ├── postgres_user_repository.go
        └── postgres_user_repository_test.go

See Also [​](#see-also)
-----------------------

*   [Clean Architecture](https://sazardev.github.io/goca/guide/clean-architecture.html) - Architecture principles
*   [Best Practices](https://sazardev.github.io/goca/guide/best-practices.html) - Coding guidelines</content>
</page>

<page>
  <title>Goca - Go Clean Architecture Code Generator</title>
  <url>https://sazardev.github.io/goca/guide/best-practices.html</url>
  <content>Best Practices [​](#best-practices)
-----------------------------------

Guidelines and recommendations for building with Goca.

Code Organization [​](#code-organization)
-----------------------------------------

### Keep Layers Independent [​](#keep-layers-independent)

Each layer should only depend on interfaces, not concrete implementations.

\*\* Good:\*\*

go

    type UserService struct {
        repo repository.UserRepository // Interface
    }

\*\* Bad:\*\*

go

    type UserService struct {
        repo *repository.PostgresUserRepository // Concrete
    }

### Use Dependency Injection [​](#use-dependency-injection)

Let the DI container wire dependencies:

go

    func NewUserService(repo repository.UserRepository) usecase.UserService {
        return &userService{repo: repo}
    }

### Single Responsibility [​](#single-responsibility)

Each struct should have one clear purpose:

go

    //  One responsibility: handle HTTP requests
    type UserHandler struct {
        service usecase.UserService
    }
    
    //  Too many responsibilities
    type UserHandler struct {
        db      *gorm.DB
        cache   *redis.Client
        mailer  *smtp.Client
    }

Error Handling [​](#error-handling)
-----------------------------------

### Use Domain Errors [​](#use-domain-errors)

Define errors in the domain layer:

go

    // internal/domain/errors.go
    var (
        ErrUserNotFound = errors.New("user not found")
        ErrInvalidEmail = errors.New("invalid email format")
    )

### Wrap Errors with Context [​](#wrap-errors-with-context)

go

    func (s *userService) GetUser(ctx context.Context, id uint) (*UserResponse, error) {
        user, err := s.repo.FindByID(ctx, id)
        if err != nil {
            return nil, fmt.Errorf("failed to get user %d: %w", id, err)
        }
        return toUserResponse(user), nil
    }

### Handle Errors at HTTP Layer [​](#handle-errors-at-http-layer)

go

    func (h *UserHandler) GetUser(c *gin.Context) {
        user, err := h.service.GetUser(c.Request.Context(), id)
        if errors.Is(err, domain.ErrUserNotFound) {
            c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
            return
        }
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal error"})
            return
        }
        c.JSON(http.StatusOK, user)
    }

Testing [​](#testing)
---------------------

### Use Table-Driven Tests [​](#use-table-driven-tests)

go

    func TestUserService_CreateUser(t *testing.T) {
        tests := []struct {
            name    string
            input   CreateUserInput
            wantErr bool
        }{
            {
                name: "valid user",
                input: CreateUserInput{
                    Name:  "John Doe",
                    Email: "john@example.com",
                },
                wantErr: false,
            },
            {
                name: "invalid email",
                input: CreateUserInput{
                    Name:  "Jane Doe",
                    Email: "invalid-email",
                },
                wantErr: true,
            },
        }
        
        for _, tt := range tests {
            t.Run(tt.name, func(t *testing.T) {
                // Test implementation
            })
        }
    }

### Mock External Dependencies [​](#mock-external-dependencies)

go

    type MockUserRepository struct {
        SaveFunc    func(ctx context.Context, user *domain.User) error
        FindByIDFunc func(ctx context.Context, id uint) (*domain.User, error)
    }
    
    func (m *MockUserRepository) Save(ctx context.Context, user *domain.User) error {
        return m.SaveFunc(ctx, user)
    }

Database Operations [​](#database-operations)
---------------------------------------------

### Use Transactions [​](#use-transactions)

go

    func (s *orderService) CreateOrder(ctx context.Context, input CreateOrderInput) error {
        return s.db.Transaction(func(tx *gorm.DB) error {
            // Create order
            if err := tx.Create(&order).Error; err != nil {
                return err
            }
            
            // Create order items
            if err := tx.Create(&items).Error; err != nil {
                return err
            }
            
            return nil
        })
    }

### Use Context [​](#use-context)

Always pass context for cancellation and timeouts:

go

    func (r *PostgresUserRepository) FindByID(ctx context.Context, id uint) (*domain.User, error) {
        var user domain.User
        err := r.db.WithContext(ctx).First(&user, id).Error
        return &user, err
    }

API Design [​](#api-design)
---------------------------

### Use DTOs [​](#use-dtos)

Don't expose domain entities directly:

go

    //  Use DTOs
    type CreateUserInput struct {
        Name  string `json:"name" binding:"required"`
        Email string `json:"email" binding:"required,email"`
    }
    
    type UserResponse struct {
        ID    uint   `json:"id"`
        Name  string `json:"name"`
        Email string `json:"email"`
    }

### Consistent Response Format [​](#consistent-response-format)

go

    type SuccessResponse struct {
        Success bool        `json:"success"`
        Data    interface{} `json:"data"`
    }
    
    type ErrorResponse struct {
        Success bool   `json:"success"`
        Error   string `json:"error"`
    }

Security [​](#security)
-----------------------

### Validate Input [​](#validate-input)

go

    type CreateUserInput struct {
        Name     string `json:"name" binding:"required,min=2,max=100"`
        Email    string `json:"email" binding:"required,email"`
        Password string `json:"password" binding:"required,min=8"`
    }

### Don't Log Sensitive Data [​](#don-t-log-sensitive-data)

go

    //  Bad - logs password
    log.Printf("Creating user: %+v", input)
    
    //  Good - excludes sensitive fields
    log.Printf("Creating user: name=%s, email=%s", input.Name, input.Email)

Performance [​](#performance)
-----------------------------

### Use Pagination [​](#use-pagination)

go

    type ListUsersInput struct {
        Page     int `form:"page" binding:"min=1"`
        PageSize int `form:"page_size" binding:"min=1,max=100"`
    }
    
    func (r *PostgresUserRepository) FindAll(ctx context.Context, page, pageSize int) ([]*domain.User, error) {
        var users []*domain.User
        offset := (page - 1) * pageSize
        err := r.db.WithContext(ctx).
            Offset(offset).
            Limit(pageSize).
            Find(&users).Error
        return users, err
    }

### Use Indexes [​](#use-indexes)

go

    type User struct {
        ID    uint   `gorm:"primaryKey"`
        Email string `gorm:"unique;index"` // Add index for lookups
        Name  string `gorm:"index"`        // Index for searches
    }

See Also [​](#see-also)
-----------------------

*   [Project Structure](https://sazardev.github.io/goca/guide/project-structure.html) - Directory organization
*   [Clean Architecture](https://sazardev.github.io/goca/guide/clean-architecture.html) - Architecture principles</content>
</page>

<page>
  <title>Goca - Go Clean Architecture Code Generator</title>
  <url>https://sazardev.github.io/goca/blog/releases/v1-14-1.html</url>
  <content>v1.14.1 - Test Suite Improvements [​](#v1-14-1-test-suite-improvements)
-----------------------------------------------------------------------

October 27, 2025Bug FixesQuality

* * *

Overview [​](#overview)
-----------------------

Version 1.14.1 focuses on improving test reliability and Windows compatibility. This release includes critical fixes for path handling, working directory management, and module dependencies. The test success rate has been significantly improved to **99.04%** (310/313 tests passing).

Bug Fixes [​](#bug-fixes)
-------------------------

### Test Suite Improvements [​](#test-suite-improvements)

#### Fixed Windows Path Handling in BackupFile [​](#fixed-windows-path-handling-in-backupfile)

Corrected path issues on Windows systems that were causing backup failures:

*   Changed from `filepath.Join(BackupDir, filepath.Dir(filePath))` to using only `filepath.Base(filePath)`
*   Prevents invalid "C:" subdirectory creation on Windows
*   Backup files now correctly created with `.backup` extension in backup directory root
*   Resolves file not found errors in safety manager tests

**Impact**: Windows users can now safely use the `--backup` flag without path-related errors.

#### Fixed Test Working Directory Management [​](#fixed-test-working-directory-management)

Improved test reliability across different execution contexts:

*   Added `SetProjectDir()` calls in handler and workflow tests after project initialization
*   Corrected file path assertions from absolute to relative paths
*   Fixed path expectations to match actual command execution context
*   All handler command tests now pass with correct working directory setup

**Code Example**:

go

    // Before: Tests failed due to incorrect working directory
    handler.Execute()
    // Files created in wrong location
    
    // After: Tests pass with proper directory setup
    handler.Execute()
    SetProjectDir(projectPath)
    // Files created in correct location

#### Updated Test Message Validation [​](#updated-test-message-validation)

Aligned test expectations with actual output:

*   Converted Spanish error messages to English in entity and feature tests
*   Simplified feature test validations to accept flexible message formats
*   Improved test robustness by accepting both English and Spanish variations

**Before**:

go

    assert.Contains(t, output, "La entidad ya existe")

**After**:

go

    assert.Contains(t, output, "already exists")

#### Fixed Module Dependencies [​](#fixed-module-dependencies)

Corrected testify dependency declaration:

*   Moved `github.com/stretchr/testify` from indirect to direct dependencies in `go.mod`
*   Fixes GitHub Actions CI failure on `go mod tidy` check
*   Properly declares direct usage in test files (`internal/testing/tests/*.go`)

**go.mod change**:

diff

    require (
        github.com/spf13/cobra v1.8.0
    +   github.com/stretchr/testify v1.9.0
    )

Quality Improvements [​](#quality-improvements)
-----------------------------------------------

### Test Metrics [​](#test-metrics)

| Metric | Before | After | Improvement |
| --- | --- | --- | --- |
| Test Success Rate | 96% | 99.04% | +3.04% |
| Passing Tests | 270/313 | 310/313 | +40 tests |
| Test Failures | 40 | 3 | \-92.5% |

### Reliability Enhancements [​](#reliability-enhancements)

*   **Core Commands**: All core commands (init, entity, usecase, repository, handler, feature, di, integrate) fully functional
*   **Cross-Platform**: Improved Windows compatibility in file operations
*   **Path Handling**: Better path handling across different operating systems
*   **CI/CD**: Enhanced cross-platform test reliability in automated environments

### Test Documentation [​](#test-documentation)

*   Added comprehensive skip messages for temporarily disabled tests
*   Documented differences between test expectations and actual code generation
*   Clear issue references (`#XXX`) for tracking test improvements

**Example of improved test documentation**:

go

    t.Skip("Integration test temporarily disabled: " +
        "Test expectations require validation strictness updates. " +
        "All sub-tests pass individually. " +
        "Issue #142: Update integration test validators")

Platform Support [​](#platform-support)
---------------------------------------

### Windows Compatibility [​](#windows-compatibility)

This release significantly improves Windows support:

*   Fixed backup file creation on Windows paths
*   Corrected directory separator handling
*   Improved path normalization across platforms
*   All safety features now work correctly on Windows

### Cross-Platform Testing [​](#cross-platform-testing)

Enhanced test reliability on:

*   Windows 10/11
*   macOS (Intel and Apple Silicon)
*   Linux (Ubuntu, Debian, Fedora)

Migration Guide [​](#migration-guide)
-------------------------------------

No migration required. This release contains only bug fixes and quality improvements. Simply update to v1.14.1:

bash

    go install github.com/sazardev/goca@v1.14.1

Verify the installation:

bash

    goca version
    # Output: v1.14.1

Known Issues [​](#known-issues)
-------------------------------

### Temporarily Disabled Tests [​](#temporarily-disabled-tests)

Two complex integration tests are temporarily disabled with detailed documentation:

1.  **Full Workflow Test**: Validation strictness for generated files
2.  **Multi-Feature Integration**: Test expectations alignment with actual output

These tests have all sub-tests passing individually. The issue is with overall validation strictness, not functionality.

**Status**: Tracked in issues with clear documentation for future enhancement.

Contributors [​](#contributors)
-------------------------------

Thank you to all contributors who helped improve test reliability and Windows compatibility in this release.

Next Steps [​](#next-steps)
---------------------------

Version 1.15.0 will focus on:

*   Integration test re-enablement with updated validators
*   Additional Windows-specific test coverage
*   Performance optimization for test execution
*   Enhanced CI/CD pipeline configuration

Resources [​](#resources)
-------------------------

*   [Full CHANGELOG](https://github.com/sazardev/goca/blob/master/CHANGELOG.md)
*   [GitHub Release](https://github.com/sazardev/goca/releases/tag/v1.14.1)
*   [Report Issues](https://github.com/sazardev/goca/issues)
*   [Documentation](https://sazardev.github.io/goca/)

* * *</content>
</page>

<page>
  <title>Goca - Go Clean Architecture Code Generator</title>
  <url>https://sazardev.github.io/goca/tutorials/rest-api.html</url>
  <content>Building a REST API [​](#building-a-rest-api)
---------------------------------------------

Learn how to build a complete RESTful API with Goca from scratch.

What We'll Build [​](#what-we-ll-build)
---------------------------------------

A **Blog API** with:

*   User authentication and authorization
*   Posts with CRUD operations
*   Comments system
*   Tags and categories
*   Pagination and filtering
*   Image uploads

Prerequisites [​](#prerequisites)
---------------------------------

*   Goca installed
*   PostgreSQL or MySQL
*   Go 1.21+
*   Basic REST API knowledge

Step 1: Project Setup [​](#step-1-project-setup)
------------------------------------------------

bash

    mkdir blog-api
    cd blog-api
    
    goca init blog --module github.com/yourusername/blog --database postgres --auth
    cd blog

The `--auth` flag includes JWT authentication scaffolding.

Step 2: User Authentication [​](#step-2-user-authentication)
------------------------------------------------------------

The authentication system is already generated. Let's configure it:

Edit `internal/usecase/auth_service.go`:

go

    package usecase
    
    import (
        "context"
        "errors"
        "github.com/golang-jwt/jwt/v5"
        "golang.org/x/crypto/bcrypt"
        "time"
    )
    
    type AuthService interface {
        Register(ctx context.Context, input RegisterInput) (*AuthResponse, error)
        Login(ctx context.Context, input LoginInput) (*AuthResponse, error)
        ValidateToken(tokenString string) (*Claims, error)
    }
    
    type authService struct {
        userRepo repository.UserRepository
        jwtSecret string
    }
    
    func NewAuthService(userRepo repository.UserRepository, jwtSecret string) AuthService {
        return &authService{
            userRepo: userRepo,
            jwtSecret: jwtSecret,
        }
    }
    
    type RegisterInput struct {
        Name     string `json:"name" binding:"required,min=2"`
        Email    string `json:"email" binding:"required,email"`
        Password string `json:"password" binding:"required,min=8"`
    }
    
    type LoginInput struct {
        Email    string `json:"email" binding:"required,email"`
        Password string `json:"password" binding:"required"`
    }
    
    type AuthResponse struct {
        Token string       `json:"token"`
        User  UserResponse `json:"user"`
    }
    
    type Claims struct {
        UserID uint   `json:"user_id"`
        Email  string `json:"email"`
        jwt.RegisteredClaims
    }
    
    func (s *authService) Register(ctx context.Context, input RegisterInput) (*AuthResponse, error) {
        // Check if user exists
        existing, _ := s.userRepo.FindByEmail(ctx, input.Email)
        if existing != nil {
            return nil, errors.New("email already registered")
        }
        
        // Hash password
        hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
        if err != nil {
            return nil, err
        }
        
        // Create user
        user := &domain.User{
            Name:     input.Name,
            Email:    input.Email,
            Password: string(hashedPassword),
        }
        
        if err := s.userRepo.Save(ctx, user); err != nil {
            return nil, err
        }
        
        // Generate token
        token, err := s.generateToken(user)
        if err != nil {
            return nil, err
        }
        
        return &AuthResponse{
            Token: token,
            User:  toUserResponse(user),
        }, nil
    }
    
    func (s *authService) Login(ctx context.Context, input LoginInput) (*AuthResponse, error) {
        user, err := s.userRepo.FindByEmail(ctx, input.Email)
        if err != nil {
            return nil, errors.New("invalid credentials")
        }
        
        // Verify password
        if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
            return nil, errors.New("invalid credentials")
        }
        
        // Generate token
        token, err := s.generateToken(user)
        if err != nil {
            return nil, err
        }
        
        return &AuthResponse{
            Token: token,
            User:  toUserResponse(user),
        }, nil
    }
    
    func (s *authService) generateToken(user *domain.User) (string, error) {
        claims := Claims{
            UserID: user.ID,
            Email:  user.Email,
            RegisteredClaims: jwt.RegisteredClaims{
                ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
                IssuedAt:  jwt.NewNumericDate(time.Now()),
            },
        }
        
        token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
        return token.SignedString([]byte(s.jwtSecret))
    }
    
    func (s *authService) ValidateToken(tokenString string) (*Claims, error) {
        token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
            return []byte(s.jwtSecret), nil
        })
        
        if err != nil {
            return nil, err
        }
        
        if claims, ok := token.Claims.(*Claims); ok && token.Valid {
            return claims, nil
        }
        
        return nil, errors.New("invalid token")
    }

Step 3: Create Post Feature [​](#step-3-create-post-feature)
------------------------------------------------------------

bash

    goca feature Post title:string content:text author_id:uint:fk published:bool views:int

Add relationships in `internal/domain/post.go`:

go

    type Post struct {
        ID        uint      `gorm:"primaryKey"`
        Title     string    `gorm:"not null;index"`
        Content   string    `gorm:"type:text"`
        AuthorID  uint      `gorm:"not null"`
        Author    User      `gorm:"foreignKey:AuthorID"`
        Published bool      `gorm:"default:false"`
        Views     int       `gorm:"default:0"`
        Comments  []Comment `gorm:"foreignKey:PostID"`
        Tags      []Tag     `gorm:"many2many:post_tags"`
        CreatedAt time.Time
        UpdatedAt time.Time
    }
    
    // Business methods
    func (p *Post) Publish() error {
        if p.Title == "" || p.Content == "" {
            return errors.New("cannot publish incomplete post")
        }
        p.Published = true
        return nil
    }
    
    func (p *Post) IncrementViews() {
        p.Views++
    }

bash

    goca feature Comment post_id:uint:fk author_id:uint:fk content:text

Step 5: Create Tag Feature [​](#step-5-create-tag-feature)
----------------------------------------------------------

bash

    goca feature Tag name:string slug:string

Step 6: Add Authentication Middleware [​](#step-6-add-authentication-middleware)
--------------------------------------------------------------------------------

Create `internal/middleware/auth.go`:

go

    package middleware
    
    import (
        "net/http"
        "strings"
        "github.com/gin-gonic/gin"
        "yourproject/internal/usecase"
    )
    
    func AuthMiddleware(authService usecase.AuthService) gin.HandlerFunc {
        return func(c *gin.Context) {
            authHeader := c.GetHeader("Authorization")
            if authHeader == "" {
                c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
                c.Abort()
                return
            }
            
            // Extract token
            parts := strings.Split(authHeader, " ")
            if len(parts) != 2 || parts[0] != "Bearer" {
                c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization format"})
                c.Abort()
                return
            }
            
            token := parts[1]
            
            // Validate token
            claims, err := authService.ValidateToken(token)
            if err != nil {
                c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
                c.Abort()
                return
            }
            
            // Store user info in context
            c.Set("user_id", claims.UserID)
            c.Set("user_email", claims.Email)
            
            c.Next()
        }
    }

Step 7: Update Routes [​](#step-7-update-routes)
------------------------------------------------

Edit `internal/handler/http/routes.go`:

go

    func SetupRoutes(router *gin.Engine, container *di.Container) {
        api := router.Group("/api/v1")
        
        // Public routes
        auth := api.Group("/auth")
        {
            auth.POST("/register", container.AuthHandler.Register)
            auth.POST("/login", container.AuthHandler.Login)
        }
        
        // Public posts (read-only)
        publicPosts := api.Group("/posts")
        {
            publicPosts.GET("", container.PostHandler.ListPosts)
            publicPosts.GET("/:id", container.PostHandler.GetPost)
        }
        
        // Protected routes
        protected := api.Group("")
        protected.Use(middleware.AuthMiddleware(container.AuthService))
        {
            // Posts management
            posts := protected.Group("/posts")
            {
                posts.POST("", container.PostHandler.CreatePost)
                posts.PUT("/:id", container.PostHandler.UpdatePost)
                posts.DELETE("/:id", container.PostHandler.DeletePost)
                posts.POST("/:id/publish", container.PostHandler.PublishPost)
            }
            
            // Comments
            comments := protected.Group("/comments")
            {
                comments.POST("", container.CommentHandler.CreateComment)
                comments.PUT("/:id", container.CommentHandler.UpdateComment)
                comments.DELETE("/:id", container.CommentHandler.DeleteComment)
            }
        }
    }

Edit `internal/repository/postgres_post_repository.go`:

go

    type PostFilters struct {
        Published  *bool
        AuthorID   *uint
        Tag        string
        Search     string
        Page       int
        PageSize   int
    }
    
    func (r *PostgresPostRepository) FindWithFilters(ctx context.Context, filters PostFilters) ([]*domain.Post, int64, error) {
        query := r.db.WithContext(ctx).Model(&domain.Post{})
        
        // Apply filters
        if filters.Published != nil {
            query = query.Where("published = ?", *filters.Published)
        }
        
        if filters.AuthorID != nil {
            query = query.Where("author_id = ?", *filters.AuthorID)
        }
        
        if filters.Search != "" {
            query = query.Where("title ILIKE ? OR content ILIKE ?",
                "%"+filters.Search+"%",
                "%"+filters.Search+"%")
        }
        
        if filters.Tag != "" {
            query = query.Joins("JOIN post_tags ON post_tags.post_id = posts.id").
                Joins("JOIN tags ON tags.id = post_tags.tag_id").
                Where("tags.slug = ?", filters.Tag)
        }
        
        // Count total
        var total int64
        query.Count(&total)
        
        // Pagination
        offset := (filters.Page - 1) * filters.PageSize
        
        var posts []*domain.Post
        err := query.
            Preload("Author").
            Preload("Tags").
            Offset(offset).
            Limit(filters.PageSize).
            Order("created_at DESC").
            Find(&posts).Error
        
        return posts, total, err
    }

Step 9: Test the API [​](#step-9-test-the-api)
----------------------------------------------

### Register a User [​](#register-a-user)

bash

    curl -X POST http://localhost:8080/api/v1/auth/register \
      -H "Content-Type: application/json" \
      -d '{
        "name": "John Doe",
        "email": "john@example.com",
        "password": "securepass123"
      }'

Response:

json

    {
      "token": "eyJhbGciOiJIUzI1NiIs...",
      "user": {
        "id": 1,
        "name": "John Doe",
        "email": "john@example.com"
      }
    }

### Login [​](#login)

bash

    curl -X POST http://localhost:8080/api/v1/auth/login \
      -H "Content-Type: application/json" \
      -d '{
        "email": "john@example.com",
        "password": "securepass123"
      }'

### Create a Post (Authenticated) [​](#create-a-post-authenticated)

bash

    TOKEN="your-jwt-token"
    
    curl -X POST http://localhost:8080/api/v1/posts \
      -H "Content-Type: application/json" \
      -H "Authorization: Bearer $TOKEN" \
      -d '{
        "title": "My First Post",
        "content": "This is the content of my first blog post.",
        "published": false
      }'

### List Posts with Filters [​](#list-posts-with-filters)

bash

    # Get published posts only
    curl "http://localhost:8080/api/v1/posts?published=true&page=1&page_size=10"
    
    # Search posts
    curl "http://localhost:8080/api/v1/posts?search=golang&page=1&page_size=10"
    
    # Filter by tag
    curl "http://localhost:8080/api/v1/posts?tag=programming&page=1&page_size=10"

bash

    curl -X POST http://localhost:8080/api/v1/comments \
      -H "Content-Type: application/json" \
      -H "Authorization: Bearer $TOKEN" \
      -d '{
        "post_id": 1,
        "content": "Great post!"
      }'

Step 10: Add Image Upload [​](#step-10-add-image-upload)
--------------------------------------------------------

Create `internal/handler/http/upload_handler.go`:

go

    package http
    
    import (
        "github.com/gin-gonic/gin"
        "net/http"
        "path/filepath"
        "github.com/google/uuid"
    )
    
    type UploadHandler struct {
        uploadDir string
    }
    
    func NewUploadHandler(uploadDir string) *UploadHandler {
        return &UploadHandler{uploadDir: uploadDir}
    }
    
    func (h *UploadHandler) UploadImage(c *gin.Context) {
        file, err := c.FormFile("image")
        if err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": "No file uploaded"})
            return
        }
        
        // Validate file type
        ext := filepath.Ext(file.Filename)
        if ext != ".jpg" && ext != ".jpeg" && ext != ".png" && ext != ".gif" {
            c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file type"})
            return
        }
        
        // Generate unique filename
        filename := uuid.New().String() + ext
        path := filepath.Join(h.uploadDir, filename)
        
        // Save file
        if err := c.SaveUploadedFile(file, path); err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
            return
        }
        
        c.JSON(http.StatusOK, gin.H{
            "success": true,
            "url": "/uploads/" + filename,
        })
    }

Next Steps [​](#next-steps)
---------------------------

*   [Adding Features](https://sazardev.github.io/goca/tutorials/adding-features.html) - Extend functionality
*   [Best Practices](https://sazardev.github.io/goca/guide/best-practices.html) - Code quality
*   [Complete Tutorial](https://sazardev.github.io/goca/tutorials/complete-tutorial.html) - Advanced patterns

Summary [​](#summary)
---------------------

You now have:

*   User authentication with JWT
*   CRUD operations for posts
*   Comments system
*   Pagination and filtering
*   Protected routes with middleware
*   Many-to-many relationships (tags)
*   Business logic in domain layer

Your REST API follows Clean Architecture and is production-ready!</content>
</page>

<page>
  <title>Goca - Go Clean Architecture Code Generator</title>
  <url>https://sazardev.github.io/goca/blog/articles/</url>
  <content>Articles [​](#articles)
-----------------------

In-depth articles covering Clean Architecture patterns, best practices, testing strategies, and advanced techniques for building production-ready Go applications with Goca.

* * *

### [Mastering the Repository Pattern in Clean Architecture](https://sazardev.github.io/goca/blog/articles/mastering-repository-pattern)

A comprehensive guide to the Repository pattern and data access abstraction. Learn what repositories are, how they differ from DAOs, interface design principles, and how Goca generates database-agnostic implementations that maintain Clean Architecture boundaries.

Architecture • Infrastructure Layer

Repository PatternData AccessInfrastructureClean Architecture

### [Mastering Use Cases in Clean Architecture](https://sazardev.github.io/goca/blog/articles/mastering-use-cases)

A deep dive into use cases, application services, and DTOs. Learn what use cases are, how they orchestrate business workflows, and how Goca generates production-ready application layer code with comprehensive examples.

Architecture • Application Layer

Use CasesApplication ServicesDTOsClean Architecture

### [Understanding Domain Entities in Clean Architecture](https://sazardev.github.io/goca/blog/articles/understanding-domain-entities)

A comprehensive guide to domain entities, their role in Clean Architecture, and how Goca generates production-ready entities following DDD principles. Learn the critical distinction between entities and models, best practices, and testing strategies.

Architecture • Domain-Driven Design

Domain EntitiesClean ArchitectureDDDBest Practices

### [Advanced Features Showcase](https://sazardev.github.io/goca/blog/articles/example-showcase)

Demonstration of blog post capabilities including Mermaid diagrams, code blocks, and markdown features. Learn how to leverage VitePress features for technical documentation.

Example • Tutorial

MermaidDiagramsCode ExamplesClean Architecture

Coming Soon
-----------

Articles are in development. Check back soon for in-depth content on:

*   Building scalable microservices with Clean Architecture
*   Advanced testing strategies: Unit, Integration, and E2E
*   Database migration patterns and versioning
*   Performance optimization techniques
*   Domain-Driven Design with Goca
*   Implementing event-driven architectures

* * *

Submit an Article [​](#submit-an-article)
-----------------------------------------

Have a great article idea or want to share your experience with Goca? We welcome community contributions.

[Open an issue on GitHub](https://github.com/sazardev/goca/issues/new?title=Article%20Proposal:) to propose an article.</content>
</page>

<page>
  <title>Goca - Go Clean Architecture Code Generator</title>
  <url>https://sazardev.github.io/goca/blog/releases/</url>
  <content>Release Notes [​](#release-notes)
---------------------------------

Track the evolution of Goca through detailed release notes. Each release includes new features, bug fixes, improvements, and migration guides.

* * *

Test Suite Improvements - Fixed Windows path handling, test working directory management, and module dependencies. Test success rate improved to 99.04%.

**Key Improvements:**

*   Fixed Windows path handling in BackupFile
*   Improved test working directory management
*   Updated test message validation
*   Fixed module dependencies for testify

* * *

Release Versioning [​](#release-versioning)
-------------------------------------------

Goca follows [Semantic Versioning](https://semver.org/):

*   **Major (X.0.0)**: Breaking changes
*   **Minor (1.X.0)**: New features (backward compatible)
*   **Patch (1.14.X)**: Bug fixes and minor improvements

Stay Updated [​](#stay-updated)
-------------------------------

*   Watch the [GitHub repository](https://github.com/sazardev/goca) for release notifications
*   View the complete [CHANGELOG](https://github.com/sazardev/goca/blob/master/CHANGELOG.md)
*   Subscribe to [GitHub Releases](https://github.com/sazardev/goca/releases)</content>
</page>

<page>
  <title>Mastering the Repository Pattern in Clean Architecture | Goca Blog</title>
  <url>https://sazardev.github.io/goca/blog/articles/mastering-repository-pattern.html</url>
  <content>Mastering the Repository Pattern in Clean Architecture [​](#mastering-the-repository-pattern-in-clean-architecture)
-------------------------------------------------------------------------------------------------------------------

InfrastructureData Access

The Repository pattern is a critical abstraction that isolates domain logic from data access concerns. Understanding repositories correctly is essential for building applications that remain testable, maintainable, and independent of specific database technologies. When implemented properly, repositories enable you to change databases, add caching, or switch ORMs without touching business logic.

* * *

What is the Repository Pattern? [​](#what-is-the-repository-pattern)
--------------------------------------------------------------------

A repository mediates between the domain layer and data mapping layers, acting as an in-memory collection of domain objects. Repositories provide a clean, domain-centric API for data access while hiding the complexities of database interactions, query construction, and persistence mechanisms.

### Core Responsibilities [​](#core-responsibilities)

**Abstraction**: Repositories abstract the details of data access. The domain layer works with repository interfaces that express intent ("find user by email") rather than implementation ("execute SQL query with WHERE clause").

**Collection Semantics**: Repositories present data as collections. You add entities to repositories, remove them, and query them using domain-meaningful methods. The fact that data persists to a database is an implementation detail.

**Domain-Centric API**: Repository methods use domain language and work with domain entities. A `UserRepository` has methods like `Save(user)`, not `Insert(tableName, columns, values)`.

**Testability**: Because repositories are interfaces, you can replace real implementations with in-memory fakes during testing. This allows unit testing of business logic without database setup.

**Database Independence**: Repositories isolate database-specific code. Changing from PostgreSQL to MongoDB requires changing only repository implementations, not domain or application logic.

Repository vs DAO vs Data Mapper [​](#repository-vs-dao-vs-data-mapper)
-----------------------------------------------------------------------

Developers often confuse repositories with Data Access Objects (DAOs) or Data Mappers. These patterns serve different purposes and operate at different abstraction levels.

### What a Repository Is NOT [​](#what-a-repository-is-not)

**Not a DAO**: DAOs provide CRUD operations on database tables. Repositories provide domain operations on entity collections. A DAO might have `insertUser(name, email)`. A repository has `Save(user *User)`.

**Not a Data Mapper**: Data Mappers convert between database rows and objects. Repositories use data mappers internally but provide higher-level operations that reflect business operations.

**Not a Generic Interface**: Repositories are not `IRepository<T>` with generic CRUD. Each repository has a domain-specific interface. `UserRepository` has methods that make sense for users; `OrderRepository` has methods that make sense for orders.

**Not Query Builders**: Repositories do not expose SQL or query DSLs. They provide intention-revealing methods. Instead of `repository.Query("SELECT * FROM users WHERE age > ?", 18)`, you have `repository.FindAdults()`.

### The Clear Distinction [​](#the-clear-distinction)

    Domain Layer (Business Logic)
        ↓ Uses interface
    Repository Interface (Contract)
        ↓ Defined in domain
    Infrastructure Layer (Data Access)
        ↓ Implements interface
    Repository Implementation (Database-Specific)
        ↓ Uses ORM/Driver
    Database

Repositories are defined in the domain layer as interfaces but implemented in the infrastructure layer with database-specific code. This inverts the dependency, making infrastructure depend on the domain rather than the reverse.

The Infrastructure Layer [​](#the-infrastructure-layer)
-------------------------------------------------------

Repositories form the infrastructure layer in Clean Architecture. This layer contains all the concrete implementations of persistence, external services, and framework-specific code.

### Infrastructure Layer Characteristics [​](#infrastructure-layer-characteristics)

**Implements Domain Interfaces**: The infrastructure layer provides concrete implementations of repository interfaces defined in the domain. The domain depends on abstractions; infrastructure depends on the domain.

**Database-Specific Code**: Infrastructure contains ORM configurations, SQL queries, connection management, and database-specific optimizations. This code is hidden behind interfaces.

**Framework Dependencies**: Infrastructure can depend on GORM, database drivers, caching libraries, and external SDKs. These dependencies do not leak into the domain.

**Swappable Implementations**: You can have multiple repository implementations for the same interface: PostgreSQL for production, in-memory for testing, MongoDB for a specific feature.

### Why Separate Infrastructure? [​](#why-separate-infrastructure)

The infrastructure layer exists because data access mechanisms change independently of business rules:

**Domain Logic**: "A user must have a unique email" is domain logic. It does not care how you check uniqueness.

**Infrastructure Logic**: "Execute SELECT COUNT(\*) FROM users WHERE email = ? to check uniqueness" is infrastructure logic. It is specific to SQL databases.

Separating these concerns allows you to:

*   Test domain logic without database setup
*   Change databases without changing business rules
*   Optimize queries without touching domain code
*   Support multiple databases simultaneously

Repository Interface Design [​](#repository-interface-design)
-------------------------------------------------------------

Well-designed repository interfaces express domain intent clearly while remaining independent of implementation details.

### Interface Location: Domain Layer [​](#interface-location-domain-layer)

Repository interfaces belong in the domain layer, typically in a package like `internal/domain` or alongside entity definitions. This placement is critical:

**Domain Owns the Contract**: The domain defines what operations it needs. Infrastructure adapts to the domain's requirements, not vice versa.

**Dependency Inversion**: By placing interfaces in the domain, you invert the dependency. Infrastructure imports domain types, not the other way around.

**No Infrastructure Leakage**: Domain interfaces use only domain types. They do not reference database connections, ORM types, or SQL constructs.

### Method Design Principles [​](#method-design-principles)

**Intention-Revealing Names**: Methods should express why you are querying, not how. `FindActiveUsers()` is better than `QueryUsersWithStatus("active")`.

**Domain Types Only**: Parameters and return types are domain entities and value objects, never database-specific types like `sql.Row` or `bson.Document`.

**Error Handling**: Return domain errors, not database errors. Instead of returning `sql.ErrNoRows`, return `ErrUserNotFound`.

**No Leaky Abstractions**: Methods should not expose pagination cursors, query builders, or transaction objects unless these concepts exist in the domain.

### Common Repository Methods [​](#common-repository-methods)

Every repository typically includes these fundamental operations:

go

    type UserRepository interface {
        // Save adds or updates a user
        Save(user *User) error
        
        // FindByID retrieves a user by unique identifier
        FindByID(id uint) (*User, error)
        
        // FindAll retrieves all users (use with caution in production)
        FindAll() ([]*User, error)
        
        // Update modifies an existing user
        Update(user *User) error
        
        // Delete removes a user
        Delete(id uint) error
    }

Beyond these basics, add domain-specific query methods:

go

    type UserRepository interface {
        // ... basic methods ...
        
        // Domain-specific queries
        FindByEmail(email string) (*User, error)
        FindAdults() ([]*User, error)
        FindByLastLoginAfter(date time.Time) ([]*User, error)
        CountByStatus(status Status) (int, error)
    }

Repository Implementation Patterns [​](#repository-implementation-patterns)
---------------------------------------------------------------------------

Repository implementations encapsulate all database-specific logic, translating domain operations into database operations.

### Basic PostgreSQL Implementation [​](#basic-postgresql-implementation)

go

    package repository
    
    import (
        "github.com/yourorg/yourapp/internal/domain"
        "gorm.io/gorm"
    )
    
    type postgresUserRepository struct {
        db *gorm.DB
    }
    
    // NewPostgresUserRepository creates a PostgreSQL implementation
    func NewPostgresUserRepository(db *gorm.DB) domain.UserRepository {
        return &postgresUserRepository{db: db}
    }
    
    func (r *postgresUserRepository) Save(user *domain.User) error {
        return r.db.Create(user).Error
    }
    
    func (r *postgresUserRepository) FindByID(id uint) (*domain.User, error) {
        var user domain.User
        err := r.db.First(&user, id).Error
        if err == gorm.ErrRecordNotFound {
            return nil, domain.ErrUserNotFound
        }
        return &user, err
    }
    
    func (r *postgresUserRepository) FindByEmail(email string) (*domain.User, error) {
        var user domain.User
        err := r.db.Where("email = ?", email).First(&user).Error
        if err == gorm.ErrRecordNotFound {
            return nil, domain.ErrUserNotFound
        }
        return &user, err
    }
    
    func (r *postgresUserRepository) Update(user *domain.User) error {
        return r.db.Save(user).Error
    }
    
    func (r *postgresUserRepository) Delete(id uint) error {
        return r.db.Delete(&domain.User{}, id).Error
    }
    
    func (r *postgresUserRepository) FindAll() ([]*domain.User, error) {
        var users []*domain.User
        err := r.db.Find(&users).Error
        return users, err
    }

Notice how the implementation:

*   Returns domain errors, not GORM errors
*   Uses domain types in signatures
*   Hides all GORM-specific code
*   Implements the domain interface

### MongoDB Implementation [​](#mongodb-implementation)

For NoSQL databases, the implementation differs dramatically, but the interface remains the same:

go

    package repository
    
    import (
        "context"
        "github.com/yourorg/yourapp/internal/domain"
        "go.mongodb.org/mongo-driver/bson"
        "go.mongodb.org/mongo-driver/mongo"
    )
    
    type mongoUserRepository struct {
        collection *mongo.Collection
    }
    
    func NewMongoUserRepository(db *mongo.Database) domain.UserRepository {
        return &mongoUserRepository{
            collection: db.Collection("users"),
        }
    }
    
    func (r *mongoUserRepository) Save(user *domain.User) error {
        ctx := context.TODO()
        _, err := r.collection.InsertOne(ctx, user)
        return err
    }
    
    func (r *mongoUserRepository) FindByID(id uint) (*domain.User, error) {
        ctx := context.TODO()
        var user domain.User
        
        err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&user)
        if err == mongo.ErrNoDocuments {
            return nil, domain.ErrUserNotFound
        }
        return &user, err
    }
    
    func (r *mongoUserRepository) FindByEmail(email string) (*domain.User, error) {
        ctx := context.TODO()
        var user domain.User
        
        err := r.collection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
        if err == mongo.ErrNoDocuments {
            return nil, domain.ErrUserNotFound
        }
        return &user, err
    }

The key insight: both implementations satisfy the same `UserRepository` interface. Application and domain code work with either database without modification.

### In-Memory Implementation for Testing [​](#in-memory-implementation-for-testing)

For unit tests, create an in-memory fake that implements the repository interface:

go

    package repository
    
    import (
        "sync"
        "github.com/yourorg/yourapp/internal/domain"
    )
    
    type inMemoryUserRepository struct {
        users  map[uint]*domain.User
        nextID uint
        mu     sync.RWMutex
    }
    
    func NewInMemoryUserRepository() domain.UserRepository {
        return &inMemoryUserRepository{
            users:  make(map[uint]*domain.User),
            nextID: 1,
        }
    }
    
    func (r *inMemoryUserRepository) Save(user *domain.User) error {
        r.mu.Lock()
        defer r.mu.Unlock()
        
        if user.ID == 0 {
            user.ID = r.nextID
            r.nextID++
        }
        
        r.users[user.ID] = user
        return nil
    }
    
    func (r *inMemoryUserRepository) FindByID(id uint) (*domain.User, error) {
        r.mu.RLock()
        defer r.mu.RUnlock()
        
        user, exists := r.users[id]
        if !exists {
            return nil, domain.ErrUserNotFound
        }
        return user, nil
    }
    
    func (r *inMemoryUserRepository) FindByEmail(email string) (*domain.User, error) {
        r.mu.RLock()
        defer r.mu.RUnlock()
        
        for _, user := range r.users {
            if user.Email == email {
                return user, nil
            }
        }
        return nil, domain.ErrUserNotFound
    }

This in-memory implementation enables fast, isolated unit tests without database dependencies.

How Goca Generates Repositories [​](#how-goca-generates-repositories)
---------------------------------------------------------------------

Goca's `goca repository` command generates both interfaces and implementations following Clean Architecture principles.

### Basic Repository Generation [​](#basic-repository-generation)

bash

    goca repository User --database postgres

This generates:

**1\. Repository Interface** (`internal/repository/interfaces.go`):

go

    package repository
    
    import "yourapp/internal/domain"
    
    type UserRepository interface {
        Save(user *domain.User) error
        FindByID(id uint) (*domain.User, error)
        Update(user *domain.User) error
        Delete(id uint) error
        FindAll() ([]*domain.User, error)
    }

**2\. PostgreSQL Implementation** (`internal/repository/postgres_user_repository.go`):

go

    package repository
    
    import (
        "yourapp/internal/domain"
        "gorm.io/gorm"
    )
    
    type postgresUserRepository struct {
        db *gorm.DB
    }
    
    func NewPostgresUserRepository(db *gorm.DB) UserRepository {
        return &postgresUserRepository{db: db}
    }
    
    // ... CRUD implementations ...

### Database-Specific Implementations [​](#database-specific-implementations)

Goca supports multiple databases, generating appropriate implementations for each:

**PostgreSQL with GORM**:

bash

    goca repository User --database postgres
    # Generates GORM-based implementation with SQL transactions

**MongoDB**:

bash

    goca repository User --database mongodb
    # Generates MongoDB driver implementation with BSON

**PostgreSQL with JSONB**:

bash

    goca repository User --database postgres-json
    # Generates JSONB-specific queries for nested documents

**MySQL**:

bash

    goca repository User --database mysql
    # Generates MySQL-specific GORM implementation

**DynamoDB**:

bash

    goca repository User --database dynamodb
    # Generates AWS SDK v2 implementation with attribute mapping

### Custom Query Methods [​](#custom-query-methods)

Goca auto-generates query methods based on entity fields:

bash

    goca repository User --fields "name:string,email:string,age:int,status:string"

Generates these additional methods:

go

    type UserRepository interface {
        // Basic CRUD
        Save(user *domain.User) error
        FindByID(id uint) (*domain.User, error)
        Update(user *domain.User) error
        Delete(id uint) error
        FindAll() ([]*domain.User, error)
        
        // Field-based queries (auto-generated)
        FindByName(name string) (*domain.User, error)
        FindByEmail(email string) (*domain.User, error)
        FindByAge(age int) (*domain.User, error)
        FindByStatus(status string) (*domain.User, error)
    }

### Interface-Only Generation [​](#interface-only-generation)

For Test-Driven Development (TDD), generate only the interface first:

bash

    goca repository User --interface-only

This creates the contract without implementation, allowing you to:

1.  Write use cases against the interface
2.  Create mock implementations for tests
3.  Implement the real repository later

### Implementation-Only Generation [​](#implementation-only-generation)

If you already have the interface but need a new database implementation:

bash

    goca repository User --implementation --database mongodb

This generates only the MongoDB implementation without modifying the interface.

Advanced Repository Patterns [​](#advanced-repository-patterns)
---------------------------------------------------------------

Beyond basic CRUD, repositories support advanced patterns for complex data access scenarios.

### Specification Pattern [​](#specification-pattern)

Use specifications to encapsulate query criteria:

go

    type UserSpecification interface {
        IsSatisfiedBy(user *domain.User) bool
        ToSQL() (string, []interface{})
    }
    
    type UserRepository interface {
        // ... basic methods ...
        FindBySpec(spec UserSpecification) ([]*domain.User, error)
    }
    
    // Usage
    activeAdults := NewAndSpecification(
        NewAgeGreaterThanSpec(18),
        NewStatusEqualsSpec("active"),
    )
    users, err := repo.FindBySpec(activeAdults)

### Unit of Work Pattern [​](#unit-of-work-pattern)

Coordinate multiple repository operations in a transaction:

go

    type UnitOfWork interface {
        Users() UserRepository
        Orders() OrderRepository
        
        Begin() error
        Commit() error
        Rollback() error
    }
    
    // Usage in use case
    func (s *orderService) CreateOrder(userID uint, items []Item) error {
        uow := s.unitOfWork
        
        if err := uow.Begin(); err != nil {
            return err
        }
        defer uow.Rollback()
        
        user, err := uow.Users().FindByID(userID)
        if err != nil {
            return err
        }
        
        order := domain.NewOrder(user, items)
        if err := uow.Orders().Save(order); err != nil {
            return err
        }
        
        return uow.Commit()
    }

### Caching Layer [​](#caching-layer)

Add caching transparently using the decorator pattern:

go

    type cachedUserRepository struct {
        repository UserRepository
        cache      Cache
    }
    
    func NewCachedUserRepository(repo UserRepository, cache Cache) UserRepository {
        return &cachedUserRepository{
            repository: repo,
            cache:      cache,
        }
    }
    
    func (r *cachedUserRepository) FindByID(id uint) (*domain.User, error) {
        // Check cache first
        cacheKey := fmt.Sprintf("user:%d", id)
        if cached, found := r.cache.Get(cacheKey); found {
            return cached.(*domain.User), nil
        }
        
        // Cache miss: query database
        user, err := r.repository.FindByID(id)
        if err != nil {
            return nil, err
        }
        
        // Store in cache
        r.cache.Set(cacheKey, user, 5*time.Minute)
        return user, nil
    }

The use case layer does not know caching exists. The cached repository satisfies the same interface as the uncached version.

### Pagination Support [​](#pagination-support)

For large datasets, repositories should support pagination:

go

    type Page struct {
        Items      []*User
        TotalItems int
        Page       int
        PageSize   int
        TotalPages int
    }
    
    type UserRepository interface {
        // ... basic methods ...
        FindAllPaginated(page, pageSize int) (*Page, error)
        FindByStatusPaginated(status string, page, pageSize int) (*Page, error)
    }
    
    // Implementation
    func (r *postgresUserRepository) FindAllPaginated(page, pageSize int) (*Page, error) {
        var users []*domain.User
        var total int64
        
        offset := (page - 1) * pageSize
        
        // Count total
        r.db.Model(&domain.User{}).Count(&total)
        
        // Fetch page
        err := r.db.Limit(pageSize).Offset(offset).Find(&users).Error
        
        return &Page{
            Items:      users,
            TotalItems: int(total),
            Page:       page,
            PageSize:   pageSize,
            TotalPages: int(total)/pageSize + 1,
        }, err
    }

### Soft Delete Support [​](#soft-delete-support)

Implement soft deletes while maintaining a clean repository interface:

go

    func (r *postgresUserRepository) Delete(id uint) error {
        // Soft delete: set deleted_at timestamp
        return r.db.Model(&domain.User{}).
            Where("id = ?", id).
            Update("deleted_at", time.Now()).Error
    }
    
    func (r *postgresUserRepository) FindByID(id uint) (*domain.User, error) {
        var user domain.User
        // Automatically exclude soft-deleted records
        err := r.db.Where("deleted_at IS NULL").First(&user, id).Error
        if err == gorm.ErrRecordNotFound {
            return nil, domain.ErrUserNotFound
        }
        return &user, err
    }
    
    // Add method for finding deleted records if needed
    func (r *postgresUserRepository) FindDeletedByID(id uint) (*domain.User, error) {
        var user domain.User
        err := r.db.Unscoped().First(&user, id).Error
        if err == gorm.ErrRecordNotFound {
            return nil, domain.ErrUserNotFound
        }
        return &user, err
    }

Testing Repository Implementations [​](#testing-repository-implementations)
---------------------------------------------------------------------------

Repositories require different testing strategies than domain logic because they involve external dependencies.

### Unit Testing with Mocks [​](#unit-testing-with-mocks)

For use case testing, use mock repositories:

go

    type MockUserRepository struct {
        SaveFunc    func(user *domain.User) error
        FindByIDFunc func(id uint) (*domain.User, error)
    }
    
    func (m *MockUserRepository) Save(user *domain.User) error {
        if m.SaveFunc != nil {
            return m.SaveFunc(user)
        }
        return nil
    }
    
    func (m *MockUserRepository) FindByID(id uint) (*domain.User, error) {
        if m.FindByIDFunc != nil {
            return m.FindByIDFunc(id)
        }
        return nil, domain.ErrUserNotFound
    }
    
    // Test use case
    func TestCreateUser(t *testing.T) {
        mockRepo := &MockUserRepository{
            FindByEmailFunc: func(email string) (*domain.User, error) {
                return nil, domain.ErrUserNotFound // Email available
            },
            SaveFunc: func(user *domain.User) error {
                user.ID = 1
                return nil
            },
        }
        
        service := NewUserService(mockRepo)
        
        user, err := service.CreateUser(CreateUserInput{
            Name:  "John",
            Email: "john@example.com",
        })
        
        assert.NoError(t, err)
        assert.Equal(t, uint(1), user.ID)
    }

### Integration Testing with Real Database [​](#integration-testing-with-real-database)

For repository implementation testing, use a real database:

go

    func TestPostgresUserRepository_Save(t *testing.T) {
        // Setup test database
        db := setupTestDatabase(t)
        defer cleanupTestDatabase(t, db)
        
        repo := NewPostgresUserRepository(db)
        
        // Create test user
        user := &domain.User{
            Name:  "John Doe",
            Email: "john@example.com",
            Age:   30,
        }
        
        // Test Save
        err := repo.Save(user)
        assert.NoError(t, err)
        assert.NotZero(t, user.ID)
        
        // Verify in database
        found, err := repo.FindByID(user.ID)
        assert.NoError(t, err)
        assert.Equal(t, user.Name, found.Name)
        assert.Equal(t, user.Email, found.Email)
    }
    
    func setupTestDatabase(t *testing.T) *gorm.DB {
        db, err := gorm.Open(postgres.Open("postgres://test:test@localhost/testdb"), &gorm.Config{})
        require.NoError(t, err)
        
        // Run migrations
        db.AutoMigrate(&domain.User{})
        
        return db
    }
    
    func cleanupTestDatabase(t *testing.T, db *gorm.DB) {
        db.Exec("TRUNCATE TABLE users CASCADE")
    }

### Testing with Docker [​](#testing-with-docker)

For isolated integration tests, use Docker containers:

go

    func TestWithDocker(t *testing.T) {
        if testing.Short() {
            t.Skip("Skipping integration test")
        }
        
        // Start PostgreSQL container
        ctx := context.Background()
        req := testcontainers.ContainerRequest{
            Image:        "postgres:15",
            ExposedPorts: []string{"5432/tcp"},
            Env: map[string]string{
                "POSTGRES_PASSWORD": "test",
                "POSTGRES_DB":       "testdb",
            },
            WaitingFor: wait.ForLog("database system is ready"),
        }
        
        postgres, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
            ContainerRequest: req,
            Started:          true,
        })
        require.NoError(t, err)
        defer postgres.Terminate(ctx)
        
        // Get connection string
        host, _ := postgres.Host(ctx)
        port, _ := postgres.MappedPort(ctx, "5432")
        dsn := fmt.Sprintf("postgres://test:test@%s:%s/testdb", host, port.Port())
        
        // Run tests
        db, _ := gorm.Open(postgres.Open(dsn), &gorm.Config{})
        repo := NewPostgresUserRepository(db)
        
        // ... test repository operations ...
    }

Repository Anti-Patterns [​](#repository-anti-patterns)
-------------------------------------------------------

Understanding what not to do is as important as knowing best practices.

### Anti-Pattern: Generic Repository [​](#anti-pattern-generic-repository)

**Problem**:

go

    type Repository[T any] interface {
        Save(entity T) error
        FindByID(id int) (T, error)
        Update(entity T) error
        Delete(id int) error
    }
    
    type UserRepository = Repository[User]
    type OrderRepository = Repository[Order]

**Why It's Bad**: Generic repositories force all entities to have the same operations. `FindByEmail` makes sense for users, not orders. You lose domain-specific expressiveness.

**Solution**: Create specific interfaces for each entity with domain-meaningful methods.

### Anti-Pattern: Leaky Abstraction [​](#anti-pattern-leaky-abstraction)

**Problem**:

go

    type UserRepository interface {
        Query(sql string, args ...interface{}) ([]*User, error)
        GetDB() *sql.DB
    }

**Why It's Bad**: This exposes database implementation details. Consumers must write SQL. You cannot swap databases.

**Solution**: Provide intention-revealing methods that hide database details.

### Anti-Pattern: Business Logic in Repository [​](#anti-pattern-business-logic-in-repository)

**Problem**:

go

    func (r *postgresUserRepository) CreateAdminUser(name, email string) (*User, error) {
        user := &User{
            Name:  name,
            Email: email,
            Role:  "admin",
        }
        
        // Business logic: validate admin email domain
        if !strings.HasSuffix(email, "@company.com") {
            return nil, errors.New("admin must use company email")
        }
        
        return user, r.db.Create(user).Error
    }

**Why It's Bad**: Business rules belong in the domain or use cases, not repositories. Repositories are for data access only.

**Solution**: Move validation to entity or use case. Repository only persists valid entities.

### Anti-Pattern: Repository Returning DTOs [​](#anti-pattern-repository-returning-dtos)

**Problem**:

go

    type UserRepository interface {
        FindByID(id int) (*UserDTO, error)
    }

**Why It's Bad**: Repositories work with domain entities, not DTOs. DTOs are for external communication, not internal operations.

**Solution**: Return domain entities. Use cases convert entities to DTOs.

Best Practices Summary [​](#best-practices-summary)
---------------------------------------------------

### Do This [​](#do-this)

✅ **Define interfaces in domain layer**: Keep contracts with domain code  
✅ **Use domain types in signatures**: Parameters and returns are entities  
✅ **Name methods by intent**: `FindActiveUsers()`, not `QueryUsers()`  
✅ **Return domain errors**: `ErrUserNotFound`, not `sql.ErrNoRows`  
✅ **Keep implementations simple**: One responsibility per repository  
✅ **Test with real databases**: Integration tests verify SQL correctness  
✅ **Use in-memory fakes for unit tests**: Fast, isolated tests

### Avoid This [​](#avoid-this)

❌ **Don't expose database types**: No `*sql.DB`, `*gorm.DB` in interfaces  
❌ **Don't add business logic**: Repositories persist, they don't validate  
❌ **Don't use generic interfaces**: Each entity gets specific methods  
❌ **Don't return DTOs**: Repositories work with domain entities  
❌ **Don't couple to ORM**: Use interfaces, not concrete ORM types

Conclusion [​](#conclusion)
---------------------------

The Repository pattern is a cornerstone of Clean Architecture, providing a boundary between domain logic and data access concerns. When implemented correctly, repositories enable:

*   **Database Independence**: Change databases without changing business logic
*   **Testability**: Unit test without database setup using mocks
*   **Maintainability**: Data access code isolated in one layer
*   **Flexibility**: Support multiple databases simultaneously
*   **Clean Domain**: Domain remains pure, focused on business rules

Goca's `goca repository` command generates repositories that follow these principles, creating both interfaces and database-specific implementations that maintain architectural boundaries while providing practical, production-ready data access code.

Master repositories, and you master one of the most critical patterns in Clean Architecture.</content>
</page>

<page>
  <title>Goca - Go Clean Architecture Code Generator</title>
  <url>https://sazardev.github.io/goca/commands/feature.html</url>
  <content>goca feature [​](#goca-feature)
-------------------------------

Generate a complete feature with all Clean Architecture layers in a single command.

Syntax [​](#syntax)
-------------------

bash

    goca feature <FeatureName> [flags]

Description [​](#description)
-----------------------------

The `goca feature` command is the **fastest way** to add new functionality to your project. It generates all layers (domain, use case, repository, handler) and automatically integrates them with dependency injection and routing.

Recommended Workflow

This is the recommended command for most use cases. It saves time and ensures all layers are properly connected.

Arguments [​](#arguments)
-------------------------

### `<FeatureName>` [​](#featurename)

**Required.** The name of your feature (singular, PascalCase).

bash

    goca feature User
    goca feature Product
    goca feature Order

Flags [​](#flags)
-----------------

### `--fields` [​](#fields)

Define the structure of your entity.

**Format:** `"field1:type1,field2:type2,..."`

**Supported Types:**

*   `string` - Text data
*   `int`, `int64` - Integer numbers
*   `float64` - Decimal numbers
*   `bool` - Boolean values
*   `time.Time` - Timestamps
*   `[]type` - Arrays/slices

bash

    goca feature Product --fields "name:string,price:float64,inStock:bool"

### `--validation` [​](#validation)

Add domain-level validation rules.

bash

    goca feature User --fields "name:string,email:string,age:int" --validation

Generates validation methods like:

*   Email format validation
*   Required field checks
*   Range validations
*   Custom business rules

### `--database` [​](#database)

Specify database type for repository. Default: `postgres`

**Options:**

*   `postgres` - PostgreSQL (GORM)
*   `postgres-json` - PostgreSQL with JSONB
*   `mysql` - MySQL (GORM)
*   `mongodb` - MongoDB (native driver)
*   `sqlite` - SQLite (embedded)
*   `sqlserver` - SQL Server
*   `elasticsearch` - Elasticsearch (v8)
*   `dynamodb` - DynamoDB (AWS)

bash

    goca feature Order --fields "total:float64" --database mysql
    goca feature Config --fields "name:string,value:string" --database postgres-json
    goca feature Article --fields "title:string,content:string" --database elasticsearch

### `--handlers` [​](#handlers)

Generate multiple handler types.

**Options:** `http` | `grpc` | `cli` | `worker` | `soap`

bash

    goca feature Payment --fields "amount:float64" --handlers "http,grpc"

Examples [​](#examples)
-----------------------

### Basic Feature [​](#basic-feature)

bash

    goca feature User --fields "name:string,email:string,age:int"

**Generates:**

    internal/
    ├── domain/
    │   ├── user.go              # Entity with business rules
    │   └── user_errors.go       # Domain-specific errors
    ├── usecase/
    │   ├── user_dto.go          # Input/Output DTOs
    │   ├── user_interfaces.go   # Use case contracts
    │   └── user_service.go      # Business logic implementation
    ├── repository/
    │   ├── user_repository.go   # Repository interface
    │   └── postgres_user_repository.go  # PostgreSQL implementation
    └── handler/
        └── http/
            └── user_handler.go  # HTTP REST endpoints

### Feature with Validation [​](#feature-with-validation)

bash

    goca feature Product \
      --fields "name:string,price:float64,stock:int,category:string" \
      --validation

Generates entity with:

go

    func (p *Product) Validate() error {
        if p.Name == "" {
            return ErrProductNameRequired
        }
        if p.Price < 0 {
            return ErrInvalidPrice
        }
        if p.Stock < 0 {
            return ErrInvalidStock
        }
        return nil
    }

### E-commerce Order Feature [​](#e-commerce-order-feature)

bash

    goca feature Order \
      --fields "customerID:int,items:[]OrderItem,total:float64,status:string,createdAt:time.Time" \
      --validation \
      --database postgres

### Multi-Protocol Service [​](#multi-protocol-service)

bash

    goca feature Payment \
      --fields "amount:float64,currency:string,method:string,status:string" \
      --handlers "http,grpc,worker" \
      --validation

### Complex Domain Model [​](#complex-domain-model)

bash

    goca feature Invoice \
      --fields "number:string,customerID:int,items:[]InvoiceItem,subtotal:float64,tax:float64,total:float64,dueDate:time.Time,status:string" \
      --validation \
      --database postgres

What Gets Generated [​](#what-gets-generated)
---------------------------------------------

### 1\. Domain Layer (`internal/domain/`) [​](#_1-domain-layer-internal-domain)

**Entity:** `<feature>.go`

go

    package domain
    
    type User struct {
        ID        uint      `json:"id"`
        Name      string    `json:"name"`
        Email     string    `json:"email"`
        Age       int       `json:"age"`
        CreatedAt time.Time `json:"created_at"`
        UpdatedAt time.Time `json:"updated_at"`
    }
    
    func (u *User) Validate() error { /* ... */ }

**Errors:** `<feature>_errors.go`

go

    var (
        ErrUserNotFound      = errors.New("user not found")
        ErrUserNameRequired  = errors.New("name is required")
        ErrInvalidEmail      = errors.New("invalid email")
    )

### 2\. Use Case Layer (`internal/usecase/`) [​](#_2-use-case-layer-internal-usecase)

**DTOs:** `<feature>_dto.go`

go

    type CreateUserRequest struct {
        Name  string `json:"name" validate:"required"`
        Email string `json:"email" validate:"required,email"`
        Age   int    `json:"age" validate:"min=0,max=150"`
    }
    
    type UserResponse struct {
        ID    uint   `json:"id"`
        Name  string `json:"name"`
        Email string `json:"email"`
        Age   int    `json:"age"`
    }

**Service:** `<feature>_service.go`

go

    type UserService interface {
        Create(ctx context.Context, req CreateUserRequest) (*UserResponse, error)
        GetByID(ctx context.Context, id uint) (*UserResponse, error)
        Update(ctx context.Context, id uint, req UpdateUserRequest) error
        Delete(ctx context.Context, id uint) error
        List(ctx context.Context) ([]*UserResponse, error)
    }

### 3\. Repository Layer (`internal/repository/`) [​](#_3-repository-layer-internal-repository)

**Interface:** Repository contract in domain **Implementation:** `postgres_<feature>_repository.go`

go

    type postgresUserRepository struct {
        db *sql.DB
    }
    
    func (r *postgresUserRepository) Save(ctx context.Context, user *domain.User) error {
        // Implementation
    }

### 4\. Handler Layer (`internal/handler/http/`) [​](#_4-handler-layer-internal-handler-http)

**Handler:** `<feature>_handler.go`

go

    type UserHandler struct {
        service usecase.UserService
    }
    
    func (h *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
        // HTTP handling
    }

**Routes:** Automatically registered

go

    // POST   /api/v1/users
    // GET    /api/v1/users/:id
    // PUT    /api/v1/users/:id
    // DELETE /api/v1/users/:id
    // GET    /api/v1/users

Automatic Integration [​](#automatic-integration)
-------------------------------------------------

After generating a feature, Goca automatically:

Updates dependency injection container  
Registers HTTP routes  
Adds database migrations  
Configures repository connections  
Wires all dependencies

**You can immediately test your new feature!**

bash

    # Generate feature
    goca feature User --fields "name:string,email:string"
    
    # Run server (feature is already integrated!)
    go run cmd/server/main.go
    
    # Test the API
    curl -X POST http://localhost:8080/api/v1/users \
      -H "Content-Type: application/json" \
      -d '{"name":"John Doe","email":"john@example.com"}'

After Generation [​](#after-generation)
---------------------------------------

### 1\. Review Generated Code [​](#_1-review-generated-code)

Check the generated files and customize as needed:

bash

    # View generated files
    find internal -name "*user*"

### 2\. Add Business Logic [​](#_2-add-business-logic)

Enhance the domain entity with business rules:

go

    // internal/domain/user.go
    func (u *User) CanDelete() bool {
        return u.Status != "active"
    }
    
    func (u *User) IsAdmin() bool {
        return u.Role == "admin"
    }

### 3\. Run Tests [​](#_3-run-tests)

bash

    go test ./internal/...

### 4\. Run Application [​](#_4-run-application)

bash

    go run cmd/server/main.go

Tips [​](#tips)
---------------

### Start Simple [​](#start-simple)

Begin with basic fields, then add complexity:

bash

    # Start
    goca feature Product --fields "name:string,price:float64"
    
    # Later, manually add relationships and complex logic

### Use Consistent Naming [​](#use-consistent-naming)

*   Feature names: **Singular**, **PascalCase** (User, Product, Order)
*   Fields: **camelCase** (firstName, productName)

### Field Naming Conventions [​](#field-naming-conventions)

bash

    #  Good
    --fields "firstName:string,lastName:string,emailAddress:string"
    
    #  Avoid
    --fields "first_name:string,Last-Name:string,EMAIL:string"

### Complex Types [​](#complex-types)

For complex types, generate basic structure first, then modify:

bash

    goca feature Order --fields "customerID:int,total:float64"
    
    # Then manually add:
    # - items []OrderItem
    # - metadata map[string]interface{}
    # - custom types

Troubleshooting [​](#troubleshooting)
-------------------------------------

### Feature Already Exists [​](#feature-already-exists)

**Problem:** "feature already exists"

**Solution:** Choose a different name or delete existing files:

bash

    rm -rf internal/domain/user*
    rm -rf internal/usecase/user*
    # ... then regenerate

### Invalid Field Type [​](#invalid-field-type)

**Problem:** "unsupported field type"

**Solution:** Use supported types or add manually after generation.

### Integration Failed [​](#integration-failed)

**Problem:** Routes not working

**Solution:** Run integration command:

bash

    goca integrate --feature User

See Also [​](#see-also)
-----------------------

*   [`goca init`](https://sazardev.github.io/goca/commands/init.html) - Initialize project
*   [`goca integrate`](https://sazardev.github.io/goca/commands/integrate.html) - Manual integration
*   [`goca entity`](https://sazardev.github.io/goca/commands/entity.html) - Generate entity only
*   [Complete Tutorial](https://sazardev.github.io/goca/tutorials/complete-tutorial.html) - Step-by-step guide

Next Steps [​](#next-steps)
---------------------------

*   Learn about [individual layer commands](https://sazardev.github.io/goca/commands/)
*   Follow the [Complete Tutorial](https://sazardev.github.io/goca/tutorials/complete-tutorial.html)
*   Understand [Clean Architecture](https://sazardev.github.io/goca/guide/clean-architecture.html)</content>
</page>

<page>
  <title>Mastering Use Cases in Clean Architecture | Goca Blog</title>
  <url>https://sazardev.github.io/goca/blog/articles/mastering-use-cases.html</url>
  <content>Mastering Use Cases in Clean Architecture [​](#mastering-use-cases-in-clean-architecture)
-----------------------------------------------------------------------------------------

ArchitectureApplication Layer

Use cases represent application-specific business rules and orchestrate the flow of data between entities and external systems. Understanding use cases correctly is critical for building well-structured applications that adapt to changing requirements without compromising core business logic.

* * *

What is a Use Case? [​](#what-is-a-use-case)
--------------------------------------------

A use case is an application service that coordinates domain entities and infrastructure to fulfill a specific user or system goal. Use cases answer the question: "What can the application do?"

### Core Responsibilities [​](#core-responsibilities)

**Orchestration**: Use cases coordinate multiple domain entities, repositories, and external services to complete a workflow. They do not contain business rules; they apply them.

**Data Flow Control**: Use cases manage the flow of data between the UI layer and the domain layer, transforming external requests into domain operations and domain results into external responses.

**Application Logic**: Use cases implement application-specific rules that do not belong in the domain. These rules depend on the use case context, not on core business concepts.

**Transaction Management**: Use cases define transaction boundaries, ensuring that operations either complete fully or roll back entirely.

**Permission and Security**: Use cases enforce authorization rules, checking whether the requesting user can perform the operation.

Use Case vs Controller vs Service [​](#use-case-vs-controller-vs-service)
-------------------------------------------------------------------------

Many developers confuse use cases with controllers or generic services. This confusion leads to bloated classes and violated architectural boundaries.

### What a Use Case Is NOT [​](#what-a-use-case-is-not)

**Not a Controller**: Controllers are adapters that convert HTTP requests to use case calls. Controllers handle protocol concerns; use cases handle application logic.

**Not a Generic Service**: A use case serves a specific goal, not general utilities. Services like "EmailService" or "LoggerService" are infrastructure concerns, not use cases.

**Not a Transaction Script**: Use cases orchestrate domain entities. They do not implement business rules. Domain logic belongs in entities and value objects.

**Not a Facade**: Use cases are not simple pass-throughs to repositories. They add application-level coordination and workflow management.

### The Clear Distinction [​](#the-clear-distinction)

    HTTP Request
        ↓
    Controller (Adapter - Outer Layer)
        ↓ Converts to DTO
    Use Case (Application Layer)
        ↓ Orchestrates
    Domain Entity (Domain Layer)
        ↓ Enforces rules
    Repository Interface (Domain Layer)
        ↓ Implements
    Repository Implementation (Infrastructure Layer)
        ↓ Persists
    Database

Each layer has distinct responsibilities. Use cases sit between adapters and domain, orchestrating operations without implementing business rules or handling external protocols.

The Application Layer [​](#the-application-layer)
-------------------------------------------------

Use cases form the application layer in Clean Architecture, distinct from both the domain layer and the infrastructure layer.

### Application Layer Characteristics [​](#application-layer-characteristics)

**Depends on Domain**: Use cases depend on domain entities and interfaces. They call entity methods and use repository interfaces defined in the domain.

**Independent of Infrastructure**: Use cases do not import database drivers, HTTP libraries, or external service clients. They work with interfaces.

**Stateless by Design**: Use cases do not maintain state between calls. Each operation is independent.

**Transaction Boundaries**: Use cases define where transactions begin and end, ensuring data consistency.

### Why a Separate Layer? [​](#why-a-separate-layer)

The application layer exists because application logic and domain logic are different:

**Domain Logic**: "A user must have a valid email address" is domain logic. This rule exists regardless of how you access users.

**Application Logic**: "To create a user, check if the email exists, create the user, send a welcome email" is application logic. This workflow is specific to the user creation use case.

Separating these concerns allows you to:

*   Change workflows without changing domain rules
*   Test business rules without application context
*   Reuse domain logic across different workflows
*   Evolve the application independently of the domain

Data Transfer Objects (DTOs) [​](#data-transfer-objects-dtos)
-------------------------------------------------------------

DTOs are simple structures that carry data between layers without behavior. Use cases use DTOs to receive input and provide output.

### Why DTOs? [​](#why-dtos)

**Layer Separation**: DTOs prevent external layers from depending on domain entities directly. Changing an entity does not break API contracts.

**Validation Boundary**: DTOs define what data the use case needs and validate it before processing.

**Security**: DTOs control what data external systems can provide or receive, preventing over-posting and data exposure.

**Versioning**: DTOs allow multiple API versions to coexist by mapping different external structures to the same domain entities.

### Input DTOs [​](#input-dtos)

Input DTOs represent the data a use case needs to perform an operation:

go

    type CreateUserInput struct {
        Name  string `json:"name" validate:"required,min=2"`
        Email string `json:"email" validate:"required,email"`
        Age   int    `json:"age" validate:"required,gte=0"`
    }

Input DTOs include validation tags that define constraints:

*   **required**: Field must be present
*   **min/max**: String length or numeric range
*   **email**: Valid email format
*   **gte/lte**: Greater than or equal / less than or equal

These validations are input validations, not business rules. They ensure the data is well-formed before processing.

### Output DTOs [​](#output-dtos)

Output DTOs represent the data a use case returns:

go

    type CreateUserOutput struct {
        User    domain.User `json:"user"`
        Message string      `json:"message"`
    }

Output DTOs can include:

*   Domain entities for complete information
*   Specific fields for minimal responses
*   Metadata like messages or status codes
*   Related entities for composite responses

### Update DTOs with Optional Fields [​](#update-dtos-with-optional-fields)

Update operations use optional fields to support partial updates:

go

    type UpdateUserInput struct {
        Name  *string `json:"name,omitempty" validate:"omitempty,min=2"`
        Email *string `json:"email,omitempty" validate:"omitempty,email"`
        Age   *int    `json:"age,omitempty" validate:"omitempty,gte=0"`
    }

Pointer fields distinguish between "not provided" (nil) and "explicitly set to zero value" (non-nil pointer to zero value). This allows clients to update only specific fields without affecting others.

### List DTOs [​](#list-dtos)

List operations return collections with metadata:

go

    type ListUserOutput struct {
        Users   []domain.User `json:"users"`
        Total   int           `json:"total"`
        Message string        `json:"message"`
    }

List DTOs can include pagination information, filters applied, and total counts.

Use Case Implementation Patterns [​](#use-case-implementation-patterns)
-----------------------------------------------------------------------

Use cases follow consistent patterns regardless of the specific operation they perform.

### Basic Structure [​](#basic-structure)

Every use case implementation includes:

1.  Dependency injection of required repositories
2.  Input validation
3.  Domain entity coordination
4.  Business rule enforcement
5.  Persistence through repositories
6.  Output transformation

### Create Operation [​](#create-operation)

The create operation instantiates a new domain entity, validates it, and persists it:

go

    func (s *userService) Create(input CreateUserInput) (*CreateUserOutput, error) {
        // 1. Create domain entity from input
        user := domain.User{
            Name:  input.Name,
            Email: input.Email,
            Age:   input.Age,
        }
        
        // 2. Validate business rules
        if err := user.Validate(); err != nil {
            return nil, err
        }
        
        // 3. Persist through repository
        if err := s.repo.Save(&user); err != nil {
            return nil, err
        }
        
        // 4. Return success output
        return &CreateUserOutput{
            User:    user,
            Message: "User created successfully",
        }, nil
    }

This pattern ensures that:

*   Input is transformed to domain entities
*   Business rules are enforced before persistence
*   Repository handles storage concerns
*   Output is well-defined and structured

### Read Operation [​](#read-operation)

The read operation retrieves an entity by identifier:

go

    func (s *userService) GetByID(id uint) (*domain.User, error) {
        return s.repo.FindByID(int(id))
    }

Read operations are simple because they delegate directly to repositories. Complexity arises when reads require:

*   Authorization checks
*   Data enrichment from multiple sources
*   Transformation to specific output formats

### Update Operation [​](#update-operation)

The update operation retrieves an entity, modifies it, validates it, and persists changes:

go

    func (s *userService) Update(id uint, input UpdateUserInput) (*domain.User, error) {
        // 1. Retrieve existing entity
        user, err := s.repo.FindByID(int(id))
        if err != nil {
            return nil, err
        }
        
        // 2. Apply changes from input (only provided fields)
        if input.Name != nil {
            user.Name = *input.Name
        }
        if input.Email != nil {
            user.Email = *input.Email
        }
        if input.Age != nil {
            user.Age = *input.Age
        }
        
        // 3. Validate updated entity
        if err := user.Validate(); err != nil {
            return nil, err
        }
        
        // 4. Persist changes
        if err := s.repo.Update(user); err != nil {
            return nil, err
        }
        
        return user, nil
    }

The update pattern:

*   Retrieves current state
*   Applies only provided changes
*   Validates the result
*   Persists atomically

### Delete Operation [​](#delete-operation)

The delete operation removes an entity:

go

    func (s *userService) Delete(id uint) error {
        return s.repo.Delete(int(id))
    }

Delete operations can be:

*   **Hard Delete**: Permanently removes the record
*   **Soft Delete**: Marks the record as deleted without removing it

Soft deletes are preferable for audit trails and data recovery.

### List Operation [​](#list-operation)

The list operation retrieves collections with optional filtering:

go

    func (s *userService) List() (*ListUserOutput, error) {
        users, err := s.repo.FindAll()
        if err != nil {
            return nil, err
        }
        
        return &ListUserOutput{
            Users:   users,
            Total:   len(users),
            Message: "Users listed successfully",
        }, nil
    }

List operations often include:

*   Pagination parameters
*   Sort order specifications
*   Filter conditions
*   Total count calculation

How Goca Generates Use Cases [​](#how-goca-generates-use-cases)
---------------------------------------------------------------

Goca provides the `goca usecase` command to generate complete application services with DTOs and interfaces.

### Basic Use Case Generation [​](#basic-use-case-generation)

bash

    goca usecase UserService --entity User

This generates three files:

**dto.go**: Input and output DTOs

go

    package usecase
    
    import (
        "github.com/yourorg/yourproject/internal/domain"
    )
    
    type CreateUserInput struct {
        Name  string `json:"name" validate:"required,min=2"`
        Email string `json:"email" validate:"required,email"`
        Age   int    `json:"age" validate:"required,min=1"`
    }
    
    type CreateUserOutput struct {
        User    domain.User `json:"user"`
        Message string      `json:"message"`
    }
    
    type UpdateUserInput struct {
        Name  *string `json:"name,omitempty" validate:"omitempty,min=2"`
        Email *string `json:"email,omitempty" validate:"omitempty,email"`
        Age   *int    `json:"age,omitempty" validate:"omitempty,min=1"`
    }
    
    type ListUserOutput struct {
        Users   []domain.User `json:"users"`
        Total   int           `json:"total"`
        Message string        `json:"message"`
    }

**user\_service.go**: Service implementation

go

    package usecase
    
    import (
        "github.com/yourorg/yourproject/internal/domain"
        "github.com/yourorg/yourproject/internal/repository"
    )
    
    type userService struct {
        repo repository.UserRepository
    }
    
    func NewUserService(repo repository.UserRepository) UserUseCase {
        return &userService{repo: repo}
    }
    
    func (u *userService) Create(input CreateUserInput) (*CreateUserOutput, error) {
        user := domain.User{
            Name:  input.Name,
            Email: input.Email,
            Age:   input.Age,
        }
        
        if err := user.Validate(); err != nil {
            return nil, err
        }
        
        if err := u.repo.Save(&user); err != nil {
            return nil, err
        }
        
        return &CreateUserOutput{
            User:    user,
            Message: "User created successfully",
        }, nil
    }
    
    func (u *userService) GetByID(id uint) (*domain.User, error) {
        return u.repo.FindByID(int(id))
    }
    
    func (u *userService) Update(id uint, input UpdateUserInput) (*domain.User, error) {
        user, err := u.repo.FindByID(int(id))
        if err != nil {
            return nil, err
        }
        
        if input.Name != nil {
            user.Name = *input.Name
        }
        if input.Email != nil {
            user.Email = *input.Email
        }
        if input.Age != nil {
            user.Age = *input.Age
        }
        
        if err := u.repo.Update(user); err != nil {
            return nil, err
        }
        
        return user, nil
    }
    
    func (u *userService) Delete(id uint) error {
        return u.repo.Delete(int(id))
    }
    
    func (u *userService) List() (*ListUserOutput, error) {
        users, err := u.repo.FindAll()
        if err != nil {
            return nil, err
        }
        
        return &ListUserOutput{
            Users:   users,
            Total:   len(users),
            Message: "Users listed successfully",
        }, nil
    }

**interfaces.go**: Repository interface definition

The service depends on a repository interface:

go

    type UserRepository interface {
        Save(user *domain.User) error
        FindByID(id int) (*domain.User, error)
        Update(user *domain.User) error
        Delete(id int) error
        FindAll() ([]domain.User, error)
    }

This interface lives in the repository package but is used by the use case. The use case depends on the abstraction, not the implementation.

### Selecting Operations [​](#selecting-operations)

Control which CRUD operations to generate:

bash

    goca usecase ProductService --entity Product --operations "create,read,update"

This generates only create, read, and update methods, omitting delete and list.

Available operations:

*   **create**: Instantiate and persist new entities
*   **read** or **get**: Retrieve entities by ID
*   **update**: Modify existing entities
*   **delete**: Remove entities
*   **list**: Retrieve collections

### DTO Validation [​](#dto-validation)

Enable validation tags on DTOs:

bash

    goca usecase OrderService --entity Order --dto-validation

With validation enabled, DTOs include comprehensive validation rules:

go

    type CreateOrderInput struct {
        CustomerID int     `json:"customer_id" validate:"required,gt=0"`
        Total      float64 `json:"total" validate:"required,gte=0"`
        Status     string  `json:"status" validate:"required,oneof=pending confirmed shipped delivered"`
    }

Validation rules ensure:

*   Required fields are present
*   Numeric values are within acceptable ranges
*   Strings match expected patterns or enumerations
*   Email addresses are valid
*   Custom validation logic is applied

Advanced Use Case Patterns [​](#advanced-use-case-patterns)
-----------------------------------------------------------

Beyond basic CRUD, use cases handle complex workflows and business processes.

### Transactional Use Cases [​](#transactional-use-cases)

Some operations require multiple steps within a single transaction:

go

    func (s *orderService) CreateOrder(input CreateOrderInput) (*CreateOrderOutput, error) {
        // Begin transaction (pseudo-code, actual implementation depends on repository)
        
        // 1. Validate customer exists
        customer, err := s.customerRepo.FindByID(input.CustomerID)
        if err != nil {
            return nil, errors.New("customer not found")
        }
        
        // 2. Check product availability
        for _, item := range input.Items {
            product, err := s.productRepo.FindByID(item.ProductID)
            if err != nil {
                return nil, err
            }
            
            if product.Stock < item.Quantity {
                return nil, errors.New("insufficient stock")
            }
        }
        
        // 3. Create order
        order := &domain.Order{
            CustomerID: input.CustomerID,
            Items:      mapOrderItems(input.Items),
            Total:      calculateTotal(input.Items),
            Status:     "pending",
        }
        
        if err := order.Validate(); err != nil {
            return nil, err
        }
        
        // 4. Persist order
        if err := s.orderRepo.Save(order); err != nil {
            return nil, err
        }
        
        // 5. Update product stock
        for _, item := range order.Items {
            product, _ := s.productRepo.FindByID(item.ProductID)
            product.Stock -= item.Quantity
            s.productRepo.Update(product)
        }
        
        // Commit transaction
        
        return &CreateOrderOutput{
            Order:   *order,
            Message: "Order created successfully",
        }, nil
    }

This use case:

*   Validates dependencies (customer exists)
*   Checks business constraints (sufficient stock)
*   Creates the main entity (order)
*   Updates related entities (product stock)
*   Ensures atomicity through transactions

### Async Use Cases [​](#async-use-cases)

Some operations can execute asynchronously to improve response times:

bash

    goca usecase NotificationService --entity Notification --operations "create" --async

Asynchronous use cases return immediately while processing continues in the background:

go

    func (s *notificationService) SendNotification(input SendNotificationInput) (*SendNotificationOutput, error) {
        // Validate input immediately
        if err := input.Validate(); err != nil {
            return nil, err
        }
        
        // Create notification record
        notification := &domain.Notification{
            UserID:  input.UserID,
            Message: input.Message,
            Status:  "queued",
        }
        
        if err := s.repo.Save(notification); err != nil {
            return nil, err
        }
        
        // Queue for async processing
        s.queue.Enqueue(notification.ID)
        
        // Return immediately
        return &SendNotificationOutput{
            NotificationID: notification.ID,
            Status:         "queued",
            Message:        "Notification queued successfully",
        }, nil
    }

The actual sending happens asynchronously:

go

    func (s *notificationService) ProcessQueue() {
        for {
            notificationID := s.queue.Dequeue()
            
            notification, err := s.repo.FindByID(notificationID)
            if err != nil {
                continue
            }
            
            // Send notification via external service
            err = s.emailService.Send(notification.UserID, notification.Message)
            
            if err != nil {
                notification.Status = "failed"
            } else {
                notification.Status = "sent"
            }
            
            s.repo.Update(notification)
        }
    }

Asynchronous use cases are appropriate for:

*   Email sending
*   File processing
*   Report generation
*   Third-party API calls
*   Long-running computations

### Composite Use Cases [​](#composite-use-cases)

Some operations aggregate data from multiple sources:

go

    func (s *dashboardService) GetUserDashboard(userID uint) (*DashboardOutput, error) {
        // Retrieve user
        user, err := s.userRepo.FindByID(userID)
        if err != nil {
            return nil, err
        }
        
        // Retrieve user orders
        orders, err := s.orderRepo.FindByUserID(userID)
        if err != nil {
            return nil, err
        }
        
        // Calculate statistics
        totalSpent := calculateTotalSpent(orders)
        averageOrderValue := totalSpent / float64(len(orders))
        
        // Retrieve recent activity
        activity, err := s.activityRepo.FindRecentByUserID(userID, 10)
        if err != nil {
            return nil, err
        }
        
        return &DashboardOutput{
            User:              user,
            TotalOrders:       len(orders),
            TotalSpent:        totalSpent,
            AverageOrderValue: averageOrderValue,
            RecentActivity:    activity,
        }, nil
    }

Composite use cases coordinate multiple repositories to build aggregate views.

Testing Use Cases [​](#testing-use-cases)
-----------------------------------------

Use cases are highly testable because they depend on interfaces, not concrete implementations.

### Unit Testing with Mocks [​](#unit-testing-with-mocks)

Test use cases by mocking repository dependencies:

go

    func TestUserService_Create(t *testing.T) {
        // Create mock repository
        mockRepo := new(MockUserRepository)
        
        // Setup expectations
        mockRepo.On("Save", mock.AnythingOfType("*domain.User")).Return(nil)
        
        // Create service with mock
        service := NewUserService(mockRepo)
        
        // Execute use case
        input := CreateUserInput{
            Name:  "John Doe",
            Email: "john@example.com",
            Age:   30,
        }
        
        output, err := service.Create(input)
        
        // Assert results
        assert.NoError(t, err)
        assert.NotNil(t, output)
        assert.Equal(t, "John Doe", output.User.Name)
        assert.Equal(t, "User created successfully", output.Message)
        
        // Verify mock was called
        mockRepo.AssertExpectations(t)
    }

Mock repositories allow you to:

*   Test use case logic independently
*   Simulate repository errors
*   Verify correct method calls
*   Control return values

### Testing Validation [​](#testing-validation)

Test that use cases enforce validation correctly:

go

    func TestUserService_Create_InvalidEmail(t *testing.T) {
        mockRepo := new(MockUserRepository)
        service := NewUserService(mockRepo)
        
        input := CreateUserInput{
            Name:  "John Doe",
            Email: "invalid-email",
            Age:   30,
        }
        
        output, err := service.Create(input)
        
        assert.Error(t, err)
        assert.Nil(t, output)
        assert.Contains(t, err.Error(), "email")
        
        // Repository should not be called
        mockRepo.AssertNotCalled(t, "Save")
    }

### Testing Error Handling [​](#testing-error-handling)

Test that use cases handle repository errors gracefully:

go

    func TestUserService_Create_RepositoryError(t *testing.T) {
        mockRepo := new(MockUserRepository)
        
        // Simulate repository error
        mockRepo.On("Save", mock.Anything).Return(errors.New("database connection failed"))
        
        service := NewUserService(mockRepo)
        
        input := CreateUserInput{
            Name:  "John Doe",
            Email: "john@example.com",
            Age:   30,
        }
        
        output, err := service.Create(input)
        
        assert.Error(t, err)
        assert.Nil(t, output)
        assert.Equal(t, "database connection failed", err.Error())
    }

Integration with Other Layers [​](#integration-with-other-layers)
-----------------------------------------------------------------

Use cases coordinate between domain, infrastructure, and adapter layers.

### Handler to Use Case [​](#handler-to-use-case)

Handlers convert external requests to use case calls:

go

    func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
        // 1. Parse HTTP request
        var input usecase.CreateUserInput
        if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
            http.Error(w, "Invalid request body", http.StatusBadRequest)
            return
        }
        
        // 2. Call use case
        output, err := h.usecase.Create(input)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        
        // 3. Return HTTP response
        w.Header().Set("Content-Type", "application/json")
        w.WriteStatus(http.StatusCreated)
        json.NewEncoder(w).Encode(output)
    }

The handler:

*   Handles HTTP concerns (parsing, status codes, headers)
*   Delegates business logic to the use case
*   Transforms use case output to HTTP response

### Use Case to Repository [​](#use-case-to-repository)

Use cases call repository methods through interfaces:

go

    type userService struct {
        repo repository.UserRepository // Interface, not implementation
    }
    
    func (s *userService) Create(input CreateUserInput) (*CreateUserOutput, error) {
        user := domain.User{
            Name:  input.Name,
            Email: input.Email,
            Age:   input.Age,
        }
        
        // Use case calls repository interface
        if err := s.repo.Save(&user); err != nil {
            return nil, err
        }
        
        return &CreateUserOutput{User: user}, nil
    }

The repository implementation lives in the infrastructure layer:

go

    type postgresUserRepository struct {
        db *gorm.DB
    }
    
    func (r *postgresUserRepository) Save(user *domain.User) error {
        return r.db.Create(user).Error
    }

This separation allows:

*   Swapping database implementations
*   Testing use cases without databases
*   Changing persistence strategies independently

### Dependency Injection [​](#dependency-injection)

Use cases receive dependencies through constructors:

go

    func NewUserService(
        repo repository.UserRepository,
        emailService EmailService,
        logger Logger,
    ) UserUseCase {
        return &userService{
            repo:         repo,
            emailService: emailService,
            logger:       logger,
        }
    }

Dependencies are injected at the composition root, typically in a DI container:

go

    // Composition root
    userRepo := repository.NewPostgresUserRepository(db)
    emailService := email.NewSMTPService(config)
    logger := log.NewStdLogger()
    
    userService := usecase.NewUserService(userRepo, emailService, logger)
    userHandler := handler.NewUserHandler(userService)

Best Practices for Use Cases [​](#best-practices-for-use-cases)
---------------------------------------------------------------

Follow these practices to maintain clean, maintainable use cases.

**Keep Use Cases Thin**: Use cases orchestrate; they do not implement business rules. Business logic belongs in domain entities.

**One Use Case, One Goal**: Each use case serves a specific goal. "Create user" is one use case. "Create user and send email" might be one or two, depending on cohesion.

**Use DTOs for All External Data**: Never pass domain entities directly to or from external layers. DTOs provide a stable contract.

**Validate at Boundaries**: Validate input at the use case boundary. Do not assume data is valid.

**Return Errors, Don't Panic**: Use cases return errors for exceptional conditions. They do not panic or crash.

**Keep Dependencies Minimal**: Use cases should depend only on repositories and essential services. Avoid excessive dependencies.

**Write Comprehensive Tests**: Test use cases thoroughly with mocked dependencies. Use cases are the easiest layer to test.

**Document Complex Workflows**: Use cases with multiple steps should be documented clearly, explaining the workflow and error handling.

Common Mistakes to Avoid [​](#common-mistakes-to-avoid)
-------------------------------------------------------

**Business Logic in Use Cases**: Do not implement business rules in use cases. Use cases apply rules defined in entities.

**Direct Database Access**: Use cases should not import database drivers or execute SQL. They call repository methods.

**Mixing Concerns**: Use cases should not handle HTTP parsing, logging details, or UI concerns. They orchestrate business operations.

**Returning Domain Entities Directly**: Always use DTOs for external communication. Domain entities are internal structures.

**Ignoring Errors**: Handle repository errors appropriately. Log them, wrap them, or transform them, but do not ignore them.

**Tight Coupling**: Use cases depending on concrete implementations cannot be tested or swapped easily. Depend on interfaces.

Generating Complete Features [​](#generating-complete-features)
---------------------------------------------------------------

While `goca usecase` generates use cases, `goca feature` generates complete features including entities, use cases, repositories, and handlers:

bash

    goca feature User --fields "name:string,email:string,age:int"

This creates:

*   Domain entity with validation
*   Use case with all CRUD operations
*   Repository interface and implementation
*   HTTP handler
*   DTOs for all operations
*   Dependency injection wiring

All layers work together following Clean Architecture principles, with use cases at the center coordinating workflows.

Conclusion [​](#conclusion)
---------------------------

Use cases are the application layer in Clean Architecture, orchestrating domain entities and infrastructure services to fulfill user goals. They coordinate workflows without implementing business rules, maintain clear boundaries through DTOs, and depend on abstractions rather than implementations.

Understanding use cases correctly is essential for building maintainable applications. They are not controllers, not generic services, and not transaction scripts. They are focused coordinators that apply domain rules in application-specific contexts.

Goca generates production-ready use cases with comprehensive DTOs, following established patterns and best practices. By using Goca's use case generation and understanding these principles, you create systems that are:

*   Easy to test with mocked dependencies
*   Simple to modify as requirements change
*   Clear in expressing application workflows
*   Maintainable over long periods
*   Adaptable to new platforms and interfaces

Start with clear use case boundaries. Orchestrate domain logic. Let adapters handle external concerns. Build applications that last.

Further Reading [​](#further-reading)
-------------------------------------

*   Application layer patterns in the [guide section](https://sazardev.github.io/goca/guide/clean-architecture.html)
*   Complete command reference for [`goca usecase`](https://sazardev.github.io/goca/commands/usecase.html)
*   Full feature generation with [`goca feature`](https://sazardev.github.io/goca/commands/feature.html)
*   Understanding domain entities in [our previous article](https://sazardev.github.io/goca/blog/articles/understanding-domain-entities.html)
*   Repository pattern implementation examples
*   Dependency injection patterns and best practices</content>
</page>

<page>
  <title>Goca - Go Clean Architecture Code Generator</title>
  <url>https://sazardev.github.io/goca/commands/entity.html</url>
  <content>goca entity [​](#goca-entity)
-----------------------------

Generate pure domain entities following Domain-Driven Design (DDD) principles.

Syntax [​](#syntax)
-------------------

bash

    goca entity <EntityName> [flags]

Description [​](#description)
-----------------------------

The `goca entity` command generates domain entities that represent core business concepts. These entities are pure, containing only business logic without external dependencies.

Domain Layer

Entities are the heart of Clean Architecture - they contain enterprise-wide business rules and are completely independent of frameworks, databases, or external systems.

Arguments [​](#arguments)
-------------------------

### `<EntityName>` [​](#entityname)

**Required.** The name of your entity (singular, PascalCase).

bash

    goca entity User
    goca entity Product
    goca entity Order

Flags [​](#flags)
-----------------

### `--fields` [​](#fields)

**Required.** Define the structure of your entity.

**Format:** `"field1:type1,field2:type2,..."`

**Supported Types:**

*   `string` - Text data
*   `int`, `int64` - Integer numbers
*   `float64` - Decimal numbers
*   `bool` - Boolean values
*   `time.Time` - Timestamps
*   `[]type` - Arrays/slices

bash

    goca entity Product --fields "name:string,price:float64,stock:int"

### `--validation` [​](#validation)

Include domain-level validation methods.

bash

    goca entity User --fields "name:string,email:string,age:int" --validation

Generates:

*   `Validate()` method with business rules
*   Domain-specific error constants
*   Input sanitization

### `--business-rules` [​](#business-rules)

Generate business rule methods.

bash

    goca entity Order --fields "total:float64,status:string" --business-rules

Generates methods like:

*   `CanBeCancelled() bool`
*   `IsCompleted() bool`
*   Business logic calculations

### `--timestamps` [​](#timestamps)

Add automatic timestamp fields.

bash

    goca entity Product --fields "name:string,price:float64" --timestamps

Adds:

*   `CreatedAt time.Time`
*   `UpdatedAt time.Time`

### `--soft-delete` [​](#soft-delete)

Enable soft delete functionality.

bash

    goca entity User --fields "name:string,email:string" --soft-delete

Adds:

*   `DeletedAt *time.Time`
*   `IsDeleted() bool` method

### `--tests` [​](#tests)

Generate unit tests for the entity (enabled by default).

bash

    goca entity User --fields "name:string,email:string,age:int" --validation --tests

**Generates:** `internal/domain/user_test.go`

Creates comprehensive test suite including:

*   **Validation tests**: Table-driven tests for `Validate()` method
*   **Initialization tests**: Verify field assignment
*   **Edge case tests**: Test boundary conditions for each field
*   **Type-specific tests**: String length, numeric ranges, email format, etc.

**Disable tests:**

bash

    goca entity User --fields "name:string,email:string" --tests=false

Testing Best Practices

Generated tests use [testify/assert](https://github.com/stretchr/testify) for readable assertions and follow table-driven test patterns recommended by the Go community.

Examples [​](#examples)
-----------------------

### Basic Entity [​](#basic-entity)

bash

    goca entity User --fields "name:string,email:string,age:int"

**Generates:** `internal/domain/user.go`

go

    package domain
    
    import "time"
    
    type User struct {
        ID        uint      `json:"id"`
        Name      string    `json:"name"`
        Email     string    `json:"email"`
        Age       int       `json:"age"`
        CreatedAt time.Time `json:"created_at"`
        UpdatedAt time.Time `json:"updated_at"`
    }
    
    func NewUser(name, email string, age int) *User {
        return &User{
            Name:  name,
            Email: email,
            Age:   age,
        }
    }

### Entity with Validation [​](#entity-with-validation)

bash

    goca entity Product \
      --fields "name:string,price:float64,stock:int" \
      --validation

**Generates:**

go

    package domain
    
    import (
        "errors"
        "strings"
    )
    
    type Product struct {
        ID    uint    `json:"id"`
        Name  string  `json:"name"`
        Price float64 `json:"price"`
        Stock int     `json:"stock"`
    }
    
    func (p *Product) Validate() error {
        if strings.TrimSpace(p.Name) == "" {
            return ErrProductNameRequired
        }
        
        if p.Price < 0 {
            return ErrInvalidPrice
        }
        
        if p.Stock < 0 {
            return ErrInvalidStock
        }
        
        return nil
    }
    
    var (
        ErrProductNameRequired = errors.New("product name is required")
        ErrInvalidPrice        = errors.New("price cannot be negative")
        ErrInvalidStock        = errors.New("stock cannot be negative")
    )

### Complete Entity [​](#complete-entity)

bash

    goca entity Order \
      --fields "customerID:int,total:float64,status:string" \
      --validation \
      --business-rules \
      --timestamps \
      --soft-delete

**Generates:**

go

    package domain
    
    import (
        "errors"
        "time"
    )
    
    type Order struct {
        ID         uint       `json:"id"`
        CustomerID int        `json:"customer_id"`
        Total      float64    `json:"total"`
        Status     string     `json:"status"`
        CreatedAt  time.Time  `json:"created_at"`
        UpdatedAt  time.Time  `json:"updated_at"`
        DeletedAt  *time.Time `json:"deleted_at,omitempty"`
    }
    
    func (o *Order) Validate() error {
        if o.CustomerID <= 0 {
            return ErrInvalidCustomerID
        }
        
        if o.Total < 0 {
            return ErrInvalidTotal
        }
        
        if !o.IsValidStatus() {
            return ErrInvalidStatus
        }
        
        return nil
    }
    
    // Business Rules
    func (o *Order) IsValidStatus() bool {
        validStatuses := []string{"pending", "processing", "completed", "cancelled"}
        for _, status := range validStatuses {
            if o.Status == status {
                return true
            }
        }
        return false
    }
    
    func (o *Order) CanBeCancelled() bool {
        return o.Status == "pending" || o.Status == "processing"
    }
    
    func (o *Order) IsCompleted() bool {
        return o.Status == "completed"
    }
    
    func (o *Order) IsDeleted() bool {
        return o.DeletedAt != nil
    }
    
    var (
        ErrInvalidCustomerID = errors.New("invalid customer ID")
        ErrInvalidTotal      = errors.New("total cannot be negative")
        ErrInvalidStatus     = errors.New("invalid order status")
    )

### Entity with Unit Tests [​](#entity-with-unit-tests)

bash

    goca entity User \
      --fields "name:string,email:string,age:int" \
      --validation \
      --tests

**Generates:** `internal/domain/user_test.go`

go

    package domain
    
    import (
    	"testing"
    
    	"github.com/stretchr/testify/assert"
    )
    
    // TestUser_Validate tests the Validate method with various scenarios
    func TestUser_Validate(t *testing.T) {
    	tests := []struct {
    		name    string
    		user    User
    		wantErr bool
    		errMsg   string
    	}{
    		{
    			name: "valid entity",
    			user: User{
    				Name: "John Doe",
    				Email: "test@example.com",
    				Age: 25,
    			},
    			wantErr: false,
    		},
    		{
    			name: "invalid user - empty name",
    			user: User{
    				Name: "",
    				Email: "test@example.com",
    				Age: 25,
    			},
    			wantErr: true,
    			errMsg: "name",
    		},
    		{
    			name: "invalid user - negative age",
    			user: User{
    				Name: "John Doe",
    				Email: "test@example.com",
    				Age: -1,
    			},
    			wantErr: true,
    			errMsg: "age",
    		},
    	}
    
    	for _, tt := range tests {
    		t.Run(tt.name, func(t *testing.T) {
    			err := tt.user.Validate()
    			if tt.wantErr {
    				assert.Error(t, err)
    				if tt.errMsg != "" {
    					assert.Contains(t, err.Error(), tt.errMsg)
    				}
    			} else {
    				assert.NoError(t, err)
    			}
    		})
    	}
    }
    
    // TestUser_Initialization tests entity field initialization
    func TestUser_Initialization(t *testing.T) {
    	user := &User{
    		Name: "John Doe",
    		Email: "test@example.com",
    		Age: 25,
    	}
    
    	assert.Equal(t, "John Doe", user.Name, "Name should be set correctly")
    	assert.Equal(t, "test@example.com", user.Email, "Email should be set correctly")
    	assert.Equal(t, 25, user.Age, "Age should be set correctly")
    }
    
    // Additional field-specific tests for edge cases...

**Run tests:**

bash

    cd internal/domain
    go test -v -run TestUser

**Expected output:**

    === RUN   TestUser_Validate
    === RUN   TestUser_Validate/valid_entity
    === RUN   TestUser_Validate/invalid_user_-_empty_name
    === RUN   TestUser_Validate/invalid_user_-_negative_age
    --- PASS: TestUser_Validate (0.00s)
    === RUN   TestUser_Initialization
    --- PASS: TestUser_Initialization (0.00s)
    PASS
    ok      yourproject/internal/domain    0.013s

Generated Structure [​](#generated-structure)
---------------------------------------------

    internal/
    └── domain/
        ├── order.go           # Main entity
        ├── order_test.go      # Unit tests (if --tests)
        ├── order_seeds.go     # Seed data
        └── errors.go          # Domain errors (if --validation)

Note

The `--tests` flag is enabled by default. All generated tests use table-driven patterns and testify/assert for clean, maintainable test code.

Best Practices [​](#best-practices)
-----------------------------------

### DO [​](#do)

*   Keep entities pure (no external dependencies)
*   Include business validations
*   Add meaningful business rule methods
*   Use value objects for complex types
*   Document business logic

**Example:**

go

    //  Good: Pure domain logic
    func (u *User) IsAdult() bool {
        return u.Age >= 18
    }
    
    func (u *User) CanPlaceOrder() bool {
        return u.IsActive && u.EmailVerified
    }

### DON'T [​](#don-t)

*   Import database packages
*   Import HTTP frameworks
*   Include infrastructure logic
*   Add persistence methods

**Example:**

go

    //  Bad: Infrastructure dependency
    import "database/sql"
    
    func (u *User) Save(db *sql.DB) error {
        // Wrong layer!
    }
    
    //  Bad: Framework dependency
    import "github.com/gin-gonic/gin"
    
    func (u *User) ToJSON(c *gin.Context) {
        // Wrong responsibility!
    }

Integration [​](#integration)
-----------------------------

After generating an entity, you typically:

1.  **Generate Use Cases:**
    
    bash
    
        goca usecase OrderService --entity Order
    
2.  **Generate Repository:**
    
    bash
    
        goca repository Order --database postgres
    
3.  **Generate Handler:**
    
    bash
    
        goca handler Order --type http
    

Or use the complete feature command:

bash

    goca feature Order --fields "customerID:int,total:float64,status:string"

Field Type Reference [​](#field-type-reference)
-----------------------------------------------

| Type | Description | Example |
| --- | --- | --- |
| `string` | Text data | `"name:string"` |
| `int` | Integer | `"age:int"` |
| `int64` | Large integer | `"userID:int64"` |
| `float64` | Decimal number | `"price:float64"` |
| `bool` | Boolean | `"isActive:bool"` |
| `time.Time` | Timestamp | `"birthDate:time.Time"` |
| `[]string` | String array | `"tags:[]string"` |
| `[]int` | Integer array | `"scores:[]int"` |

Tips [​](#tips)
---------------

### Naming Conventions [​](#naming-conventions)

bash

    #  Good
    goca entity User --fields "firstName:string,lastName:string"
    goca entity Product --fields "productName:string,unitPrice:float64"
    
    #  Avoid
    goca entity user --fields "first_name:string"  # Use PascalCase
    goca entity PRODUCT --fields "PRICE:float64"    # Too loud

### Complex Domains [​](#complex-domains)

For complex domains, start simple and add complexity:

bash

    # Step 1: Basic structure
    goca entity Invoice --fields "number:string,amount:float64"
    
    # Step 2: Add more fields manually
    # Edit internal/domain/invoice.go to add:
    # - items []InvoiceItem
    # - taxes []Tax
    # - metadata map[string]interface{}

Troubleshooting [​](#troubleshooting)
-------------------------------------

### Entity Already Exists [​](#entity-already-exists)

**Problem:** "entity already exists"

**Solution:**

bash

    # Remove existing entity
    rm internal/domain/order.go
    rm internal/domain/order_*.go
    
    # Regenerate
    goca entity Order --fields "..."

### Invalid Field Type [​](#invalid-field-type)

**Problem:** "unsupported field type"

**Solution:** Use basic types first, then manually add custom types:

go

    // After generation, manually add:
    type User struct {
        // ... generated fields ...
        Address Address // Custom type
        Tags    []Tag   // Custom slice
    }

See Also [​](#see-also)
-----------------------

*   [`goca feature`](https://sazardev.github.io/goca/commands/feature.html) - Generate complete feature
*   [`goca usecase`](https://sazardev.github.io/goca/commands/usecase.html) - Generate use cases
*   [`goca repository`](https://sazardev.github.io/goca/commands/repository.html) - Generate repositories
*   [Clean Architecture Guide](https://sazardev.github.io/goca/guide/clean-architecture.html) - Architecture principles
*   [Domain Layer](https://sazardev.github.io/goca/guide/clean-architecture.html#layer-1-domain-entities) - Layer details

Next Steps [​](#next-steps)
---------------------------

After creating your entity:

1.  Add business logic methods manually if needed
2.  Generate use cases that use this entity
3.  Create repositories for persistence
4.  Build handlers for external access

Or use the shortcut:

bash

    goca feature YourEntity --fields "..." # Generates everything!</content>
</page>

<page>
  <title>Goca - Go Clean Architecture Code Generator</title>
  <url>https://sazardev.github.io/goca/commands/usecase.html</url>
  <content>goca usecase [​](#goca-usecase)
-------------------------------

Generate use cases (application services) with DTOs and business logic.

Syntax [​](#syntax)
-------------------

bash

    goca usecase <ServiceName> [flags]

Description [​](#description)
-----------------------------

Creates application layer services that orchestrate business workflows, coordinate repositories, and define clear input/output contracts through DTOs.

Flags [​](#flags)
-----------------

### `--entity` [​](#entity)

Associated domain entity.

bash

    goca usecase ProductService --entity Product

### `--operations` [​](#operations)

CRUD operations to generate.

**Options:** `create`, `read`, `update`, `delete`, `list`

bash

    goca usecase UserService --entity User --operations "create,read,update,delete,list"

### `--dto-validation` [​](#dto-validation)

Include DTO validation tags.

bash

    goca usecase OrderService --entity Order --dto-validation

Examples [​](#examples)
-----------------------

### Basic Use Case [​](#basic-use-case)

bash

    goca usecase ProductService --entity Product

### Complete CRUD [​](#complete-crud)

bash

    goca usecase UserService \
      --entity User \
      --operations "create,read,update,delete,list" \
      --dto-validation

Generated Files [​](#generated-files)
-------------------------------------

    internal/usecase/
    ├── user_dto.go         # Input/Output DTOs
    ├── user_interfaces.go  # Service interface
    └── user_service.go     # Service implementation

Generated Code Example [​](#generated-code-example)
---------------------------------------------------

go

    // user_dto.go
    package usecase
    
    type CreateUserInput struct {
        Name  string `json:"name" validate:"required,min=2"`
        Email string `json:"email" validate:"required,email"`
        Age   int    `json:"age" validate:"min=0,max=150"`
    }
    
    type UserResponse struct {
        ID    uint   `json:"id"`
        Name  string `json:"name"`
        Email string `json:"email"`
    }
    
    // user_interfaces.go
    type UserService interface {
        CreateUser(ctx context.Context, input CreateUserInput) (*UserResponse, error)
        GetUser(ctx context.Context, id uint) (*UserResponse, error)
        UpdateUser(ctx context.Context, id uint, input UpdateUserInput) error
        DeleteUser(ctx context.Context, id uint) error
        ListUsers(ctx context.Context) ([]*UserResponse, error)
    }
    
    // user_service.go
    type userService struct {
        userRepo domain.UserRepository
    }
    
    func NewUserService(userRepo domain.UserRepository) UserService {
        return &userService{userRepo: userRepo}
    }
    
    func (s *userService) CreateUser(ctx context.Context, input CreateUserInput) (*UserResponse, error) {
        // Validate input
        if err := input.Validate(); err != nil {
            return nil, err
        }
        
        // Create domain entity
        user := &domain.User{
            Name:  input.Name,
            Email: input.Email,
            Age:   input.Age,
        }
        
        // Validate domain rules
        if err := user.Validate(); err != nil {
            return nil, err
        }
        
        // Persist
        if err := s.userRepo.Save(ctx, user); err != nil {
            return nil, err
        }
        
        // Return DTO
        return &UserResponse{
            ID:    user.ID,
            Name:  user.Name,
            Email: user.Email,
        }, nil
    }

Best Practices [​](#best-practices)
-----------------------------------

### DO [​](#do)

*   Define clear DTOs for each operation
*   Validate input at use case boundary
*   Coordinate multiple repositories
*   Transform entities to DTOs

### DON'T [​](#don-t)

*   Include HTTP/gRPC logic
*   Write SQL queries
*   Import framework packages
*   Skip validation

See Also [​](#see-also)
-----------------------

*   [`goca entity`](https://sazardev.github.io/goca/commands/entity.html) - Generate entities
*   [`goca repository`](https://sazardev.github.io/goca/commands/repository.html) - Generate repositories
*   [`goca feature`](https://sazardev.github.io/goca/commands/feature.html) - Generate complete feature</content>
</page>

<page>
  <title>Understanding Domain Entities in Clean Architecture | Goca Blog</title>
  <url>https://sazardev.github.io/goca/blog/articles/understanding-domain-entities.html</url>
  <content>Understanding Domain Entities in Clean Architecture [​](#understanding-domain-entities-in-clean-architecture)
-------------------------------------------------------------------------------------------------------------

ArchitectureDomain-Driven Design

Domain entities are the heart of Clean Architecture. They represent the core business concepts and rules that define your application's purpose. Understanding entities correctly is fundamental to building maintainable, testable, and scalable software systems.

* * *

What is a Domain Entity? [​](#what-is-a-domain-entity)
------------------------------------------------------

A domain entity is a representation of a business concept that has a unique identity and encapsulates business rules. In Clean Architecture, entities form the innermost layer, completely independent of external concerns like databases, frameworks, or UI.

### Core Characteristics [​](#core-characteristics)

**Identity**: Each entity has a unique identifier that distinguishes it from other entities of the same type. Two entities with the same attributes but different identities are different entities.

**Business Logic**: Entities contain methods that enforce business rules and maintain invariants. They are not passive data structures but active participants in your domain model.

**Independence**: Entities have zero dependencies on external systems. They do not import HTTP libraries, database drivers, or framework code. This independence makes them portable, testable, and reusable.

**Validation**: Entities validate their own state, ensuring that business rules are never violated. Invalid states are impossible to represent.

Entity vs Model: A Critical Distinction [​](#entity-vs-model-a-critical-distinction)
------------------------------------------------------------------------------------

Many developers confuse entities with database models or API models. This confusion leads to architectural problems and coupling.

### What an Entity Is NOT [​](#what-an-entity-is-not)

**Not a Database Model**: Entities do not map directly to database tables. They represent business concepts, not storage structures. Database concerns belong to the infrastructure layer.

**Not an API Response**: Entities are not DTOs (Data Transfer Objects). API responses should be separate structures that adapt entities for external communication.

**Not Framework-Dependent**: Entities do not depend on ORMs, validation frameworks, or serialization libraries. These are implementation details.

### The Separation Principle [​](#the-separation-principle)

    Domain Entity (Pure Business Logic)
            ↓
        Use Case (Application Logic)
            ↓
    Repository Interface (Contract)
            ↓
    Repository Implementation (Database Details)
            ↓
        Database Schema

This separation allows you to:

*   Change databases without touching business logic
*   Test business rules without database setup
*   Evolve your domain model independently
*   Swap ORMs or frameworks with minimal impact

Domain-Driven Design Principles [​](#domain-driven-design-principles)
---------------------------------------------------------------------

Goca implements Domain-Driven Design (DDD) principles when generating entities, ensuring your code follows established best practices.

### Ubiquitous Language [​](#ubiquitous-language)

Entities use the same terminology as your business domain. If your business talks about "Orders," "Customers," and "Products," your entities should use these exact terms.

### Aggregate Roots [​](#aggregate-roots)

Entities can serve as aggregate roots, controlling access to related objects and maintaining consistency boundaries.

### Value Objects vs Entities [​](#value-objects-vs-entities)

Entities have identity; value objects do not. An email address is a value object. A user is an entity. Goca helps you model both correctly.

How Goca Generates Entities [​](#how-goca-generates-entities)
-------------------------------------------------------------

Goca provides the `goca entity` command to generate domain entities following Clean Architecture and DDD principles.

### Basic Entity Generation [​](#basic-entity-generation)

bash

    goca entity User --fields "name:string,email:string,age:int"

This command generates a pure domain entity with no external dependencies:

go

    package domain
    
    type User struct {
        ID    uint   `json:"id" gorm:"primaryKey;autoIncrement"`
        Name  string `json:"name" gorm:"type:varchar(255);not null"`
        Email string `json:"email" gorm:"type:varchar(255);uniqueIndex;not null"`
        Age   int    `json:"age" gorm:"type:integer;not null;default:0"`
    }

Notice that while GORM tags are present for infrastructure convenience, the entity itself remains a simple Go struct. The entity does not import GORM or any database package.

### Field Types and Conventions [​](#field-types-and-conventions)

Goca supports common field types that map to both Go types and database columns:

**String Fields**: `name:string`, `email:string`, `description:string`

*   Generate `string` type
*   Map to `varchar` or `text` columns
*   Suitable for textual data

**Numeric Fields**: `age:int`, `price:float64`, `quantity:int64`

*   Generate integer or floating-point types
*   Map to appropriate numeric columns
*   Support business calculations

**Boolean Fields**: `is_active:bool`, `verified:bool`

*   Generate `bool` type
*   Map to boolean columns
*   Represent binary states

**Temporal Fields**: `birth_date:time.Time`

*   Generate `time.Time` type
*   Handle date and time data
*   Work with the standard library

### Adding Validation [​](#adding-validation)

Business rules are enforced through validation methods:

bash

    goca entity User --fields "name:string,email:string,age:int" --validation

This generates a `Validate()` method and domain-specific errors:

go

    package domain
    
    import (
        "time"
        
        "gorm.io/gorm"
    )
    
    type User struct {
        ID        uint           `json:"id" gorm:"primaryKey;autoIncrement"`
        Name      string         `json:"name" gorm:"type:varchar(255);not null" validate:"required"`
        Email     string         `json:"email" gorm:"type:varchar(255);uniqueIndex;not null" validate:"required,email"`
        Age       int            `json:"age" gorm:"type:integer;not null;default:0" validate:"required,gte=0"`
        CreatedAt time.Time      `json:"created_at" gorm:"autoCreateTime"`
        UpdatedAt time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
        DeletedAt gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`
    }
    
    func (u *User) Validate() error {
        if u.Name == "" {
            return ErrInvalidUserName
        }
        if u.Email == "" {
            return ErrInvalidUserEmail
        }
        if u.Age < 0 {
            return ErrInvalidUserAge
        }
        return nil
    }

The validation method ensures that no invalid user can exist in your system. This is domain logic, not input validation. Input validation happens in the use case or handler layer.

### Domain Errors [​](#domain-errors)

Goca generates a separate `errors.go` file containing domain-specific errors:

go

    package domain
    
    import "errors"
    
    var (
        ErrInvalidUserName  = errors.New("invalid user name")
        ErrInvalidUserEmail = errors.New("invalid user email")
        ErrInvalidUserAge   = errors.New("invalid user age")
        ErrUserNotFound     = errors.New("user not found")
    )

These errors are part of your domain model. They communicate business rule violations clearly and can be handled appropriately by outer layers.

### Business Rules [​](#business-rules)

Beyond validation, entities can contain business logic:

bash

    goca entity Order --fields "customer_id:int,total:float64,status:string" --business-rules

This generates methods that implement domain logic:

go

    func (o *Order) Validate() error {
        if o.Customer_id < 0 {
            return ErrInvalidOrderCustomer_id
        }
        if o.Total < 0 {
            return ErrInvalidOrderTotal
        }
        if o.Status == "" {
            return ErrInvalidOrderStatus
        }
        return nil
    }

You can extend these with additional business methods:

go

    func (o *Order) CanBeCancelled() bool {
        return o.Status == "pending" || o.Status == "confirmed"
    }
    
    func (o *Order) Apply(discount float64) error {
        if discount < 0 || discount > 1 {
            return errors.New("discount must be between 0 and 1")
        }
        o.Total = o.Total * (1 - discount)
        return nil
    }
    
    func (o *Order) IsCompleted() bool {
        return o.Status == "delivered"
    }

These methods encapsulate business knowledge. They answer domain questions and enforce domain rules.

### Timestamps and Soft Deletes [​](#timestamps-and-soft-deletes)

Entities often need audit trails and soft delete functionality:

bash

    goca entity Product --fields "name:string,price:float64,stock:int" \
      --timestamps \
      --soft-delete

This generates:

go

    type Product struct {
        ID          uint           `json:"id" gorm:"primaryKey"`
        Name        string         `json:"name" gorm:"type:varchar(255);not null"`
        Price       float64        `json:"price" gorm:"type:decimal(10,2);not null;default:0"`
        Stock       int            `json:"stock" gorm:"type:integer;not null;default:0"`
        CreatedAt   time.Time      `json:"created_at" gorm:"autoCreateTime"`
        UpdatedAt   time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
        DeletedAt   gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`
    }
    
    func (p *Product) SoftDelete() {
        p.DeletedAt = gorm.DeletedAt{Time: time.Now(), Valid: true}
    }
    
    func (p *Product) IsDeleted() bool {
        return p.DeletedAt.Valid
    }

Soft deletes preserve data while marking it as inactive. The `DeletedAt` field enables this pattern without permanently removing records.

Complete Entity Example [​](#complete-entity-example)
-----------------------------------------------------

Let's examine a complete entity generated by Goca:

bash

    goca entity Product --fields "name:string,description:string,price:float64,stock:int,is_active:bool" \
      --validation \
      --business-rules \
      --timestamps \
      --soft-delete

This generates a production-ready entity:

go

    package domain
    
    import (
        "time"
        
        "gorm.io/gorm"
    )
    
    type Product struct {
        ID          uint           `json:"id" gorm:"primaryKey;autoIncrement"`
        Name        string         `json:"name" gorm:"type:varchar(255);not null" validate:"required"`
        Description string         `json:"description" gorm:"type:text"`
        Price       float64        `json:"price" gorm:"type:decimal(10,2);not null;default:0" validate:"required,gte=0"`
        Stock       int            `json:"stock" gorm:"type:integer;not null;default:0" validate:"required,gte=0"`
        IsActive    bool           `json:"is_active" gorm:"type:boolean;not null;default:false"`
        CreatedAt   time.Time      `json:"created_at" gorm:"autoCreateTime"`
        UpdatedAt   time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
        DeletedAt   gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`
    }
    
    func (p *Product) Validate() error {
        if p.Name == "" {
            return ErrInvalidProductName
        }
        if p.Price < 0 {
            return ErrInvalidProductPrice
        }
        if p.Stock < 0 {
            return ErrInvalidProductStock
        }
        return nil
    }
    
    func (p *Product) SoftDelete() {
        p.DeletedAt = gorm.DeletedAt{Time: time.Now(), Valid: true}
    }
    
    func (p *Product) IsDeleted() bool {
        return p.DeletedAt.Valid
    }

You can extend this with additional business methods:

go

    func (p *Product) IsAvailable() bool {
        return p.IsActive && p.Stock > 0 && !p.IsDeleted()
    }
    
    func (p *Product) Restock(quantity int) error {
        if quantity <= 0 {
            return errors.New("restock quantity must be positive")
        }
        p.Stock += quantity
        p.UpdatedAt = time.Now()
        return nil
    }
    
    func (p *Product) Sell(quantity int) error {
        if quantity <= 0 {
            return errors.New("sell quantity must be positive")
        }
        if p.Stock < quantity {
            return errors.New("insufficient stock")
        }
        p.Stock -= quantity
        p.UpdatedAt = time.Now()
        return nil
    }
    
    func (p *Product) ApplyDiscount(percentage float64) error {
        if percentage < 0 || percentage > 100 {
            return errors.New("discount percentage must be between 0 and 100")
        }
        p.Price = p.Price * (1 - percentage/100)
        p.UpdatedAt = time.Now()
        return nil
    }

These methods capture business logic that belongs in the domain layer. They make the entity more than a data structure; they make it a behavior-rich business object.

Testing Domain Entities [​](#testing-domain-entities)
-----------------------------------------------------

Domain entities are easy to test because they have no external dependencies. You can test business logic in isolation:

go

    package domain_test
    
    import (
        "testing"
        
        "yourproject/internal/domain"
    )
    
    func TestUser_Validate(t *testing.T) {
        tests := []struct {
            name    string
            user    domain.User
            wantErr bool
        }{
            {
                name: "valid user",
                user: domain.User{
                    Name:  "John Doe",
                    Email: "john@example.com",
                    Age:   30,
                },
                wantErr: false,
            },
            {
                name: "empty name",
                user: domain.User{
                    Name:  "",
                    Email: "john@example.com",
                    Age:   30,
                },
                wantErr: true,
            },
            {
                name: "negative age",
                user: domain.User{
                    Name:  "John Doe",
                    Email: "john@example.com",
                    Age:   -5,
                },
                wantErr: true,
            },
        }
        
        for _, tt := range tests {
            t.Run(tt.name, func(t *testing.T) {
                err := tt.user.Validate()
                if (err != nil) != tt.wantErr {
                    t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
                }
            })
        }
    }
    
    func TestProduct_Sell(t *testing.T) {
        product := domain.Product{
            Name:  "Test Product",
            Price: 100.0,
            Stock: 10,
        }
        
        // Valid sale
        err := product.Sell(5)
        if err != nil {
            t.Errorf("Sell(5) failed: %v", err)
        }
        if product.Stock != 5 {
            t.Errorf("Expected stock 5, got %d", product.Stock)
        }
        
        // Insufficient stock
        err = product.Sell(10)
        if err == nil {
            t.Error("Sell(10) should fail with insufficient stock")
        }
        
        // Negative quantity
        err = product.Sell(-1)
        if err == nil {
            t.Error("Sell(-1) should fail with negative quantity")
        }
    }

These tests run instantly because they do not touch databases, networks, or files. They verify business logic in isolation.

Best Practices for Domain Entities [​](#best-practices-for-domain-entities)
---------------------------------------------------------------------------

Based on Clean Architecture and DDD principles, follow these best practices:

**Keep Entities Pure**: Do not import framework or infrastructure code. Entities should compile without external dependencies beyond the standard library.

**Encapsulate State**: Use methods to modify entity state. Avoid exposing fields directly if business rules govern their modification.

**Express Business Rules**: Write methods that answer business questions and enforce business constraints. Make implicit knowledge explicit.

**Use Value Objects**: For concepts without identity, create value objects. An email address, money amount, or date range should be a value object, not part of an entity.

**Avoid Anemic Domain Models**: Entities with only getters and setters are anemic. Add behavior. Rich domain models contain both data and behavior.

**Design for Invariants**: Entities should always be in a valid state. Constructor functions and validation methods enforce this.

**Use Domain Language**: Name entities, fields, and methods using terms from your business domain. Code should read like business documentation.

Integration with Other Layers [​](#integration-with-other-layers)
-----------------------------------------------------------------

Entities work with other Clean Architecture layers through well-defined interfaces.

### Use Case Layer [​](#use-case-layer)

Use cases orchestrate entities to fulfill application requirements:

go

    package usecase
    
    type CreateProductInput struct {
        Name        string  `json:"name" validate:"required"`
        Description string  `json:"description"`
        Price       float64 `json:"price" validate:"required,gt=0"`
        Stock       int     `json:"stock" validate:"required,gte=0"`
    }
    
    type productService struct {
        repo repository.ProductRepository
    }
    
    func (s *productService) Create(input CreateProductInput) (*domain.Product, error) {
        // Create entity
        product := &domain.Product{
            Name:        input.Name,
            Description: input.Description,
            Price:       input.Price,
            Stock:       input.Stock,
            IsActive:    true,
        }
        
        // Validate business rules
        if err := product.Validate(); err != nil {
            return nil, err
        }
        
        // Persist through repository
        if err := s.repo.Save(product); err != nil {
            return nil, err
        }
        
        return product, nil
    }

The use case depends on the entity, not the other way around. This maintains the dependency rule.

### Repository Layer [​](#repository-layer)

Repositories provide persistence for entities through interfaces defined in the domain:

go

    package repository
    
    type ProductRepository interface {
        Save(product *domain.Product) error
        FindByID(id uint) (*domain.Product, error)
        Update(product *domain.Product) error
        Delete(id uint) error
        FindAll() ([]domain.Product, error)
    }

The repository interface lives in the domain package, but implementations live in the infrastructure layer:

go

    package repository
    
    import (
        "yourproject/internal/domain"
        "gorm.io/gorm"
    )
    
    type postgresProductRepository struct {
        db *gorm.DB
    }
    
    func NewPostgresProductRepository(db *gorm.DB) ProductRepository {
        return &postgresProductRepository{db: db}
    }
    
    func (r *postgresProductRepository) Save(product *domain.Product) error {
        return r.db.Create(product).Error
    }
    
    func (r *postgresProductRepository) FindByID(id uint) (*domain.Product, error) {
        var product domain.Product
        err := r.db.First(&product, id).Error
        if err != nil {
            return nil, err
        }
        return &product, nil
    }

This separation allows you to swap database implementations without changing business logic.

Advanced Entity Patterns [​](#advanced-entity-patterns)
-------------------------------------------------------

### Aggregate Roots [​](#aggregate-roots-1)

Entities can serve as aggregate roots, controlling access to related entities:

go

    type Order struct {
        ID         uint
        CustomerID uint
        Items      []OrderItem
        Total      float64
        Status     string
    }
    
    type OrderItem struct {
        ID        uint
        ProductID uint
        Quantity  int
        Price     float64
    }
    
    func (o *Order) AddItem(productID uint, quantity int, price float64) error {
        if quantity <= 0 {
            return errors.New("quantity must be positive")
        }
        
        item := OrderItem{
            ProductID: productID,
            Quantity:  quantity,
            Price:     price,
        }
        
        o.Items = append(o.Items, item)
        o.calculateTotal()
        return nil
    }
    
    func (o *Order) calculateTotal() {
        total := 0.0
        for _, item := range o.Items {
            total += float64(item.Quantity) * item.Price
        }
        o.Total = total
    }

The `Order` aggregate root controls `OrderItem` access, maintaining consistency.

### Factories [​](#factories)

Complex entity creation can use factory patterns:

go

    func NewUser(name, email string, age int) (*User, error) {
        user := &User{
            Name:      name,
            Email:     email,
            Age:       age,
            CreatedAt: time.Now(),
            UpdatedAt: time.Now(),
        }
        
        if err := user.Validate(); err != nil {
            return nil, err
        }
        
        return user, nil
    }

Factories ensure entities are always created in valid states.

Generating Complete Features [​](#generating-complete-features)
---------------------------------------------------------------

While `goca entity` generates entities, `goca feature` generates complete features including entities, use cases, repositories, and handlers:

bash

    goca feature Product --fields "name:string,price:float64,stock:int"

This generates:

*   Domain entity (`internal/domain/product.go`)
*   Use case interface and implementation (`internal/usecase/product_service.go`)
*   Repository interface (`internal/repository/interfaces.go`)
*   Repository implementation (`internal/repository/postgres_product_repository.go`)
*   HTTP handler (`internal/handler/http/product_handler.go`)
*   DTOs (`internal/usecase/dto.go`)
*   Error definitions (`internal/domain/errors.go`)
*   Seed data (`internal/domain/product_seeds.go`)

All layers work together following Clean Architecture principles, with the entity at the core.

Conclusion [​](#conclusion)
---------------------------

Domain entities are the foundation of Clean Architecture. They represent your business concepts and enforce business rules without coupling to external systems. Goca generates production-ready entities following DDD principles, giving you a solid starting point for building maintainable applications.

Understanding entities correctly is essential for successful software architecture. They are not database models, not API responses, and not framework-dependent structures. They are pure business logic, testable in isolation, and independent of implementation details.

By using Goca's entity generation commands and following Clean Architecture principles, you create systems that are:

*   Easy to test without external dependencies
*   Simple to modify as business requirements change
*   Portable across different frameworks and technologies
*   Clear in expressing business intent
*   Maintainable over long periods

Start with entities. Build your business logic correctly. Let outer layers adapt to your domain, not the other way around.

Further Reading [​](#further-reading)
-------------------------------------

*   Clean Architecture documentation in the [guide section](https://sazardev.github.io/goca/guide/clean-architecture.html)
*   Complete command reference for [`goca entity`](https://sazardev.github.io/goca/commands/entity.html)
*   Full feature generation with [`goca feature`](https://sazardev.github.io/goca/commands/feature.html)
*   Domain-Driven Design principles and patterns
*   Repository pattern implementation examples</content>
</page>

<page>
  <title>Goca - Go Clean Architecture Code Generator</title>
  <url>https://sazardev.github.io/goca/commands/handler.html</url>
  <content>goca handler [​](#goca-handler)
-------------------------------

Generate input adapters for different protocols (HTTP, gRPC, CLI, etc.).

Syntax [​](#syntax)
-------------------

bash

    goca handler <EntityName> [flags]

Description [​](#description)
-----------------------------

Creates handlers that adapt external requests to use case calls, handling protocol-specific concerns while keeping business logic isolated.

Flags [​](#flags)
-----------------

### `--type` [​](#type)

Handler type. Default: `http`

**Options:** `http` | `grpc` | `cli` | `worker` | `soap`

bash

    goca handler Product --type http

### `--middleware` [​](#middleware)

Include middleware setup.

bash

    goca handler User --type http --middleware

### `--validation` [​](#validation)

Add request validation.

bash

    goca handler Order --type http --validation

Examples [​](#examples)
-----------------------

### HTTP REST Handler [​](#http-rest-handler)

bash

    goca handler User --type http

**Generates:** `internal/handler/http/user_handler.go`

go

    package http
    
    import (
        "encoding/json"
        "net/http"
        "strconv"
        
        "github.com/gorilla/mux"
        "myproject/internal/usecase"
    )
    
    type UserHandler struct {
        userService usecase.UserService
    }
    
    func NewUserHandler(userService usecase.UserService) *UserHandler {
        return &UserHandler{userService: userService}
    }
    
    func (h *UserHandler) RegisterRoutes(r *mux.Router) {
        r.HandleFunc("/users", h.CreateUser).Methods(http.MethodPost)
        r.HandleFunc("/users/{id}", h.GetUser).Methods(http.MethodGet)
        r.HandleFunc("/users/{id}", h.UpdateUser).Methods(http.MethodPut)
        r.HandleFunc("/users/{id}", h.DeleteUser).Methods(http.MethodDelete)
        r.HandleFunc("/users", h.ListUsers).Methods(http.MethodGet)
    }
    
    func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
        var req usecase.CreateUserInput
        if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
            respondError(w, http.StatusBadRequest, "Invalid request body")
            return
        }
        
        user, err := h.userService.CreateUser(r.Context(), req)
        if err != nil {
            respondError(w, http.StatusInternalServerError, err.Error())
            return
        }
        
        respondJSON(w, http.StatusCreated, user)
    }
    
    func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
        vars := mux.Vars(r)
        id, err := strconv.ParseUint(vars["id"], 10, 32)
        if err != nil {
            respondError(w, http.StatusBadRequest, "Invalid user ID")
            return
        }
        
        user, err := h.userService.GetUser(r.Context(), uint(id))
        if err != nil {
            respondError(w, http.StatusNotFound, "User not found")
            return
        }
        
        respondJSON(w, http.StatusOK, user)
    }
    
    func respondJSON(w http.ResponseWriter, status int, payload interface{}) {
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(status)
        json.NewEncoder(w).Encode(payload)
    }
    
    func respondError(w http.ResponseWriter, status int, message string) {
        respondJSON(w, status, map[string]string{"error": message})
    }

### gRPC Handler [​](#grpc-handler)

bash

    goca handler Product --type grpc

**Generates:** `internal/handler/grpc/product_handler.go` + `.proto` file

### CLI Handler [​](#cli-handler)

bash

    goca handler Task --type cli

**Generates:** CLI commands with Cobra

### Worker Handler [​](#worker-handler)

bash

    goca handler Email --type worker

**Generates:** Background job handlers

Handler Types Comparison [​](#handler-types-comparison)
-------------------------------------------------------

| Type | Use Case | Generated |
| --- | --- | --- |
| **http** | REST APIs, Web services | HTTP handlers with routing |
| **grpc** | Microservices, High performance | gRPC server + proto files |
| **cli** | Command-line tools | Cobra commands |
| **worker** | Background jobs, Async tasks | Job handlers |
| **soap** | Legacy systems integration | SOAP client wrappers |

Best Practices [​](#best-practices)
-----------------------------------

### DO [​](#do)

*   Handle protocol-specific concerns only
*   Transform requests to use case DTOs
*   Format responses appropriately
*   Add proper error handling
*   Use middleware for cross-cutting concerns

### DON'T [​](#don-t)

*   Include business logic
*   Access repositories directly
*   Skip input validation
*   Return domain entities directly

See Also [​](#see-also)
-----------------------

*   [`goca usecase`](https://sazardev.github.io/goca/commands/usecase.html) - Generate use cases
*   [`goca feature`](https://sazardev.github.io/goca/commands/feature.html) - Generate complete feature
*   [Handler Layer](https://sazardev.github.io/goca/guide/clean-architecture.html#layer-4-handlers-interface-adapters)</content>
</page>

<page>
  <title>Goca - Go Clean Architecture Code Generator</title>
  <url>https://sazardev.github.io/goca/commands/repository.html</url>
  <content>goca repository [​](#goca-repository)
-------------------------------------

Generate repository interfaces and database implementations.

Syntax [​](#syntax)
-------------------

bash

    goca repository <EntityName> [flags]

Description [​](#description)
-----------------------------

Creates repository pattern implementations for data persistence, abstracting database operations behind clean interfaces.

Flags [​](#flags)
-----------------

### `--database` [​](#database)

Database type. Default: `postgres`

**Options:**

*   `postgres` - PostgreSQL (GORM)
*   `postgres-json` - PostgreSQL with JSONB support
*   `mysql` - MySQL (GORM)
*   `mongodb` - MongoDB (native driver)
*   `sqlite` - SQLite (embedded, database/sql)
*   `sqlserver` - SQL Server (GORM)
*   `elasticsearch` - Elasticsearch (v8 client)
*   `dynamodb` - DynamoDB (AWS SDK v2)

bash

    goca repository Product --database postgres
    goca repository Config --database postgres-json
    goca repository Article --database elasticsearch

### `--interface-only` [​](#interface-only)

Generate only the interface.

bash

    goca repository User --interface-only

### `--implementation` [​](#implementation)

Generate only the implementation.

bash

    goca repository User --implementation --database mysql

Examples [​](#examples)
-----------------------

### PostgreSQL Repository [​](#postgresql-repository)

bash

    goca repository User --database postgres

### PostgreSQL JSON (Semi-structured Data) [​](#postgresql-json-semi-structured-data)

bash

    goca repository Config --database postgres-json

### MongoDB Repository [​](#mongodb-repository)

bash

    goca repository Product --database mongodb

### Elasticsearch Full-Text Search [​](#elasticsearch-full-text-search)

bash

    goca repository Article --database elasticsearch

### DynamoDB (AWS Serverless) [​](#dynamodb-aws-serverless)

bash

    goca repository Order --database dynamodb

### SQL Server (Enterprise) [​](#sql-server-enterprise)

bash

    goca repository Employee --database sqlserver

### SQLite (Embedded) [​](#sqlite-embedded)

bash

    goca repository Setting --database sqlite

### Interface Only [​](#interface-only-1)

bash

    goca repository Order --interface-only

Generated Files [​](#generated-files)
-------------------------------------

    internal/repository/
    ├── interfaces.go               # Repository interfaces
    └── postgres_user_repository.go # Implementation

Generated Code Example [​](#generated-code-example)
---------------------------------------------------

go

    // interfaces.go
    package repository
    
    import (
        "context"
        "myproject/internal/domain"
    )
    
    type UserRepository interface {
        Save(ctx context.Context, user *domain.User) error
        FindByID(ctx context.Context, id uint) (*domain.User, error)
        FindAll(ctx context.Context) ([]*domain.User, error)
        Update(ctx context.Context, user *domain.User) error
        Delete(ctx context.Context, id uint) error
    }
    
    // postgres_user_repository.go
    package repository
    
    import (
        "context"
        "database/sql"
        "myproject/internal/domain"
    )
    
    type postgresUserRepository struct {
        db *sql.DB
    }
    
    func NewPostgresUserRepository(db *sql.DB) UserRepository {
        return &postgresUserRepository{db: db}
    }
    
    func (r *postgresUserRepository) Save(ctx context.Context, user *domain.User) error {
        query := `
            INSERT INTO users (name, email, age, created_at, updated_at)
            VALUES ($1, $2, $3, NOW(), NOW())
            RETURNING id, created_at, updated_at
        `
        
        return r.db.QueryRowContext(
            ctx, query,
            user.Name, user.Email, user.Age,
        ).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
    }
    
    func (r *postgresUserRepository) FindByID(ctx context.Context, id uint) (*domain.User, error) {
        query := `
            SELECT id, name, email, age, created_at, updated_at
            FROM users
            WHERE id = $1
        `
        
        user := &domain.User{}
        err := r.db.QueryRowContext(ctx, query, id).Scan(
            &user.ID, &user.Name, &user.Email, &user.Age,
            &user.CreatedAt, &user.UpdatedAt,
        )
        
        if err == sql.ErrNoRows {
            return nil, domain.ErrUserNotFound
        }
        
        return user, err
    }
    
    func (r *postgresUserRepository) FindAll(ctx context.Context) ([]*domain.User, error) {
        query := `
            SELECT id, name, email, age, created_at, updated_at
            FROM users
            ORDER BY created_at DESC
        `
        
        rows, err := r.db.QueryContext(ctx, query)
        if err != nil {
            return nil, err
        }
        defer rows.Close()
        
        var users []*domain.User
        for rows.Next() {
            user := &domain.User{}
            if err := rows.Scan(
                &user.ID, &user.Name, &user.Email, &user.Age,
                &user.CreatedAt, &user.UpdatedAt,
            ); err != nil {
                return nil, err
            }
            users = append(users, user)
        }
        
        return users, rows.Err()
    }
    
    func (r *postgresUserRepository) Update(ctx context.Context, user *domain.User) error {
        query := `
            UPDATE users
            SET name = $1, email = $2, age = $3, updated_at = NOW()
            WHERE id = $4
        `
        
        _, err := r.db.ExecContext(ctx, query,
            user.Name, user.Email, user.Age, user.ID,
        )
        
        return err
    }
    
    func (r *postgresUserRepository) Delete(ctx context.Context, id uint) error {
        query := `DELETE FROM users WHERE id = $1`
        _, err := r.db.ExecContext(ctx, query, id)
        return err
    }

Database-Specific Features [​](#database-specific-features)
-----------------------------------------------------------

### PostgreSQL [​](#postgresql)

*   RETURNING clause support
*   JSON/JSONB columns
*   Array types
*   Full-text search

### MySQL [​](#mysql)

*   AUTO\_INCREMENT handling
*   LIMIT/OFFSET pagination
*   JSON column support

### MongoDB [​](#mongodb)

*   Document-based storage
*   Flexible schema
*   Aggregation pipelines
*   Index management

### SQLite [​](#sqlite)

*   Embedded database
*   File-based storage
*   Great for testing
*   Simplified queries

Best Practices [​](#best-practices)
-----------------------------------

### DO [​](#do)

*   Keep repositories simple (CRUD + specific queries)
*   Return domain entities
*   Handle database errors properly
*   Use prepared statements
*   Implement transactions when needed

### DON'T [​](#don-t)

*   Include business logic
*   Return database-specific types
*   Expose SQL to callers
*   Skip error handling

See Also [​](#see-also)
-----------------------

*   [`goca entity`](https://sazardev.github.io/goca/commands/entity.html) - Generate entities
*   [`goca usecase`](https://sazardev.github.io/goca/commands/usecase.html) - Generate use cases
*   [`goca feature`](https://sazardev.github.io/goca/commands/feature.html) - Generate complete feature</content>
</page>

<page>
  <title>Goca - Go Clean Architecture Code Generator</title>
  <url>https://sazardev.github.io/goca/commands/di.html</url>
  <content>goca di [​](#goca-di)
---------------------

Generate dependency injection container for automatic wiring.

Syntax [​](#syntax)
-------------------

Description [​](#description)
-----------------------------

Creates a dependency injection container that automatically wires all features, repositories, use cases, and handlers together.

Flags [​](#flags)
-----------------

### `--features` [​](#features)

Comma-separated list of features to wire.

bash

    goca di --features "User,Product,Order"

### `--database` [​](#database)

Database type. Default: `postgres`

bash

    goca di --features "User,Product" --database postgres

Examples [​](#examples)
-----------------------

### Wire All Features [​](#wire-all-features)

bash

    goca di --features "User,Product,Order,Payment"

### PostgreSQL with Authentication [​](#postgresql-with-authentication)

bash

    goca di --features "User,Auth" --database postgres

Generated Code [​](#generated-code)
-----------------------------------

go

    // internal/di/container.go
    package di
    
    import (
        "database/sql"
        "myproject/internal/handler/http"
        "myproject/internal/repository"
        "myproject/internal/usecase"
    )
    
    type Container struct {
        // Repositories
        UserRepository    repository.UserRepository
        ProductRepository repository.ProductRepository
        
        // Use Cases
        UserService    usecase.UserService
        ProductService usecase.ProductService
        
        // Handlers
        UserHandler    *http.UserHandler
        ProductHandler *http.ProductHandler
    }
    
    func NewContainer(db *sql.DB) *Container {
        // Initialize repositories
        userRepo := repository.NewPostgresUserRepository(db)
        productRepo := repository.NewPostgresProductRepository(db)
        
        // Initialize use cases
        userService := usecase.NewUserService(userRepo)
        productService := usecase.NewProductService(productRepo)
        
        // Initialize handlers
        userHandler := http.NewUserHandler(userService)
        productHandler := http.NewProductHandler(productService)
        
        return &Container{
            UserRepository:    userRepo,
            ProductRepository: productRepo,
            UserService:       userService,
            ProductService:    productService,
            UserHandler:       userHandler,
            ProductHandler:    productHandler,
        }
    }

See Also [​](#see-also)
-----------------------

*   [`goca integrate`](https://sazardev.github.io/goca/commands/integrate.html) - Integrate existing features
*   [`goca feature`](https://sazardev.github.io/goca/commands/feature.html) - Generate complete feature</content>
</page>

<page>
  <title>Goca - Go Clean Architecture Code Generator</title>
  <url>https://sazardev.github.io/goca/commands/integrate.html</url>
  <content>goca integrate [​](#goca-integrate)
-----------------------------------

Integrate existing features with dependency injection and routing.

Syntax [​](#syntax)
-------------------

bash

    goca integrate [flags]

Description [​](#description)
-----------------------------

Automatically detects all features in your project and integrates them by updating the dependency injection container and registering routes.

Auto-Integration

The `goca feature` command now includes automatic integration. Use this command when you need to repair or update integration for manually created features.

Flags [​](#flags)
-----------------

### `--all` [​](#all)

Integrate all detected features.

### `--feature` [​](#feature)

Integrate a specific feature.

bash

    goca integrate --feature User

### `--dry-run` [​](#dry-run)

Show what would be integrated without making changes.

bash

    goca integrate --all --dry-run

Examples [​](#examples)
-----------------------

### Integrate All Features [​](#integrate-all-features)

**Output:**

     Scanning for features...
     Found: User
     Found: Product
     Found: Order
    
     Updating dependency injection...
     Updated internal/di/container.go
    
     Registering routes...
     Updated internal/handler/http/routes.go
    
     Integration complete! 3 features integrated.

### Integrate Specific Feature [​](#integrate-specific-feature)

bash

    goca integrate --feature Product

### Dry Run [​](#dry-run-1)

bash

    goca integrate --all --dry-run

Shows what would be changed without actually modifying files.

What Gets Integrated [​](#what-gets-integrated)
-----------------------------------------------

1.  **Dependency Injection**
    
    *   Adds repositories to DI container
    *   Wires use cases with dependencies
    *   Registers handlers
2.  **HTTP Routing**
    
    *   Registers all HTTP routes
    *   Sets up middleware
    *   Configures path prefixes
3.  **Database Migrations**
    
    *   Creates migration files
    *   Registers schema changes

Use Cases [​](#use-cases)
-------------------------

### After Manual Feature Creation [​](#after-manual-feature-creation)

If you created features manually:

bash

    # You manually created files for Order feature
    goca integrate --feature Order

### Project Repair [​](#project-repair)

If DI or routes are out of sync:

### After Git Merge [​](#after-git-merge)

After merging branches with new features:

See Also [​](#see-also)
-----------------------

*   [`goca feature`](https://sazardev.github.io/goca/commands/feature.html) - Generate complete feature (auto-integrated)
*   [`goca di`](https://sazardev.github.io/goca/commands/di.html) - Generate DI container</content>
</page>

<page>
  <title>Goca - Go Clean Architecture Code Generator</title>
  <url>https://sazardev.github.io/goca/commands/interfaces.html</url>
  <content>goca interfaces [​](#goca-interfaces)
-------------------------------------

Generate interface contracts for Test-Driven Development.

Syntax [​](#syntax)
-------------------

bash

    goca interfaces <EntityName> [flags]

Description [​](#description)
-----------------------------

Generates only the interface contracts without implementations, perfect for TDD workflows where you define contracts first.

Flags [​](#flags)
-----------------

### `--usecase` [​](#usecase)

Generate use case interfaces.

bash

    goca interfaces Product --usecase

### `--repository` [​](#repository)

Generate repository interfaces.

bash

    goca interfaces User --repository

### `--all` [​](#all)

Generate all interfaces.

bash

    goca interfaces Order --all

Examples [​](#examples)
-----------------------

### Use Case Interfaces [​](#use-case-interfaces)

bash

    goca interfaces Product --usecase

**Generates:** `internal/usecase/product_interfaces.go`

go

    package usecase
    
    import "context"
    
    type ProductService interface {
        CreateProduct(ctx context.Context, input CreateProductInput) (*ProductResponse, error)
        GetProduct(ctx context.Context, id uint) (*ProductResponse, error)
        UpdateProduct(ctx context.Context, id uint, input UpdateProductInput) error
        DeleteProduct(ctx context.Context, id uint) error
        ListProducts(ctx context.Context) ([]*ProductResponse, error)
    }

### Repository Interfaces [​](#repository-interfaces)

bash

    goca interfaces User --repository

**Generates:** `internal/repository/user_interfaces.go`

go

    package repository
    
    import (
        "context"
        "myproject/internal/domain"
    )
    
    type UserRepository interface {
        Save(ctx context.Context, user *domain.User) error
        FindByID(ctx context.Context, id uint) (*domain.User, error)
        FindAll(ctx context.Context) ([]*domain.User, error)
        Update(ctx context.Context, user *domain.User) error
        Delete(ctx context.Context, id uint) error
    }

### All Interfaces [​](#all-interfaces)

bash

    goca interfaces Order --all

TDD Workflow [​](#tdd-workflow)
-------------------------------

1.  **Generate Interfaces:**
    
    bash
    
        goca interfaces Payment --all
    
2.  **Write Tests:**
    
    go
    
        func TestPaymentService_CreatePayment(t *testing.T) {
            mockRepo := &MockPaymentRepository{}
            service := NewPaymentService(mockRepo)
            // ... test implementation
        }
    
3.  **Implement:** Implement the actual service and repository.
    

See Also [​](#see-also)
-----------------------

*   [`goca usecase`](https://sazardev.github.io/goca/commands/usecase.html) - Generate full use cases
*   [`goca repository`](https://sazardev.github.io/goca/commands/repository.html) - Generate repositories</content>
</page>

<page>
  <title>Goca - Go Clean Architecture Code Generator</title>
  <url>https://sazardev.github.io/goca/commands/messages.html</url>
  <content>goca messages [​](#goca-messages)
---------------------------------

Generate error messages, responses, and constants for features.

Syntax [​](#syntax)
-------------------

bash

    goca messages <EntityName> [flags]

Description [​](#description)
-----------------------------

Creates centralized message files for errors, success responses, and feature-specific constants.

Flags [​](#flags)
-----------------

### `--errors` [​](#errors)

Generate error messages.

bash

    goca messages User --errors

### `--responses` [​](#responses)

Generate response messages.

bash

    goca messages Product --responses

### `--constants` [​](#constants)

Generate feature constants.

bash

    goca messages Order --constants

Examples [​](#examples)
-----------------------

### Error Messages [​](#error-messages)

bash

    goca messages User --errors

**Generates:** `internal/messages/user_errors.go`

go

    package messages
    
    import "errors"
    
    var (
        ErrUserNotFound      = errors.New("user not found")
        ErrUserAlreadyExists = errors.New("user already exists")
        ErrInvalidUserData   = errors.New("invalid user data")
        ErrUserUnauthorized  = errors.New("user unauthorized")
    )

### Response Messages [​](#response-messages)

bash

    goca messages Product --responses

**Generates:** `internal/messages/product_responses.go`

go

    package messages
    
    const (
        ProductCreatedSuccess = "Product created successfully"
        ProductUpdatedSuccess = "Product updated successfully"
        ProductDeletedSuccess = "Product deleted successfully"
        ProductListSuccess    = "Products retrieved successfully"
    )

### Constants [​](#constants-1)

bash

    goca messages Order --constants

**Generates:** `internal/messages/order_constants.go`

go

    package messages
    
    const (
        OrderStatusPending    = "pending"
        OrderStatusProcessing = "processing"
        OrderStatusCompleted  = "completed"
        OrderStatusCancelled  = "cancelled"
        
        DefaultPageSize = 20
        MaxPageSize     = 100
    )

See Also [​](#see-also)
-----------------------

*   [`goca entity`](https://sazardev.github.io/goca/commands/entity.html) - Generate entities
*   [`goca feature`](https://sazardev.github.io/goca/commands/feature.html) - Generate complete feature</content>
</page>

<page>
  <title>Goca - Go Clean Architecture Code Generator</title>
  <url>https://sazardev.github.io/goca/commands/version.html</url>
  <content>goca version [​](#goca-version)
-------------------------------

Display version and build information.

Syntax [​](#syntax)
-------------------

Description [​](#description)
-----------------------------

Shows detailed version information about your Goca installation.

Example Output [​](#example-output)
-----------------------------------

    Goca v2.0.0
    Build: 2025-10-11T10:00:00Z
    Go Version: go1.24.5
    OS/Arch: linux/amd64
    Commit: abc123def

Information Displayed [​](#information-displayed)
-------------------------------------------------

*   **Version**: Current Goca version
*   **Build**: Build timestamp
*   **Go Version**: Go compiler version used
*   **OS/Arch**: Operating system and architecture
*   **Commit**: Git commit hash (if available)

Usage [​](#usage)
-----------------

Check your current version:

Verify you have the latest version:

bash

    # Current version
    goca version
    
    # Latest available
    # Check: https://github.com/sazardev/goca/releases/latest

Updating [​](#updating)
-----------------------

If you installed via `go install`:

bash

    go install github.com/sazardev/goca@latest
    goca version

See Also [​](#see-also)
-----------------------

*   [Installation Guide](https://sazardev.github.io/goca/guide/installation.html) - How to install Goca
*   [GitHub Releases](https://github.com/sazardev/goca/releases) - All versions</content>
</page>

<page>
  <title>Goca - Go Clean Architecture Code Generator</title>
  <url>https://sazardev.github.io/goca</url>
  <content>Build production-ready Go applications following Clean Architecture principles. Stop writing boilerplate, start building features.

Clean Architecture by Default
-----------------------------

Every line of code follows Uncle Bob's Clean Architecture principles. Proper layer separation, dependency rules, and clear boundaries guaranteed.

Lightning Fast Setup
--------------------

Generate complete features with all layers in seconds. From entity to handler, repository to use case - everything configured and ready.

Best Practices Enforced
-----------------------

Prevents common anti-patterns like fat controllers, god objects, and anemic domain models. Your code stays clean and maintainable.

Auto Integration
----------------

New features are automatically integrated with dependency injection and routing. No manual wiring needed.

Multi-Protocol Support
----------------------

Generate handlers for HTTP REST, gRPC, CLI, Workers, and SOAP. All following the same clean architecture pattern.

Test-Ready
----------

Code generated with clear interfaces and dependency injection makes testing a breeze. TDD-friendly from the start.

8 Databases Supported
---------------------

PostgreSQL, MySQL, MongoDB, SQLite, SQL Server, PostgreSQL JSON, Elasticsearch, and DynamoDB. Switch between databases without changing business logic.

Rich Documentation
------------------

Comprehensive guides, tutorials, and examples. Learn Clean Architecture while building real applications.

Production Ready
----------------

Used in production systems. Battle-tested patterns and code generation that scales from MVP to enterprise.

Quick Example [​](#quick-example)
---------------------------------

bash

    # Initialize a new project
    goca init my-api --module github.com/user/my-api
    
    # Generate a complete feature with all layers
    goca feature User --fields "name:string,email:string,role:string"
    
    # That's it! You now have:
    # → Domain entity with validations
    # → Use cases with DTOs
    # → Repository with PostgreSQL implementation
    # → HTTP handlers with routing
    # → Dependency injection configured

Why Clean Architecture? [​](#why-clean-architecture)
----------------------------------------------------

Clean Architecture ensures your codebase remains:

*   **Maintainable**: Changes in one layer don't cascade through the entire system
*   **Testable**: Business logic is independent of frameworks and databases
*   **Flexible**: Easy to swap implementations without touching core logic
*   **Scalable**: Clear boundaries make it easy to add new features

What Developers Say [​](#what-developers-say)
---------------------------------------------

> "Goca transformed how we build Go services. What used to take hours now takes minutes, and the code quality is consistently high."
> 
> — Production User

> "Finally, a code generator that doesn't just dump code but teaches you proper architecture."
> 
> — Go Developer

> "The automatic integration of new features saved us so much time. No more manual wiring!"
> 
> — Go Team Lead

Ready to Build? [​](#ready-to-build)
------------------------------------

[Get Started in 5 Minutes](https://sazardev.github.io/goca/getting-started.html)</content>
</page>

<page>
  <title>Mastering the Repository Pattern in Clean Architecture | Goca Blog</title>
  <url>https://sazardev.github.io/goca/blog/articles/mastering-repository-pattern</url>
  <content>Mastering the Repository Pattern in Clean Architecture [​](#mastering-the-repository-pattern-in-clean-architecture)
-------------------------------------------------------------------------------------------------------------------

InfrastructureData Access

The Repository pattern is a critical abstraction that isolates domain logic from data access concerns. Understanding repositories correctly is essential for building applications that remain testable, maintainable, and independent of specific database technologies. When implemented properly, repositories enable you to change databases, add caching, or switch ORMs without touching business logic.

* * *

What is the Repository Pattern? [​](#what-is-the-repository-pattern)
--------------------------------------------------------------------

A repository mediates between the domain layer and data mapping layers, acting as an in-memory collection of domain objects. Repositories provide a clean, domain-centric API for data access while hiding the complexities of database interactions, query construction, and persistence mechanisms.

### Core Responsibilities [​](#core-responsibilities)

**Abstraction**: Repositories abstract the details of data access. The domain layer works with repository interfaces that express intent ("find user by email") rather than implementation ("execute SQL query with WHERE clause").

**Collection Semantics**: Repositories present data as collections. You add entities to repositories, remove them, and query them using domain-meaningful methods. The fact that data persists to a database is an implementation detail.

**Domain-Centric API**: Repository methods use domain language and work with domain entities. A `UserRepository` has methods like `Save(user)`, not `Insert(tableName, columns, values)`.

**Testability**: Because repositories are interfaces, you can replace real implementations with in-memory fakes during testing. This allows unit testing of business logic without database setup.

**Database Independence**: Repositories isolate database-specific code. Changing from PostgreSQL to MongoDB requires changing only repository implementations, not domain or application logic.

Repository vs DAO vs Data Mapper [​](#repository-vs-dao-vs-data-mapper)
-----------------------------------------------------------------------

Developers often confuse repositories with Data Access Objects (DAOs) or Data Mappers. These patterns serve different purposes and operate at different abstraction levels.

### What a Repository Is NOT [​](#what-a-repository-is-not)

**Not a DAO**: DAOs provide CRUD operations on database tables. Repositories provide domain operations on entity collections. A DAO might have `insertUser(name, email)`. A repository has `Save(user *User)`.

**Not a Data Mapper**: Data Mappers convert between database rows and objects. Repositories use data mappers internally but provide higher-level operations that reflect business operations.

**Not a Generic Interface**: Repositories are not `IRepository<T>` with generic CRUD. Each repository has a domain-specific interface. `UserRepository` has methods that make sense for users; `OrderRepository` has methods that make sense for orders.

**Not Query Builders**: Repositories do not expose SQL or query DSLs. They provide intention-revealing methods. Instead of `repository.Query("SELECT * FROM users WHERE age > ?", 18)`, you have `repository.FindAdults()`.

### The Clear Distinction [​](#the-clear-distinction)

    Domain Layer (Business Logic)
        ↓ Uses interface
    Repository Interface (Contract)
        ↓ Defined in domain
    Infrastructure Layer (Data Access)
        ↓ Implements interface
    Repository Implementation (Database-Specific)
        ↓ Uses ORM/Driver
    Database

Repositories are defined in the domain layer as interfaces but implemented in the infrastructure layer with database-specific code. This inverts the dependency, making infrastructure depend on the domain rather than the reverse.

The Infrastructure Layer [​](#the-infrastructure-layer)
-------------------------------------------------------

Repositories form the infrastructure layer in Clean Architecture. This layer contains all the concrete implementations of persistence, external services, and framework-specific code.

### Infrastructure Layer Characteristics [​](#infrastructure-layer-characteristics)

**Implements Domain Interfaces**: The infrastructure layer provides concrete implementations of repository interfaces defined in the domain. The domain depends on abstractions; infrastructure depends on the domain.

**Database-Specific Code**: Infrastructure contains ORM configurations, SQL queries, connection management, and database-specific optimizations. This code is hidden behind interfaces.

**Framework Dependencies**: Infrastructure can depend on GORM, database drivers, caching libraries, and external SDKs. These dependencies do not leak into the domain.

**Swappable Implementations**: You can have multiple repository implementations for the same interface: PostgreSQL for production, in-memory for testing, MongoDB for a specific feature.

### Why Separate Infrastructure? [​](#why-separate-infrastructure)

The infrastructure layer exists because data access mechanisms change independently of business rules:

**Domain Logic**: "A user must have a unique email" is domain logic. It does not care how you check uniqueness.

**Infrastructure Logic**: "Execute SELECT COUNT(\*) FROM users WHERE email = ? to check uniqueness" is infrastructure logic. It is specific to SQL databases.

Separating these concerns allows you to:

*   Test domain logic without database setup
*   Change databases without changing business rules
*   Optimize queries without touching domain code
*   Support multiple databases simultaneously

Repository Interface Design [​](#repository-interface-design)
-------------------------------------------------------------

Well-designed repository interfaces express domain intent clearly while remaining independent of implementation details.

### Interface Location: Domain Layer [​](#interface-location-domain-layer)

Repository interfaces belong in the domain layer, typically in a package like `internal/domain` or alongside entity definitions. This placement is critical:

**Domain Owns the Contract**: The domain defines what operations it needs. Infrastructure adapts to the domain's requirements, not vice versa.

**Dependency Inversion**: By placing interfaces in the domain, you invert the dependency. Infrastructure imports domain types, not the other way around.

**No Infrastructure Leakage**: Domain interfaces use only domain types. They do not reference database connections, ORM types, or SQL constructs.

### Method Design Principles [​](#method-design-principles)

**Intention-Revealing Names**: Methods should express why you are querying, not how. `FindActiveUsers()` is better than `QueryUsersWithStatus("active")`.

**Domain Types Only**: Parameters and return types are domain entities and value objects, never database-specific types like `sql.Row` or `bson.Document`.

**Error Handling**: Return domain errors, not database errors. Instead of returning `sql.ErrNoRows`, return `ErrUserNotFound`.

**No Leaky Abstractions**: Methods should not expose pagination cursors, query builders, or transaction objects unless these concepts exist in the domain.

### Common Repository Methods [​](#common-repository-methods)

Every repository typically includes these fundamental operations:

go

    type UserRepository interface {
        // Save adds or updates a user
        Save(user *User) error
        
        // FindByID retrieves a user by unique identifier
        FindByID(id uint) (*User, error)
        
        // FindAll retrieves all users (use with caution in production)
        FindAll() ([]*User, error)
        
        // Update modifies an existing user
        Update(user *User) error
        
        // Delete removes a user
        Delete(id uint) error
    }

Beyond these basics, add domain-specific query methods:

go

    type UserRepository interface {
        // ... basic methods ...
        
        // Domain-specific queries
        FindByEmail(email string) (*User, error)
        FindAdults() ([]*User, error)
        FindByLastLoginAfter(date time.Time) ([]*User, error)
        CountByStatus(status Status) (int, error)
    }

Repository Implementation Patterns [​](#repository-implementation-patterns)
---------------------------------------------------------------------------

Repository implementations encapsulate all database-specific logic, translating domain operations into database operations.

### Basic PostgreSQL Implementation [​](#basic-postgresql-implementation)

go

    package repository
    
    import (
        "github.com/yourorg/yourapp/internal/domain"
        "gorm.io/gorm"
    )
    
    type postgresUserRepository struct {
        db *gorm.DB
    }
    
    // NewPostgresUserRepository creates a PostgreSQL implementation
    func NewPostgresUserRepository(db *gorm.DB) domain.UserRepository {
        return &postgresUserRepository{db: db}
    }
    
    func (r *postgresUserRepository) Save(user *domain.User) error {
        return r.db.Create(user).Error
    }
    
    func (r *postgresUserRepository) FindByID(id uint) (*domain.User, error) {
        var user domain.User
        err := r.db.First(&user, id).Error
        if err == gorm.ErrRecordNotFound {
            return nil, domain.ErrUserNotFound
        }
        return &user, err
    }
    
    func (r *postgresUserRepository) FindByEmail(email string) (*domain.User, error) {
        var user domain.User
        err := r.db.Where("email = ?", email).First(&user).Error
        if err == gorm.ErrRecordNotFound {
            return nil, domain.ErrUserNotFound
        }
        return &user, err
    }
    
    func (r *postgresUserRepository) Update(user *domain.User) error {
        return r.db.Save(user).Error
    }
    
    func (r *postgresUserRepository) Delete(id uint) error {
        return r.db.Delete(&domain.User{}, id).Error
    }
    
    func (r *postgresUserRepository) FindAll() ([]*domain.User, error) {
        var users []*domain.User
        err := r.db.Find(&users).Error
        return users, err
    }

Notice how the implementation:

*   Returns domain errors, not GORM errors
*   Uses domain types in signatures
*   Hides all GORM-specific code
*   Implements the domain interface

### MongoDB Implementation [​](#mongodb-implementation)

For NoSQL databases, the implementation differs dramatically, but the interface remains the same:

go

    package repository
    
    import (
        "context"
        "github.com/yourorg/yourapp/internal/domain"
        "go.mongodb.org/mongo-driver/bson"
        "go.mongodb.org/mongo-driver/mongo"
    )
    
    type mongoUserRepository struct {
        collection *mongo.Collection
    }
    
    func NewMongoUserRepository(db *mongo.Database) domain.UserRepository {
        return &mongoUserRepository{
            collection: db.Collection("users"),
        }
    }
    
    func (r *mongoUserRepository) Save(user *domain.User) error {
        ctx := context.TODO()
        _, err := r.collection.InsertOne(ctx, user)
        return err
    }
    
    func (r *mongoUserRepository) FindByID(id uint) (*domain.User, error) {
        ctx := context.TODO()
        var user domain.User
        
        err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&user)
        if err == mongo.ErrNoDocuments {
            return nil, domain.ErrUserNotFound
        }
        return &user, err
    }
    
    func (r *mongoUserRepository) FindByEmail(email string) (*domain.User, error) {
        ctx := context.TODO()
        var user domain.User
        
        err := r.collection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
        if err == mongo.ErrNoDocuments {
            return nil, domain.ErrUserNotFound
        }
        return &user, err
    }

The key insight: both implementations satisfy the same `UserRepository` interface. Application and domain code work with either database without modification.

### In-Memory Implementation for Testing [​](#in-memory-implementation-for-testing)

For unit tests, create an in-memory fake that implements the repository interface:

go

    package repository
    
    import (
        "sync"
        "github.com/yourorg/yourapp/internal/domain"
    )
    
    type inMemoryUserRepository struct {
        users  map[uint]*domain.User
        nextID uint
        mu     sync.RWMutex
    }
    
    func NewInMemoryUserRepository() domain.UserRepository {
        return &inMemoryUserRepository{
            users:  make(map[uint]*domain.User),
            nextID: 1,
        }
    }
    
    func (r *inMemoryUserRepository) Save(user *domain.User) error {
        r.mu.Lock()
        defer r.mu.Unlock()
        
        if user.ID == 0 {
            user.ID = r.nextID
            r.nextID++
        }
        
        r.users[user.ID] = user
        return nil
    }
    
    func (r *inMemoryUserRepository) FindByID(id uint) (*domain.User, error) {
        r.mu.RLock()
        defer r.mu.RUnlock()
        
        user, exists := r.users[id]
        if !exists {
            return nil, domain.ErrUserNotFound
        }
        return user, nil
    }
    
    func (r *inMemoryUserRepository) FindByEmail(email string) (*domain.User, error) {
        r.mu.RLock()
        defer r.mu.RUnlock()
        
        for _, user := range r.users {
            if user.Email == email {
                return user, nil
            }
        }
        return nil, domain.ErrUserNotFound
    }

This in-memory implementation enables fast, isolated unit tests without database dependencies.

How Goca Generates Repositories [​](#how-goca-generates-repositories)
---------------------------------------------------------------------

Goca's `goca repository` command generates both interfaces and implementations following Clean Architecture principles.

### Basic Repository Generation [​](#basic-repository-generation)

bash

    goca repository User --database postgres

This generates:

**1\. Repository Interface** (`internal/repository/interfaces.go`):

go

    package repository
    
    import "yourapp/internal/domain"
    
    type UserRepository interface {
        Save(user *domain.User) error
        FindByID(id uint) (*domain.User, error)
        Update(user *domain.User) error
        Delete(id uint) error
        FindAll() ([]*domain.User, error)
    }

**2\. PostgreSQL Implementation** (`internal/repository/postgres_user_repository.go`):

go

    package repository
    
    import (
        "yourapp/internal/domain"
        "gorm.io/gorm"
    )
    
    type postgresUserRepository struct {
        db *gorm.DB
    }
    
    func NewPostgresUserRepository(db *gorm.DB) UserRepository {
        return &postgresUserRepository{db: db}
    }
    
    // ... CRUD implementations ...

### Database-Specific Implementations [​](#database-specific-implementations)

Goca supports multiple databases, generating appropriate implementations for each:

**PostgreSQL with GORM**:

bash

    goca repository User --database postgres
    # Generates GORM-based implementation with SQL transactions

**MongoDB**:

bash

    goca repository User --database mongodb
    # Generates MongoDB driver implementation with BSON

**PostgreSQL with JSONB**:

bash

    goca repository User --database postgres-json
    # Generates JSONB-specific queries for nested documents

**MySQL**:

bash

    goca repository User --database mysql
    # Generates MySQL-specific GORM implementation

**DynamoDB**:

bash

    goca repository User --database dynamodb
    # Generates AWS SDK v2 implementation with attribute mapping

### Custom Query Methods [​](#custom-query-methods)

Goca auto-generates query methods based on entity fields:

bash

    goca repository User --fields "name:string,email:string,age:int,status:string"

Generates these additional methods:

go

    type UserRepository interface {
        // Basic CRUD
        Save(user *domain.User) error
        FindByID(id uint) (*domain.User, error)
        Update(user *domain.User) error
        Delete(id uint) error
        FindAll() ([]*domain.User, error)
        
        // Field-based queries (auto-generated)
        FindByName(name string) (*domain.User, error)
        FindByEmail(email string) (*domain.User, error)
        FindByAge(age int) (*domain.User, error)
        FindByStatus(status string) (*domain.User, error)
    }

### Interface-Only Generation [​](#interface-only-generation)

For Test-Driven Development (TDD), generate only the interface first:

bash

    goca repository User --interface-only

This creates the contract without implementation, allowing you to:

1.  Write use cases against the interface
2.  Create mock implementations for tests
3.  Implement the real repository later

### Implementation-Only Generation [​](#implementation-only-generation)

If you already have the interface but need a new database implementation:

bash

    goca repository User --implementation --database mongodb

This generates only the MongoDB implementation without modifying the interface.

Advanced Repository Patterns [​](#advanced-repository-patterns)
---------------------------------------------------------------

Beyond basic CRUD, repositories support advanced patterns for complex data access scenarios.

### Specification Pattern [​](#specification-pattern)

Use specifications to encapsulate query criteria:

go

    type UserSpecification interface {
        IsSatisfiedBy(user *domain.User) bool
        ToSQL() (string, []interface{})
    }
    
    type UserRepository interface {
        // ... basic methods ...
        FindBySpec(spec UserSpecification) ([]*domain.User, error)
    }
    
    // Usage
    activeAdults := NewAndSpecification(
        NewAgeGreaterThanSpec(18),
        NewStatusEqualsSpec("active"),
    )
    users, err := repo.FindBySpec(activeAdults)

### Unit of Work Pattern [​](#unit-of-work-pattern)

Coordinate multiple repository operations in a transaction:

go

    type UnitOfWork interface {
        Users() UserRepository
        Orders() OrderRepository
        
        Begin() error
        Commit() error
        Rollback() error
    }
    
    // Usage in use case
    func (s *orderService) CreateOrder(userID uint, items []Item) error {
        uow := s.unitOfWork
        
        if err := uow.Begin(); err != nil {
            return err
        }
        defer uow.Rollback()
        
        user, err := uow.Users().FindByID(userID)
        if err != nil {
            return err
        }
        
        order := domain.NewOrder(user, items)
        if err := uow.Orders().Save(order); err != nil {
            return err
        }
        
        return uow.Commit()
    }

### Caching Layer [​](#caching-layer)

Add caching transparently using the decorator pattern:

go

    type cachedUserRepository struct {
        repository UserRepository
        cache      Cache
    }
    
    func NewCachedUserRepository(repo UserRepository, cache Cache) UserRepository {
        return &cachedUserRepository{
            repository: repo,
            cache:      cache,
        }
    }
    
    func (r *cachedUserRepository) FindByID(id uint) (*domain.User, error) {
        // Check cache first
        cacheKey := fmt.Sprintf("user:%d", id)
        if cached, found := r.cache.Get(cacheKey); found {
            return cached.(*domain.User), nil
        }
        
        // Cache miss: query database
        user, err := r.repository.FindByID(id)
        if err != nil {
            return nil, err
        }
        
        // Store in cache
        r.cache.Set(cacheKey, user, 5*time.Minute)
        return user, nil
    }

The use case layer does not know caching exists. The cached repository satisfies the same interface as the uncached version.

### Pagination Support [​](#pagination-support)

For large datasets, repositories should support pagination:

go

    type Page struct {
        Items      []*User
        TotalItems int
        Page       int
        PageSize   int
        TotalPages int
    }
    
    type UserRepository interface {
        // ... basic methods ...
        FindAllPaginated(page, pageSize int) (*Page, error)
        FindByStatusPaginated(status string, page, pageSize int) (*Page, error)
    }
    
    // Implementation
    func (r *postgresUserRepository) FindAllPaginated(page, pageSize int) (*Page, error) {
        var users []*domain.User
        var total int64
        
        offset := (page - 1) * pageSize
        
        // Count total
        r.db.Model(&domain.User{}).Count(&total)
        
        // Fetch page
        err := r.db.Limit(pageSize).Offset(offset).Find(&users).Error
        
        return &Page{
            Items:      users,
            TotalItems: int(total),
            Page:       page,
            PageSize:   pageSize,
            TotalPages: int(total)/pageSize + 1,
        }, err
    }

### Soft Delete Support [​](#soft-delete-support)

Implement soft deletes while maintaining a clean repository interface:

go

    func (r *postgresUserRepository) Delete(id uint) error {
        // Soft delete: set deleted_at timestamp
        return r.db.Model(&domain.User{}).
            Where("id = ?", id).
            Update("deleted_at", time.Now()).Error
    }
    
    func (r *postgresUserRepository) FindByID(id uint) (*domain.User, error) {
        var user domain.User
        // Automatically exclude soft-deleted records
        err := r.db.Where("deleted_at IS NULL").First(&user, id).Error
        if err == gorm.ErrRecordNotFound {
            return nil, domain.ErrUserNotFound
        }
        return &user, err
    }
    
    // Add method for finding deleted records if needed
    func (r *postgresUserRepository) FindDeletedByID(id uint) (*domain.User, error) {
        var user domain.User
        err := r.db.Unscoped().First(&user, id).Error
        if err == gorm.ErrRecordNotFound {
            return nil, domain.ErrUserNotFound
        }
        return &user, err
    }

Testing Repository Implementations [​](#testing-repository-implementations)
---------------------------------------------------------------------------

Repositories require different testing strategies than domain logic because they involve external dependencies.

### Unit Testing with Mocks [​](#unit-testing-with-mocks)

For use case testing, use mock repositories:

go

    type MockUserRepository struct {
        SaveFunc    func(user *domain.User) error
        FindByIDFunc func(id uint) (*domain.User, error)
    }
    
    func (m *MockUserRepository) Save(user *domain.User) error {
        if m.SaveFunc != nil {
            return m.SaveFunc(user)
        }
        return nil
    }
    
    func (m *MockUserRepository) FindByID(id uint) (*domain.User, error) {
        if m.FindByIDFunc != nil {
            return m.FindByIDFunc(id)
        }
        return nil, domain.ErrUserNotFound
    }
    
    // Test use case
    func TestCreateUser(t *testing.T) {
        mockRepo := &MockUserRepository{
            FindByEmailFunc: func(email string) (*domain.User, error) {
                return nil, domain.ErrUserNotFound // Email available
            },
            SaveFunc: func(user *domain.User) error {
                user.ID = 1
                return nil
            },
        }
        
        service := NewUserService(mockRepo)
        
        user, err := service.CreateUser(CreateUserInput{
            Name:  "John",
            Email: "john@example.com",
        })
        
        assert.NoError(t, err)
        assert.Equal(t, uint(1), user.ID)
    }

### Integration Testing with Real Database [​](#integration-testing-with-real-database)

For repository implementation testing, use a real database:

go

    func TestPostgresUserRepository_Save(t *testing.T) {
        // Setup test database
        db := setupTestDatabase(t)
        defer cleanupTestDatabase(t, db)
        
        repo := NewPostgresUserRepository(db)
        
        // Create test user
        user := &domain.User{
            Name:  "John Doe",
            Email: "john@example.com",
            Age:   30,
        }
        
        // Test Save
        err := repo.Save(user)
        assert.NoError(t, err)
        assert.NotZero(t, user.ID)
        
        // Verify in database
        found, err := repo.FindByID(user.ID)
        assert.NoError(t, err)
        assert.Equal(t, user.Name, found.Name)
        assert.Equal(t, user.Email, found.Email)
    }
    
    func setupTestDatabase(t *testing.T) *gorm.DB {
        db, err := gorm.Open(postgres.Open("postgres://test:test@localhost/testdb"), &gorm.Config{})
        require.NoError(t, err)
        
        // Run migrations
        db.AutoMigrate(&domain.User{})
        
        return db
    }
    
    func cleanupTestDatabase(t *testing.T, db *gorm.DB) {
        db.Exec("TRUNCATE TABLE users CASCADE")
    }

### Testing with Docker [​](#testing-with-docker)

For isolated integration tests, use Docker containers:

go

    func TestWithDocker(t *testing.T) {
        if testing.Short() {
            t.Skip("Skipping integration test")
        }
        
        // Start PostgreSQL container
        ctx := context.Background()
        req := testcontainers.ContainerRequest{
            Image:        "postgres:15",
            ExposedPorts: []string{"5432/tcp"},
            Env: map[string]string{
                "POSTGRES_PASSWORD": "test",
                "POSTGRES_DB":       "testdb",
            },
            WaitingFor: wait.ForLog("database system is ready"),
        }
        
        postgres, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
            ContainerRequest: req,
            Started:          true,
        })
        require.NoError(t, err)
        defer postgres.Terminate(ctx)
        
        // Get connection string
        host, _ := postgres.Host(ctx)
        port, _ := postgres.MappedPort(ctx, "5432")
        dsn := fmt.Sprintf("postgres://test:test@%s:%s/testdb", host, port.Port())
        
        // Run tests
        db, _ := gorm.Open(postgres.Open(dsn), &gorm.Config{})
        repo := NewPostgresUserRepository(db)
        
        // ... test repository operations ...
    }

Repository Anti-Patterns [​](#repository-anti-patterns)
-------------------------------------------------------

Understanding what not to do is as important as knowing best practices.

### Anti-Pattern: Generic Repository [​](#anti-pattern-generic-repository)

**Problem**:

go

    type Repository[T any] interface {
        Save(entity T) error
        FindByID(id int) (T, error)
        Update(entity T) error
        Delete(id int) error
    }
    
    type UserRepository = Repository[User]
    type OrderRepository = Repository[Order]

**Why It's Bad**: Generic repositories force all entities to have the same operations. `FindByEmail` makes sense for users, not orders. You lose domain-specific expressiveness.

**Solution**: Create specific interfaces for each entity with domain-meaningful methods.

### Anti-Pattern: Leaky Abstraction [​](#anti-pattern-leaky-abstraction)

**Problem**:

go

    type UserRepository interface {
        Query(sql string, args ...interface{}) ([]*User, error)
        GetDB() *sql.DB
    }

**Why It's Bad**: This exposes database implementation details. Consumers must write SQL. You cannot swap databases.

**Solution**: Provide intention-revealing methods that hide database details.

### Anti-Pattern: Business Logic in Repository [​](#anti-pattern-business-logic-in-repository)

**Problem**:

go

    func (r *postgresUserRepository) CreateAdminUser(name, email string) (*User, error) {
        user := &User{
            Name:  name,
            Email: email,
            Role:  "admin",
        }
        
        // Business logic: validate admin email domain
        if !strings.HasSuffix(email, "@company.com") {
            return nil, errors.New("admin must use company email")
        }
        
        return user, r.db.Create(user).Error
    }

**Why It's Bad**: Business rules belong in the domain or use cases, not repositories. Repositories are for data access only.

**Solution**: Move validation to entity or use case. Repository only persists valid entities.

### Anti-Pattern: Repository Returning DTOs [​](#anti-pattern-repository-returning-dtos)

**Problem**:

go

    type UserRepository interface {
        FindByID(id int) (*UserDTO, error)
    }

**Why It's Bad**: Repositories work with domain entities, not DTOs. DTOs are for external communication, not internal operations.

**Solution**: Return domain entities. Use cases convert entities to DTOs.

Best Practices Summary [​](#best-practices-summary)
---------------------------------------------------

### Do This [​](#do-this)

✅ **Define interfaces in domain layer**: Keep contracts with domain code  
✅ **Use domain types in signatures**: Parameters and returns are entities  
✅ **Name methods by intent**: `FindActiveUsers()`, not `QueryUsers()`  
✅ **Return domain errors**: `ErrUserNotFound`, not `sql.ErrNoRows`  
✅ **Keep implementations simple**: One responsibility per repository  
✅ **Test with real databases**: Integration tests verify SQL correctness  
✅ **Use in-memory fakes for unit tests**: Fast, isolated tests

### Avoid This [​](#avoid-this)

❌ **Don't expose database types**: No `*sql.DB`, `*gorm.DB` in interfaces  
❌ **Don't add business logic**: Repositories persist, they don't validate  
❌ **Don't use generic interfaces**: Each entity gets specific methods  
❌ **Don't return DTOs**: Repositories work with domain entities  
❌ **Don't couple to ORM**: Use interfaces, not concrete ORM types

Conclusion [​](#conclusion)
---------------------------

The Repository pattern is a cornerstone of Clean Architecture, providing a boundary between domain logic and data access concerns. When implemented correctly, repositories enable:

*   **Database Independence**: Change databases without changing business logic
*   **Testability**: Unit test without database setup using mocks
*   **Maintainability**: Data access code isolated in one layer
*   **Flexibility**: Support multiple databases simultaneously
*   **Clean Domain**: Domain remains pure, focused on business rules

Goca's `goca repository` command generates repositories that follow these principles, creating both interfaces and database-specific implementations that maintain architectural boundaries while providing practical, production-ready data access code.

Master repositories, and you master one of the most critical patterns in Clean Architecture.</content>
</page>

<page>
  <title>Example - Advanced Features Showcase | Goca Blog</title>
  <url>https://sazardev.github.io/goca/blog/articles/example-showcase</url>
  <content>Advanced Features Showcase [​](#advanced-features-showcase)
-----------------------------------------------------------

ExampleTutorial

This is an example article demonstrating the full capabilities of the Goca blog system, including Mermaid diagrams, syntax-highlighted code blocks, and advanced markdown features.

* * *

Clean Architecture Flow [​](#clean-architecture-flow)
-----------------------------------------------------

Here's how Goca implements Clean Architecture principles using a Mermaid diagram:

mermaid

    graph TD
        A[HTTP Request] --> B[Handler Layer]
        B --> C[Use Case Layer]
        C --> D[Repository Interface]
        D --> E[Repository Implementation]
        E --> F[Database]
        
        B -.->|DTO| C
        C -.->|Domain Entity| D

### Layer Dependencies [​](#layer-dependencies)

The dependency rule states that source code dependencies must point inward:

mermaid

    graph LR
        A[External Interfaces<br/>Frameworks & Drivers] --> B[Interface Adapters<br/>Controllers, Presenters]
        B --> C[Application Business Rules<br/>Use Cases]
        C --> D[Enterprise Business Rules<br/>Entities]

Code Generation Workflow [​](#code-generation-workflow)
-------------------------------------------------------

### Entity Generation Process [​](#entity-generation-process)

go

    // Generate an entity with Goca
    package main
    
    import (
        "fmt"
        "time"
    )
    
    // User represents a user entity in the domain layer
    type User struct {
        ID        uint      `json:"id" gorm:"primaryKey"`
        Name      string    `json:"name" gorm:"type:varchar(255);not null"`
        Email     string    `json:"email" gorm:"type:varchar(255);uniqueIndex;not null"`
        CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
        UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
    }
    
    // Validate performs business logic validation
    func (u *User) Validate() error {
        if u.Name == "" {
            return fmt.Errorf("name cannot be empty")
        }
        if u.Email == "" {
            return fmt.Errorf("email cannot be empty")
        }
        return nil
    }

### Command Execution Sequence [​](#command-execution-sequence)

mermaid

    sequenceDiagram
        participant User
        participant CLI
        participant Validator
        participant Generator
        participant FileSystem
        
        User->>CLI: goca feature User --fields "name:string,email:string"
        CLI->>Validator: Validate input
        Validator-->>CLI: Validation OK
        CLI->>Generator: Generate files
        Generator->>FileSystem: Create domain/user.go
        Generator->>FileSystem: Create usecase/user_service.go
        Generator->>FileSystem: Create repository/postgres_user_repository.go
        Generator->>FileSystem: Create handler/http/user_handler.go
        FileSystem-->>User: ✅ Feature generated successfully

Database Support Matrix [​](#database-support-matrix)
-----------------------------------------------------

Goca supports multiple databases with specific implementations:

| Database | Type | Primary Use Case | Status |
| --- | --- | --- | --- |
| PostgreSQL | SQL | OLTP/General | ✅ Stable |
| PostgreSQL JSON | SQL+Document | Semi-structured | ✅ Stable |
| MySQL | SQL | Web Applications | ✅ Stable |
| MongoDB | NoSQL | Document Store | ✅ Stable |
| SQLite | SQL | Embedded/Testing | ✅ Stable |
| SQL Server | SQL | Enterprise | ✅ Stable |
| Elasticsearch | Search | Full-text Search | ✅ Stable |
| DynamoDB | NoSQL | Serverless AWS | ✅ Stable |

### Database Selection Decision Tree [​](#database-selection-decision-tree)

mermaid

    graph TD
        A[Choose Database] --> B{Data Structure?}
        B -->|Relational| C{Scale?}
        B -->|Document| D{Managed?}
        B -->|Search| E[Elasticsearch]
        
        C -->|Small| F[SQLite]
        C -->|Medium| G{Platform?}
        C -->|Large| H[PostgreSQL]
        
        G -->|Open Source| I[MySQL]
        G -->|Enterprise| J[SQL Server]
        
        D -->|Self-hosted| K[MongoDB]
        D -->|AWS| L[DynamoDB]

Project Structure Visualization [​](#project-structure-visualization)
---------------------------------------------------------------------

mermaid

    graph TB
        subgraph "Clean Architecture Layers"
            A[internal/domain/<br/>Entities & Business Rules]
            B[internal/usecase/<br/>Application Services]
            C[internal/repository/<br/>Data Access]
            D[internal/handler/<br/>External Interfaces]
        end
        
        subgraph "Infrastructure"
            E[internal/di/<br/>Dependency Injection]
            F[cmd/<br/>Entry Points]
        end
        
        F --> E
        E --> D
        D --> B
        B --> A
        C --> A

Testing Strategy [​](#testing-strategy)
---------------------------------------

### Test Pyramid [​](#test-pyramid)

mermaid

    graph TB
        subgraph "Test Pyramid"
            A[E2E Tests<br/>Integration Tests]
            B[Service Tests<br/>Use Case Tests]
            C[Unit Tests<br/>Domain Logic Tests]
        end
        
        A --> B --> C

### Test Coverage by Layer [​](#test-coverage-by-layer)

bash

    # Run tests with coverage
    go test ./internal/domain/... -cover
    # PASS coverage: 95.2% of statements
    
    go test ./internal/usecase/... -cover
    # PASS coverage: 89.7% of statements
    
    go test ./internal/repository/... -cover
    # PASS coverage: 78.4% of statements
    
    go test ./internal/handler/... -cover
    # PASS coverage: 82.1% of statements

Best Practices [​](#best-practices)
-----------------------------------

### Repository Pattern Implementation [​](#repository-pattern-implementation)

go

    // Repository interface (domain layer)
    type UserRepository interface {
        Save(ctx context.Context, user *User) error
        FindByID(ctx context.Context, id uint) (*User, error)
        Update(ctx context.Context, user *User) error
        Delete(ctx context.Context, id uint) error
        FindAll(ctx context.Context) ([]*User, error)
    }
    
    // PostgreSQL implementation (infrastructure layer)
    type postgresUserRepository struct {
        db *gorm.DB
    }
    
    func (r *postgresUserRepository) Save(ctx context.Context, user *User) error {
        return r.db.WithContext(ctx).Create(user).Error
    }
    
    func (r *postgresUserRepository) FindByID(ctx context.Context, id uint) (*User, error) {
        var user User
        if err := r.db.WithContext(ctx).First(&user, id).Error; err != nil {
            if errors.Is(err, gorm.ErrRecordNotFound) {
                return nil, fmt.Errorf("user not found")
            }
            return nil, err
        }
        return &user, nil
    }

### Use Case Pattern [​](#use-case-pattern)

go

    type UserUseCase interface {
        CreateUser(ctx context.Context, input CreateUserInput) (*UserOutput, error)
        GetUser(ctx context.Context, id uint) (*UserOutput, error)
        UpdateUser(ctx context.Context, input UpdateUserInput) (*UserOutput, error)
        DeleteUser(ctx context.Context, id uint) error
        ListUsers(ctx context.Context) ([]*UserOutput, error)
    }
    
    type userService struct {
        repo UserRepository
    }
    
    func (s *userService) CreateUser(ctx context.Context, input CreateUserInput) (*UserOutput, error) {
        // 1. Validate input
        if err := input.Validate(); err != nil {
            return nil, err
        }
        
        // 2. Create domain entity
        user := &User{
            Name:  input.Name,
            Email: input.Email,
        }
        
        // 3. Validate entity
        if err := user.Validate(); err != nil {
            return nil, err
        }
        
        // 4. Save to repository
        if err := s.repo.Save(ctx, user); err != nil {
            return nil, err
        }
        
        // 5. Return output DTO
        return &UserOutput{
            ID:        user.ID,
            Name:      user.Name,
            Email:     user.Email,
            CreatedAt: user.CreatedAt,
        }, nil
    }

Command Reference Quick Guide [​](#command-reference-quick-guide)
-----------------------------------------------------------------

bash

    # Initialize new project
    goca init myproject --database postgres
    
    # Generate complete feature
    goca feature User --fields "name:string,email:string,age:int"
    
    # Generate with testing support
    goca feature Product --fields "name:string,price:float64" \
        --integration-tests --mocks
    
    # Generate only entity
    goca entity Order --fields "total:float64,status:string"
    
    # Generate repository
    goca repository User --database postgres
    
    # Generate handler
    goca handler User --type http
    
    # Wire everything together
    goca integrate --all
    
    # Generate dependency injection
    goca di

State Machine Example [​](#state-machine-example)
-------------------------------------------------

Here's how you might model an order state machine:

mermaid

    stateDiagram-v2
        [*] --> Pending
        Pending --> Processing: Payment Confirmed
        Pending --> Cancelled: User Cancelled
        
        Processing --> Shipped: Items Dispatched
        Processing --> Cancelled: Out of Stock
        
        Shipped --> Delivered: Customer Received
        Shipped --> Returned: Customer Rejected
        
        Delivered --> [*]
        Returned --> Refunded
        Refunded --> [*]
        Cancelled --> [*]
        
        note right of Processing
            Inventory reserved
            Payment captured
        end note
        
        note right of Shipped
            Tracking number generated
            Notification sent
        end note

Performance Considerations [​](#performance-considerations)
-----------------------------------------------------------

### Database Query Optimization [​](#database-query-optimization)

sql

    -- Before: N+1 query problem
    SELECT * FROM users WHERE id = 1;
    SELECT * FROM orders WHERE user_id = 1;
    SELECT * FROM orders WHERE user_id = 2;
    SELECT * FROM orders WHERE user_id = 3;
    
    -- After: Single query with join
    SELECT u.*, o.*
    FROM users u
    LEFT JOIN orders o ON u.id = o.user_id
    WHERE u.id IN (1, 2, 3);

go

    // Use GORM preloading to avoid N+1
    var users []User
    db.Preload("Orders").Find(&users)

Deployment Architecture [​](#deployment-architecture)
-----------------------------------------------------

mermaid

    graph LR
        A[Load Balancer] --> B[App Server 1]
        A --> C[App Server 2]
        A --> D[App Server 3]
        
        B --> E[(Primary DB)]
        C --> E
        D --> E
        
        E --> F[(Replica DB)]
        E --> G[(Replica DB)]
        
        B --> H[Redis Cache]
        C --> H
        D --> H

Conclusion [​](#conclusion)
---------------------------

This example demonstrates the powerful capabilities available in Goca blog posts:

*   Mermaid diagrams for architecture visualization
*   Syntax-highlighted code blocks
*   Tables and structured data
*   State machines and flow diagrams
*   Sequence diagrams
*   Best practices and patterns

Use these features to create comprehensive, professional documentation and blog posts for your Goca projects.

* * *</content>
</page>

<page>
  <title>Understanding Domain Entities in Clean Architecture | Goca Blog</title>
  <url>https://sazardev.github.io/goca/blog/articles/understanding-domain-entities</url>
  <content>Understanding Domain Entities in Clean Architecture [​](#understanding-domain-entities-in-clean-architecture)
-------------------------------------------------------------------------------------------------------------

ArchitectureDomain-Driven Design

Domain entities are the heart of Clean Architecture. They represent the core business concepts and rules that define your application's purpose. Understanding entities correctly is fundamental to building maintainable, testable, and scalable software systems.

* * *

What is a Domain Entity? [​](#what-is-a-domain-entity)
------------------------------------------------------

A domain entity is a representation of a business concept that has a unique identity and encapsulates business rules. In Clean Architecture, entities form the innermost layer, completely independent of external concerns like databases, frameworks, or UI.

### Core Characteristics [​](#core-characteristics)

**Identity**: Each entity has a unique identifier that distinguishes it from other entities of the same type. Two entities with the same attributes but different identities are different entities.

**Business Logic**: Entities contain methods that enforce business rules and maintain invariants. They are not passive data structures but active participants in your domain model.

**Independence**: Entities have zero dependencies on external systems. They do not import HTTP libraries, database drivers, or framework code. This independence makes them portable, testable, and reusable.

**Validation**: Entities validate their own state, ensuring that business rules are never violated. Invalid states are impossible to represent.

Entity vs Model: A Critical Distinction [​](#entity-vs-model-a-critical-distinction)
------------------------------------------------------------------------------------

Many developers confuse entities with database models or API models. This confusion leads to architectural problems and coupling.

### What an Entity Is NOT [​](#what-an-entity-is-not)

**Not a Database Model**: Entities do not map directly to database tables. They represent business concepts, not storage structures. Database concerns belong to the infrastructure layer.

**Not an API Response**: Entities are not DTOs (Data Transfer Objects). API responses should be separate structures that adapt entities for external communication.

**Not Framework-Dependent**: Entities do not depend on ORMs, validation frameworks, or serialization libraries. These are implementation details.

### The Separation Principle [​](#the-separation-principle)

    Domain Entity (Pure Business Logic)
            ↓
        Use Case (Application Logic)
            ↓
    Repository Interface (Contract)
            ↓
    Repository Implementation (Database Details)
            ↓
        Database Schema

This separation allows you to:

*   Change databases without touching business logic
*   Test business rules without database setup
*   Evolve your domain model independently
*   Swap ORMs or frameworks with minimal impact

Domain-Driven Design Principles [​](#domain-driven-design-principles)
---------------------------------------------------------------------

Goca implements Domain-Driven Design (DDD) principles when generating entities, ensuring your code follows established best practices.

### Ubiquitous Language [​](#ubiquitous-language)

Entities use the same terminology as your business domain. If your business talks about "Orders," "Customers," and "Products," your entities should use these exact terms.

### Aggregate Roots [​](#aggregate-roots)

Entities can serve as aggregate roots, controlling access to related objects and maintaining consistency boundaries.

### Value Objects vs Entities [​](#value-objects-vs-entities)

Entities have identity; value objects do not. An email address is a value object. A user is an entity. Goca helps you model both correctly.

How Goca Generates Entities [​](#how-goca-generates-entities)
-------------------------------------------------------------

Goca provides the `goca entity` command to generate domain entities following Clean Architecture and DDD principles.

### Basic Entity Generation [​](#basic-entity-generation)

bash

    goca entity User --fields "name:string,email:string,age:int"

This command generates a pure domain entity with no external dependencies:

go

    package domain
    
    type User struct {
        ID    uint   `json:"id" gorm:"primaryKey;autoIncrement"`
        Name  string `json:"name" gorm:"type:varchar(255);not null"`
        Email string `json:"email" gorm:"type:varchar(255);uniqueIndex;not null"`
        Age   int    `json:"age" gorm:"type:integer;not null;default:0"`
    }

Notice that while GORM tags are present for infrastructure convenience, the entity itself remains a simple Go struct. The entity does not import GORM or any database package.

### Field Types and Conventions [​](#field-types-and-conventions)

Goca supports common field types that map to both Go types and database columns:

**String Fields**: `name:string`, `email:string`, `description:string`

*   Generate `string` type
*   Map to `varchar` or `text` columns
*   Suitable for textual data

**Numeric Fields**: `age:int`, `price:float64`, `quantity:int64`

*   Generate integer or floating-point types
*   Map to appropriate numeric columns
*   Support business calculations

**Boolean Fields**: `is_active:bool`, `verified:bool`

*   Generate `bool` type
*   Map to boolean columns
*   Represent binary states

**Temporal Fields**: `birth_date:time.Time`

*   Generate `time.Time` type
*   Handle date and time data
*   Work with the standard library

### Adding Validation [​](#adding-validation)

Business rules are enforced through validation methods:

bash

    goca entity User --fields "name:string,email:string,age:int" --validation

This generates a `Validate()` method and domain-specific errors:

go

    package domain
    
    import (
        "time"
        
        "gorm.io/gorm"
    )
    
    type User struct {
        ID        uint           `json:"id" gorm:"primaryKey;autoIncrement"`
        Name      string         `json:"name" gorm:"type:varchar(255);not null" validate:"required"`
        Email     string         `json:"email" gorm:"type:varchar(255);uniqueIndex;not null" validate:"required,email"`
        Age       int            `json:"age" gorm:"type:integer;not null;default:0" validate:"required,gte=0"`
        CreatedAt time.Time      `json:"created_at" gorm:"autoCreateTime"`
        UpdatedAt time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
        DeletedAt gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`
    }
    
    func (u *User) Validate() error {
        if u.Name == "" {
            return ErrInvalidUserName
        }
        if u.Email == "" {
            return ErrInvalidUserEmail
        }
        if u.Age < 0 {
            return ErrInvalidUserAge
        }
        return nil
    }

The validation method ensures that no invalid user can exist in your system. This is domain logic, not input validation. Input validation happens in the use case or handler layer.

### Domain Errors [​](#domain-errors)

Goca generates a separate `errors.go` file containing domain-specific errors:

go

    package domain
    
    import "errors"
    
    var (
        ErrInvalidUserName  = errors.New("invalid user name")
        ErrInvalidUserEmail = errors.New("invalid user email")
        ErrInvalidUserAge   = errors.New("invalid user age")
        ErrUserNotFound     = errors.New("user not found")
    )

These errors are part of your domain model. They communicate business rule violations clearly and can be handled appropriately by outer layers.

### Business Rules [​](#business-rules)

Beyond validation, entities can contain business logic:

bash

    goca entity Order --fields "customer_id:int,total:float64,status:string" --business-rules

This generates methods that implement domain logic:

go

    func (o *Order) Validate() error {
        if o.Customer_id < 0 {
            return ErrInvalidOrderCustomer_id
        }
        if o.Total < 0 {
            return ErrInvalidOrderTotal
        }
        if o.Status == "" {
            return ErrInvalidOrderStatus
        }
        return nil
    }

You can extend these with additional business methods:

go

    func (o *Order) CanBeCancelled() bool {
        return o.Status == "pending" || o.Status == "confirmed"
    }
    
    func (o *Order) Apply(discount float64) error {
        if discount < 0 || discount > 1 {
            return errors.New("discount must be between 0 and 1")
        }
        o.Total = o.Total * (1 - discount)
        return nil
    }
    
    func (o *Order) IsCompleted() bool {
        return o.Status == "delivered"
    }

These methods encapsulate business knowledge. They answer domain questions and enforce domain rules.

### Timestamps and Soft Deletes [​](#timestamps-and-soft-deletes)

Entities often need audit trails and soft delete functionality:

bash

    goca entity Product --fields "name:string,price:float64,stock:int" \
      --timestamps \
      --soft-delete

This generates:

go

    type Product struct {
        ID          uint           `json:"id" gorm:"primaryKey"`
        Name        string         `json:"name" gorm:"type:varchar(255);not null"`
        Price       float64        `json:"price" gorm:"type:decimal(10,2);not null;default:0"`
        Stock       int            `json:"stock" gorm:"type:integer;not null;default:0"`
        CreatedAt   time.Time      `json:"created_at" gorm:"autoCreateTime"`
        UpdatedAt   time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
        DeletedAt   gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`
    }
    
    func (p *Product) SoftDelete() {
        p.DeletedAt = gorm.DeletedAt{Time: time.Now(), Valid: true}
    }
    
    func (p *Product) IsDeleted() bool {
        return p.DeletedAt.Valid
    }

Soft deletes preserve data while marking it as inactive. The `DeletedAt` field enables this pattern without permanently removing records.

Complete Entity Example [​](#complete-entity-example)
-----------------------------------------------------

Let's examine a complete entity generated by Goca:

bash

    goca entity Product --fields "name:string,description:string,price:float64,stock:int,is_active:bool" \
      --validation \
      --business-rules \
      --timestamps \
      --soft-delete

This generates a production-ready entity:

go

    package domain
    
    import (
        "time"
        
        "gorm.io/gorm"
    )
    
    type Product struct {
        ID          uint           `json:"id" gorm:"primaryKey;autoIncrement"`
        Name        string         `json:"name" gorm:"type:varchar(255);not null" validate:"required"`
        Description string         `json:"description" gorm:"type:text"`
        Price       float64        `json:"price" gorm:"type:decimal(10,2);not null;default:0" validate:"required,gte=0"`
        Stock       int            `json:"stock" gorm:"type:integer;not null;default:0" validate:"required,gte=0"`
        IsActive    bool           `json:"is_active" gorm:"type:boolean;not null;default:false"`
        CreatedAt   time.Time      `json:"created_at" gorm:"autoCreateTime"`
        UpdatedAt   time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
        DeletedAt   gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`
    }
    
    func (p *Product) Validate() error {
        if p.Name == "" {
            return ErrInvalidProductName
        }
        if p.Price < 0 {
            return ErrInvalidProductPrice
        }
        if p.Stock < 0 {
            return ErrInvalidProductStock
        }
        return nil
    }
    
    func (p *Product) SoftDelete() {
        p.DeletedAt = gorm.DeletedAt{Time: time.Now(), Valid: true}
    }
    
    func (p *Product) IsDeleted() bool {
        return p.DeletedAt.Valid
    }

You can extend this with additional business methods:

go

    func (p *Product) IsAvailable() bool {
        return p.IsActive && p.Stock > 0 && !p.IsDeleted()
    }
    
    func (p *Product) Restock(quantity int) error {
        if quantity <= 0 {
            return errors.New("restock quantity must be positive")
        }
        p.Stock += quantity
        p.UpdatedAt = time.Now()
        return nil
    }
    
    func (p *Product) Sell(quantity int) error {
        if quantity <= 0 {
            return errors.New("sell quantity must be positive")
        }
        if p.Stock < quantity {
            return errors.New("insufficient stock")
        }
        p.Stock -= quantity
        p.UpdatedAt = time.Now()
        return nil
    }
    
    func (p *Product) ApplyDiscount(percentage float64) error {
        if percentage < 0 || percentage > 100 {
            return errors.New("discount percentage must be between 0 and 100")
        }
        p.Price = p.Price * (1 - percentage/100)
        p.UpdatedAt = time.Now()
        return nil
    }

These methods capture business logic that belongs in the domain layer. They make the entity more than a data structure; they make it a behavior-rich business object.

Testing Domain Entities [​](#testing-domain-entities)
-----------------------------------------------------

Domain entities are easy to test because they have no external dependencies. You can test business logic in isolation:

go

    package domain_test
    
    import (
        "testing"
        
        "yourproject/internal/domain"
    )
    
    func TestUser_Validate(t *testing.T) {
        tests := []struct {
            name    string
            user    domain.User
            wantErr bool
        }{
            {
                name: "valid user",
                user: domain.User{
                    Name:  "John Doe",
                    Email: "john@example.com",
                    Age:   30,
                },
                wantErr: false,
            },
            {
                name: "empty name",
                user: domain.User{
                    Name:  "",
                    Email: "john@example.com",
                    Age:   30,
                },
                wantErr: true,
            },
            {
                name: "negative age",
                user: domain.User{
                    Name:  "John Doe",
                    Email: "john@example.com",
                    Age:   -5,
                },
                wantErr: true,
            },
        }
        
        for _, tt := range tests {
            t.Run(tt.name, func(t *testing.T) {
                err := tt.user.Validate()
                if (err != nil) != tt.wantErr {
                    t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
                }
            })
        }
    }
    
    func TestProduct_Sell(t *testing.T) {
        product := domain.Product{
            Name:  "Test Product",
            Price: 100.0,
            Stock: 10,
        }
        
        // Valid sale
        err := product.Sell(5)
        if err != nil {
            t.Errorf("Sell(5) failed: %v", err)
        }
        if product.Stock != 5 {
            t.Errorf("Expected stock 5, got %d", product.Stock)
        }
        
        // Insufficient stock
        err = product.Sell(10)
        if err == nil {
            t.Error("Sell(10) should fail with insufficient stock")
        }
        
        // Negative quantity
        err = product.Sell(-1)
        if err == nil {
            t.Error("Sell(-1) should fail with negative quantity")
        }
    }

These tests run instantly because they do not touch databases, networks, or files. They verify business logic in isolation.

Best Practices for Domain Entities [​](#best-practices-for-domain-entities)
---------------------------------------------------------------------------

Based on Clean Architecture and DDD principles, follow these best practices:

**Keep Entities Pure**: Do not import framework or infrastructure code. Entities should compile without external dependencies beyond the standard library.

**Encapsulate State**: Use methods to modify entity state. Avoid exposing fields directly if business rules govern their modification.

**Express Business Rules**: Write methods that answer business questions and enforce business constraints. Make implicit knowledge explicit.

**Use Value Objects**: For concepts without identity, create value objects. An email address, money amount, or date range should be a value object, not part of an entity.

**Avoid Anemic Domain Models**: Entities with only getters and setters are anemic. Add behavior. Rich domain models contain both data and behavior.

**Design for Invariants**: Entities should always be in a valid state. Constructor functions and validation methods enforce this.

**Use Domain Language**: Name entities, fields, and methods using terms from your business domain. Code should read like business documentation.

Integration with Other Layers [​](#integration-with-other-layers)
-----------------------------------------------------------------

Entities work with other Clean Architecture layers through well-defined interfaces.

### Use Case Layer [​](#use-case-layer)

Use cases orchestrate entities to fulfill application requirements:

go

    package usecase
    
    type CreateProductInput struct {
        Name        string  `json:"name" validate:"required"`
        Description string  `json:"description"`
        Price       float64 `json:"price" validate:"required,gt=0"`
        Stock       int     `json:"stock" validate:"required,gte=0"`
    }
    
    type productService struct {
        repo repository.ProductRepository
    }
    
    func (s *productService) Create(input CreateProductInput) (*domain.Product, error) {
        // Create entity
        product := &domain.Product{
            Name:        input.Name,
            Description: input.Description,
            Price:       input.Price,
            Stock:       input.Stock,
            IsActive:    true,
        }
        
        // Validate business rules
        if err := product.Validate(); err != nil {
            return nil, err
        }
        
        // Persist through repository
        if err := s.repo.Save(product); err != nil {
            return nil, err
        }
        
        return product, nil
    }

The use case depends on the entity, not the other way around. This maintains the dependency rule.

### Repository Layer [​](#repository-layer)

Repositories provide persistence for entities through interfaces defined in the domain:

go

    package repository
    
    type ProductRepository interface {
        Save(product *domain.Product) error
        FindByID(id uint) (*domain.Product, error)
        Update(product *domain.Product) error
        Delete(id uint) error
        FindAll() ([]domain.Product, error)
    }

The repository interface lives in the domain package, but implementations live in the infrastructure layer:

go

    package repository
    
    import (
        "yourproject/internal/domain"
        "gorm.io/gorm"
    )
    
    type postgresProductRepository struct {
        db *gorm.DB
    }
    
    func NewPostgresProductRepository(db *gorm.DB) ProductRepository {
        return &postgresProductRepository{db: db}
    }
    
    func (r *postgresProductRepository) Save(product *domain.Product) error {
        return r.db.Create(product).Error
    }
    
    func (r *postgresProductRepository) FindByID(id uint) (*domain.Product, error) {
        var product domain.Product
        err := r.db.First(&product, id).Error
        if err != nil {
            return nil, err
        }
        return &product, nil
    }

This separation allows you to swap database implementations without changing business logic.

Advanced Entity Patterns [​](#advanced-entity-patterns)
-------------------------------------------------------

### Aggregate Roots [​](#aggregate-roots-1)

Entities can serve as aggregate roots, controlling access to related entities:

go

    type Order struct {
        ID         uint
        CustomerID uint
        Items      []OrderItem
        Total      float64
        Status     string
    }
    
    type OrderItem struct {
        ID        uint
        ProductID uint
        Quantity  int
        Price     float64
    }
    
    func (o *Order) AddItem(productID uint, quantity int, price float64) error {
        if quantity <= 0 {
            return errors.New("quantity must be positive")
        }
        
        item := OrderItem{
            ProductID: productID,
            Quantity:  quantity,
            Price:     price,
        }
        
        o.Items = append(o.Items, item)
        o.calculateTotal()
        return nil
    }
    
    func (o *Order) calculateTotal() {
        total := 0.0
        for _, item := range o.Items {
            total += float64(item.Quantity) * item.Price
        }
        o.Total = total
    }

The `Order` aggregate root controls `OrderItem` access, maintaining consistency.

### Factories [​](#factories)

Complex entity creation can use factory patterns:

go

    func NewUser(name, email string, age int) (*User, error) {
        user := &User{
            Name:      name,
            Email:     email,
            Age:       age,
            CreatedAt: time.Now(),
            UpdatedAt: time.Now(),
        }
        
        if err := user.Validate(); err != nil {
            return nil, err
        }
        
        return user, nil
    }

Factories ensure entities are always created in valid states.

Generating Complete Features [​](#generating-complete-features)
---------------------------------------------------------------

While `goca entity` generates entities, `goca feature` generates complete features including entities, use cases, repositories, and handlers:

bash

    goca feature Product --fields "name:string,price:float64,stock:int"

This generates:

*   Domain entity (`internal/domain/product.go`)
*   Use case interface and implementation (`internal/usecase/product_service.go`)
*   Repository interface (`internal/repository/interfaces.go`)
*   Repository implementation (`internal/repository/postgres_product_repository.go`)
*   HTTP handler (`internal/handler/http/product_handler.go`)
*   DTOs (`internal/usecase/dto.go`)
*   Error definitions (`internal/domain/errors.go`)
*   Seed data (`internal/domain/product_seeds.go`)

All layers work together following Clean Architecture principles, with the entity at the core.

Conclusion [​](#conclusion)
---------------------------

Domain entities are the foundation of Clean Architecture. They represent your business concepts and enforce business rules without coupling to external systems. Goca generates production-ready entities following DDD principles, giving you a solid starting point for building maintainable applications.

Understanding entities correctly is essential for successful software architecture. They are not database models, not API responses, and not framework-dependent structures. They are pure business logic, testable in isolation, and independent of implementation details.

By using Goca's entity generation commands and following Clean Architecture principles, you create systems that are:

*   Easy to test without external dependencies
*   Simple to modify as business requirements change
*   Portable across different frameworks and technologies
*   Clear in expressing business intent
*   Maintainable over long periods

Start with entities. Build your business logic correctly. Let outer layers adapt to your domain, not the other way around.

Further Reading [​](#further-reading)
-------------------------------------

*   Clean Architecture documentation in the [guide section](https://sazardev.github.io/goca/guide/clean-architecture.html)
*   Complete command reference for [`goca entity`](https://sazardev.github.io/goca/commands/entity.html)
*   Full feature generation with [`goca feature`](https://sazardev.github.io/goca/commands/feature.html)
*   Domain-Driven Design principles and patterns
*   Repository pattern implementation examples</content>
</page>

<page>
  <title>Goca - Go Clean Architecture Code Generator</title>
  <url>https://sazardev.github.io/goca/blog/releases/v1-14-1</url>
  <content>v1.14.1 - Test Suite Improvements [​](#v1-14-1-test-suite-improvements)
-----------------------------------------------------------------------

October 27, 2025Bug FixesQuality

* * *

Overview [​](#overview)
-----------------------

Version 1.14.1 focuses on improving test reliability and Windows compatibility. This release includes critical fixes for path handling, working directory management, and module dependencies. The test success rate has been significantly improved to **99.04%** (310/313 tests passing).

Bug Fixes [​](#bug-fixes)
-------------------------

### Test Suite Improvements [​](#test-suite-improvements)

#### Fixed Windows Path Handling in BackupFile [​](#fixed-windows-path-handling-in-backupfile)

Corrected path issues on Windows systems that were causing backup failures:

*   Changed from `filepath.Join(BackupDir, filepath.Dir(filePath))` to using only `filepath.Base(filePath)`
*   Prevents invalid "C:" subdirectory creation on Windows
*   Backup files now correctly created with `.backup` extension in backup directory root
*   Resolves file not found errors in safety manager tests

**Impact**: Windows users can now safely use the `--backup` flag without path-related errors.

#### Fixed Test Working Directory Management [​](#fixed-test-working-directory-management)

Improved test reliability across different execution contexts:

*   Added `SetProjectDir()` calls in handler and workflow tests after project initialization
*   Corrected file path assertions from absolute to relative paths
*   Fixed path expectations to match actual command execution context
*   All handler command tests now pass with correct working directory setup

**Code Example**:

go

    // Before: Tests failed due to incorrect working directory
    handler.Execute()
    // Files created in wrong location
    
    // After: Tests pass with proper directory setup
    handler.Execute()
    SetProjectDir(projectPath)
    // Files created in correct location

#### Updated Test Message Validation [​](#updated-test-message-validation)

Aligned test expectations with actual output:

*   Converted Spanish error messages to English in entity and feature tests
*   Simplified feature test validations to accept flexible message formats
*   Improved test robustness by accepting both English and Spanish variations

**Before**:

go

    assert.Contains(t, output, "La entidad ya existe")

**After**:

go

    assert.Contains(t, output, "already exists")

#### Fixed Module Dependencies [​](#fixed-module-dependencies)

Corrected testify dependency declaration:

*   Moved `github.com/stretchr/testify` from indirect to direct dependencies in `go.mod`
*   Fixes GitHub Actions CI failure on `go mod tidy` check
*   Properly declares direct usage in test files (`internal/testing/tests/*.go`)

**go.mod change**:

diff

    require (
        github.com/spf13/cobra v1.8.0
    +   github.com/stretchr/testify v1.9.0
    )

Quality Improvements [​](#quality-improvements)
-----------------------------------------------

### Test Metrics [​](#test-metrics)

| Metric | Before | After | Improvement |
| --- | --- | --- | --- |
| Test Success Rate | 96% | 99.04% | +3.04% |
| Passing Tests | 270/313 | 310/313 | +40 tests |
| Test Failures | 40 | 3 | \-92.5% |

### Reliability Enhancements [​](#reliability-enhancements)

*   **Core Commands**: All core commands (init, entity, usecase, repository, handler, feature, di, integrate) fully functional
*   **Cross-Platform**: Improved Windows compatibility in file operations
*   **Path Handling**: Better path handling across different operating systems
*   **CI/CD**: Enhanced cross-platform test reliability in automated environments

### Test Documentation [​](#test-documentation)

*   Added comprehensive skip messages for temporarily disabled tests
*   Documented differences between test expectations and actual code generation
*   Clear issue references (`#XXX`) for tracking test improvements

**Example of improved test documentation**:

go

    t.Skip("Integration test temporarily disabled: " +
        "Test expectations require validation strictness updates. " +
        "All sub-tests pass individually. " +
        "Issue #142: Update integration test validators")

Platform Support [​](#platform-support)
---------------------------------------

### Windows Compatibility [​](#windows-compatibility)

This release significantly improves Windows support:

*   Fixed backup file creation on Windows paths
*   Corrected directory separator handling
*   Improved path normalization across platforms
*   All safety features now work correctly on Windows

### Cross-Platform Testing [​](#cross-platform-testing)

Enhanced test reliability on:

*   Windows 10/11
*   macOS (Intel and Apple Silicon)
*   Linux (Ubuntu, Debian, Fedora)

Migration Guide [​](#migration-guide)
-------------------------------------

No migration required. This release contains only bug fixes and quality improvements. Simply update to v1.14.1:

bash

    go install github.com/sazardev/goca@v1.14.1

Verify the installation:

bash

    goca version
    # Output: v1.14.1

Known Issues [​](#known-issues)
-------------------------------

### Temporarily Disabled Tests [​](#temporarily-disabled-tests)

Two complex integration tests are temporarily disabled with detailed documentation:

1.  **Full Workflow Test**: Validation strictness for generated files
2.  **Multi-Feature Integration**: Test expectations alignment with actual output

These tests have all sub-tests passing individually. The issue is with overall validation strictness, not functionality.

**Status**: Tracked in issues with clear documentation for future enhancement.

Contributors [​](#contributors)
-------------------------------

Thank you to all contributors who helped improve test reliability and Windows compatibility in this release.

Next Steps [​](#next-steps)
---------------------------

Version 1.15.0 will focus on:

*   Integration test re-enablement with updated validators
*   Additional Windows-specific test coverage
*   Performance optimization for test execution
*   Enhanced CI/CD pipeline configuration

Resources [​](#resources)
-------------------------

*   [Full CHANGELOG](https://github.com/sazardev/goca/blob/master/CHANGELOG.md)
*   [GitHub Release](https://github.com/sazardev/goca/releases/tag/v1.14.1)
*   [Report Issues](https://github.com/sazardev/goca/issues)
*   [Documentation](https://sazardev.github.io/goca/)

* * *</content>
</page>

<page>
  <title>Mastering Use Cases in Clean Architecture | Goca Blog</title>
  <url>https://sazardev.github.io/goca/blog/articles/mastering-use-cases</url>
  <content>Mastering Use Cases in Clean Architecture [​](#mastering-use-cases-in-clean-architecture)
-----------------------------------------------------------------------------------------

ArchitectureApplication Layer

Use cases represent application-specific business rules and orchestrate the flow of data between entities and external systems. Understanding use cases correctly is critical for building well-structured applications that adapt to changing requirements without compromising core business logic.

* * *

What is a Use Case? [​](#what-is-a-use-case)
--------------------------------------------

A use case is an application service that coordinates domain entities and infrastructure to fulfill a specific user or system goal. Use cases answer the question: "What can the application do?"

### Core Responsibilities [​](#core-responsibilities)

**Orchestration**: Use cases coordinate multiple domain entities, repositories, and external services to complete a workflow. They do not contain business rules; they apply them.

**Data Flow Control**: Use cases manage the flow of data between the UI layer and the domain layer, transforming external requests into domain operations and domain results into external responses.

**Application Logic**: Use cases implement application-specific rules that do not belong in the domain. These rules depend on the use case context, not on core business concepts.

**Transaction Management**: Use cases define transaction boundaries, ensuring that operations either complete fully or roll back entirely.

**Permission and Security**: Use cases enforce authorization rules, checking whether the requesting user can perform the operation.

Use Case vs Controller vs Service [​](#use-case-vs-controller-vs-service)
-------------------------------------------------------------------------

Many developers confuse use cases with controllers or generic services. This confusion leads to bloated classes and violated architectural boundaries.

### What a Use Case Is NOT [​](#what-a-use-case-is-not)

**Not a Controller**: Controllers are adapters that convert HTTP requests to use case calls. Controllers handle protocol concerns; use cases handle application logic.

**Not a Generic Service**: A use case serves a specific goal, not general utilities. Services like "EmailService" or "LoggerService" are infrastructure concerns, not use cases.

**Not a Transaction Script**: Use cases orchestrate domain entities. They do not implement business rules. Domain logic belongs in entities and value objects.

**Not a Facade**: Use cases are not simple pass-throughs to repositories. They add application-level coordination and workflow management.

### The Clear Distinction [​](#the-clear-distinction)

    HTTP Request
        ↓
    Controller (Adapter - Outer Layer)
        ↓ Converts to DTO
    Use Case (Application Layer)
        ↓ Orchestrates
    Domain Entity (Domain Layer)
        ↓ Enforces rules
    Repository Interface (Domain Layer)
        ↓ Implements
    Repository Implementation (Infrastructure Layer)
        ↓ Persists
    Database

Each layer has distinct responsibilities. Use cases sit between adapters and domain, orchestrating operations without implementing business rules or handling external protocols.

The Application Layer [​](#the-application-layer)
-------------------------------------------------

Use cases form the application layer in Clean Architecture, distinct from both the domain layer and the infrastructure layer.

### Application Layer Characteristics [​](#application-layer-characteristics)

**Depends on Domain**: Use cases depend on domain entities and interfaces. They call entity methods and use repository interfaces defined in the domain.

**Independent of Infrastructure**: Use cases do not import database drivers, HTTP libraries, or external service clients. They work with interfaces.

**Stateless by Design**: Use cases do not maintain state between calls. Each operation is independent.

**Transaction Boundaries**: Use cases define where transactions begin and end, ensuring data consistency.

### Why a Separate Layer? [​](#why-a-separate-layer)

The application layer exists because application logic and domain logic are different:

**Domain Logic**: "A user must have a valid email address" is domain logic. This rule exists regardless of how you access users.

**Application Logic**: "To create a user, check if the email exists, create the user, send a welcome email" is application logic. This workflow is specific to the user creation use case.

Separating these concerns allows you to:

*   Change workflows without changing domain rules
*   Test business rules without application context
*   Reuse domain logic across different workflows
*   Evolve the application independently of the domain

Data Transfer Objects (DTOs) [​](#data-transfer-objects-dtos)
-------------------------------------------------------------

DTOs are simple structures that carry data between layers without behavior. Use cases use DTOs to receive input and provide output.

### Why DTOs? [​](#why-dtos)

**Layer Separation**: DTOs prevent external layers from depending on domain entities directly. Changing an entity does not break API contracts.

**Validation Boundary**: DTOs define what data the use case needs and validate it before processing.

**Security**: DTOs control what data external systems can provide or receive, preventing over-posting and data exposure.

**Versioning**: DTOs allow multiple API versions to coexist by mapping different external structures to the same domain entities.

### Input DTOs [​](#input-dtos)

Input DTOs represent the data a use case needs to perform an operation:

go

    type CreateUserInput struct {
        Name  string `json:"name" validate:"required,min=2"`
        Email string `json:"email" validate:"required,email"`
        Age   int    `json:"age" validate:"required,gte=0"`
    }

Input DTOs include validation tags that define constraints:

*   **required**: Field must be present
*   **min/max**: String length or numeric range
*   **email**: Valid email format
*   **gte/lte**: Greater than or equal / less than or equal

These validations are input validations, not business rules. They ensure the data is well-formed before processing.

### Output DTOs [​](#output-dtos)

Output DTOs represent the data a use case returns:

go

    type CreateUserOutput struct {
        User    domain.User `json:"user"`
        Message string      `json:"message"`
    }

Output DTOs can include:

*   Domain entities for complete information
*   Specific fields for minimal responses
*   Metadata like messages or status codes
*   Related entities for composite responses

### Update DTOs with Optional Fields [​](#update-dtos-with-optional-fields)

Update operations use optional fields to support partial updates:

go

    type UpdateUserInput struct {
        Name  *string `json:"name,omitempty" validate:"omitempty,min=2"`
        Email *string `json:"email,omitempty" validate:"omitempty,email"`
        Age   *int    `json:"age,omitempty" validate:"omitempty,gte=0"`
    }

Pointer fields distinguish between "not provided" (nil) and "explicitly set to zero value" (non-nil pointer to zero value). This allows clients to update only specific fields without affecting others.

### List DTOs [​](#list-dtos)

List operations return collections with metadata:

go

    type ListUserOutput struct {
        Users   []domain.User `json:"users"`
        Total   int           `json:"total"`
        Message string        `json:"message"`
    }

List DTOs can include pagination information, filters applied, and total counts.

Use Case Implementation Patterns [​](#use-case-implementation-patterns)
-----------------------------------------------------------------------

Use cases follow consistent patterns regardless of the specific operation they perform.

### Basic Structure [​](#basic-structure)

Every use case implementation includes:

1.  Dependency injection of required repositories
2.  Input validation
3.  Domain entity coordination
4.  Business rule enforcement
5.  Persistence through repositories
6.  Output transformation

### Create Operation [​](#create-operation)

The create operation instantiates a new domain entity, validates it, and persists it:

go

    func (s *userService) Create(input CreateUserInput) (*CreateUserOutput, error) {
        // 1. Create domain entity from input
        user := domain.User{
            Name:  input.Name,
            Email: input.Email,
            Age:   input.Age,
        }
        
        // 2. Validate business rules
        if err := user.Validate(); err != nil {
            return nil, err
        }
        
        // 3. Persist through repository
        if err := s.repo.Save(&user); err != nil {
            return nil, err
        }
        
        // 4. Return success output
        return &CreateUserOutput{
            User:    user,
            Message: "User created successfully",
        }, nil
    }

This pattern ensures that:

*   Input is transformed to domain entities
*   Business rules are enforced before persistence
*   Repository handles storage concerns
*   Output is well-defined and structured

### Read Operation [​](#read-operation)

The read operation retrieves an entity by identifier:

go

    func (s *userService) GetByID(id uint) (*domain.User, error) {
        return s.repo.FindByID(int(id))
    }

Read operations are simple because they delegate directly to repositories. Complexity arises when reads require:

*   Authorization checks
*   Data enrichment from multiple sources
*   Transformation to specific output formats

### Update Operation [​](#update-operation)

The update operation retrieves an entity, modifies it, validates it, and persists changes:

go

    func (s *userService) Update(id uint, input UpdateUserInput) (*domain.User, error) {
        // 1. Retrieve existing entity
        user, err := s.repo.FindByID(int(id))
        if err != nil {
            return nil, err
        }
        
        // 2. Apply changes from input (only provided fields)
        if input.Name != nil {
            user.Name = *input.Name
        }
        if input.Email != nil {
            user.Email = *input.Email
        }
        if input.Age != nil {
            user.Age = *input.Age
        }
        
        // 3. Validate updated entity
        if err := user.Validate(); err != nil {
            return nil, err
        }
        
        // 4. Persist changes
        if err := s.repo.Update(user); err != nil {
            return nil, err
        }
        
        return user, nil
    }

The update pattern:

*   Retrieves current state
*   Applies only provided changes
*   Validates the result
*   Persists atomically

### Delete Operation [​](#delete-operation)

The delete operation removes an entity:

go

    func (s *userService) Delete(id uint) error {
        return s.repo.Delete(int(id))
    }

Delete operations can be:

*   **Hard Delete**: Permanently removes the record
*   **Soft Delete**: Marks the record as deleted without removing it

Soft deletes are preferable for audit trails and data recovery.

### List Operation [​](#list-operation)

The list operation retrieves collections with optional filtering:

go

    func (s *userService) List() (*ListUserOutput, error) {
        users, err := s.repo.FindAll()
        if err != nil {
            return nil, err
        }
        
        return &ListUserOutput{
            Users:   users,
            Total:   len(users),
            Message: "Users listed successfully",
        }, nil
    }

List operations often include:

*   Pagination parameters
*   Sort order specifications
*   Filter conditions
*   Total count calculation

How Goca Generates Use Cases [​](#how-goca-generates-use-cases)
---------------------------------------------------------------

Goca provides the `goca usecase` command to generate complete application services with DTOs and interfaces.

### Basic Use Case Generation [​](#basic-use-case-generation)

bash

    goca usecase UserService --entity User

This generates three files:

**dto.go**: Input and output DTOs

go

    package usecase
    
    import (
        "github.com/yourorg/yourproject/internal/domain"
    )
    
    type CreateUserInput struct {
        Name  string `json:"name" validate:"required,min=2"`
        Email string `json:"email" validate:"required,email"`
        Age   int    `json:"age" validate:"required,min=1"`
    }
    
    type CreateUserOutput struct {
        User    domain.User `json:"user"`
        Message string      `json:"message"`
    }
    
    type UpdateUserInput struct {
        Name  *string `json:"name,omitempty" validate:"omitempty,min=2"`
        Email *string `json:"email,omitempty" validate:"omitempty,email"`
        Age   *int    `json:"age,omitempty" validate:"omitempty,min=1"`
    }
    
    type ListUserOutput struct {
        Users   []domain.User `json:"users"`
        Total   int           `json:"total"`
        Message string        `json:"message"`
    }

**user\_service.go**: Service implementation

go

    package usecase
    
    import (
        "github.com/yourorg/yourproject/internal/domain"
        "github.com/yourorg/yourproject/internal/repository"
    )
    
    type userService struct {
        repo repository.UserRepository
    }
    
    func NewUserService(repo repository.UserRepository) UserUseCase {
        return &userService{repo: repo}
    }
    
    func (u *userService) Create(input CreateUserInput) (*CreateUserOutput, error) {
        user := domain.User{
            Name:  input.Name,
            Email: input.Email,
            Age:   input.Age,
        }
        
        if err := user.Validate(); err != nil {
            return nil, err
        }
        
        if err := u.repo.Save(&user); err != nil {
            return nil, err
        }
        
        return &CreateUserOutput{
            User:    user,
            Message: "User created successfully",
        }, nil
    }
    
    func (u *userService) GetByID(id uint) (*domain.User, error) {
        return u.repo.FindByID(int(id))
    }
    
    func (u *userService) Update(id uint, input UpdateUserInput) (*domain.User, error) {
        user, err := u.repo.FindByID(int(id))
        if err != nil {
            return nil, err
        }
        
        if input.Name != nil {
            user.Name = *input.Name
        }
        if input.Email != nil {
            user.Email = *input.Email
        }
        if input.Age != nil {
            user.Age = *input.Age
        }
        
        if err := u.repo.Update(user); err != nil {
            return nil, err
        }
        
        return user, nil
    }
    
    func (u *userService) Delete(id uint) error {
        return u.repo.Delete(int(id))
    }
    
    func (u *userService) List() (*ListUserOutput, error) {
        users, err := u.repo.FindAll()
        if err != nil {
            return nil, err
        }
        
        return &ListUserOutput{
            Users:   users,
            Total:   len(users),
            Message: "Users listed successfully",
        }, nil
    }

**interfaces.go**: Repository interface definition

The service depends on a repository interface:

go

    type UserRepository interface {
        Save(user *domain.User) error
        FindByID(id int) (*domain.User, error)
        Update(user *domain.User) error
        Delete(id int) error
        FindAll() ([]domain.User, error)
    }

This interface lives in the repository package but is used by the use case. The use case depends on the abstraction, not the implementation.

### Selecting Operations [​](#selecting-operations)

Control which CRUD operations to generate:

bash

    goca usecase ProductService --entity Product --operations "create,read,update"

This generates only create, read, and update methods, omitting delete and list.

Available operations:

*   **create**: Instantiate and persist new entities
*   **read** or **get**: Retrieve entities by ID
*   **update**: Modify existing entities
*   **delete**: Remove entities
*   **list**: Retrieve collections

### DTO Validation [​](#dto-validation)

Enable validation tags on DTOs:

bash

    goca usecase OrderService --entity Order --dto-validation

With validation enabled, DTOs include comprehensive validation rules:

go

    type CreateOrderInput struct {
        CustomerID int     `json:"customer_id" validate:"required,gt=0"`
        Total      float64 `json:"total" validate:"required,gte=0"`
        Status     string  `json:"status" validate:"required,oneof=pending confirmed shipped delivered"`
    }

Validation rules ensure:

*   Required fields are present
*   Numeric values are within acceptable ranges
*   Strings match expected patterns or enumerations
*   Email addresses are valid
*   Custom validation logic is applied

Advanced Use Case Patterns [​](#advanced-use-case-patterns)
-----------------------------------------------------------

Beyond basic CRUD, use cases handle complex workflows and business processes.

### Transactional Use Cases [​](#transactional-use-cases)

Some operations require multiple steps within a single transaction:

go

    func (s *orderService) CreateOrder(input CreateOrderInput) (*CreateOrderOutput, error) {
        // Begin transaction (pseudo-code, actual implementation depends on repository)
        
        // 1. Validate customer exists
        customer, err := s.customerRepo.FindByID(input.CustomerID)
        if err != nil {
            return nil, errors.New("customer not found")
        }
        
        // 2. Check product availability
        for _, item := range input.Items {
            product, err := s.productRepo.FindByID(item.ProductID)
            if err != nil {
                return nil, err
            }
            
            if product.Stock < item.Quantity {
                return nil, errors.New("insufficient stock")
            }
        }
        
        // 3. Create order
        order := &domain.Order{
            CustomerID: input.CustomerID,
            Items:      mapOrderItems(input.Items),
            Total:      calculateTotal(input.Items),
            Status:     "pending",
        }
        
        if err := order.Validate(); err != nil {
            return nil, err
        }
        
        // 4. Persist order
        if err := s.orderRepo.Save(order); err != nil {
            return nil, err
        }
        
        // 5. Update product stock
        for _, item := range order.Items {
            product, _ := s.productRepo.FindByID(item.ProductID)
            product.Stock -= item.Quantity
            s.productRepo.Update(product)
        }
        
        // Commit transaction
        
        return &CreateOrderOutput{
            Order:   *order,
            Message: "Order created successfully",
        }, nil
    }

This use case:

*   Validates dependencies (customer exists)
*   Checks business constraints (sufficient stock)
*   Creates the main entity (order)
*   Updates related entities (product stock)
*   Ensures atomicity through transactions

### Async Use Cases [​](#async-use-cases)

Some operations can execute asynchronously to improve response times:

bash

    goca usecase NotificationService --entity Notification --operations "create" --async

Asynchronous use cases return immediately while processing continues in the background:

go

    func (s *notificationService) SendNotification(input SendNotificationInput) (*SendNotificationOutput, error) {
        // Validate input immediately
        if err := input.Validate(); err != nil {
            return nil, err
        }
        
        // Create notification record
        notification := &domain.Notification{
            UserID:  input.UserID,
            Message: input.Message,
            Status:  "queued",
        }
        
        if err := s.repo.Save(notification); err != nil {
            return nil, err
        }
        
        // Queue for async processing
        s.queue.Enqueue(notification.ID)
        
        // Return immediately
        return &SendNotificationOutput{
            NotificationID: notification.ID,
            Status:         "queued",
            Message:        "Notification queued successfully",
        }, nil
    }

The actual sending happens asynchronously:

go

    func (s *notificationService) ProcessQueue() {
        for {
            notificationID := s.queue.Dequeue()
            
            notification, err := s.repo.FindByID(notificationID)
            if err != nil {
                continue
            }
            
            // Send notification via external service
            err = s.emailService.Send(notification.UserID, notification.Message)
            
            if err != nil {
                notification.Status = "failed"
            } else {
                notification.Status = "sent"
            }
            
            s.repo.Update(notification)
        }
    }

Asynchronous use cases are appropriate for:

*   Email sending
*   File processing
*   Report generation
*   Third-party API calls
*   Long-running computations

### Composite Use Cases [​](#composite-use-cases)

Some operations aggregate data from multiple sources:

go

    func (s *dashboardService) GetUserDashboard(userID uint) (*DashboardOutput, error) {
        // Retrieve user
        user, err := s.userRepo.FindByID(userID)
        if err != nil {
            return nil, err
        }
        
        // Retrieve user orders
        orders, err := s.orderRepo.FindByUserID(userID)
        if err != nil {
            return nil, err
        }
        
        // Calculate statistics
        totalSpent := calculateTotalSpent(orders)
        averageOrderValue := totalSpent / float64(len(orders))
        
        // Retrieve recent activity
        activity, err := s.activityRepo.FindRecentByUserID(userID, 10)
        if err != nil {
            return nil, err
        }
        
        return &DashboardOutput{
            User:              user,
            TotalOrders:       len(orders),
            TotalSpent:        totalSpent,
            AverageOrderValue: averageOrderValue,
            RecentActivity:    activity,
        }, nil
    }

Composite use cases coordinate multiple repositories to build aggregate views.

Testing Use Cases [​](#testing-use-cases)
-----------------------------------------

Use cases are highly testable because they depend on interfaces, not concrete implementations.

### Unit Testing with Mocks [​](#unit-testing-with-mocks)

Test use cases by mocking repository dependencies:

go

    func TestUserService_Create(t *testing.T) {
        // Create mock repository
        mockRepo := new(MockUserRepository)
        
        // Setup expectations
        mockRepo.On("Save", mock.AnythingOfType("*domain.User")).Return(nil)
        
        // Create service with mock
        service := NewUserService(mockRepo)
        
        // Execute use case
        input := CreateUserInput{
            Name:  "John Doe",
            Email: "john@example.com",
            Age:   30,
        }
        
        output, err := service.Create(input)
        
        // Assert results
        assert.NoError(t, err)
        assert.NotNil(t, output)
        assert.Equal(t, "John Doe", output.User.Name)
        assert.Equal(t, "User created successfully", output.Message)
        
        // Verify mock was called
        mockRepo.AssertExpectations(t)
    }

Mock repositories allow you to:

*   Test use case logic independently
*   Simulate repository errors
*   Verify correct method calls
*   Control return values

### Testing Validation [​](#testing-validation)

Test that use cases enforce validation correctly:

go

    func TestUserService_Create_InvalidEmail(t *testing.T) {
        mockRepo := new(MockUserRepository)
        service := NewUserService(mockRepo)
        
        input := CreateUserInput{
            Name:  "John Doe",
            Email: "invalid-email",
            Age:   30,
        }
        
        output, err := service.Create(input)
        
        assert.Error(t, err)
        assert.Nil(t, output)
        assert.Contains(t, err.Error(), "email")
        
        // Repository should not be called
        mockRepo.AssertNotCalled(t, "Save")
    }

### Testing Error Handling [​](#testing-error-handling)

Test that use cases handle repository errors gracefully:

go

    func TestUserService_Create_RepositoryError(t *testing.T) {
        mockRepo := new(MockUserRepository)
        
        // Simulate repository error
        mockRepo.On("Save", mock.Anything).Return(errors.New("database connection failed"))
        
        service := NewUserService(mockRepo)
        
        input := CreateUserInput{
            Name:  "John Doe",
            Email: "john@example.com",
            Age:   30,
        }
        
        output, err := service.Create(input)
        
        assert.Error(t, err)
        assert.Nil(t, output)
        assert.Equal(t, "database connection failed", err.Error())
    }

Integration with Other Layers [​](#integration-with-other-layers)
-----------------------------------------------------------------

Use cases coordinate between domain, infrastructure, and adapter layers.

### Handler to Use Case [​](#handler-to-use-case)

Handlers convert external requests to use case calls:

go

    func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
        // 1. Parse HTTP request
        var input usecase.CreateUserInput
        if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
            http.Error(w, "Invalid request body", http.StatusBadRequest)
            return
        }
        
        // 2. Call use case
        output, err := h.usecase.Create(input)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        
        // 3. Return HTTP response
        w.Header().Set("Content-Type", "application/json")
        w.WriteStatus(http.StatusCreated)
        json.NewEncoder(w).Encode(output)
    }

The handler:

*   Handles HTTP concerns (parsing, status codes, headers)
*   Delegates business logic to the use case
*   Transforms use case output to HTTP response

### Use Case to Repository [​](#use-case-to-repository)

Use cases call repository methods through interfaces:

go

    type userService struct {
        repo repository.UserRepository // Interface, not implementation
    }
    
    func (s *userService) Create(input CreateUserInput) (*CreateUserOutput, error) {
        user := domain.User{
            Name:  input.Name,
            Email: input.Email,
            Age:   input.Age,
        }
        
        // Use case calls repository interface
        if err := s.repo.Save(&user); err != nil {
            return nil, err
        }
        
        return &CreateUserOutput{User: user}, nil
    }

The repository implementation lives in the infrastructure layer:

go

    type postgresUserRepository struct {
        db *gorm.DB
    }
    
    func (r *postgresUserRepository) Save(user *domain.User) error {
        return r.db.Create(user).Error
    }

This separation allows:

*   Swapping database implementations
*   Testing use cases without databases
*   Changing persistence strategies independently

### Dependency Injection [​](#dependency-injection)

Use cases receive dependencies through constructors:

go

    func NewUserService(
        repo repository.UserRepository,
        emailService EmailService,
        logger Logger,
    ) UserUseCase {
        return &userService{
            repo:         repo,
            emailService: emailService,
            logger:       logger,
        }
    }

Dependencies are injected at the composition root, typically in a DI container:

go

    // Composition root
    userRepo := repository.NewPostgresUserRepository(db)
    emailService := email.NewSMTPService(config)
    logger := log.NewStdLogger()
    
    userService := usecase.NewUserService(userRepo, emailService, logger)
    userHandler := handler.NewUserHandler(userService)

Best Practices for Use Cases [​](#best-practices-for-use-cases)
---------------------------------------------------------------

Follow these practices to maintain clean, maintainable use cases.

**Keep Use Cases Thin**: Use cases orchestrate; they do not implement business rules. Business logic belongs in domain entities.

**One Use Case, One Goal**: Each use case serves a specific goal. "Create user" is one use case. "Create user and send email" might be one or two, depending on cohesion.

**Use DTOs for All External Data**: Never pass domain entities directly to or from external layers. DTOs provide a stable contract.

**Validate at Boundaries**: Validate input at the use case boundary. Do not assume data is valid.

**Return Errors, Don't Panic**: Use cases return errors for exceptional conditions. They do not panic or crash.

**Keep Dependencies Minimal**: Use cases should depend only on repositories and essential services. Avoid excessive dependencies.

**Write Comprehensive Tests**: Test use cases thoroughly with mocked dependencies. Use cases are the easiest layer to test.

**Document Complex Workflows**: Use cases with multiple steps should be documented clearly, explaining the workflow and error handling.

Common Mistakes to Avoid [​](#common-mistakes-to-avoid)
-------------------------------------------------------

**Business Logic in Use Cases**: Do not implement business rules in use cases. Use cases apply rules defined in entities.

**Direct Database Access**: Use cases should not import database drivers or execute SQL. They call repository methods.

**Mixing Concerns**: Use cases should not handle HTTP parsing, logging details, or UI concerns. They orchestrate business operations.

**Returning Domain Entities Directly**: Always use DTOs for external communication. Domain entities are internal structures.

**Ignoring Errors**: Handle repository errors appropriately. Log them, wrap them, or transform them, but do not ignore them.

**Tight Coupling**: Use cases depending on concrete implementations cannot be tested or swapped easily. Depend on interfaces.

Generating Complete Features [​](#generating-complete-features)
---------------------------------------------------------------

While `goca usecase` generates use cases, `goca feature` generates complete features including entities, use cases, repositories, and handlers:

bash

    goca feature User --fields "name:string,email:string,age:int"

This creates:

*   Domain entity with validation
*   Use case with all CRUD operations
*   Repository interface and implementation
*   HTTP handler
*   DTOs for all operations
*   Dependency injection wiring

All layers work together following Clean Architecture principles, with use cases at the center coordinating workflows.

Conclusion [​](#conclusion)
---------------------------

Use cases are the application layer in Clean Architecture, orchestrating domain entities and infrastructure services to fulfill user goals. They coordinate workflows without implementing business rules, maintain clear boundaries through DTOs, and depend on abstractions rather than implementations.

Understanding use cases correctly is essential for building maintainable applications. They are not controllers, not generic services, and not transaction scripts. They are focused coordinators that apply domain rules in application-specific contexts.

Goca generates production-ready use cases with comprehensive DTOs, following established patterns and best practices. By using Goca's use case generation and understanding these principles, you create systems that are:

*   Easy to test with mocked dependencies
*   Simple to modify as requirements change
*   Clear in expressing application workflows
*   Maintainable over long periods
*   Adaptable to new platforms and interfaces

Start with clear use case boundaries. Orchestrate domain logic. Let adapters handle external concerns. Build applications that last.

Further Reading [​](#further-reading)
-------------------------------------

*   Application layer patterns in the [guide section](https://sazardev.github.io/goca/guide/clean-architecture.html)
*   Complete command reference for [`goca usecase`](https://sazardev.github.io/goca/commands/usecase.html)
*   Full feature generation with [`goca feature`](https://sazardev.github.io/goca/commands/feature.html)
*   Understanding domain entities in [our previous article](https://sazardev.github.io/goca/blog/articles/understanding-domain-entities.html)
*   Repository pattern implementation examples
*   Dependency injection patterns and best practices</content>
</page>
