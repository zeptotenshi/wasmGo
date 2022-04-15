//+build tinygo wasm,js

package firebase

import (
	"fmt"
	"syscall/js"

	"github.com/zeptotenshi/wasmGo/web"
)

const (
	auth = "auth"

	AUTH__user             = "user"
	AUTH__user_displayName = "displayName"
	AUTH__user_idToken     = "id-token"

	AUTH__error_emailInUse    = "auth/email-already-in-use"
	AUTH__error_invalidEmail  = "auth/invalid-email"
	AUTH__error_weakPassword  = "auth/weak-password"
	AUTH__error_wrongPassword = "auth/wrong-password"
	AUTH__error_userDisabled  = "auth/user-disabled"
	AUTH__error_userNotFound  = "auth/user-not-found"

	function__onAuthStateChanged             = "onAuthStateChanged"
	function__createUserWithEmailAndPassword = "createUserWithEmailAndPassword"
	function__signInWithEmailAndPassword     = "signInWithEmailAndPassword"
	function__signOut                        = "signOut"

	function__user_getIdToken = "getIdToken"
)

type Auth struct {
	value js.Value
	User  js.Value
}

func (a *Auth) SetAuthStateChangedCallback(_cb js.Func) error {
	if err := web.ValidJSValue(auth, a.value); err != nil {
		return fmt.Errorf("[%s] [%s] [AuthStateChanged] [error]: %v", firebase, auth, err)
	}
	a.value.Call(function__onAuthStateChanged, _cb)
	return nil
}

func (a *Auth) CreateUser(_uname, _pword string, _errCB js.Func) error {
	err := web.ValidJSValue(auth, a.value)
	if err != nil {
		return fmt.Errorf("[%s] [%s] [CreateUser] [error]: %v", firebase, auth, err)
	}
	prom := a.value.Call(function__createUserWithEmailAndPassword, _uname, _pword)
	if err = web.ValidJSValue(web.PROMISE, prom); err != nil {
		return fmt.Errorf("[%s] [%s] [CreateUser] [error]: %v", firebase, auth, err)
	}
	prom.Call(web.PROMISE__catch, _errCB)
	return nil
}

func (a *Auth) SignIn(_uname, _pword string, _errCB js.Func) error {
	err := web.ValidJSValue(auth, a.value)
	if err != nil {
		return fmt.Errorf("[%s] [%s] [SignIn] [error]: %v", firebase, auth, err)
	}
	prom := a.value.Call(function__signInWithEmailAndPassword, _uname, _pword)
	if err = web.ValidJSValue(web.PROMISE, prom); err != nil {
		return fmt.Errorf("[%s] [%s] [SignIn] [error]: %v", firebase, auth, err)
	}
	prom.Call(web.PROMISE__catch, _errCB)
	return nil
}

func (a *Auth) SignOut() error {
	if err := web.ValidJSValue(auth, a.value); err != nil {
		return fmt.Errorf("[%s] [%s] [SignOut] [error]: %v", firebase, auth, err)
	}
	a.value.Call(function__signOut)
	return nil
}

func (a *Auth) GetIdToken(_successCB, _errorCB js.Func) error {
	err := web.ValidJSValue(AUTH__user, a.User)
	if err != nil {
		return fmt.Errorf("[%s] [%s] [GetIdToken] [error]: %v", firebase, auth, err)
	}
	if err = web.ValidJSValue(fmt.Sprintf("%s.%s.success-callback", AUTH__user, function__user_getIdToken), js.ValueOf(_successCB)); err != nil {
		return fmt.Errorf("[%s] [%s] [GetIdToken] [error]: %v", firebase, auth, err)
	}

	prom := a.User.Call(function__user_getIdToken, true)
	if err = web.ValidJSValue(fmt.Sprintf("%s.%s", function__user_getIdToken, web.PROMISE), prom); err != nil {
		return fmt.Errorf("[%s] [%s] [GetIdToken] [error]: %v", firebase, auth, err)
	}
	prom.Call(web.PROMISE__then, _successCB)
	if err = web.ValidJSValue(fmt.Sprintf("%s.%s.error-callback", AUTH__user, function__user_getIdToken), js.ValueOf(_errorCB)); err == nil {
		prom.Call(web.PROMISE__catch, _errorCB)
	}

	return nil
}
