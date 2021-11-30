import sqlite3

def connect():
    return sqlite3.connect("registry.db")

def init_db():
    connection = connect()
    cursor = connection.cursor()
    cursor.execute("CREATE TABLE IF NOT EXISTS USER(phone TEXT, name TEXT, role TEXT, password TEXT, created_at TIMESTAMP)")