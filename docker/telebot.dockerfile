FROM python:3.8-slim

WORKDIR /app

COPY requirements.txt .
COPY rampart rampart
COPY templates templates

RUN apt-get update && \
    apt-get install -y libgomp1 && \
    pip install -U pip && \
    pip install -r requirements.txt
