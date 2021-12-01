from flask import Blueprint
from flask_restful import Api, Resource, reqparse
from flask_jwt_extended import create_access_token, jwt_required, get_jwt_identity
from flask_expects_json import expects_json
from database import get_connection
from util.response import create_response
import datetime, random, string

bp_user = Blueprint('user', __name__)
api = Api(bp_user)

class User(Resource):
    def __init__(self):
        self.db = get_connection()
        
    schema = {
        'type': 'object',
        'properties': {
            'phone': {'type': 'string', 'minLength':8, 'maxLength':12},
            'name': {'type': 'string'},
            'role': {'type': 'string'}
        },
    'required': ['phone', 'name', 'role']
    }
        
    @jwt_required
    def get(self):
        current_user = get_jwt_identity()
        return create_response(current_user, "Success", "ok"), 200
    
    @expects_json(schema)    
    def post(self):
        parse = reqparse.RequestParser()
        parse.add_argument('phone', location='json', type=str, required=True)
        parse.add_argument('name', location='json', type=str, required=True)
        parse.add_argument('role', location='json', type=str, required=True)
        args = parse.parse_args()
        
        check = args["phone"].isnumeric()
        if check == False:
            return create_response([], "Phone must be number", "failed")
                
        user = self.db.execute("SELECT name FROM user WHERE name='" + args["name"] + "'").fetchone()
        if user != None:
            return create_response([], "Name already used, use different name", "failed")
        
        password = self.create_password()
        created_at = datetime.datetime.now().strftime("%Y-%m-%d %H:%M:%S")
        
        self.db.execute("INSERT INTO user VALUES (?, ?, ?, ?, ?)", [args["phone"], args["name"], args["role"], password, created_at])
        self.db.commit()
        self.db.close()
        data = {
            "phone": args["phone"],
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
class UserToken(Resource):
    def __init__(self):
        self.db = get_connection()
        
    schema = {
        'type': 'object',
        'properties': {
            'phone': {'type': 'string', 'minLength':8, 'maxLength':12},
            'password': {'type': 'string', 'minLength':4, 'maxLength':4},
        },
    'required': ['phone', 'password']
    }
    
    @expects_json(schema)  
    def post(self):
        parse = reqparse.RequestParser()
        parse.add_argument('phone', location='json', type=str, required=True)
        parse.add_argument('password', location='json', type=str, required=True)
        args = parse.parse_args()
        
        check = args["phone"].isnumeric()
        if check == False:
            return create_response([], "Phone must be number", "failed")
        
        user = self.db.execute("SELECT name, role, created_at FROM user WHERE phone='" + args["phone"] + "'" + " AND password='" + args["password"] + "'").fetchone()
        if user == None:
            return create_response([], "Invalid credential, please check phone number and password", "Unauthorized"), 401
        claims= {
            "name": user[0],
            "phone": args["phone"],
            "role": user[1],
            "created_at": user[2],
        }
        jwt = create_access_token(claims)
        return create_response({"token": jwt}, "Success", "ok"), 200
    
api.add_resource(User,'')
api.add_resource(UserToken,'/login')
