FROM ubuntu:14.04
MAINTAINER Sudharshan S <sudharsh@gmail.com>

ENV HENCHMAN_URL github.com/sudharsh/henchman

ENV HOME /root

ENV GOLANG_VERSION 1.3
ENV GOLANG_OS linux
ENV GOLANG_ARCH amd64
ENV GOLANG_TARBALL go${GOLANG_VERSION}.${GOLANG_OS}-${GOLANG_ARCH}.tar.gz
ENV GOLANG_URL http://golang.org/dl/${GOLANG_TARBALL}

RUN apt-get update
RUN apt-get install -y wget git make mercurial
RUN ssh-keygen -f $HOME/.ssh/id_rsa

RUN wget -q $GOLANG_URL -O $HOME/${GOLANG_TARBALL}
RUN tar -C /usr/local -xzf $HOME/${GOLANG_TARBALL} 

ENV PATH $PATH:/usr/local/go/bin

ENV GOPATH $HOME/go
ENV PATH $PATH:$GOPATH/bin

ENV HENCHMAN_MODULES_PATH $HOME

ADD test_plan.yaml $HOME/plans/
ADD insecure_private_key $HOME/.ssh/
ADD insecure_private_key.pub $HOME/.ssh/

CMD go get -v $HENCHMAN_URL && \
    henchman -private-keyfile ~/.ssh/id_rsa $HOME/plans/test_plan.yaml





