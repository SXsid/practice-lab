# the db way
# but in oops the Cusotem has a relation wiht Bacncomp and it will be less so we can make it composite
from typing import List


class Customer:

    def __init__(self, name, email, accountId) -> None:
        self.__email = email
        self.__name = name
        self.__acount_ids = accountId or []

    def getter_name(self):
        return self.__name

    def getter_email(self):
        return self.__email

    def get_accounts(self):
        return self.__acount_ids


class BankAccount:
    def __init__(self, email, balance):
        self.__email = email
        self.__balance = balance  # exposed

    def withdraw(self, amount):
        if self.__balance < amount:
            raise ValueError("gareeb")
        self.__balance -= amount  # no checks at all

    def deposit(self, amount):
        if amount < 0:
            raise ValueError("kmal krte ho pandey ji")
        self.__balance += amount  # no checks at all

    def get_balance(self):
        return self.__balance


# composite implem tation


class CustomerComposite:

    def __init__(self, name, email, accounts: List[BankAccount]) -> None:
        self.__email = email
        self.__name = name
        self.__accounts: List[BankAccount] = accounts or []

    def getter_name(self):
        return self.__name

    def getter_email(self):
        return self.__email

    def get_account(self):
        return self.__accounts
