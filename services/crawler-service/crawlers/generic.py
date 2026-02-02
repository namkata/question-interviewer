from typing import Dict, Any
from bs4 import BeautifulSoup
from .base import BaseCrawler

class GenericCrawler(BaseCrawler):
    async def can_handle(self, url: str) -> bool:
        return True # Fallback

    async def extract(self, url: str) -> Dict[str, Any]:
        html = await self.fetch_html(url)
        soup = BeautifulSoup(html, 'lxml')

        title = soup.title.get_text().strip() if soup.title else url
        
        # Try to find common content containers
        content_container = soup.find('article') or soup.find('main') or soup.find('div', class_='content') or soup.body
        
        content = ""
        if content_container:
            content = content_container.get_text(separator="\n").strip()

        return {
            "source": "generic",
            "title": title,
            "content": content[:5000],
            "url": url
        }
