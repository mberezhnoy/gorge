ARG PROJ_VERSION=6.3.1

FROM osgeo/proj:${PROJ_VERSION}

RUN printenv
RUN ls -la /github
RUN ls -la /github/home
RUN ls -la /github/workspace
RUN ls -la /github/workflow

RUN mkdir -p /github/home/libproj && \
    cp /usr/lib/libproj.*a /github/home/libproj/
