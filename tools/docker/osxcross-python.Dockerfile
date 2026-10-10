# Local input ID: sha256:933b46a028fd786c9c3d426ebabc237e29a15912231ea8de576e95f0e4f41a4c
FROM python:3.13-bookworm AS python-runtime
# Local input ID: sha256:d4c50a5cf3cf35164b83fb1531fe66f52b9224e826944363bf1bc4d7ba3c8d0c
FROM psychicwar-osxcross:go1.24.13-15.5-r2
COPY --from=python-runtime /usr/local /opt/python
ENV PATH="/opt/python/bin:${PATH}" \
    LD_LIBRARY_PATH="/opt/python/lib"
