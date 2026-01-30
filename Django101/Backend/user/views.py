from django.http import HttpRequest, HttpResponse, HttpResponseNotAllowed
from django.shortcuts import render
from django.views.decorators.csrf import csrf_exempt


def get_signup_page(request: HttpRequest):
    return render(request=request, template_name="user/signup.html")


def create_account(request: HttpRequest):
    if request.method != "POST":
        return HttpResponseNotAllowed(["POST"])
    name = request.POST.get("name")
    password = request.POST.get("password")
    return HttpResponse(f"{name} you account is created successully!!".encode("utf-8"))
