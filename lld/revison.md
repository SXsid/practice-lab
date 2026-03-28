# OOP & Design Patterns — Revision Notes

> Not a book. Just the ideas that matter, the traps to avoid, and how they connect.

---

## The One Rule Behind Everything

> **Abstractions should be pulled out by pain, not pushed in by anticipation.**

Don't build the factory until you have two things to put in it.
Don't build the interface until you have two implementations.
Write the simplest code that handles today, structured so tomorrow's change is easy.

---

## Phase 1 — OOP Foundations

### Encapsulation

Bind data + behavior together. Hide data, expose only safe operations.

```python
class BankAccount:
    def __init__(self, balance):
        self.__balance = balance      # hidden

    def deposit(self, amount):        # controlled access
        if amount <= 0:
            raise ValueError("positive only")
        self.__balance += amount
```

**Pitfall** — `_single` underscore is convention only. `__double` actually hides via name mangling. Java enforces with `private`. Go uses lowercase.

**Key question before hiding** — who should control this data? If the answer is "only this class" → hide it.

---

### Composition vs Inheritance

```
HAS A  →  Composition   (Customer has BankAccount)
IS A   →  Inheritance   (SavingsAccount is a BankAccount)
```

**Composition — two types:**

- **Strong (Composition)** — child dies with parent. Customer creates Account internally.
- **Weak (Aggregation)** — child survives parent. Account passed in from outside.

**Inheritance pitfall** — Penguin IS AN Animal but can't fly. If child can't do everything parent does → wrong tool. Use interfaces instead.

**Rule of thumb** — prefer composition over inheritance. Inheritance locks you in. Composition stays flexible.

---

### Interfaces — Contracts

Define WHAT an object can do. Not HOW.

```python
class PaymentProvider(ABC):
    @abstractmethod
    def charge(self, amount): pass

    @abstractmethod
    def refund(self, amount): pass
```

| Language | Enforcement                                     |
| -------- | ----------------------------------------------- |
| Python   | ABC + @abstractmethod → fails on instantiation  |
| Java     | interface + implements → compile time           |
| Go       | implicit → if methods match, contract satisfied |

**Key insight** — depend on interfaces, not concrete types. Your OrderService shouldn't know Razorpay exists.

---

## Phase 2 — Design Thinking

### Single Responsibility

Every class has one reason to change.

```
Customer     →  owns identity data
BankAccount  →  owns balance + money operations
Transaction  →  owns record of what happened
```

**Smell** — if you describe a class using "and" → it has too many responsibilities.

---

### Dependency Injection

Don't reach for dependencies. Receive them.

```python
# bad — reaches for its own dependency
class OrderService:
    def __init__(self):
        self.db = DatabaseConnection()    # hidden, untestable

# good — receives dependency
class OrderService:
    def __init__(self, db: Database):     # explicit, swappable
        self.db = db
```

**DI vs Singleton philosophy:**

```
Singleton  →  object finds its own dependency  (hidden)
DI         →  dependency handed to the object  (explicit)
```

---

### Composition Root

One place where everything is created and wired. Your `main.go` or `main.py`.

```go
func main() {
    cfg    := config.New()
    db     := postgres.NewPool(cfg.DB)
    logger := logger.New(cfg.LogLevel)

    userService    := user.NewService(db, logger)
    paymentService := payment.NewService(db, logger)

    app := app.New(userService, paymentService)
}
```

**Pitfall** — passing entire `app` struct everywhere. Each handler should only receive what it needs.

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

**Valid use cases** — Logger, Feature flags.
**Avoid for** — Database, Services. Create once and pass around instead.

**Node.js gotcha** — module cache gives accidental singleton. NestJS `@Injectable` gives intentional singleton. Not the same thing.

---

### Factory

One place that decides WHICH object to create. Contains the if/else so nobody else has to.

```go
func NewPayment(region string) PaymentProvider {
    switch region {
    case "IN": return NewRazorpay()
    case "US": return NewStripe()
    }
}
```

**Key insight** — complexity doesn't disappear. Good design **contains** it in one place.

**When implementations are totally different per env** → Interface + Factory. Each env has its own struct implementing the interface. No generic builder.

---

### Builder

When object needs many parameters. Named, explicit, validated before construction.

```go
db := NewDatabase().
    WithHost("localhost").
    WithSSL(true).
    WithPoolSize(100).
    Build()        // validates here, fails fast
```

**Why not constructor** — `NewDatabase("localhost", 5432, true, 100, 30, 3)` — what is 100? what is 30? Wrong order = silent bug.

**Factory vs Builder:**

```
Factory  →  WHICH object/config     (based on condition)
Builder  →  HOW to construct        (step by step, validated)
```

They compose — Factory decides which values, Builder constructs safely.

---

## Phase 4 — Structural Patterns

### Decorator

Same interface. Adds behavior. Neither side changes.

