FROM jupyter/scipy-notebook:95ccda3619d0

USER root

RUN apt-get update && \
    apt-get install -y --no-install-recommends libgomp1 && \
    rm -rf /var/lib/apt/lists/*

USER $NB_UID

COPY --chown=${NB_UID}:${NB_GID} requirements /requirements

RUN pip install -r /requirements/notebook.txt && \
    fix-permissions $CONDA_DIR && \
    fix-permissions /home/$NB_USER
