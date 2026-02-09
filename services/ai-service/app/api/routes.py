import json

from fastapi import APIRouter, HTTPException, Depends
from app.models.schemas import (
    EvaluationRequest,
    EvaluationResponse,
    GenerationRequest,
    GenerationResponse,
    GeneratedQuestion,
    InterviewSummaryRequest,
    InterviewSummaryResponse,
)
from app.services.evaluator import AnswerEvaluator

router = APIRouter()

def get_evaluator():
    return AnswerEvaluator()

@router.post("/evaluate", response_model=EvaluationResponse)
async def evaluate_answer(
    request: EvaluationRequest,
    evaluator: AnswerEvaluator = Depends(get_evaluator)
):
    """
    Evaluate a candidate's answer using AI.
    """
    try:
        response = await evaluator.evaluate(
            question=request.question_content,
            user_answer=request.user_answer,
            correct_answer=request.correct_answer,
            topic=request.topic,
            level=request.level,
            language=request.language
        )
        return response
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))

@router.post("/generate", response_model=GenerationResponse)
async def generate_questions(
    request: GenerationRequest,
    evaluator: AnswerEvaluator = Depends(get_evaluator)
):
    """
    Generate interview questions using AI.
    """
    try:
        questions = await evaluator.generate_questions(
            topic=request.topic,
            count=request.count,
            level=request.level
        )
        return GenerationResponse(questions=[GeneratedQuestion(**q) for q in questions])
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))

@router.post("/summarize", response_model=InterviewSummaryResponse)
async def summarize_interview(
    request: InterviewSummaryRequest,
    evaluator: AnswerEvaluator = Depends(get_evaluator)
):
    try:
        attempts_json = json.dumps(
            [(a.model_dump() if hasattr(a, "model_dump") else a.dict()) for a in request.attempts],
            ensure_ascii=False
        )
        result = await evaluator.summarize_interview(
            role=request.role or "General",
            attempts_json=attempts_json,
            language=request.language or "en"
        )
        return InterviewSummaryResponse(**result)
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))
