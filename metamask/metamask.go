//+build tinygo wasm,js

package metamask

import (
	"fmt"
	"syscall/js"

	"github.com/zeptotenshi/wasmGo/web"
)

type ChainID string
type ProviderEvent int

const (
	Mainnet ChainID = "0x1"  //	1  Ethereum Main Network (Mainnet)
	Ropsten         = "0x3"  //	3  Ropsten Test Network
	Rinkeby         = "0x4"  //	4  Rinkeby Test Network
	Goerli          = "0x5"  //	5  Goerli Test Network
	Kovan           = "0x2a" // 42 Kovan Test Network

	Connect ProviderEvent = iota
	Disconnect
	AccountsChanged
	ChainChanged
	Message

	METAMASK  = "metamask"
	ethereum  = "ethereum"
	chain__id = "chainId"

	function__isConnected = "isConnected"
	function__on          = "on"
	function__request     = "request"

	request__method         = "method"
	method__requestAccounts = "eth_requestAccounts"

	error__message = "message"
	error__code    = "code"
)

var (
	chain__names = map[ChainID]string{
		Mainnet: "mainnet",
		Ropsten: "ropsten",
		Rinkeby: "rinkeby",
		Goerli:  "goerli",
		Kovan:   "kovan",
	}

	event__names = map[ProviderEvent]string{
		Connect:         "connect",
		Disconnect:      "disconnect",
		AccountsChanged: "accountsChanged",
		ChainChanged:    "chainChanged",
		Message:         "message",
	}
)

type MetaMask struct {
	win      *web.Window
	provider js.Value

	address string

	ConnectEl *web.Element
	AddressEl *web.Element

	callbacks map[string]js.Func
}

func NewMetaMask(_win *web.Window) *MetaMask {
	mm := &MetaMask{
		provider: js.ValueOf(nil),
		win:      _win,

		callbacks: map[string]js.Func{},
	}

	eth, err := _win.GetGlobal(ethereum)
	if err != nil {
		_win.Error(fmt.Errorf("[MetaMask] [new] [error]: %v", err))
		return mm
	}
	mm.provider = eth

	return mm
}

func (mm *MetaMask) Init() {
	for pe, pen := range event__names {
		switch pe {
		case Connect:
			mm.provider.Call(function__on, pen, js.FuncOf(mm.connect))
		case Disconnect:
			mm.provider.Call(function__on, pen, js.FuncOf(mm.disconnect))
		case AccountsChanged:
			mm.provider.Call(function__on, pen, js.FuncOf(mm.accountsChanged))
		case ChainChanged:
			mm.provider.Call(function__on, pen, js.FuncOf(mm.chainChanged))
		case Message:
			mm.provider.Call(function__on, pen, js.FuncOf(mm.message))
		default:
			continue
		}
	}
}

func (mm *MetaMask) IsConnected() bool {
	tv := mm.provider.Call(function__isConnected)
	if err := web.ValidJSValue(function__isConnected, tv); err != nil {
		mm.win.Error(fmt.Errorf("[MetaMask] [IsConnected] [error]: %v", err))
		return false
	}
	return tv.Bool()
}

func (mm *MetaMask) EnableEthereum() {
	prom := mm.provider.Call(function__request,
		map[string]interface{}{
			request__method: method__requestAccounts,
		},
	)
	mm.win.Debug("[MetaMask] [EnableEthereum] request:")
	mm.win.LogValue(prom)

	prom.Call(web.PROMISE__then, js.FuncOf(mm.accountsChanged))
	prom.Call(web.PROMISE__catch, js.FuncOf(mm.handleError))
}

///////////////////////////////////// DEFAULT EVENT HANDLERS /////////////////////////////////////

func (mm *MetaMask) connect(_this js.Value, _args []js.Value) interface{} {
	// mm.win.Debug("[MetaMask] [connect] args[0] js.Value:")
	// mm.win.LogValue(_args[0])
	tv := _args[0].Get(chain__id)
	if err := web.ValidJSValue(fmt.Sprintf("ConnectInfo.%s", chain__id), tv); err != nil {
		mm.win.Error(fmt.Errorf("[MetaMask] [connect] [error]: %v", err))
		return nil
	}
	id := tv.String()
	name, ok := chain__names[ChainID(id)]
	if !ok {
		mm.win.Error(fmt.Errorf("[MetaMask] [connect] [error]: invalid chain[%s]", id))
		return nil
	}

	mm.win.Info(fmt.Sprintf("[MetaMask] [connect] chain[%s]", name))

	// release the function (?) mm.connect.Release()
	return nil
}

