import json
import urllib.request
import urllib.error
import time

AI_SERVICE_URL = "http://localhost:8000/api/v1/generate"
PRACTICE_SERVICE_URL = "http://localhost:8080/api/v1/practice/questions"

TOPICS = {
    "Leadership": ["Leadership Principles", "Conflict Resolution", "Team Management", "Mentorship"],
    "Behavioral": ["Resume Walkthrough", "Strengths and Weaknesses", "Project Experience", "Motivation"],
    "Algorithms": ["Arrays", "LinkedLists", "Trees", "Graphs", "Sorting", "Searching", "Dynamic Programming"],
    "System Design": ["Scalability", "Load Balancing", "Database Design", "Caching", "Microservices"],
    "Golang": ["Goroutines", "Channels", "Interfaces", "Error Handling", "Garbage Collection"],
    "Python": ["GIL", "Decorators", "Generators", "Memory Management", "Asyncio"],
    "NodeJS": ["Event Loop", "Streams", "Modules", "Buffers", "Async/Await"],
    "NestJS": ["Dependency Injection", "Modules", "Controllers", "Providers", "Middleware", "Guards", "Interceptors"]
}

LEVELS = ["Junior", "Medium", "Senior"]

def generate_questions(topic, count=3, level="Medium"):
    # Optionally use subtopics to make it more specific, but for now just pass the main topic
    payload = {
        "topic": topic,
        "count": count,
        "level": level
    }
    data = json.dumps(payload).encode('utf-8')
    req = urllib.request.Request(AI_SERVICE_URL, data=data, headers={'Content-Type': 'application/json'})
    
    try:
        with urllib.request.urlopen(req) as response:
            if response.status == 200:
                return json.loads(response.read().decode('utf-8'))
    except urllib.error.URLError as e:
        print(f"Error generating questions for {topic}: {e}")
        return None
    return None

def save_question(question_data):
    payload = {
        "content": question_data["content"],
        "topic": question_data["topic"],
        "level": question_data["level"],
        "correct_answer": question_data.get("correct_answer", "")
    }
    data = json.dumps(payload).encode('utf-8')
    req = urllib.request.Request(PRACTICE_SERVICE_URL, data=data, headers={'Content-Type': 'application/json'})
    
    try:
        with urllib.request.urlopen(req) as response:
            if response.status == 201:
                print(f"Saved: {question_data['content'][:30]}...")
                return True
    except urllib.error.URLError as e:
        print(f"Error saving question: {e}")
        return False
    return False

def main():
    print("Starting data crawl/generation...")
    
    for topic in TOPICS:
        for level in LEVELS:
            print(f"\nGenerating {level} questions for {topic}...")
            response = generate_questions(topic, count=3, level=level)
            
            if response and "questions" in response:
                questions = response["questions"]
                print(f"Got {len(questions)} questions. Saving...")
                
                for q in questions:
                    save_question(q)
            else:
                print("No questions generated.")
            
            # Sleep to avoid rate limits
            time.sleep(1)

    print("\nDone!")

if __name__ == "__main__":
    main()
