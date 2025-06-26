# ORM Selection Guide

A comprehensive guide to choosing and using the right ORM (Object-Relational Mapping) solution in your go-starter projects.

## Current ORM Support

go-starter currently supports two database interaction approaches:

| ORM | Status | Best For | Performance | Learning Curve |
|-----|--------|----------|-------------|----------------|
| **GORM** | ✅ Fully Implemented | Rapid development, CRUD operations | Good | Low |
| **Raw SQL** | ✅ Fully Implemented | Complex queries, full control | Excellent | Medium |

Future releases will include: sqlx, sqlc, ent, and more options.

## Quick Decision Guide

```
Choose GORM when:
✅ Building CRUD-heavy applications
✅ Need rapid prototyping
✅ Want automatic migrations
✅ Prefer convention over configuration
✅ Team is familiar with ORMs

Choose Raw SQL when:
✅ Need maximum performance
✅ Have complex queries
✅ Want full control over SQL
✅ Working with existing database schemas
✅ Team has strong SQL skills
```

## GORM - The Full-Featured ORM

### Overview

GORM is a developer-friendly ORM that provides:
- Auto-migrations
- Associations (Has One, Has Many, Belongs To, Many To Many)
- Hooks (Before/After Create/Update/Delete/Find)
- Preloading
- Transactions
- SQL Builder
- Auto-increment IDs
- Composite Primary Keys

### When to Use GORM

#### Perfect For:
- **CRUD Applications**: Simple Create, Read, Update, Delete operations
- **Rapid Development**: Get up and running quickly
- **Admin Panels**: Quick interfaces for data management
- **Prototypes**: Fast iteration on data models
- **Standard Web Apps**: Typical web application needs

#### Example Use Cases:
- E-commerce platforms
- Content management systems
- User management systems
- Blog platforms
- API backends with standard operations

### GORM Code Examples

#### Model Definition
```go
package models

import (
    "time"
    "gorm.io/gorm"
)

type User struct {
    ID        uint           `gorm:"primarykey"`
    CreatedAt time.Time
    UpdatedAt time.Time
    DeletedAt gorm.DeletedAt `gorm:"index"`
    
    // User fields
    Email     string         `gorm:"uniqueIndex;not null"`
    Username  string         `gorm:"uniqueIndex;not null"`
    Password  string         `gorm:"not null"`
    
    // Associations
    Profile   Profile        `gorm:"constraint:OnDelete:CASCADE;"`
    Posts     []Post         `gorm:"foreignKey:AuthorID"`
}

type Profile struct {
    ID        uint      `gorm:"primarykey"`
    UserID    uint      `gorm:"uniqueIndex"`
    FirstName string
    LastName  string
    Bio       string    `gorm:"type:text"`
    Avatar    string
}

type Post struct {
    ID        uint           `gorm:"primarykey"`
    CreatedAt time.Time
    UpdatedAt time.Time
    
    Title     string         `gorm:"not null"`
    Content   string         `gorm:"type:text"`
    AuthorID  uint
    Author    User           `gorm:"foreignKey:AuthorID"`
    Tags      []Tag          `gorm:"many2many:post_tags;"`
}

type Tag struct {
    ID    uint   `gorm:"primarykey"`
    Name  string `gorm:"uniqueIndex"`
    Posts []Post `gorm:"many2many:post_tags;"`
}
```

