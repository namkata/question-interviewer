import os
import psycopg2

DB_HOST = "localhost"
DB_PORT = "5432"
DB_NAME = "question_db"
DB_USER = "user"
DB_PASSWORD = "password"

# In docker-compose, host is 'postgres', but from host machine (where I am), it is mapped to 5432.
# Wait, the environment says "Operating system: macos". 
# The docker-compose ports are "5432:5432". So localhost:5432 should work.

def apply_migration():
    try:
        conn = psycopg2.connect(
            host=DB_HOST,
            port=DB_PORT,
            dbname=DB_NAME,
            user=DB_USER,
            password=DB_PASSWORD
        )
        cur = conn.cursor()
        
        print("Applying migration: Add hint column to questions table")
        cur.execute("ALTER TABLE questions ADD COLUMN IF NOT EXISTS hint TEXT;")
        
        conn.commit()
        cur.close()
        conn.close()
        print("Migration applied successfully.")
    except Exception as e:
        print(f"Error applying migration: {e}")

if __name__ == "__main__":
    apply_migration()
