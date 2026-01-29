from django.http import HttpRequest, HttpResponse
from django.shortcuts import render


# Create your views here.
async def hello_world(request: HttpRequest):
    return HttpResponse(b"hello world i am learning django")
