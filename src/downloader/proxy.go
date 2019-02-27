package downloader

import (
	"bufio"
	"io"
	"os"
	"strconv"
	"strings"
	"sunteng/commons/util"
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
	id       int
	host     string
	port     int
	username string
	password string
}

func NewProxy(id int, host string, port int, username, password string) Proxy {
	return Proxy{id: id, host: host, port: port, username: username, password: password}
}

func (this Proxy) GetId() int {
	return this.id
}

func (this Proxy) GetHost() string {
	return this.host
}

func (this Proxy) GetPort() int {
	return this.port
}

func (this Proxy) GetUsername() string {
	return this.username
}

func (this Proxy) GetPassword() string {
	return this.password
}

func (this Proxy) equals(proxy Proxy) bool {
	if this == proxy {
		return true
	}
	if this.host != proxy.host {
		return false
	}
	if this.port != proxy.port {
		return false
	}
	if this.username != proxy.username {
		return false
	}
	if this.password != proxy.password {
		return false
	}
	return true
}

type IProxyFactory interface {
	GetProxy() (IProxy, error)
	ReturnProxy(proxy IProxy)
}

// read proxy file and produce Proxy
// line format: {address} {port}
type FileProxyFactory struct {
	fileUrl    string //file path
	proxyCache []IProxy
	inited     bool
	index      *int32
	proxyQueue *util.AsyncQueue
}

func NewFileProxyFactory(fileUrl string) (*FileProxyFactory, error) {
	var i int32
	fileProxyFactory := &FileProxyFactory{
		fileUrl:    fileUrl,
		proxyCache: make([]IProxy, 0),
		inited:     false,
		index:      &i,
		proxyQueue: &util.AsyncQueue{},
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
			this.proxyCache = append(this.proxyCache, Proxy{host: host, port: port})
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
