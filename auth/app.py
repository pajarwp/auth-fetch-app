from flask import Flask, make_response, jsonify
from flask_jwt_extended import JWTManager
from database import init_db
import os
from dotenv import load_dotenv
from jsonschema import ValidationError

def create_app():
    app = Flask(__name__)
    app.config['BUNDLE_ERRORS'] = True
    load_dotenv()
    app.config["JWT_SECRET_KEY"] = os.getenv("JWT_SECRET")
    app.config["JWT_TOKEN_LOCATION"] = 'headers'
    app.register_error_handler(400, handle_bad_request)
    JWTManager(app)
    
    init_db()
    
    from api.user.user import bp_user
    app.register_blueprint(bp_user, url_prefix='/user')

    return app

def handle_bad_request(error):
    if isinstance(error.description, ValidationError):
        original_error = error.description
        return make_response(jsonify({'error': original_error.message}), 400)
    return error

if __name__ == "__main__":
    app = create_app()
    app.run(debug=True)