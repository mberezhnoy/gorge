ARG PROJ_VERSION=6.3.1

FROM osgeo/proj:${PROJ_VERSION}

RUN mkdir -p ${GITHUB_WORKSPACE}/libproj && \
    cp /usr/lib/libproj.*a ${GITHUB_WORKSPACE}/libproj/ && \
    echo $GITHUB_WORKSPACE && \
    echo $HOME && \
    mkdir -p ${HOME}/libproj && \
    cp /usr/lib/libproj.*a ${HOME}/libproj/ && \
    mkdir -p /github/home/libproj && \
    cp /usr/lib/libproj.*a /github/home/libproj/
