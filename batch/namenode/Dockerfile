FROM bde2020/hadoop-base:2.0.0-hadoop3.2.1-java8

MAINTAINER Ivan Ermilov <ivan.s.ermilov@gmail.com>

COPY PopularTopicsJob.jar /home/PopularTopicsJob.jar

HEALTHCHECK CMD curl -f http://localhost:9870/ || exit 1

ENV HDFS_CONF_dfs_namenode_name_dir=file:///hadoop/dfs/name
ENV JAR_FILEPATH="/home/PopularTopicsJob.jar"
ENV CLASS_TO_RUN="PopularTopicsJob"
ENV PARAMS="/input /output"
RUN mkdir -p /hadoop/dfs/name
VOLUME /hadoop/dfs/name

ADD run.sh /run.sh
RUN chmod a+x /run.sh

RUN apt update
RUN apt install netcat-openbsd
RUN rm -f /tmp/f; mkfifo /tmp/f

EXPOSE 9870

CMD ["/run.sh"]
