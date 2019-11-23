package http

import (
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"gopkg.in/oauth2.v3"
	"gopkg.in/oauth2.v3/generates"
	"gopkg.in/oauth2.v3/manage"
	"gopkg.in/oauth2.v3/models"
	"gopkg.in/oauth2.v3/server"
	"gopkg.in/oauth2.v3/store"
)

//OAuth2 desc
//@struct OAuth2 desc:
type OAuth2 struct {
	_m *manage.Manager
	_c *store.ClientStore
	_s *server.Server

	_tokenExp          int
	_refreshTokenExp   int
	_isGenerateRefresh bool
	_s256KeyValue      string
	_accessURI         string
}

//SetTokenExp desc
//@method SetTokenExp desc: Setting oauth2 token exp time
//@param (int) token exp time sec
func (slf *OAuth2) SetTokenExp(v int) {
	slf._tokenExp = v
}

//SetRefreshTokenExp desc
//@method SetRefreshTokenExp desc: Setting oauth2 refresh token exp time
//@param (int) refresh token exp time sec
func (slf *OAuth2) SetRefreshTokenExp(v int) {
	slf._refreshTokenExp = v
}

//SetRefresh desc
//@method SetRefresh desc: Setting oauth2 token reset refresh
//@param (bool) is refresh
func (slf *OAuth2) SetRefresh(v bool) {
	slf._isGenerateRefresh = v
}

//SetKey desc
//@method SetKey desc: Setting oauth2 token s256 key
//@param (string) key
func (slf *OAuth2) SetKey(v string) {
	slf._s256KeyValue = v
}

//SetURI desc
//@method SetURI desc: Setting oauth2 access address
//@param (string) oauth2 access address
func (slf *OAuth2) SetURI(v string) {
	slf._accessURI = v
}

//Initial desc
//@method Initial desc:
func (slf *OAuth2) Initial() {
	codeTokenCfg := &manage.Config{
		AccessTokenExp:    time.Second * time.Duration(slf._tokenExp),
		RefreshTokenExp:   time.Second * time.Duration(slf._refreshTokenExp),
		IsGenerateRefresh: slf._isGenerateRefresh,
	}

	clientTokenCfg := &manage.Config{AccessTokenExp: time.Second * time.Duration(slf._tokenExp)}

	slf._m = manage.NewDefaultManager()
	slf._m.SetAuthorizeCodeTokenCfg(codeTokenCfg)
	slf._m.SetClientTokenCfg(clientTokenCfg)

	if slf._s256KeyValue != "" {
		slf._m.MapAccessGenerate(generates.NewJWTAccessGenerate([]byte(slf._s256KeyValue), jwt.SigningMethodHS256))
	}

	slf._c = store.NewClientStore()
	slf._m.MapClientStorage(slf._c)

	slf._s = server.NewDefaultServer(slf._m)

	slf._s.SetAllowGetAccessRequest(true)
	slf._s.SetAllowedGrantType(oauth2.ClientCredentials, oauth2.Refreshing)
	slf._s.SetClientInfoHandler(server.ClientFormHandler)
	slf._m.SetRefreshTokenCfg(manage.DefaultRefreshTokenCfg)
}

//onRequestAuth desc
//@method onRequestAuth desc: http request method
//@param (http.ResponseWriter)
//@param (http.Request)
func (slf *OAuth2) onRequestAuth(w http.ResponseWriter, r *http.Request) {
	err := slf._s.HandleTokenRequest(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}

//AddAuthClient desc
//@method AddAuthClient desc: append authorization client
//@param (string) client id
//@param (string) secret
//@param (string) domain
//@param (string) user id
func (slf *OAuth2) AddAuthClient(id string, secret string, domain string, userid string) {
	slf._c.Set(id, &models.Client{ID: id, Secret: secret, Domain: domain, UserID: userid})
}

//ValidateToken desc
func ValidateToken(fun func(w http.ResponseWriter, r *http.Request),
	f interface{},
	srv *server.Server) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		_, err := srv.ValidationBearerToken(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}
	//TODO: 找到目标方法并回传
	return nil
}
