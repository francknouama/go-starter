# Example Projects

Complete example projects demonstrating different use cases and patterns with go-starter.

## Quick Examples

### 1. Simple REST API

**Generate the project:**
```bash
go-starter new todo-api --type=web-api --framework=gin
```

**Key implementation:**
```go
// main.go
package main

import (
    "github.com/gin-gonic/gin"
    "github.com/yourusername/todo-api/internal/config"
    "github.com/yourusername/todo-api/internal/handlers"
    "github.com/yourusername/todo-api/internal/middleware"
)

func main() {
    cfg := config.Load()
    
    r := gin.New()
    r.Use(middleware.Logger())
    r.Use(middleware.ErrorHandler())
    
    // Routes
    api := r.Group("/api/v1")
    {
        api.GET("/todos", handlers.ListTodos)
        api.POST("/todos", handlers.CreateTodo)
        api.GET("/todos/:id", handlers.GetTodo)
        api.PUT("/todos/:id", handlers.UpdateTodo)
        api.DELETE("/todos/:id", handlers.DeleteTodo)
    }
    
    r.Run(cfg.Server.Port)
}
```

**Example handler:**
```go
// internal/handlers/todo.go
type Todo struct {
    ID          string    `json:"id"`
    Title       string    `json:"title" binding:"required"`
    Description string    `json:"description"`
    Completed   bool      `json:"completed"`
    CreatedAt   time.Time `json:"created_at"`
}

func CreateTodo(c *gin.Context) {
    var todo Todo
    if err := c.ShouldBindJSON(&todo); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }
    
    todo.ID = uuid.New().String()
    todo.CreatedAt = time.Now()
    
    // Save to database
    if err := db.Create(&todo); err != nil {
        c.JSON(500, gin.H{"error": "Failed to create todo"})
        return
    }
    
    c.JSON(201, todo)
}
```

---

### 2. CLI Tool with Subcommands

**Generate the project:**
```bash
go-starter new devtools --type=cli --framework=cobra
```

**Command structure:**
```go
// cmd/root.go
var rootCmd = &cobra.Command{
    Use:   "devtools",
    Short: "A collection of development tools",
    Long:  `DevTools provides utilities for everyday development tasks.`,
}

func init() {
    rootCmd.AddCommand(encodeCmd)
    rootCmd.AddCommand(formatCmd)
    rootCmd.AddCommand(generateCmd)
}
```

**Example subcommand:**
```go
// cmd/encode.go
var encodeCmd = &cobra.Command{
    Use:   "encode [type] [input]",
    Short: "Encode data in various formats",
    Args:  cobra.ExactArgs(2),
    Run: func(cmd *cobra.Command, args []string) {
        encodeType := args[0]
        input := args[1]
        
        switch encodeType {
        case "base64":
            encoded := base64.StdEncoding.EncodeToString([]byte(input))
            fmt.Println(encoded)
        case "url":
            encoded := url.QueryEscape(input)
            fmt.Println(encoded)
        case "hex":
            encoded := hex.EncodeToString([]byte(input))
            fmt.Println(encoded)
        default:
            fmt.Printf("Unknown encoding type: %s\n", encodeType)
        }
    },
}
```

---

### 3. Microservice with gRPC

**Generate the project:**
```bash
go-starter new user-service --type=microservice
```

**Proto definition:**
```proto
// proto/user.proto
syntax = "proto3";

package user.v1;

service UserService {
    rpc CreateUser(CreateUserRequest) returns (User);
    rpc GetUser(GetUserRequest) returns (User);
    rpc UpdateUser(UpdateUserRequest) returns (User);
    rpc ListUsers(ListUsersRequest) returns (ListUsersResponse);
}

message User {
    string id = 1;
    string email = 2;
    string name = 3;
    google.protobuf.Timestamp created_at = 4;
}
```

**Service implementation:**
```go
// internal/service/user.go
type userService struct {
    user.UnimplementedUserServiceServer
    repo repository.UserRepository
}

func (s *userService) CreateUser(ctx context.Context, req *user.CreateUserRequest) (*user.User, error) {
    // Validate request
    if err := validateEmail(req.Email); err != nil {
        return nil, status.Errorf(codes.InvalidArgument, "invalid email: %v", err)
    }
    
    // Create user
    u := &model.User{
        ID:        uuid.New().String(),
        Email:     req.Email,
        Name:      req.Name,
        CreatedAt: time.Now(),
    }
    
    if err := s.repo.Create(ctx, u); err != nil {
        return nil, status.Errorf(codes.Internal, "failed to create user: %v", err)
    }
    
    return toProtoUser(u), nil
}
```

---

### 4. AWS Lambda Function

**Generate the project:**
```bash
go-starter new image-processor --type=lambda
```

