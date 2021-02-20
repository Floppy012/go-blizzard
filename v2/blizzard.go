// Package blizzard is a client library designed to make calling and processing Blizzard Game APIs simple
package blizzard

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
)

// For testing
var c *Client

// Client regional API URLs, locale, client ID, client secret
type Client struct {
	client                                          *http.Client
	cfg                                             clientcredentials.Config
	authorizedCfg                                   oauth2.Config
	oauth                                           OAuth
	oauthHost                                       string
	apiHost                                         string
	dynamicNamespace, staticNamespace               string
	profileNamespace                                string
	dynamicClassicNamespace, staticClassicNamespace string
	region                                          Region
	locale                                          Locale
}

// Region type
type Region int

// Region constants (1=US, 2=EU, 3=KO and TW, 5=CN) DO NOT REARRANGE
const (
	_ Region = iota
	US
	EU
	KR
	TW
	CN
)

func (region Region) String() string {
	var rr = []string{
		"",
		"us",
		"eu",
		"kr",
		"tw",
		"cn",
	}

	return rr[region]
}

// Locale generic locale string
// enUS, esMX, ptBR, enGB, esES, frFR, ruRU, deDE, ptPT, itIT, koKR, zhTW, zhCN
type Locale string

func (locale Locale) String() string {
	return string(locale)
}

// Locale constants
const (
	DeDE = Locale("de_DE")
	EnUS = Locale("en_US")
	EsES = Locale("es_ES")
	EsMX = Locale("es_MX")
	FrFR = Locale("fr_FR")
	ItIT = Locale("it_IT")
	JaJP = Locale("ja_JP")
	KoKR = Locale("ko_KR")
	PlPL = Locale("pl_PL")
	PtBR = Locale("pt_BR")
	RuRU = Locale("ru_RU")
	ThTH = Locale("th_TH")
	ZhCN = Locale("zh_CN")
	ZhTW = Locale("zh_TW")
)

// NewClient create new Blizzard structure. This structure will be used to acquire your access token and make API calls.
func NewClient(clientID, clientSecret string, region Region, locale Locale) *Client {
	var c = Client{
		oauth: OAuth{
			ClientID:     clientID,
			ClientSecret: clientSecret,
		},
		locale: locale,
	}

	c.cfg = clientcredentials.Config{
		ClientID:     c.oauth.ClientID,
		ClientSecret: c.oauth.ClientSecret,
	}

	c.SetRegion(region)

	return &c
}

// GetLocale returns the Locale of the client
func (c *Client) GetLocale() Locale {
	return c.locale
}

// SetLocale changes the Locale of the client
func (c *Client) SetLocale(locale Locale) {
	c.locale = locale
}

// GetRegion returns the Region of the client
func (c *Client) GetRegion() Region {
	return c.region
}

// SetRegion changes the Region of the client
func (c *Client) SetRegion(region Region) {
	c.region = region

	switch region {
	case CN:
		c.oauthHost = "https://www.battlenet.com.cn"
		c.apiHost = "https://gateway.battlenet.com.cn"
		c.dynamicNamespace = "dynamic-zh"
		c.dynamicClassicNamespace = "dynamic-classic-zh"
		c.profileNamespace = "profile-zh"
		c.staticNamespace = "static-zh"
		c.staticClassicNamespace = "static-classic-zh"
	default:
		c.oauthHost = fmt.Sprintf("https://%s.battle.net", region)
		c.apiHost = fmt.Sprintf("https://%s.api.blizzard.com", region)
		c.dynamicNamespace = fmt.Sprintf("dynamic-%s", region)
		c.dynamicClassicNamespace = fmt.Sprintf("dynamic-classic-%s", region)
		c.profileNamespace = fmt.Sprintf("profile-%s", region)
		c.staticNamespace = fmt.Sprintf("static-%s", region)
		c.staticClassicNamespace = fmt.Sprintf("static-classic-%s", region)
	}

	c.cfg.TokenURL = c.oauthHost + "/oauth/token"
	c.client = c.cfg.Client(context.Background())
}

