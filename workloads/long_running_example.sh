#!/usr/bin/env bash
#SBATCH --job-name=caravan-long-test
#SBATCH --output=caravan-long-test.out
#SBATCH --time=00:10:00
#SBATCH --ntasks=1
#SBATCH --gres=gpu:1

echo "Hello from Caravan long-running job on $(hostname)"
echo "Holding gpu:1 for 5 minutes — long enough to watch in 'caravan cluster status' or squint."

for i in $(seq 1 30); do
    echo "[$(date -u +%H:%M:%S)] still running ($i/30)"
    sleep 10
done

echo "Done"
