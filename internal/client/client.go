package client

import (
	"context"

	eg "github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/edgegrid"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/session"
	"github.com/pkg/errors"
	"gopkg.in/ini.v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"

	xpv1 "github.com/crossplane/crossplane-runtime/apis/common/v1"
	"github.com/crossplane/crossplane-runtime/pkg/resource"

	"github.com/crossplane-contrib/provider-akamai/apis/v1alpha1"
)

type Config struct {
	Host         string `ini:"host"`
	ClientToken  string `ini:"client_token"`
	ClientSecret string `ini:"client_secret"`
	AccessToken  string `ini:"access_token"`
}

func NewAkamaiClient(c Config) *session.Session {
	edgerc := eg.Must(eg.New(
		eg.WithEnv(true),
	))

	s, err := session.New(
		session.WithSigner(edgerc),
	)
	if err != nil {
		panic(err)
	}

	return &s
}

// UseProviderConfig to produce a config that can be used to authenticate to Akamai.
func UseProviderConfig(ctx context.Context, c client.Client, mg resource.Managed) (*Config, error) {
	// var (
	// 	requiredOptions = []string{"host", "client_token", "client_secret", "access_token"}
	// )
	pc := &v1alpha1.ProviderConfig{}
	if err := c.Get(ctx, types.NamespacedName{Name: mg.GetProviderConfigReference().Name}, pc); err != nil {
		return nil, errors.Wrap(err, "cannot get referenced Provider")
	}

	t := resource.NewProviderConfigUsageTracker(c, &v1alpha1.ProviderConfigUsage{})
	if err := t.Track(ctx, mg); err != nil {
		return nil, errors.Wrap(err, "cannot track ProviderConfig usage")
	}

	switch s := pc.Spec.Credentials.Source; s { //nolint:exhaustive
	case xpv1.CredentialsSourceSecret:
		csr := pc.Spec.Credentials.SecretRef
		if csr == nil {
			return nil, errors.New("no credentials secret referenced")
		}
		s := &corev1.Secret{}
		if err := c.Get(ctx, types.NamespacedName{Namespace: csr.Namespace, Name: csr.Name}, s); err != nil {
			return nil, errors.Wrap(err, "cannot get credentials secret")
		}

		edgerc, err := ini.Load(s.Data[csr.Key])
		if err != nil {
			return nil, errors.Wrap(err, "cannot parse credentials secret")
		}

		sec, err := edgerc.GetSection("default")
		if err != nil {
			return nil, errors.Wrap(err, "default section does not exist")
		}

		config := &Config{}
		err = sec.MapTo(&config)
		if err != nil {
			return nil, err
		}

		// for _, opt := range requiredOptions {
		// 	if !(edgerc.Section("default").HasKey(opt)) {
		// 		return nil, errors.Wrap(err, "required key not found")
		// 	}
		// }

		return &Config{
			Host:         config.Host,
			ClientSecret: config.ClientSecret,
			AccessToken:  config.AccessToken,
			ClientToken:  config.ClientToken,
		}, nil
	default:
		return nil, errors.Errorf("credentials source %s is not currently supported", s)
	}
}
