package v1

import (
	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/traPtitech/traQ/repository"
	"github.com/traPtitech/traQ/router/session"
	random2 "github.com/traPtitech/traQ/utils/random"
	"net/http"
	"testing"
)

func TestHandlers_GetUserGroups(t *testing.T) {
	t.Parallel()
	env, _, _, s, _, _, adminUser := setupWithUsers(t, s1)

	env.mustMakeUserGroup(t, rand, adminUser.GetID())
	env.mustMakeUserGroup(t, rand, adminUser.GetID())
	env.mustMakeUserGroup(t, rand, adminUser.GetID())

	t.Run("NotLoggedIn", func(t *testing.T) {
		t.Parallel()
		e := env.makeExp(t)
		e.GET("/api/1.0/groups").
			Expect().
			Status(http.StatusUnauthorized)
	})

	t.Run("ok", func(t *testing.T) {
		t.Parallel()
		e := env.makeExp(t)
		e.GET("/api/1.0/groups").
			WithCookie(session.CookieName, s).
			Expect().
			Status(http.StatusOK).
			JSON().
			Array().
			Length().
			Equal(3)
	})
}

func TestHandlers_PostUserGroups(t *testing.T) {
	t.Parallel()
	env, _, _, s, _, user, adminUser := setupWithUsers(t, common5)

	t.Run("NotLoggedIn", func(t *testing.T) {
		t.Parallel()
		e := env.makeExp(t)
		name := random2.AlphaNumeric(20)
		e.POST("/api/1.0/groups").
			WithJSON(map[string]interface{}{"name": name, "description": name}).
			Expect().
			Status(http.StatusUnauthorized)
	})

	t.Run("bad request", func(t *testing.T) {
		t.Parallel()
		e := env.makeExp(t)
		e.POST("/api/1.0/groups").
			WithCookie(session.CookieName, s).
			WithJSON(map[string]interface{}{"name": true}).
			Expect().
			Status(http.StatusBadRequest)
	})

	t.Run("conflict", func(t *testing.T) {
		t.Parallel()
		e := env.makeExp(t)
		name := random2.AlphaNumeric(20)
		env.mustMakeUserGroup(t, name, adminUser.GetID())
		e.POST("/api/1.0/groups").
			WithCookie(session.CookieName, s).
			WithJSON(map[string]interface{}{"name": name, "description": name}).
			Expect().
			Status(http.StatusConflict)
	})

	t.Run("ok", func(t *testing.T) {
		t.Parallel()
		e := env.makeExp(t)
		name := random2.AlphaNumeric(20)
		obj := e.POST("/api/1.0/groups").
			WithCookie(session.CookieName, s).
			WithJSON(map[string]interface{}{"name": name, "description": name}).
			Expect().
			Status(http.StatusCreated).
			JSON().
			Object()

		obj.Value("groupId").String().NotEmpty()
		obj.Value("name").String().Equal(name)
		obj.Value("description").String().Equal(name)
		obj.Value("adminUserId").String().Equal(user.GetID().String())
		obj.Value("members").Array().Empty()
	})
}

func TestHandlers_GetUserGroup(t *testing.T) {
	t.Parallel()
	env, _, _, s, _, user, adminUser := setupWithUsers(t, common5)

	g := env.mustMakeUserGroup(t, rand, adminUser.GetID())
	env.mustAddUserToGroup(t, user.GetID(), g.ID)

	t.Run("NotLoggedIn", func(t *testing.T) {
		t.Parallel()
		e := env.makeExp(t)
		e.GET("/api/1.0/groups/{groupID}", g.ID.String()).
			Expect().
			Status(http.StatusUnauthorized)
	})

	t.Run("not found", func(t *testing.T) {
		t.Parallel()
		e := env.makeExp(t)
		e.GET("/api/1.0/groups/{groupID}", uuid.Must(uuid.NewV4())).
			WithCookie(session.CookieName, s).
			Expect().
			Status(http.StatusNotFound)
	})

	t.Run("ok", func(t *testing.T) {
		t.Parallel()
		e := env.makeExp(t)
		obj := e.GET("/api/1.0/groups/{groupID}", g.ID.String()).
			WithCookie(session.CookieName, s).
			Expect().
			Status(http.StatusOK).
			JSON().
			Object()

		obj.Value("groupId").String().Equal(g.ID.String())
		obj.Value("name").String().Equal(g.Name)
		obj.Value("description").String().Equal(g.Description)
		obj.Value("adminUserId").String().Equal(g.Admins[0].UserID.String())
		obj.Value("members").Array().ContainsOnly(user.GetID().String())
	})
}

