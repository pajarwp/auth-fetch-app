import sqlite3

def get_connection():
    return sqlite3.connect("registry.db")

def init_db():  
    connection = get_connection()
    cursor = connection.cursor()
    cursor.execute("CREATE TABLE IF NOT EXISTS user(phone TEXT, name TEXT, role TEXT, password TEXT, created_at TIMESTAMP)")