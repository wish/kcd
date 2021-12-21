package harbor

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/heroku/docker-registry-client/registry"
	"github.com/pkg/errors"
	kcdRegistry "github.com/wish/kcd/registry"
)

// Default Wish Harbor url
const harborUrl = "https://harbor.s.wish.site/"

type Options struct {
	HarborUrl, User, Password string
}

// Allows for base Harbor url to be overridden
func WithUrl(url string) func(*Options) {
	return func(opts *Options) {
		opts.HarborUrl = url
	}
}

// Allows for username/password based credentials to be passed
func WithCreds(user, password string) func(*Options) {
	return func(opts *Options) {
		opts.User = user
		opts.Password = password
	}
}

type Provider struct {
	registry   string
	repository string
	client     *registry.Registry
	opts       *Options
}

// Creates a new Harbor registry provider
func NewHarbor(imageUrl, versionExp string, options ...func(*Options)) (*Provider, error) {
	// Defaults
	opts := &Options{
		User:      "",
		Password:  "",
		HarborUrl: harborUrl,
	}
	for _, opt := range options {
		opt(opts)
	}

	client, err := registry.New(opts.HarborUrl, opts.User, opts.Password)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to connect to Harbor")
	}

	reg, repo, err := parseRegRepoFromUrl(imageUrl)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to parse repository from image url: %s", imageUrl)
	}
	fmt.Printf("parsed registry: %s, repository: %s\n", reg, repo)

	return &Provider{
		client:     client,
		registry:   reg,
		repository: repo,
	}, nil
}

// Implements the kcdRegistry.Provider interface
func (p *Provider) RegistryFor(imageUrl string) (kcdRegistry.Registry, error) {
	reg, repo, err := parseRegRepoFromUrl(imageUrl)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to parse repository from image url: %s", imageUrl)
	}

	return &Provider{
		client:     p.client,
		registry:   reg,
		repository: repo,
	}, nil
}

// Implements the kcdRegistry.Registry interface
func (p *Provider) Versions(ctx context.Context, tag string) ([]string, error) {
	tags := make([]string, 0, 5)
	version, err := p.getDigest(tag)
	if err != nil {
		return tags, errors.Errorf("No version found for tag %s", tag)
	}
	tags = append(tags, version)
	return tags, nil
}

// Implements the kcdRegistry.Tagger interface. Fetches image manifest of specified version
// and adds additional tags to the manifest
func (p *Provider) Add(version string, tags ...string) error {
	manifest, err := p.client.Manifest(p.repository, version)
	if err != nil {
		return errors.Wrapf(err, "Failed to find manifest for image version %s on repository %s", version, p.repository)
	}

	for _, tag := range tags {
		err := p.client.PutManifest(p.repository, tag, manifest)
		if err != nil {
			return errors.Wrapf(err, "Failed to add tags %s on image version %s on repository %s", tag, version, p.repository)
		}
	}
	return nil
}

// Implements the kcdRegistry.Tagger interface. Currently not supported
func (p *Provider) Remove(tags ...string) error {
	return errors.New("Not supported")
}

// Implements the kcdRegistry.Tagger interface. Gets list of tags on the image with given version tag
func (p *Provider) Get(version string) ([]string, error) {
	fmt.Printf("Getting version %s\n", version)
	digest, err := p.getDigest(version)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to get tags")
	}
	return []string{version, digest}, nil
}

// Gets the image digest of a tagged docker image
func (p *Provider) getDigest(tag string) (string, error) {
	fmt.Printf("Getting digest for %s on %s\n", tag, p.repository)
	digest, err := p.client.ManifestDigest(p.repository, tag)
	if err != nil {
		return "", errors.Wrapf(err, "Failed to get tag %s on repository %s", tag, p.repository)
	}
	fmt.Printf("Found digest: %s", digest.String())
	return digest.String(), nil
}

func parseRegRepoFromUrl(imageUrl string) (string, string, error) {
	// Repository needs to be parsed from the image url. Assumes image is of form: harbor.s.wish.site/<namespace>/<image>
	u, err := url.Parse(imageUrl)
	if err != nil || u.Host == "" {
		// url.Parse expects scheme
		u, err = url.Parse("https://" + imageUrl)
		if err != nil {
			return "", "", errors.Wrapf(err, "Failed to parse repository from image url: %s", imageUrl)
		}
	}

	return u.Host, strings.TrimLeft(u.Path, "/"), nil
}
