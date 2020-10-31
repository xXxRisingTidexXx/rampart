# TODO: add distroless image (https://github.com/GoogleContainerTools/distroless).
# TODO: add non-root user.
FROM python:3.8-slim

WORKDIR /app

ENV PYTHONUNBUFFERED 1

COPY auge.pth .
COPY requirements requirements
COPY rampart rampart

RUN pip install -U pip && pip install -r requirements/recognition.txt