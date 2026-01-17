# 代码评审规范

## 1. SOLID 原则检查清单

### 1.1 Single Responsibility Principle (单一职责原则)

**原则定义**：一个类或模块应该有且只有一个改变的理由。

#### 评审要点
- 每个结构体/类是否只负责一个功能领域？
- 方法是否只做一件事？
- 是否有将多个职责混在一起的情况？

#### ✅ 好的做法：职责清晰分离

```go
// UserValidator 只负责验证逻辑
type UserValidator struct{}

func (v *UserValidator) Validate CreateUser(req request.CreateUserRequest) error {
    if len(req.Username) < 3 {
        return ErrInvalidUsername
    }
    return nil
}

// UserRepository 只负责数据持久化
type UserRepository struct {
    db *gorm.DB
}

func (r *UserRepository) Save(ctx context.Context, user *domain.User) error {
    return r.db.Save(user).Error
}

// UserService 只负责业务逻辑编排
type UserService struct {
    validator *UserValidator
    repo      *UserRepository
}

func (s *UserService) CreateUser(ctx context.Context, req request.CreateUserRequest) (*domain.User, error) {
    if err := s.validator.ValidateCreateUser(req); err != nil {
        return nil, err
    }

    user := &domain.User{Username: req.Username}
    if err := s.repo.Save(ctx, user); err != nil {
        return nil, err
    }

    return user, nil
}
```

#### ❌ 避免：职责混乱

```go
// ❌ 一个结构体承担了多个职责：验证、持久化、日志、通知
type UserService struct {
    db *gorm.DB
}

func (s *UserService) CreateUser(req request.CreateUserRequest) error {
    // 职责1：验证
    if len(req.Username) < 3 {
        return ErrInvalidUsername
    }

    // 职责2：数据库操作
    user := &User{Username: req.Username}
    if err := s.db.Create(user).Error; err != nil {
        return err
    }

    // 职责3：日志记录
    log.Printf("User created: %d", user.ID)

    // 职责4：发送通知
    sendEmail(user.Email, "Welcome!")

    return nil
}
```

#### 评审检查项
- [ ] Service 层是否只包含业务逻辑编排？
- [ ] Repository 层是否只包含数据访问？
- [ ] Handler 层是否只包含 HTTP 处理逻辑？
- [ ] 是否有"上帝对象"（God Object）做太多事情？

---

### 1.2 Open/Closed Principle (开闭原则)

**原则定义**：软件实体应该对扩展开放，对修改关闭。

#### 评审要点
- 添加新功能时是否需要修改已有代码？
- 是否使用接口和抽象来支持扩展？
- 是否使用策略模式、模板方法模式等支持扩展？

#### ✅ 好的做法：通过接口扩展

```go
// 定义支付策略接口
type PaymentStrategy interface {
    Pay(ctx context.Context, amount float64) error
}

// 各种支付实现
type AlipayStrategy struct{}
func (s *AlipayStrategy) Pay(ctx context.Context, amount float64) error {
    // 支付宝支付逻辑
    return nil
}

type WechatPayStrategy struct{}
func (s *WechatPayStrategy) Pay(ctx context.Context, amount float64) error {
    // 微信支付逻辑
    return nil
}

type CreditCardStrategy struct{}
func (s *CreditCardStrategy) Pay(ctx context.Context, amount float64) error {
    // 信用卡支付逻辑
    return nil
}

// PaymentService 使用接口，对扩展开放
type PaymentService struct {
    strategies map[string]PaymentStrategy
}

func (s *PaymentService) RegisterStrategy(name string, strategy PaymentStrategy) {
    s.strategies[name] = strategy
}

func (s *PaymentService) Pay(ctx context.Context, method string, amount float64) error {
    strategy, ok := s.strategies[method]
    if !ok {
        return ErrUnsupportedPaymentMethod
    }
    return strategy.Pay(ctx, amount)
}

// ✅ 添加新的支付方式时，只需新增实现，无需修改 PaymentService
```

#### ❌ 避免：修改现有代码添加功能

```go
// ❌ 每次添加新支付方式都需要修改 Pay 方法
type PaymentService struct{}

func (s *PaymentService) Pay(method string, amount float64) error {
    switch method {
    case "alipay":
        return s.payAlipay(amount)
    case "wechat":
        return s.payWechat(amount)
    case "credit_card":
        return s.payCreditCard(amount)
    // ❌ 添加新支付方式需要修改这里，违反开闭原则
    case "new_payment":
        return s.payNewPayment(amount)
    default:
        return ErrUnsupportedPaymentMethod
    }
}
```

