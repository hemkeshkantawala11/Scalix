from locust import HttpUser, task, between
import random
import string

def random_ascii_string(length=32):
    return ''.join(random.choices(string.ascii_letters + string.digits, k=length))

class CacheServerLoadTest(HttpUser):
    wait_time = between(1, 2)  # Simulate real-world wait times
    keys = []  # Store successfully inserted keys

    @task(2)  # Higher weight for POST requests
    def set_value(self):
        """ Inserts a random key-value pair into the cache """
        key = random_ascii_string(10)
        value = random_ascii_string(50)

        response = self.client.post("/api/v1/cache", json={"Key": key, "Value": value})
        if response.status_code == 200:
            self.keys.append(key)  # Store successfully inserted keys

    @task(1)
    def get_value(self):
        """ Fetches a value for a previously inserted key """
        if self.keys:
            key = random.choice(self.keys)
            self.client.get(f"/api/v1/cache?key={key}")

