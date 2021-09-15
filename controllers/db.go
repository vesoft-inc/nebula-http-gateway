package controllers

import (
	"encoding/json"

	"github.com/astaxie/beego"
	"github.com/vesoft-inc/nebula-http-gateway/common"
	"github.com/vesoft-inc/nebula-http-gateway/service/dao"
)

type DatabaseController struct {
	beego.Controller
}

type Response struct {
	Code    int        `json:"code"`
	Data    common.Any `json:"data"`
	Message string     `json:"message"`
}

type Request struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Address  string `json:"address"`
	Port     int    `json:"port"`
}

type ExecuteRequest struct {
	Gql       string             `json:"gql"`
	ParamsMap common.ParamterMap `json:"paramsMap"`
}

type Data map[string]interface{}

func (this *DatabaseController) Connect() {
	var res Response
	var params Request
	json.Unmarshal(this.Ctx.Input.RequestBody, &params)
	nsid, err := dao.Connect(params.Address, params.Port, params.Username, params.Password)
	if err == nil {
		res.Code = 0
		m := make(map[string]common.Any)
		m["nsid"] = nsid
		res.Data = nsid
		this.Ctx.SetCookie("Secure", "true")
		this.Ctx.SetCookie("SameSite", "None")
		this.SetSession("nsid", nsid)

		res.Message = "Login successfully"
	} else {
		res.Code = -1
		res.Message = err.Error()
	}
	this.Data["json"] = &res
	this.ServeJSON()
}

func (this *DatabaseController) Home() {
	var res Response
	res.Code = 0
	this.ServeJSON()
}

func (this *DatabaseController) Disconnect() {
	var res Response
	nsid := this.GetSession("nsid")
	if nsid != nil {
		dao.Disconnect(nsid.(string))
	}
	res.Code = 0
	res.Message = "Disconnect successfully"
	this.Data["json"] = &res
	this.ServeJSON()
}

func (this *DatabaseController) Execute() {
	var res Response
	var params ExecuteRequest
	nsid := this.GetSession("nsid")
	if nsid == nil {
		res.Code = -1
		res.Message = "connection refused for lack of session"
	} else {
		json.Unmarshal(this.Ctx.Input.RequestBody, &params)
		result, err := dao.Execute(nsid.(string), params.Gql, params.ParamsMap)
		if err == nil {
			res.Code = 0
			res.Data = &result
		} else {
			res.Code = -1
			res.Message = err.Error()
		}
	}
	this.Data["json"] = &res
	this.ServeJSON()
}
