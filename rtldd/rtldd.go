package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/spf13/cobra"
	"golang.org/x/sys/unix"
)

var (
	sysroot string
	verify  bool
)

func chdir(dir string) error {
	if dir != "." {
		return unix.Chdir(dir)
	}
	return nil
}

func chroot(root string) error {
	if root != "/" {
		return unix.Chroot(root)
	}
	return nil
}

func run(cmd *exec.Cmd) int {
	var discard bytes.Buffer
	if cmd.Stdin == nil {
		cmd.Stdin = &discard
	}
	if cmd.Stdout == nil {
		cmd.Stdout = &discard
	}
	if cmd.Stderr == nil {
		cmd.Stderr = &discard
	}

	cmd.Run()

	return cmd.ProcessState.ExitCode()
}

func main() {
	rootCmd := &cobra.Command{
		Use:  "rtldd [flags] file",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			dir := filepath.Dir(args[0])

			root, err := filepath.Abs(sysroot)
			if err != nil {
				return err
			}

			file := fmt.Sprintf("./%s", filepath.Base(args[0]))

			if err := chdir(dir); err != nil {
				return err
			}

			if err := chroot(root); err != nil {
				return err
			}

			info, err := os.Stat(file)
			if err != nil {
				return err
			}
			if !info.Mode().IsRegular() {
				return fmt.Errorf("not regular file")
			}

			for _, rtld := range []string{
				"/lib/ld-linux.so.2",
				"/lib64/ld-linux-x86-64.so.2",
				"/libx32/ld-linux-x32.so.2",
			} {
				cmd := exec.Command(rtld, "--verify", file)
				var verifyOut bytes.Buffer
				cmd.Stdout = &verifyOut

				switch run(cmd) {
				case 0, 2:
					if verify {
						rtld, err := filepath.EvalSymlinks(rtld)
						if err != nil {
							return err
						}

						fmt.Println(rtld)
					} else {
						cmd := exec.Command(rtld, file)
						cmd.Env = append(os.Environ(),
							fmt.Sprintf("LD_LIBRARY_VERSION=%s", verifyOut.Bytes()),
							"LD_TRACE_LOADED_OBJECTS=1",
						)
						cmd.Stdout = os.Stdout

						switch run(cmd) {
						case 0:
						default:
							return fmt.Errorf("exited with unknown exit code")
						}
					}

					return nil
				}
			}

			return fmt.Errorf("not a dynamic executable")
		},
	}

	rootCmd.Flags().StringVar(&sysroot, "sysroot", "/", "use string as the location of the sysroot")
	rootCmd.MarkFlagDirname("sysroot")
	rootCmd.Flags().BoolVar(&verify, "verify", false, "verify that file is dynamically linked and a dynamic linker can handle it")

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
