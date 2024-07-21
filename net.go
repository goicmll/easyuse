package easyuse

import (
	"net"
	"net/http"
)

type SSLCertInfo struct {
	Url       string
	Subject   string
	NotBefore string
	NotAfter  string
}

// GetDomainCertInfo 获取Domain 证书信息
func GetDomainCertInfo(domain string) (*SSLCertInfo, error) {
	resp, err := http.Get(domain)
	if err != nil {
		return nil, err
	}
	sci := SSLCertInfo{}
	certInfo := resp.TLS.PeerCertificates[0]
	sci.NotAfter = certInfo.NotAfter.Format("2006-01-02 15:04:05")
	sci.NotBefore = certInfo.NotBefore.Format("2006-01-02 15:04:05")
	sci.Subject = certInfo.Subject.CommonName
	return &sci, nil
}

// GetDomainResolve 通过域名解析获取IP地址列表。
// 该函数接受一个域名作为参数，返回该域名对应的IP地址列表以及可能出现的错误。
// 如果域名无法解析，函数将返回一个空切片和一个非nil的错误。
func GetDomainResolve(domain string) ([]string, error) {
	// 初始化一个IP地址切片，预分配空间为2，以减少内存分配。
	// 在实际编程中，预分配大小应根据实际情况进行调整。
	var ipss = make([]string, 0, 2)

	// 使用net包的LookupIP函数解析域名，返回一个IP地址切片和一个错误。
	// 如果域名无法解析，err将非nil。
	ips, err := net.LookupIP(domain)
	if err != nil {
		// 如果解析出错，返回空切片和错误。
		return ipss, err
	}

	// 遍历解析得到的IP地址切片，将每个IP地址转换为字符串形式并添加到ipss切片中。
	// 这里使用了append函数来动态添加元素到切片。
	for _, ip := range ips {
		ipss = append(ipss, ip.String())
	}

	// 解析成功，返回IP地址字符串切片和nil错误。
	return ipss, nil
}
