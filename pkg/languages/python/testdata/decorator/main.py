from flask import Flask, make_response

app = Flask(__name__)

@app.route("/api/v2/test_response")
def users():
    headers = {"Content-Type": "application/json"}
    return make_response(
        'Test worked!',
        200,
        headers=headers
    )