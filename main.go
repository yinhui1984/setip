package main

import (
	"fmt"
	"setip/helpers"
	"strings"
	"time"
)

func getTargetAdapter() string {
	adatpers, err := helpers.ListAllNetWorkAdapters()
	if err != nil {
		helpers.LogFatal(err)
	}
	selectedAdapterIndex, err := helpers.MenuSelect(adatpers, "选择网卡")
	if err != nil {
		helpers.LogFatal(err)
	}

	if selectedAdapterIndex < 0 {
		return ""
	}

	return adatpers[selectedAdapterIndex]
}

func getDhcpType() int {
	options := []string{"自动获取IP (DHCP)", "手动设置"}
	selectedIndex, err := helpers.MenuSelect(options, "选择IP获取方式")
	if err != nil {
		helpers.LogFatal(err)
	}
	return selectedIndex
}

func setAdapterDhcp(adapter string) {
	cmd := `sudo networksetup -setdhcp ` + adapter
	output, err := helpers.RunCommandLine2(cmd)
	if err != nil {
		helpers.LogFatal(err, output)
	}

	checkAdapterStatus(adapter, "")

	helpers.LogOK("设置完成")
	pingRouter(adapter)
}

func setAdapterStaticIP(adapter, ip, netmarsk, router string) {
	info := "将会设置网卡 " + adapter + ": \nIP: " + ip + " 子网掩码: " + netmarsk + " 网关: " + router + " \n是否继续？"
	if !helpers.AskForYes(info) {
		return
	}

	cmd := `sudo networksetup -setmanual ` + adapter + ` ` + ip + ` ` + netmarsk + ` ` + router
	output, err := helpers.RunCommandLine2(cmd)
	if err != nil {
		helpers.LogFatal(err, output)
	}

	checkAdapterStatus(adapter, ip)

	helpers.LogOK("设置完成")

	pingRouter(adapter)
}

func showAdapterInfo(adapter string) {
	cmd := `networksetup -getinfo ` + adapter
	output, err := helpers.RunCommandLine2(cmd)
	if err != nil {
		helpers.LogFatal(err, output)
	}

	fmt.Println(output)
}

func checkAdapterStatus(adapter string, staticIP string) {

	var output string
	var err error
	cmd := `networksetup -getinfo ` + adapter
	fmt.Println(">> " + cmd)
	for {
		output, err = helpers.RunCommandLine2(cmd)
		time.Sleep(100 * time.Millisecond)
		if err != nil {
			helpers.LogFatal(err)
		}

		if (strings.Contains(output, "DHCP Configuration") && staticIP == "") ||
			(strings.Contains(output, "Manual Configuration") && staticIP != "" && strings.Contains(output, staticIP)) {
			fmt.Println(output)
			break
		}

	}

}

func pingRouter(adapter string) {
	// networksetup -getinfo Wi-Fi |grep '^Router' | awk -F'[: ]+' '/Router:/{print $2}'

	time.Sleep(2 * time.Second)

	router, err := helpers.RunCommandLine2(`networksetup -getinfo ` + adapter + ` |grep '^Router' | awk -F'[: ]+' '/Router:/{print $2}'`)
	if err != nil {
		return
	}

	router = strings.TrimSpace(router)
	if router != "" && helpers.CheckIPAddress(router) {
		helpers.LogInfo("正在ping网关: " + router)
		if helpers.Ping(router, 1) {
			helpers.LogOK("ping网关成功")
		} else {
			helpers.LogWarning("ping网关失败")
		}
	}
}

func main() {
	targetAdapter := getTargetAdapter()
	if targetAdapter == "" {
		return
	}

	targetAdapter = "'" + targetAdapter + "'" // 名称中可能有空格，需要加上单引号

	helpers.PrintTitle("网卡: " + targetAdapter)

	showAdapterInfo(targetAdapter)

	dhcpType := getDhcpType()
	if dhcpType < 0 {
		return
	}

	if dhcpType == 0 {
		setAdapterDhcp(targetAdapter)

		return
	} else if dhcpType == 1 {
		fields := []string{"IP Address", "Netmask", "Default Gateway"}
		values, err := helpers.StaticIpInputs(fields)
		if err != nil {
			helpers.LogFatal(err)
		}

		if len(values) != 3 {
			helpers.LogError("输入的参数不正确")
			return
		}

		if values[1] == "" {
			values[1] = "255.255.255.0"
		}

		if helpers.CheckIPAddress(values[0]) && values[2] == "" {
			// for example, ip is 192.168.1.100, but gateway is emtpy, then set gateway to 192.168.1.1
			ip := values[0]
			ip = ip[:strings.LastIndex(ip, ".")]
			values[2] = ip + ".1"
		}

		if !helpers.CheckIPAddress(values[0]) {
			helpers.LogError("IP地址格式不正确")
			return
		}
		if !helpers.CheckIPAddress(values[1]) {
			helpers.LogError("子网掩码格式不正确")
			return
		}
		if !helpers.CheckIPAddress(values[2]) {
			helpers.LogError("网关格式不正确")
			return
		}

		setAdapterStaticIP(targetAdapter, values[0], values[1], values[2])
		return
	}

}
