package clients

import (
	"fmt"
	"regexp"
)

func convert(name string) string {
	var re = regexp.MustCompile(`(?m)^cloud-manager-client-(.*)\.(rpm|deb|exe)\.sum$`)
	return fmt.Sprintf("cloud-manager-%s-%s", getOs(re.FindStringSubmatch(name)[2]), getArch(re.FindStringSubmatch(name)[1]))
}

func getArch(arch string) string {
	switch arch {
	case "x86":
		return "386"
	case "x64":
		return "amd64"
	case "aarch64":
		return "arm64"
	case "arm":
		return "arm"
	}
	return ""
}

func convertArch(arch string) string {
	switch arch {
	case "386":
		return "x86"
	case "amd64":
		return "x64"
	case "arm64":
		return "aarch64"
	case "arm":
		return "arm"
	}
	return ""
}

func getOs(os string) string {
	switch os {
	case "exe":
		return "windows"
	case "deb":
		return "linux"
	case "rpm":
		return "rpm-linux"
	}
	return ""
}

func convertOs(os string) string {
	switch os {
	case "windows":
		return "exe"
	case "linux":
		return "deb"
	case "rpm-linux":
		return "rpm"
	}
	return ""
}
