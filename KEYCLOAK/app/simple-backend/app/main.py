from fastapi import FastAPI
from dotenv import load_dotenv
import os

load_dotenv()

app = FastAPI()

# READ ENV VALUE
APP_NAME = os.getenv("APP_NAME", "Simple Backend")

@app.get("/")
def root():
    return {"message": f"Welcome to {APP_NAME}!"}

@app.get("/hello/{name}")
def say_hello(name: str):
    return {"message": f"Hello, {name}!"}

@app.get("/add")
def add(a: int, b: int):
    return {"result": a + b}
