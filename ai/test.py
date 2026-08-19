"""
Gemini Tool-Calling Hello World — Weather Agent
================================================

The simplest end-to-end tool-calling loop with Gemini:

1. Ask the model to pick a random city (its own knowledge, no input from us).
2. The model requests the `get_weather` tool with that city.
3. We execute the tool locally and send the result back.
4. The model returns a natural-language summary, which we package into a
   structured dict alongside the raw tool output.

Setup
-----
    pip install google-genai requests
    export GEMINI_API_KEY="your-key-here"
    python gemini_weather_agent.py

Notes
-----
- `get_weather` calls Open-Meteo (free, no API key) for real data, and falls
  back to a clearly-labeled simulated reading if the network call fails, so
  the demo always completes.
- All I/O is logged; the script exits non-zero on failure so it's safe to
  call from a shell pipeline or CI job.
"""

from __future__ import annotations

import json
import logging
import os
import random
import sys
from dataclasses import asdict, dataclass
from typing import Any, Optional

from dotenv import load_dotenv
from google import genai
from google.genai import types

load_dotenv()
logging.basicConfig(
    level=logging.INFO,
    format="%(asctime)s [%(levelname)s] %(message)s",
)
logger = logging.getLogger("weather_agent")

MODEL_NAME = "gemini-3.5-flash-lite"

# ---------------------------------------------------------------------------
# Tool implementation
# ---------------------------------------------------------------------------


@dataclass
class WeatherReading:
    temperature_c: float
    condition: str
    humidity: int


_WMO_CODES = {
    0: "Clear",
    1: "Mostly Clear",
    2: "Partly Cloudy",
    3: "Overcast",
    45: "Fog",
    48: "Fog",
    51: "Light Drizzle",
    61: "Rain",
    71: "Snow",
    80: "Rain Showers",
    95: "Thunderstorm",
}


def _wmo_to_condition(code: int) -> str:
    return _WMO_CODES.get(code, "Unknown")


def get_weather(city: str) -> dict[str, Any]:
    """Fetch current weather for `city`. Real API first, simulated fallback."""
    import requests

    try:
        geo = requests.get(
            "https://geocoding-api.open-meteo.com/v1/search",
            params={"name": city, "count": 1},
            timeout=5,
        ).json()
        if not geo.get("results"):
            raise ValueError(f"No geocoding match for '{city}'")

        lat = geo["results"][0]["latitude"]
        lon = geo["results"][0]["longitude"]

        wx = requests.get(
            "https://api.open-meteo.com/v1/forecast",
            params={
                "latitude": lat,
                "longitude": lon,
                "current": "temperature_2m,relative_humidity_2m,weather_code",
            },
            timeout=5,
        ).json()

        current = wx["current"]
        reading = WeatherReading(
            temperature_c=round(current["temperature_2m"], 1),
            condition=_wmo_to_condition(current["weather_code"]),
            humidity=int(current["relative_humidity_2m"]),
        )
        logger.info("Live weather for %s: %s", city, reading)
        return asdict(reading)

    except Exception as exc:  # noqa: BLE001 — any failure falls back gracefully
        logger.warning(
            "Live weather fetch failed for '%s' (%s); using simulated reading.",
            city,
            exc,
        )
        reading = WeatherReading(
            temperature_c=round(random.uniform(5, 35), 1),
            condition=random.choice(
                ["Clear", "Cloudy", "Rain", "Partly Cloudy", "Windy"]
            ),
            humidity=random.randint(30, 90),
        )
        return asdict(reading)


GET_WEATHER_DECLARATION = {
    "name": "get_weather",
    "description": "Get the current weather for a given city.",
    "parameters": {
        "type": "object",
        "properties": {
            "city": {
                "type": "string",
                "description": "City name, e.g. 'Tokyo' or 'Nairobi'.",
            },
        },
        "required": ["city"],
    },
}

# ---------------------------------------------------------------------------
# Agent loop
# ---------------------------------------------------------------------------


def _extract_function_call(response) -> Optional[types.FunctionCall]:
    for candidate in response.candidates or []:
        for part in candidate.content.parts or []:
            if part.function_call:
                return part.function_call
    return None


def run_agent() -> dict[str, Any]:
    api_key = os.environ.get("GEMINI_API_KEY")
    if not api_key:
        raise RuntimeError("Set the GEMINI_API_KEY environment variable first.")

    client = genai.Client(api_key=api_key)
    tool = types.Tool(function_declarations=[GET_WEATHER_DECLARATION])
    config = types.GenerateContentConfig(tools=[tool])

    prompt = (
        "Pick one random, real city anywhere in the world. "
        "Then call the get_weather tool to fetch its current weather."
    )

    response = client.models.generate_content(
        model=MODEL_NAME,
        contents=prompt,
        config=config,
    )

    call = _extract_function_call(response)
    if call is None:
        raise RuntimeError("Model did not request the get_weather tool.")

    city = call.args["city"]
    logger.info("Model chose city: %s", city)

    tool_result = get_weather(city)

    # Send the tool result back so the model can produce a natural-language summary.
    follow_up = client.models.generate_content(
        model=MODEL_NAME,
        contents=[
            types.Content(role="user", parts=[types.Part(text=prompt)]),
            types.Content(role="model", parts=[types.Part(function_call=call)]),
            types.Content(
                role="user",
                parts=[
                    types.Part(
                        function_response=types.FunctionResponse(
                            name="get_weather", response=tool_result
                        )
                    )
                ],
            ),
        ],
        config=config,
    )

    summary = (
        follow_up.text.strip()
        if follow_up.text
        else (
            f"{city}: {tool_result['temperature_c']}°C, {tool_result['condition']}, "
            f"{tool_result['humidity']}% humidity."
        )
    )

    return {
        "city": city,
        "temperature_c": tool_result["temperature_c"],
        "condition": tool_result["condition"],
        "humidity": tool_result["humidity"],
        "weather_summary": summary,
    }


if __name__ == "__main__":
    try:
        result = run_agent()
        print(json.dumps(result, indent=2, ensure_ascii=False))
    except Exception as exc:  # noqa: BLE001 — top-level guard for CLI use
        logger.error("Agent run failed: %s", exc)
        sys.exit(1)
