FROM migrate/migrate:v4.12.2

ARG RAMPART_UID
ARG RAMPART_GID

RUN addgroup --gid $RAMPART_GID rampart && adduser --disabled-password --gecos '' --uid $RAMPART_UID --ingroup rampart rampart

USER rampart
