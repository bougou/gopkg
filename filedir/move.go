package filedir

type CrossFileSystemBehavior string

const (
	CrossFileSystemBehaviorError    CrossFileSystemBehavior = "error"
	CrossFileSystemBehaviorContinue CrossFileSystemBehavior = "continue"
	CrossFileSystemBehaviorNoop     CrossFileSystemBehavior = "noop"
)

// valid src and dest may be:
//
//	somefile -> somefile2 (√)
//	  eg: /opt/file.txt -> /usr/local/file2.txt
//
//	somedir -> somedir2 (√)
//	  eg:             /opt/somedir       -> /usr/local/somedir2
//	  result is:      /opt/somedir/{...} -> /usr/local/somedir2/{...}
//	  result is not:  /opt/somedir/{...} -> /usr/local/somedir2/somedir/{...}
//
//	somefile -> somedir (√)
//	  eg:         /opt/file.txt -> /usr/local
//	  result is:  /opt/file.txt -> /usr/local/file.txt
//
// invalid
//
//	somedir -> somefile (X)
func Rename(src string, dest string, crossFileSystemBehavior CrossFileSystemBehavior) error {

}
