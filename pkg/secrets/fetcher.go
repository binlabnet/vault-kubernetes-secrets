package secrets

import (
	"github.com/spf13/viper"
	"github.com/nais/vault-kubernetes-secrets/pkg/vault"
	"io/ioutil"
	"fmt"
)

const (
	envSecretsDestPath = "VS_SECRET_DEST_PATH"
	envJwtFile         = "VS_SERVICE_ACCOUNT_TOKEN"
)

func init() {
	viper.BindEnv(envJwtFile)
	viper.BindEnv(envSecretsDestPath)
	viper.SetDefault(envSecretsDestPath, "/var/run/secrets/naisd.io/vault")
	viper.SetDefault(envJwtFile, "/var/run/secrets/kubernetes.io/serviceaccount/token")
}

func New(options vault.ClientOptions) SecretFetcher {
	return secretFetcher{
		auth:         vault.NewAuthClient(options),
		kv:           vault.NewKVClient(options),
		jwtRetriever: jwtFromFile(viper.GetString(envJwtFile)),
		secretWriter: writeToFile(viper.GetString(envSecretsDestPath)),
	}
}

type SecretFetcher interface {
	FetchSecrets(role, authPath, kvPath string) error
}

type secretFetcher struct {
	auth         vault.Auth
	kv           vault.KV
	jwtRetriever func() (string, error)
	secretWriter func(map[string]string) error
}

func (s secretFetcher) FetchSecrets(role, authPath, kvPath string) error {
	jwt, jwtError := s.jwtRetriever()
	if jwtError != nil {
		return jwtError
	}

	accessToken, loginError := s.auth.LoginK8s(role, jwt, authPath)
	if loginError != nil {
		return loginError
	}

	secrets, fetchError := s.kv.Get(kvPath, accessToken)
	if fetchError != nil {
		return fetchError
	}

	return s.secretWriter(secrets)
}

func jwtFromFile(jwtFile string) func() (token string, err error) {
	return func() (token string, err error) {
		if b, e := ioutil.ReadFile(jwtFile); e != nil {
			return "", e
		} else {
			return string(b), nil
		}
	}
}

func writeToFile(destDir string) func(secrets map[string]string) error {
	return func(secrets map[string]string) error {
		for k, v := range secrets {
			dest := destDir + "/" + k
			if err := ioutil.WriteFile(dest, []byte(v), 0644); err != nil {
				return fmt.Errorf("Fail to write secret %s to  %s. Error: ", k, err.Error())
			}
		}
		return nil
	}
}
