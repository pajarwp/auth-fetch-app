from sqlite3.dbapi2 import connect
from flask import Flask
import sqlite3

app = Flask(__name__)
connection = sqlite3.connect("registry.db")
cursor = connection.cursor()
cursor.execute("CREATE TABLE IF NOT EXISTS USER(phone TEXT, name TEXT, role TEXT, password TEXT, created_at TIMESTAMP)")

@app.route("/")
def hello():
    return "Hello World"