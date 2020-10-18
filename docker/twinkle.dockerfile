FROM python:3.8-slim

WORKDIR /app

COPY requirements.txt .
COPY rampart rampart

RUN apt-get update && \
    apt-get install -y libgomp1 && \
    pip install -r requirements.txt
