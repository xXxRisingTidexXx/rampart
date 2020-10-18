FROM jupyter/scipy-notebook:ea01ec4d9f57

USER root

RUN apt-get update && \
    apt-get install -y --no-install-recommends libgomp1 && \
    rm -rf /var/lib/apt/lists/*

USER $NB_UID

RUN pip install psycopg2-binary tabulate shapely ppscore lightgbm
