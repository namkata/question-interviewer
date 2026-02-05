import React, { useState, useRef, useEffect } from 'react';
import { 
  Mic, MicOff, Send, Play, RefreshCw, 
  CheckCircle, AlertCircle, ChevronRight, BookOpen,
  Code, Database, Layout, Server, Terminal,
  Cpu, Shield, MessageSquare, Lightbulb, XCircle
} from 'lucide-react';
import axios from 'axios';

// --- Constants & Types ---

const API_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080';

type Language = 'vi' | 'en';
type Level = 'Fresher' | 'Junior' | 'Mid' | 'Senior';
type Role = 'BackEnd' | 'FrontEnd' | 'FullStack' | 'DevOps' | 'Data Engineer';

interface Round {
  id: string;
  title: { en: string; vi: string };
  description: { en: string; vi: string };
  icon: React.ReactNode;
  count: number;
}

const ROUNDS: Round[] = [
  {
    id: 'cv_screening',
    title: { en: 'CV Screening', vi: 'Sàng lọc CV' },
    description: { en: 'Background & Fit', vi: 'Kiểm tra nền tảng & mức độ phù hợp' },
    icon: <BookOpen className="w-5 h-5" />,
    count: 2
  },
  {
    id: 'core_tech',
    title: { en: 'Language & Core', vi: 'Ngôn ngữ & Cốt lõi' },
    description: { en: 'Deep dive into main language', vi: 'Đào sâu ngôn ngữ chính' },
    icon: <Code className="w-5 h-5" />,
    count: 3 // Reduced for demo flow, normally 5-7
  },
  {
    id: 'database',
    title: { en: 'Database', vi: 'Cơ sở dữ liệu' },
    description: { en: 'Data modeling & SQL', vi: 'Mô hình dữ liệu & SQL' },
    icon: <Database className="w-5 h-5" />,
    count: 2
  },
  {
    id: 'system_design',
    title: { en: 'System Design', vi: 'Thiết kế hệ thống' },
    description: { en: 'Architecture & Scalability', vi: 'Kiến trúc & Khả năng mở rộng' },
    icon: <Layout className="w-5 h-5" />,
    count: 1
  },
  {
    id: 'coding',
    title: { en: 'Coding', vi: 'Lập trình' },
    description: { en: 'Algorithms & Data Structures', vi: 'Thuật toán & Cấu trúc dữ liệu' },
    icon: <Terminal className="w-5 h-5" />,
    count: 1
  },
  {
    id: 'testing',
    title: { en: 'Testing', vi: 'Kiểm thử' },
    description: { en: 'Quality Assurance', vi: 'Đảm bảo chất lượng' },
    icon: <CheckCircle className="w-5 h-5" />,
    count: 1
  },
  {
    id: 'devops_round',
    title: { en: 'DevOps', vi: 'DevOps' },
    description: { en: 'Deployment & Infrastructure', vi: 'Triển khai & Hạ tầng' },
    icon: <Server className="w-5 h-5" />,
    count: 1
  },
  {
    id: 'behavioral',
    title: { en: 'Behavioral', vi: 'Hành vi' },
    description: { en: 'Culture Fit & Soft Skills', vi: 'Văn hóa & Kỹ năng mềm' },
    icon: <MessageSquare className="w-5 h-5" />,
    count: 1
  }
];

