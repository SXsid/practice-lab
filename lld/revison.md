# OOP & Design Patterns — Complete Revision

> Not a book. The ideas that matter, the traps to avoid, how they connect.
> Read this before solving problems. Confusion clears through building, not reading.

---

## The One Rule Behind Everything

> **Abstractions should be pulled out by pain, not pushed in by anticipation.**

Ask before every abstraction:

1. How often will this change?
2. How many places touch this?
3. What breaks if I'm wrong?
4. Who else reads this?

---

## The Smell Test

```
Scattered if/else growing?         → Factory or Strategy
One class doing too many things?   → SRP + split it
Hard to test?                      → hidden dependency, use DI
Adding feature breaks other code?  → wrong abstraction, use interfaces
10 param constructor?              → Builder
Copied code in 3 places?           → extract to base or interface
Object behaves differently by status? → State
Treating tree structures?          → Composite
Two incompatible interfaces?       → Adapter
Too many things caller must know?  → Facade
Adding behavior without changing?  → Decorator
One event, many reactions?         → Observer
Action needs audit/undo/queue?     → Command
```

---

## Phase 1 — OOP Foundations

### Object

Data + behavior bundled together. Data hidden, behavior exposed.

### Encapsulation

```python
class BankAccount:
    def __init__(self, balance):
        self.__balance = balance      # hidden

    def deposit(self, amount):
        if amount <= 0:
            raise ValueError("positive only")
        self.__balance += amount

    def get_balance(self):
        return self.__balance
```

| Language | How                                   |
| -------- | ------------------------------------- |
| Python   | `__double` underscore — name mangling |
| Java     | `private` keyword — compiler enforced |
| Go       | lowercase field — package level       |

---

### Composition vs Inheritance

```
HAS A  →  Composition   Customer HAS A BankAccount
IS A   →  Inheritance   SavingsAccount IS A BankAccount
```

**Composition types:**

- Strong (Composition) — child dies with parent
- Weak (Aggregation) — child survives parent

**Inheritance pitfall** — Penguin IS AN Animal but can't fly.
If child can't do everything parent promises → wrong tool.

**Rule** — prefer composition. Inheritance locks you in.

---

### Interfaces — Contracts

WHAT an object can do. Not HOW.

```go
type PaymentProvider interface {
    Charge(amount float64) error
    Refund(amount float64) error
}
```

| Language | Enforcement                                    |
| -------- | ---------------------------------------------- |
| Python   | ABC + @abstractmethod → fails on instantiation |
| Java     | interface + implements → compile time          |
| Go       | implicit → methods match = contract satisfied  |

**Key** — depend on interfaces not concrete types.

---

## Phase 2 — Design Thinking

### Single Responsibility

One class, one reason to change.

```
Customer     →  identity data
BankAccount  →  balance + money operations
Transaction  →  record of what happened
```

**Smell** — describing class using "and" → too many responsibilities.

---

### Dependency Injection

Don't reach for dependencies. Receive them.

```go
// bad — reaches internally
type OrderService struct {
    db *PostgresDB    // concrete, unreplaceable
}

// good — receives interface
type OrderService struct {
    db Database    // interface, swappable, testable
}
```

**DI vs Singleton:**

```
Singleton  →  object finds its own dependency  (hidden)
DI         →  dependency handed to the object  (explicit)
```

---

### Composition Root

One place. Everything created and wired. Your `main.go`.

```go
func main() {
    cfg    := config.New()
    pool   := postgres.NewPool(cfg.DB)
    logger := logger.New(cfg.LogLevel)

    userRepo  := postgres.NewUserRepo(pool)
    orderRepo := postgres.NewOrderRepo(pool)

    userService  := service.NewUserService(userRepo)
    orderService := service.NewOrderService(orderRepo, userService)

    handler := api.New(orderService, userService)
}
```

**Pitfall** — passing entire app struct everywhere. Each component takes only what it needs.

---

## Phase 3 — Creational Patterns

### Singleton

One instance throughout the process.

```python
class Config:
    __instance = None
    def __new__(cls):
        if cls.__instance is None:
            cls.__instance = super().__new__(cls)
        return cls.__instance
```

**Valid uses** — Logger, Feature flags.
**Avoid for** — Database, Services. Create once, pass around instead.

---

### Factory

One place decides WHICH object to create.

```go
func NewPayment(region string) PaymentProvider {
    switch region {
    case "IN": return NewRazorpay()
    case "US": return NewStripe()
    case "EU": return NewPaypal()
    }
}
```

**Key insight** — complexity doesn't disappear. Good design contains it in one place.

**When implementations totally differ per env:**
→ Interface + Factory. Each env has own struct. No generic builder.

---

### Builder

Many params. Named. Explicit. Validated before construction.

```go
db := NewDatabase().
    WithHost("localhost").
    WithSSL(true).
    WithPoolSize(100).
    Build()        // validates here
```

**Why not constructor** — positional args = silent wrong-order bugs.

**Factory vs Builder:**

```
Factory  →  WHICH object      (which condition)
Builder  →  HOW to construct  (step by step, validated)
```

