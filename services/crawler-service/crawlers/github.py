import re
from typing import Dict, Any
from .base import BaseCrawler

class GitHubCrawler(BaseCrawler):
    async def can_handle(self, url: str) -> bool:
        return "github.com" in url

    async def extract(self, url: str) -> Dict[str, Any]:
        # Handle GitHub URLs
        # 1. Blob/File URL: https://github.com/user/repo/blob/main/README.md
        # 2. Repo Root URL: https://github.com/user/repo (defaults to README)
        
        # Simple strategy: Try to convert to raw content URL if it's a blob
        raw_url = url
        if "github.com" in url and "/blob/" in url:
            raw_url = url.replace("github.com", "raw.githubusercontent.com").replace("/blob/", "/")
        
        # If it's a root repo url, we might need to guess the README, but for now let's assume direct file links or handle generic HTML parsing if raw fails.
        
        try:
            content = await self.fetch_html(raw_url)
            # If we fetched raw content successfully, title is the filename or repo name
            title = url.split("/")[-1]
            if not title or title == "":
                title = "GitHub Content"
                
            return {
                "source": "github",
                "title": title,
                "content": content,
                "url": url
            }
        except Exception as e:
            # Fallback to HTML parsing if raw fetch fails (e.g. it's a repo root)
            # For now, let's return error or basic info
             return {
                "source": "github",
                "title": "Error fetching GitHub content",
                "content": str(e),
                "url": url,
                "error": True
            }
