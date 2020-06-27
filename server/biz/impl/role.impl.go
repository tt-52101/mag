package impl

import (
	"context"

	"github.com/casbin/casbin/v2"
	"github.com/google/wire"
	"github.com/key7men/mag/pkg/errs"
	"github.com/key7men/mag/server/assist/uuid"
	"github.com/key7men/mag/server/biz"
	"github.com/key7men/mag/server/model"
	"github.com/key7men/mag/server/schema"
)

var _ biz.IRole = (*Role)(nil)

// RoleSet 注入Role
var RoleSet = wire.NewSet(wire.Struct(new(Role), "*"), wire.Bind(new(biz.IRole), new(*Role)))

// Role 角色管理
type Role struct {
	Enforcer      *casbin.SyncedEnforcer
	TransModel    model.ITrans
	RoleModel     model.IRole
	RoleMenuModel model.IRoleMenu
	UserModel     model.IUser
}

// Query 查询数据
func (a *Role) Query(ctx context.Context, params schema.RoleQueryParam, opts ...schema.RoleQueryOptions) (*schema.RoleQueryResult, error) {
	return a.RoleModel.Query(ctx, params, opts...)
}

// Get 查询指定数据
func (a *Role) Get(ctx context.Context, id string, opts ...schema.RoleQueryOptions) (*schema.Role, error) {
	item, err := a.RoleModel.Get(ctx, id, opts...)
	if err != nil {
		return nil, err
	} else if item == nil {
		return nil, errs.ErrNotFound
	}

	roleMenus, err := a.QueryRoleMenus(ctx, id)
	if err != nil {
		return nil, err
	}
	item.RoleMenus = roleMenus

	return item, nil
}

// QueryRoleMenus 查询角色菜单列表
func (a *Role) QueryRoleMenus(ctx context.Context, roleID string) (schema.RoleMenus, error) {
	result, err := a.RoleMenuModel.Query(ctx, schema.RoleMenuQueryParam{
		RoleID: roleID,
	})
	if err != nil {
		return nil, err
	}
	return result.Data, nil
}

// Create 创建数据
func (a *Role) Create(ctx context.Context, item schema.Role) (*schema.IDResult, error) {
	err := a.checkName(ctx, item)
	if err != nil {
		return nil, err
	}

	item.ID = uuid.NewID()
	err = ExecTrans(ctx, a.TransModel, func(ctx context.Context) error {
		for _, rmItem := range item.RoleMenus {
			rmItem.ID = uuid.NewID()
			rmItem.RoleID = item.ID
			err := a.RoleMenuModel.Create(ctx, *rmItem)
			if err != nil {
				return err
			}
		}
		return a.RoleModel.Create(ctx, item)
	})
	if err != nil {
		return nil, err
	}
	LoadCasbinPolicy(ctx, a.Enforcer)
	return schema.NewIDResult(item.ID), nil
}

func (a *Role) checkName(ctx context.Context, item schema.Role) error {
	result, err := a.RoleModel.Query(ctx, schema.RoleQueryParam{
		PaginationParam: schema.PaginationParam{OnlyCount: true},
		Name:            item.Name,
	})
	if err != nil {
		return err
	} else if result.PageResult.Total > 0 {
		return errs.New400Response("角色名称已经存在")
	}
	return nil
}

// Update 更新数据
func (a *Role) Update(ctx context.Context, id string, item schema.Role) error {
	oldItem, err := a.Get(ctx, id)
	if err != nil {
		return err
	} else if oldItem == nil {
		return errs.ErrNotFound
	} else if oldItem.Name != item.Name {
		err := a.checkName(ctx, item)
		if err != nil {
			return err
		}
	}

	item.ID = oldItem.ID
	item.Creator = oldItem.Creator
	item.CreatedAt = oldItem.CreatedAt
	err = ExecTrans(ctx, a.TransModel, func(ctx context.Context) error {
		addRoleMenus, delRoleMenus := a.compareRoleMenus(ctx, oldItem.RoleMenus, item.RoleMenus)
		for _, rmitem := range addRoleMenus {
			rmitem.ID = uuid.NewID()
			rmitem.RoleID = id
			err := a.RoleMenuModel.Create(ctx, *rmitem)
			if err != nil {
				return err
			}
		}

		for _, rmitem := range delRoleMenus {
			err := a.RoleMenuModel.Delete(ctx, rmitem.ID)
			if err != nil {
				return err
			}
		}

		return a.RoleModel.Update(ctx, id, item)
	})
	if err != nil {
		return err
	}
	LoadCasbinPolicy(ctx, a.Enforcer)
	return nil
}

func (a *Role) compareRoleMenus(ctx context.Context, oldRoleMenus, newRoleMenus schema.RoleMenus) (addList, delList schema.RoleMenus) {
	mOldRoleMenus := oldRoleMenus.ToMap()
	mNewRoleMenus := newRoleMenus.ToMap()

	for k, item := range mNewRoleMenus {
		if _, ok := mOldRoleMenus[k]; ok {
			delete(mOldRoleMenus, k)
			continue
		}
		addList = append(addList, item)
	}

	for _, item := range mOldRoleMenus {
		delList = append(delList, item)
	}
	return
}

// Delete 删除数据
func (a *Role) Delete(ctx context.Context, id string) error {
	oldItem, err := a.RoleModel.Get(ctx, id)
	if err != nil {
		return err
	} else if oldItem == nil {
		return errs.ErrNotFound
	}

	userResult, err := a.UserModel.Query(ctx, schema.UserQueryParam{
		PaginationParam: schema.PaginationParam{OnlyCount: true},
		RoleIDs:         []string{id},
	})
	if err != nil {
		return err
	} else if userResult.PageResult.Total > 0 {
		return errs.New400Response("该角色已被赋予用户，不允许删除")
	}

	err = ExecTrans(ctx, a.TransModel, func(ctx context.Context) error {
		err := a.RoleMenuModel.DeleteByRoleID(ctx, id)
		if err != nil {
			return err
		}

		return a.RoleModel.Delete(ctx, id)
	})
	if err != nil {
		return err
	}

	LoadCasbinPolicy(ctx, a.Enforcer)
	return nil
}

// UpdateStatus 更新状态
func (a *Role) UpdateStatus(ctx context.Context, id string, status int) error {
	oldItem, err := a.RoleModel.Get(ctx, id)
	if err != nil {
		return err
	} else if oldItem == nil {
		return errs.ErrNotFound
	}

	err = a.RoleModel.UpdateStatus(ctx, id, status)
	if err != nil {
		return err
	}
	LoadCasbinPolicy(ctx, a.Enforcer)
	return nil
}