func TestHandlers_PatchUserGroup(t *testing.T) {
	t.Parallel()
	env, _, _, s, _, user, adminUser := setupWithUsers(t, common5)

	user2 := env.mustMakeUser(t, rand)
	g := env.mustMakeUserGroup(t, rand, user.GetID())

	t.Run("NotLoggedIn", func(t *testing.T) {
		t.Parallel()
		e := env.makeExp(t)
		e.PATCH("/api/1.0/groups/{groupID}", g.ID.String()).
			Expect().
			Status(http.StatusUnauthorized)
	})

	t.Run("not found", func(t *testing.T) {
		t.Parallel()
		e := env.makeExp(t)
		e.PATCH("/api/1.0/groups/{groupID}", uuid.Must(uuid.NewV4())).
			WithCookie(session.CookieName, s).
			WithJSON(map[string]interface{}{"description": "aaa"}).
			Expect().
			Status(http.StatusNotFound)
	})

	t.Run("bad request", func(t *testing.T) {
		t.Parallel()
		e := env.makeExp(t)
		e.PATCH("/api/1.0/groups/{groupID}", g.ID.String()).
			WithCookie(session.CookieName, s).
			WithJSON(map[string]interface{}{"name": true}).
			Expect().
			Status(http.StatusBadRequest)
	})

	t.Run("conflict", func(t *testing.T) {
		t.Parallel()
		e := env.makeExp(t)
		name := random2.AlphaNumeric(20)
		env.mustMakeUserGroup(t, name, adminUser.GetID())
		e.PATCH("/api/1.0/groups/{groupID}", g.ID.String()).
			WithCookie(session.CookieName, s).
			WithJSON(map[string]interface{}{"name": name}).
			Expect().
			Status(http.StatusConflict)
	})

	t.Run("forbidden", func(t *testing.T) {
		t.Parallel()
		e := env.makeExp(t)
		e.PATCH("/api/1.0/groups/{groupID}", g.ID.String()).
			WithCookie(session.CookieName, env.generateSession(t, user2.GetID())).
			WithJSON(map[string]interface{}{"description": "aaa"}).
			Expect().
			Status(http.StatusForbidden)
	})

	t.Run("ok", func(t *testing.T) {
		t.Parallel()
		e := env.makeExp(t)
		g := env.mustMakeUserGroup(t, rand, user.GetID())
		name := random2.AlphaNumeric(20)
		e.PATCH("/api/1.0/groups/{groupID}", g.ID.String()).
			WithCookie(session.CookieName, s).
			WithJSON(map[string]interface{}{"name": name, "description": "aaa"}).
			Expect().
			Status(http.StatusNoContent)

		a, err := env.Repository.GetUserGroup(g.ID)
		if assert.NoError(t, err) {
			assert.Equal(t, a.Name, name)
			assert.Equal(t, a.Description, "aaa")
		}
	})

}

func TestHandlers_DeleteUserGroup(t *testing.T) {
	t.Parallel()
	env, _, _, s, _, user, _ := setupWithUsers(t, common5)

	g := env.mustMakeUserGroup(t, rand, user.GetID())

	t.Run("NotLoggedIn", func(t *testing.T) {
		t.Parallel()
		e := env.makeExp(t)
		e.DELETE("/api/1.0/groups/{groupID}", g.ID.String()).
			Expect().
			Status(http.StatusUnauthorized)
	})

	t.Run("not found", func(t *testing.T) {
		t.Parallel()
		e := env.makeExp(t)
		e.DELETE("/api/1.0/groups/{groupID}", uuid.Must(uuid.NewV4())).
			WithCookie(session.CookieName, s).
			Expect().
			Status(http.StatusNotFound)
	})

	t.Run("forbidden", func(t *testing.T) {
		t.Parallel()
		user2 := env.mustMakeUser(t, rand)
		e := env.makeExp(t)
		e.DELETE("/api/1.0/groups/{groupID}", g.ID.String()).
			WithCookie(session.CookieName, env.generateSession(t, user2.GetID())).
			Expect().
			Status(http.StatusForbidden)
	})

	t.Run("ok", func(t *testing.T) {
		t.Parallel()
		g := env.mustMakeUserGroup(t, rand, user.GetID())
		e := env.makeExp(t)
		e.DELETE("/api/1.0/groups/{groupID}", g.ID.String()).
			WithCookie(session.CookieName, s).
			Expect().
			Status(http.StatusNoContent)

		_, err := env.Repository.GetUserGroup(g.ID)
		assert.EqualError(t, err, repository.ErrNotFound.Error())
	})
}

