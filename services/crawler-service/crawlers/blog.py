from typing import Dict, Any
from bs4 import BeautifulSoup
from .base import BaseCrawler

class BlogCrawler(BaseCrawler):
    async def can_handle(self, url: str) -> bool:
        # Generic handler for known blogs
        domains = ["medium.com", "dev.to", "engineering.atspotify.com", "eng.uber.com", "netflixtechblog.com"]
        return any(d in url for d in domains) or "blog" in url

    async def extract(self, url: str) -> Dict[str, Any]:
        html = await self.fetch_html(url)
        soup = BeautifulSoup(html, 'lxml')

        # Heuristic for title
        title = ""
        if soup.h1:
            title = soup.h1.get_text().strip()
        elif soup.title:
            title = soup.title.get_text().strip()
        
        # Heuristic for content (very basic)
        # 1. Try to find article tag
        article = soup.find('article')
        if article:
            content = article.get_text(separator="\n\n").strip()
        else:
            # Fallback to main tag
            main = soup.find('main')
            if main:
                content = main.get_text(separator="\n\n").strip()
            else:
                # Fallback to body (messy)
                content = soup.body.get_text(separator="\n\n").strip() if soup.body else ""

        return {
            "source": "blog",
            "title": title,
            "content": content[:10000], # Limit content length for safety
            "url": url
        }
