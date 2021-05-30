FROM python:3.8-slim
WORKDIR /app
ENV PYTHONUNBUFFERED 1
COPY requirements requirements
RUN apt-get update && apt-get install -y libgomp1 && pip install -r requirements/python.txt