func TestHandlers_GetUserGroupMembers(t *testing.T) {
	t.Parallel()
	env, _, _, s, _, user, adminUser := setupWithUsers(t, common5)

	g := env.mustMakeUserGroup(t, rand, adminUser.GetID())
	env.mustAddUserToGroup(t, user.GetID(), g.ID)

	t.Run("NotLoggedIn", func(t *testing.T) {
		t.Parallel()
		e := env.makeExp(t)
		e.GET("/api/1.0/groups/{groupID}/members", g.ID.String()).
			Expect().
			Status(http.StatusUnauthorized)
	})

	t.Run("not found", func(t *testing.T) {
		t.Parallel()
		e := env.makeExp(t)
		e.GET("/api/1.0/groups/{groupID}/members", uuid.Must(uuid.NewV4())).
			WithCookie(session.CookieName, s).
			Expect().
			Status(http.StatusNotFound)
	})

	t.Run("ok", func(t *testing.T) {
		t.Parallel()
		e := env.makeExp(t)
		e.GET("/api/1.0/groups/{groupID}/members", g.ID.String()).
			WithCookie(session.CookieName, s).
			Expect().
			Status(http.StatusOK).
			JSON().
			Array().
			ContainsOnly(user.GetID().String())
	})
}

func TestHandlers_PostUserGroupMembers(t *testing.T) {
	t.Parallel()
	env, _, _, s, _, user, _ := setupWithUsers(t, common5)
	g := env.mustMakeUserGroup(t, rand, user.GetID())
	user2 := env.mustMakeUser(t, rand)

	t.Run("NotLoggedIn", func(t *testing.T) {
		t.Parallel()
		e := env.makeExp(t)
		e.POST("/api/1.0/groups/{groupID}/members", g.ID.String()).
			WithJSON(map[string]interface{}{"userId": user.GetID()}).
			Expect().
			Status(http.StatusUnauthorized)
	})

	t.Run("not found", func(t *testing.T) {
		t.Parallel()
		e := env.makeExp(t)
		e.POST("/api/1.0/groups/{groupID}/members", uuid.Must(uuid.NewV4())).
			WithCookie(session.CookieName, s).
			WithJSON(map[string]interface{}{"userId": user.GetID()}).
			Expect().
			Status(http.StatusNotFound)
	})

	t.Run("bad request", func(t *testing.T) {
		t.Parallel()
		e := env.makeExp(t)
		e.POST("/api/1.0/groups/{groupID}/members", g.ID.String()).
			WithCookie(session.CookieName, s).
			WithJSON(map[string]interface{}{"userId": true}).
			Expect().
			Status(http.StatusBadRequest)
	})

	t.Run("forbidden", func(t *testing.T) {
		t.Parallel()
		e := env.makeExp(t)
		e.POST("/api/1.0/groups/{groupID}/members", g.ID.String()).
			WithCookie(session.CookieName, env.generateSession(t, user2.GetID())).
			WithJSON(map[string]interface{}{"userId": user.GetID()}).
			Expect().
			Status(http.StatusForbidden)
	})

	t.Run("unknown user", func(t *testing.T) {
		t.Parallel()
		e := env.makeExp(t)
		e.POST("/api/1.0/groups/{groupID}/members", g.ID.String()).
			WithCookie(session.CookieName, s).
			WithJSON(map[string]uuid.UUID{"userId": uuid.Must(uuid.NewV4())}).
			Expect().
			Status(http.StatusBadRequest)
	})

	t.Run("ok", func(t *testing.T) {
		t.Parallel()
		e := env.makeExp(t)
		e.POST("/api/1.0/groups/{groupID}/members", g.ID.String()).
			WithCookie(session.CookieName, s).
			WithJSON(map[string]interface{}{"userId": user.GetID()}).
			Expect().
			Status(http.StatusNoContent)

		ids, err := env.Repository.GetUserIDs(repository.UsersQuery{}.GMemberOf(g.ID))
		if assert.NoError(t, err) {
			assert.ElementsMatch(t, ids, []uuid.UUID{user.GetID()})
		}
	})
}

