package modules

import (
	"covid19kalteng/components/basemodel"
	"strconv"
	"strings"

	"github.com/labstack/echo"
)

//Pagination helper for short step init pagination filter
type Pagination struct {
	*basemodel.BaseModel //must pointer
}

func (pc *Pagination) SetPaginationFilter(c echo.Context) {
	pc.Rows, _ = strconv.Atoi(c.QueryParam("rows"))
	pc.Page, _ = strconv.Atoi(c.QueryParam("page"))
	pc.Order = strings.Split(c.QueryParam("orderby"), ",")
	pc.Sort = strings.Split(c.QueryParam("sort"), ",")
}

func SetPaginationFilter(obj interface{}, c echo.Context) {
	//pointer to object PaginationContext to handle default filter (pagination)
	//pass by reference so changes in this func fill affected after step out from this
	page := Pagination{obj.(*basemodel.BaseModel)}

	//call pagination set up from context
	page.SetPaginationFilter(c)
}
