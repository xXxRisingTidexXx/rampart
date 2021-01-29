FROM python:3.8-slim

WORKDIR /app

ENV FLASK_ENV development
ENV PYTHONUNBUFFERED 1
ENV MYPYPATH "${MYPYPATH}:/app"

COPY tox.ini .
COPY requirements requirements
COPY rampart rampart
COPY config config

RUN apt-get update && \
    apt-get install -y libgomp1 && \
    pip install -U pip && \
    pip install -r requirements/common.txt && \
    pip install -r requirements/twinkle.txt && \
    pip install -r requirements/auge.txt && \
    pip install -r requirements/pychecks.txt
