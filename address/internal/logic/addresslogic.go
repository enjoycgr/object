package logic

import (
	"context"

	"github.com/jinzhu/copier"
	"github.com/tal-tech/go-zero/core/logx"
	"go-zero-demo/address/internal/model"
	"go-zero-demo/address/internal/svc"
	"go-zero-demo/address/internal/types"
)

type AddressLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

type Tree struct {
	JingUuid string
	ParentId string
	Name     string
	Code     string
	Child    []*Tree
}

func NewAddressLogic(ctx context.Context, svcCtx *svc.ServiceContext) AddressLogic {
	return AddressLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AddressLogic) Address() (resp *types.Response, err error) {
	resp = new(types.Response)
	list, err := l.svcCtx.AddressModel.List()
	if err != nil && err != model.ErrNotFound {
		return nil, err
	}

	tree := l.tree(list)
	err = copier.Copy(resp, tree)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (l *AddressLogic) tree(list []*model.Address) *Tree {
	listMap := make(map[string]*Tree, len(list))
	treeMap := make(map[string][]*Tree, len(list))
	root := &Tree{}

	for _, l := range list {
		t := &Tree{
			JingUuid: l.JingUuid,
			Name:     l.Value,
			ParentId: l.ParentId,
			Code:     l.Code,
		}
		listMap[l.JingUuid] = t
		treeMap[l.ParentId] = append(treeMap[l.ParentId], t)
	}

	for _, l := range list {
		if l.ParentId == "" {
			root.Child = append(root.Child, listMap[l.JingUuid])
		}
		listMap[l.JingUuid].Child = treeMap[l.JingUuid]
	}

	return root
}
