import logging
from datetime import datetime
from time import perf_counter
from typing import Any, Callable

logging.basicConfig(level=logging.INFO)


def performace(func: Callable[..., Any]) -> Callable[..., Any]:
    def wrapper(*args: Any, **kwargs: Any) -> Any:
        start = perf_counter()
        res = func(*args, **kwargs)
        logging.info(f"Time it took :{perf_counter()-start:.6f}")
        return res

    return wrapper


@performace
def hello_world():
    print("hello world")


hello_world()
