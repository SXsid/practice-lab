"""
PRACTICAL CELERY EXAMPLE: E-commerce Order Processing System

This is a working example that demonstrates:
1. Multiple queues (high priority, normal, emails)
2. Task chains and groups
3. Retry logic with exponential backoff
4. Progress tracking
5. Error handling
6. Real-world use case

TO RUN THIS EXAMPLE:
1. Install: pip install celery redis
2. Start Redis: redis-server
3. In terminal 1: celery -A practical_example worker -Q high_priority,orders,emails --loglevel=info
4. In terminal 2: python practical_example.py
"""

import random
import time
from datetime import datetime, timedelta

from celery import Celery, chain, chord, group

# Initialize Celery
app = Celery(
    "ecommerce", broker="redis://localhost:6379/0", backend="redis://localhost:6379/1"
)

# Configuration
app.conf.update(
    task_serializer="json",
    accept_content=["json"],
    result_serializer="json",
    timezone="UTC",
    enable_utc=True,
    task_track_started=True,
)

# Task routing
app.conf.task_routes = {
    "ecommerce.validate_order": {"queue": "high_priority"},
    "ecommerce.charge_payment": {"queue": "high_priority"},
    "ecommerce.process_order": {"queue": "orders"},
    "ecommerce.update_inventory": {"queue": "orders"},
    "ecommerce.send_email": {"queue": "emails"},
}

# ============================================================================
# TASKS
# ============================================================================


@app.task(name="ecommerce.validate_order")
def validate_order(order_id):
    """Validate order details"""
    print(f"Validating order {order_id}...")
    time.sleep(1)  # Simulate validation

    # Simulate validation logic
    order_data = {
        "order_id": order_id,
        "items": [
            {"product_id": 101, "quantity": 2, "price": 29.99},
            {"product_id": 102, "quantity": 1, "price": 49.99},
        ],
        "total": 109.97,
        "customer_id": "CUST-12345",
        "email": "customer@example.com",
    }

    print(f"✓ Order {order_id} validated")
    return order_data


@app.task(name="ecommerce.charge_payment", bind=True, max_retries=3)
def charge_payment(self, order_data):
    """
    Charge payment with retry logic
    Simulates payment gateway that might fail
    """
    order_id = order_data["order_id"]
    amount = order_data["total"]

    print(f"Charging ${amount} for order {order_id}...")

    try:
        # Simulate payment gateway call
        time.sleep(1)

        # Simulate occasional failures (20% chance)
        if random.random() < 0.2:
            raise Exception("Payment gateway timeout")

        payment_result = {
            "payment_id": f"PAY-{order_id}-{int(time.time())}",
            "status": "success",
            "amount": amount,
            "charged_at": datetime.utcnow().isoformat(),
        }

        order_data["payment"] = payment_result
        print(f"✓ Payment charged: {payment_result['payment_id']}")
        return order_data

    except Exception as exc:
        # Retry with exponential backoff
        retry_count = self.request.retries
        countdown = 2**retry_count  # 1s, 2s, 4s

        print(
            f"⚠ Payment failed, retry {retry_count + 1}/{self.max_retries} in {countdown}s"
        )

        if retry_count >= self.max_retries:
            # Max retries reached, send alert
            send_payment_failure_alert.delay(order_id, str(exc))
            raise

        raise self.retry(exc=exc, countdown=countdown)


@app.task(name="ecommerce.update_inventory")
def update_inventory(order_data):
    """Update inventory for ordered items"""
    order_id = order_data["order_id"]
    print(f"Updating inventory for order {order_id}...")

    for item in order_data["items"]:
        product_id = item["product_id"]
        quantity = item["quantity"]
        print(f"  - Product {product_id}: -{quantity} units")
        time.sleep(0.5)

    print(f"✓ Inventory updated")
    return order_data


@app.task(name="ecommerce.create_shipment")
def create_shipment(order_data):
    """Create shipment record"""
    order_id = order_data["order_id"]
    print(f"Creating shipment for order {order_id}...")
    time.sleep(1)

    shipment = {
        "shipment_id": f"SHIP-{order_id}",
        "tracking_number": f"TRK{random.randint(100000, 999999)}",
        "estimated_delivery": (datetime.utcnow() + timedelta(days=5)).isoformat(),
    }

    order_data["shipment"] = shipment
    print(f"✓ Shipment created: {shipment['tracking_number']}")
    return order_data


