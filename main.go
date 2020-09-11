package main

import (
	"SRUN_LOGIN/jsVM"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

var host string

func httpGet(url string) (resp *http.Response, err error) {
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", `Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/78.0.3904.70 Safari/537.36`)
	req.Header.Set("Accept", `*/*`)
	req.Header.Set("Host", host)
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	req.Header.Set("Accept-Encoding", "gzip, deflate")
	req.Header.Set("Accept-Language", "en,zh-CN;q=0.9,zh;q=0.8,pl;q=0.7")
	fmt.Println(req.Header)
	return http.DefaultClient.Do(req)
}

func timestamp() string {
	return strconv.Itoa(int(time.Now().UnixNano() / 1e6))
}
func hereyouare(body io.ReadCloser) (m map[string]interface{}, err error) {
	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(body)
	if err != nil {
		return
	}
	r := buf.String()
	pre := "hereyouare("
	if !strings.HasPrefix(r, pre) {
		err = errors.New("hereyouare: 不是正确的jsonp格式，可能是网络请求出错")
		return
	}
	r = r[len(pre) : len(r)-1]
	m = make(map[string]interface{})
	err = json.Unmarshal([]byte(r), &m)
	return
}
func getChallenge(uname, ip string) (challenge string, err error) {
	resp, err := httpGet(fmt.Sprintf("http://%v/cgi-bin/get_challenge?callback=hereyouare&username=%v&ip=%v&_=%v",
		host,
		uname,
		ip,
		timestamp()),
	)
	if err != nil {
		return
	}
	m, err := hereyouare(resp.Body)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if m["res"].(string) != "ok" {
		return "", errors.New("get_challenge: res不ok - " + m["error"].(string))
	}
	return m["challenge"].(string), nil
}

func genParams(cha, uname, passwd, acid, ip string) (password, chksum, info string, err error) {
	utils := jsVM.NewUtils()
	base64 := jsVM.NewBase64()
	md5 := jsVM.NewMd5()
	sha1 := jsVM.NewSha1()
	b, _ := json.Marshal(struct {
		Username string `json:"username"`
		Password string `json:"password"`
		IP       string `json:"ip"`
		Acid     string `json:"acid"`
		EncVer   string `json:"enc_ver"`
	}{
		uname,
		passwd,
		ip,
		acid,
		"s" + "run" + "_bx1",
	})
	v, err := utils.Call("xEncode", nil, string(b), cha)
	if err != nil {
		return
	}
	v, err = base64.Call("_encode", nil, v.String())
	if err != nil {
		return
	}
	info = "{SRBX1}" + v.String()
	chkstr := cha + uname
	v, err = md5.Call("md5", nil, passwd, cha)
	if err != nil {
		return
	}
	hmd5 := v.String()
	chkstr += cha + hmd5
	chkstr += cha + acid
	chkstr += cha + ip
	chkstr += cha + "200"
	chkstr += cha + "1"
	chkstr += cha + info
	v, err = sha1.Call("sha1", nil, chkstr)
	if err != nil {
		return
	}
	chksum = v.String()
	password = "{MD5}" + hmd5
	return
}

func getViewParam() (acid, ip string, err error) {
	resp, err := httpGet(fmt.Sprintf("http://%v", host))
	if err != nil {
		return
	}
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	acid, _ = doc.Find("#ac_id").First().Attr("value")
	ip, _ = doc.Find("#user_ip").First().Attr("value")
	return
}

func login(uname, passwd string) (sucMessage string, err error) {
	acid, ip, err := getViewParam()
	if err != nil {
		return
	}
	cha, err := getChallenge(uname, ip)
	if err != nil {
		return
	}
	password, chksum, info, err := genParams(cha, uname, passwd, acid, ip)
	if err != nil {
		return
	}
	u := fmt.Sprintf(`http://%v/cgi-bin/srun_portal?callback=hereyouare&action=login&username=%v&password=%v&ac_id=%v&ip=%v&chksum=%v&info=%v&n=200&type=1&os=Linux&name=Linux&double_stack=0&_=%v`,
		host,
		url.QueryEscape(uname),
		url.QueryEscape(password),
		acid,
		ip,
		chksum,
		url.QueryEscape(info),
		timestamp(),
	)
	resp, err := httpGet(u)
	if err != nil {
		return
	}
	m, err := hereyouare(resp.Body)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if m["res"] != "ok" {
		return "", errors.New("登录失败" + m["error_msg"].(string) + m["error"].(string))
	}
	return m["suc_msg"].(string), nil
}

func main() {
	host = os.Getenv("SRUN_HOST")
	if len(host) <= 0 {
		host = "10.248.98.2"
	}
	uname := os.Getenv("SRUN_UNAME")
	passwd := os.Getenv("SRUN_PASSWD")
	if len(uname) <= 0 || len(passwd) <= 0 {
		log.Fatal("账密输了吗？")
	}
	sucMessage, err := login(uname, passwd)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(sucMessage)
	if keepAlive := os.Getenv("SRUN_KEEP_ALIVE"); keepAlive == "1" || strings.EqualFold(keepAlive, "true") {
		c := http.Client{Timeout: 5 * time.Second}
		// 每180秒一个心跳http包
		tick := time.Tick(180 * time.Second)
		for range tick {
			resp, err := c.Get("http://www.msftconnecttest.com/connecttest.txt")
			if err == nil {
				_ = resp.Body.Close()
			}
		}
	}
}
