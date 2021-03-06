// See https://cloud.google.com/docs/authentication/.
// Use GOOGLE_APPLICATION_CREDENTIALS environment variable to specify
// a service account key file to authenticate to the API.

package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/kyma-project/test-infra/development/tools/pkg/release"
	log "github.com/sirupsen/logrus"
)

var (
	targetCommit           = flag.String("targetCommit", "", "Target commitish [Required]")
	bucketName             = flag.String("bucketName", "kyma-prow-artifacts", "Google bucket name where artifacts are stored [Optional]")
	kymaInstaller          = flag.String("kymaInstaller", "kyma-installer.yaml", "Filename for list of CRDs and deployment artifact [Optional]")
	kymaConfigLocal        = flag.String("kymaConfigLocal", "kyma-config-local.yaml", "Filename for local config artifact [Optional]")
	kymaInstallerCRLocal   = flag.String("kymaInstallerCRLocal", "kyma-installer-cr-local.yaml", "Filename for list of componets kyma installer would install for local [Optional]")
	kymaInstallerCRCluster = flag.String("kymaInstallerCRCluster", "kyma-installer-cr-cluster.yaml", "Filename for list of componets kyma installer would install for cluster [Optional]")
	kymaChangelog          = flag.String("kymaChangelog", "release-changelog.md", "Filename for release changelog [Optional]")
	githubRepoOwner        = flag.String("githubRepoOwner", "", "Github repository owner [Required]")
	githubRepoName         = flag.String("githubRepoName", "", "Github repository name [Required]")
	githubAccessToken      = flag.String("githubAccessToken", "", "Github access token [Required]")
	releaseVersionFilePath = flag.String("releaseVersionFilePath", "", "Full path to a file containing release version [Required]")
)

func main() {
	flag.Parse()

	if *targetCommit == "" {
		fmt.Fprintln(os.Stderr, "missing -targetCommit flag")
		flag.Usage()
		os.Exit(2)
	}

	if *githubRepoOwner == "" {
		fmt.Fprintln(os.Stderr, "missing -githubRepoOwner flag")
		flag.Usage()
		os.Exit(2)
	}

	if *githubRepoName == "" {
		fmt.Fprintln(os.Stderr, "missing -githubRepoName flag")
		flag.Usage()
		os.Exit(2)
	}

	if *githubAccessToken == "" {
		fmt.Fprintln(os.Stderr, "missing -githubAccessToken flag")
		flag.Usage()
		os.Exit(2)
	}

	if *releaseVersionFilePath == "" {
		fmt.Fprintln(os.Stderr, "missing -releaseVersionFilePath flag")
		flag.Usage()
		os.Exit(2)
	}

	ctx := context.Background()

	ga := release.NewGithubAPI(ctx, *githubAccessToken, *githubRepoOwner, *githubRepoName)

	sa, err := release.NewStorageAPI(ctx, *bucketName)
	if err != nil {
		log.Fatal(err)
	}

	c := release.NewCreator(ga, sa)

	relOpts, err := release.NewOptions(ctx, sa, *releaseVersionFilePath, *kymaChangelog, *targetCommit, nil)
	if err != nil {
		log.Fatal(err)
	}

	// Github release
	err = c.CreateNewRelease(ctx, relOpts, *kymaConfigLocal, *kymaInstallerCRLocal, *kymaInstallerCRCluster, *kymaInstaller)
	if err != nil {
		log.Fatal(err)
	}
}