#### ✅ 好的做法：使用依赖注入支持扩展

```go
// NotificationSender 接口
type NotificationSender interface {
    Send(ctx context.Context, msg Message) error
}

// UserService 通过依赖注入使用 NotificationSender
type UserService struct {
    notificationSender NotificationSender
}

func NewUserService(sender NotificationSender) *UserService {
    return &UserService{
        notificationSender: sender,
    }
}

// ✅ 可以注入不同的实现：EmailSender、SMSSender、PushSender 等
```

#### 评审检查项
- [ ] 添加新功能时是否通过新增代码而非修改现有代码？
- [ ] 是否使用接口来定义抽象？
- [ ] switch/if-else 是否可以通过策略模式消除？
- [ ] 配置项是否可以外部化而非硬编码？

---

### 1.3 Liskov Substitution Principle (里氏替换原则)

**原则定义**：子类对象应该可以替换父类对象而不影响程序正确性。

#### 评审要点
- 子类是否实现了父类的所有约定？
- 子类是否加强了前置条件或减弱了后置条件？
- 子类是否抛出了父类未声明的异常？

#### ✅ 好的做法：完全遵循接口契约

```go
// Repository 接口定义数据访问契约
type Repository interface {
    FindByID(ctx context.Context, id int) (interface{}, error)
    Save(ctx context.Context, entity interface{}) error
}

// MySQLRepository 完全实现接口
type MySQLRepository struct {
    db *gorm.DB
}

func (r *MySQLRepository) FindByID(ctx context.Context, id int) (interface{}, error) {
    var entity Entity
    err := r.db.WithContext(ctx).First(&entity, id).Error
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, ErrNotFound // 遵循契约，返回约定错误
        }
        return nil, err
    }
    return &entity, nil
}

// ✅ 其他实现（PostgreSQLRepository、MongoRepository）都可以替换使用
```

#### ❌ 避免：违反接口契约

```go
// ❌ 子类改变了接口的语义
type CachedRepository struct {
    repo Repository
    cache Cache
}

func (r *CachedRepository) FindByID(ctx context.Context, id int) (interface{}, error) {
    // ❌ 问题1：缓存未命中时不返回错误，返回 nil, nil
    // 这与原接口契约不同（原接口返回 error 表示查找失败）
    if val := r.cache.Get(id); val != nil {
        return val, nil
    }

    // ❌ 问题2：静默失败，不返回错误
    entity, err := r.repo.FindByID(ctx, id)
    if err != nil {
        return nil, nil // 违反契约
    }

    return entity, nil
}
```

#### ✅ 好的做法：矩形继承关系

```go
// Bird 接口定义鸟类行为
type Bird interface {
    Fly() error
}

// ✅ 所有鸟类都能飞，可以安全替换
type Sparrow struct{}
func (s *Sparrow) Fly() error { return nil }

type Eagle struct{}
func (e *Eagle) Fly() error { return nil }

// ❌ 企鹅不能飞，不应该实现 Bird 接口
// type Penguin struct{}
// func (p *Penguin) Fly() error { return ErrCannotFly }

// ✅ 正确做法：重新设计接口层次
type Bird interface{}

type FlyingBird interface {
    Bird
    Fly() error
}

type SwimmingBird interface {
    Bird
    Swim() error
}
```

#### 评审检查项
- [ ] 接口的所有实现是否遵循相同的契约？
- [ ] 子类是否没有"偷懒"实现（空实现、返回默认值）？
- [ ] 接口方法的前置/后置条件是否一致？
- [ ] 是否存在逻辑上不该实现某接口的类型？

---

### 1.4 Interface Segregation Principle (接口隔离原则)

**原则定义**：客户端不应该依赖它不需要的接口。

#### 评审要点
- 接口是否过于臃肿？
- 客户端是否被迫实现不需要的方法？
- 接口是否按照职责拆分？

#### ✅ 好的做法：接口职责单一

```go
// ✅ 将大的 Repository 接口拆分为多个小接口

// 只读操作接口
type Reader interface {
    FindByID(ctx context.Context, id int) (*Entity, error)
    FindAll(ctx context.Context) ([]*Entity, error)
    Query(ctx context.Context, criteria QueryCriteria) ([]*Entity, error)
}

// 只写操作接口
type Writer interface {
    Create(ctx context.Context, entity *Entity) error
    Update(ctx context.Context, entity *Entity) error
    Delete(ctx context.Context, id int) error
}

// 批量操作接口
type BatchWriter interface {
    CreateBatch(ctx context.Context, entities []*Entity) error
    UpdateBatch(ctx context.Context, entities []*Entity) error
    DeleteBatch(ctx context.Context, ids []int) error
}

// ✅ 客户端按需依赖接口
type ReadOnlyService struct {
    repo Reader // 只需要读取功能
}

type WriteOnlyService struct {
    repo Writer // 只需要写入功能
}

type FullService struct {
    repo ReaderWriter // 需要读写功能
}
```

