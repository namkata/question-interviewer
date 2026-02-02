FROM python:3.10-slim

WORKDIR /app

# Install system dependencies for Playwright (if needed later)
# RUN apt-get update && apt-get install -y ...

# Copy requirements
COPY requirements.txt .

# Install dependencies
RUN pip install --no-cache-dir -r requirements.txt

# Copy source code
COPY . .

# Expose port
EXPOSE 8000

# Command to run
CMD ["uvicorn", "main:app", "--host", "0.0.0.0", "--port", "8000"]
