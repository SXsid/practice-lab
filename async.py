import asyncio


async def TestFunction():
    print("first line")
    await asyncio.sleep(2)
    print("second line")


def som():
    print(20 + 3)


asyncio.run(TestFunction())
som()
