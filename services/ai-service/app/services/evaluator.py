try:
    from langchain_google_genai import ChatGoogleGenerativeAI
except ImportError:
    ChatGoogleGenerativeAI = None
try:
    from langchain_openai import ChatOpenAI
except ImportError:
    ChatOpenAI = None
try:
    from langchain_groq import ChatGroq
except ImportError:
    ChatGroq = None
from langchain_core.prompts import PromptTemplate
from langchain_core.output_parsers import JsonOutputParser
from pydantic import BaseModel, Field
from app.core.config import settings
from app.models.schemas import EvaluationResponse
from typing import List, Optional

# Define the internal Pydantic model for LangChain parser (v1 compatible if needed, but let's try to match schema)
class EvaluationOutput(BaseModel):
    score: int = Field(description="Score from 0 to 100")
    feedback: str = Field(description="Detailed feedback on the answer")
    suggestions: List[str] = Field(description="List of suggestions for improvement")
    improved_answer: str = Field(description="An example of a better answer")

class AnswerEvaluator:
    def __init__(self):
        self.llm = self._get_llm()
        self.parser = JsonOutputParser(pydantic_object=EvaluationOutput)
        self.prompt = PromptTemplate(
            template="""
            You are an expert technical interviewer. Evaluate the candidate's answer to the following interview question.
            
            Question: {question}
            Topic: {topic}
            Level: {level}
            
            Candidate's Answer: {user_answer}
            
            Correct Answer / Key Points (Reference): {correct_answer}
            
            Provide a fair score (0-100), detailed feedback explaining what was good and what was missing, 
            concrete suggestions for improvement, and an example of a better/ideal answer.

            If the Candidate's Answer is empty or "N/A", treat this as a request for a sample answer:
            - Give a score of 0.
            - In feedback, explain the key points that a strong answer should cover.
            - In improved_answer, provide a complete, high-quality sample answer the candidate can learn from.
            
            SCORING GUIDELINES:
            - If the answer matches the Key Points/Correct Answer substantially, give a high score (90-100).
            - If the answer is perfect or near-perfect, give 100.
            - Do NOT artificially cap the score at 80. Reward good answers.

            IMPORTANT: The response (feedback, suggestions, improved_answer) MUST be in {language} language.
            
            {format_instructions}
            """,
            input_variables=["question", "topic", "level", "user_answer", "correct_answer", "language"],
            partial_variables={"format_instructions": self.parser.get_format_instructions()},
        )
        self.chain = self.prompt | self.llm | self.parser

    def _get_llm(self):
        if settings.LLM_PROVIDER == "openai":
            if ChatOpenAI is None:
                raise ValueError("langchain-openai is not installed. Please install it to use OpenAI provider.")
            return ChatOpenAI(api_key=settings.OPENAI_API_KEY, model="gpt-3.5-turbo")
        elif settings.LLM_PROVIDER == "groq":
            if ChatGroq is None:
                raise ValueError("langchain-groq is not installed. Please install it to use Groq provider.")
            return ChatGroq(api_key=settings.GROQ_API_KEY, model="llama-3.3-70b-versatile")
        else:
            return ChatGoogleGenerativeAI(google_api_key=settings.GOOGLE_API_KEY, model="gemini-pro")

    async def evaluate(self, question: str, user_answer: str, correct_answer: str = None, topic: str = "General", level: str = "Medium", language: str = "en") -> EvaluationResponse:
        try:
            # Map language code to full name for clearer prompt
            lang_map = {"vi": "Vietnamese", "en": "English"}
            full_lang = lang_map.get(language, "English")

            result = await self.chain.ainvoke({
                "question": question,
                "user_answer": user_answer,
                "correct_answer": correct_answer if correct_answer else "N/A - Use your expert knowledge",
                "topic": topic if topic else "General",
                "level": level if level else "Medium",
                "language": full_lang
            })
            
            return EvaluationResponse(
                score=result.get("score", 0),
                feedback=result.get("feedback", "No feedback generated"),
                suggestions=result.get("suggestions", []),
                improved_answer=result.get("improved_answer", "")
            )
        except Exception as e:
            print(f"Error evaluating answer: {e}")
            # Fallback in case of parsing error or LLM failure
            return EvaluationResponse(
                score=0,
                feedback=f"Failed to evaluate answer due to error: {str(e)}",
                suggestions=["Please try again later."],
                improved_answer=""
            )

    async def generate_questions(self, topic: str, count: int = 5, level: str = "Medium") -> List[dict]:
        try:
            gen_parser = JsonOutputParser()
            gen_prompt = PromptTemplate(
                template="""
                You are a technical interview expert. Generate {count} interview questions for the topic "{topic}" at {level} level.
                
                Return the result as a JSON array of objects, where each object has:
                - "content": The question text
                - "correct_answer": A brief summary of the correct answer
                - "topic": "{topic}"
                - "level": "{level}"
                
                Example output format:
                [
                    {{
                        "content": "What is X?",
                        "correct_answer": "X is Y...",
                        "topic": "{topic}",
                        "level": "{level}"
                    }}
                ]
                
                Ensure the JSON is valid.
                """,
                input_variables=["topic", "count", "level"],
            )
            
            chain = gen_prompt | self.llm | gen_parser
            
            result = await chain.ainvoke({
                "topic": topic,
                "count": count,
                "level": level
            })
            
            # Ensure result is a list
            if isinstance(result, list):
                return result
            elif isinstance(result, dict) and "questions" in result:
                return result["questions"]
            else:
                return []
                
        except Exception as e:
            print(f"Error generating questions: {e}")
            return []
