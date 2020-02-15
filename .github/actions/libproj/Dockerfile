ARG PROJ_VERSION=6.3.1

FROM osgeo/proj:${PROJ_VERSION}

RUN printenv
RUN ls -la /home

RUN mkdir -p /home/libproj && \
    cp /usr/lib/libproj.*a /home/libproj/
