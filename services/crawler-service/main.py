from fastapi import FastAPI, HTTPException, Depends
from pydantic import BaseModel
from typing import Optional, List, Union, Any
from sqlalchemy.orm import Session

from database import get_db, engine, Base
from models import CrawledQuestion, Question, Topic
from normalizer import Normalizer

from crawlers.github import GitHubCrawler
from crawlers.blog import BlogCrawler
from crawlers.generic import GenericCrawler
from crawlers.vietnam_market import VietnamMarketCrawler

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
    detected_role: Optional[str]
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
        VietnamMarketCrawler(),
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
        raw_data_result = await selected_crawler.extract(url)
        
        # Ensure list for consistent processing
        items = raw_data_result if isinstance(raw_data_result, list) else [raw_data_result]
        
        saved_items = []
        
        for raw_data in items:
            # 3. Normalize & Classify
            normalized_content = Normalizer.normalize_content(raw_data.get("content", ""))
            
            # Use metadata from crawler if available, otherwise detect
            detected_topic = raw_data.get("meta_topic")
            if not detected_topic:
                detected_topic = Normalizer.detect_topic(raw_data.get("title", ""), normalized_content)
            
            detected_level = raw_data.get("meta_level")
            if not detected_level:
                detected_level = Normalizer.detect_level(raw_data.get("title", ""), normalized_content)
                
            detected_role = raw_data.get("meta_role")
            if not detected_role:
                detected_role = Normalizer.detect_role(raw_data.get("title", ""), normalized_content)
            
            # 4. Store to Staging DB
            db_item = CrawledQuestion(
                source=raw_data.get("source"),
                raw_title=raw_data.get("title"),
                raw_content=normalized_content,
                url=raw_data.get("url"),
                detected_topic=detected_topic,
                detected_role=detected_role,
                detected_level=detected_level,
                status="pending"
            )
            db.add(db_item)
            saved_items.append(db_item)
            
        db.commit()
        
        return {"message": f"Successfully crawled {len(saved_items)} items"}
        
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))

@app.post("/crawled-questions/{question_id}/approve")
def approve_question(question_id: str, db: Session = Depends(get_db)):
    # Find crawled question
    crawled = db.query(CrawledQuestion).filter(CrawledQuestion.id == question_id).first()
    if not crawled:
        raise HTTPException(status_code=404, detail="Crawled question not found")
    
    if crawled.status == "approved":
        return {"message": "Question already approved"}
        
    # Find or Create Topic
    topic_name = crawled.detected_topic or "General"
    topic = db.query(Topic).filter(Topic.name == topic_name).first()
    if not topic:
        topic = Topic(name=topic_name, description=f"Topic {topic_name}")
        db.add(topic)
        db.flush() # get ID
        
    # Determine Language
    language = "en"
    if crawled.source == "Vietnam IT Job Market":
        language = "vi"
        
    # Create Question
    new_question = Question(
        title=crawled.raw_title,
        content=crawled.raw_content or "",
        level=crawled.detected_level,
        language=language,
        role=crawled.detected_role,
        topic_id=topic.id,
        status="published"
    )
    db.add(new_question)
    
    # Update Status
    crawled.status = "approved"
    
    db.commit()
    
    return {"message": "Question approved and published", "question_id": str(new_question.id)}