#### ❌ 避免：臃肿的接口

```go
// ❌ 一个接口包含太多方法
type UserRepository interface {
    FindByID(ctx context.Context, id int) (*User, error)
    FindByName(ctx context.Context, name string) (*User, error)
    FindByEmail(ctx context.Context, email string) (*User, error)
    Create(ctx context.Context, user *User) error
    Update(ctx context.Context, user *User) error
    Delete(ctx context.Context, id int) error
    UpdatePassword(ctx context.Context, id int, password string) error
    UpdateEmail(ctx context.Context, id int, email string) error
    UpdateProfile(ctx context.Context, id int, profile Profile) error
    FindByRole(ctx context.Context, role string) ([]*User, error)
    FindActiveUsers(ctx context.Context) ([]*User, error)
    FindInactiveUsers(ctx context.Context) ([]*User, error)
    // ... 太多方法
}

// ❌ 客户端被迫实现不需要的方法
type MockUserRepositoryForSearch struct {
    // 只需要查询功能，但必须实现所有方法
    mock.Mock
}

func (m *MockUserRepositoryForSearch) Create(ctx context.Context, user *User) error {
    // ❌ 不需要这个方法，但必须实现
    return nil
}
```

#### ✅ 好的做法：接口按角色拆分

```go
// 普通用户只需要基本操作
type UserReader interface {
    GetProfile(userID int) (*Profile, error)
}

// 管理员需要管理操作
type UserAdmin interface {
    UserReader
    CreateUser(user *User) error
    DeleteUser(userID int) error
    UpdateUserStatus(userID int, status Status) error
}

// 审计员只需要审计功能
type UserAuditor interface {
    GetUserActivity(userID int) ([]Activity, error)
    GetLoginHistory(userID int) ([]LoginRecord, error)
}
```

#### 评审检查项
- [ ] 接口方法数量是否合理（建议 < 10 个）？
- [ ] 是否有"胖接口"包含太多不相关的方法？
- [ ] 客户端是否被迫实现空方法？
- [ ] 接口是否按照使用场景拆分？

---

### 1.5 Dependency Inversion Principle (依赖倒置原则)

**原则定义**：
1. 高层模块不应依赖低层模块，二者都应依赖抽象
2. 抽象不应依赖细节，细节应依赖抽象

#### 评审要点
- 高层模块是否直接依赖具体实现？
- 是否通过接口进行依赖注入？
- 依赖关系是否指向抽象？

#### ✅ 好的做法：依赖抽象

```go
// Domain 层：定义抽象接口（高层）
type UserRepository interface {
    FindByID(ctx context.Context, id int) (*User, error)
    Save(ctx context.Context, user *User) error
}

type UserService struct {
    // ✅ 依赖接口抽象，而非具体实现
    repo UserRepository
}

func NewUserService(repo UserRepository) *UserService {
    return &UserService{
        repo: repo, // 依赖注入
    }
}

func (s *UserService) GetUser(ctx context.Context, id int) (*User, error) {
    // ✅ 调用接口方法，不关心具体实现
    return s.repo.FindByID(ctx, id)
}

// Infrastructure 层：实现接口（低层）
type MySQLUserRepository struct {
    db *gorm.DB
}

func (r *MySQLUserRepository) FindByID(ctx context.Context, id int) (*User, error) {
    // MySQL 具体实现
    return &User{}, nil
}

// ✅ 依赖方向：UserService -> UserRepository <- MySQLUserRepository
```

#### ❌ 避免：高层依赖低层

```go
// ❌ Service 层直接依赖 Infrastructure 层的具体实现
type UserService struct {
    // ❌ 直接依赖 MySQL，耦合度太高
    db *gorm.DB
}

func NewUserService(db *gorm.DB) *UserService {
    return &UserService{
        db: db,
    }
}

func (s *UserService) GetUser(ctx context.Context, id int) (*User, error) {
    // ❌ 直接使用 GORM API，难以替换为其他数据库
    var user User
    err := s.db.WithContext(ctx).First(&user, id).Error
    return &user, err
}
```

