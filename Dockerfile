FROM ubuntu:16.04
WORKDIR /data

#install needed tools
RUN apt update && apt -y install wget && apt -y install vim && apt -y install git && apt -y install zip && apt -y install build-essential && apt -y install make

#install go
RUN wget https://storage.googleapis.com/golang/go1.11.5.linux-amd64.tar.gz
RUN tar -C /usr/local -xzf go1.11.5.linux-amd64.tar.gz
ENV PATH $PATH:/usr/local/go/bin
RUN echo "export PATH=$PATH:/usr/loca/go/bin" >> /etc/environment
RUN go version

#install terraform 
RUN wget https://releases.hashicorp.com/terraform/0.12.12/terraform_0.12.12_linux_amd64.zip
RUN unzip terraform_0.12.12_linux_amd64.zip
RUN mkdir /root/go && mkdir /root/go/bin && mv terraform /root/go/bin
ENV PATH=$PATH:/root/go/bin
RUN echo "PATH=$PATH:/root/go/bin" >> /etc/environment
RUN terraform version

#install huaweicloudstack provider
RUN git clone https://github.com/huaweicloud/terraform-provider-huaweicloudstack /root/go/src/github.com/terraform-providers/terraform-provider-huaweicloudstack
WORKDIR  /root/go/src/github.com/terraform-providers/terraform-provider-huaweicloudstack/
RUN make build
