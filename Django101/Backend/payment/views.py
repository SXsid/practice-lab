from django.http import HttpRequest, HttpResponse
from django.shortcuts import render


def payment_home(request: HttpRequest):
    return HttpResponse(b"this is payment page")


# Create your views here.
