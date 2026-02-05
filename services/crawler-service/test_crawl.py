import asyncio
import sys
import os

# Add current dir to path to import modules
sys.path.append(os.path.join(os.getcwd(), 'services/crawler-service'))

from crawlers.vietnam_market import VietnamMarketCrawler

async def main():
    crawler = VietnamMarketCrawler()
    
    test_urls = [
        "vn-market://backend",
        "vn-market://questions?role=BackEnd&level=Mid",
        "vn-market://questions?lang=Golang",
        "vn-market://questions?lang=Python",
        "vn-market://questions?role=DevOps&level=Senior",
        # Multi-params: multiple langs
        "vn-market://questions?role=BackEnd&level=Mid&lang=Python&lang=Golang",
        # Multi-params: multiple roles (BackEnd + DevOps) and lang filter
        "vn-market://questions?role=BackEnd&role=DevOps&level=Junior&lang=Kubernetes",
        # Path-based role with lang filter
        "vn-market://frontend?lang=TypeScript"
    ]
    
    print("=== Testing VietnamMarketCrawler Filters ===")
    
    for url in test_urls:
        print(f"\nTesting URL: {url}")
        if await crawler.can_handle(url):
            results = await crawler.extract(url)
            print(f"Items fetched: {len(results)}")
            for i, item in enumerate(results):
                print(f"  Item {i+1}: {item['title']} | Role: {item.get('meta_role')} | Level: {item.get('meta_level')} | Tags: {item.get('meta_tags')}")
        else:
            print("  Cannot handle this URL")

if __name__ == "__main__":
    asyncio.run(main())
