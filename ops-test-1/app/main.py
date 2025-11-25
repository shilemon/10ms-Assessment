from flask import Flask, jsonify
import time

app = Flask(__name__)

@app.route('/')
def home():
    # FIX: Remove or reduce the sleep that causes slowness
    # Original might have: time.sleep(5)
    return jsonify({"message": "Hello from SRE Assessment!", "status": "healthy"})

@app.route('/healthz')
def health():
    # FIX: Return proper 200 status code
    # Original might have: return jsonify({"status": "healthy"}), 500
    return jsonify({"status": "healthy"}), 200

if __name__ == '__main__':
    app.run(host='0.0.0.0', port=8080)
