import os
import json
import requests


INGRESS_HOST = os.getenv("INGRESS_HOST", "localhost")
INGRESS_PORT = os.getenv("INGRESS_PORT", "80")
URL = f"http://{INGRESS_HOST}:{INGRESS_PORT}/v1/models/sklearn-iris:predict"


if __name__ == "__main__":
    with open("iris-input.json", "r") as f:
        data = json.load(f)

    headers = {"Host": "sklearn-iris.default.example.com"}

    res = requests.post(URL, json=data, headers=headers)

    print(f"Status code: {res.status_code}")
    print(f"Predictions: {json.loads(res.text)['predictions']}")
