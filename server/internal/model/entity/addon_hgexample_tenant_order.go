// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// AddonHgexampleTenantOrder is the golang structure for table addon_hgexample_tenant_order.
type AddonHgexampleTenantOrder struct {
	Id          int64       `json:"id"          description:"主键"`
	TenantId    int64       `json:"tenantId"    description:"租户ID"`
	MerchantId  int64       `json:"merchantId"  description:"商户ID"`
	UserId      int64       `json:"userId"      description:"用户ID"`
	ProductName string      `json:"productName" description:"购买产品"`
	OrderSn     string      `json:"orderSn"     description:"订单号"`
	Money       float64     `json:"money"       description:"充值金额"`
	Remark      string      `json:"remark"      description:"备注"`
	Status      int         `json:"status"      description:"订单状态"`
	CreatedAt   *gtime.Time `json:"createdAt"   description:"创建时间"`
	UpdatedAt   *gtime.Time `json:"updatedAt"   description:"修改时间"`
}
