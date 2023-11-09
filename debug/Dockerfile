FROM debian

# Install a number of debug utilities
RUN apt update && apt upgrade
RUN apt install -y procps gdb gdbserver strace ltrace net-tools lsof dnsutils iputils-ping \
    host nmap traceroute tcpdump iproute2 netcat-openbsd \
    socat curl wget less jq gpg ca-certificates \
    postgresql-client-15 \
    man manpages man-db manpages-dev

# Install kubectl, per instructions at
# https://kubernetes.io/docs/tasks/tools/install-kubectl-linux/#install-using-native-package-management
RUN curl -fsSL https://packages.cloud.google.com/apt/doc/apt-key.gpg | gpg --dearmor -o /etc/apt/keyrings/kubernetes-archive-keyring.gpg
RUN echo "deb [signed-by=/etc/apt/keyrings/kubernetes-archive-keyring.gpg] https://apt.kubernetes.io/ kubernetes-xenial main" > /etc/apt/sources.list.d/kubernetes.list
RUN apt update
RUN apt install -y kubectl

# Install go
COPY --from=golang:1.21.3 /usr/local/go/ /usr/local/go/
ENV PATH="/usr/local/go/bin:${PATH}"

# Install delve
RUN CGO_ENABLED=0 go install -ldflags "-s -w -extldflags '-static'" github.com/go-delve/delve/cmd/dlv@v1.20.2
ENV PATH="/root/go/bin/:${PATH}}"

# Cool hack to access the root filesystem of the main container
RUN ln -s /proc/1/root /container
