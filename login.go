package main

import (
	"github.com/keybase/go-libkb"
	"github.com/keybase/protocol/go"
	"github.com/maxtaco/go-framed-msgpack-rpc/rpc2"
)

type LoginHandler struct {
	xp         *rpc2.Transport
	cli        *rpc2.Client
	loginCli   *keybase_1.LoginUiClient
	identifyUi libkb.IdentifyUI
}

func (h *LoginHandler) getRpcClient() *rpc2.Client {
	if h.cli == nil {
		h.cli = rpc2.NewClient(h.xp, libkb.UnwrapError)
	}
	return h.cli
}

func (h *LoginHandler) getLoginUiCli() *keybase_1.LoginUiClient {
	if h.loginCli == nil {
		h.loginCli = &keybase_1.LoginUiClient{h.getRpcClient()}
	}
	return h.loginCli
}

type LoginUI struct {
	cli *keybase_1.LoginUiClient
}

type IdentifyUI struct {
	sessionId int
	cli       *keybase_1.IdentifyUiClient
}

func (h *LoginHandler) getLoginUi() libkb.LoginUI {
	return &LoginUI{h.getLoginUiCli()}
}

func (h *LoginHandler) getIdentifyUi() libkb.IdentifyUI {
	if h.identifyUi == nil {
		h.identifyUi = NextRemoteSelfIdentifyUI(h.getRpcClient())
	}
	return h.identifyUi
}

func (u *LoginUI) GetEmailOrUsername() (ret string, err error) {
	return u.cli.GetEmailOrUsername()
}

func (u *LoginUI) GetKeybasePassphrase(username string, retry string) (string, error) {
	arg := keybase_1.GetKeybasePassphraseArg{Username : username, Retry :retry}
	return u.cli.GetKeybasePassphrase(arg)
}

func (h *LoginHandler) Logout() error {
	return G.LoginState.Logout()
}

func (h *LoginHandler) PassphraseLogin(doIdentify bool) error {
	loginui := h.getLoginUi()
	var idui libkb.IdentifyUI
	if doIdentify {
		idui = h.getIdentifyUi()
	}
	return libkb.LoginAndIdentify(loginui, idui)
}

func (h *LoginHandler) PassphraseLoginNoIdentify(arg keybase_1.PassphraseLoginNoIdentifyArg) error {
	largs := libkb.LoginArg{ Username : arg.Username, Passphrase : arg.Passphrase, NoUi : true }
	return libkb.Login(largs)
}

func (h *LoginHandler) PubkeyLogin() error {
	return nil
}

func (h *LoginHandler) SwitchUser(username string) error {
	return nil
}
