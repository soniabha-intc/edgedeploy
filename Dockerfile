# SPDX-License-Identifier: Apache-2.0
# Copyright (c) 2019-2020 Intel Corporation

FROM alpine:3.12.0 AS app-deps-image

RUN printf "http://nl.alpinelinux.org/alpine/v3.12/main\nhttp://nl.alpinelinux.org/alpine/v3.12/community" >> /etc/apk/repositories
RUN apk add --no-cache sudo
FROM app-deps-image
ARG username=edgedeployusr
ARG user_dir=/home/$username

RUN addgroup -S sudo && adduser -S $username -G sudo
#RUN addgroup docker
#RUN adduser $username docker

RUN echo '%sudo ALL=(ALL) NOPASSWD:ALL' >> /etc/sudoers
# Workaround for sudo error: https://gitlab.alpinelinux.org/alpine/aports/issues/11122
RUN echo 'Set disable_coredump false' >> /etc/sudo.conf

USER $username
WORKDIR /home/edgedeploy
#WORKDIR /root/
COPY ./edgedeploy ./edgedeploy
CMD ["sleep", "100000"]
