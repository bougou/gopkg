package filedirs

import "os"

// Returns true if:
//   - item exists
//   - and the item is not directory
func FileExists(item string) (bool, error) {
	info, err := os.Stat(item)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}

	// item exists
	return !info.IsDir(), nil
}

// Returns true if:
//   - no item exists
//   - or the item is directoy
func FileNotExists(item string) (bool, error) {
	info, err := os.Stat(item)
	if err != nil {
		if os.IsNotExist(err) {
			return true, nil
		}
		return false, err
	}

	// item exists
	return info.IsDir(), nil
}

// Returns true if:
//   - item exists
//   - and the item is directory
func DirExist(item string) (bool, error) {
	info, err := os.Stat(item)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}

	// item exists
	return info.IsDir(), nil
}

// Returns true if:
//   - item does not exists
//   - or the itme is not directory
func DirNotExists(item string) (bool, error) {
	info, err := os.Stat(item)
	if err != nil {
		if os.IsNotExist(err) {
			return true, nil
		}
		return false, err
	}

	// item exists
	return !info.IsDir(), nil
}

// Returns true if:
//   - item exists
func FirOrDirExists(item string) (bool, error) {
	_, err := os.Stat(item)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}

	// item exists
	return true, nil
}

// Returns true if:
//   - item does not exists
func FileOrDirNotExists(item string) (bool, error) {
	_, err := os.Stat(item)
	if err != nil {
		if os.IsNotExist(err) {
			return true, nil
		}
		return false, nil
	}

	// item exists
	return false, nil
}
