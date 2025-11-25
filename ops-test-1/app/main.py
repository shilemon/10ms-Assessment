from flask import Flask, jsonify
import time

app = Flask(__name__)

@app.route('/')
def home():
    return jsonify({"message": "Hello from SRE Assessment!", "status": "healthy"})

@app.route('/healthz')
def health():
    return jsonify({"status": "healthy"}), 200

if __name__ == '__main__':
    app.run(host='0.0.0.0', port=8080)
