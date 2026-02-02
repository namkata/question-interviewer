from fastapi import FastAPI, HTTPException, Depends
from pydantic import BaseModel
from typing import Optional, List
from sqlalchemy.orm import Session

from database import get_db, engine, Base
from models import CrawledQuestion
from normalizer import Normalizer

from crawlers.github import GitHubCrawler
from crawlers.blog import BlogCrawler
from crawlers.generic import GenericCrawler

# Create tables if not exist (for dev simplicity, usually managed by migrations)
# Base.metadata.create_all(bind=engine)

app = FastAPI()

class CrawlRequest(BaseModel):
    url: str

class CrawledQuestionResponse(BaseModel):
    id: str
    title: str
    source: str
    detected_topic: Optional[str]
    detected_level: Optional[str]
    status: str
    
    class Config:
        orm_mode = True

@app.get("/")
def read_root():
    return {"Hello": "Crawler Service"}

@app.post("/crawl")
async def crawl_url(request: CrawlRequest, db: Session = Depends(get_db)):
    url = request.url
    
    # 1. Factory Pattern to select crawler
    crawlers = [
        GitHubCrawler(),
        BlogCrawler(),
        GenericCrawler()
    ]
    
    selected_crawler = None
    for crawler in crawlers:
        if await crawler.can_handle(url):
            selected_crawler = crawler
            break
    
    if not selected_crawler:
        raise HTTPException(status_code=400, detail="No suitable crawler found")
        
    try:
        # 2. Extract
        raw_data = await selected_crawler.extract(url)
        
        # 3. Normalize & Classify
        normalized_content = Normalizer.normalize_content(raw_data.get("content", ""))
        detected_topic = Normalizer.detect_topic(raw_data.get("title", ""), normalized_content)
        detected_level = Normalizer.detect_level(raw_data.get("title", ""), normalized_content)
        
        # 4. Store to Staging DB
        db_item = CrawledQuestion(
            source=raw_data.get("source"),
            raw_title=raw_data.get("title"),
            raw_content=normalized_content,
            url=raw_data.get("url"),
            detected_topic=detected_topic,
            detected_level=detected_level,
            status="pending"
        )
        db.add(db_item)
        db.commit()
        db.refresh(db_item)
        
        return {"status": "success", "data": {
            "id": str(db_item.id),
            "title": db_item.raw_title,
            "topic": db_item.detected_topic,
            "level": db_item.detected_level
        }}
        
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))

@app.get("/crawled-questions", response_model=List[CrawledQuestionResponse])
def list_crawled_questions(status: str = "pending", limit: int = 10, db: Session = Depends(get_db)):
    questions = db.query(CrawledQuestion).filter(CrawledQuestion.status == status).limit(limit).all()
    # Manual conversion to pydantic friendly dict if needed, or rely on orm_mode
    return [
        CrawledQuestionResponse(
            id=str(q.id),
            title=q.raw_title,
            source=q.source,
            detected_topic=q.detected_topic,
            detected_level=q.detected_level,
            status=q.status
        ) for q in questions
    ]

@app.post("/crawled-questions/{question_id}/approve")
def approve_question(question_id: str, db: Session = Depends(get_db)):
    # In a real system, this might trigger a call to Question Service to create a real question
    # For now, we just mark it as approved in staging
    question = db.query(CrawledQuestion).filter(CrawledQuestion.id == question_id).first()
    if not question:
        raise HTTPException(status_code=404, detail="Question not found")
    
    question.status = "approved"
    db.commit()
    return {"status": "approved", "id": question_id}

@app.post("/crawled-questions/{question_id}/reject")
def reject_question(question_id: str, db: Session = Depends(get_db)):
    question = db.query(CrawledQuestion).filter(CrawledQuestion.id == question_id).first()
    if not question:
        raise HTTPException(status_code=404, detail="Question not found")
    
    question.status = "rejected"
    db.commit()
    return {"status": "rejected", "id": question_id}
