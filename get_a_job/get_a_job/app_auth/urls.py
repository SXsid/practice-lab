from django.contrib.auth import views as auth_view
from django.urls import path

from . import views

urlpatterns = [
    path("signup", view=views.signup, name="signup"),
    path("logout", view=auth_view.LogoutView.as_view(), name="logout"),
    path("login", view=auth_view.LoginView.as_view(), name="login"),
]