@app.task(name="ecommerce.send_email", bind=True, max_retries=3)
def send_email(self, recipient, subject, body, order_id=None):
    """Send email with retry logic"""
    print(f"Sending email to {recipient}: {subject}")

    try:
        # Simulate email service
        time.sleep(1)

        # Simulate occasional failures (10% chance)
        if random.random() < 0.1:
            raise Exception("Email service unavailable")

        print(f"✓ Email sent to {recipient}")
        return f"Email sent to {recipient}"

    except Exception as exc:
        retry_count = self.request.retries
        print(f"⚠ Email failed, retry {retry_count + 1}/{self.max_retries}")

        if retry_count >= self.max_retries:
            # Log failure but don't break the workflow
            print(f"✗ Email to {recipient} failed after {self.max_retries} retries")
            return f"Email failed to {recipient}"

        raise self.retry(exc=exc, countdown=30)


@app.task(name="ecommerce.process_order", bind=True)
def process_order(self, order_id):
    """
    Main order processing task with progress tracking
    Coordinates the entire workflow
    """
    print(f"\n{'='*60}")
    print(f"PROCESSING ORDER: {order_id}")
    print(f"{'='*60}\n")

    try:
        # Step 1: Validate
        self.update_state(state="PROGRESS", meta={"step": "validation", "progress": 20})
        order_data = validate_order(order_id)

        # Step 2: Charge payment
        self.update_state(state="PROGRESS", meta={"step": "payment", "progress": 40})
        order_data = charge_payment(order_data)

        # Step 3: Update inventory and create shipment in parallel
        self.update_state(
            state="PROGRESS", meta={"step": "fulfillment", "progress": 60}
        )

        # Run these two tasks in parallel using group
        parallel_tasks = group(
            [update_inventory.s(order_data), create_shipment.s(order_data)]
        )
        results = parallel_tasks.apply_async().get()

        # Merge results (both return order_data with updates)
        order_data = results[0]  # Base data from inventory update
        order_data["shipment"] = results[1]["shipment"]  # Add shipment info

        # Step 4: Send confirmation email
        self.update_state(
            state="PROGRESS", meta={"step": "notification", "progress": 80}
        )

        email_subject = f"Order Confirmation - {order_id}"
        email_body = f"""
        Thank you for your order!
        
        Order ID: {order_id}
        Total: ${order_data['total']}
        Payment ID: {order_data['payment']['payment_id']}
        Tracking Number: {order_data['shipment']['tracking_number']}
        Estimated Delivery: {order_data['shipment']['estimated_delivery']}
        """

        send_email.delay(
            recipient=order_data["email"],
            subject=email_subject,
            body=email_body,
            order_id=order_id,
        )

        # Complete
        self.update_state(state="PROGRESS", meta={"step": "completed", "progress": 100})

        print(f"\n{'='*60}")
        print(f"✓ ORDER {order_id} COMPLETED SUCCESSFULLY")
        print(f"{'='*60}\n")

        return {
            "status": "completed",
            "order_id": order_id,
            "payment_id": order_data["payment"]["payment_id"],
            "tracking_number": order_data["shipment"]["tracking_number"],
        }

    except Exception as e:
        print(f"\n✗ ORDER {order_id} FAILED: {str(e)}\n")

        # Send failure notification
        send_email.delay(
            recipient="support@example.com",
            subject=f"Order Processing Failed - {order_id}",
            body=f"Order {order_id} failed: {str(e)}",
            order_id=order_id,
        )

        raise


@app.task(name="ecommerce.send_payment_failure_alert")
def send_payment_failure_alert(order_id, error):
    """Send alert when payment fails"""
    print(f"🚨 ALERT: Payment failed for order {order_id}: {error}")
    # In real app: send to monitoring system, Slack, email, etc.


