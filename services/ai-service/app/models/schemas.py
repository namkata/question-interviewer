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
    score: int = Field(..., description="Score from 1 to 10")
    feedback: str = Field(..., description="Short feedback on the answer")
    suggestions: List[str] = Field(..., description="List of suggestions for improvement")
    improved_answer: Optional[str] = Field(None, description="An example of a better answer")

class InterviewSummaryAttempt(BaseModel):
    question_content: str
    user_answer: str
    score: int
    feedback: str
    round_name: Optional[str] = None
    round_index: Optional[int] = None

class InterviewSummaryRequest(BaseModel):
    role: Optional[str] = Field(None, description="Role of the interview (optional)")
    language: Optional[str] = Field("en", description="Language for summary (en or vi)")
    attempts: List[InterviewSummaryAttempt]

class InterviewSummaryResponse(BaseModel):
    strengths: str
    weaknesses: str
    readiness: str
    overall_score: int = Field(..., description="Overall readiness score from 1 to 10")

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