const TRANSLATIONS = {
  en: {
    title: 'AI Mock Interview',
    subtitle: 'Practice with AI across 8 standard interview rounds',
    start_btn: 'Start Interview',
    role_label: 'Target Role',
    level_label: 'Seniority Level',
    tech_stack_label: 'Tech Stack',
    language_label: 'Interview Language',
    enable_voice: 'Voice Mode',
    enable_hints: 'Answer Hints',
    listening: 'Listening...',
    processing: 'Thinking...',
    your_answer: 'Your Answer',
    submit: 'Submit Answer',
    next: 'Next Question',
    finish: 'Finish Session',
    feedback: 'AI Feedback',
    score: 'Score',
    suggested_answer: 'Suggested Answer',
    error_start: 'Failed to start session',
    error_answer: 'Failed to submit answer',
    round_complete: 'Round Complete!',
    session_complete: 'Interview Session Completed',
    summary_title: 'Performance Summary',
    pass: 'PASS',
    fail: 'NEEDS IMPROVEMENT',
    recording_error: 'Microphone access denied or not supported',
    hint_title: 'Hint',
    restart: 'Start New Session'
  },
  vi: {
    title: 'Phỏng vấn thử AI',
    subtitle: 'Luyện tập với AI qua 8 vòng phỏng vấn tiêu chuẩn',
    start_btn: 'Bắt đầu phỏng vấn',
    role_label: 'Vị trí ứng tuyển',
    level_label: 'Cấp độ',
    tech_stack_label: 'Tech Stack',
    language_label: 'Ngôn ngữ phỏng vấn',
    enable_voice: 'Chế độ giọng nói',
    enable_hints: 'Gợi ý câu trả lời',
    listening: 'Đang nghe...',
    processing: 'Đang suy nghĩ...',
    your_answer: 'Câu trả lời của bạn',
    submit: 'Gửi câu trả lời',
    next: 'Câu tiếp theo',
    finish: 'Kết thúc phiên',
    feedback: 'Đánh giá của AI',
    score: 'Điểm',
    suggested_answer: 'Gợi ý trả lời',
    error_start: 'Không thể bắt đầu phiên',
    error_answer: 'Không thể gửi câu trả lời',
    round_complete: 'Hoàn thành vòng!',
    session_complete: 'Đã hoàn thành buổi phỏng vấn',
    summary_title: 'Tổng kết hiệu suất',
    pass: 'ĐẠT',
    fail: 'CẦN CẢI THIỆN',
    recording_error: 'Không thể truy cập microphone',
    hint_title: 'Gợi ý',
    restart: 'Bắt đầu phiên mới'
  }
};

const ROLES: { id: Role; label: string }[] = [
  { id: 'BackEnd', label: 'Backend Developer' },
  { id: 'FrontEnd', label: 'Frontend Developer' },
  { id: 'FullStack', label: 'Fullstack Developer' },
  { id: 'DevOps', label: 'DevOps Engineer' },
  { id: 'Data Engineer', label: 'Data Engineer' }
];

const LEVELS: { id: Level; label: string }[] = [
  { id: 'Fresher', label: 'Fresher' },
  { id: 'Junior', label: 'Junior' },
  { id: 'Mid', label: 'Mid-Level' },
  { id: 'Senior', label: 'Senior' }
];

const STACKS = [
  'JavaScript', 'TypeScript', 'Python', 'Java', 'Go', 'C#', 
  'React', 'Vue', 'Angular', 'Node.js', 'Django', 'Spring Boot',
  'PostgreSQL', 'MongoDB', 'Redis', 'Docker', 'Kubernetes', 'AWS'
];