#### ✅ 好的做法：使用工厂和依赖注入容器

```go
// 定义接口
type Config interface {
    Get(key string) string
}

type Logger interface {
    Info(msg string, fields ...Field)
    Error(msg string, fields ...Field)
}

type Database interface {
    Query(ctx context.Context, query string, args ...interface{}) (Result, error)
}

// ✅ 高层模块依赖接口
type Service struct {
    config   Config
    logger   Logger
    database Database
}

func NewService(config Config, logger Logger, database Database) *Service {
    return &Service{
        config:   config,
        logger:   logger,
        database: database,
    }
}

// ✅ 在 main 函数中组装依赖（应用入口）
func main() {
    // 创建具体实现
    config := NewEnvConfig()
    logger := NewZapLogger()
    database := NewPostgresDB(config.Get("DB_URL"))

    // 注入到高层模块
    service := NewService(config, logger, database)

    // ✅ 可以轻松替换实现
    // logger := NewConsoleLogger()
    // database := NewMySQLDB(config.Get("DB_URL"))
}
```

#### ✅ 好的做法：DDD 分层中的依赖倒置

```go
// Domain 层（最底层，不依赖任何层）
package domain

type UserRepository interface { // 领域接口
    FindByID(ctx context.Context, id int) (*User, error)
}

type UserService struct { // 领域服务
    repo UserRepository
}

// Application 层（依赖 Domain）
package application

type UserApplicationService struct {
    userRepo domain.UserRepository // 依赖领域接口
}

// Infrastructure 层（依赖 Domain）
package infrastructure

type MySQLUserRepository struct { // 实现领域接口
    db *gorm.DB
}

func (r *MySQLUserRepository) FindByID(ctx context.Context, id int) (*domain.User, error) {
    // 实现...
}

// ✅ 依赖关系：
// Application -> Domain <- Infrastructure
```

#### 评审检查项
- [ ] 高层模块是否只依赖接口？
- [ ] 是否通过构造函数注入依赖？
- [ ] 是否可以在不修改高层代码的情况下替换低层实现？
- [ ] DDD 分层中，Domain 层是否不依赖其他层？

---

## 2. 代码质量评审规范

### 2.1 命名规范

#### ✅ 好的做法
```go
// 包名：小写、简洁、描述性
package user
package auth
package payment

// 接口名：动词+名词或名词
type UserReader interface {}
type PaymentProcessor interface {}

// 结构体：名词
type UserService struct {}
type UserRepository struct {}

// 常量：大驼峰或全大写
const MaxRetryCount = 3
type UserStatus string
const (
    UserStatusActive   UserStatus = "active"
    UserStatusInactive UserStatus = "inactive"
)

// 函数：动词开头
func CreateUser() {}
func ValidateInput() {}
func ProcessPayment() {}

// 变量：驼峰命名，语义清晰
var userCount int
var isActive bool
```

#### ❌ 避免
```go
// ❌ 包名过长或使用下划线
package user_service
package mypackage

// ❌ 接口名以 I 开头（C# 风格）
type IUserRepository interface {}

// ❌ 缩写不当
type Usr struct {}
type UsrSvc struct {}

// ❌ 变量名无意义
var a int
var tmp string
var data interface{}
```

#### 评审检查项
- [ ] 包名是否简洁且具有描述性？
- [ ] 接口是否避免使用 "I" 前缀？
- [ ] 常量是否使用大写或大驼峰？
- [ ] 变量名是否清晰表达用途？
- [ ] 是否避免单字母变量（除循环变量 i, j, k）？
- [ ] 是否避免使用缩写？

---

### 2.2 错误处理

#### ✅ 好的做法：完整的错误处理

```go
// ✅ 创建有意义的错误变量
var (
    ErrUserNotFound     = errors.New("user not found")
    ErrInvalidInput     = errors.New("invalid input")
    ErrUnauthorized     = errors.New("unauthorized access")
)

// ✅ 包装错误上下文
func (s *UserService) GetUser(id int) (*User, error) {
    user, err := s.repo.FindByID(id)
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, fmt.Errorf("get user by id %d: %w", id, ErrUserNotFound)
        }
        return nil, fmt.Errorf("get user by id %d: %w", id, err)
    }
    return user, nil
}

// ✅ 调用者检查错误
user, err := userService.GetUser(userID)
if err != nil {
    if errors.Is(err, ErrUserNotFound) {
        return c.JSON(404, ErrorResponse{Message: "User not found"})
    }
    return c.JSON(500, ErrorResponse{Message: "Internal server error"})
}

// ✅ 自定义错误类型
type ValidationError struct {
    Field   string
    Message string
}

func (e *ValidationError) Error() string {
    return fmt.Sprintf("validation failed for field %s: %s", e.Field, e.Message)
}
```

