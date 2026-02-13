from django.urls import path

from . import views

urlspatterns = [path("", view=views.available_job)]
