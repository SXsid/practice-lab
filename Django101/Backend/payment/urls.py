from django.urls import path

from . import views

urlpatterns = [
    path(
        "",
        view=views.payment_home,
        name="payment home page",
    ),
    path(
        "/",
        view=views.get_payment_by_id,
    ),
    path(
        "/pay",
        view=views.pay,
        name="pay me bitches",
    ),
]
