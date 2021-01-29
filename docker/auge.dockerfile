FROM python:3.8-slim

WORKDIR /app

ENV PYTHONUNBUFFERED 1

COPY requirements requirements
COPY rampart rampart
COPY config config

RUN pip install -U pip && \
    pip install -r requirements/common.txt && \
    pip install -r requirements/auge.txt
