package router

import (
	"github.com/alexedwards/scs/v2"
	"github.com/go-fuego/fuego"
)

func Get[T, B, P any](
	s *fuego.Server,
	path string,
	handler func(fuego.Context[B, P]) (T, error),
	doc Doc,
	session *scs.SessionManager,
) *fuego.Route[T, B, P] {
	return fuego.Get(s, path, handler, ToOptions(doc, session)...)
}

func Post[T, B, P any](
	s *fuego.Server,
	path string,
	handler func(fuego.Context[B, P]) (T, error),
	doc Doc,
	session *scs.SessionManager,
) *fuego.Route[T, B, P] {
	return fuego.Post(s, path, handler, ToOptions(doc, session)...)
}

func Put[T, B, P any](
	s *fuego.Server,
	path string,
	handler func(fuego.Context[B, P]) (T, error),
	doc Doc,
	session *scs.SessionManager,
) *fuego.Route[T, B, P] {
	return fuego.Put(s, path, handler, ToOptions(doc, session)...)
}

func Delete[T, B, P any](
	s *fuego.Server,
	path string,
	handler func(fuego.Context[B, P]) (T, error),
	doc Doc,
	session *scs.SessionManager,
) *fuego.Route[T, B, P] {
	return fuego.Delete(s, path, handler, ToOptions(doc, session)...)
}
