// Copyright (C) 2015  TF2Stadium
// Use of this source code is governed by the GPLv3
// that can be found in the COPYING file.

package routes

import (
	"net/http"

	"github.com/TF2Stadium/Helen/config"
	"github.com/TF2Stadium/Helen/controllers"
	"github.com/TF2Stadium/Helen/controllers/admin"
	chelpers "github.com/TF2Stadium/Helen/controllers/controllerhelpers"
	"github.com/TF2Stadium/Helen/controllers/login"
	"github.com/TF2Stadium/Helen/helpers"
)

type route struct {
	pattern string
	handler http.HandlerFunc
}

var routes = []route{
	{"/", controllers.MainHandler},
	{"/openidcallback", login.LoginCallbackHandler},
	{"/startLogin", login.LoginHandler},
	{"/startTwitchLogin", login.TwitchLogin},
	{"/twitchAuth", login.TwitchAuth},
	{"/twitchLogout", login.TwitchLogout},
	{"/logout", login.LogoutHandler},
	{"/websocket/", controllers.SocketHandler},
	{"/startMockLogin", login.MockLoginHandler},

	{"/admin", chelpers.FilterHTTPRequest(helpers.ActionViewPage, admin.ServeAdminPage)},
	{"/admin/roles", chelpers.FilterHTTPRequest(helpers.ActionViewPage, admin.ChangeRole)},
	{"/admin/ban", chelpers.FilterHTTPRequest(helpers.ActionViewPage, admin.BanPlayer)},
	{"/admin/chatlogs", chelpers.FilterHTTPRequest(helpers.ActionViewLogs, admin.GetChatLogs)},
	{"/admin/banlogs", chelpers.FilterHTTPRequest(helpers.ActionViewLogs, admin.GetBanLogs)},
	{"/admin/server/", chelpers.FilterHTTPRequest(helpers.ModifyServers, admin.ViewServerPage)},
	{"/admin/server/add", chelpers.FilterHTTPRequest(helpers.ModifyServers, admin.AddServer)},
	{"/admin/server/remove", chelpers.FilterHTTPRequest(helpers.ModifyServers, admin.RemoveServer)},
	{"/admin/lobbies", chelpers.FilterHTTPRequest(helpers.ActionViewLogs, admin.ViewOpenLobbies)},

	{"/health", controllers.Health},
	{"/badge/", controllers.TwitchBadge},
}

func SetupHTTP(mux *http.ServeMux) {
	for _, route := range routes {
		mux.HandleFunc(route.pattern, route.handler)
	}

	if config.Constants.ServeStatic {
		mux.HandleFunc("/static/", func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, "views/static.html")
		})

	}
}