#### Repository Pattern with GORM
```go
package repository

import (
    "context"
    "errors"
    "gorm.io/gorm"
)

type UserRepository struct {
    db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
    return &UserRepository{db: db}
}

// Create a new user with profile
func (r *UserRepository) Create(ctx context.Context, user *User) error {
    return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
        if err := tx.Create(user).Error; err != nil {
            return err
        }
        
        // Create associated profile
        if user.Profile.UserID == 0 {
            user.Profile.UserID = user.ID
        }
        
        return tx.Create(&user.Profile).Error
    })
}

// Find user by ID with associations
func (r *UserRepository) FindByID(ctx context.Context, id uint) (*User, error) {
    var user User
    err := r.db.WithContext(ctx).
        Preload("Profile").
        Preload("Posts").
        First(&user, id).Error
        
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, ErrUserNotFound
        }
        return nil, err
    }
    
    return &user, nil
}

// Update user with optimistic locking
func (r *UserRepository) Update(ctx context.Context, user *User) error {
    result := r.db.WithContext(ctx).
        Model(user).
        Updates(user)
        
    if result.Error != nil {
        return result.Error
    }
    
    if result.RowsAffected == 0 {
        return ErrUserNotFound
    }
    
    return nil
}

// Complex query example
func (r *UserRepository) FindActiveUsersWithPosts(ctx context.Context, limit int) ([]User, error) {
    var users []User
    
    err := r.db.WithContext(ctx).
        Preload("Posts", func(db *gorm.DB) *gorm.DB {
            return db.Order("created_at DESC").Limit(5)
        }).
        Joins("LEFT JOIN posts ON posts.author_id = users.id").
        Where("users.deleted_at IS NULL").
        Group("users.id").
        Having("COUNT(posts.id) > ?", 0).
        Limit(limit).
        Find(&users).Error
        
    return users, err
}

// Pagination example
func (r *UserRepository) Paginate(ctx context.Context, page, pageSize int) ([]User, int64, error) {
    var users []User
    var total int64
    
    offset := (page - 1) * pageSize
    
    // Count total records
    if err := r.db.Model(&User{}).Count(&total).Error; err != nil {
        return nil, 0, err
    }
    
    // Fetch page
    err := r.db.WithContext(ctx).
        Preload("Profile").
        Offset(offset).
        Limit(pageSize).
        Find(&users).Error
        
    return users, total, err
}
```

#### Migration Management
```go
package database

import (
    "fmt"
    "gorm.io/gorm"
)

func AutoMigrate(db *gorm.DB) error {
    return db.AutoMigrate(
        &User{},
        &Profile{},
        &Post{},
        &Tag{},
    )
}

// Custom migration for complex changes
func MigrateAddUserStatus(db *gorm.DB) error {
    type User struct {
        Status string `gorm:"default:'active'"`
    }
    
    return db.AutoMigrate(&User{})
}

// Add indexes
func AddIndexes(db *gorm.DB) error {
    // Composite index
    if err := db.Exec("CREATE INDEX idx_users_email_username ON users(email, username)").Error; err != nil {
        return err
    }
    
    // Partial index (PostgreSQL)
    if err := db.Exec("CREATE INDEX idx_active_users ON users(email) WHERE deleted_at IS NULL").Error; err != nil {
        return err
    }
    
    return nil
}
```

### GORM Performance Tips

1. **Use Preloading Wisely**
```go
// Bad: N+1 queries
users, _ := db.Find(&users)
for _, user := range users {
    db.Model(&user).Association("Posts").Find(&user.Posts)
}

// Good: 2 queries total
db.Preload("Posts").Find(&users)

// Better: Selective preloading
db.Preload("Posts", "status = ?", "published").Find(&users)
```

2. **Select Specific Fields**
```go
// Select only needed fields
var users []UserDTO
db.Model(&User{}).
    Select("id", "email", "username").
    Find(&users)
```

3. **Use Batch Operations**
```go
// Batch insert
users := []User{{Name: "user1"}, {Name: "user2"}}
db.CreateInBatches(users, 100)

// Batch update
db.Model(&User{}).Where("status = ?", "inactive").
    UpdateColumn("status", "archived")
```

## Raw SQL - Full Control Approach

### Overview

Raw SQL with database/sql provides:
- Complete control over queries
- Maximum performance
- Direct database feature access
- No ORM overhead
- Explicit transaction management

### When to Use Raw SQL

#### Perfect For:
- **High-Performance Requirements**: When every millisecond counts
- **Complex Queries**: Advanced JOINs, CTEs, window functions
- **Database-Specific Features**: Using PostgreSQL arrays, JSON operations
- **Legacy Databases**: Working with existing schemas
- **Data Warehousing**: Complex analytical queries

