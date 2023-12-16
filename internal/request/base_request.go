package request

import (
	"html"
	"strings"

	"github.com/rizalarfiyan/be-revend/constants"
	"github.com/rizalarfiyan/be-revend/utils"
)

type BasePagination struct {
	Page    int                        `json:"page" field:"Page" validate:"omitempty,min=1" example:"1"`
	Limit   int                        `json:"limit" field:"Limit" example:"10"`
	Search  string                     `json:"search" field:"Search"`
	OrderBy string                     `json:"order_by" field:"Order By"`
	Order   string                     `json:"order" field:"Order"`
	Status  constants.FilterListStatus `json:"status" field:"Status"`
}

func (bp *BasePagination) Normalize() {
	if bp.Search != "" {
		bp.Search = html.EscapeString(utils.Str(bp.Search))
	}

	if bp.Status != "" {
		status := constants.FilterListStatus(utils.Str(string(bp.Status)))
		if status.IsValid() {
			bp.Status = status
		}
	}
}

func (bp *BasePagination) ValidateAndUpdateOrderBy(data map[string]string) {
	keyOrderBy := strings.ToLower(bp.OrderBy)
	getOrderBy, isFoundOrderBy := data[keyOrderBy]
	if !isFoundOrderBy {
		bp.OrderBy = ""
		bp.Order = ""
		return
	}

	bp.OrderBy = getOrderBy
	orderMap := map[string]bool{
		"asc":  true,
		"desc": true,
	}

	keyOrder := strings.ToLower(bp.Order)
	if _, ok := orderMap[keyOrder]; ok {
		bp.Order = keyOrder
	} else {
		bp.Order = "asc"
	}
}
