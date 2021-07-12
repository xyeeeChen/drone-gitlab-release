package main

import (
	"fmt"
	"log"
	"net/url"
	"os"

	"github.com/urfave/cli"
	gitlab "github.com/xanzy/go-gitlab"
)

const (
	apiUrl      = "https://%s/api/v4/"
	registryUrl = "https://%s/"
)

func main() {
	app := cli.NewApp()
	app.Name = "Gitlab Release"
	app.Usage = "Create a release to gitlab"
	app.Action = action
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:     "access_token",
			Usage:    "Gitlab access token",
			EnvVar:   "DRONE_ACCESS_TOKEN",
			Required: true,
		},
		cli.StringFlag{
			Name:     "domain",
			Usage:    "Gitlab domain name",
			EnvVar:   "DRONE_DOMAIN",
			Required: true,
		},
		cli.StringFlag{
			Name:     "repo",
			Usage:    "The gitlab project ID or URL-encoded path of the project",
			EnvVar:   "DRONE_REPO",
			Required: true,
		},
		cli.StringFlag{
			Name:     "release",
			Usage:    "The release name",
			EnvVar:   "DRONE_RELEASE",
			Required: true,
		},
		cli.StringFlag{
			Name:     "tag",
			Usage:    "The tag name",
			EnvVar:   "DRONE_TAG",
			Required: true,
		},
		cli.StringFlag{
			Name:     "description",
			Usage:    "The description of the release that supports Markdown",
			EnvVar:   "DRONE_DESCRIPTION",
			Required: true,
		},
		cli.StringFlag{
			Name:     "ref",
			Usage:    "The release is created from ref and tagged with tag_name. It can be a commit SHA, another tag name, or a branch name.",
			EnvVar:   "DRONE_REF",
			Required: true,
		},
		cli.StringSliceFlag{
			Name:   "assets",
			Usage:  "Optional path for a direct asset link.",
			EnvVar: "DRONE_ASSETS",
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatalln(err)
	}
}

func action(c *cli.Context) error {
	api, err := url.Parse(fmt.Sprintf(apiUrl, c.String("domain")))
	if err != nil {
		return err
	}
	registry, err := url.Parse(fmt.Sprintf(registryUrl, c.String("domain")))
	if err != nil {
		return err
	}

	client, err := login(c.String("access_token"), api.String())
	if err != nil {
		return err
	}

	err = releaseExist(client, c.String("repo"), c.String("tag"))
	if err != nil {
		return err
	}

	assets := c.StringSlice("assets")
	var assetLinks []*gitlab.ReleaseAssetLink
	if len(assets) > 0 {
		assetLinks, err = uploadAssets(client, c.String("repo"), registry.String(), assets)
		if err != nil {
			return err
		}
	}

	return createRelease(client, assetLinks, c)
}

func login(token, api string) (*gitlab.Client, error) {
	client, err := gitlab.NewClient(token, gitlab.WithBaseURL(api))
	if err != nil {
		return nil, err
	}

	user, _, err := client.Users.CurrentUser()
	if err != nil {
		return nil, err
	}
	log.Println("Login with", user.Name)
	return client, err
}

func releaseExist(client *gitlab.Client, repo, tag string) error {
	rel, _, _ := client.Releases.GetRelease(repo, tag)
	if rel != nil {
		return fmt.Errorf("Error: The release %s exists.", tag)
	}
	return nil
}

func uploadAssets(client *gitlab.Client, repo, registry string, assets []string) ([]*gitlab.ReleaseAssetLink, error) {
	log.Println("Uploading assets...")

	var assetLinks []*gitlab.ReleaseAssetLink

	for _, asset := range assets {
		log.Println(fmt.Sprintf("%-20s: %s", "Uploading asset", asset))

		projectFile, _, err := client.Projects.UploadFile(repo, asset)
		if err != nil {
			return nil, fmt.Errorf("Upload Error: %w", err)
		}

		log.Println(fmt.Sprintf("%-20s: %s", "Done", projectFile.URL))

		assetURL := fmt.Sprintf("%s%s%s", registry, repo, projectFile.URL)
		var ral = gitlab.ReleaseAssetLink{Name: projectFile.Alt, URL: assetURL}
		assetLinks = append(assetLinks, &ral)
	}

	log.Println("Upload successful.")
	return assetLinks, nil
}

func createRelease(client *gitlab.Client, assetLinks []*gitlab.ReleaseAssetLink, c *cli.Context) error {
	repo := c.String("repo")
	description := c.String("description")
	tag := c.String("tag")
	release := c.String("release")
	ref := c.String("ref")

	opts := &gitlab.CreateReleaseOptions{
		Description: &description,
		TagName:     &tag,
		Name:        &release,
		Assets:      &gitlab.ReleaseAssets{Links: assetLinks},
		Ref:         &ref,
	}
	_, _, err := client.Releases.CreateRelease(repo, opts)
	if err != nil {
		return err
	}

	log.Println("The release", repo, tag, "is created.")
	return nil
}
