from flask import Flask
from database import init_db

def create_app():
    app = Flask(__name__)
    app.config['BUNDLE_ERRORS'] = True
    
    init_db()
    
    from api.user.user import bp_user
    app.register_blueprint(bp_user, url_prefix='/user')

    return app

if __name__ == "__main__":
    app = create_app()
    app.run(debug=True)