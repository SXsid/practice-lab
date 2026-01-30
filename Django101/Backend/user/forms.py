from django import forms


class createUserRequest(forms.Form):
    name = forms.CharField(max_length=100)
    age = forms.IntegerField(max_value=150, required=False, initial=10)
    password = forms.CharField(max_length=15, min_length=8)
