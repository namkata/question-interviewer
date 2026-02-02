import uuid
from sqlalchemy import Column, String, Text, DateTime, func
from sqlalchemy.dialects.postgresql import UUID
from database import Base

class CrawledQuestion(Base):
    __tablename__ = "crawled_questions"

    id = Column(UUID(as_uuid=True), primary_key=True, default=uuid.uuid4)
    source = Column(String, nullable=False)
    raw_title = Column(Text, nullable=False)
    raw_content = Column(Text, nullable=True)
    url = Column(String, nullable=True)
    detected_topic = Column(String, nullable=True)
    detected_level = Column(String, nullable=True)
    status = Column(String, default="pending") # pending, approved, rejected
    created_at = Column(DateTime(timezone=True), server_default=func.now())
