package coroutine

type SortDirection string

const (
	SortDirectionAsc  SortDirection = "asc"
	SortDirectionDesc SortDirection = "desc"
)

type Optional func(*options)

type options struct {
	sortDirection SortDirection
}

func SetSortDirection(direction SortDirection) Optional {
	return func(o *options) {
		o.sortDirection = direction
	}
}
