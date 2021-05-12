FROM python:3.8-slim

WORKDIR /app

ENV PYTHONUNBUFFERED 1

COPY requirements requirements
COPY rampart rampart
COPY config config
COPY templates templates

RUN apt-get update && \
    apt-get install -y libgomp1 && \
    pip install -U pip && \
    pip install -r requirements/python.txt