**Handler implementation:**
```go
// handler.go
package main

import (
    "context"
    "encoding/json"
    "fmt"
    
    "github.com/aws/aws-lambda-go/events"
    "github.com/aws/aws-lambda-go/lambda"
    "github.com/aws/aws-sdk-go-v2/config"
    "github.com/aws/aws-sdk-go-v2/service/s3"
)

type ImageEvent struct {
    Bucket string `json:"bucket"`
    Key    string `json:"key"`
    Size   string `json:"size"`
}

func handler(ctx context.Context, s3Event events.S3Event) error {
    cfg, err := config.LoadDefaultConfig(ctx)
    if err != nil {
        return fmt.Errorf("failed to load config: %w", err)
    }
    
    s3Client := s3.NewFromConfig(cfg)
    
    for _, record := range s3Event.Records {
        bucket := record.S3.Bucket.Name
        key := record.S3.Object.Key
        
        // Process image
        if err := processImage(ctx, s3Client, bucket, key); err != nil {
            return fmt.Errorf("failed to process %s/%s: %w", bucket, key, err)
        }
    }
    
    return nil
}

func main() {
    lambda.Start(handler)
}
```

---

### 5. Library Package

**Generate the project:**
```bash
go-starter new go-validators --type=library
```

**Library implementation:**
```go
// validators.go
package validators

import (
    "net/mail"
    "regexp"
)

// Email validates email addresses
func Email(email string) error {
    _, err := mail.ParseAddress(email)
    if err != nil {
        return &ValidationError{
            Field:   "email",
            Message: "invalid email format",
        }
    }
    return nil
}

// PhoneNumber validates phone numbers
func PhoneNumber(phone string) error {
    pattern := `^[+]?[(]?[0-9]{3}[)]?[-\s.]?[(]?[0-9]{3}[)]?[-\s.]?[0-9]{4,6}$`
    matched, _ := regexp.MatchString(pattern, phone)
    if !matched {
        return &ValidationError{
            Field:   "phone",
            Message: "invalid phone number format",
        }
    }
    return nil
}

// URL validates URLs
func URL(url string) error {
    pattern := `^(https?://)?([\da-z.-]+)\.([a-z.]{2,6})([/\w .-]*)*/?$`
    matched, _ := regexp.MatchString(pattern, url)
    if !matched {
        return &ValidationError{
            Field:   "url",
            Message: "invalid URL format",
        }
    }
    return nil
}
```

---

## Complete Example Projects

### E-Commerce API (Clean Architecture)

**Project Structure:**
```
ecommerce-api/
├── cmd/
│   └── api/
│       └── main.go
├── internal/
│   ├── domain/
│   │   ├── product.go
│   │   ├── order.go
│   │   └── user.go
│   ├── usecase/
│   │   ├── product_usecase.go
│   │   └── order_usecase.go
│   ├── repository/
│   │   ├── product_repository.go
│   │   └── order_repository.go
│   └── delivery/
│       └── http/
│           ├── product_handler.go
│           └── order_handler.go
└── pkg/
    └── errors/
