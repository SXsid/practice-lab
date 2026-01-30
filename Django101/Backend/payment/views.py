from django.http import HttpRequest, HttpResponse, HttpResponseNotAllowed
from django.shortcuts import render
from django.views.decorators.csrf import csrf_exempt


def payment_home(request: HttpRequest):
    return render(request=request, template_name="payment/home.html")


def get_payment_by_id(request: HttpRequest):
    id = request.GET.get("id")
    data = f"This is the response for {id}"
    return HttpResponse(data.encode("utf-8"))


@csrf_exempt
def pay(request: HttpRequest):
    if request.method == "GET":
        return render(request=request, template_name="payment/home.html")
    elif request.method == "POST":
        name = request.POST.get("name")
        amount = request.POST.get("amount")
        return HttpResponse(
            f"{name} your payment of {amount}INR is successfull".encode("utf-8")
        )
    else:
        return HttpResponseNotAllowed(["POST,GET"])
