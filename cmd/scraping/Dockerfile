ARG PYTHON_VERSION

FROM public.ecr.aws/docker/library/python:${PYTHON_VERSION}-slim-bookworm
RUN apt-get update && apt-get install -y --no-install-recommends \
  ca-certificates \
  && rm -rf /var/lib/apt/lists/*
ENV APPDIR /go/src/github.com/acu4git/gimme-scholarship-scraping
COPY . $APPDIR/
WORKDIR $APPDIR
RUN pip install --no-cache-dir -r requirements.txt
ENTRYPOINT [ "python", "main.py" ]