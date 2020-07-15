FROM jupyter/scipy-notebook:ea01ec4d9f57

RUN pip install 'psycopg2-binary==2.8.5'
