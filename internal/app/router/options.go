package router

import (
	"github.com/alexedwards/scs/v2"
	"github.com/caiolandgraf/go-project-base/internal/app/middleware"
	"github.com/go-fuego/fuego"
	"github.com/go-fuego/fuego/option"
)

func ToOptions(doc Doc, session *scs.SessionManager) []fuego.RouteOption {
	opts := make([]fuego.RouteOption, 0, 8+len(doc.Responses)+len(doc.Query))

	if doc.Detail.Summary != "" {
		opts = append(opts, option.Summary(doc.Detail.Summary))
	}

	if doc.Detail.Description != "" {
		opts = append(opts, option.Description(doc.Detail.Description))
	}

	if len(doc.Detail.Tags) > 0 {
		opts = append(opts, option.Tags(doc.Detail.Tags...))
		for _, tag := range doc.Detail.Tags {
			opts = append(opts, option.TagInfo(tag, ""))
		}
	}

	if doc.Detail.OperationID != "" {
		opts = append(opts, option.OperationID(doc.Detail.OperationID))
	}

	for code, resp := range doc.Responses {
		opts = append(opts, option.AddResponse(
			code,
			resp.Description,
			fuego.Response{Type: resp.Type},
		))
	}

	for _, q := range doc.Query {
		opts = append(opts, queryOption(q))
	}

	if doc.Auth && session != nil {
		opts = append(opts, option.Middleware(middleware.AuthRequired(session)))
	}

	return opts
}

func queryOption(q QueryParam) fuego.RouteOption {
	switch q.Type {
	case QueryInt:
		params := queryParams(q.Default)
		return option.QueryInt(q.Name, q.Description, params...)
	case QueryBool:
		params := queryParams(q.Default)
		return option.QueryBool(q.Name, q.Description, params...)
	default:
		params := queryParams(q.Default)
		return option.Query(q.Name, q.Description, params...)
	}
}

func queryParams(defaultValue any) []fuego.ParamOption {
	if defaultValue == nil {
		return nil
	}

	switch v := defaultValue.(type) {
	case int:
		return []fuego.ParamOption{fuego.ParamDefault(v)}
	case string:
		return []fuego.ParamOption{fuego.ParamDefault(v)}
	case bool:
		return []fuego.ParamOption{fuego.ParamDefault(v)}
	default:
		return nil
	}
}
