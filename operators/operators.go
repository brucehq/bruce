package operators

import (
	"bruce/system"
	"bytes"
	"fmt"
	"github.com/rs/zerolog/log"
	"os"
	"strings"
	"text/template"
)

type Operator interface {
	Execute() error
}

type NullOperator struct {
}

func (n *NullOperator) Execute() error {
	return fmt.Errorf("invalid operator")
}

func GetValueForOSHandler(value string) string {
	log.Debug().Msgf("OS Handler value iteration: %#v", value)
	if system.Get().PackageHandler == "" {
		log.Error().Err(fmt.Errorf("cannot retrieve os handler value without a known package handler"))
		return ""
	}
	log.Debug().Msgf("testing for my package handler: %s", system.Get().PackageHandler)
	if strings.Contains(value, "|") {
		managerList := strings.Split(value, "|")
		var basePackage = ""
		var usablePackage = ""
		for _, mpkg := range managerList {
			log.Debug().Msgf("os handler iteration for manager: %#v", mpkg)
			if strings.Contains(mpkg, "=") {
				pmSplit := strings.Split(mpkg, "=")
				log.Debug().Msgf("handler [%s] specific value: %s", pmSplit[0], pmSplit[1])
				if pmSplit[0] == system.Get().PackageHandler {
					usablePackage = pmSplit[1]
				}
			} else {
				basePackage = mpkg
			}
		}
		if usablePackage != "" {
			log.Debug().Msgf("returning package manager value: %s", usablePackage)
			return usablePackage
		}
		log.Debug().Msgf("returning base value: %s", basePackage)
		return basePackage
	}
	log.Debug().Msgf("returning original value: %s", value)
	return value
}

func RenderEnvString(o string) string {
	// now we iterate over the steps and replace the --== and ==-- with {{ and }} respectively and store it in s
	s := strings.ReplaceAll(o, "--==", "{{")
	s = strings.ReplaceAll(s, "==--", "}}")
	log.Debug().Msgf("rendering env string: %s", s)
	envVars := make(map[string]string)
	for _, e := range os.Environ() {
		pair := strings.SplitN(e, "=", 2)
		envVars[pair[0]] = pair[1]
	}

	s, err := StringFromTemplate(s, envVars)
	if err != nil {
		log.Error().Err(err).Msg("failed to render env string")
	}
	return s
}

// StringFromTemplate applies a map of fields to a template string and returns the resulting string.
func StringFromTemplate(InTemplate string, fields map[string]string) (string, error) {
	tmpl, err := template.New("stringTemplate").Parse(InTemplate)
	if err != nil {
		return "", err
	}

	var result bytes.Buffer
	err = tmpl.Execute(&result, fields)
	if err != nil {
		return "", err
	}

	return result.String(), nil
}
