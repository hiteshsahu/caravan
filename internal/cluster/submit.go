package cluster

import (
	"fmt"
	"os"
)

func Submit(scriptPath string) error {
	dir, err := scaffoldDir()
	if err != nil {
		return err
	}
	engine, err := NewEngine()
	if err != nil {
		return err
	}
	if _, err := os.Stat(scriptPath); err != nil {
		return err
	}
	scriptFile, err := os.Open(scriptPath)
	if err != nil {
		return err
	}
	defer scriptFile.Close()

	fmt.Printf("→ submitting %s to local Slurm cluster in %s\n", scriptPath, dir)
	return engine.Exec("slurmctld", scriptFile, "sh", "-c", "cat > /tmp/caravan-job.sh && sbatch --parsable /tmp/caravan-job.sh")
}
