import os

import requests
from dotenv import load_dotenv
from openai import OpenAI
from pydantic import BaseModel

load_dotenv()
API_KEY = os.getenv("OPEN_WHEATHER_KEY")
MODELS = "gemini-3.5-flash"
GEMINI_API_KEY = os.getenv("GEMINI_API_KEY")
if not GEMINI_API_KEY:
    print("key is not set")

model = OpenAI(
    api_key=GEMINI_API_KEY,
    base_url="https://generativelanguage.googleapis.com/v1beta/openai/",
)

Prompt1 = (
    "you have to pick one city in the world randomly ,expected output format  name:str"
)
prompt2 = "you have to help me writing this api response more in deatil ok input \n"


class City_Weather(BaseModel):
    condition: str
    temperatre_c: float
    humidity: int


class Response(BaseModel):

    city: str
    temperature_c: float
    condition: str
    humidity: int
    weather_summary: str


class City_name(BaseModel):
    name: str


def get_weather(CITY_NAME: str) -> City_Weather:
    weather_url = (
        f"https://api.openweathermap.org/data/2.5/weather"
        f"?q={CITY_NAME}&units=metric&appid={API_KEY}"
    )

    weather_data = requests.get(weather_url).json()

    temp = weather_data["main"]["temp"]
    humid = weather_data["main"]["humidity"]
    condi = weather_data["weather"][0]["description"]
    return City_Weather(temperatre_c=temp, humidity=humid, condition=condi)


def llm(respo_format, prompt: str):

    res = (
        (
            model.beta.chat.completions.parse(
                model=MODELS,
                temperature=0.3,
                messages=[
                    {"role": "system", "content": prompt},
                    {"role": "user", "content": "go "},
                ],
                response_format=respo_format,
            )
        ).choices[0]
    ).message.parsed
    return res


def main():
    name = llm(City_name, Prompt1).name
    data = get_weather(name)
    print(llm(Response, prompt2 + "city :" + name + f"{data}"))


if __name__ == "__main__":
    main()
