from django.urls import path

from . import views

urlpatterns = [
    path(
        "/hello",
        view=views.hello_world,
        name="hello world",
    ),
]
