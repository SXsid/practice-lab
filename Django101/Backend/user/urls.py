from django.urls import path

from . import views

urlpatterns = [
    path(
        "signup",
        view=views.get_signup_page,
        name="Render singup page",
    ),
    path(
        "create",
        view=views.create_account,
        name="create_account",
    ),
]
