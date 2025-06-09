package test

import (
	"context"
	"net/http"
	"selfhosted/database"
	"selfhosted/database/store"
	"testing"
)

func TestUserMenu(t *testing.T) {
	t.Run("guests cannot access user modal", func(t *testing.T) {
		tc := NewTestCase(t)
		defer tc.Close()

		tc.SetupAdmin()

		tc.Get("/user")
		tc.AssertRedirect(http.StatusSeeOther, "/login")
	})

	t.Run("users can access user modal", func(t *testing.T) {
		tc := NewTestCase(t)
		defer tc.Close()

		tc.SetupAdmin()

		tc.LogIn("admin@example.com", "password123")
		tc.Get("/user")
		tc.AssertStatus(http.StatusOK)
		tc.AssertNoRedirect()
		tc.AssertElementVisible("nav")
	})

	t.Run("user can see all their teams", func(t *testing.T) {
		tc := NewTestCase(t)
		defer tc.Close()

		adminId := tc.SetupAdmin()

		team1Id, _ := store.New(database.DB).CreateTeam(context.Background(), store.CreateTeamParams{
			Name: "Team #1",
		})
		store.New(database.DB).AddMemberToTeam(context.Background(), store.AddMemberToTeamParams{
			UserID: adminId,
			TeamID: team1Id,
		})

		team2Id, _ := store.New(database.DB).CreateTeam(context.Background(), store.CreateTeamParams{
			Name: "Team #2",
		})
		store.New(database.DB).AddMemberToTeam(context.Background(), store.AddMemberToTeamParams{
			UserID: adminId,
			TeamID: team2Id,
		})

		store.New(database.DB).CreateTeam(context.Background(), store.CreateTeamParams{
			Name: "Team #3",
		})

		tc.LogIn("admin@example.com", "password123")

		tc.Get("/user")
		tc.AssertStatus(http.StatusOK)
		tc.AssertSee("Team #1")
		tc.AssertSee("Team #2")
		tc.AssertNotSee("Team #3")
	})
}
