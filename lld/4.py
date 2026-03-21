# wraping hte provider in logger claass
from abc import ABC, abstractmethod


class PaymentProvider(ABC):
    @abstractmethod
    def charge(self, amount):
        pass

    @abstractmethod
    def refund(self, amount):
        pass

    @abstractmethod
    def get_status(self):
        pass


class LoggedPayment(PaymentProvider):
    def __init__(self, provider: PaymentProvider):
        self.__provider = provider  # upi or razaorpay
        # what goes here?

    def charge(self, amount):
        # log here
        self.__provider.charge(amount)
        # log here

    def refund(self, amount):
        self.__provider.refund(amount)

    def get_status(self):
        self.__provider.get_status()


class OrderService:
    def __init__(self, paymentProvider: PaymentProvider) -> None:
        self.paymentProvider = paymentProvider

    def checkout(self, amount):
        self.paymentProvider.charge(amount)


# INFO:
# razorpay = RazorpayPayment()  acutal payment impletion of payemprovide rwiht razpoay keys
# logged = LoggedPayment(razorpay) just a genric logger wraper where whic can logg any payemn provider
# order = OrderService(logged)depend on PaymentProvider interface
# order.checkout(1000)
# # logs before charge
# # actually charges via razorpay
# # logs after charge
