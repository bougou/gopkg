- `func WaitTimeout(c *exec.Cmd, timeout time.Duration) error`
  - `WaitTimeout` 等待 exec.Cmd 结束，超时后 Kill 掉 cmd
  - 如果直接使用 `WaitTimeout` 函数，需要在调用 `WaitTimeout` 之前调用 `c.Start()` 将 cmd 启动， `WaitTimeout` 不负责启动 cmd
  - 大部分时候，应用程序应该调用下面提供的基于 `WaitTimeout` 之上的函数

## Based on `WaitTimeout`

- `func RunTimeout(c *exec.Cmd, timeout time.Duration) error` 不返回输出
- `func CombinedOutputTimeout(c *exec.Cmd, timeout time.Duration) ([]byte, error)` 返回混合输出（不区分 stdout，stderr）
- `func StdOutputTimeout(c *exec.Cmd, timeout time.Duration) ([]byte, error)` 只返回标准输出
- `func SeparatedOutputTimeout(c *exec.Cmd, timeout time.Duration) (stdout []byte, stderr []byte, err error)` 分开返回标准输出和错误输出

## ShellCommand

- `func (sh *ShellCommand) SimpleExec() error` 执行 shell 命令，使用系统的标准输入，标准输出和标准错误输出。
  效果类似于：直接在 shell 终端中敲入命令并执行。使用该函数无法获取命令的输出。
- `func (sh *ShellCommand) Exec() ([]byte, error)` 执行 shell 命令，并返回命令的输出内容（不区分 stdout, stderr）
- `func RunShellCommand(command string, silent bool, trim bool) (output []byte, err error)`
- `func RunShellCommandTimeout(command string, silent bool, trim bool, timeout time.Duration) (output []byte, err error)`

## Base on

- `func RunCmd(timeout time.Duration, sudo bool, command string, args ...string) ([]byte, error)`