#### ❌ 避免：忽略或不当处理错误

```go
// ❌ 忽略错误
user, _ := userService.GetUser(userID)

// ❌ 吞掉错误
func (s *UserService) GetUser(id int) *User {
    user, err := s.repo.FindByID(id)
    if err != nil {
        return nil // 丢失错误信息
    }
    return user
}

// ❌ 使用 panic 处理业务错误
func (s *UserService) CreateUser(req CreateUserRequest) *User {
    if req.Username == "" {
        panic("username is required") // ❌ 不要使用 panic
    }
    // ...
}

// ❌ 创建错误但不使用
func (s *UserService) DeleteUser(id int) error {
    if err := s.repo.Delete(id); err != nil {
        errors.New("delete failed") // ❌ 应该包装原始错误
        return err
    }
    return nil
}
```

#### 评审检查项
- [ ] 所有错误是否都被检查？
- [ ] 错误是否包含有用的上下文信息？
- [ ] 是否使用 `%w` 包装错误而不是 `%v`？
- [ ] 是否避免使用 panic 处理业务错误？
- [ ] 自定义错误类型是否实现了 Error() 方法？

---

### 2.3 并发安全

#### ✅ 好的做法：并发安全的设计

```go
// ✅ 使用 sync.Mutex 保护共享状态
type SafeCounter struct {
    mu    sync.RWMutex
    count int
}

func (c *SafeCounter) Increment() {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.count++
}

func (c *SafeCounter) Get() int {
    c.mu.RLock()
    defer c.mu.RUnlock()
    return c.count
}

// ✅ 使用通道进行通信
type WorkerPool struct {
    tasks   chan Task
    workers int
    wg      sync.WaitGroup
}

func NewWorkerPool(workers int) *WorkerPool {
    pool := &WorkerPool{
        tasks:   make(chan Task, 100),
        workers: workers,
    }
    pool.start()
    return pool
}

func (p *WorkerPool) Submit(task Task) {
    p.tasks <- task
}

// ✅ 使用 context 管理生命周期
func (s *Service) ProcessWithTimeout(ctx context.Context) error {
    ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
    defer cancel()

    resultChan := make(chan Result, 1)

    go func() {
        resultChan <- s.process()
    }()

    select {
    case result := <-resultChan:
        return result.Err
    case <-ctx.Done():
        return ctx.Err()
    }
}
```

#### ❌ 避免：并发安全问题

```go
// ❌ 数据竞争
type UnsafeCounter struct {
    count int
}

func (c *UnsafeCounter) Increment() {
    c.count++ // ❌ 并发不安全
}

// ❌ 死锁风险
func deadlock() {
    m1 := sync.Mutex{}
    m2 := sync.Mutex{}

    go func() {
        m1.Lock()
        time.Sleep(100 * time.Millisecond)
        m2.Lock()
    }()

    m2.Lock()
    time.Sleep(100 * time.Millisecond)
    m1.Lock() // ❌ 可能死锁
}

// ❌ Goroutine 泄漏
func leak() {
    ch := make(chan int)

    go func() {
        val := <-ch
        fmt.Println(val)
    }()

    // ❌ 没有往 ch 发送数据，goroutine 永远阻塞
}
```

#### 评审检查项
- [ ] 共享状态是否使用 mutex 保护？
- [ ] 是否避免了数据竞争？
- [ ] 是否有死锁风险？
- [ ] Goroutine 是否会泄漏？
- [ ] 是否使用 context 管理 goroutine 生命周期？

---

### 2.4 资源管理

#### ✅ 好的做法：正确管理资源

```go
// ✅ 使用 defer 确保资源释放
func ProcessFile(filepath string) error {
    file, err := os.Open(filepath)
    if err != nil {
        return err
    }
    defer file.Close() // 确保文件被关闭

    // 处理文件...
    return nil
}

// ✅ 使用 defer 清理临时资源
func DownloadAndProcess(url string) error {
    tmpFile, err := os.CreateTemp("", "download")
    if err != nil {
        return err
    }
    defer os.Remove(tmpFile.Name()) // 确保删除临时文件

    // 下载和处理...
    return nil
}

// ✅ 数据库连接池管理
func NewDBPool(config *Config) (*sql.DB, error) {
    db, err := sql.Open("postgres", config.DSN)
    if err != nil {
        return nil, err
    }

    // 设置连接池参数
    db.SetMaxOpenConns(25)
    db.SetMaxIdleConns(5)
    db.SetConnMaxLifetime(5 * time.Minute)

    // 验证连接
    if err := db.Ping(); err != nil {
        db.Close()
        return nil, err
    }

    return db, nil
}
```

