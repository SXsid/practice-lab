from django.http import HttpRequest, HttpResponse, JsonResponse
from django.shortcuts import render

# Create your views here.


def createUser(requset: HttpRequest):
    return JsonResponse({"helo": 2})