@app.task(name="ecommerce.batch_process_orders")
def batch_process_orders(order_ids):
    """
    Process multiple orders in parallel
    Demonstrates group pattern
    """
    print(f"\n{'='*60}")
    print(f"BATCH PROCESSING {len(order_ids)} ORDERS")
    print(f"{'='*60}\n")

    # Process all orders in parallel
    job = group([process_order.s(order_id) for order_id in order_ids])

    result = job.apply_async()
    print(f"Queued {len(order_ids)} orders for processing")

    return f"Batch of {len(order_ids)} orders queued"


# ============================================================================
# WORKFLOW EXAMPLES
# ============================================================================


def example_simple_order():
    """Example 1: Process a single order"""
    print("\n=== EXAMPLE 1: Simple Order Processing ===\n")

    result = process_order.delay("ORD-001")
    print(f"Task ID: {result.id}")
    print("Waiting for completion...")

    # Poll for progress
    while not result.ready():
        if result.state == "PROGRESS":
            info = result.info
            print(
                f"Progress: {info.get('progress', 0)}% - {info.get('step', 'processing')}"
            )
        time.sleep(2)

    if result.successful():
        print(f"\nResult: {result.get()}")
    else:
        print(f"\nFailed: {result.info}")


def example_batch_orders():
    """Example 2: Process multiple orders in parallel"""
    print("\n=== EXAMPLE 2: Batch Order Processing ===\n")

    order_ids = [f"ORD-{i:03d}" for i in range(100, 105)]
    result = batch_process_orders.delay(order_ids)

    print(f"Batch task ID: {result.id}")
    print(result.get(timeout=60))


def example_workflow_chain():
    """Example 3: Chain of tasks"""
    print("\n=== EXAMPLE 3: Task Chain ===\n")

    # Sequential workflow: validate → charge → fulfill
    workflow = chain(
        validate_order.s("ORD-200"),
        charge_payment.s(),
        update_inventory.s(),
        create_shipment.s(),
    )

    result = workflow.apply_async()
    print(f"Chain task ID: {result.id}")
    final_result = result.get(timeout=30)
    print(f"Final result: {final_result}")


def example_scheduled_order():
    """Example 4: Schedule order for future processing"""
    print("\n=== EXAMPLE 4: Scheduled Order ===\n")

    # Process order in 10 seconds
    eta = datetime.utcnow() + timedelta(seconds=10)
    result = process_order.apply_async(args=["ORD-300"], eta=eta)

    print(f"Order scheduled for {eta}")
    print(f"Task ID: {result.id}")


def check_task_status(task_id):
    """Utility: Check status of any task"""
    from celery.result import AsyncResult

    result = AsyncResult(task_id, app=app)

    print(f"\nTask ID: {task_id}")
    print(f"Status: {result.status}")
    print(f"Ready: {result.ready()}")

    if result.state == "PROGRESS":
        print(f"Info: {result.info}")
    elif result.successful():
        print(f"Result: {result.get()}")
    elif result.failed():
        print(f"Error: {result.info}")


# ============================================================================
# MAIN - Run Examples
# ============================================================================

if __name__ == "__main__":
    print("""
    CELERY PRACTICAL EXAMPLE
    ========================
    
    Make sure you have:
    1. Redis running: redis-server
    2. Worker running: celery -A practical_example worker -Q high_priority,orders,emails --loglevel=info
    
    Choose an example to run:
    """)

    print("1. Process single order (with progress tracking)")
    print("2. Process batch of orders (parallel)")
    print("3. Task chain workflow")
    print("4. Schedule order for future")
    print("5. Just queue an order (async)")

    choice = input("\nEnter choice (1-5): ").strip()

    if choice == "1":
        example_simple_order()
    elif choice == "2":
        example_batch_orders()
    elif choice == "3":
        example_workflow_chain()
    elif choice == "4":
        example_scheduled_order()
    elif choice == "5":
        result = process_order.delay("ORD-ASYNC-001")
        print(f"\nOrder queued! Task ID: {result.id}")
        print(f"Check status: result.ready() → {result.ready()}")
        print(f"Get result: result.get() when ready")
    else:
        print("Invalid choice")

    print("\n✓ Done! Check the worker terminal to see task execution.")
