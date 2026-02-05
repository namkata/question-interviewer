import uuid
from typing import Any, Dict, List
from urllib.parse import urlparse, parse_qs
from .base import BaseCrawler

class VietnamMarketCrawler(BaseCrawler):
    def __init__(self):
        self.data = [
            # FrontEnd - Fresher/Junior
            {
                "title": "Phân biệt `var`, `let` và `const` trong JavaScript?",
                "content": """
                Câu hỏi phỏng vấn vị trí Frontend Developer (Fresher/Junior).
                
                Hãy giải thích sự khác nhau giữa `var`, `let` và `const` về phạm vi (scope), hoisting và khả năng gán lại giá trị.
                Khi nào nên dùng cái nào?
                """,
                "source": "Vietnam IT Job Market",
                "role": "FrontEnd",
                "level": "Fresher",
                "tags": ["JavaScript", "Frontend Basic"],
                "hint": "Tập trung vào Scope (Function vs Block), Hoisting (có được khởi tạo không?), và Reassignment (gán lại).",
                "correct_answer": """
                1. **Scope**: `var` có function scope, `let` và `const` có block scope.
                2. **Hoisting**: `var` được hoisted và khởi tạo với `undefined`. `let` và `const` được hoisted nhưng nằm trong 'Temporal Dead Zone' cho đến khi được khai báo.
                3. **Reassignment**: `var` và `let` có thể gán lại, `const` không thể gán lại (nhưng object gán cho const vẫn có thể mutate).
                4. **Best Practice**: Ưu tiên dùng `const`, nếu cần gán lại thì dùng `let`. Tránh dùng `var`.
                """
            },
            {
                "title": "React Lifecycle Methods (Class vs Hooks)",
                "content": """
                Câu hỏi phỏng vấn ReactJS.
                
                So sánh Lifecycle trong Class Component (componentDidMount, componentDidUpdate, componentWillUnmount) với `useEffect` trong Functional Component.
                """,
                "source": "Vietnam IT Job Market",
                "role": "FrontEnd",
                "level": "Junior",
                "tags": ["React", "JavaScript"],
                "hint": "Liên hệ `useEffect` với dependency array rỗng [], có dependency [prop], và return function (cleanup).",
                "correct_answer": """
                - **componentDidMount**: Tương đương `useEffect(() => { ... }, [])` (chạy 1 lần sau render).
                - **componentDidUpdate**: Tương đương `useEffect(() => { ... }, [prop])` (chạy khi dependency thay đổi).
                - **componentWillUnmount**: Tương đương function return trong `useEffect` (cleanup function).
                - `useEffect` linh hoạt hơn vì có thể gom logic liên quan lại với nhau thay vì tách rời theo lifecycle method.
                """
            },
            {
                "title": "Vue.js Computed vs Watchers",
                "content": """
                Câu hỏi phỏng vấn Vue.js.
                
                Khi nào nên sử dụng Computed Properties? Khi nào nên sử dụng Watchers?
                Sự khác biệt về caching của chúng.
                """,
                "source": "Vietnam IT Job Market",
                "role": "FrontEnd",
                "level": "Junior",
                "tags": ["Vue", "JavaScript"],
                "hint": "Computed dựa trên sự phụ thuộc và có caching. Watcher dùng cho side-effects (API call, timer).",
                "correct_answer": """
                - **Computed Properties**: Được cache dựa trên reactive dependencies. Chỉ tính toán lại khi dependency thay đổi. Dùng cho việc biến đổi dữ liệu để hiển thị.
                - **Watchers**: Không có caching. Chạy mỗi khi data thay đổi. Dùng cho side-effects như gọi API không đồng bộ, thao tác DOM thủ công, hoặc logic phức tạp khi data đổi.
                """
            },
             {
                "title": "Angular Dependency Injection",
                "content": """
                Câu hỏi phỏng vấn Angular.
                
                Giải thích cơ chế Dependency Injection trong Angular. 
                Sự khác biệt giữa `providedIn: 'root'` và providers trong module/component.
                """,
                "source": "Vietnam IT Job Market",
                "role": "FrontEnd",
                "level": "Mid",
                "tags": ["Angular", "TypeScript"],
                "hint": "Singleton service vs Multiple instances. Tree shaking.",
                "correct_answer": """
                - **Dependency Injection (DI)**: Design pattern nơi class yêu cầu dependencies từ bên ngoài thay vì tự tạo. Angular có DI system tích hợp sẵn.
                - **providedIn: 'root'**: Tạo singleton service cho toàn bộ ứng dụng. Hỗ trợ Tree Shaking (loại bỏ code thừa nếu không dùng).
                - **providers trong Module/Component**: Tạo instance riêng cho Module/Component đó (và con của nó). Không phải singleton toàn cục.
                """
            },
            
            # BackEnd - Junior/Mid/Senior
            {
                "title": "Sự khác biệt giữa TCP và UDP?",
                "content": """
                Câu hỏi mạng máy tính cơ bản cho Backend Developer.
                
                Giải thích sự khác biệt chính giữa giao thức TCP và UDP. Khi nào nên dùng TCP, khi nào dùng UDP?
                Ví dụ ứng dụng thực tế.
                """,
                "source": "Vietnam IT Job Market",
                "role": "BackEnd",
                "level": "Junior",
                "tags": ["Network", "General"],
                "hint": "Độ tin cậy vs Tốc độ. Connection-oriented vs Connectionless.",
                "correct_answer": """
                - **TCP (Transmission Control Protocol)**: Hướng kết nối, đảm bảo độ tin cậy (gửi lại gói tin lỗi, đúng thứ tự), chậm hơn. Dùng cho Web (HTTP), Email (SMTP), File Transfer (FTP).
                - **UDP (User Datagram Protocol)**: Không kết nối, không đảm bảo tin cậy (có thể mất gói), nhanh hơn. Dùng cho Streaming, Gaming, VoIP, DNS.
                """
            },
            {
                "title": "Giải thích về Database Indexing",
                "content": """
                Câu hỏi tối ưu hóa cơ sở dữ liệu.
                
                Index là gì? Tại sao Index giúp tăng tốc độ truy vấn?
                Nhược điểm của việc đánh quá nhiều Index là gì?
                B-Tree Index hoạt động như thế nào?
                """,
                "source": "Vietnam IT Job Market",
                "role": "BackEnd",
                "level": "Mid",
                "tags": ["Database", "SQL"],
                "hint": "Giống như mục lục sách. Trade-off giữa Read speed và Write speed.",
                "correct_answer": """
                - **Index**: Cấu trúc dữ liệu (thường là B-Tree) giúp tìm kiếm nhanh hơn mà không cần quét toàn bộ bảng (Full Table Scan).
                - **Ưu điểm**: Tăng tốc độ SELECT (WHERE, JOIN, ORDER BY).
                - **Nhược điểm**: Giảm tốc độ INSERT/UPDATE/DELETE (vì phải cập nhật cả Index), tốn dung lượng lưu trữ.
                - **B-Tree**: Cấu trúc cây cân bằng, giúp tìm kiếm với độ phức tạp O(log n).
                """
            },
             {
                "title": "Golang Goroutines vs Threads",
                "content": """
                Câu hỏi chuyên sâu về Golang.
                
                Goroutine khác gì so với OS Thread? Tại sao Goroutine lại nhẹ hơn (lightweight)?
                Cơ chế M:N scheduling trong Go Runtime.
                """,
                "source": "Vietnam IT Job Market",
                "role": "BackEnd",
                "level": "Mid",
                "tags": ["Golang"],
                "hint": "Stack size, Scheduling (OS vs Runtime).",
                "correct_answer": """
                - **Stack Size**: Goroutine bắt đầu rất nhỏ (~2KB) và grow/shrink động. OS Thread lớn hơn nhiều (1-2MB).
                - **Scheduling**: Goroutine được quản lý bởi Go Runtime (User space), context switch rất nhanh. OS Thread quản lý bởi Kernel, context switch tốn kém.
                - **M:N Scheduling**: Go Runtime map M Goroutines lên N OS Threads, tận dụng đa nhân nhưng vẫn giữ chi phí thấp.
                """
            },
            {
                "title": "Python Global Interpreter Lock (GIL)",
                "content": """
                Câu hỏi phỏng vấn Python.
                
                GIL là gì và nó ảnh hưởng thế nào đến multi-threading trong Python?
                Làm sao để vượt qua giới hạn của GIL để tận dụng đa nhân CPU?
                """,
                "source": "Vietnam IT Job Market",
                "role": "BackEnd",
                "level": "Mid",
                "tags": ["Python"],
                "hint": "Mutex lock. CPU-bound vs I/O-bound. Multiprocessing.",
                "correct_answer": """
                - **GIL**: Mutex bảo vệ truy cập vào Python objects, ngăn cản nhiều thread thực thi Python bytecodes cùng lúc trên một process.
                - **Ảnh hưởng**: Multi-threading trong Python không tăng tốc được các tác vụ CPU-bound (tính toán), chỉ tốt cho I/O-bound.
                - **Giải pháp**: Sử dụng `multiprocessing` (tạo process riêng, mỗi process có GIL riêng) để tận dụng đa nhân cho CPU-bound tasks.
                """
            },
            {
                "title": "Java Memory Management & Garbage Collection",
                "content": """
                Câu hỏi phỏng vấn Java Backend.
                
                Giải thích cơ chế quản lý bộ nhớ trong Java (Heap, Stack).
                Garbage Collection hoạt động như thế nào? Các thuật toán GC phổ biến (G1, CMS).
                """,
                "source": "Vietnam IT Job Market",
                "role": "BackEnd",
                "level": "Senior",
                "tags": ["Java"],
                "hint": "Stack lưu local vars/method calls. Heap lưu Objects. Mark and Sweep.",
                "correct_answer": """
                - **Stack**: Lưu trữ biến cục bộ, tham chiếu method. Tự động giải phóng khi method kết thúc (LIFO).
                - **Heap**: Lưu trữ Objects (new Keyword). Được quản lý bởi GC.
                - **GC (Garbage Collection)**: Tự động tìm và xóa object không còn được tham chiếu.
                - **Thuật toán**: Mark-and-Sweep (đánh dấu và quét), Generational (Young/Old Gen). G1 GC chia heap thành các vùng nhỏ để tối ưu pause time.
                """
            },
            {
                "title": "Microservices Design Patterns",
                "content": """
                Câu hỏi System Design cho vị trí Senior/Architect.
                
                Các pattern phổ biến: API Gateway, Circuit Breaker, Saga Pattern (cho distributed transactions).
                Khi nào nên chuyển từ Monolith sang Microservices?
                """,
                "source": "Vietnam IT Job Market",
                "role": "BackEnd",
                "level": "Senior",
                "tags": ["System Design", "Architecture"],
                "hint": "Scale độc lập, Fault isolation. Complexity management.",
                "correct_answer": """
                - **API Gateway**: Cổng vào duy nhất cho client, xử lý routing, auth, rate limiting.
                - **Circuit Breaker**: Ngắt kết nối khi service con bị lỗi để tránh cascade failure.
                - **Saga Pattern**: Quản lý transaction phân tán bằng chuỗi các local transactions (có bù trừ - compensation).
                - **Chuyển đổi khi**: Monolith quá lớn, khó maintain, deploy chậm, cần scale từng phần riêng biệt, team size lớn cần làm việc độc lập.
                """
            },
            {
                "title": "Node.js Event Loop",
                "content": """
                Câu hỏi phỏng vấn Node.js.
                
                Giải thích cơ chế Event Loop trong Node.js. 
                Các giai đoạn (Phases) của Event Loop: Timers, Pending Callbacks, Poll, Check, Close Callbacks.
                Microtasks (Promise) vs Macrotasks (setTimeout).
                """,
                "source": "Vietnam IT Job Market",
                "role": "BackEnd",
                "level": "Mid",
                "tags": ["Node.js", "JavaScript"],
                "hint": "Single Threaded nhưng Non-blocking I/O. Call Stack -> Microtasks -> Event Loop Phases.",
                "correct_answer": """
                - **Event Loop**: Cơ chế giúp Node.js thực hiện non-blocking I/O bằng cách offload operations cho system kernel.
                - **Phases**: Timers (setTimeout) -> Pending Callbacks -> Poll (I/O) -> Check (setImmediate) -> Close.
                - **Microtasks (Promise.then, process.nextTick)**: Có độ ưu tiên cao hơn, được thực thi ngay sau khi operation hiện tại hoàn thành và trước khi chuyển phase tiếp theo.
                """
            },

            # DevOps
            {
                "title": "Docker vs Virtual Machine",
                "content": """
                Câu hỏi cơ bản về Containerization.
                
                So sánh Docker Container và Virtual Machine (VM).
                Lợi ích của việc sử dụng Docker trong quy trình CI/CD.
                """,
                "source": "Vietnam IT Job Market",
                "role": "DevOps",
                "level": "Junior",
                "tags": ["Docker", "DevOps"],
                "hint": "OS Kernel sharing vs Full OS guest. Lightweight vs Heavy.",
                "correct_answer": """
                - **Docker Container**: Chia sẻ OS Kernel của Host, nhẹ, khởi động nhanh (giây), portable.
                - **VM**: Có Guest OS riêng, nặng, khởi động chậm (phút), cách ly tốt hơn về bảo mật (hardware level).
                - **Lợi ích CI/CD**: Môi trường nhất quán (Dev=Prod), build nhanh, dễ scale.
                """
            },
            {
                "title": "Kubernetes Pod Lifecycle",
                "content": """
                Câu hỏi về Kubernetes.
                
                Mô tả vòng đời của một Pod trong Kubernetes. Các trạng thái (Pending, Running, Succeeded, Failed, Unknown).
                Làm thế nào để debug khi Pod bị CrashLoopBackOff?
                """,
                "source": "Vietnam IT Job Market",
                "role": "DevOps",
                "level": "Mid",
                "tags": ["Kubernetes", "DevOps"],
                "hint": "Pending -> ContainerCreating -> Running. kubectl describe/logs.",
                "correct_answer": """
                - **Lifecycle**: Pending (đang schedule/pull image) -> Running (container đã start) -> Succeeded/Failed (kết thúc).
                - **CrashLoopBackOff**: Container start rồi crash liên tục.
                - **Debug**: Dùng `kubectl describe pod <name>` để xem events. Dùng `kubectl logs <name>` để xem log lỗi của ứng dụng (thường do config sai, thiếu env, code lỗi).
                """
            },
            {
                "title": "CI/CD Pipeline Best Practices",
                "content": """
                Câu hỏi về quy trình CI/CD.
                
                Thiết kế một pipeline CI/CD an toàn và hiệu quả.
                Khái niệm Blue-Green Deployment và Canary Deployment.
                """,
                "source": "Vietnam IT Job Market",
                "role": "DevOps",
                "level": "Senior",
                "tags": ["CI/CD", "DevOps"],
                "hint": "Fail fast. Zero downtime deployment.",
                "correct_answer": """
                - **Best Practices**: Commit code thường xuyên, build once deploy anywhere (artifacts), test tự động (unit, integration), security scan, môi trường staging giống prod.
                - **Blue-Green**: 2 môi trường song song (Blue=Live, Green=New). Switch traffic sang Green khi test xong. Rollback nhanh.
                - **Canary**: Deploy bản mới cho một lượng nhỏ user trước (ví dụ 5%), monitor, nếu ổn thì roll out toàn bộ.
                """
            },
            {
                "title": "Infrastructure as Code (Terraform)",
                "content": """
                Câu hỏi về IaC.
                
                Tại sao nên dùng Terraform (hoặc IaC nói chung) thay vì cấu hình thủ công?
                Giải thích về Terraform State và cách quản lý State trong team (Remote State).
                """,
                "source": "Vietnam IT Job Market",
                "role": "DevOps",
                "level": "Mid",
                "tags": ["Terraform", "DevOps"],
                "hint": "Consistency, Version Control. State file lưu mapping resource.",
                "correct_answer": """
                - **Lợi ích IaC**: Nhất quán, tránh cấu hình sai (human error), có version control (Git), tái sử dụng, tài liệu hóa hạ tầng.
                - **Terraform State**: File JSON lưu trạng thái hiện tại của hạ tầng. Terraform dùng nó để so sánh với thực tế và plan thay đổi.
                - **Remote State (S3 + DynamoDB)**: Để chia sẻ state giữa các thành viên trong team, tránh conflict và mất mát state file (locking).
                """
            },

            # Data Engineer
            {
                "title": "ETL vs ELT",
                "content": """
                Câu hỏi cơ bản về Data Engineering.
                
                Phân biệt quy trình ETL (Extract, Transform, Load) và ELT (Extract, Load, Transform).
                Khi nào nên dùng ELT?
                """,
                "source": "Vietnam IT Job Market",
                "role": "Data Engineer",
                "level": "Junior",
                "tags": ["Data Engineering", "ETL"],
                "hint": "Thứ tự xử lý. Sức mạnh của Modern Data Warehouse.",
                "correct_answer": """
                - **ETL**: Transform dữ liệu trên server riêng trước khi Load vào Warehouse. Thường dùng cho hệ thống cũ, tốn resource transform.
                - **ELT**: Load dữ liệu thô vào Warehouse trước, sau đó Transform bằng sức mạnh của Warehouse (BigQuery, Snowflake).
                - **Dùng ELT khi**: Dữ liệu lớn (Big Data), sử dụng Cloud Data Warehouse hiện đại, cần tốc độ load nhanh.
                """
            },
            {
                "title": "Spark RDD vs DataFrame vs Dataset",
                "content": """
                Câu hỏi về Apache Spark.
                
                So sánh RDD, DataFrame và Dataset trong Spark.
                Tại sao DataFrame/Dataset thường nhanh hơn RDD?
                """,
                "source": "Vietnam IT Job Market",
                "role": "Data Engineer",
                "level": "Mid",
                "tags": ["Spark", "Big Data"],
                "hint": "Low-level vs High-level API. Optimization (Catalyst Optimizer).",
                "correct_answer": """
                - **RDD**: Low-level API, mạnh mẽ nhưng không tối ưu tự động, code dài dòng.
                - **DataFrame**: High-level API, dữ liệu có cấu trúc (schema), được tối ưu bởi Catalyst Optimizer.
                - **Dataset**: Type-safe (như RDD) nhưng có tối ưu (như DataFrame). Có trong Scala/Java.
                - **Hiệu năng**: DataFrame/Dataset nhanh hơn vì Spark hiểu cấu trúc dữ liệu và tối ưu kế hoạch thực thi (Query Plan).
                """
            },
             {
                "title": "Data Warehouse vs Data Lake",
                "content": """
                Câu hỏi kiến trúc dữ liệu.
                
                Sự khác biệt giữa Data Warehouse và Data Lake.
                Mô hình Data Lakehouse là gì?
                """,
                "source": "Vietnam IT Job Market",
                "role": "Data Engineer",
                "level": "Senior",
                "tags": ["Data Architecture"],
                "hint": "Structured vs Unstructured. Schema-on-write vs Schema-on-read.",
                "correct_answer": """
                - **Data Warehouse**: Lưu dữ liệu có cấu trúc (Structured), đã qua xử lý sạch. Dùng cho BI, Reporting. (Schema-on-write).
                - **Data Lake**: Lưu mọi loại dữ liệu (Raw, Structured, Unstructured). Chi phí rẻ. Dùng cho ML, Exploration. (Schema-on-read).
                - **Data Lakehouse**: Kết hợp ưu điểm cả hai: Lưu trữ rẻ của Data Lake + Quản lý transaction/schema của Warehouse (ACID).
                """
            },
            {
                "title": "SQL Window Functions",
                "content": """
                Câu hỏi SQL nâng cao.
                
                Window Functions là gì? Khác gì với GROUP BY?
                Ví dụ về `ROW_NUMBER()`, `RANK()`, `DENSE_RANK()`.
                """,
                "source": "Vietnam IT Job Market",
                "role": "Data Engineer",
                "level": "Mid",
                "tags": ["SQL", "Data Engineering"],
                "hint": "Tính toán trên tập hợp dòng liên quan mà không gom nhóm (collapse) như GROUP BY.",
                "correct_answer": """
                - **Window Functions**: Thực hiện tính toán trên một tập hợp các dòng (window) liên quan đến dòng hiện tại. Không làm giảm số lượng dòng trả về như `GROUP BY`.
                - **ROW_NUMBER()**: Đánh số thứ tự liên tục (1, 2, 3, 4).
                - **RANK()**: Đánh số thứ tự có nhảy cóc khi trùng giá trị (1, 2, 2, 4).
                - **DENSE_RANK()**: Đánh số thứ tự không nhảy cóc khi trùng giá trị (1, 2, 2, 3).
                """
            },
            {
                "title": "Apache Kafka Basics",
                "content": """
                Câu hỏi về hệ thống message streaming.
                
                Apache Kafka là gì? Các thành phần chính: Producer, Consumer, Broker, Topic, Partition, Offset.
                Tại sao Kafka lại có throughput cao?
                """,
                "source": "Vietnam IT Job Market",
                "role": "Data Engineer",
                "level": "Mid",
                "tags": ["Kafka", "Big Data"],
                "hint": "Distributed Event Streaming. Sequential I/O. Zero Copy.",
                "correct_answer": """
                - **Kafka**: Nền tảng phân tán để xử lý luồng sự kiện (event streaming).
                - **Thành phần**: Producer (gửi tin), Consumer (nhận tin), Broker (server), Topic (kênh), Partition (chia nhỏ topic để scale), Offset (vị trí tin nhắn).
                - **Hiệu năng cao**: Nhờ ghi đĩa tuần tự (Sequential I/O) thay vì ngẫu nhiên, và sử dụng Zero Copy để chuyển dữ liệu từ disk sang network buffer.
                """
            },
            {
                "title": "Airflow & Workflow Orchestration",
                "content": """
                Câu hỏi về công cụ điều phối workflow.
                
                Apache Airflow là gì? Khái niệm DAG (Directed Acyclic Graph).
                Làm thế nào để xử lý backfill dữ liệu trong Airflow?
                """,
                "source": "Vietnam IT Job Market",
                "role": "Data Engineer",
                "level": "Senior",
                "tags": ["Airflow", "Data Engineering"],
                "hint": "Python-based. Scheduler. Operators.",
                "correct_answer": """
                - **Airflow**: Platform để lập lịch và giám sát workflow.
                - **DAG**: Đồ thị có hướng không chu trình, biểu diễn luồng công việc.
                - **Backfill**: Chạy lại các task trong quá khứ (ví dụ khi sửa lỗi logic hoặc thêm metric mới) bằng cách chỉ định ngày bắt đầu và kết thúc cũ.
                """
            },
            {
                "title": "Redis Use Cases",
                "content": """
                Câu hỏi về Caching/NoSQL.
                
                Redis là gì? Các kiểu dữ liệu chính (String, List, Set, Hash, Sorted Set).
                Khi nào nên dùng Redis làm Cache? Khi nào dùng làm Message Broker?
                """,
                "source": "Vietnam IT Job Market",
                "role": "BackEnd",
                "level": "Mid",
                "tags": ["Redis", "NoSQL"],
                "hint": "In-memory. Key-Value. Pub/Sub.",
                "correct_answer": """
                - **Redis**: In-memory data structure store. Rất nhanh.
                - **Data Types**: String (cache), List (queue), Set (unique), Hash (object), Sorted Set (leaderboard).
                - **Cache**: Dùng khi cần truy xuất nhanh dữ liệu ít thay đổi hoặc tính toán tốn kém.
                - **Message Broker**: Dùng Pub/Sub hoặc List/Stream cho các tác vụ thời gian thực nhẹ nhàng (chat, notification).
                """
            },
            {
                "title": "Frontend Performance Optimization",
                "content": """
                Câu hỏi tối ưu hiệu năng Frontend.
                
                Các kỹ thuật tối ưu performance cho web app: Lazy Loading, Code Splitting, Caching, Image Optimization.
                Critical Rendering Path là gì?
                """,
                "source": "Vietnam IT Job Market",
                "role": "FrontEnd",
                "level": "Senior",
                "tags": ["Performance", "Frontend"],
                "hint": "Giảm size bundle. Giảm request. Tối ưu render blocking.",
                "correct_answer": """
                - **Lazy Loading**: Chỉ tải resource (ảnh, component) khi cần thiết (khi scroll tới).
                - **Code Splitting**: Chia nhỏ JS bundle để tải song song hoặc theo route.
                - **Critical Rendering Path**: Chuỗi các bước trình duyệt phải làm để render pixel đầu tiên (HTML -> DOM -> CSSOM -> Render Tree -> Layout -> Paint). Tối ưu bằng cách inline CSS quan trọng, defer JS.
                """
            }
        ]

    async def can_handle(self, url: str) -> bool:
        return url.startswith("vn-market://")

    async def extract(self, url: str) -> List[Dict[str, Any]]:
        parsed = urlparse(url)
        params = parse_qs(parsed.query)
        
        # Filters from query params
        filter_roles = params.get('role', [])
        filter_levels = params.get('level', [])
        filter_langs = params.get('lang', []) # e.g. ?lang=Python&lang=Java
        
        # Also support path-based filtering for backward compatibility if needed
        # But prefer query params for multi-select
        
        results = []
        for item in self.data:
            # 1. Role Filter
            if filter_roles:
                # Check if item['role'] matches any of the requested roles (case-insensitive)
                if not any(r.lower() == item["role"].lower() for r in filter_roles):
                    continue
            
            # 2. Level Filter
            if filter_levels:
                if not any(l.lower() == item["level"].lower() for l in filter_levels):
                    continue

            # 3. Language/Tag Filter (User requests "choose multiple programming languages")
            # We map "lang" param to "tags" in our data
            if filter_langs:
                item_tags = [t.lower() for t in item.get("tags", [])]
                # If item has ANY of the requested languages in its tags
                has_lang = False
                for lang in filter_langs:
                    if lang.lower() in item_tags:
                        has_lang = True
                        break
                # Special case: If user selects a language, but the question is generic for the role (e.g. "General" tag)
                # we might still want to show it? Or strictly filter?
                # User said: "choose multiple programming languages they know to optimize questions"
                # Strict filtering seems safer to start.
                if not has_lang:
                    continue

            # Add metadata for import script
            item["meta_role"] = item["role"]
            item["meta_level"] = item["level"]
            item["meta_tags"] = item["tags"]
            # Hint and Correct Answer are already in item
            
            results.append(item)
            
        return results
