from flask import Flask, request, jsonify

app = Flask(__name__)

# 1) Health check
@app.get("/ping")
def ping():
    return jsonify(status="ok"), 200

# 2) Path + query params
#    Example: /greet/Shashank?title=Mr
@app.get("/greet/<name>")
def greet(name):
    title = request.args.get("title", "")
    full = f"{title} {name}".strip()
    return jsonify(message=f"Hello, {full}!"), 202

# 3) Accept JSON body (POST)
#    Send: {"message": "hello"}
@app.post("/echo")
def echo():
    data = request.get_json(silent=True) or {}
    if "message" not in data:
        # bad request example
        return jsonify(error="`message` is required"), 400
    return jsonify(you_sent=data["message"]), 201

if __name__ == "__main__":
    # Run the dev server
    app.run(debug=True)
