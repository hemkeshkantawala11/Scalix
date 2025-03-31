from locust import HttpUser, task, between
import random
import string
import json

class CacheLoadTest(HttpUser):
    wait_time = between(1, 2)  # Simulate realistic user wait times
    
    def on_start(self):
        self.keys = []  # Store inserted keys
    
    def random_string(self, length=10):
        return ''.join(random.choices(string.ascii_letters + string.digits, k=length))
    
    @task(2)  # 2x weight for setting cache values
    def set_cache(self):
        key = self.random_string()
        value = self.random_string(20)
        ttl = random.randint(10, 600)
        payload = {"Key": key, "Value": value}
        headers = {"Content-Type": "application/json"}
        
        response = self.client.post("/put", data=json.dumps(payload), headers=headers)
        if response.status_code == 200:
            self.keys.append(key)  # Store successful keys
    
    @task(3)  # 3x weight for getting cache values
    def get_cache(self):
        if not self.keys:
            return  # Avoid requests if no keys exist
        key = random.choice(self.keys)  # Choose a stored key
        self.client.get(f"/get?key={key}")