#### Example Use Cases:
- Financial systems
- Analytics platforms
- High-frequency trading systems
- Report generation
- Data migration tools

### Raw SQL Code Examples

#### Database Connection Setup
```go
package database

import (
    "database/sql"
    "fmt"
    "time"
    
    _ "github.com/lib/pq"        // PostgreSQL
    _ "github.com/go-sql-driver/mysql" // MySQL
)

type DB struct {
    *sql.DB
}

func NewConnection(driver, dsn string) (*DB, error) {
    db, err := sql.Open(driver, dsn)
    if err != nil {
        return nil, fmt.Errorf("failed to open database: %w", err)
    }
    
    // Configure connection pool
    db.SetMaxOpenConns(25)
    db.SetMaxIdleConns(5)
    db.SetConnMaxLifetime(5 * time.Minute)
    db.SetConnMaxIdleTime(5 * time.Minute)
    
    // Verify connection
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    
    if err := db.PingContext(ctx); err != nil {
        return nil, fmt.Errorf("failed to ping database: %w", err)
    }
    
    return &DB{db}, nil
}
```

#### Repository Pattern with Raw SQL
```go
package repository

import (
    "context"
    "database/sql"
    "fmt"
    "strings"
)

type UserRepository struct {
    db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
    return &UserRepository{db: db}
}

// Create user with prepared statement
func (r *UserRepository) Create(ctx context.Context, user *User) error {
    tx, err := r.db.BeginTx(ctx, nil)
    if err != nil {
        return err
    }
    defer tx.Rollback()
    
    // Insert user
    query := `
        INSERT INTO users (email, username, password, created_at, updated_at)
        VALUES ($1, $2, $3, NOW(), NOW())
        RETURNING id, created_at, updated_at
    `
    
    err = tx.QueryRowContext(ctx, query,
        user.Email,
        user.Username,
        user.Password,
    ).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
    
    if err != nil {
        return fmt.Errorf("failed to insert user: %w", err)
    }
    
    // Insert profile
    profileQuery := `
        INSERT INTO profiles (user_id, first_name, last_name, bio, avatar)
        VALUES ($1, $2, $3, $4, $5)
        RETURNING id
    `
    
    err = tx.QueryRowContext(ctx, profileQuery,
        user.ID,
        user.Profile.FirstName,
        user.Profile.LastName,
        user.Profile.Bio,
        user.Profile.Avatar,
    ).Scan(&user.Profile.ID)
    
    if err != nil {
        return fmt.Errorf("failed to insert profile: %w", err)
    }
    
    return tx.Commit()
}

// Find by ID with JOIN
func (r *UserRepository) FindByID(ctx context.Context, id int64) (*User, error) {
    query := `
        SELECT 
            u.id, u.email, u.username, u.created_at, u.updated_at,
            p.id, p.first_name, p.last_name, p.bio, p.avatar
        FROM users u
        LEFT JOIN profiles p ON p.user_id = u.id
        WHERE u.id = $1 AND u.deleted_at IS NULL
    `
    
    var user User
    var profile Profile
    
    err := r.db.QueryRowContext(ctx, query, id).Scan(
        &user.ID, &user.Email, &user.Username, &user.CreatedAt, &user.UpdatedAt,
        &profile.ID, &profile.FirstName, &profile.LastName, &profile.Bio, &profile.Avatar,
    )
    
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, ErrUserNotFound
        }
        return nil, err
    }
    
    user.Profile = profile
    return &user, nil
}

// Complex query with dynamic conditions
func (r *UserRepository) Search(ctx context.Context, filters SearchFilters) ([]User, error) {
    query := strings.Builder{}
    args := []interface{}{}
    argCount := 0
    
    query.WriteString(`
        SELECT u.id, u.email, u.username, u.created_at
        FROM users u
        WHERE u.deleted_at IS NULL
    `)
    
    // Dynamic filters
    if filters.Email != "" {
        argCount++
        query.WriteString(fmt.Sprintf(" AND u.email LIKE $%d", argCount))
        args = append(args, "%"+filters.Email+"%")
    }
    
    if filters.Username != "" {
        argCount++
        query.WriteString(fmt.Sprintf(" AND u.username LIKE $%d", argCount))
        args = append(args, "%"+filters.Username+"%")
    }
    
    if filters.CreatedAfter != nil {
        argCount++
        query.WriteString(fmt.Sprintf(" AND u.created_at > $%d", argCount))
        args = append(args, filters.CreatedAfter)
    }
    
    // Sorting
    if filters.SortBy != "" {
        query.WriteString(fmt.Sprintf(" ORDER BY %s %s", filters.SortBy, filters.SortOrder))
    }
    
    // Pagination
    if filters.Limit > 0 {
        argCount++
        query.WriteString(fmt.Sprintf(" LIMIT $%d", argCount))
        args = append(args, filters.Limit)
        
        if filters.Offset > 0 {
            argCount++
            query.WriteString(fmt.Sprintf(" OFFSET $%d", argCount))
            args = append(args, filters.Offset)
        }
    }
    
    rows, err := r.db.QueryContext(ctx, query.String(), args...)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    
    var users []User
    for rows.Next() {
        var user User
        if err := rows.Scan(&user.ID, &user.Email, &user.Username, &user.CreatedAt); err != nil {
            return nil, err
        }
        users = append(users, user)
    }
    
    return users, rows.Err()
}

// Bulk operations
func (r *UserRepository) BulkInsert(ctx context.Context, users []User) error {
    tx, err := r.db.BeginTx(ctx, nil)
    if err != nil {
        return err
    }
    defer tx.Rollback()
    
    stmt, err := tx.PrepareContext(ctx, `
        INSERT INTO users (email, username, password, created_at, updated_at)
        VALUES ($1, $2, $3, NOW(), NOW())
    `)
    if err != nil {
        return err
    }
    defer stmt.Close()
    
    for _, user := range users {
        _, err := stmt.ExecContext(ctx, user.Email, user.Username, user.Password)
        if err != nil {
            return err
        }
    }
    
    return tx.Commit()
}
```

