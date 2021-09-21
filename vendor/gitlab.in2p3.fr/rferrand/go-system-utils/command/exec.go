// Copyright (c) IN2P3 Computing Centre, IN2P3, CNRS
//
// Contributor(s): Remi Ferrand <remi.ferrand_at_cc.in2p3.fr>, 2017
//
// This software is governed by the CeCILL  license under French law and
// abiding by the rules of distribution of free software.  You can  use,
// modify and/ or redistribute the software under the terms of the CeCILL
// license as circulated by CEA, CNRS and INRIA at the following URL
// "http://www.cecill.info".
//
// As a counterpart to the access to the source code and  rights to copy,
// modify and redistribute granted by the license, users are provided only
// with a limited warranty  and the software's author,  the holder of the
// economic rights,  and the successive licensors  have only  limited
// liability.
//
// In this respect, the user's attention is drawn to the risks associated
// with loading,  using,  modifying and/or developing or reproducing the
// software by the user in light of its specific status of free software,
// that may mean  that it is complicated to manipulate,  and  that  also
// therefore means  that it is reserved for developers  and  experienced
// professionals having in-depth computer knowledge. Users are therefore
// encouraged to load and test the software's suitability as regards their
// requirements in conditions enabling the security of their systems and/or
// data to be ensured and,  more generally, to use and operate it in the
// same conditions as regards security.
//
// The fact that you are presently reading this means that you have had
// knowledge of the CeCILL license and that you accept its terms.
package command

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"time"
)

var execCommand = exec.Command

// ExecArgv is the structure used to represent
// arguments passed to exec.Command
type ExecArgv struct {
	Command                 []string
	AllowedRc               []int
	Environ                 []string
	CleanupEnviron          bool
	ProcessStartedCallbacks ExecCallbacks
	Stdin                   io.Reader
	Stdout                  io.Writer
	Stderr                  io.Writer
}

// ExecResult is the structure used to
// model result of a process execution
type ExecResult struct {
	Stdout   string
	Stderr   string
	Rc       int
	Duration time.Duration
}

// ExecError is the structure used to
// model an error that occured during
// a process execution
type ExecError struct {
	Argv     *ExecArgv
	Stdout   string
	Stderr   string
	Rc       int
	Duration time.Duration
	Cause    error
}

func (e *ExecError) Error() string {
	cmdStr := strings.Join(e.Argv.Command, " ")
	return fmt.Sprintf("command '%s' exited with status %d. stdout='%s'. stderr='%s'\n", cmdStr, e.Rc, e.Stdout, e.Stderr)
}

func (r *ExecResult) String() string {
	return fmt.Sprintf("stdout='%s', stderr='%s', rc=%d, duration=%s\n", r.Stdout, r.Stderr, r.Rc, r.Duration)
}

func Execute(argv *ExecArgv) (*ExecResult, error) {
	var stdoutBuf bytes.Buffer
	var stderrBuf bytes.Buffer
	var rc = -1
	var err error

	cmd := execCommand(argv.Command[0], argv.Command[1:]...)
	if argv.Stdout != nil {
		cmd.Stdout = argv.Stdout
	} else {
		cmd.Stdout = &stdoutBuf
	}

	if argv.Stderr != nil {
		cmd.Stderr = argv.Stderr
	} else {
		cmd.Stderr = &stderrBuf
	}

	if len(argv.Environ) > 0 {
		if argv.CleanupEnviron {
			cmd.Env = argv.Environ
		} else {
			cmd.Env = append(os.Environ(), argv.Environ...)
		}
	}

	if argv.Stdin != nil {
		cmd.Stdin = argv.Stdin
		if err != nil {
			return nil, fmt.Errorf("getting command stdin pipe failed: %s", err.Error())
		}
	}

	start := time.Now()
	if err = cmd.Start(); err != nil {
		return nil, fmt.Errorf("starting command failed: %s", err.Error())
	}
	if !argv.ProcessStartedCallbacks.Empty() {
		go func() {
			argv.ProcessStartedCallbacks.RunCallbacks(cmd)
		}()
	}
	err = cmd.Wait()
	duration := time.Since(start)

	if err != nil {
		// If the command starts but does not complete successfully, the error is of type *ExitError.

		if exiterr, status := err.(*exec.ExitError); status {
			waitStatus, _ := exiterr.Sys().(syscall.WaitStatus)
			rc = waitStatus.ExitStatus()
		}

		allowedRc := false
		for _, arc := range argv.AllowedRc {
			if arc == rc {
				allowedRc = true
			}
		}

		if !allowedRc {
			// return an error
			detailedErr := ExecError{
				Stdout:   stdoutBuf.String(),
				Stderr:   stderrBuf.String(),
				Rc:       rc,
				Duration: duration,
				Cause:    err,
				Argv:     argv,
			}
			return nil, &detailedErr
		}
	} else {
		rc = 0
	}

	result := ExecResult{
		Stdout:   stdoutBuf.String(),
		Stderr:   stderrBuf.String(),
		Rc:       rc,
		Duration: duration,
	}
	return &result, nil
}
