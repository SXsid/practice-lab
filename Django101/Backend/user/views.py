from django.http import HttpRequest, HttpResponse, HttpResponseNotAllowed
from django.shortcuts import render

from . import forms


def get_signup_page(request: HttpRequest):
    return render(request=request, template_name="user/signup.html")


def create_account(request: HttpRequest):
    if request.method != "POST":
        return HttpResponseNotAllowed(["POST"])
    form_data = forms.createUserRequest(request.POST)
    if not form_data.is_valid():
        return HttpResponse(str(form_data.errors).encode("utf-8"), status=400)
    name = form_data.cleaned_data["name"]
    age = form_data.cleaned_data["age"]

    return HttpResponse(
        f"{name}@{age} you account is created successully!!".encode("utf-8")
    )
