package downloader

import (
	"bufio"
	"github.com/gogodjzhu/gogoscrapy/src/utils"
	"io"
	"os"
	"strconv"
	"strings"
	"sync/atomic"
)

type IProxy interface {
	GetId() int
	GetHost() string
	GetPort() int
	GetUsername() string
	GetPassword() string
}

type Proxy struct {
	Id       int    `json:"id"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func NewProxy(id int, host string, port int, username, password string) Proxy {
	return Proxy{Id: id, Host: host, Port: port, Username: username, Password: password}
}

func (this Proxy) GetId() int {
	return this.Id
}

func (this Proxy) GetHost() string {
	return this.Host
}

func (this Proxy) GetPort() int {
	return this.Port
}

func (this Proxy) GetUsername() string {
	return this.Username
}

func (this Proxy) GetPassword() string {
	return this.Password
}

func (this Proxy) equals(proxy Proxy) bool {
	if this == proxy {
		return true
	}
	if this.Host != proxy.Host {
		return false
	}
	if this.Port != proxy.Port {
		return false
	}
	if this.Username != proxy.Username {
		return false
	}
	if this.Password != proxy.Password {
		return false
	}
	return true
}

type IProxyFactory interface {
	GetProxy() (IProxy, error)
	ReturnProxy(proxy IProxy)
}

// read proxy file and produce Proxy
// line format: {address} {Port}
type FileProxyFactory struct {
	fileUrl    string //file path
	proxyCache []IProxy
	inited     bool
	index      *int32
	proxyQueue *utils.AsyncQueue
}

func NewFileProxyFactory(fileUrl string) (*FileProxyFactory, error) {
	var i int32
	fileProxyFactory := &FileProxyFactory{
		fileUrl:    fileUrl,
		proxyCache: make([]IProxy, 0),
		inited:     false,
		index:      &i,
		proxyQueue: &utils.AsyncQueue{},
	}
	if err := fileProxyFactory.init(); err != nil {
		return nil, err
	}
	return fileProxyFactory, nil
}

//not thread safe
func (this *FileProxyFactory) init() error {
	if this.inited {
		return nil
	}
	fi, err := os.Open(this.fileUrl)
	if err != nil {
		return err
	}
	defer fi.Close()

	br := bufio.NewReader(fi)

	for {
		if line, err := br.ReadString('\n'); err != nil {
			if err == io.EOF {
				break
			}
			return err
		} else {
			line = strings.TrimSpace(line)
			arr := strings.Split(line, "\t")
			host := arr[0]
			port, err := strconv.Atoi(arr[1])
			if err != nil {
				return err
			}
			this.proxyCache = append(this.proxyCache, Proxy{Host: host, Port: port})
		}
	}
	this.inited = true
	return nil
}

func (this *FileProxyFactory) GetProxy() (IProxy, error) {
	i := atomic.AddInt32(this.index, 1)
	cleanIndex := i % int32(len(this.proxyCache))
	return this.proxyCache[cleanIndex], nil
}

func (this *FileProxyFactory) ReturnProxy(proxy IProxy) {
	//FileProxyFactory just reuse the proxy circularly, not need to return proxy.
}
