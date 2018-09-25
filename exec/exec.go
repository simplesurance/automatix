package exec

import (
	"bufio"
	"os/exec"
	"strings"
	"syscall"

	"github.com/pkg/errors"

	"github.com/simplesurance/automatix/log"
)

// Command runs the passed command in a shell in the passed dir.
// If dir is empty, the command is run in /.
// If the command exits with a code != 0, an error is returned
func Command(dir, command string) (output string, err error) {
	cmd := exec.Command("sh", "-c", command)

	outReader, err := cmd.StdoutPipe()
	if err != nil {
		return
	}

	cmd.Stderr = cmd.Stdout
	if len(dir) != 0 {
		cmd.Dir = dir
	} else {
		cmd.Dir = "/"
	}

	log.Debugf("exec: running %q in directory %q", strings.Join(cmd.Args, " "), cmd.Dir)

	err = cmd.Start()
	if err != nil {
		err = errors.Wrap(err, "running command failed")
		return
	}

	in := bufio.NewScanner(outReader)
	for in.Scan() {
		o := in.Text()
		log.Debugln("exec: " + o)
		output += o + "\n"
	}

	err = cmd.Wait()
	if err != nil {
		if ee, ok := err.(*exec.ExitError); ok {
			if status, ok := ee.Sys().(syscall.WaitStatus); ok {
				return output, errors.Wrapf(err, "command exited with code %d: %s", status.ExitStatus(), output)
			}
		}

		return output, errors.Wrapf(err, "command execution failed: %s", output)
	}

	return output, nil
}
