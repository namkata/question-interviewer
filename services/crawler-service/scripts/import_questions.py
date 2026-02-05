import asyncio
import sys
import os
import httpx
import json

# Add parent directory to path to import crawlers
sys.path.append(os.path.join(os.path.dirname(__file__), '..'))

from crawlers.vietnam_market import VietnamMarketCrawler

QUESTION_SERVICE_URL = os.getenv("QUESTION_SERVICE_URL", "http://question-service:8080/api/v1")
SEED_USER_ID = "123e4567-e89b-12d3-a456-426614174000"

async def main():
    print("Starting import process...")
    
    # 1. Run Crawler
    crawler = VietnamMarketCrawler()
    # Crawl everything (no specific filter to get all sample data)
    questions = await crawler.extract("vn-market://questions")
    print(f"Crawled {len(questions)} questions.")

    async with httpx.AsyncClient() as client:
        # 2. Process and Import
        for q in questions:
            title = q["title"]
            print(f"Processing: {title}")
            
            # Determine Topic
            tags = q.get("meta_tags", [])
            topic_name = "General"
            if tags:
                topic_name = tags[0] # Use first tag as primary topic
            
            topic_id = await get_or_create_topic(client, topic_name)
            if not topic_id:
                print(f"Failed to get topic for {title}, skipping.")
                continue
                
            # Create Question
            payload = {
                "title": title,
                "content": q["content"],
                "level": q["meta_level"],
                "language": detect_language(title, q["content"]), # Simple detection or default
                "role": q["meta_role"],
                "hint": q.get("hint", ""),
                "correct_answer": q.get("correct_answer", ""),
                "topic_id": topic_id,
                "created_by": SEED_USER_ID
            }
            
            await create_question(client, payload)

async def get_or_create_topic(client, name):
    # Try to get
    try:
        resp = await client.get(f"{QUESTION_SERVICE_URL}/topics", params={"name": name})
        if resp.status_code == 200:
            data = resp.json()
            if data and len(data) > 0:
                return data[0]["id"]
    except Exception as e:
        print(f"Error getting topic {name}: {e}")
        return None

    # Create
    try:
        print(f"Creating topic: {name}")
        resp = await client.post(f"{QUESTION_SERVICE_URL}/topics", json={
            "name": name,
            "description": f"Topic for {name}"
        })
        if resp.status_code == 201:
            return resp.json()["id"]
        else:
            print(f"Failed to create topic {name}: {resp.text}")
            return None
    except Exception as e:
        print(f"Error creating topic {name}: {e}")
        return None

async def create_question(client, payload):
    try:
        resp = await client.post(f"{QUESTION_SERVICE_URL}/questions", json=payload)
        if resp.status_code == 201:
            print(f"Successfully imported: {payload['title']}")
        else:
            print(f"Failed to import {payload['title']}: {resp.text}")
    except Exception as e:
        print(f"Error importing question {payload['title']}: {e}")

def detect_language(title, content):
    # Simple heuristic
    text = (title + content).lower()
    if "vietnamese" in text or "tiếng việt" in text or "là gì" in text:
        return "vi"
    return "en" # Default

if __name__ == "__main__":
    asyncio.run(main())