#### Advanced SQL Patterns
```go
// Window functions (PostgreSQL)
func (r *UserRepository) GetUserRankings(ctx context.Context) ([]UserRanking, error) {
    query := `
        WITH user_stats AS (
            SELECT 
                u.id,
                u.username,
                COUNT(p.id) as post_count,
                RANK() OVER (ORDER BY COUNT(p.id) DESC) as rank,
                LAG(COUNT(p.id), 1) OVER (ORDER BY COUNT(p.id) DESC) as prev_count
            FROM users u
            LEFT JOIN posts p ON p.author_id = u.id
            GROUP BY u.id, u.username
        )
        SELECT id, username, post_count, rank, 
               COALESCE(post_count - prev_count, 0) as difference
        FROM user_stats
        ORDER BY rank
        LIMIT 10
    `
    
    rows, err := r.db.QueryContext(ctx, query)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    
    var rankings []UserRanking
    for rows.Next() {
        var r UserRanking
        err := rows.Scan(&r.ID, &r.Username, &r.PostCount, &r.Rank, &r.Difference)
        if err != nil {
            return nil, err
        }
        rankings = append(rankings, r)
    }
    
    return rankings, rows.Err()
}

// JSON operations (PostgreSQL)
func (r *UserRepository) UpdateUserSettings(ctx context.Context, userID int64, settings map[string]interface{}) error {
    settingsJSON, err := json.Marshal(settings)
    if err != nil {
        return err
    }
    
    _, err = r.db.ExecContext(ctx, `
        UPDATE users 
        SET settings = settings || $1::jsonb,
            updated_at = NOW()
        WHERE id = $2
    `, settingsJSON, userID)
    
    return err
}

// Common Table Expressions (CTE)
func (r *UserRepository) GetUserHierarchy(ctx context.Context, userID int64) ([]User, error) {
    query := `
        WITH RECURSIVE user_tree AS (
            SELECT id, username, manager_id, 0 as level
            FROM users
            WHERE id = $1
            
            UNION ALL
            
            SELECT u.id, u.username, u.manager_id, ut.level + 1
            FROM users u
            INNER JOIN user_tree ut ON u.manager_id = ut.id
        )
        SELECT id, username, manager_id, level
        FROM user_tree
        ORDER BY level
    `
    
    rows, err := r.db.QueryContext(ctx, query, userID)
    // ... handle rows
}
```

