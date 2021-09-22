from getpass import getpass
import sqlite3
from Crypto.Protocol.KDF import bcrypt


db_path = input("Enter the db path: ")
username = input("Enter the user's name: ")
password = getpass("Enter the user's password: ")

con = sqlite3.connect(db_path)
cur = con.cursor()

hashed_password = bcrypt(password, 14).decode()

cur.execute("""
INSERT INTO users(name, hashed_password) values(?, ?)
""", (username, hashed_password))

con.commit()
con.close()

