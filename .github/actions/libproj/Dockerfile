ARG PROJ_VERSION=6.3.1

FROM osgeo/proj:${PROJ_VERSION}

RUN ls -la /
RUN ls -la $HOME

RUN mkdir -p $HOME/libproj && \
    cp /usr/lib/libproj.*a $HOME/libproj/