### Raw SQL Performance Tips

1. **Use Prepared Statements**
```go
// Prepare once, execute many times
stmt, err := db.Prepare("SELECT * FROM users WHERE email = $1")
defer stmt.Close()

for _, email := range emails {
    rows, err := stmt.Query(email)
    // Process rows
}
```

2. **Batch Operations**
```go
// Use COPY for bulk inserts (PostgreSQL)
func BulkInsertWithCopy(db *sql.DB, users []User) error {
    txn, err := db.Begin()
    if err != nil {
        return err
    }
    
    stmt, err := txn.Prepare(pq.CopyIn("users", "email", "username", "password"))
    if err != nil {
        return err
    }
    
    for _, user := range users {
        _, err = stmt.Exec(user.Email, user.Username, user.Password)
        if err != nil {
            return err
        }
    }
    
    _, err = stmt.Exec()
    if err != nil {
        return err
    }
    
    return txn.Commit()
}
```

3. **Connection Pooling**
```go
// Monitor pool stats
stats := db.Stats()
log.Printf("Open connections: %d", stats.OpenConnections)
log.Printf("In use: %d", stats.InUse)
log.Printf("Idle: %d", stats.Idle)
```

## GORM vs Raw SQL Comparison

### Performance Comparison

| Operation | GORM | Raw SQL | Notes |
|-----------|------|---------|-------|
| Simple INSERT | ~2ms | ~1ms | GORM adds minimal overhead |
| Complex JOIN | ~15ms | ~8ms | Raw SQL more efficient |
| Bulk INSERT (1000 rows) | ~200ms | ~50ms | Raw SQL with COPY is fastest |
| Simple SELECT | ~3ms | ~2ms | Similar performance |
| Update with conditions | ~5ms | ~3ms | GORM generates optimal SQL |

### Code Comparison

#### Simple Query
```go
// GORM
var user User
db.First(&user, "email = ?", "user@example.com")

// Raw SQL
var user User
row := db.QueryRow("SELECT * FROM users WHERE email = $1", "user@example.com")
err := row.Scan(&user.ID, &user.Email, &user.Username, /*...*/)
```

#### Complex Query
```go
// GORM
var results []struct {
    UserID    uint
    PostCount int64
}
db.Model(&User{}).
    Select("users.id as user_id, count(posts.id) as post_count").
    Joins("LEFT JOIN posts ON posts.author_id = users.id").
    Group("users.id").
    Having("count(posts.id) > ?", 5).
    Find(&results)

// Raw SQL
rows, err := db.Query(`
    SELECT u.id as user_id, COUNT(p.id) as post_count
    FROM users u
    LEFT JOIN posts p ON p.author_id = u.id
    GROUP BY u.id
    HAVING COUNT(p.id) > $1
`, 5)
// ... scan rows
```

## Migration Strategies

### GORM to Raw SQL Migration

If you need to migrate from GORM to Raw SQL:

1. **Keep Models**: Use struct tags for both
```go
type User struct {
    ID       int64     `gorm:"primaryKey" db:"id"`
    Email    string    `gorm:"uniqueIndex" db:"email"`
    Username string    `gorm:"uniqueIndex" db:"username"`
}
```

2. **Gradual Migration**: Implement repository interface
```go
type UserRepository interface {
    Create(ctx context.Context, user *User) error
    FindByID(ctx context.Context, id int64) (*User, error)
    Update(ctx context.Context, user *User) error
    Delete(ctx context.Context, id int64) error
}

// Can swap implementations
var repo UserRepository = NewGORMUserRepository(db)
// Later...
var repo UserRepository = NewSQLUserRepository(db)
```

