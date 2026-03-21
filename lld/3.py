# aggreation
# when owner ship is not as storng as composion like


# child is storigly realed to parent
class BankAccount:
    def __init__(self, initial_balance) -> None:
        self.__intial_balance = initial_balance


class Customer:
    def __init__(self, name, email):
        self.__name = name
        self.__accounts = []

    def open_account(self, initial_balance):
        account = BankAccount(initial_balance)  # Customer CREATES the account
        self.__accounts.append(account)
        return account


# aggrestion


class AggreBankAccount:
    def __init__(self, initBalance, owners) -> None:
        self.__owners = owners or []
        self.__initBalance = initBalance

    def add_owner(self, customer):
        self.__owners.append(customer)

    def remove_owner(self, customer):
        self.__owners.remove(customer)
        if len(self.__owners) == 0:
            raise ValueError("Account must have at least one owner")


class AggreCustomer:
    def __init__(self, name, email) -> None:
        self.__name = name
        self.email = email
        self.__accounts = (
            []
        )  # the parent clas don't create child hence child can exists multiple placess

    def add_account(self, account: AggreBankAccount):
        self.__accounts.append(account)


## Zoom out — what you've learned so far
#
# Customer
#     │
#     ├── [Composition]  sole BankAccount  →  dies with customer
#     │
#     └── [Aggregation]  shared BankAccount → survives, transfers ownership
#
# `
