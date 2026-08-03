import os

from openai import OpenAI


def main():
    MODELS = "gemini-2.5-flash-lite"
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
    print(
        client.chat.completions.create(
            model=MODELS,
            messages=[
                {"role": "system", "content": SYSTEM_PROMPT},
                {"role": "user", "content": "go find me random itenery."},
            ],
        )
        .choices[0]
        .message.content
    )


if __name__ == "__main__":
    from dotenv import load_dotenv

    load_dotenv()

    main()
