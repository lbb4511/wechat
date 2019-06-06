package main

/*微信公众号本地模拟测试工具也*/
import (
	"encoding/xml"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"
)

const usage = `
Usage: wechatImitateTool [dir|file [filter]]
       wechatImitateTool -h

Example:
  wechatImitateTool url token id id content

  Report bugs to <lbb4511{AT}126.com>.
`

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {
	if len(os.Args) != 6 || os.Args[1] == "-h" {
		fmt.Fprintln(os.Stderr, usage[1:len(usage)-1])
		os.Exit(0)
	}

	url := os.Args[1]
	token := os.Args[2]
	from := os.Args[3]
	to := os.Args[4]
	text := os.Args[5]

	timestamp := time.Now().Unix()
	timestampStr := strconv.FormatInt(timestamp, 10)
	nonce := randStringRunes(8)

	sign := signature(timestampStr, nonce, token)

	url = fmt.Sprintf("%s?signature=%s&timestamp=%s&nonce=%s", url, sign, timestampStr, nonce)

	var message Text
	message.FromUserName = strToCDATA(from)
	message.ToUserName = strToCDATA(to)
	message.MsgType = strToCDATA("text")
	message.MsgId = rand.Int63()
	message.Content = strToCDATA(text)
	message.CreateTime = timestamp

	xml, err := xml.Marshal(message)

	if err != nil {
		fmt.Println(err)
		return
	}

	resp, err := send(url, string(xml))

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("URL:", url)
	fmt.Println("--------------------")
	fmt.Println("Send Message:")

	x, err := formatXML(xml)

	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(string(x))
	}

	fmt.Println("--------------------")
	fmt.Println("Response:")

	if resp == nil {
		fmt.Println("--------------------")
		return
	}

	x, err = formatXML(resp)

	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(string(x))
	}

	fmt.Println("--------------------")
}