```

**Domain Model:**
```go
// internal/domain/product.go
type Product struct {
    ID          string    `json:"id"`
    Name        string    `json:"name"`
    Description string    `json:"description"`
    Price       float64   `json:"price"`
    Stock       int       `json:"stock"`
    CategoryID  string    `json:"category_id"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
}

type ProductRepository interface {
    Create(ctx context.Context, product *Product) error
    GetByID(ctx context.Context, id string) (*Product, error)
    Update(ctx context.Context, product *Product) error
    Delete(ctx context.Context, id string) error
    List(ctx context.Context, filter ProductFilter) ([]*Product, error)
}

type ProductUseCase interface {
    CreateProduct(ctx context.Context, product *Product) error
    GetProduct(ctx context.Context, id string) (*Product, error)
    UpdateProduct(ctx context.Context, product *Product) error
    DeleteProduct(ctx context.Context, id string) error
    ListProducts(ctx context.Context, filter ProductFilter) ([]*Product, error)
}
```

---

### Real-Time Chat Server (WebSocket)

**Generate base:**
```bash
go-starter new chat-server --type=web-api --framework=gin
```

**WebSocket Implementation:**
```go
// internal/websocket/hub.go
type Hub struct {
    clients    map[*Client]bool
    broadcast  chan []byte
    register   chan *Client
    unregister chan *Client
}

func NewHub() *Hub {
    return &Hub{
        broadcast:  make(chan []byte),
        register:   make(chan *Client),
        unregister: make(chan *Client),
        clients:    make(map[*Client]bool),
    }
}

func (h *Hub) Run() {
    for {
        select {
        case client := <-h.register:
            h.clients[client] = true
            log.Printf("Client connected: %s", client.id)
            
        case client := <-h.unregister:
            if _, ok := h.clients[client]; ok {
                delete(h.clients, client)
                close(client.send)
                log.Printf("Client disconnected: %s", client.id)
            }
            
        case message := <-h.broadcast:
            for client := range h.clients {
                select {
                case client.send <- message:
                default:
                    close(client.send)
                    delete(h.clients, client)
                }
            }
        }
    }
}
```

---

### Event-Driven Order Processing

**Generate base:**
```bash
go-starter new order-processor --type=event-driven
```

**Event Handler:**
```go
// internal/events/order_events.go
type OrderCreatedEvent struct {
    OrderID    string    `json:"order_id"`
    CustomerID string    `json:"customer_id"`
    Items      []Item    `json:"items"`
    Total      float64   `json:"total"`
    CreatedAt  time.Time `json:"created_at"`
}

type OrderEventHandler struct {
    orderRepo    repository.OrderRepository
    inventoryService service.InventoryService
    paymentService   service.PaymentService
    eventBus         EventBus
}

func (h *OrderEventHandler) HandleOrderCreated(ctx context.Context, event OrderCreatedEvent) error {
    // Validate inventory
    for _, item := range event.Items {
        available, err := h.inventoryService.CheckAvailability(ctx, item.ProductID, item.Quantity)
        if err != nil || !available {
            return h.publishOrderFailed(ctx, event.OrderID, "Insufficient inventory")
        }
    }
    
    // Process payment
    paymentResult, err := h.paymentService.ProcessPayment(ctx, event.CustomerID, event.Total)
    if err != nil {
        return h.publishOrderFailed(ctx, event.OrderID, "Payment failed")
    }
    
    // Reserve inventory
    for _, item := range event.Items {
        if err := h.inventoryService.Reserve(ctx, item.ProductID, item.Quantity); err != nil {
            // Compensate previous reservations
            return h.compensateReservations(ctx, event.Items, item)
        }
    }
    
    // Publish success event
    return h.eventBus.Publish(ctx, OrderConfirmedEvent{
        OrderID:   event.OrderID,
        PaymentID: paymentResult.ID,
    })
}
```

---

## Testing Examples

### Unit Testing

```go
// internal/service/user_service_test.go
func TestUserService_CreateUser(t *testing.T) {
    tests := []struct {
        name    string
        input   *CreateUserInput
        mockFn  func(*mocks.MockUserRepository)
        wantErr bool
    }{
        {
            name: "successful creation",
            input: &CreateUserInput{
                Email: "test@example.com",
                Name:  "Test User",
            },
            mockFn: func(m *mocks.MockUserRepository) {
                m.EXPECT().
                    Create(gomock.Any(), gomock.Any()).
                    Return(nil)
            },
            wantErr: false,
        },
        {
            name: "duplicate email",
            input: &CreateUserInput{
                Email: "existing@example.com",
                Name:  "Test User",
            },
            mockFn: func(m *mocks.MockUserRepository) {
                m.EXPECT().
                    Create(gomock.Any(), gomock.Any()).
                    Return(ErrDuplicateEmail)
            },
            wantErr: true,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            ctrl := gomock.NewController(t)
            defer ctrl.Finish()
            
            mockRepo := mocks.NewMockUserRepository(ctrl)
            tt.mockFn(mockRepo)
            
            svc := NewUserService(mockRepo)
            _, err := svc.CreateUser(context.Background(), tt.input)
            
            if (err != nil) != tt.wantErr {
                t.Errorf("CreateUser() error = %v, wantErr %v", err, tt.wantErr)
            }
        })
    }
}
```

### Integration Testing

```go
// tests/integration/api_test.go
func TestAPI_CreateProduct(t *testing.T) {
    // Setup test database
    db := setupTestDB(t)
    defer cleanupTestDB(t, db)
    
    // Create test server
    app := createTestApp(db)
    
    // Test cases
    tests := []struct {
        name       string
        payload    string
        wantStatus int
    }{
        {
            name: "valid product",
            payload: `{
                "name": "Test Product",
                "price": 99.99,
                "stock": 100
            }`,
            wantStatus: http.StatusCreated,
        },
        {
            name:       "invalid payload",
            payload:    `{"name": ""}`,
            wantStatus: http.StatusBadRequest,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            req := httptest.NewRequest("POST", "/api/v1/products", strings.NewReader(tt.payload))
            req.Header.Set("Content-Type", "application/json")
            
            rec := httptest.NewRecorder()
            app.ServeHTTP(rec, req)
            
            assert.Equal(t, tt.wantStatus, rec.Code)
        })
    }
}
```

---

## Running the Examples

### Clone Example Repository

```bash
# Clone all examples
git clone https://github.com/go-starter/examples.git
cd examples

# List available examples
ls -la

# Run specific example
cd todo-api
go mod download
go run main.go
```

### Try Live Examples

1. **Todo API**: https://todo-api.go-starter.dev
2. **Chat Demo**: https://chat.go-starter.dev
3. **GraphQL Playground**: https://graphql.go-starter.dev

### Build Your Own

1. Start with a blueprint:
   ```bash
   go-starter new myproject --type=web-api
   ```

2. Follow the patterns in these examples

3. Customize for your needs

4. Share with the community!

---

*Have an interesting example? [Contribute it](https://github.com/go-starter/examples)!*