### Raw SQL to GORM Migration

If you need to migrate from Raw SQL to GORM:

1. **Map existing tables**:
```go
type User struct {
    ID       uint   `gorm:"column:user_id;primaryKey"`
    Email    string `gorm:"column:email_address"`
    Username string `gorm:"column:user_name"`
}

func (User) TableName() string {
    return "app_users" // Custom table name
}
```

2. **Custom SQL in GORM**:
```go
// Use raw SQL when needed
var users []User
db.Raw("SELECT * FROM users WHERE custom_function(email) = ?", value).Scan(&users)
```

## Best Practices

### General Database Best Practices

1. **Always use contexts**:
```go
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()
db.WithContext(ctx).Find(&users)
```

2. **Handle errors properly**:
```go
if err != nil {
    if errors.Is(err, gorm.ErrRecordNotFound) || err == sql.ErrNoRows {
        return nil, ErrNotFound
    }
    logger.Error("Database error", "error", err)
    return nil, ErrInternal
}
```

3. **Use transactions for consistency**:
```go
err := db.Transaction(func(tx *gorm.DB) error {
    // All operations in transaction
    return nil
})
```

4. **Index your queries**:
```sql
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_posts_author_created ON posts(author_id, created_at DESC);
```

### GORM-Specific Best Practices

1. **Disable default transaction** for read queries:
```go
db.Session(&gorm.Session{SkipDefaultTransaction: true}).Find(&users)
```

2. **Use scopes** for common queries:
```go
func Active(db *gorm.DB) *gorm.DB {
    return db.Where("status = ?", "active")
}

db.Scopes(Active).Find(&users)
```

3. **Avoid Select(*)**:
```go
// Bad
db.Select("*").Find(&users)

// Good
db.Select("id", "email", "username").Find(&users)
```

### Raw SQL Best Practices

1. **Use named parameters** (with sqlx):
```go
query := `SELECT * FROM users WHERE email = :email AND status = :status`
rows, err := db.NamedQuery(query, map[string]interface{}{
    "email": "user@example.com",
    "status": "active",
})
```

2. **Close resources**:
```go
rows, err := db.Query(query)
if err != nil {
    return err
}
defer rows.Close() // Always close rows
```

3. **Check for iteration errors**:
```go
for rows.Next() {
    // ... scan
}
if err := rows.Err(); err != nil {
    return err
}
```

## Performance Monitoring

### Query Performance Monitoring

```go
// Middleware for timing queries
func TimeQuery(name string, fn func() error) error {
    start := time.Now()
    err := fn()
    duration := time.Since(start)
    
    if duration > 100*time.Millisecond {
        logger.Warn("Slow query detected", 
            "name", name,
            "duration", duration,
        )
    }
    
    metrics.RecordQueryDuration(name, duration)
    return err
}

// Usage
err := TimeQuery("user.findByID", func() error {
    return db.First(&user, id).Error
})
```

### Database Metrics

```go
// Regular monitoring
func MonitorDatabase(db *sql.DB) {
    ticker := time.NewTicker(30 * time.Second)
    defer ticker.Stop()
    
    for range ticker.C {
        stats := db.Stats()
        
        metrics.SetGauge("db.connections.open", float64(stats.OpenConnections))
        metrics.SetGauge("db.connections.in_use", float64(stats.InUse))
        metrics.SetGauge("db.connections.idle", float64(stats.Idle))
        metrics.SetGauge("db.wait_count", float64(stats.WaitCount))
        metrics.SetGauge("db.wait_duration", stats.WaitDuration.Seconds())
    }
}
```

## Conclusion

Choose your ORM based on:

1. **Project Requirements**: Performance needs, query complexity
2. **Team Experience**: SQL knowledge, ORM familiarity
3. **Development Speed**: Time to market vs optimization needs
4. **Maintenance**: Long-term maintainability and debugging

Remember: You can always start with GORM for rapid development and optimize specific queries with raw SQL where needed. The generated code structure supports both approaches seamlessly.