# Slurm Cluster Scaffold

This is the GPU Slurm cluster that Caravan carries inside its binary. It's
embedded via `//go:embed slurm-cluster/*` in [embed.go](../embed.go) at the
repo root, then extracted to `CARAVAN_DIR` (default `~/.caravan/cluster`) and
run with `docker`/`podman compose` by `caravan cluster up`.

Editing files here changes what gets embedded in the next build — there's no
separate generated copy to keep in sync.

## Files

- **Dockerfile** — Ubuntu 24.04 image with `slurm-wlm` + `munge` installed.
- **entrypoint.sh** — starts `munged`, then execs `slurmctld` or `slurmd`
  depending on the container's `command`.
- **docker-compose.yml** — one `slurmctld` controller plus two `slurmd`
  compute nodes (`c1`, `c2`), each advertising `gpu:4` as fake, count-only
  GPUs (no real hardware, no `nvidia-smi`).
- **slurm.conf** — cluster config. `NodeName` CPUs/RealMemory must not
  exceed what Docker actually gives the container, or the node registers
  invalid/drained.
- **cgroup.conf** — `CgroupPlugin=cgroup/v2` + `IgnoreSystemd=yes`, since
  these containers have no systemd/dbus to manage cgroup scopes for them.
- **gres.conf** — declares the fake `gpu:4` GRES per node.

## Known quirks

- `c1`/`c2` run `privileged: true` + `cgroup: host` + `cgroup_parent: "/"`.
  slurmd's `cgroup/v2` plugin assumes its own container cgroup sits directly
  under the cgroup mount root (so it can create a `system.slice/<scope>`
  sibling) — true by default on hosts using the `systemd` cgroup driver, but
  not on hosts using `cgroupfs` (e.g. Docker Desktop on Windows/WSL2). The
  compose settings force that layout regardless of host driver.
- Nodes report `gres/gpu count reported lower than configured (0 < 4)` and
  drain a few seconds after slurmd starts, since these are file-less fake
  GPUs. Submitted jobs will queue but stay pending — not yet fixed.
