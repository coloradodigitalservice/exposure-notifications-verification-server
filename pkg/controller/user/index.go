// Copyright 2020 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package user

import (
	"context"
	"net/http"

	"github.com/google/exposure-notifications-verification-server/pkg/controller"
	"github.com/google/exposure-notifications-verification-server/pkg/database"
	"github.com/google/exposure-notifications-verification-server/pkg/pagination"
)

const (
	// QueryKeySearch is the query key where the search query exists.
	QueryKeySearch = "q"
)

func (c *Controller) HandleIndex() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		realm := controller.RealmFromContext(ctx)
		if realm == nil {
			controller.MissingRealm(w, r, c.h)
			return
		}

		pageParams, err := pagination.FromRequest(r)
		if err != nil {
			controller.InternalError(w, r, c.h, err)
			return
		}

		var scopes []database.Scope
		q := r.FormValue(QueryKeySearch)
		scopes = append(scopes, database.WithUserSearch(q))

		users, paginator, err := realm.ListUsers(c.db, pageParams, scopes...)
		if err != nil {
			controller.InternalError(w, r, c.h, err)
			return
		}

		c.renderIndex(ctx, w, users, paginator, q)
	})
}

func (c *Controller) renderIndex(
	ctx context.Context, w http.ResponseWriter,
	users []*database.User, paginator *pagination.Paginator, query string) {
	m := controller.TemplateMapFromContext(ctx)
	m.Title("Users")
	m["users"] = users
	m["paginator"] = paginator
	m["query"] = query
	c.h.RenderHTML(w, "users/index", m)
}