They compose — Factory picks values, Builder constructs safely.

---

## Phase 4 — Structural Patterns

### Decorator

Same interface. Adds behavior. Neither side changes.

```python
class LoggedPayment(PaymentProvider):
    def __init__(self, provider: PaymentProvider):
        self.__provider = provider

    def charge(self, amount):
        log("before")
        self.__provider.charge(amount)
        log("after")

    def refund(self, amount):
        self.__provider.refund(amount)    # passthrough
```

Stack in any order:

```python
razorpay = RazorpayPayment()
retried  = RetryPayment(razorpay, retries=3)
logged   = LoggedPayment(retried)
```

**Real world** — HTTP middleware, @decorators, @Transactional.

---

### Adapter

Different interface. Translates. Neither side changes.

```python
class SendGridAdapter(Notifier):
    def __init__(self, client: SendGridClient):
        self.__client = client

    def send(self, to: str, message: str):      # your interface
        self.__client.send_email(               # their interface
            from=DEFAULT_FROM,
            to=[to],                            # str → list
            subject=DEFAULT_SUBJECT,
            body=message,
            is_html=False
        )
```

**Real world** — third party SDK wrappers, legacy bridges.

---

### Facade

Hides many interfaces behind one simple one.

```python
class OrderFacade:
    def place_order(self, user_id, item_id, amount):
        user    = self.user_service.get(user_id)
        _       = self.inventory.check(item_id)
        payment = self.payment.charge(user, amount)
        order   = self.order_repo.create(user, item_id, payment)
        self.notifier.send(user.email, "confirmed!")
        return order
```

Handler just calls one thing. Doesn't know what's inside.

---

### Composite

Individual objects and groups treated uniformly. Tree structures.

```python
class FileSystemItem(ABC):
    @abstractmethod
    def get_size(self) -> int: pass

class File(FileSystemItem):
    def get_size(self) -> int:
        return self.__size

class Folder(FileSystemItem):
    def get_size(self) -> int:
        return sum(child.get_size() for child in self.__children)
```

No `isinstance` checks. Infinite nesting works automatically.

**Real world** — file systems, UI components, org charts, HTML DOM.

---

### Three Wrappers — Never Confuse

```
Decorator  →  same interface    + adds behavior      LoggedPayment
Adapter    →  different interface + translates        SendGridAdapter
Facade     →  hides many interfaces behind one        OrderFacade
```

---

## Phase 5 — Behavioral Patterns

### Strategy

Swappable algorithm behind an interface.

```python
class SortStrategy(ABC):
    @abstractmethod
    def sort(self, orders): pass

class SortByDate(SortStrategy):
    def sort(self, orders):
        return sorted(orders, key=lambda o: o.date)
```

**Use class** — algorithm needs state or config.
**Use function** — algorithm is simple and stateless.

**Always combine with DI** — inject strategy, don't reach for it.

**Smell** — growing if/else choosing between algorithms → Strategy.

---

### Observer

Publisher fires event. Subscribers react independently. Publisher doesn't know subscribers exist.

```python
class EventBus:
    def publish(self, event, data):
        for handler in self.listeners[event]:
            handler(data)

event_bus.publish("ORDER_PLACED", order)
event_bus.subscribe("ORDER_PLACED", email_handler)
event_bus.subscribe("ORDER_PLACED", inventory_handler)
```

**Error handling:**

```
db.save fails    →  stop everything, non negotiable
email fails      →  retry independently
inventory fails  →  retry independently
```

**Type safety** — typed event structs not generic maps.

**Scale progression:**

```
In-memory EventBus  →  single process
RabbitMQ / SQS      →  across services, async retries
Kafka               →  massive scale, event replay
Temporal            →  complex workflows, state
```

---

### Command

Wraps request as object. Enables audit, undo, queue.

```python
class SuspendUserCommand(Command):
    def execute(self): self.user_service.suspend(self.user_id)
    def undo(self):    self.user_service.unsuspend(self.user_id)
    def name(self):    return f"suspend_user:{self.user_id}"

class CommandExecutor:
    def execute(self, command: Command):
        command.execute()
        self.history.append(command)
        audit_log.save(command.name())

    def undo_last(self):
        self.history.pop().undo()
```

**Responsibilities:**

```
Factory          →  WHICH command
Command          →  WHAT it does + how to undo
CommandExecutor  →  HOW to run, audit, queue, undo
```

**Real world** — job queues, Git commits, DB migrations, admin panels.

---

### State

Object changes behavior when internal state changes. Object manages its own transitions.

```python
class PendingState(OrderState):
    def confirm(self, order):
        order.state = ConfirmedState()    # transitions itself

    def cancel(self, order):
        order.state = CancelledState()

    def ship(self, order):
        raise Exception("can't ship pending order")

class Order:
    def __init__(self):
        self.state = PendingState()

    def confirm(self): self.state.confirm(self)
    def cancel(self):  self.state.cancel(self)
    def ship(self):    self.state.ship(self)
```

**State vs Strategy:**

