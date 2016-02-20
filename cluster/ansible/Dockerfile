FROM gcr.io/google_containers/deploy-base:v1


RUN mkdir /ansible
WORKDIR /ansible

ADD roles roles
ADD books.json .
ADD configure-local.sh .
