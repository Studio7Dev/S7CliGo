import os
from flask import Flask, request, jsonify, make_response
import json

app = Flask(__name__)
DATABASE_FILE = 'users.json'
if not os.path.exists(DATABASE_FILE):
    with open(DATABASE_FILE, 'w') as f:
        json.dump([], f)

class UserModel:
    @staticmethod
    def load_users():
        """Load users from JSON."""
        with open(DATABASE_FILE, 'r') as f:
            return json.load(f)

    @staticmethod
    def save_users(users):
        """Save users into JSON."""
        with open(DATABASE_FILE, 'w') as f:
            json.dump(users, f)

    @staticmethod
    def get(username):
        """Get a single user by username."""
        users = UserModel.load_users()
        for user in users:
            if user['username'] == username:
                return user
        return None

    @staticmethod
    def create(username, password, role):
        """Create a new user."""
        users = UserModel.load_users()

        if any(user['username'] == username for user in users):
            response = {
                "statusCode": 409,
                "message": "Username already taken.",
            }
            return make_response(jsonify(response)), 409

        new_user = {"username": username, "password": password, "role": role}
        users.append(new_user)
        UserModel.save_users(users)

        response = {
            "statusCode": 201,
            "message": "Successfully created account.",
        }
        return make_response(jsonify(response)), 201

@app.route("/login", methods=["POST"])
def login():
    data = request.get_json()
    user = UserModel.get(data['username'])
    if user and user['password'] == data['password']:
        auth_response = {
            "statusCode": 200,
            "authenticated": True,
            "username": user['username'],
            "role": user['role'],
        }
    else:
        auth_response = {
            "statusCode": 401,
            "authenticated": False,
            "errorMessage": "Incorrect username or password.",
        }
    resp = make_response(jsonify(auth_response))
    resp.headers['Content-Type'] = "application/json"
    return resp

@app.route("/register", methods=["POST"])
def register():
    data = request.get_json()
    username = data['username']
    password = data['password']
    role = "user"

    return UserModel.create(username, password, role)

if __name__ == "__main__":
    app.run(debug=False, )