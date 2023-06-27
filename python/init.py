import os
import sys
import hashlib

import base64


def hash(data: str) -> str:
    h = hashlib.sha3_512()
    h.update(data.encode("utf-8"))
    hash_value = h.digest()
    return base64.b64encode(hash_value)


def sqlite(login, password):

    os.system(f"""
    rm sqlite/database.db >/dev/null 2>&1
    mkdir sqlite >/dev/null 2>&1
    touch sqlite/database.db >/dev/null 2>&1

    sqlite3 sqlite/database.db "CREATE TABLE IF NOT EXISTS users (id INTEGER PRIMARY KEY, login TEXT, password TEXT)"

    sqlite3 sqlite/database.db "CREATE TABLE IF NOT EXISTS allowed (id INTEGER PRIMARY KEY, ip TEXT)"
    sqlite3 sqlite/database.db "INSERT INTO users(login, password) VALUES('{login}', '{password}')"
    """)

    print("Success")


if __name__ == "__main__":

    try:
        login = sys.argv[1]
    except Exception as e:
        print(
            "python3 python/init.py <LOGIN> <PASSWORD>\nExample:\npython3 python/init.py root root")
        exit()

    if login == "hash":
        try:
            pas = sys.argv[3]
            log = sys.argv[2]
        except:
            print(
                "python3 python/init.py hash <LOGIN> <PASSWORD>\nExample:\npython3 python/init.py hash root root")
            exit()

        print(
            "\n=================================\n\t      Login\n=================================\n")
        print(hash(log).decode("utf-8"))

        print(
            "\n=================================\n\t     Password\n=================================\n")

        print(hash(pas).decode("utf-8"))
        exit()

    try:
        password = sys.argv[2]
    except:
        print("python3 python/init.py <LOGIN> <PASSWORD>\nExample:\npython3 python/init.py root root")
        exit()

    sqlite(hash(login).decode("utf-8"), hash(password).decode("utf-8"))