#### ❌ 避免：资源泄漏

```go
// ❌ 文件未关闭
func ReadFile(filepath string) ([]byte, error) {
    file, err := os.Open(filepath)
    if err != nil {
        return nil, err
    }
    // ❌ 忘记 defer file.Close()

    return ioutil.ReadAll(file)
}

// ❌ HTTP 响应体未关闭
func FetchData(url string) ([]byte, error) {
    resp, err := http.Get(url)
    if err != nil {
        return nil, err
    }
    // ❌ 忘记 defer resp.Body.Close()

    return ioutil.ReadAll(resp.Body)
}

// ❌ 连接未释放
func QueryData() error {
    conn, err := grpc.Dial(addr, opts)
    if err != nil {
        return err
    }
    // ❌ 忘记 defer conn.Close()

    client := NewClient(conn)
    return client.DoSomething()
}
```

#### 评审检查项
- [ ] 所有打开的文件是否使用 defer 关闭？
- [ ] HTTP 响应体是否关闭？
- [ ] 数据库连接是否正确释放？
- [ ] 临时文件是否删除？
- [ ] 锁是否在 defer 中释放？

---

## 3. DDD 架构评审规范

### 3.1 分层架构检查

#### ✅ 好的做法：清晰的分层结构

```
project/
├── domain/              # 领域层（最底层，不依赖其他层）
│   ├── entity/         # 实体
│   ├── valueobject/    # 值对象
│   ├── repository/     # 仓储接口
│   └── service/        # 领域服务
├── application/         # 应用层（依赖 domain）
│   ├── service/        # 应用服务
│   └── dto/            # 数据传输对象
├── infrastructure/      # 基础设施层（依赖 domain）
│   ├── persistence/    # 数据持久化
│   ├── messaging/      # 消息队列
│   └── cache/          # 缓存实现
└── interfaces/          # 接口层（依赖 application）
    ├── http/           # HTTP 控制器
    └── grpc/           # gRPC 服务
```

#### 依赖规则
```
interfaces -> application -> domain
infrastructure -> domain
```

#### 评审检查项
- [ ] 是否遵循依赖方向（上层依赖下层）？
- [ ] Domain 层是否不依赖任何其他层？
- [ ] 是否有跨层直接调用的情况？
- [ ] 层与层之间是否通过接口交互？

---

### 3.2 领域模型评审

#### ✅ 好的做法：充血模型

```go
// ✅ User 是充血模型，包含业务逻辑
type User struct {
    ID       int
    Username string
    Password string
    Email    string
    Status   UserStatus
}

// 业务逻辑封装在实体内部
func (u *User) ChangeEmail(email string) error {
    if !isValidEmail(email) {
        return ErrInvalidEmail
    }
    u.Email = email
    return nil
}

func (u *User) Activate() error {
    if u.Status == StatusActive {
        return ErrAlreadyActive
    }
    u.Status = StatusActive
    return nil
}

func (u *User) ValidatePassword(password string) error {
    if !comparePassword(u.Password, password) {
        return ErrInvalidPassword
    }
    return nil
}
```

#### ❌ 避免：贫血模型

```go
// ❌ User 只是数据容器，没有业务逻辑
type User struct {
    ID       int
    Username string
    Password string
    Email    string
    Status   UserStatus
}

// ❌ 业务逻辑散落在 Service 层
func (s *UserService) ChangeEmail(userID int, email string) error {
    user, err := s.repo.FindByID(userID)
    if err != nil {
        return err
    }

    // ❌ 业务逻辑应该在 User 实体中
    if !isValidEmail(email) {
        return ErrInvalidEmail
    }
    user.Email = email

    return s.repo.Save(user)
}
```

#### 评审检查项
- [ ] 实体是否包含业务逻辑？
- [ ] 业务规则是否封装在领域模型中？
- [ ] 值对象是否不可变？
- [ ] 聚合根是否维护一致性边界？

---

## 4. 性能评审规范

### 4.1 数据库查询

#### ✅ 好的做法：高效的数据库操作

