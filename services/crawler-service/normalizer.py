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
        if has_word("java", text) or "jvm" in text:
            return "Java"
        if has_word("javascript", text) or has_word("typescript", text) or "react" in text or "vue" in text:
            return "JavaScript/TypeScript"
        if "system design" in text or "scalability" in text:
            return "System Design"
        if "database" in text or "sql" in text:
            return "Database"
        if "kubernetes" in text or "docker" in text or "ci/cd" in text:
            return "DevOps"
        
        return "General"

    @staticmethod
    def detect_level(title: str, content: str) -> str:
        # Very naive heuristic
        text = (title + " " + content).lower()
        
        if any(w in text for w in ["fresher", "intern", "new grad", "mới tốt nghiệp", "thực tập"]):
            return "Fresher"
        if any(w in text for w in ["basic", "what is", "define", "cơ bản", "là gì"]):
            return "Junior"
        if any(w in text for w in ["design", "architecture", "trade-off", "optimize", "tối ưu", "thiết kế"]):
            return "Senior"
            
        return "Mid"

    @staticmethod
    def detect_role(title: str, content: str) -> str:
        text = (title + " " + content).lower()
        
        if any(w in text for w in ["frontend", "front-end", "react", "vue", "angular", "css", "html", "browser"]):
            return "FrontEnd"
        if any(w in text for w in ["backend", "back-end", "api", "database", "server", "golang", "java", "python"]):
            return "BackEnd"
        if any(w in text for w in ["devops", "sre", "cloud", "aws", "docker", "kubernetes", "ci/cd"]):
            return "DevOps"
        if any(w in text for w in ["data engineer", "etl", "spark", "hadoop", "warehouse", "big data"]):
            return "Data Engineer"
            
        return "BackEnd" # Default
