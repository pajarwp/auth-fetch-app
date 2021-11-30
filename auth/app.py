from sqlite3.dbapi2 import connect
from flask import Flask
from database import init_db

app = Flask(__name__)
init_db()

@app.route("/")
def hello():
    return "Hello World"