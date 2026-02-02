import React, { useState } from 'react';
import axios from 'axios';
import { Send, Play, BookOpen, AlertCircle, Settings, Globe, Code, Briefcase, Mic, MicOff } from 'lucide-react';

// Types
interface Session {
  id: string;
  user_id: string;
  status: string;
  score: number;
}

interface Question {
  id: string;
  content: string;
  topic: string;
  level: string;
}

interface Attempt {
  id: string;
  score: number;
  feedback: string;
}

const API_URL = 'http://localhost:8081/api/v1';

// Translations
const TRANSLATIONS = {
  en: {
    title: "Mock Interview Practice",
    subtitle: "Practice your technical interview answers with AI feedback",
    setup_title: "Interview Setup",
    language_label: "Language",
    stack_label: "Select Tech Stack / Field",
    round_label: "Select Interview Round",
    start_session: "Start Session",
    starting: "Starting...",
    question_label: "Question",
    your_answer: "Your Answer",
    placeholder_answer: "Type your answer here...",
    submit: "Submit Answer",
    evaluating: "Evaluating...",
    next_question: "Next Question",
    ai_evaluation: "AI Evaluation",
    feedback_detail: "Detailed Feedback",
    good_job: "Great job! Keep it up.",
    needs_improvement: "Good attempt, but there is room for improvement.",
    error_start: "Failed to start session",
    error_fetch: "Failed to fetch question",
    error_submit: "Failed to submit answer",
    enable_voice: "Enable Voice Input",
    listening: "Listening...",
    click_to_speak: "Click to speak",
    rounds: {
      recruiter: "Round 1: Recruiter Screen",
      technical: "Round 2: Technical Phone Screen",
      ds_algo: "Round 3: Data Structures & Algorithms",
      system_design: "Round 3: System Design",
      leadership: "Round 3: Googleyness & Leadership"
    },
    suggestions_label: "Suggestions for improvement",
    download_summary: "Download Summary",
    avg_label: "Avg",
    threshold_label: "Threshold",
    pass_label: "Pass",
    fail_label: "Fail",
    overall_passed: "Result: Passed Interview",
    overall_failed: "Result: Failed Interview",
    return_setup: "Return to Setup",
    ai_enabled_label: "Enable AI Evaluation (Uncheck to use standard answers)"
  },
  vi: {
    title: "Phỏng Vấn Thử Nghiệm",
    subtitle: "Luyện tập trả lời phỏng vấn kỹ thuật với phản hồi từ AI",
    setup_title: "Thiết lập phỏng vấn",
    language_label: "Ngôn ngữ",
    stack_label: "Chọn Công Nghệ / Lĩnh Vực",
    round_label: "Chọn Vòng Phỏng Vấn",
    start_session: "Bắt Đầu Phỏng Vấn",
    starting: "Đang bắt đầu...",
    question_label: "Câu hỏi",
    your_answer: "Câu trả lời của bạn",
    placeholder_answer: "Nhập câu trả lời của bạn vào đây...",
    submit: "Gửi câu trả lời",
    evaluating: "Đang đánh giá...",
    next_question: "Câu hỏi tiếp theo",
    ai_evaluation: "Đánh giá từ AI",
    feedback_detail: "Chi tiết phản hồi",
    good_job: "Làm tốt lắm! Hãy tiếp tục phát huy.",
    needs_improvement: "Nỗ lực tốt, nhưng vẫn còn chỗ cần cải thiện.",
    error_start: "Không thể bắt đầu phiên",
    error_fetch: "Không thể tải câu hỏi",
    error_submit: "Không thể gửi câu trả lời",
    enable_voice: "Bật nhập liệu bằng giọng nói",
    listening: "Đang nghe...",
    click_to_speak: "Nhấn để nói",
    rounds: {
      recruiter: "Vòng 1: Sàng lọc (Recruiter Screen)",
      technical: "Vòng 2: Kỹ thuật (Technical Phone Screen)",
      ds_algo: "Vòng 3: Thuật toán & Cấu trúc dữ liệu",
      system_design: "Vòng 3: Thiết kế hệ thống",
      leadership: "Vòng 3: Googleyness & Leadership"
    },
    suggestions_label: "Gợi ý cải thiện",
    download_summary: "Tải báo cáo",
    avg_label: "Điểm TB",
    threshold_label: "Ngưỡng",
    pass_label: "Đạt",
    fail_label: "Không đạt",
    overall_passed: "Kết quả: Đậu phỏng vấn",
    overall_failed: "Kết quả: Trượt phỏng vấn",
    return_setup: "Quay lại thiết lập",
    ai_enabled_label: "Bật đánh giá AI (Bỏ chọn để dùng câu trả lời mẫu)"
  }
};