```go
// ✅ 批量查询，避免 N+1 问题
func (r *UserRepository) FindByIDs(ctx context.Context, ids []int) ([]*User, error) {
    var users []*User
    err := r.db.WithContext(ctx).Where("id IN ?", ids).Find(&users).Error
    return users, err
}

// ✅ 使用索引字段查询
func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*User, error) {
    var user User
    // email 应该有索引
    err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error
    return &user, err
}

// ✅ 分页查询
func (r *UserRepository) List(ctx context.Context, page, pageSize int) ([]*User, int64, error) {
    var users []*User
    var total int64

    err := r.db.WithContext(ctx).
        Model(&User{}).
        Count(&total).Error
    if err != nil {
        return nil, 0, err
    }

    err = r.db.WithContext(ctx).
        Offset((page - 1) * pageSize).
        Limit(pageSize).
        Find(&users).Error

    return users, total, err
}
```

#### ❌ 避免：性能问题

```go
// ❌ N+1 查询问题
func (r *UserRepository) FindWithOrders(userID int) (*User, error) {
    var user User
    r.db.First(&user, userID)

    // ❌ 在循环中查询，产生 N+1 问题
    var orders []Order
    r.db.Where("user_id = ?", user.ID).Find(&orders)

    return &user, nil
}

// ❌ 查询所有数据
func (r *UserRepository) ListAll() ([]*User, error) {
    var users []*User
    // ❌ 可能返回百万级数据
    err := r.db.Find(&users).Error
    return users, err
}

// ❌ 在循环中执行 SQL
func (s *UserService) UpdateUsers(userIDs []int, data UpdateData) error {
    for _, id := range userIDs {
        // ❌ 每次循环都执行一次 UPDATE
        s.db.Model(&User{}).Where("id = ?", id).Update(data)
    }
    return nil
}
```

#### 评审检查项
- [ ] 是否有 N+1 查询问题？
- [ ] 查询是否使用了索引？
- [ ] 是否有批量查询替代循环查询？
- [ ] 大数据量查询是否分页？
- [ ] 是否避免 SELECT * ？

---

### 4.2 缓存使用

#### ✅ 好的做法：合理使用缓存

```go
// ✅ 查询时使用缓存
func (s *UserService) GetUser(ctx context.Context, id int) (*User, error) {
    // 先查缓存
    cacheKey := fmt.Sprintf("user:%d", id)
    val, err := s.cache.Get(ctx, cacheKey)
    if err == nil {
        var user User
        if json.Unmarshal(val, &user) == nil {
            return &user, nil
        }
    }

    // 缓存未命中，查数据库
    user, err := s.repo.FindByID(ctx, id)
    if err != nil {
        return nil, err
    }

    // 写入缓存
    data, _ := json.Marshal(user)
    s.cache.Set(ctx, cacheKey, data, 5*time.Minute)

    return user, nil
}

// ✅ 更新时删除缓存
func (s *UserService) UpdateUser(ctx context.Context, user *User) error {
    if err := s.repo.Save(ctx, user); err != nil {
        return err
    }

    // 删除缓存，保证一致性
    cacheKey := fmt.Sprintf("user:%d", user.ID)
    s.cache.Delete(ctx, cacheKey)

    return nil
}
```

#### ❌ 避免：缓存滥用

```go
// ❌ 缓存所有数据，内存占用过高
func (s *UserService) ListAllUsers(ctx context.Context) ([]*User, error) {
    cacheKey := "all_users"

    // ❌ 缓存所有用户，可能百万级
    val, err := s.cache.Get(ctx, cacheKey)
    if err == nil {
        var users []*User
        json.Unmarshal(val, &users)
        return users, nil
    }

    // ...
}

// ❌ 缓存不一致
func (s *UserService) UpdateUser(user *User) error {
    if err := s.repo.Save(user); err != nil {
        return err
    }

    // ❌ 更新了数据库但没有删除缓存
    return nil
}
```

#### 评审检查项
- [ ] 是否只缓存热点数据？
- [ ] 缓存过期时间是否合理？
- [ ] 更新数据时是否删除/更新缓存？
- [ ] 是否考虑了缓存穿透、雪崩、击穿？

---

## 5. 安全评审规范

### 5.1 输入验证

#### ✅ 好的做法：严格的输入验证

```go
// ✅ 使用验证器
type CreateUserRequest struct {
    Username string `json:"username" validate:"required,min=3,max=20,alphanum"`
    Email    string `json:"email" validate:"required,email"`
    Password string `json:"password" validate:"required,min=8"`
    Age      int    `json:"age" validate:"gte=0,lte=150"`
}

func (s *UserService) CreateUser(ctx context.Context, req CreateUserRequest) error {
    // 验证输入
    if err := s.validator.Validate(req); err != nil {
        return ErrInvalidInput
    }

    // 业务验证
    if exists, _ := s.repo.ExistsByUsername(req.Username); exists {
        return ErrUsernameExists
    }

    // 创建用户...
    return nil
}
```

