import os
from typing import List

from openai import OpenAI
from pydantic import BaseModel


class DayPlan(BaseModel):

    day: int
    activities: List[str]


class Itenery(BaseModel):
    source: str
    destination: str
    trip_duration_days: int
    budget_category: str
    top_attractions: List[str]
    daily_plan: List[DayPlan]


def main():
    MODELS = "gemini-3.5-flash"
    gemini_api_key = os.getenv("GEMINI_API_KEY")
    if not gemini_api_key:
        print("key is not set")
        return

    SYSTEM_PROMPT = """pick two city randomly , one as source city and other as destination city and then try to answer or fill this 
    field of the dict and out put in this format by you knowlege,try to mainting the type sanity
    expected output format:
        destination: str
        trip_duration_days: int
        budget_category: str
        top_attractions: list[str]
        daily_plan: list[{day: int, activities: list[str]}].
    each as json key and value  decide by you
    """
    client = OpenAI(
        api_key=gemini_api_key,
        base_url="https://generativelanguage.googleapis.com/v1beta/openai/",
    )
    try:
        res = (
            (
                client.beta.chat.completions.parse(
                    model=MODELS,
                    temperature=2,
                    messages=[
                        {"role": "system", "content": SYSTEM_PROMPT},
                        {"role": "user", "content": "go find me random itenery."},
                    ],
                    response_format=Itenery,
                )
            ).choices[0]
        ).message.parsed
        # manul type informenct / validation -> making no determintc a lil bit predictable
        if not res:
            print("sorry not ")
            return
        print(res.source, "->", res.destination)

    except Exception as e:
        print(e)


if __name__ == "__main__":
    from dotenv import load_dotenv

    load_dotenv()

    main()
