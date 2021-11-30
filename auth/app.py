from flask import Flask
from flask_jwt_extended import JWTManager
from database import init_db
import os
from dotenv import load_dotenv

def create_app():
    app = Flask(__name__)
    app.config['BUNDLE_ERRORS'] = True
    load_dotenv()
    app.config["JWT_SECRET_KEY"] = os.getenv("JWT_SECRET")
    JWTManager(app)
    
    init_db()
    
    from api.user.user import bp_user
    app.register_blueprint(bp_user, url_prefix='/user')

    return app

if __name__ == "__main__":
    app = create_app()
    app.run(debug=True)