#### ❌ 避免：不验证输入

```go
// ❌ 不验证用户输入
func (s *UserService) CreateUser(req CreateUserRequest) error {
    // ❌ 直接使用用户输入，没有验证
    user := &User{
        Username: req.Username,
        Email:    req.Email,
        Password: req.Password, // ❌ 没有加密
    }
    return s.repo.Save(user)
}
```

#### 评审检查项
- [ ] 所有用户输入是否都经过验证？
- [ ] 是否使用白名单而非黑名单？
- [ ] SQL 查询是否使用参数化查询？
- [ ] 密码是否加密存储？

---

### 5.2 敏感信息保护

#### ✅ 好的做法：保护敏感信息

```go
// ✅ 敏感信息不记录到日志
func (s *UserService) Login(username, password string) error {
    // ✅ 不记录密码
    logger.Info("user login attempt",
        logger.String("username", username),
        // ❌ 不要记录: logger.String("password", password)
    )

    user, err := s.repo.FindByUsername(username)
    if err != nil {
        return err
    }

    if !user.CheckPassword(password) {
        return ErrInvalidCredentials
    }

    return nil
}

// ✅ 响应中过滤敏感字段
type UserResponse struct {
    ID       int    `json:"id"`
    Username string `json:"username"`
    Email    string `json:"email"`
    // Password 不导出
}

func ToResponse(user *User) UserResponse {
    return UserResponse{
        ID:       user.ID,
        Username: user.Username,
        Email:    user.Email,
        // 不包含密码
    }
}
```

#### 评审检查项
- [ ] 日志中是否包含密码、token 等敏感信息？
- [ ] API 响应是否过滤敏感字段？
- [ ] 是否使用 HTTPS 传输敏感数据？
- [ ] 密码是否使用哈希算法存储？

---

## 6. 测试评审规范

### 6.1 单元测试

#### ✅ 好的做法：完整的单元测试

```go
func TestUserService_CreateUser_Success(t *testing.T) {
    // Arrange
    mockRepo := new(MockUserRepository)
    service := NewUserService(mockRepo)

    req := request.CreateUserRequest{
        Username: "testuser",
        Email:    "test@example.com",
        Password: "password123",
    }

    mockRepo.On("Save", mock.Anything, mock.AnythingOfType("*domain.User")).Return(nil)

    // Act
    user, err := service.CreateUser(context.Background(), req)

    // Assert
    assert.NoError(t, err)
    assert.NotNil(t, user)
    assert.Equal(t, "testuser", user.Username)
    mockRepo.AssertExpectations(t)
}

func TestUserService_CreateUser_DuplicateUsername(t *testing.T) {
    // 测试错误场景...
}
```

#### 评审检查项
- [ ] 核心业务逻辑是否有单元测试？
- [ ] 测试覆盖率是否 > 80%？
- [ ] 是否测试了边界情况？
- [ ] 是否使用 mock 隔离外部依赖？

---

## 7. 代码评审检查清单总览

### 7.1 设计原则
- [ ] **S**RP: 单一职责，每个模块只负责一件事
- [ ] **O**CP: 开闭原则，通过扩展而非修改来添加功能
- [ ] **L**SP: 里氏替换，子类可以替换父类
- [ ] **I**SP: 接口隔离，接口小而专注
- [ ] **D**IP: 依赖倒置，依赖抽象而非具体实现

### 7.2 代码质量
- [ ] 命名清晰、语义化
- [ ] 错误处理完善
- [ ] 并发安全
- [ ] 资源正确管理（defer、close）
- [ ] 无内存泄漏

### 7.3 架构设计
- [ ] 遵循 DDD 分层架构
- [ ] 使用充血模型
- [ ] 依赖方向正确
- [ ] 接口定义合理

### 7.4 性能
- [ ] 无 N+1 查询
- [ ] 合理使用缓存
- [ ] 数据库查询优化
- [ ] 批量操作替代循环

### 7.5 安全
- [ ] 输入验证
- [ ] 敏感信息保护
- [ ] SQL 注入防护
- [ ] 密码加密

### 7.6 测试
- [ ] 单元测试覆盖率
- [ ] 集成测试
- [ ] 边界测试

### 7.7 文档
- [ ] 导出类型有注释
- [ ] 复杂逻辑有注释
- [ ] API 文档完整
