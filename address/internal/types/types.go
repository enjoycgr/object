// Code generated by goctl. DO NOT EDIT.
package types

type Request struct {
	Name string `path:"name,options=you|me"`
}

type Response struct {
	JingUuid string `json:"jing_uuid"`
	ParentId string `json:"parent_id"`
	Name     string `json:"name"`
	Code     string `json:"code"`
	Child    []*Response `json:"child"`
}
