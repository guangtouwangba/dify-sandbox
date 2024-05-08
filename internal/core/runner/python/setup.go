package python

import (
	_ "embed"
	"os"
	"os/exec"
	"path"

	"github.com/langgenius/dify-sandbox/internal/core/runner"
	"github.com/langgenius/dify-sandbox/internal/utils/log"
)

//go:embed python.so
var python_lib []byte

func init() {
	log.Info("initializing python runner environment...")
	// remove /tmp/sandbox-python
	os.RemoveAll("/tmp/sandbox-python")
	os.Remove("/tmp/sandbox-python")

	err := os.MkdirAll("/tmp/sandbox-python", 0755)
	if err != nil {
		log.Panic("failed to create /tmp/sandbox-python")
	}
	err = os.WriteFile("/tmp/sandbox-python/python.so", python_lib, 0755)
	if err != nil {
		log.Panic("failed to write /tmp/sandbox-python/python.so")
	}
	log.Info("python runner environment initialized")
}

func InstallDependencies(requirements string) error {
	if requirements == "" {
		return nil
	}

	runner := runner.TempDirRunner{}
	return runner.WithTempDir([]string{}, func(root_path string) error {
		// create a requirements file
		err := os.WriteFile(path.Join(root_path, "requirements.txt"), []byte(requirements), 0644)
		if err != nil {
			log.Panic("failed to create requirements.txt")
		}

		// install dependencies
		cmd := exec.Command("pip3", "install", "-r", "requirements.txt")

		reader, err := cmd.StdoutPipe()
		if err != nil {
			log.Panic("failed to get stdout pipe of pip3")
		}
		defer reader.Close()

		err = cmd.Start()
		if err != nil {
			log.Panic("failed to start pip3")
		}
		defer cmd.Wait()

		for {
			buf := make([]byte, 1024)
			n, err := reader.Read(buf)
			if err != nil {
				break
			}
			log.Info(string(buf[:n]))
		}

		return nil
	})
}
