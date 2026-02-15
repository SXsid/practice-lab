import time
from typing import cast

from celery import Celery, Task

celery_app = Celery(
    "myapp", backend="redis://localhost:6379", broker="redis://localhost:6379"
)
celery_app.conf.update(
    task_serializer="json",  # How tasks are serialized
    accept_content=["json"],  # Accepted content types
    result_serializer="json",  # How results are serialized
    timezone="UTC",  # Timezone for scheduled tasks
    enable_utc=True,  # Use UTC
    task_track_started=True,  # Track when task starts
    task_time_limit=30 * 60,  # Hard time limit (30 min)
    task_soft_time_limit=25 * 60,  # Soft time limit (25 min)
    worker_prefetch_multiplier=1,  # How many tasks worker prefetches
    worker_max_tasks_per_child=1000,  # Restart worker after N tasks
)


@celery_app.task(queue="math")
def add(num1, num2):
    time.sleep(2)
    return num1 + num2


add_task = cast(Task, add)


@celery_app.task(name="task.multiply", queue="math")
def multiply(x, y):
    time.sleep(4)
    return x * y


@celery_app.task(max_retries=3, default_retry_delay=60, queue="email", bind=True)
def send_dummy_email(self, email, subject, body):
    try:
        print(f"sending email to {email}")
        time.sleep(5)
        return f"Email sent to {email}"
    except Exception as e:
        raise self.retry(exc=e, coutdown=60 * (2**self.request.retries))


def main():
    print("here")


if __name__ == "__main__":
    main()