const STACKS = ["Golang", "Python", "NodeJS", "NestJS", "Design Patterns"];

const ROUNDS = [
  { id: "recruiter", labelKey: "recruiter", count: 3 },
  { id: "technical", labelKey: "technical", count: 2 },
  { id: "ds_algo", labelKey: "ds_algo", count: 2 },
  { id: "system_design", labelKey: "system_design", count: 1 },
  { id: "leadership", labelKey: "leadership", count: 3 }
];

const PracticeSession: React.FC = () => {
  const [session, setSession] = useState<Session | null>(null);
  const [question, setQuestion] = useState<Question | null>(null);
  const [answer, setAnswer] = useState('');
  const [attempt, setAttempt] = useState<Attempt | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');
  const [currentQuestionIndex, setCurrentQuestionIndex] = useState(0);
  const [totalQuestions, setTotalQuestions] = useState(0);
  const [isSessionComplete, setIsSessionComplete] = useState(false);
  const [nextQuestionId, setNextQuestionId] = useState<string | null>(null);
  const [isListening, setIsListening] = useState(false);
  const [recognition, setRecognition] = useState<any>(null);
  const [currentRoundIndex, setCurrentRoundIndex] = useState(0);
  const [roundStats, setRoundStats] = useState<{[key:string]: {sum:number, count:number}}>({});
  const [finalSummary, setFinalSummary] = useState<null | {overallPass: boolean, details: Array<{id:string, avg:number, threshold:number, pass:boolean}>}>(null);
  const [attemptsHistory, setAttemptsHistory] = useState<Array<{roundId:string; questionId:string; score:number; feedback:string}>>([]);

  // Setup State
  const [language, setLanguage] = useState<'en' | 'vi'>('en');
  const [stack, setStack] = useState('Golang');
  const [round, setRound] = useState('technical');
  const [step, setStep] = useState<'setup' | 'session'>('setup');
  const [enableVoice, setEnableVoice] = useState(false);
  const [aiEnabled, setAiEnabled] = useState(true);
  const [customThresholds, setCustomThresholds] = useState<{[key:string]: number}>({});

  const t = TRANSLATIONS[language];

  const getTechSpecificTips = (stackName: string, roundId: string): string[] => {
    if (language === 'en') {
       const techTipsEn: {[key: string]: {[key: string]: string[]}} = {
        "Golang": {
          "technical": ["Use goroutines/channels for concurrency", "Explain defer/panic/recover", "Understand interfaces and struct embedding"],
          "ds_algo": ["Optimize slice operations", "Use sync package for synchronization", "Avoid memory leaks with goroutines"],
          "system_design": ["Design microservices with Go kit or Gin", "Handle concurrent requests efficiently", "Go routines pooling"]
        },
        "Python": {
          "technical": ["Explain GIL and concurrency impact", "Use decorators for code optimization", "Memory management with garbage collection"],
          "ds_algo": ["Use list comprehensions for performance", "Optimize dictionary lookups", "Use generators for large datasets"],
          "system_design": ["Use Celery for async tasks", "Django/FastAPI scalability", "Python multiprocessing vs multithreading"]
        },
        "NodeJS": {
          "technical": ["Understand Event Loop and non-blocking I/O", "Use Streams for large data", "Manage memory leaks in Node.js"],
          "ds_algo": ["Optimize V8 engine execution", "Use Buffer efficiently", "Async/Await vs Promises patterns"],
          "system_design": ["Design scalable architecture with Clustering", "Microservices with NestJS/Express", "Use Redis for caching/pub-sub"]
        },
        "NestJS": {
          "technical": ["Dependency Injection in NestJS", "Use Interceptors, Guards, Pipes", "Module architecture patterns"],
          "ds_algo": ["Optimize TypeORM/Prisma queries", "Reactive programming with RxJS", "Custom decorators implementation"],
          "system_design": ["Microservices with NestJS Transport layers", "CQRS pattern implementation", "WebSocket gateways scaling"]
        }
      };
      return techTipsEn[stackName]?.[roundId] || [];
    }

    const techTips: {[key: string]: {[key: string]: string[]}} = {
      "Golang": {
        "technical": ["Sử dụng goroutines/channels cho concurrency", "Giải thích defer/panic/recover", "Hiểu rõ về interface và struct embedding"],
        "ds_algo": ["Tối ưu hóa slice operations", "Sử dụng sync package cho synchronization", "Tránh memory leaks với goroutines"],
        "system_design": ["Thiết kế microservices với Go kit hoặc Gin", "Xử lý concurrent requests hiệu quả", "Go routines pooling"]
      },
      "Python": {
        "technical": ["Giải thích GIL và tác động đến concurrency", "Sử dụng decorators để tối ưu code", "Quản lý memory với garbage collection"],
        "ds_algo": ["Sử dụng list comprehensions cho hiệu suất", "Tối ưu hóa dictionary lookups", "Sử dụng generators cho large datasets"],
        "system_design": ["Sử dụng Celery cho async tasks", "Django/FastAPI scalability", "Python multiprocessing vs multithreading"]
      },
      "NodeJS": {
        "technical": ["Hiểu rõ Event Loop và non-blocking I/O", "Sử dụng Streams cho dữ liệu lớn", "Quản lý memory leaks trong Node.js"],
        "ds_algo": ["Tối ưu hóa V8 engine execution", "Sử dụng Buffer hiệu quả", "Async/Await vs Promises patterns"],
        "system_design": ["Thiết kế scalable architecture với Clustering", "Microservices với NestJS/Express", "Sử dụng Redis cho caching/pub-sub"]
      },
      "NestJS": {
        "technical": ["Dependency Injection trong NestJS", "Sử dụng Interceptors, Guards, Pipes", "Module architecture patterns"],
        "ds_algo": ["Tối ưu hóa TypeORM/Prisma queries", "Reactive programming với RxJS", "Custom decorators implementation"],
        "system_design": ["Microservices với NestJS Transport layers", "CQRS pattern implementation", "WebSocket gateways scaling"]
      }
    };
    return techTips[stackName]?.[roundId] || [];
  };

  const getThresholdForRound = (rId: string) => {
    if (customThresholds[rId] !== undefined) {
      return customThresholds[rId];
    }
    switch (rId) {
      case 'recruiter': return 60;
      case 'leadership': return 60;
      case 'technical': return 65;
      case 'ds_algo': return 65;
      case 'system_design': return 70;
      default: return 65;
    }
  };

  const addScoreForRound = (rId: string, score: number) => {
    setRoundStats(prev => {
      const cur = prev[rId] || { sum: 0, count: 0 };
      return { ...prev, [rId]: { sum: cur.sum + score, count: cur.count + 1 } };
    });
  };

  const getTopicAndLevel = () => {
    let topicName = stack;
    let levelName = 'Medium';
    let questionCount = 2;

    const selectedRound = ROUNDS.find(r => r.id === round);
    if (selectedRound) {
        questionCount = selectedRound.count;
    }

    switch (round) {
      case 'recruiter':
        topicName = 'Behavioral';
        levelName = 'Junior';
        break;
      case 'technical':
        topicName = stack; 
        levelName = 'Medium';
        break;
      case 'ds_algo':
        topicName = 'Algorithms';
        levelName = 'Medium';
        break;
      case 'system_design':
        topicName = 'System Design';
        levelName = 'Senior';
        break;
      case 'leadership':
        topicName = 'Leadership';
        levelName = 'Senior';
        break;
      default:
        topicName = stack;
    }
    return { topicName, levelName, questionCount };
  };

  const startSession = async () => {
    setLoading(true);
    setError('');
    setNextQuestionId(null);
    setIsSessionComplete(false);
    setCurrentQuestionIndex(1);
    setFinalSummary(null);
    setAttemptsHistory([]);
    const idx = ROUNDS.findIndex(r => r.id === round);
    setCurrentRoundIndex(idx >= 0 ? idx : 0);
    
    const { topicName, levelName, questionCount } = getTopicAndLevel();
    setTotalQuestions(questionCount);

    try {
      const userId = '123e4567-e89b-12d3-a456-426614174000';
      const response = await axios.post(`${API_URL}/sessions`, {
        user_id: userId,
        topic_id: topicName,
        level: levelName
      });

      const { session, first_question_id } = response.data;
      setSession(session);
      setStep('session');
      await fetchQuestion(first_question_id);
    } catch (err: any) {
      setError(err.response?.data?.error || t.error_start);
    } finally {
      setLoading(false);
    }
  };

  const fetchQuestion = async (questionId: string) => {
    try {
      const response = await axios.get(`${API_URL}/questions/${questionId}`);
      setQuestion(response.data);
      setAnswer('');
      setAttempt(null);
    } catch (err: any) {
      setError(err.response?.data?.error || t.error_fetch);
    }
  };

  const submitAnswer = async () => {
    if (!session || !question || !answer.trim()) return;

    setLoading(true);
    setError('');
    try {
      const response = await axios.post(`${API_URL}/sessions/${session.id}/answers`, {
        question_id: question.id,
        content: answer,
        language: language,
        ai_enabled: aiEnabled
      });
      const { attempt, next_question_id } = response.data;
      setAttempt(attempt);
      setNextQuestionId(next_question_id);
      addScoreForRound(round, attempt.score);
      setAttemptsHistory(prev => [...prev, {roundId: round, questionId: question.id, score: attempt.score, feedback: attempt.feedback}]);
    } catch (err: any) {
      setError(err.response?.data?.error || t.error_submit);
    } finally {
      setLoading(false);
    }
  };

  const getImprovementTips = (rId: string): string[] => {
    const techTips = getTechSpecificTips(stack, rId);
    let baseTips: string[] = [];
    if (language === 'vi') {
      switch (rId) {
        case 'recruiter':
          baseTips = ['Trình bày câu chuyện rõ ràng, theo STAR', 'Liên hệ kinh nghiệm với JD', 'Nhấn mạnh động lực và giá trị phù hợp văn hoá'];
          break;
        case 'technical':
          baseTips = ['Cấu trúc câu trả lời theo Problem–Approach–Tradeoffs', 'Đưa ví dụ code ngắn gọn', 'Nêu rõ độ phức tạp và tối ưu hoá'];
          break;
        case 'ds_algo':
          baseTips = ['Phân tích độ phức tạp trước khi code', 'Vẽ test case biên', 'Giải thích vì sao chọn cấu trúc dữ liệu'];
          break;
        case 'system_design':
          baseTips = ['Bắt đầu từ requirements', 'Vẽ high-level architecture (API, storage, cache, queue)', 'Phân tích bottlenecks và scaling plan'];
          break;
        case 'leadership':
          baseTips = ['Minh hoạ bằng tình huống thực tế', 'Nêu rõ vai trò, quyết định và kết quả', 'Phản tư và bài học rút ra'];
          break;
        default:
          baseTips = [];
      }
    } else {
      switch (rId) {
        case 'recruiter':
          baseTips = ['Tell coherent stories using STAR', 'Map experience to the JD', 'Emphasize motivation and culture fit'];
          break;
        case 'technical':
          baseTips = ['Structure answers: Problem–Approach–Tradeoffs', 'Provide concise code snippets', 'State complexity and optimizations'];
          break;
        case 'ds_algo':
          baseTips = ['Analyze complexity before coding', 'Design edge-case tests', 'Justify chosen data structures'];
          break;
        case 'system_design':
          baseTips = ['Start from requirements', 'Draw high-level architecture (API, storage, cache, queue)', 'Discuss bottlenecks and scaling plan'];
          break;
        case 'leadership':
          baseTips = ['Use real scenarios', 'Clarify role, decisions, outcomes', 'Reflect and extract lessons'];
          break;
        default:
          baseTips = [];
      }
    }
    return [...baseTips, ...techTips];
  };

  const proceedToNextRoundOrFinish = () => {
    const nextIndex = currentRoundIndex + 1;
    if (nextIndex < ROUNDS.length) {
      const nextRoundId = ROUNDS[nextIndex].id;
      setRound(nextRoundId);
      setCurrentRoundIndex(nextIndex);
      setQuestion(null);
      setAnswer('');
      setAttempt(null);
      const { questionCount } = (() => {
        const sr = ROUNDS[nextIndex];
        return { questionCount: sr.count };
      })();
      setTotalQuestions(questionCount);
      setCurrentQuestionIndex(1);
      setIsSessionComplete(false);
      startSession();
    } else {
      const details = ROUNDS.map(r => {
        const stat = roundStats[r.id] || { sum: 0, count: 0 };
        const avg = stat.count ? Math.round((stat.sum / stat.count) * 10) / 10 : 0;
        const threshold = getThresholdForRound(r.id);
        const pass = avg >= threshold;
        return { id: r.id, avg, threshold, pass };
      });
      const overallPass = details.every(d => d.pass);
      setFinalSummary({ overallPass, details });
      setIsSessionComplete(true);
    }
  };

  const handleNextQuestion = () => {
    if (currentQuestionIndex >= totalQuestions) {
      proceedToNextRoundOrFinish();
      return;
    }
    if (nextQuestionId && nextQuestionId !== '00000000-0000-0000-0000-000000000000') {
      setCurrentQuestionIndex(prev => prev + 1);
      fetchQuestion(nextQuestionId);
    } else {
      proceedToNextRoundOrFinish();
    }
  };

  const handleFinishSession = () => {
    setStep('setup');
    setSession(null);
    setQuestion(null);
    setAnswer('');
    setAttempt(null);
    setCurrentQuestionIndex(0);
    setTotalQuestions(0);
    setIsSessionComplete(false);
    setNextQuestionId(null);
  };

  const toggleListening = () => {
    if (isListening) {
      if (recognition) {
        recognition.stop();
      }
      setIsListening(false);
    } else {
      const SpeechRecognition = (window as any).SpeechRecognition || (window as any).webkitSpeechRecognition;
      if (!SpeechRecognition) {
        alert("Browser does not support speech recognition.");
        return;
      }

      const newRecognition = new SpeechRecognition();
      newRecognition.lang = language === 'vi' ? 'vi-VN' : 'en-US';
      newRecognition.continuous = false; // Stop after one sentence/pause
      newRecognition.interimResults = true;

      newRecognition.onstart = () => {
        setIsListening(true);
      };

      newRecognition.onresult = (event: any) => {
        let finalTranscript = '';
        for (let i = event.resultIndex; i < event.results.length; ++i) {
          if (event.results[i].isFinal) {
            finalTranscript += event.results[i][0].transcript;
          }
        }
        if (finalTranscript) {
          setAnswer(prev => {
             // Add space if there is already text and it doesn't end with space
             const prefix = prev && !prev.endsWith(' ') ? ' ' : '';
             return prev + prefix + finalTranscript;
          });
        }
      };

      newRecognition.onerror = (event: any) => {
        console.error("Speech recognition error", event.error);
        setIsListening(false);
      };

      newRecognition.onend = () => {
        setIsListening(false);
      };

      newRecognition.start();
      setRecognition(newRecognition);
    }
  };

  const getRoundDescription = () => {
    switch (round) {
      case 'recruiter': return language === 'vi' ? 'Phù hợp với mọi vị trí, tập trung vào kỹ năng mềm và văn hóa.' : 'Suitable for all roles, focusing on soft skills and culture.';
      case 'technical': return language === 'vi' ? `Các câu hỏi chuyên sâu về ${stack}.` : `In-depth questions about ${stack}.`;
      case 'ds_algo': return language === 'vi' ? 'Kiểm tra tư duy thuật toán và giải quyết vấn đề (LeetCode style).' : 'Tests algorithmic thinking and problem solving (LeetCode style).';
      case 'system_design': return language === 'vi' ? 'Thiết kế hệ thống quy mô lớn (Scalability, Reliability).' : 'Designing large-scale systems (Scalability, Reliability).';
      case 'leadership': return language === 'vi' ? 'Kỹ năng lãnh đạo, quản lý và xử lý tình huống.' : 'Leadership, management, and situational handling skills.';
      default: return '';
    }
  };

  // Render Setup Step
  if (step === 'setup') {
    return (
      <div className="max-w-2xl mx-auto p-6 bg-gray-50 min-h-screen font-sans flex flex-col justify-center">
        <header className="mb-8 text-center">
          <h1 className="text-3xl font-bold text-gray-800 flex items-center justify-center gap-2">
            <BookOpen className="w-8 h-8 text-blue-600" />
            {t.title}
          </h1>
          <p className="text-gray-600 mt-2">{t.subtitle}</p>
        </header>

        <div className="bg-white rounded-xl shadow-sm border border-gray-200 p-8">
          <h2 className="text-xl font-semibold text-gray-800 mb-6 flex items-center gap-2">
            <Settings className="w-5 h-5 text-gray-500" />
            {t.setup_title}
          </h2>

          <div className="space-y-6">
            {/* Language Selection */}
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-2 flex items-center gap-2">
                <Globe className="w-4 h-4" />
                {t.language_label}
              </label>
              <div className="flex gap-4">
                <button
                  onClick={() => setLanguage('en')}
                  className={`flex-1 py-2 px-4 rounded-lg border ${
                    language === 'en'
                      ? 'bg-blue-50 border-blue-500 text-blue-700'
                      : 'border-gray-200 text-gray-600 hover:bg-gray-50'
                  }`}
                >
                  English
                </button>
                <button
                  onClick={() => setLanguage('vi')}
                  className={`flex-1 py-2 px-4 rounded-lg border ${
                    language === 'vi'
                      ? 'bg-blue-50 border-blue-500 text-blue-700'
                      : 'border-gray-200 text-gray-600 hover:bg-gray-50'
                  }`}
                >
                  Tiếng Việt
                </button>
              </div>
            </div>

            {/* Stack Selection */}
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-2 flex items-center gap-2">
                <Code className="w-4 h-4" />
                {t.stack_label}
              </label>
              <select
                value={stack}
                onChange={(e) => setStack(e.target.value)}
                className="w-full p-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
              >
                {STACKS.map((s) => (
                  <option key={s} value={s}>{s}</option>
                ))}
              </select>
            </div>

            {/* Round Selection */}
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-2 flex items-center gap-2">
                <Briefcase className="w-4 h-4" />
                {t.round_label}
              </label>
              <div className="space-y-2">
                {ROUNDS.map((r) => (
                  <label
                    key={r.id}
                    className={`flex items-center p-3 rounded-lg border cursor-pointer transition ${
                      round === r.id
                        ? 'bg-blue-50 border-blue-500 ring-1 ring-blue-500'
                        : 'border-gray-200 hover:bg-gray-50'
                    }`}
                  >
                    <input
                      type="radio"
                      name="round"
                      value={r.id}
                      checked={round === r.id}
                      onChange={(e) => setRound(e.target.value)}
                      className="w-4 h-4 text-blue-600 border-gray-300 focus:ring-blue-500"
                    />
                    <span className="ml-3 text-gray-700">
                      {t.rounds[r.labelKey as keyof typeof t.rounds]}
                    </span>
                  </label>
                ))}
              </div>
            </div>

            {/* Threshold Configuration */}
            <div className="mt-4 pt-4 border-t border-gray-100">
               <button 
                  className="flex items-center gap-2 text-sm text-blue-600 font-medium mb-2 hover:underline"
                  onClick={() => {
                     const el = document.getElementById('threshold-config');
                     if (el) el.classList.toggle('hidden');
                  }}
               >
                  <Settings className="w-4 h-4" />
                  {language === 'vi' ? 'Cấu hình nâng cao (Ngưỡng đậu)' : 'Advanced Configuration (Pass Thresholds)'}
               </button>
               <div id="threshold-config" className="hidden space-y-3 bg-gray-50 p-3 rounded-lg border border-gray-200">
                  {ROUNDS.map(r => (
                    <div key={r.id} className="flex items-center justify-between">
                       <label className="text-sm text-gray-700">{t.rounds[r.labelKey as keyof typeof t.rounds]}</label>
                       <input 
                          type="number" 
                          min="0" max="100"
                          value={getThresholdForRound(r.id)}
                          onChange={(e) => setCustomThresholds(prev => ({...prev, [r.id]: parseInt(e.target.value) || 0}))}
                          className="w-20 p-1 text-sm border border-gray-300 rounded text-center"
                       />
                    </div>
                  ))}
               </div>
            </div>

            {/* Description Box */}
            <div className="bg-blue-50 p-4 rounded-lg text-blue-800 text-sm border border-blue-100">
              <strong>{t.round_label}: </strong> {getRoundDescription()}
            </div>

            {/* AI Enable Toggle */}
            <div className="flex items-center gap-3 py-2 border-t border-gray-100 mt-4 pt-4">
                <label className="relative inline-flex items-center cursor-pointer">
                    <input 
                        type="checkbox" 
                        checked={aiEnabled}
                        onChange={(e) => setAiEnabled(e.target.checked)}
                        className="sr-only peer"
                    />
                    <div className="w-11 h-6 bg-gray-200 peer-focus:outline-none peer-focus:ring-4 peer-focus:ring-blue-300 rounded-full peer peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-blue-600"></div>
                    <span className="ml-3 text-sm font-medium text-gray-900 flex items-center gap-2">
                        <Settings className="w-4 h-4 text-gray-500" />
                        {t.ai_enabled_label}
                    </span>
                </label>
            </div>

            {/* Voice Input Toggle */}
            <div className="flex items-center gap-3 py-2">
                <label className="relative inline-flex items-center cursor-pointer">
                    <input 
                        type="checkbox" 
                        checked={enableVoice}
                        onChange={(e) => setEnableVoice(e.target.checked)}
                        className="sr-only peer"
                    />
                    <div className="w-11 h-6 bg-gray-200 peer-focus:outline-none peer-focus:ring-4 peer-focus:ring-blue-300 rounded-full peer peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-blue-600"></div>
                    <span className="ml-3 text-sm font-medium text-gray-900 flex items-center gap-2">
                        <Mic className="w-4 h-4 text-gray-500" />
                        {t.enable_voice}
                    </span>
                </label>
            </div>

            <button
              onClick={startSession}
              disabled={loading}
              className="w-full py-3 bg-blue-600 text-white rounded-lg font-medium hover:bg-blue-700 transition disabled:opacity-50 flex justify-center items-center gap-2"
            >
              {loading ? (
                <span>{t.starting}</span>
              ) : (
                <>
                  <Play className="w-5 h-5" />
                  {t.start_session}
                </>
              )}
            </button>

            {error && (
              <div className="p-3 bg-red-50 text-red-700 rounded-lg text-sm flex items-center gap-2">
                <AlertCircle className="w-4 h-4" />
                {error}
              </div>
            )}
          </div>
        </div>
      </div>
    );
  }

  // Render Session Step
  return (
    <div className="max-w-4xl mx-auto p-6 bg-gray-50 min-h-screen font-sans">
      <header className="mb-6 text-center">
        <h1 className="text-2xl font-bold text-gray-800 flex items-center justify-center gap-2">
          <BookOpen className="w-6 h-6 text-blue-600" />
          {t.title}
        </h1>
        <p className="text-gray-600 mt-2">{t.subtitle}</p>
        <div className="mt-4 flex items-center justify-center gap-2 text-sm text-gray-500 font-medium">
             <span>{stack}</span>
             <span>•</span>
             <span>{getRoundDescription()}</span>
             <span>•</span>
             <span className="bg-blue-100 text-blue-800 px-2 py-0.5 rounded-full">
               {t.question_label} {currentQuestionIndex} / {totalQuestions}
             </span>
        </div>
      </header>

      {error && (
        <div className="mb-6 p-4 bg-red-100 border border-red-300 text-red-700 rounded-lg flex items-center gap-2">
          <AlertCircle className="w-5 h-5" />
          {error}
        </div>
      )}

      {isSessionComplete ? (
         <div className="max-w-2xl mx-auto bg-white rounded-xl shadow-sm border border-gray-200 p-8" id="final-summary">
             <div className="flex items-center gap-3 mb-6">
               <Briefcase className={`w-6 h-6 ${finalSummary?.overallPass ? 'text-green-600' : 'text-red-600'}`} />
               <h2 className="text-2xl font-bold text-gray-800">
                 {finalSummary?.overallPass ? t.overall_passed : t.overall_failed}
               </h2>
             </div>
             <div className="space-y-3">
               {finalSummary?.details.map(d => (
                 <div key={d.id} className="flex items-center justify-between p-3 rounded-lg border">
                   <div className="text-gray-700 font-medium">
                     {t.rounds[d.id as keyof typeof t.rounds]}
                   </div>
                   <div className="flex items-center gap-3">
                     <span className="text-sm text-gray-600">
                       {t.avg_label}: <strong>{d.avg}</strong> / {t.threshold_label}: {d.threshold}
                     </span>
                     <span className={`px-2 py-0.5 rounded text-sm ${d.pass ? 'bg-green-100 text-green-700' : 'bg-red-100 text-red-700'}`}>
                       {d.pass ? t.pass_label : t.fail_label}
                     </span>
                   </div>
                 </div>
               ))}
             </div>
             <div className="mt-8 flex justify-center gap-3">
               <button
                 onClick={() => {
                   const lines: string[] = [];
                   lines.push(`# ${finalSummary?.overallPass ? t.overall_passed : t.overall_failed}`);
                   lines.push(`Date: ${new Date().toLocaleString()}`);
                   lines.push(`Tech Stack: ${stack}`);
                   lines.push('');
                   
                   lines.push('## Summary Details');
                   finalSummary?.details.forEach(d => {
                     const name = t.rounds[d.id as keyof typeof t.rounds];
                     lines.push(`- ${name}: ${t.avg_label} ${d.avg} / ${t.threshold_label} ${d.threshold} — ${d.pass ? t.pass_label : t.fail_label}`);
                   });
                   lines.push('');

                   lines.push('## Improvement Tips');
                   ROUNDS.forEach(r => {
                      const tips = getImprovementTips(r.id);
                      if (tips.length > 0) {
                          lines.push(`### ${t.rounds[r.labelKey as keyof typeof t.rounds]}`);
                          tips.forEach(tip => lines.push(`- ${tip}`));
                          lines.push('');
                      }
                   });

                   lines.push('## Attempts History');
                   attemptsHistory.forEach((a, i) => {
                     const name = t.rounds[a.roundId as keyof typeof t.rounds];
                     lines.push(`### Q${i+1} (${name}) — Score: ${a.score}`);
                     lines.push(`**Feedback:**`);
                     lines.push(`${a.feedback}`);
                     lines.push('');
                   });
                   
                   const blob = new Blob([lines.join('\n')], { type: 'text/markdown;charset=utf-8' });
                   const url = URL.createObjectURL(blob);
                   const aEl = document.createElement('a');
                   aEl.href = url;
                   aEl.download = `interview-summary-${new Date().toISOString().split('T')[0]}.md`;
                   document.body.appendChild(aEl);
                   aEl.click();
                   document.body.removeChild(aEl);
                   URL.revokeObjectURL(url);
                 }}
                 className="px-6 py-2 bg-gray-100 text-gray-800 rounded-lg hover:bg-gray-200 transition border"
               >
                 {t.download_summary}
               </button>
               <button
                 onClick={handleFinishSession}
                 className="px-6 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition"
               >
                 {t.return_setup}
               </button>
             </div>
         </div>
      ) : (
      <div className="space-y-6">
        {/* Question Card */}
        {question && (
          <div className="bg-white rounded-xl shadow-sm border border-gray-200 p-6">
            <div className="flex justify-between items-start mb-4">
              <span className="px-3 py-1 bg-blue-100 text-blue-800 text-xs font-semibold rounded-full uppercase tracking-wide">
                {question.topic}
              </span>
              <span className="px-3 py-1 bg-gray-100 text-gray-600 text-xs font-semibold rounded-full uppercase tracking-wide">
                {question.level}
              </span>
            </div>
            <h3 className="text-xl font-medium text-gray-900 mb-2">{t.question_label}:</h3>
            <p className="text-gray-800 text-lg leading-relaxed whitespace-pre-wrap">{question.content}</p>
          </div>
        )}

        {/* Answer Section */}
        <div className="bg-white rounded-xl shadow-sm border border-gray-200 p-6">
          <label className="block text-sm font-medium text-gray-700 mb-2">{t.your_answer}</label>
          <div className="relative">
            <textarea
              value={answer}
              onChange={(e) => setAnswer(e.target.value)}
              disabled={loading || !!attempt}
              className="w-full h-40 p-4 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500 resize-none transition"
              placeholder={t.placeholder_answer}
            />
            {enableVoice && !attempt && (
                <button
                    onClick={toggleListening}
                    className={`absolute right-4 bottom-4 p-2 rounded-full transition-colors ${
                        isListening 
                        ? 'bg-red-100 text-red-600 hover:bg-red-200 animate-pulse' 
                        : 'bg-gray-100 text-gray-600 hover:bg-gray-200'
                    }`}
                    title={isListening ? "Stop listening" : t.click_to_speak}
                >
                    {isListening ? <MicOff className="w-5 h-5" /> : <Mic className="w-5 h-5" />}
                </button>
            )}
          </div>
          {!attempt && (
            <div className="mt-4 flex justify-end">
              <button
                onClick={submitAnswer}
                disabled={loading || !answer.trim()}
                className="flex items-center gap-2 px-6 py-2.5 bg-green-600 text-white rounded-lg font-medium hover:bg-green-700 transition disabled:opacity-50"
              >
                {loading ? (
                  t.evaluating
                ) : (
                  <>
                    <Send className="w-4 h-4" />
                    {t.submit}
                  </>
                )}
              </button>
            </div>
          )}
        </div>

        {/* Feedback Section */}
        {attempt && (
          <div className="bg-white rounded-xl shadow-sm border border-gray-200 overflow-hidden">
            <div className={`p-6 border-b ${attempt.score >= 70 ? 'bg-green-50 border-green-100' : 'bg-orange-50 border-orange-100'}`}>
              <div className="flex items-center gap-4">
                <div className={`flex items-center justify-center w-16 h-16 rounded-full border-4 text-2xl font-bold ${
                  attempt.score >= 70 ? 'border-green-500 text-green-700 bg-white' : 'border-orange-500 text-orange-700 bg-white'
                }`}>
                  {attempt.score}
                </div>
                <div>
                  <h3 className="text-lg font-bold text-gray-900">{t.ai_evaluation}</h3>
                  <p className="text-gray-600">
                    {attempt.score >= 70 ? t.good_job : t.needs_improvement}
                  </p>
                </div>
              </div>
            </div>
            <div className="p-6 bg-gray-50">
              <h4 className="text-sm font-semibold text-gray-500 uppercase tracking-wider mb-3">{t.feedback_detail}</h4>
              <div className="prose prose-blue max-w-none text-gray-800 whitespace-pre-wrap">
                {attempt.feedback}
              </div>
              <div className="mt-4">
                <h4 className="text-sm font-semibold text-gray-700 mb-2">{t.suggestions_label}</h4>
                <ul className="list-disc ml-5 text-gray-800">
                  {getImprovementTips(round).map((tip, idx) => (
                    <li key={idx}>{tip}</li>
                  ))}
                </ul>
              </div>
            </div>
            <div className="p-4 bg-gray-100 border-t border-gray-200 flex justify-center">
              <button
                onClick={handleNextQuestion}
                className="flex items-center gap-2 px-6 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition"
              >
                {t.next_question}
              </button>
            </div>
          </div>
        )}
      </div>
      )}
    </div>
  );
};

export default PracticeSession;
