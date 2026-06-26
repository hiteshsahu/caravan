FROM ubuntu:24.04

ENV DEBIAN_FRONTEND=noninteractive

RUN apt-get update && apt-get install -y --no-install-recommends \
        slurm-wlm slurm-wlm-basic-plugins \
        munge libmunge2 \
        util-linux iproute2 ca-certificates procps less vim-tiny \
    && rm -rf /var/lib/apt/lists/*

# Shared munge key baked into the image, so every container trusts every other.
# Insecure on purpose — this is a throwaway local dev cluster.
RUN install -d -o munge -g munge -m 0700 /etc/munge /var/lib/munge /var/log/munge /run/munge \
    && dd if=/dev/urandom of=/etc/munge/munge.key bs=1024 count=1 \
    && chown munge:munge /etc/munge/munge.key \
    && chmod 400 /etc/munge/munge.key

RUN install -d -m 0755 /var/spool/slurmctld /var/spool/slurmd /var/log/slurm /run/slurm

COPY entrypoint.sh /usr/local/bin/entrypoint.sh
RUN chmod +x /usr/local/bin/entrypoint.sh

ENTRYPOINT ["/usr/local/bin/entrypoint.sh"]
