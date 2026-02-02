from pydantic import BaseModel, Field
from typing import Optional, List

class EvaluationRequest(BaseModel):
    question_content: str = Field(..., description="The content of the question asked")
    user_answer: str = Field(..., description="The answer provided by the user")
    correct_answer: Optional[str] = Field(None, description="The correct answer or key points (optional)")
    topic: Optional[str] = Field(None, description="Topic of the question")
    level: Optional[str] = Field(None, description="Difficulty level")
    language: Optional[str] = Field("en", description="Language for feedback (en or vi)")

class EvaluationResponse(BaseModel):
    score: int = Field(..., description="Score from 0 to 100")
    feedback: str = Field(..., description="Detailed feedback on the answer")
    suggestions: List[str] = Field(..., description="List of suggestions for improvement")
    improved_answer: Optional[str] = Field(None, description="An example of a better answer")

class GenerationRequest(BaseModel):
    topic: str = Field(..., description="Topic to generate questions for")
    count: int = Field(5, description="Number of questions to generate")
    level: str = Field("Medium", description="Difficulty level")

class GeneratedQuestion(BaseModel):
    content: str
    correct_answer: str
    topic: str
    level: str

class GenerationResponse(BaseModel):
    questions: List[GeneratedQuestion]
