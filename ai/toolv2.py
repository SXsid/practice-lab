import logging
import os
from typing import Any, Optional

import requests
from dotenv import load_dotenv
from google import genai
from google.genai._gaos.types.interactions import Interaction
from pydantic import BaseModel

load_dotenv()
API_KEY = os.getenv("OPEN_WHEATHER_KEY")

# logging.basicConfig(
#     level=logging.INFO, format="%s(asctime)s [%(levelname)] %(message)s"
# )


class City_Weather(BaseModel):
    city: str
    condition: str
    temperatre_c: float
    humidity: int


class Response(BaseModel):

    city: str
    temperature_c: float
    condition: str
    humidity: int
    weather_summary: str


def get_wheather(CITY_NAME: str) -> City_Weather:
    weather_url = (
        f"https://api.openweathermap.org/data/2.5/weather"
        f"?q={CITY_NAME}&units=metric&appid={API_KEY}"
    )

    weather_data = requests.get(weather_url).json()

    temp = weather_data["main"]["temp"]
    humid = weather_data["main"]["humidity"]
    condi = weather_data["weather"][0]["description"]
    return City_Weather(
        city=CITY_NAME, temperatre_c=temp, humidity=humid, condition=condi
    )


weather_function = {
    "type": "function",
    "name": "get_current_temperature",
    "description": "Gets the current temperature for a given location.",
    "parameters": {
        "type": "object",
        "properties": {
            "location": {
                "type": "string",
                "description": "The city name, e.g. San Francisco",
            },
        },
        "required": ["location"],
    },
}


def handle_func_call(interaction):
    if interaction.steps:
        for step in interaction.steps:
            if step.type == "function_call":
                if step.name == "get_current_temperature":
                    res = get_wheather(step.arguments["location"])
                    return res


Model = "gemini-3.5-flash-lite"
client = genai.Client()


def llm(input: str, response_format: Optional[BaseModel] = None):
    kwargs = {"model": Model, "input": input, "tools": [weather_function]}
    if response_format:
        kwargs["response_format"] = {
            "type": "text",
            "mime_type": "application/json",
            "schema": response_format.model_json_schema(),
        }
    interaction = (client.interactions.create(**kwargs),)[0]
    return interaction


def main():
    prompt = "get a  random city be as random as you want  then use tools to get the weather and finally based of the rsponse you job is to genreate schema appropitate reposne ,don't assume be directy , always answer corectly no emotion "
    interaction = llm(prompt)
    res = handle_func_call(interaction)
    prompt = f"the outpur from the funton call {res} now you have to sumarize the data and give me some intersting thing of the summrize seciton and follow the sechma we have res no need for funtion call right?"
    interaction = llm(prompt, Response)
    print(Response.model_validate_json(interaction.output_text))


if __name__ == "__main__":
    main()
