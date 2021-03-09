package https_expire

import (
	"crypto/tls"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"
	"zabbix.com/pkg/conf"
    "zabbix.com/pkg/plugin"
	"zabbix.com/pkg/zbxerr"
)

const (
	paramnums  = 1
	pluginName="Https_expire"
)

var (
	impl Plugin
	opts Options
	err error
)
type Plugin struct {
	plugin.Base
	options Options
	client  *client
}


type Options struct {
	Timeout  int    `conf:"optional,range=1:30"`
}

type client struct {
	client http.Client
}
func newClient(timeout int) *client {
	transport:=&http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify:true},
	}

	client := client{}
	client.client = http.Client{
		Transport: transport,
		Timeout:   time.Duration(timeout) * time.Second,
	}

	return &client
}


// Configure implements the Configurator interface.
// Initializes configuration structures.
func (p *Plugin) Configure(global *plugin.GlobalOptions, options interface{}) {
	if err = conf.Unmarshal(options, &p.options); err != nil {
		p.Errf("cannot unmarshal configuration options: %s", err)
	}

	if p.options.Timeout == 0 {
		p.options.Timeout = global.Timeout
	}
	p.client = newClient(p.options.Timeout)

}

func (p *Plugin) Validate(options interface{}) error {

	return conf.Unmarshal(options, &opts)
}

func checkParamnums(params []string) error {
	if len(params) > paramnums {
		err:=errors.New("Too many parameters.")
		return zbxerr.ErrorTooFewParameters.Wrap(err)
	} else if len(params) ==0 {
		err:=errors.New("Missing URL parameters.")
		return zbxerr.ErrorTooFewParameters.Wrap(err)
	}
	return nil
}

func checkParams(params []string) (string, error) {
	if strings.HasPrefix(params[0], "http://") {
		errorsting:=fmt.Sprintf("Target is using http scheme: %s", params[0])
		err:=errors.New(errorsting)
		return "",zbxerr.ErrorInvalidParams.Wrap(err)
	}

	if !strings.HasPrefix(params[0], "https://") {
		params[0] = "https://" + params[0]
	}
	return string(params[0]),nil
}
func (cli *client) Query(url string) (int64, error) {
	resp, err := cli.client.Get(url)
	if err != nil {
		impl.Debugf("cannot fetch data: %s", err)
		err=errors.New("cannot fetch data")
		return 0, zbxerr.ErrorCannotFetchData.Wrap(err)
	}
	defer resp.Body.Close()
	certInfo:=resp.TLS.PeerCertificates[0]
	expiredays:=(certInfo.NotAfter.Unix()-time.Now().Unix())/60/60/24
	return expiredays,nil
}

// Export implements the Exporter interface.
func (p *Plugin) Export(key string, params []string, ctx plugin.ContextProvider) (interface{}, error) {
	if err = checkParamnums(params); err != nil {
		return nil, err
	}
	urls,err:= checkParams(params)
	if err!= nil {
		return nil,err
	}
	body, err := p.client.Query(urls)
	if err!=nil{
		return nil, err
	}
	return body,nil

}
func init() {
	plugin.RegisterMetrics(&impl, pluginName,
		"https_expire", "Returns the number of days between the HTTPS certificate expiration time and the current date.")
}


