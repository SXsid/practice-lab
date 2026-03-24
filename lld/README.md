# Desgin Patterns

## Createonal Desgin Pattern

1. Singleton pattern : One instance throughout the process — but prefer passing it around explicitly over letting classes reach for it globally. Use sparingly — logger is the classic valid use case.

2. Factory pattern : one place that decides which object or which configuration to create based on some condition (env, region, user type). Contains the if/else so nobody else has to.
3. builder pattern : when an object needs many parameters, Builder makes each parameter named and explicit, provides sensible defaults, and validates before constructing. Prevents silent bugs from wrong parameter order

```Markdown
Abstraction          ← the big idea
    │
    ├── Encapsulation    ← hide data complexity
    ├── Interfaces       ← hide implementation complexity
    ├── Factory          ← hide creation complexity
    ├── Builder          ← hide construction complexity
    └── Singleton        ← hide lifecycle complexity
```

```Markdown
Creational Patterns — HOW objects are created

Singleton  →  hide lifecycle    (one instance, don't care how)
Factory    →  hide selection    (right object, don't care which)
Builder    →  hide construction (valid object, don't care how built)

Core idea behind all three → Abstraction
```

## structual Desgin pattern

1. Decorator → same interface + adds behavior
   LoggedPayment.Send() calls RazorpayPayment.Send() + logs

2. Adapter → different interface + translates
   SendGridAdapter.Send() calls SendGridClient.SendEmail()
3. Facade → hides multiple interfaces behind one simple one (used in service orchestration)
