package osinhander

import (
	"net/http"

	"github.com/RangelReale/osin"
	ex "github.com/RangelReale/osin/example"
)

type OAuthServer struct {
	ser *osin.Server
}

func (s *OAuthServer) AuthApi(w http.ResponseWriter, r *http.Request) {
	resp := s.ser.NewResponse()
	defer resp.Close()

	if ar := s.ser.HandleAuthorizeRequest(resp, r); ar != nil {

		// HANDLE LOGIN PAGE HERE

		ar.Authorized = true
		s.ser.FinishAuthorizeRequest(resp, r, ar)
	}
	osin.OutputJSON(resp, w, r)
}

func (s *OAuthServer) TokenApi(w http.ResponseWriter, r *http.Request) {
	resp := s.ser.NewResponse()
	defer resp.Close()

	if ar := s.ser.HandleAccessRequest(resp, r); ar != nil {
		ar.Authorized = true
		s.ser.FinishAccessRequest(resp, r, ar)
	}
	osin.OutputJSON(resp, w, r)
}

func NewOsinSer() (s *OAuthServer) {
	s = new(OAuthServer)

	sconfig := osin.NewServerConfig()
	sconfig.AccessExpiration = 24 * 3600 // token过期时间改为24h
	sconfig.AllowedAuthorizeTypes = osin.AllowedAuthorizeType{osin.CODE, osin.TOKEN}
	sconfig.AllowedAccessTypes = osin.AllowedAccessType{osin.AUTHORIZATION_CODE,
		osin.REFRESH_TOKEN, osin.PASSWORD, osin.CLIENT_CREDENTIALS, osin.ASSERTION}
	sconfig.AllowGetAccessRequest = true
	sconfig.AllowClientSecretInParams = true
	sconfig.RedirectUriSeparator = " "

	s.ser = osin.NewServer(sconfig, ex.NewTestStorage())
	return s
}
