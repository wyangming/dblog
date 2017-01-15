package pagination

import (
	"math"
)

//offset查询的起始位置，rows偏移的位置，nums查询出的数据总合
////设置offset时注意数据库是从0开始计数的，如果用的是数据库的计数这里要加1
func NewPaginator(offset, rows, nums int64) *Paginator {
	paginator := &Paginator{}
	paginator.setInfo(offset, rows, nums)
	return paginator
}

//分页信息类
type Paginator struct {
	//数据条目总合
	nums int64
	//总页数
	pageNums int64
	//当前页
	pageCur int64
	//分页大小
	pageSize int64
}

func (this *Paginator) setInfo(offset, rows, nums int64) {
	this.nums = nums
	//设置当前页
	this.pageCur = offset / rows
	if (offset % rows) > 0 {
		this.pageCur += 1
	}
	//如果是第0页的处理
	if this.pageCur < 1 {
		this.pageCur = 1
	}
	this.pageSize = rows
	//设置页面总数
	this.pageNums = nums / rows
	if (nums % rows) > 0 {
		this.pageNums += 1
	}
	//如果总页数是0
	if this.pageNums < 1 {
		this.pageNums = 1
	}
}

//得到总页数
func (this *Paginator) Nums() int64 {
	return this.nums
}

//当前页
func (this *Paginator) Cur() int64 {
	return this.pageCur
}

//分页大小
func (this *Paginator) Size() int64 {
	return this.pageSize
}

//得到所有的页数
func (this *Paginator) Pages() []int64 {
	var pages []int64
	pageNums := this.pageNums
	page := this.pageCur
	switch {
	case page >= pageNums-4 && pageNums > 9:
		start := int64(pageNums - 9 + 1)
		pages = make([]int64, 9)
		for i := range pages {
			pages[i] = int64(start + int64(i))
		}
	case page >= 5 && pageNums > 9:
		start := int64(page - 5 + 1)
		pages = make([]int64, int64(math.Min(9, float64(page+4+1))))
		for i := range pages {
			pages[i] = int64(start + int64(i))
		}
	default:
		pages = make([]int64, int64(math.Min(9, float64(pageNums))))
		for i := range pages {
			pages[i] = int64(i + 1)
		}
	}
	return pages
}

//首页
func (this *Paginator) FirstPage() int64 {
	return int64(1)
}

//上一页
func (this *Paginator) PrevPage() int64 {
	return this.pageCur - 1
}

//是否还有上一页
func (this *Paginator) HasPrev() bool {
	return this.pageCur > 1
}

//是否还有下一页
func (this *Paginator) HasNext() bool {
	return this.pageCur < this.pageNums
}

//下一页
func (this *Paginator) NextPage() int64 {
	return this.pageCur + 1
}

//最后一页
func (this *Paginator) LastPage() int64 {
	return this.pageNums
}

//判断是否是当前页面
func (this *Paginator) IsActive(page int64) bool {
	return this.pageCur == page
}
