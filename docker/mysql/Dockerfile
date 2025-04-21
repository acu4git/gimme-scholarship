FROM mysql:8.0-debian

RUN apt-get update && apt-get install -y locales \
  && sed -i 's/# ja_JP.UTF-8 UTF-8/ja_JP.UTF-8 UTF-8/' /etc/locale.gen && locale-gen

ENV LANG ja_JP.UTF-8
ENV TZ Asia/Tokyo

EXPOSE 3306