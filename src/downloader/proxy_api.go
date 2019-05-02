package downloader

import (
	"gogoproxypool/src"
)

type ApiProxyFactory struct {
	proxyMapper *src.ProxyMapper
}

func NewApiProxyFactory(config *src.MySqlConfig) (*ApiProxyFactory, error) {
	persist, err := src.NewPersistence(config)
	if err != nil {
		return nil, err
	}
	return &ApiProxyFactory{
		proxyMapper: src.NewProxyMapper(persist),
	}, nil
}

func (this *ApiProxyFactory) GetProxy() (IProxy, error) {
	proxy, err := this.proxyMapper.Get()
	if err != nil {
		return nil, err
	}
	return NewProxy(proxy.Id, proxy.Host, proxy.Port, proxy.Username, proxy.Password), nil
}

func (this *ApiProxyFactory) ReturnProxy(proxy IProxy) {
	if err := this.proxyMapper.ReturnCache(proxy.GetId()); err != nil {
		LOG.Errorf("failed to return proxy, err:%+v", err)
	}
}
