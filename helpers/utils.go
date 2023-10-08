package helpers

import (
	"fmt"
	"net"
	"os/exec"
	"regexp"
	inputstui "setip/helpers/inputs_tui"
	menutui "setip/helpers/menu_tui"
	"strings"
)

// 其支持通配符
func RunCommandLine2(cmd string) (string, error) {
	out, err := exec.Command("bash", "-c", cmd).CombinedOutput()
	return string(out), err
}

func Ping(ip string, count int) bool {

	ok := false
	err := error(nil)
	output := ""
	for i := 0; i < count; i++ {
		output, err = RunCommandLine2("ping -c " + "1" + " " + ip)
		ok = ok || err == nil
	}
	fmt.Println(output)
	return ok

	// pinger 需要sudo
	// pinger, err := ping.NewPinger(ip)
	// if err != nil {
	// 	LogError("Error creating pinger: ", err)
	// 	return false
	// }

	// pinger.Count = count
	// pinger.Timeout = time.Second * 2
	// pinger.SetPrivileged(true)

	// err = pinger.Run()
	// if err != nil {
	// 	LogError("Error running pinger: ", err)
	// 	return false
	// }

	// stats := pinger.Statistics()
	// if stats.PacketsRecv > 0 {
	// 	return true
	// } else {
	// 	return false
	// }
}

func ListAllNetWorkAdapters() ([]string, error) {
	cmd := `networksetup -listallnetworkservices`
	output, err := exec.Command("bash", "-c", cmd).CombinedOutput()
	if err != nil {
		return nil, err
	}

	list := strings.Split(string(output), "\n")

	var adapters []string
	for i := 1; i < len(list)-1; i++ {
		adapters = append(adapters, list[i])
	}
	return adapters, nil
}

func MenuSelect(options []string, title string) (int, error) {

	menutui.MakeMenuList(options, title)
	return menutui.LastChoiseIndex, nil
}

func StaticIpInputs(fields []string) ([]string, error) {
	err := inputstui.MakeInputs(fields)
	return inputstui.LastFielsValues, err
}

func PrintTitle(msg string) {
	box := fmt.Sprintf("%s%s%s", strings.Repeat(" ", (len(msg)-2)/2), msg, strings.Repeat(" ", (len(msg)-2)/2))
	fmt.Println("+" + strings.Repeat("-", len(box)-2) + "+")
	fmt.Println("|" + box + "|")
	fmt.Println("+" + strings.Repeat("-", len(box)-2) + "+")
}

func CheckIPAddress(ip string) bool {
	if parsedIP := net.ParseIP(ip); parsedIP == nil {
		return false
	}
	return true
}

func AskForYes(title string) bool {
	var input string
	fmt.Print(title + "(Y/n):")
	fmt.Scanln(&input)
	input = strings.ToUpper(input)
	return (input != "N" && input != "NO")
}

func ParseIpFromString(data string) string {
	// (Router: 192.168.1.1)  //"(DNS Server: 192.168.1.1)"
	ipRegex := regexp.MustCompile(`\b\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}\b`)
	ip := ipRegex.FindString(data)
	return ip
}
