FROM python:3.10-slim

WORKDIR /app

# Install system dependencies if needed
# RUN apt-get update && apt-get install -y gcc

# Copy requirements
COPY requirements.txt .

# Install dependencies
RUN pip install --no-cache-dir -r requirements.txt

# Copy source code
COPY . .

# Expose port
EXPOSE 8000

# Command to run (using uvicorn)
# Assumes main.py has app object
CMD ["uvicorn", "main:app", "--host", "0.0.0.0", "--port", "8000"]
