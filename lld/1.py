# Encapsulation
# problem .....


# solution
class BankAccount:
    def __init__(self, owner, balance):
        self.owner = owner
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


b = BankAccount("sid", 20)
# print(b.__balance)
print(b.get_balance())
