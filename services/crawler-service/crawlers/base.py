from abc import ABC, abstractmethod
from typing import Optional, Dict, Any, Union, List
import httpx

class BaseCrawler(ABC):
    @abstractmethod
    async def can_handle(self, url: str) -> bool:
        pass

    @abstractmethod
    async def extract(self, url: str) -> Union[Dict[str, Any], List[Dict[str, Any]]]:
        """
        Returns a dictionary or list of dictionaries with:
        - title: str
        - content: str (markdown or text)
        - source: str
        - url: str
        """
        pass

    async def fetch_html(self, url: str) -> str:
        async with httpx.AsyncClient(follow_redirects=True) as client:
            # Fake User-Agent to avoid 403s
            headers = {
                "User-Agent": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36"
            }
            resp = await client.get(url, headers=headers)
            resp.raise_for_status()
            return resp.text
