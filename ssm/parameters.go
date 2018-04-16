package ssm

import (
	"encoding/json"
	"fmt"

	"github.com/apex/log"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
	multierror "github.com/hashicorp/go-multierror"
)

var localSession *session.Session

func makeSession() error {
	if localSession == nil {
		log.Debug("Creating session")
		var err error
		// create AWS session
		localSession, err = session.NewSessionWithOptions(session.Options{
			Config:            aws.Config{}, //Region: aws.String(config.Region)},
			SharedConfigState: session.SharedConfigEnable,
			Profile:           "", //config.Profile,
		})
		if err != nil {
			return fmt.Errorf("Can't get aws session.")
		}
	}
	return nil
}

func getJsonSSMParametersByPaths(paths []string, strict, recursive bool) (parameters []map[string]interface{}, err error) {
	err = makeSession()
	if err != nil {
		log.WithError(err).Fatal("Can't create session") // fail early here
	}
	s := ssm.New(localSession)
	for _, path := range paths {
		response, innerErr := s.GetParametersByPath(&ssm.GetParametersByPathInput{
			Path:           aws.String(path),
			WithDecryption: aws.Bool(true),
			Recursive:      aws.Bool(recursive),
		})
		if innerErr != nil {
			err = multierror.Append(err, fmt.Errorf("Can't get parameters from path '%s': %s", path, innerErr))
		}
		for _, parameter := range response.Parameters {
			value := make(map[string]interface{})
			innerErr := json.Unmarshal([]byte(*parameter.Value), &value)
			if innerErr != nil {
				err = multierror.Append(err, fmt.Errorf("Can't unmarshal json from '%s': %s", *parameter.Name, innerErr))
			}
			parameters = append(parameters, value)
		}

	}

	return
}

func getJsonSSMParameters(names []string, strict bool) (parameters []map[string]interface{}, err error) {
	err = makeSession()
	if err != nil {
		log.WithError(err).Fatal("Can't create session") // fail early here
	}
	s := ssm.New(localSession)
	response, err := s.GetParameters(&ssm.GetParametersInput{
		Names:          aws.StringSlice(names),
		WithDecryption: aws.Bool(true),
	})
	if err != nil {
		return nil, err
	}
	if len(response.Parameters) < len(names) {
		if strict {
			err = multierror.Append(err, fmt.Errorf("Found %d parameters from %d names", len(response.Parameters), len(names)))
		} else {
			var found []string
			for _, f := range response.Parameters {
				found = append(found, *f.Name)
			}
			diff := stringSliceDifference(names, found)
			log.WithFields(log.Fields{"missing_names": diff}).Warn("Some parameters have not been found")
		}
	}
	for _, parameter := range response.Parameters {
		value := make(map[string]interface{})
		innerErr := json.Unmarshal([]byte(*parameter.Value), &value)
		if innerErr != nil {
			err = multierror.Append(err, fmt.Errorf("Can't unmarshal json from '%s': %s", *parameter.Name, innerErr))
		}
		parameters = append(parameters, value)
	}
	return
}

func GetParameters(names, paths []string, strict, recursive bool) (parameters []map[string]interface{}, err error) {
	if len(names) > 0 {
		parametersFromNames, err := getJsonSSMParameters(names, strict)
		if err != nil {
			log.WithError(err).WithFields(
				log.Fields{"names": names},
			).Fatal("Can't get parameters by names")
		}
		for _, parameter := range parametersFromNames {
			parameters = append(parameters, parameter)
		}
	}
	if len(paths) > 0 {
		parametersFromPaths, err := getJsonSSMParametersByPaths(paths, strict, recursive)
		if err != nil {
			log.WithError(err).WithFields(
				log.Fields{"paths": paths},
			).Fatal("Can't get parameters by paths")
		}
		for _, parameter := range parametersFromPaths {
			parameters = append(parameters, parameter)
		}
	}
	return
}
