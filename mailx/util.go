package mailx

import "net/mail"

func plain(addrs []mail.Address) []string {
	res := make([]string, len(addrs))
	for i, v := range addrs {
		res[i] = v.String()
	}
	return res
}
