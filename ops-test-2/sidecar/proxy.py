# sidecar/proxy.py
from flask import Flask, request, Response
import requests
import os

app = Flask(__name__)
UPSTREAM = os.environ.get("UPSTREAM", "http://localhost:8080")

@app.route('/', defaults={'path': ''}, methods=["GET","POST","PUT","DELETE"])
@app.route('/<path:path>', methods=["GET","POST","PUT","DELETE"])
def proxy(path):
    # forward request to upstream deterministically
    try:
        resp = requests.request(
            method=request.method,
            url=f"{UPSTREAM}/{path}",
            headers={k:v for k,v in request.headers if k.lower() != 'host'},
            data=request.get_data(),
            timeout=2.0  # bounded timeout
        )
        return Response(resp.content, status=resp.status_code, headers=dict(resp.headers))
    except requests.exceptions.RequestException as e:
        app.logger.error("proxy error: %s", e)
        return Response("upstream error", status=502)

if __name__ == "__main__":
    app.run(host="0.0.0.0", port=8081)
