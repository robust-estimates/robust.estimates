FROM python:3.13.9-slim
WORKDIR /app
COPY scoring/requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt
COPY scoring/app ./app
COPY scoring/generated ./generated
COPY scoring/schema.py ./schema.py
COPY scoring/risk ./risk
COPY scoring/affordability ./affordability
ENV PYTHONPATH=/app/generated:/app PYTHONDONTWRITEBYTECODE=1 PYTHONUNBUFFERED=1
USER 65532:65532
EXPOSE 9093
CMD ["python", "-m", "app.server"]
