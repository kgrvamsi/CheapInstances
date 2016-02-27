# For scrapy updater
FROM ubuntu:14.04
MAINTAINER KGR VAMSI <kgrvamsi@yahoo.com>
RUN apt-get update &&
    apt-get clean &&
    rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*
RUN 
