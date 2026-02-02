import asyncio
import sys
import os

# Add current dir to path to import modules
sys.path.append(os.path.join(os.getcwd(), 'services/crawler-service'))

from crawlers.blog import BlogCrawler
from normalizer import Normalizer

async def main():
    url = "https://dev.to/finalroundai/40-python-interview-questions-for-2025-how-many-can-you-answer-9k6"
    print(f"Crawling {url}...")
    
    crawler = BlogCrawler()
    if await crawler.can_handle(url):
        print("Crawler: BlogCrawler selected.")
        try:
            raw_data = await crawler.extract(url)
            print("\n--- Raw Data Extracted ---")
            print(f"Title: {raw_data.get('title')}")
            
            normalized_content = Normalizer.normalize_content(raw_data.get("content", ""))
            detected_topic = Normalizer.detect_topic(raw_data.get("title", ""), normalized_content)
            detected_level = Normalizer.detect_level(raw_data.get("title", ""), normalized_content)
            
            print("\n--- Normalized & Classified ---")
            print(f"Topic: {detected_topic}")
            print(f"Level: {detected_level}")
            print(f"Content Preview (first 500 chars):\n{normalized_content[:500]}...")
            
        except Exception as e:
            print(f"Error: {e}")
            import traceback
            traceback.print_exc()
    else:
        print("No crawler found for this URL.")

if __name__ == "__main__":
    asyncio.run(main())
