import re

class Normalizer:
    @staticmethod
    def normalize_content(content: str) -> str:
        if not content:
            return ""
        
        # 1. Remove excessive whitespace
        content = re.sub(r'\n{3,}', '\n\n', content)
        content = content.strip()
        
        # 2. Convert common HTML entities if not handled by BS4
        # (BeautifulSoup usually handles this, but good to be safe)
        
        # 3. Standardize code blocks (e.g. ```go -> ```go)
        # TODO: Detect language if missing
        
        return content

    @staticmethod
    def detect_topic(title: str, content: str) -> str:
        # Simple rule-based topic detection
        text = (title + " " + content).lower()
        
        # Helper for whole word matching
        def has_word(word, text):
            return re.search(r'\b' + re.escape(word) + r'\b', text) is not None

        if has_word("golang", text) or has_word("go", text):
            return "Golang"
        if has_word("python", text) or "list comprehension" in text:
            return "Python"
        if "system design" in text or "scalability" in text:
            return "System Design"
        if "database" in text or "sql" in text:
            return "Database"
        
        return "General"

    @staticmethod
    def detect_level(title: str, content: str) -> str:
        # Very naive heuristic
        text = (title + " " + content).lower()
        
        if any(w in text for w in ["basic", "what is", "define"]):
            return "Junior"
        if any(w in text for w in ["design", "architecture", "trade-off", "optimize"]):
            return "Senior"
            
        return "Mid"
