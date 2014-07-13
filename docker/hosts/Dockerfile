FROM ubuntu:14.04
MAINTAINER Sudharshan S <sudharsh@gmail.com>

RUN apt-get update

RUN apt-get install -y openssh-server
RUN mkdir /var/run/sshd
RUN echo 'SSHD: ALL' >> /etc/hosts.allow
RUN echo 'root:password' | chpasswd

ENV HOME /root
ADD insecure_private_key.pub $HOME/.ssh/

RUN cat $HOME/.ssh/insecure_private_key.pub >> $HOME/.ssh/authorized_keys
RUN chmod 700 $HOME/.ssh
RUN chmod 600 $HOME/.ssh/authorized_keys

EXPOSE 22
CMD    ["/usr/sbin/sshd", "-D"]
