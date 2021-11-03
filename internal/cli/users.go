package cli

import (
	"fmt"
	distrybute "github.com/mmichaelb/distrybute/pkg"
	"github.com/mmichaelb/distrybute/pkg/postgresminio"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
	"time"
)

var service *postgresminio.Service

var usernameFlag = &cli.StringFlag{Name: "username", Aliases: []string{"u"}, Required: true}

var userCommand = &cli.Command{
	Name:    "user",
	Aliases: []string{"u"},
	Usage:   "manage users",
	Subcommands: []*cli.Command{
		{
			Name:   "create",
			Usage:  "create a new distrybute user",
			Action: createUser,
			Flags: []cli.Flag{
				usernameFlag,
				&cli.StringFlag{Name: "password", Aliases: []string{"p"}, Required: true},
			},
		},
		{
			Name:   "delete",
			Usage:  "delete a distrybute user",
			Action: deleteUser,
			Flags: []cli.Flag{
				&cli.StringFlag{Name: "username", Aliases: []string{"u"}},
			},
		},
		{
			Name:   "list",
			Usage:  "list all distrybute users",
			Action: listUsers,
		},
	},
}

func createUser(c *cli.Context) error {
	username := c.String("username")
	password := []byte(c.String("password"))
	log.Info().Msg("creating new user...")
	user, err := service.CreateNewUser(username, password)
	if err != nil {
		log.Err(err).Msg("could not create new user")
		return err
	}
	log.Info().Str("username", username).Str("id", user.ID.String()).Str("auth_token", user.AuthorizationToken).Msg("created user in databas!")
	return nil
}

func deleteUser(c *cli.Context) error {
	username := c.String("username")
	start := time.Now()
	log.Info().Str("username", username).Msg("searching for user...")
	user, err := service.GetUserByUsername(username)
	if err == distrybute.ErrUserNotFound {
		log.Err(err).Str("username", username).Msg("the specified user could not be found")
		return err
	} else if err != nil {
		return err
	}
	log.Info().Str("username", user.Username).Str("id", user.ID.String()).Msg("found user")
	err = service.DeleteUser(user.ID)
	if err != nil {
		log.Err(err).Msg("could not delete user")
		return err
	}
	end := time.Now()
	log.Info().Dur("duration", end.Sub(start)).Msg("successfully deleted user")
	return nil
}

func listUsers(_ *cli.Context) error {
	log.Info().Msg("listing users...")
	users, err := service.ListUsers()
	if err != nil {
		log.Err(err).Msg("could not request user list")
		return err
	}
	format := "%-32s | %-36s"
	log.Info().Msg(fmt.Sprintf(format, "Username", "ID"))
	for _, user := range users {
		log.Info().Msg(fmt.Sprintf(format, user.Username, user.ID.String()))
	}
	log.Info().Msg("done with user list")
	return nil
}
