FROM ubuntu:16.04
WORKDIR /data

#Call this before you can install anything
RUN apt update

#install needed tools
RUN apt -y install wget
RUN apt -y install vim
RUN apt -y install git
RUN apt -y install zip
RUN apt -y install build-essential
RUN apt -y install make

#install go
RUN wget https://storage.googleapis.com/golang/go1.11.5.linux-amd64.tar.gz
RUN tar -C /usr/local -xzf go1.11.5.linux-amd64.tar.gz
ENV PATH $PATH:/usr/local/go/bin
RUN echo "export PATH=$PATH:/usr/loca/go/bin" >> /etc/environment
RUN go version

#install terraform 
RUN wget https://releases.hashicorp.com/terraform/0.12.12/terraform_0.12.12_linux_amd64.zip
RUN unzip terraform_0.12.12_linux_amd64.zip
RUN mv terraform /usr/local/go/bin
RUN terraform version

#install huaweicloudstack provider
RUN git clone https://github.com/huaweicloud/terraform-provider-huaweicloudstack /data/go/src/github.com/terraform-providers/terraform-provider-huaweicloudstack
#1 cd doesn't change the work directory and you have to use WORKDIR, and I spent half an hour on this stupid Dockerfile feature
#RUN cd ~/go/src/github.com/terraform-providers/terraform-provider-huaweicloudstack/ 
#2 ~ doesn't refer to the home directory in Docerfile
#3 I need to read a docker book before writing any Dockerfile. too much time wasted in stupid stuff
WORKDIR  /data/go/src/github.com/terraform-providers/terraform-provider-huaweicloudstack/
RUN make build
ENV PATH=$PATH:/data/go/bin
RUN echo "PATH=$PATH:/data/go/bin" >> /etc/environment
