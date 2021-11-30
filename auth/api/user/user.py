from flask import Blueprint
from flask_restful import Api, Resource, reqparse
from database import get_connection
from util.response import create_response
import datetime, random, string

bp_user = Blueprint('user', __name__)
api = Api(bp_user)

class User(Resource):
    def __init__(self):
        self.db = get_connection()
        
    def post(self):
        parse = reqparse.RequestParser()
        parse.add_argument('phone_number', location='json', type=str, required=True)
        parse.add_argument('name', location='json', type=str, required=True)
        parse.add_argument('role', location='json', type=str, required=True)
        args = parse.parse_args()
                
        user = self.db.execute("SELECT phone FROM user WHERE phone=" + args["phone_number"]).fetchone()
        if user != None:
            return create_response([], "Phone number already used, use different phone number", "failed")
        
        password = self.create_password()
        created_at = datetime.datetime.now().strftime("%Y-%m-%d %H:%M:%S")
        
        self.db.execute("INSERT INTO user VALUES (?, ?, ?, ?, ?)", [args["phone_number"], args["name"], args["role"], password, created_at])
        self.db.commit()
        self.db.close()
        data = {
            "phone": args["phone_number"],
            "name": args["name"],
            "role": args["role"],
            "password": password,
            "created_at": created_at,
        }
        return create_response(data, "Register Success", "ok")
    
    def create_password(self):
        password = random.choices(string.hexdigits, k=4)
        password = ''.join(password)
        row = self.db.execute("SELECT password FROM user WHERE password='" + password + "'").fetchone()
        if row != None:
            self.create_password()
        return password
        
        
            
api.add_resource(User,'')