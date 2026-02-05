import uuid
from sqlalchemy import Column, String, Text, DateTime, func, ForeignKey
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
    detected_role = Column(String, nullable=True)
    detected_level = Column(String, nullable=True)
    status = Column(String, default="pending") # pending, approved, rejected
    created_at = Column(DateTime(timezone=True), server_default=func.now())

class Topic(Base):
    __tablename__ = "topics"

    id = Column(UUID(as_uuid=True), primary_key=True, default=uuid.uuid4)
    name = Column(String, unique=True, nullable=False)
    description = Column(Text, nullable=True)

class Question(Base):
    __tablename__ = "questions"

    id = Column(UUID(as_uuid=True), primary_key=True, default=uuid.uuid4)
    title = Column(String, nullable=False)
    content = Column(Text, nullable=False)
    level = Column(String, nullable=True)
    language = Column(String, default="en")
    role = Column(String, nullable=True)
    topic_id = Column(UUID(as_uuid=True), ForeignKey("topics.id"))
    created_by = Column(UUID(as_uuid=True), nullable=True) # Nullable for auto-generated
    status = Column(String, default="published")
    created_at = Column(DateTime(timezone=True), server_default=func.now())
    updated_at = Column(DateTime(timezone=True), server_default=func.now(), onupdate=func.now())
