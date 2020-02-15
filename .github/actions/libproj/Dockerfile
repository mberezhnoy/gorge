ARG PROJ_VERSION=6.3.1

FROM osgeo/proj:${PROJ_VERSION}

RUN mkdir -p ${GITHUB_WORKSPACE}/libproj && \
    cp /usr/lib/libproj.*a ${GITHUB_WORKSPACE}/libproj/
