FROM python:3.13.9-slim
WORKDIR /app
COPY mocks/profiles.py mocks/profiles.json mocks/server.py ./
ENV PYTHONDONTWRITEBYTECODE=1 PYTHONUNBUFFERED=1 PORT=8090
USER 65532:65532
EXPOSE 8090
CMD ["python", "server.py"]