export default function PracticeSession() {
  // --- State ---
  const [step, setStep] = useState<'setup' | 'session' | 'summary'>('setup');
  const [language, setLanguage] = useState<Language>('vi');
  
  // Setup State
  const [selectedRole, setSelectedRole] = useState<Role>('BackEnd');
  const [selectedLevel, setSelectedLevel] = useState<Level>('Junior');
  const [selectedStacks, setSelectedStacks] = useState<string[]>(['Go', 'PostgreSQL']);
  const [enableVoice, setEnableVoice] = useState(false);
  const [enableHints, setEnableHints] = useState(true);

  // Session State
  const [session, setSession] = useState<any>(null);
  const [currentRoundIndex, setCurrentRoundIndex] = useState(0);
  const [round, setRound] = useState<string>(ROUNDS[0].id);
  
  const [question, setQuestion] = useState<any>(null);
  const [answer, setAnswer] = useState('');
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');
  
  const [attempt, setAttempt] = useState<any>(null);
  const [nextQuestionId, setNextQuestionId] = useState<string | null>(null);
  const [currentQuestionIndex, setCurrentQuestionIndex] = useState(1);
  const [totalQuestions, setTotalQuestions] = useState(2);
  
  // Voice State
  const [isListening, setIsListening] = useState(false);
  const [voiceSupported, setVoiceSupported] = useState(true);
  const recognitionRef = useRef<any>(null);

  useEffect(() => {
    if (!('webkitSpeechRecognition' in window)) {
      setVoiceSupported(false);
    }
  }, []);

  // Auto-start listening when question changes if voice mode is enabled
  useEffect(() => {
    if (enableVoice && question && !attempt && !loading) {
        // Small delay to allow TTS to start or user to get ready
        const timer = setTimeout(() => {
            startListening();
        }, 1000);
        return () => clearTimeout(timer);
    }
  }, [enableVoice, question, attempt, loading]);

  // Stats
  const [roundStats, setRoundStats] = useState<Record<string, { sum: number; count: number }>>({});
  const [finalSummary, setFinalSummary] = useState<any>(null);
  const [attemptsHistory, setAttemptsHistory] = useState<any[]>([]);

  const t = TRANSLATIONS[language];

  // --- Helpers ---
  
  const toggleStack = (stack: string) => {
    setSelectedStacks(prev => 
      prev.includes(stack) 
        ? prev.filter(s => s !== stack)
        : [...prev, stack]
    );
  };

  const speak = (text: string) => {
    if (!enableVoice) return;
    window.speechSynthesis.cancel();
    const utterance = new SpeechSynthesisUtterance(text);
    utterance.lang = language === 'vi' ? 'vi-VN' : 'en-US';
    window.speechSynthesis.speak(utterance);
  };

  const startListening = () => {
    if (!('webkitSpeechRecognition' in window)) {
      setError(t.recording_error);
      return;
    }
    const SpeechRecognition = (window as any).webkitSpeechRecognition;
    const recognition = new SpeechRecognition();
    recognition.lang = language === 'vi' ? 'vi-VN' : 'en-US';
    recognition.continuous = false;
    recognition.interimResults = false;

    recognition.onstart = () => setIsListening(true);
    recognition.onend = () => setIsListening(false);
    recognition.onerror = (event: any) => {
      console.error('Speech recognition error', event.error);
      setIsListening(false);
    };
    recognition.onresult = (event: any) => {
      const transcript = event.results[0][0].transcript;
      setAnswer(prev => prev + (prev ? ' ' : '') + transcript);
    };

    recognitionRef.current = recognition;
    recognition.start();
  };

  const stopListening = () => {
    if (recognitionRef.current) {
      recognitionRef.current.stop();
    }
  };

  // --- API Calls ---

  const getApiUrl = (endpoint: string) => {
    // Ensure we construct the full path correctly
    const baseUrl = API_URL.endsWith('/') ? API_URL.slice(0, -1) : API_URL;
    return `${baseUrl}/api/v1/practice${endpoint}`;
  };

  const fetchQuestion = async (questionId: string) => {
    try {
      setLoading(true);
      const res = await axios.get(getApiUrl(`/questions/${questionId}`));
      setQuestion(res.data);
      speak(res.data.content);
      setAnswer('');
      setAttempt(null);
    } catch (err) {
      console.error(err);
      setError('Could not load question');
    } finally {
      setLoading(false);
    }
  };

  const startSession = async (roundIdOverride?: string) => {
    setLoading(true);
    setError('');
    setNextQuestionId(null);
    setCurrentQuestionIndex(1);
    
    // Determine the round to use: override -> state -> default to first round
    let activeRoundId = roundIdOverride || round;
    if (!activeRoundId && step === 'setup') {
        activeRoundId = ROUNDS[0].id;
        setRound(activeRoundId);
    }

    // If starting a fresh session (setup step), reset summary and history
    if (step === 'setup') {
        setFinalSummary(null);
        setAttemptsHistory([]);
        setRoundStats({});
    }

    const idx = ROUNDS.findIndex(r => r.id === activeRoundId);
    setCurrentRoundIndex(idx >= 0 ? idx : 0);
    
    // Logic to map round to topic/level
    const getTopicAndLevelForRound = (rId: string) => {
        let topicName: string | undefined = undefined;
        let levelName = selectedLevel;
        let questionCount = 2;

        const selectedRound = ROUNDS.find(r => r.id === rId);
        if (selectedRound) {
            questionCount = selectedRound.count;
        }

        switch (rId) {
            case 'cv_screening': topicName = 'CV Screening'; break;
            case 'core_tech': topicName = undefined; break; // Will use role/stack
            case 'database': topicName = 'Database'; break;
            case 'system_design': topicName = 'System Design'; break;
            case 'coding': topicName = 'Algorithms'; break;
            case 'testing': topicName = 'Testing'; break;
            case 'devops_round': topicName = 'DevOps'; break;
            case 'behavioral': topicName = 'Behavioral'; break;
            default: topicName = undefined;
        }
        return { topicName, levelName, questionCount };
    };

    const { topicName, levelName, questionCount } = getTopicAndLevelForRound(activeRoundId);
    setTotalQuestions(questionCount);

    try {
      // Use a fixed user ID for demo purposes
      const userId = '123e4567-e89b-12d3-a456-426614174000';
      const response = await axios.post(getApiUrl('/sessions'), {
        user_id: userId,
        topic_id: topicName, // Backend handles string topic lookup if needed or ID
        level: levelName,
        language: language,
        config: {
          role: selectedRole,
          stacks: selectedStacks,
          round_id: activeRoundId
        }
      });

      const { session, first_question_id } = response.data;
      setSession(session);
      setStep('session');
      if (first_question_id) {
          await fetchQuestion(first_question_id);
      } else {
          // Fallback if no question returned
           setError('No questions found for this criteria.');
      }
    } catch (err: any) {
      console.error(err);
      setError(err.response?.data?.error || t.error_start);
    } finally {
      setLoading(false);
    }
  };

  const submitAnswer = async () => {
    if (!answer.trim() || !session || !question) return;
    setLoading(true);
    try {
      // Use the correct endpoint /answers and payload key 'content'
      const res = await axios.post(getApiUrl(`/sessions/${session.id}/answers`), {
        question_id: question.id,
        content: answer,
        language: language,
        ai_enabled: true 
      });
      
      const result = res.data.attempt || res.data; // Handle wrapped response if any
      setAttempt(result);
      setNextQuestionId(res.data.next_question_id);
      
      // Update stats
      setRoundStats(prev => {
        const current = prev[round] || { sum: 0, count: 0 };
        return {
            ...prev,
            [round]: {
                sum: current.sum + (result.score || 0),
                count: current.count + 1
            }
        };
      });
      
      setAttemptsHistory(prev => [...prev, {
          roundId: round,
          question: question.content,
          answer: answer,
          score: result.score,
          feedback: result.feedback
      }]);

      if (enableVoice && result.feedback) {
        speak(result.feedback.split('.')[0]); // Speak first sentence of feedback
      }

    } catch (err: any) {
      console.error(err);
      setError(t.error_answer);
    } finally {
      setLoading(false);
    }
  };

  const proceedToNextQuestion = () => {
      if (nextQuestionId && currentQuestionIndex < totalQuestions) {
          setCurrentQuestionIndex(prev => prev + 1);
          fetchQuestion(nextQuestionId);
      } else {
          // End of round
          proceedToNextRoundOrFinish();
      }
  };

  const getThresholdForRound = (roundId: string) => {
      // Example thresholds
      if (roundId === 'core_tech' || roundId === 'system_design') return 7.0;
      return 6.0;
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
      // Auto start next round
      startSession(nextRoundId);
    } else {
      // All rounds complete
      const details = ROUNDS.map(r => {
        const stat = roundStats[r.id] || { sum: 0, count: 0 };
        const avg = stat.count ? Math.round((stat.sum / stat.count) * 10) / 10 : 0;
        const threshold = getThresholdForRound(r.id);
        const pass = avg >= threshold;
        return { id: r.id, title: r.title[language], avg, threshold, pass };
      });
      const overallPass = details.every(d => d.pass);
      setFinalSummary({ overallPass, details });
      setStep('summary');
    }
  };

  // --- Render ---

  if (step === 'setup') {
    return (
      <div className="max-w-4xl mx-auto p-6">
        <div className="bg-white rounded-xl shadow-lg p-8">
          <div className="flex items-center gap-3 mb-6">
            <div className="p-3 bg-blue-100 rounded-lg">
              <Cpu className="w-8 h-8 text-blue-600" />
            </div>
            <div>
              <h1 className="text-2xl font-bold text-gray-900">{t.title}</h1>
              <p className="text-gray-500">{t.subtitle}</p>
            </div>
          </div>

          <div className="grid grid-cols-1 md:grid-cols-2 gap-8">
            <div className="space-y-6">
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-2">{t.language_label}</label>
                <div className="flex gap-2">
                  <button 
                    onClick={() => setLanguage('vi')}
                    className={`px-4 py-2 rounded-lg border ${language === 'vi' ? 'bg-blue-50 border-blue-200 text-blue-700' : 'bg-white border-gray-200 text-gray-700'}`}
                  >
                    Tiếng Việt
                  </button>
                  <button 
                    onClick={() => setLanguage('en')}
                    className={`px-4 py-2 rounded-lg border ${language === 'en' ? 'bg-blue-50 border-blue-200 text-blue-700' : 'bg-white border-gray-200 text-gray-700'}`}
                  >
                    English
                  </button>
                </div>
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-700 mb-2">{t.role_label}</label>
                <select 
                  value={selectedRole}
                  onChange={(e) => setSelectedRole(e.target.value as Role)}
                  className="w-full p-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500"
                >
                  {ROLES.map(r => (
                    <option key={r.id} value={r.id}>{r.label}</option>
                  ))}
                </select>
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-700 mb-2">{t.level_label}</label>
                <select 
                  value={selectedLevel}
                  onChange={(e) => setSelectedLevel(e.target.value as Level)}
                  className="w-full p-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500"
                >
                  {LEVELS.map(l => (
                    <option key={l.id} value={l.id}>{l.label}</option>
                  ))}
                </select>
              </div>

              {/* Toggles */}
              <div className="space-y-3 pt-2">
                <div className={`flex items-center justify-between p-2 rounded-lg border transition-all group ${!voiceSupported ? 'bg-gray-100 border-gray-200 opacity-60 cursor-not-allowed' : 'hover:bg-gray-50 border-transparent hover:border-gray-100'}`}>
                  <label htmlFor="voice-toggle" className={`text-sm font-medium flex items-center gap-2 flex-1 ${!voiceSupported ? 'cursor-not-allowed text-gray-500' : 'cursor-pointer text-gray-700'}`}>
                     {enableVoice ? <Mic className="w-4 h-4 text-blue-500" /> : <MicOff className="w-4 h-4 text-gray-400" />}
                    {t.enable_voice}
                    {!voiceSupported && <span className="text-xs text-red-500 ml-2">(Not supported in this browser)</span>}
                  </label>
                  <div className="relative inline-flex items-center">
                    <input 
                      id="voice-toggle"
                      type="checkbox" 
                      checked={enableVoice}
                      disabled={!voiceSupported}
                      onChange={(e) => setEnableVoice(e.target.checked)}
                      className="sr-only peer"
                    />
                    <label htmlFor="voice-toggle" className={`w-11 h-6 bg-gray-200 rounded-full peer peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all ${voiceSupported ? 'peer-checked:bg-blue-600 peer-focus:outline-none peer-focus:ring-4 peer-focus:ring-blue-100 cursor-pointer' : 'cursor-not-allowed'}`}></label>
                  </div>
                </div>

                <div className="flex items-center justify-between p-2 hover:bg-gray-50 rounded-lg border border-transparent hover:border-gray-100 transition-all group">
                  <label htmlFor="hints-toggle" className="text-sm font-medium text-gray-700 flex items-center gap-2 cursor-pointer select-none flex-1">
                     <Lightbulb className={`w-4 h-4 ${enableHints ? 'text-amber-500' : 'text-gray-400'}`} />
                    {t.enable_hints}
                  </label>
                  <div className="relative inline-flex items-center">
                    <input 
                      id="hints-toggle"
                      type="checkbox" 
                      checked={enableHints}
                      onChange={(e) => setEnableHints(e.target.checked)}
                      className="sr-only peer"
                    />
                    <label htmlFor="hints-toggle" className="w-11 h-6 bg-gray-200 peer-focus:outline-none peer-focus:ring-4 peer-focus:ring-amber-100 rounded-full peer peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-amber-500 cursor-pointer"></label>
                  </div>
                </div>
              </div>

            </div>

            <div>
              <label className="block text-sm font-medium text-gray-700 mb-2">{t.tech_stack_label}</label>
              <div className="flex flex-wrap gap-2">
                {STACKS.map(stack => (
                  <button
                    key={stack}
                    onClick={() => toggleStack(stack)}
                    className={`px-3 py-1 rounded-full text-sm border transition-colors ${
                      selectedStacks.includes(stack)
                        ? 'bg-blue-100 border-blue-300 text-blue-800'
                        : 'bg-white border-gray-200 text-gray-600 hover:border-blue-300'
                    }`}
                  >
                    {stack}
                  </button>
                ))}
              </div>
              
              <div className="mt-8 bg-gray-50 p-4 rounded-lg">
                <h3 className="font-medium text-gray-900 mb-1">Interview Roadmap (8 Rounds)</h3>
                <p className="text-xs text-gray-500 mb-3 italic">
                    {language === 'vi' ? 'Bạn sẽ trải qua lần lượt tất cả 8 vòng phỏng vấn.' : 'You will go through all 8 interview rounds sequentially.'}
                </p>
                <div className="space-y-2">
                  {ROUNDS.map((r, i) => (
                    <div key={r.id} className="flex items-center gap-2 text-sm text-gray-600">
                      <div className="w-6 h-6 rounded-full bg-gray-200 flex items-center justify-center text-xs font-bold">
                        {i + 1}
                      </div>
                      <span className="font-medium">{r.title[language]}</span>
                    </div>
                  ))}
                </div>
              </div>
            </div>
          </div>

          <div className="mt-8 flex justify-end">
            <button
              onClick={() => startSession()}
              disabled={loading}
              className="flex items-center gap-2 px-8 py-3 bg-blue-600 text-white rounded-lg hover:bg-blue-700 disabled:opacity-50 font-medium text-lg shadow-sm transition-all hover:shadow-md"
            >
              {loading ? (
                <RefreshCw className="w-5 h-5 animate-spin" />
              ) : (
                <Play className="w-5 h-5" />
              )}
              {t.start_btn}
            </button>
          </div>
          
          {error && (
            <div className="mt-4 p-4 bg-red-50 text-red-700 rounded-lg flex items-center gap-2">
              <AlertCircle className="w-5 h-5" />
              {error}
            </div>
          )}
        </div>
      </div>
    );
  }

  // Summary View
  if (step === 'summary') {
    return (
      <div className="max-w-4xl mx-auto p-6">
        <div className="bg-white rounded-xl shadow-lg p-8">
          <div className="text-center mb-8">
            <h2 className="text-3xl font-bold text-gray-900 mb-2">{t.session_complete}</h2>
            <div className={`inline-flex items-center gap-2 px-4 py-2 rounded-full font-bold text-lg ${
              finalSummary?.overallPass ? 'bg-green-100 text-green-700' : 'bg-amber-100 text-amber-700'
            }`}>
              {finalSummary?.overallPass ? <CheckCircle className="w-6 h-6" /> : <AlertCircle className="w-6 h-6" />}
              {finalSummary?.overallPass ? t.pass : t.fail}
            </div>
          </div>

          <div className="grid grid-cols-1 md:grid-cols-2 gap-6 mb-8">
             {finalSummary?.details.map((detail: any) => (
                <div key={detail.id} className="bg-gray-50 p-4 rounded-lg border border-gray-200">
                    <div className="flex justify-between items-center mb-2">
                        <h3 className="font-bold text-gray-800">{detail.title}</h3>
                        <span className={`text-sm font-bold ${detail.pass ? 'text-green-600' : 'text-red-600'}`}>
                            {detail.pass ? 'PASS' : 'FAIL'}
                        </span>
                    </div>
                    <div className="flex items-center gap-2">
                         <div className="flex-1 h-2 bg-gray-200 rounded-full overflow-hidden">
                             <div 
                                className={`h-full ${detail.pass ? 'bg-green-500' : 'bg-red-500'}`} 
                                style={{ width: `${Math.min(detail.avg * 10, 100)}%` }}
                             ></div>
                         </div>
                         <span className="text-sm text-gray-600 font-medium w-12 text-right">{detail.avg}/10</span>
                    </div>
                </div>
             ))}
          </div>

          <div className="mb-8">
              <h3 className="text-xl font-bold text-gray-900 mb-4">Detailed Feedback History</h3>
              <div className="space-y-4 max-h-96 overflow-y-auto pr-2">
                  {attemptsHistory.map((attempt, idx) => (
                      <div key={idx} className="bg-gray-50 p-4 rounded-lg border border-gray-100">
                          <div className="flex justify-between mb-2">
                              <span className="text-sm font-medium text-gray-500">
                                {ROUNDS.find(r => r.id === attempt.roundId)?.title[language]}
                              </span>
                              <span className={`text-sm font-bold ${
                                  attempt.score >= 7 ? 'text-green-600' : attempt.score >= 5 ? 'text-amber-600' : 'text-red-600'
                              }`}>
                                  Score: {attempt.score}/10
                              </span>
                          </div>
                          <p className="font-medium text-gray-900 mb-2">{attempt.question}</p>
                          <div className="pl-4 border-l-2 border-gray-200 mb-2">
                              <p className="text-gray-600 text-sm italic">{attempt.answer}</p>
                          </div>
                          <p className="text-blue-800 text-sm bg-blue-50 p-3 rounded">{attempt.feedback}</p>
                      </div>
                  ))}
              </div>
          </div>

          <div className="flex justify-center">
            <button
              onClick={() => setStep('setup')}
              className="flex items-center gap-2 px-8 py-3 bg-blue-600 text-white rounded-lg hover:bg-blue-700 font-medium text-lg shadow-sm"
            >
              <RefreshCw className="w-5 h-5" />
              {t.restart}
            </button>
          </div>
        </div>
      </div>
    );
  }

  // Session View
  const currentRound = ROUNDS.find(r => r.id === round) || ROUNDS[0];

  return (
    <div className="max-w-4xl mx-auto p-6">
      {/* Header */}
      <div className="mb-6 flex items-center justify-between">
        <div className="flex items-center gap-4">
          <div className="p-3 bg-white rounded-lg shadow-sm border border-gray-100">
            {currentRound.icon}
          </div>
          <div>
            <h2 className="text-xl font-bold text-gray-900 flex items-center gap-2">
              {currentRound.title[language]}
              <span className="text-sm font-normal text-gray-500 px-2 py-0.5 bg-gray-100 rounded-full">
                Round {currentRoundIndex + 1}/8
              </span>
            </h2>
            <p className="text-sm text-gray-500">{currentRound.description[language]}</p>
          </div>
        </div>
        <div className="text-right">
          <div className="text-sm text-gray-500">Question</div>
          <div className="text-2xl font-bold text-blue-600">{currentQuestionIndex} <span className="text-gray-300 text-lg">/ {totalQuestions}</span></div>
        </div>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
        {/* Main Content */}
        <div className="lg:col-span-2 space-y-6">
          {/* Question Card */}
          <div className="bg-white rounded-xl shadow-sm border border-gray-100 p-6">
            <h3 className="text-lg font-medium text-gray-900 mb-4 flex items-center gap-2">
              <MessageSquare className="w-5 h-5 text-blue-500" />
              Question
            </h3>
            {loading && !question ? (
              <div className="animate-pulse space-y-3">
                <div className="h-4 bg-gray-200 rounded w-3/4"></div>
                <div className="h-4 bg-gray-200 rounded w-1/2"></div>
              </div>
            ) : error ? (
              <div className="p-4 bg-red-50 text-red-700 rounded-lg border border-red-200">
                <div className="flex items-center gap-2 mb-1">
                  <AlertCircle className="w-5 h-5" />
                  <p className="font-medium">Error loading question</p>
                </div>
                <p className="text-sm ml-7">{error}</p>
                <button 
                  onClick={() => {
                    setError('');
                    setLoading(true);
                    if (session?.id) {
                         const qId = nextQuestionId || session.current_question_id;
                         const url = qId 
                           ? getApiUrl(`/sessions/${session.id}/questions/${qId}`)
                           : getApiUrl(`/sessions/${session.id}/questions/random`);
                           
                         axios.get(url, {
                             params: { 
                               topic: round === 'core_tech' ? undefined : ROUNDS.find(r => r.id === round)?.title.en 
                             }
                         })
                         .then(res => {
                             setQuestion(res.data);
                             setLoading(false);
                         })
                         .catch(err => {
                             setError(err.response?.data?.error || 'Failed to load question');
                             setLoading(false);
                         });
                    }
                  }} 
                  className="ml-7 mt-2 text-sm font-medium hover:underline"
                >
                  Retry Loading Question
                </button>
              </div>
            ) : (
              <div className="prose prose-blue max-w-none">
                <p className="text-gray-800 text-lg leading-relaxed">
                  {question?.content || 'Loading question...'}
                </p>
                {enableHints && question?.hint && (
                  <div className="mt-4 p-3 bg-amber-50 border border-amber-100 rounded-lg text-sm text-amber-800 flex gap-2 items-start">
                    <Lightbulb className="w-4 h-4 mt-0.5 shrink-0" />
                    <div>
                      <span className="font-bold">{t.hint_title}: </span>
                      {question.hint}
                    </div>
                  </div>
                )}
              </div>
            )}
          </div>

          {/* Answer Area */}
          <div className="bg-white rounded-xl shadow-sm border border-gray-100 p-6">
            <div className="flex items-center justify-between mb-4">
              <h3 className="text-lg font-medium text-gray-900">{t.your_answer}</h3>
              {enableVoice && (
                <button
                  onClick={isListening ? stopListening : startListening}
                  className={`p-2 rounded-full transition-colors ${
                    isListening 
                      ? 'bg-red-100 text-red-600 animate-pulse' 
                      : 'bg-gray-100 text-gray-600 hover:bg-gray-200'
                  }`}
                >
                  {isListening ? <MicOff className="w-5 h-5" /> : <Mic className="w-5 h-5" />}
                </button>
              )}
            </div>
            
            {!attempt ? (
              <>
                <textarea
                  value={answer}
                  onChange={(e) => setAnswer(e.target.value)}
                  placeholder={isListening ? t.listening : "Type your answer here..."}
                  className="w-full h-40 p-4 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent resize-none"
                />
                <div className="mt-4 flex justify-end">
                  <button
                    onClick={submitAnswer}
                    disabled={!answer.trim() || loading}
                    className="flex items-center gap-2 px-6 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 disabled:opacity-50 font-medium"
                  >
                    {loading ? <RefreshCw className="w-4 h-4 animate-spin" /> : <Send className="w-4 h-4" />}
                    {t.submit}
                  </button>
                </div>
              </>
            ) : (
              <div className="space-y-6">
                <div className="p-4 bg-gray-50 rounded-lg border border-gray-200">
                  <p className="text-gray-700 whitespace-pre-wrap">{answer}</p>
                </div>
                
                {/* Feedback */}
                <div className="bg-blue-50 rounded-lg p-6 border border-blue-100">
                  <div className="flex items-center justify-between mb-4">
                    <h4 className="font-bold text-blue-900 flex items-center gap-2">
                      <Shield className="w-5 h-5" />
                      {t.feedback}
                    </h4>
                    <div className="flex items-center gap-2">
                        <span className="text-sm text-blue-700">{t.score}:</span>
                        <span className={`text-xl font-bold ${
                            (attempt.score || 0) >= 70 ? 'text-green-600' : 
                            (attempt.score || 0) >= 50 ? 'text-amber-600' : 'text-red-600'
                        }`}>
                            {attempt.score}/100
                        </span>
                    </div>
                  </div>
                  <p className="text-blue-800 leading-relaxed mb-4">{attempt.feedback}</p>
                  
                  {attempt.suggested_answer && (
                    <div className="pt-4 border-t border-blue-200">
                        <h5 className="text-sm font-bold text-blue-900 mb-1">{t.suggested_answer}:</h5>
                        <p className="text-sm text-blue-800 opacity-90">{attempt.suggested_answer}</p>
                    </div>
                  )}
                </div>

                <div className="flex justify-end">
                  <button
                    onClick={proceedToNextQuestion}
                    className="flex items-center gap-2 px-6 py-2 bg-green-600 text-white rounded-lg hover:bg-green-700 font-medium"
                  >
                    {t.next}
                    <ChevronRight className="w-4 h-4" />
                  </button>
                </div>
              </div>
            )}
          </div>
        </div>

        {/* Sidebar */}
        <div className="space-y-6">
          {/* Progress */}
          <div className="bg-white rounded-xl shadow-sm border border-gray-100 p-6">
            <h3 className="font-medium text-gray-900 mb-4">Session Progress</h3>
            <div className="space-y-4">
                {ROUNDS.map((r, idx) => (
                    <div key={r.id} className="flex items-center gap-3">
                        <div className={`w-6 h-6 rounded-full flex items-center justify-center text-xs font-bold ${
                            idx < currentRoundIndex ? 'bg-green-100 text-green-700' :
                            idx === currentRoundIndex ? 'bg-blue-600 text-white' :
                            'bg-gray-100 text-gray-400'
                        }`}>
                            {idx < currentRoundIndex ? <CheckCircle className="w-4 h-4" /> : idx + 1}
                        </div>
                        <div className="flex-1">
                            <div className={`text-sm font-medium ${
                                idx === currentRoundIndex ? 'text-blue-700' : 'text-gray-600'
                            }`}>
                                {r.title[language]}
                            </div>
                        </div>
                    </div>
                ))}
            </div>
          </div>
          
          <div className="bg-gray-50 rounded-xl p-4 border border-gray-100">
              <button 
                  onClick={() => {
                      if (confirm('Are you sure you want to end the session?')) {
                          setStep('setup');
                      }
                  }}
                  className="w-full flex items-center justify-center gap-2 text-red-600 hover:text-red-700 text-sm font-medium"
              >
                  <XCircle className="w-4 h-4" />
                  End Session
              </button>
          </div>
        </div>
      </div>
    </div>
  );
}
