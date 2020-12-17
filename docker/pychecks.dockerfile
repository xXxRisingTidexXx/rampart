FROM python:3.8-slim

WORKDIR /app

ENV FLASK_ENV development
ENV PYTHONUNBUFFERED 1
ENV MYPYPATH "${MYPYPATH}:/app"

COPY tox.ini .
COPY requirements requirements
COPY rampart rampart
COPY templates templates
COPY config config

RUN apt-get update && \
    apt-get install -y libgomp1 && \
    pip install -U pip && \
    pip install -r requirements/ranking.txt && \
    pip install -r requirements/recognition.txt && \
    pip install -r requirements/pychecks.txt