// GetRegion returns the Region of the client
func (c *Client) GetOAuthHost() string {
	return c.oauthHost
}

// GetRegion returns the Region of the client
func (c *Client) GetAPIHost() string {
	return c.apiHost
}

// GetRegion returns the Region of the client
func (c *Client) GetDynamicNamespace() string {
	return c.dynamicNamespace
}

// GetRegion returns the Region of the client
func (c *Client) GetDynamicClassicNamespace() string {
	return c.dynamicClassicNamespace
}

// GetRegion returns the Region of the client
func (c *Client) GetProfileNamespace() string {
	return c.profileNamespace
}

// GetRegion returns the Region of the client
func (c *Client) GetStaticNamespace() string {
	return c.staticNamespace
}

// GetRegion returns the Region of the client
func (c *Client) GetStaticClassicNamespace() string {
	return c.staticClassicNamespace
}

// getStructData processes simple GET request based on pathAndQuery an returns the structured data.
func (c *Client) getStructData(ctx context.Context, pathAndQuery, namespace string, dat interface{}) (interface{}, []byte, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", c.apiHost+pathAndQuery, nil)
	if err != nil {
		return dat, nil, err
	}

	req.Header.Set("Accept", "application/json")

	q := req.URL.Query()
	q.Set("locale", c.locale.String())
	req.URL.RawQuery = q.Encode()

	if namespace != "" {
		req.Header.Set("Battlenet-Namespace", namespace)
	}

	res, err := c.client.Do(req)
	if err != nil {
		return dat, nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return dat, nil, err
	}

	if res.StatusCode != http.StatusOK {
		return dat, body, errors.New(res.Status)
	}

	err = json.Unmarshal(body, &dat)
	if err != nil {
		return dat, body, err
	}

	return dat, body, nil
}

// getStructDataNoLocale processes simple GET request based on pathAndQuery an returns the structured data.
// Does not use a Locale.
func (c *Client) getStructDataNoLocale(ctx context.Context, pathAndQuery, namespace string, dat interface{}) (interface{}, []byte, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", c.apiHost+pathAndQuery, nil)
	if err != nil {
		return dat, nil, err
	}

	req.Header.Set("Accept", "application/json")

	if namespace != "" {
		req.Header.Set("Battlenet-Namespace", namespace)
	}

	res, err := c.client.Do(req)
	if err != nil {
		return dat, nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return dat, nil, err
	}

	if res.StatusCode != http.StatusOK {
		return dat, body, errors.New(res.Status)
	}

	err = json.Unmarshal(body, &dat)
	if err != nil {
		return dat, body, err
	}

	return dat, body, nil
}

// getStructDataOAuth processes simple GET request based on pathAndQuery an returns the structured data.
// Uses OAuth2.
func (c *Client) getStructDataOAuth(ctx context.Context, pathAndQuery, namespace string,
	token *oauth2.Token, dat interface{}) (interface{}, []byte, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", c.apiHost+pathAndQuery, nil)
	if err != nil {
		return dat, nil, err
	}

	req.Header.Set("Accept", "application/json")

	q := req.URL.Query()
	q.Set("locale", c.locale.String())
	req.URL.RawQuery = q.Encode()

	if namespace != "" {
		req.Header.Set("Battlenet-Namespace", namespace)
	}

	client := c.authorizedCfg.Client(context.Background(), token)

	res, err := client.Do(req)
	if err != nil {
		return dat, nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return dat, nil, err
	}

	if res.StatusCode != http.StatusOK {
		return dat, body, errors.New(res.Status)
	}

	err = json.Unmarshal(body, &dat)
	if err != nil {
		return dat, body, err
	}

	return dat, body, nil
}

func formatAccount(account string) string {
	return strings.Replace(account, "#", "-", 1)
}