func TestHandlers_DeleteUserGroupMembers(t *testing.T) {
	t.Parallel()
	env, _, _, s, _, user, _ := setupWithUsers(t, common5)
	g := env.mustMakeUserGroup(t, rand, user.GetID())
	env.mustAddUserToGroup(t, user.GetID(), g.ID)
	user2 := env.mustMakeUser(t, rand)

	t.Run("NotLoggedIn", func(t *testing.T) {
		t.Parallel()
		e := env.makeExp(t)
		e.DELETE("/api/1.0/groups/{groupID}/members/{userID}", g.ID.String(), user.GetID().String()).
			Expect().
			Status(http.StatusUnauthorized)
	})

	t.Run("not found", func(t *testing.T) {
		t.Parallel()
		e := env.makeExp(t)
		e.DELETE("/api/1.0/groups/{groupID}/members/{userID}", uuid.Must(uuid.NewV4()), user.GetID().String()).
			WithCookie(session.CookieName, s).
			Expect().
			Status(http.StatusNotFound)
	})

	t.Run("unknown user", func(t *testing.T) {
		t.Parallel()
		e := env.makeExp(t)
		e.DELETE("/api/1.0/groups/{groupID}/members/{userID}", g.ID.String(), uuid.Must(uuid.NewV4())).
			WithCookie(session.CookieName, s).
			Expect().
			Status(http.StatusNoContent)
	})

	t.Run("forbidden", func(t *testing.T) {
		t.Parallel()
		e := env.makeExp(t)
		e.DELETE("/api/1.0/groups/{groupID}/members/{userID}", g.ID.String(), user.GetID().String()).
			WithCookie(session.CookieName, env.generateSession(t, user2.GetID())).
			Expect().
			Status(http.StatusForbidden)
	})

	t.Run("ok", func(t *testing.T) {
		t.Parallel()
		e := env.makeExp(t)
		e.DELETE("/api/1.0/groups/{groupID}/members/{userID}", g.ID.String(), user.GetID().String()).
			WithCookie(session.CookieName, s).
			Expect().
			Status(http.StatusNoContent)

		ids, err := env.Repository.GetUserIDs(repository.UsersQuery{}.GMemberOf(g.ID))
		if assert.NoError(t, err) {
			assert.Len(t, ids, 0)
		}
	})
}

func TestHandlers_GetMyBelongingGroup(t *testing.T) {
	t.Parallel()
	env, _, _, s, _, user, adminUser := setupWithUsers(t, common5)

	g1 := env.mustMakeUserGroup(t, rand, adminUser.GetID())
	g2 := env.mustMakeUserGroup(t, rand, adminUser.GetID())
	env.mustAddUserToGroup(t, user.GetID(), g1.ID)
	env.mustAddUserToGroup(t, user.GetID(), g2.ID)

	t.Run("NotLoggedIn", func(t *testing.T) {
		t.Parallel()
		e := env.makeExp(t)
		e.GET("/api/1.0/users/me/groups").
			Expect().
			Status(http.StatusUnauthorized)
	})

	t.Run("ok", func(t *testing.T) {
		t.Parallel()
		e := env.makeExp(t)
		e.GET("/api/1.0/users/me/groups").
			WithCookie(session.CookieName, s).
			Expect().
			Status(http.StatusOK).
			JSON().
			Array().
			ContainsOnly(g1.ID.String(), g2.ID.String())
	})
}

func TestHandlers_GetUserBelongingGroup(t *testing.T) {
	t.Parallel()
	env, _, _, s, _, _, adminUser := setupWithUsers(t, common5)

	user := env.mustMakeUser(t, rand)
	g1 := env.mustMakeUserGroup(t, rand, adminUser.GetID())
	g2 := env.mustMakeUserGroup(t, rand, adminUser.GetID())
	env.mustAddUserToGroup(t, user.GetID(), g1.ID)
	env.mustAddUserToGroup(t, user.GetID(), g2.ID)

	t.Run("NotLoggedIn", func(t *testing.T) {
		t.Parallel()
		e := env.makeExp(t)
		e.GET("/api/1.0/users/{userID}/groups", user.GetID().String()).
			Expect().
			Status(http.StatusUnauthorized)
	})

	t.Run("unknown user", func(t *testing.T) {
		t.Parallel()
		e := env.makeExp(t)
		e.GET("/api/1.0/users/{userID}/groups", uuid.Must(uuid.NewV4())).
			WithCookie(session.CookieName, s).
			Expect().
			Status(http.StatusNotFound)
	})

	t.Run("ok", func(t *testing.T) {
		t.Parallel()
		e := env.makeExp(t)
		e.GET("/api/1.0/users/{userID}/groups", user.GetID().String()).
			WithCookie(session.CookieName, s).
			Expect().
			Status(http.StatusOK).
			JSON().
			Array().
			ContainsOnly(g1.ID.String(), g2.ID.String())
	})
}
