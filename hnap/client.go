package hnap

import (
	"net/http"
	"os"
	"strings"
	"github.com/PuerkitoBio/goquery"
	"crypto/md5"
	"crypto/hmac"
	"time"
	"encoding/hex"
	"strconv"
)

const HNAP1_XMLNS = "http://purenetworks.com/HNAP1/";
const HNAP_LOGIN_METHOD = "Login";

func hnapAddress() string {
	return os.Getenv("HNAP_ADDRESS")
}

func hnapUser() string {
	return os.Getenv("HNAP_USER")
}

func hnapPass() string {
	return os.Getenv("HNAP_PASS")
}


type HNAPCredentials struct {
	challenge string
	cookie string
	publicKey string
	privateKey string
}


func login() HNAPCredentials {
	resp1 := request(HNAP_LOGIN_METHOD, requestBody(HNAP_LOGIN_METHOD, loginRequest1), HNAPCredentials{})

	credentials := HNAPCredentials {
		challenge: resp1.Find("Challenge").Text(),
		cookie: resp1.Find("Cookie").Text(),
		publicKey: resp1.Find("PublicKey").Text()}

	credentials.privateKey = hmacGenerate(credentials.publicKey + hnapPass(), credentials.challenge)

	request(HNAP_LOGIN_METHOD, requestBody(HNAP_LOGIN_METHOD, func() string {
		return loginRequest2(credentials.privateKey, credentials.challenge)
	}), credentials)

	return credentials
}

func getHnapAuth(soapAction string, privateKey string) string {
	time_stamp := time.Now().Unix()
	return hmacGenerate(privateKey, string(time_stamp) + soapAction) + " " + string(time_stamp);
}

func hmacGenerate(key string, data string) string {
	mac := hmac.New(md5.New, []byte(key) )
	mac.Write([]byte(data) )
	str := hex.EncodeToString(mac.Sum(nil) )
	return strings.ToUpper(str[0:32])
}

func loginRequest1() string {
	return "<Action>request</Action>" +
		   "<Username>" + hnapUser() + "</Username>" +
		   "<LoginPassword />"	+
		   "<Captcha></Captcha>"
}

func loginRequest2(privateKey string, challenge string) string {
	return "<Action>login</Action>" +
		   "<Username>" + hnapUser() + "</Username>" +
		   "<LoginPassword>" + hmacGenerate(privateKey, challenge) + "</LoginPassword>" +
		   "<Captcha></Captcha>"
}


func requestBody(method string, parameters func() string) string {

	return "<?xml version=\"1.0\" encoding=\"utf-8\"?>" +
		"<soap:Envelope " +
		"xmlns:xsi=\"http://www.w3.org/2001/XMLSchema-instance\" " +
		"xmlns:xsd=\"http://www.w3.org/2001/XMLSchema\" " +
		"xmlns:soap=\"http://schemas.xmlsoap.org/soap/envelope/\">" +
			"<soap:Body>" +
			"<" + method + " xmlns=\"" + HNAP1_XMLNS + "\">" +
				parameters() +
			"</" + method + ">" +
			"</soap:Body>" +
		"</soap:Envelope>"
}

func request(soapAction string, soapBody string, credentials HNAPCredentials) (document *goquery.Document) {

	tr := &http.Transport{
		MaxIdleConns:       10,
		IdleConnTimeout:    30 * time.Second,
		DisableCompression: true,
	}
	client := &http.Client{Transport: tr}

	req, _ := http.NewRequest("POST", hnapAddress(), strings.NewReader(soapBody) )

	req.Header.Add("Content-Type", "text/xml; charset=utf-8")
	req.Header.Add("SOAPAction", "\"" + HNAP1_XMLNS + soapAction + "\"")

	// For all methods, except for login, we need an active session
	// TODO : Improve such that we don't have to login everytime
	if len(credentials.challenge) > 0 {
		hnapAuth := getHnapAuth("\"" + HNAP1_XMLNS + soapAction + "\"", credentials.privateKey)
		req.Header.Add("HNAP_AUTH", hnapAuth)
		req.Header.Add("Cookie", "uid=" + credentials.cookie)
	}

	resp, _ := client.Do(req)

	document, _ = goquery.NewDocumentFromReader(resp.Body)
	resp.Body.Close()
	return
}

func moduleParameters(module string) string {
	return "<ModuleID>" + module + "</ModuleID>";
}

func controlParameters(module string, status string) string {
	return moduleParameters(module) +
	"<NickName>Socket 1</NickName><Description>Socket 1</Description>" +
	"<OPStatus>" + status + "</OPStatus><Controller>1</Controller>"
}

func On() {
	params := login()
	body := requestBody("SetSocketSettings", func() string {
		return controlParameters("1", "true")
	})

	request("SetSocketSettings", body, params);
};

func Off() {
	params := login()
	body := requestBody("SetSocketSettings", func() string {
		return controlParameters("1", "false")
	})

	request("SetSocketSettings", body, params);
};


func State() bool {
	params := login()
	body := requestBody("GetSocketSettings", func() string {
		return moduleParameters("1")
	})

	doc := request("GetSocketSettings", body, params);
	return doc.Find("OPStatus").Text() == "true"
}

func Consumption() (res float64) {
	params := login()
	body := requestBody("GetCurrentPowerConsumption", func() string {
		return moduleParameters("1")
	})

	doc := request("GetCurrentPowerConsumption", body, params);
	res, _ = strconv.ParseFloat(doc.Find("CurrentConsumption").Text(), 64 )
	return
}
