package router

// Doc holds declarative OpenAPI metadata for a route.
// Request body schemas are inferred from the handler generic (ContextWithBody[T]).
type Doc struct {
	Auth      bool
	Detail    Detail
	Responses map[int]Response
	Query     []QueryParam
}

type Detail struct {
	Tags        []string
	Summary     string
	Description string
	OperationID string
}

type Response struct {
	Description string
	Type        any
}

type QueryType int

const (
	QueryInt QueryType = iota
	QueryString
	QueryBool
)

type QueryParam struct {
	Name        string
	Description string
	Type        QueryType
	Default     any
}
