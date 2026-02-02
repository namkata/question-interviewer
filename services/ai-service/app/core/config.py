from pydantic_settings import BaseSettings

class Settings(BaseSettings):
    PROJECT_NAME: str = "AI Service"
    API_V1_STR: str = "/api/v1"
    
    # LLM Settings
    GOOGLE_API_KEY: str = ""
    OPENAI_API_KEY: str = ""
    GROQ_API_KEY: str = ""
    
    # Default Provider: "google", "openai", or "groq"
    LLM_PROVIDER: str = "groq" 

    class Config:
        env_file = ".env"

settings = Settings()
