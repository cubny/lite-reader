package utils

import (
	"fmt"
	"log"
	"net/url"
	"path"
	"strings"
)

func ForceUrl(u string) string {
	u = strings.Trim(u, " ")
	if strings.HasPrefix(u, "feed://") {
		return fmt.Sprintf("http://%s", u[7:])
	}
	for _, proto := range []string{"http://", "https://"} {
		if strings.HasPrefix(u, proto) {
			return u
		}
	}
	return fmt.Sprintf("http://%s", u)
}
func JoinUrl(baseU, suffix string) (string, error) {
	baseUrl, err := url.Parse(baseU)
	if err != nil {
		log.Println("Error joining urls", err)
		return "", err
	}
	baseUrl.Path = path.Join(baseUrl.Path, suffix)
	return baseUrl.String(), nil
}
