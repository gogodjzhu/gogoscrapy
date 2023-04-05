package downloader

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/gogodjzhu/gogoscrapy/utils"
	"io"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync/atomic"
)

type IProxy interface {
	GetId() int
	GetTransport() *http.Transport
}

type Proxy struct {
	Id   int    `json:"id"`
	Host string `json:"host"`
	Port int    `json:"port"`
	Type string `json:"type"` //http or https or socks5
}

func NewProxy(id int, host string, port int, typ string) (IProxy, error) {
	// validate type
	if !regexp.MustCompile(`^http|https|socks5$`).MatchString(typ) {
		return nil, errors.New(fmt.Sprintf("invalid proxy type: %s", typ))
	}
	// validate port
	if port < 0 || port > 65535 {
		return nil, errors.New(fmt.Sprintf("invalid proxy port: %d", port))
	}
	return &Proxy{Id: id, Host: host, Port: port, Type: typ}, nil
}

func (p *Proxy) GetId() int {
	return p.Id
}

func (p *Proxy) GetTransport() *http.Transport {
	return &http.Transport{
		Proxy: http.ProxyURL(&url.URL{
			Host: fmt.Sprintf("%s://%s:%d", p.Type, p.Host, p.Port),
		}),
	}
}

type IProxyFactory interface {
	GetProxy() (IProxy, error)
	ReturnProxy(proxy IProxy)
}

// FileHttpProxyFactory read http proxy file and produce Proxy
// line format: {address}:{Port}
type FileHttpProxyFactory struct {
	fileUrl    string //file path
	proxyCache []IProxy
	index      *int32
	proxyQueue *utils.AsyncQueue
}

func NewFileHttpProxyFactory(fileUrl string) (IProxyFactory, error) {
	var i int32
	fileProxyFactory := &FileHttpProxyFactory{
		fileUrl:    fileUrl,
		proxyCache: make([]IProxy, 0),
		index:      &i,
		proxyQueue: &utils.AsyncQueue{},
	}
	if err := fileProxyFactory.init(); err != nil {
		return nil, err
	}
	return fileProxyFactory, nil
}

func (fp *FileHttpProxyFactory) init() error {
	fi, err := os.Open(fp.fileUrl)
	if err != nil {
		return err
	}
	defer fi.Close()

	br := bufio.NewReader(fi)

	for {
		if line, _, err := br.ReadLine(); err != nil {
			if err == io.EOF {
				break
			}
			return err
		} else {
			str := strings.TrimSpace(string(line))
			arr := strings.Split(str, ":")
			host := arr[0]
			port, err := strconv.Atoi(arr[1])
			if err != nil {
				return err
			}
			proxy, err := NewProxy(len(fp.proxyCache), host, port, "http")
			if err != nil {
				return err
			}
			fp.proxyCache = append(fp.proxyCache, proxy)
		}
	}
	return nil
}

func (fp *FileHttpProxyFactory) GetProxy() (IProxy, error) {
	i := atomic.AddInt32(fp.index, 1)
	cleanIndex := i % int32(len(fp.proxyCache))
	return fp.proxyCache[cleanIndex], nil
}

func (fp *FileHttpProxyFactory) ReturnProxy(proxy IProxy) {
	//FileHttpProxyFactory just reuse the proxy circularly, not need to return proxy.
}
