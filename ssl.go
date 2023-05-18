package habits

import (
	"net/http"
)

type SSLCertInfo struct {
	Url       string
	Subject   string
	NotBefore string
	NotAfter  string
}

// PickHttpsCertInfo 获取https 证书信息
func PickHttpsCertInfo(url string) (*SSLCertInfo, error) {
	resp, err := http.Get(url)
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
