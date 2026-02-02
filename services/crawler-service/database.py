import os
from sqlalchemy import create_engine
from sqlalchemy.ext.declarative import declarative_base
from sqlalchemy.orm import sessionmaker

# Use environment variable or default for local docker
DATABASE_URL = os.getenv("DATABASE_URL", "postgresql://user:password@postgres:5432/question_db")

# For local testing outside docker (if localhost)
if "localhost" in os.getenv("DB_HOST", ""):
    DATABASE_URL = "postgresql://user:password@localhost:5432/question_db"

engine = create_engine(DATABASE_URL)
SessionLocal = sessionmaker(autocommit=False, autoflush=False, bind=engine)

Base = declarative_base()

def get_db():
    db = SessionLocal()
    try:
        yield db
    finally:
        db.close()
