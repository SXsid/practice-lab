from django.urls import path

from . import views

urlpatterns = [
    path(
        "",
        view=views.payment_home,
        name="payment home page",
    ),
]