```python
class LoggedPayment(PaymentProvider):
    def __init__(self, provider: PaymentProvider):
        self.__provider = provider

    def charge(self, amount):
        log("before charge")
        self.__provider.charge(amount)
        log("after charge")

    def refund(self, amount):
        self.__provider.refund(amount)    # passthrough
```

Stack them in any order:

```python
razorpay = RazorpayPayment()
retry    = RetryPayment(razorpay, retries=3)
logged   = LoggedPayment(retry)
```

**Real world** — HTTP middleware chains, Python @decorators, Java @Transactional.

---

### Adapter

Different interface. Translates between two worlds. Neither side changes.

```python
class SendGridAdapter(Notifier):
    def __init__(self, client: SendGridClient):
        self.__client = client

    def send(self, to: str, message: str):     # your interface
        self.__client.send_email(              # their interface
            from=DEFAULT_FROM,
            to=[to],                           # string → list
            subject=DEFAULT_SUBJECT,
            body=message,
            is_html=False
        )
```

**Real world** — third party SDK wrappers, legacy code bridges, microservice protocol translation.

---

### Facade

Hides multiple interfaces behind one simple one. Caller doesn't know what's inside.

```python
class OrderFacade:
    def place_order(self, user_id, item_id, amount):
        user    = self.user_service.get(user_id)
        _       = self.inventory.check(item_id)
        payment = self.payment.charge(user, amount)
        order   = self.order_repo.create(user, item_id, payment)
        self.notifier.send(user.email, "confirmed!")
        self.inventory.deduct(item_id)
        return order

# handler just does:
order = order_facade.place_order(user_id, item_id, amount)
```

**Handler job** — parse, validate, return response. Not orchestrate.

---

### Three Wrappers — Never Confuse Again

```
Decorator  →  same interface    + adds behavior      (LoggedPayment)
Adapter    →  different interface + translates       (SendGridAdapter)
Facade     →  hides many interfaces behind one       (OrderFacade)
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

**Use class when** — algorithm needs state or config.
**Use plain function when** — algorithm is simple and stateless.

**Smell** — growing if/else choosing between algorithms → Strategy.

**Always combine with DI** — inject strategy, don't reach for it.

---

### Observer

Publisher fires event. Subscribers react. Publisher doesn't know subscribers exist.

```python
class EventBus:
    def publish(self, event, data):
        for handler in self.listeners[event]:
            handler(data)

# publisher
event_bus.publish("ORDER_PLACED", order)

# subscribers — independent, decoupled
event_bus.subscribe("ORDER_PLACED", email_handler)
event_bus.subscribe("ORDER_PLACED", inventory_handler)
```

**Error handling** — db.save failure = stop everything. Email failure = retry independently. Eventual consistency.

**Type safety** — typed event structs, not generic maps. Schema Registry in Kafka for cross-service events.

**Scale progression:**

```
In-memory EventBus   →  single process
RabbitMQ / SQS       →  across services, async retries
Kafka                →  massive scale, event replay
Temporal             →  complex workflows, state management
```

---

### Command

Wraps a request as an object. Enables audit, undo, queue without caller knowing.

```
Factory          →  WHICH command to create
Command          →  WHAT the action does + how to undo
CommandExecutor  →  HOW to run, audit, queue, undo
```

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
        command = self.history.pop()
        command.undo()
```

**Real world** — job queues (Celery, Sidekiq), Git commits, DB migrations (up/down), admin panels.

---

## How Patterns Compose in Real Systems

```
HTTP Request
    │
    ▼
Handler              parse + validate only (Facade consumer)
    │
    ▼
CommandFactory       WHICH command (Factory)
    │
    ▼
CommandExecutor      HOW to run (audit, queue, undo)
    │
    ▼
Command              WHAT to do (Strategy per action)
    │
    ├── Service      orchestrates business logic (Facade)
    │       │
    │       └── Repository   talks to DB (next)
    │
    └── EventBus     fires events (Observer)
            │
            ├── EmailHandler
            ├── InventoryHandler
            └── AnalyticsHandler
```

---

## The Questions to Ask Before Any Abstraction

```
1. How often will this change?
   rarely     → keep it simple
   frequently → abstract, contain the change

2. How many places touch this?
   one place  → inline it
   many places → extract it

3. What breaks if I'm wrong?
   low risk   → ship it, refactor later
   high risk  → design carefully now

4. Who else reads this?
   just me    → maybe okay
   team       → optimise for readability
```

---

## Common Mistakes

| Mistake                                 | Fix                                             |
| --------------------------------------- | ----------------------------------------------- |
| Passing entire app struct everywhere    | Each handler takes only what it needs           |
| Singleton for everything                | Singleton for logger only, pass everything else |
| Inheritance when behavior differs       | Interface + Composition                         |
| if/else scattered across files          | Factory — contain it in one place               |
| Reaching for dependencies inside class  | Inject from outside — DI                        |
| Handler orchestrating multiple services | Facade handles orchestration                    |
| Generic untyped events                  | Typed event structs                             |
| Building abstractions upfront           | Wait for the pain, then extract                 |