```
Strategy  →  behavior swapped by CALLER
State     →  behavior swapped by OBJECT ITSELF
```

**Real world** — order status, payment lifecycle, TCP connection, traffic lights.

---

## Phase 6 — Architecture Patterns

### Repository

Hides DB behind interface. Service knows nothing about SQL.

```go
// defined by consumer — service owns the interface
type OrderRepository interface {
    GetByID(ctx context.Context, id string) (*Order, error)
    ListByUser(ctx context.Context, userID string) ([]*Order, error)
    Create(ctx context.Context, order *Order) error
    UpdateStatus(ctx context.Context, id string, status string) error
}

// implementation — just satisfies interface
type PostgresOrderRepository struct {
    pool *pgxpool.Pool
}
```

**Interface lives near consumer not implementation** — Go philosophy.

**Interface segregation** — UserService only gets the OrderRepo methods it needs. Not the full interface.

**Testing:**

```go
type MockOrderRepository struct{}
func (m *MockOrderRepository) Create(...) error { return nil }
// no postgres, no docker, instant tests
```

---

### Service Layer

Business logic home. Between handler and repository.

```
Handler     →  parse, validate input, return response
Service     →  business rules, orchestrate repos + services
Repository  →  DB read/write, no business logic
Domain      →  pure structs, no logic, no DB, no HTTP
```

```go
// ✅ correct
func (s *OrderService) PlaceOrder(ctx context.Context, userID string, amount float64) (*Order, error) {
    if amount <= 0 {
        return nil, errors.New("amount must be positive")    // business rule
    }
    user, err := s.userRepo.GetByID(ctx, userID)            // delegates to repo
    // ...
}

// ❌ wrong — business logic in handler
// ❌ wrong — SQL in service
// ❌ wrong — HTTP status codes in service
```

**Who owns what:**

```
Single entity operations      →  Service
Cross entity reads             →  Service (GetUserOrders lives in OrderService)
Cross entity combined data     →  Facade (UserDashboard = user + orders)
Cross service mutations        →  Facade or Saga
```

---

## Full Architecture — Everything Connected

```
HTTP Request
    │
    ▼
Handler                    parse + validate only
    │
    ▼
Facade (optional)          orchestrate across services
    │
    ▼
Service                    business logic
    │
    ├── Repository          DB interface
    │       └── PostgresRepo  actual SQL
    │
    ├── PaymentProvider     external interface
    │       └── StripeAdapter
    │
    └── EventBus            fire and forget
            ├── EmailHandler
            ├── InventoryHandler
            └── AnalyticsHandler
```

Every arrow = interface. Nothing concrete except at edges.

---

## Patterns Working Together — Real Example

```
POST /admin/actions {"type": "suspend_user"}
    │
    ▼
Handler                     parse request
    │
    ▼
CommandFactory              WHICH command  (Factory)
    │
    ▼
CommandExecutor             HOW to run     (audit, undo, queue)
    │
    ▼
SuspendUserCommand          WHAT to do     (Command)
    │
    ▼
UserService                 business rules (Service Layer)
    │
    ├── UserRepository      DB ops         (Repository)
    └── EventBus            notify others  (Observer)
            └── AuditHandler              (Strategy per handler)
```

---

## Common Mistakes

| Mistake                                | Fix                                     |
| -------------------------------------- | --------------------------------------- |
| Passing entire app struct              | Each component takes only what it needs |
| Singleton for everything               | Singleton for logger only               |
| Inheritance when behavior differs      | Interface + Composition                 |
| if/else scattered across files         | Factory — one place                     |
| Reaching for dependencies inside class | DI — receive them                       |
| Handler orchestrating services         | Facade handles orchestration            |
| Generic untyped events                 | Typed event structs                     |
| isinstance checks everywhere           | Composite — shared interface            |
| Status checked in every method         | State pattern                           |
| Building abstractions upfront          | Wait for pain, then extract             |
| Interface near implementation          | Interface near consumer (Go)            |
| Full repo interface everywhere         | Interface segregation                   |

---

## Language Comparison — Same Idea Different Clothes

| Concept          | Python                   | Java              | Go               |
| ---------------- | ------------------------ | ----------------- | ---------------- |
| Hide data        | `__field`                | `private`         | lowercase        |
| Interface        | ABC + abstractmethod     | interface         | implicit         |
| Enforce contract | runtime                  | compile time      | compile time     |
| Composition      | object in `self`         | object in field   | struct embedding |
| DI container     | manual / FastAPI Depends | Spring @Autowired | manual in main   |
| Singleton        | module cache / **new**   | Spring @Bean      | sync.Once        |

---

## What Comes From Experience, Not Patterns

```
These you learn by shipping and breaking things:
- When to use eventual consistency
- Race conditions and concurrency
- Cache invalidation strategies
- Database performance
- Service mesh complexity
- When NOT to use microservices
```

Patterns give you vocabulary and structure.
Experience gives you judgment.
Both matter. Neither replaces the other.

## BLOG

> https://medium.com/@milon.istiyak/top-10-most-used-design-patterns-in-low-level-design-lld-a15e8bb449c6
