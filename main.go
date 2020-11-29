package main

import (
	"context"
	"fmt"
	"github.com/felipesere/probe/cmd"
	"github.com/felipesere/probe/lib"
	"github.com/shurcooL/githubv4"
	"github.com/spf13/viper"
	"golang.org/x/oauth2"
	"os"
)

func main() {
	paths := []string{"$HOME", "."}
	viper.SetConfigName(".probe")
	viper.SetConfigType("yml")
	for _, path := range paths {
		viper.AddConfigPath(path)
	}
	viper.SetEnvPrefix("")
	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			fmt.Fprintln(os.Stderr, "Config file was not found")
			os.Exit(1)
		} else {
			// Config file was found but another error was produced
		}
	}

	client := *githubv4.NewClient(oauth2.NewClient(context.Background(), oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: viper.GetString("github_token")},
	)))

	db, err := lib.NewStorage(viper.GetString("database_path"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not load DB: %s\n", err.Error())
		os.Exit(1)
	}
	cmd.Execute(client, db)
	db.Flush()
}