func (mm *MetaMask) disconnect(_this js.Value, _args []js.Value) interface{} {
	mm.win.Debug("[MetaMask] [disconnect] args[0] js.Value:")
	mm.win.LogValue(_args[0])

	// release the function (?) mm.disconnect.Release()
	return nil
}

func (mm *MetaMask) accountsChanged(_this js.Value, _args []js.Value) interface{} {
	var ts string
	if _args[0].Length() < 1 {
		ts = ""
	} else {
		tv := _args[0].Index(0)
		if err := web.ValidJSValue("ethereum.accounts[0]", tv); err != nil {
			mm.win.Error(fmt.Errorf("[MetaMask] [accountsChanged] [error]: %v", err))
			return nil
		}
		ts = tv.String()
	}

	if ts != mm.address {
		mm.address = ts
		mm.win.Debug(fmt.Sprintf("[MetaMask] [accountsChanged] address[%s]", mm.address))
		mm.AddressEl.SetAttribute("text", map[string]interface{}{"value": mm.address})

		if mm.address == "" {
			mm.ConnectEl.SetClass("trigger")
			mm.ConnectEl.SetAttribute("visible", map[string]interface{}{"var": true})
		} else {
			mm.ConnectEl.SetClass("")
			mm.ConnectEl.SetAttribute("visible", map[string]interface{}{"var": false})
		}
	}

	// release the function (?) mm.accountsChanged.Release()
	return nil
}

func (mm *MetaMask) chainChanged(_this js.Value, _args []js.Value) interface{} {
	mm.win.Debug("[MetaMask] [chainChanged] args[0] js.Value:")
	mm.win.LogValue(_args[0])

	// release the function (?) mm.chainChanged.Release()
	return nil
}

func (mm *MetaMask) message(_this js.Value, _args []js.Value) interface{} {
	mm.win.Debug("[MetaMask] [message] args[0] js.Value:")
	mm.win.LogValue(_args[0])

	// release the function (?) mm.message.Release()
	return nil
}

func (mm *MetaMask) handleError(_this js.Value, _args []js.Value) interface{} {
	tv := _args[0].Get(error__message)
	if err := web.ValidJSValue(fmt.Sprintf("error.%s", error__message), tv); err != nil {
		mm.win.Error(fmt.Errorf("[MetaMask] [handleError] [error]: %v", err))
		return nil
	}
	message := tv.String()

	tv = _args[0].Get(error__code)
	if err := web.ValidJSValue(fmt.Sprintf("error.%s", error__code), tv); err != nil {
		mm.win.Error(fmt.Errorf("[MetaMask] [handleError] [error]: %v", err))
		return nil
	}
	code := tv.Int()

	switch code {
	case 4001: // The request was rejected by the user
		mm.ConnectEl.SetClass("trigger")
		mm.ConnectEl.Emit("genRelease", map[string]interface{}{}, false)
	case 32602: // The parameters were invalid
	case 32603: // Internal error
	}

	mm.win.Error(fmt.Errorf("[MetaMask] [handleError] - {%d} %s", code, message))

	return nil
}

// func (mm *MetaMask) SetCallback(_pe ProviderEvent, _f js.Func) error {
// 	if err := goweb.ValidJSValue("window.ethereum", mm.provider); err != nil {
// 		return fmt.Errorf("[MetaMask] [SetCallback] [%s] [error]: %v", event__names[_pe], err)
// 	}

// 	var def js.Func
// 	switch _pe {
// 	case Connect:
// 	case Disconnect:
// 	case AccountsChanged:
// 	case ChainChanged:
// 	case Message:
// 	}

// 	wf := wrap()
// 	mm.provider.Call("on", event__names[_pe], _f)

// 	return nil
// }

// func wrap(_default, _ext js.Func) (func(js.Value, map[string]interface{}) interface{}) {
// 	return func(_this js.Value, _data map[string]interface{}) interface{} {
// 		_default(_this, _ext)
// 		_ext(_this, _ext)

// 		return js.ValueOf(nil)
// 	}
// }
