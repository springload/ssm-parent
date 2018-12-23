package ssm

import (
	"encoding/json"
	"fmt"
	goPath "path"

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
			Config:            aws.Config{},
			SharedConfigState: session.SharedConfigEnable,
			Profile:           "",
		})
		if err != nil {
			return fmt.Errorf("Can't get aws session.")
		}
	}
	return nil
}

func collectJsonParameters(responseParameters []*ssm.Parameter) (parameters []map[string]string, errors []error) {
	for _, parameter := range responseParameters {
		value := make(map[string]string)
		if innerErr := json.Unmarshal([]byte(aws.StringValue(parameter.Value)), &value); innerErr != nil {
			errors = append(errors, fmt.Errorf("Can't unmarshal json from '%s': %s", aws.StringValue(parameter.Name), innerErr))
		} else {
			parameters = append(parameters, value)
		}
	}
	return
}

func getJsonSSMParametersByPaths(paths []string, strict, recursive bool) (parameters []map[string]string, err error) {
	err = makeSession()
	if err != nil {
		log.WithError(err).Fatal("Can't create session") // fail early here
	}
	s := ssm.New(localSession)
	for _, path := range paths {
		innerErr := s.GetParametersByPathPages(&ssm.GetParametersByPathInput{
			Path:           aws.String(path),
			WithDecryption: aws.Bool(true),
			Recursive:      aws.Bool(recursive),
		}, func(response *ssm.GetParametersByPathOutput, last bool) bool {
			innerParameters, errs := collectJsonParameters(response.Parameters)
			for _, parseErr := range errs {
				err = multierror.Append(err, parseErr)
			}
			parameters = append(parameters, innerParameters...)

			return true
		},
		)
		if innerErr != nil {
			err = multierror.Append(err, fmt.Errorf("Can't get parameters from path '%s': %s", path, innerErr))
		}
	}

	return
}

func getJsonSSMParameters(names []string, strict bool) (parameters []map[string]string, err error) {
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
				found = append(found, aws.StringValue(f.Name))
			}
			diff := stringSliceDifference(names, found)
			log.WithFields(log.Fields{"missing_names": diff}).Warn("Some parameters have not been found")
		}
	}
	innerParameters, errs := collectJsonParameters(response.Parameters)
	for _, parseErr := range errs {
		err = multierror.Append(err, parseErr)
	}
	parameters = append(parameters, innerParameters...)
	return
}

func collectPlainParameters(responseParameters []*ssm.Parameter) (parameters []map[string]string, errors []error) {
	for _, parameter := range responseParameters {
		values := make(map[string]string)
		values[goPath.Base(aws.StringValue(parameter.Name))] = aws.StringValue(parameter.Value)
		parameters = append(parameters, values)
	}
	return
}

func getPlainSSMParametersByPaths(paths []string, strict, recursive bool) (parameters []map[string]string, err error) {
	err = makeSession()
	if err != nil {
		log.WithError(err).Fatal("Can't create session") // fail early here
	}
	s := ssm.New(localSession)
	for _, path := range paths {
		innerErr := s.GetParametersByPathPages(&ssm.GetParametersByPathInput{
			Path:           aws.String(path),
			WithDecryption: aws.Bool(true),
			Recursive:      aws.Bool(recursive),
		}, func(response *ssm.GetParametersByPathOutput, last bool) bool {
			innerParameters, errs := collectPlainParameters(response.Parameters)
			for _, parseErr := range errs {
				err = multierror.Append(err, parseErr)
			}
			parameters = append(parameters, innerParameters...)

			return true
		},
		)
		if innerErr != nil {
			err = multierror.Append(err, fmt.Errorf("Can't get parameters from path '%s': %s", path, innerErr))
		}
	}
	return
}

func getPlainSSMParameters(names []string, strict bool) (parameters []map[string]string, err error) {
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
				found = append(found, aws.StringValue(f.Name))
			}
			diff := stringSliceDifference(names, found)
			log.WithFields(log.Fields{"missing_names": diff}).Warn("Some parameters have not been found")
		}
	}
	innerParameters, errs := collectPlainParameters(response.Parameters)
	for _, parseErr := range errs {
		err = multierror.Append(err, parseErr)
	}
	parameters = append(parameters, innerParameters...)
	return
}

func GetParameters(names, paths, plainNames, plainPaths []string, expand, strict, recursive bool) (parameters []map[string]string, err error) {
	localNames := names
	localPaths := paths
	localPlainNames := plainNames
	localPlainPaths := plainPaths
	if expand {
		localNames = ExpandArgs(names)
		localPaths = ExpandArgs(paths)
		localPlainNames = ExpandArgs(plainNames)
		localPlainPaths = ExpandArgs(plainPaths)
	}

	if len(localPaths) > 0 {
		parametersFromPaths, err := getJsonSSMParametersByPaths(localPaths, strict, recursive)
		if err != nil {
			log.WithError(err).WithFields(
				log.Fields{"paths": localPaths},
			).Fatal("Can't get parameters by paths")
		}
		for _, parameter := range parametersFromPaths {
			parameters = append(parameters, parameter)
		}
	}
	if len(localNames) > 0 {
		parametersFromNames, err := getJsonSSMParameters(localNames, strict)
		if err != nil {
			log.WithError(err).WithFields(
				log.Fields{"names": localNames},
			).Fatal("Can't get parameters by names")
		}
		for _, parameter := range parametersFromNames {
			parameters = append(parameters, parameter)
		}
	}
	if len(localPlainPaths) > 0 {
		parametersFromPlainPaths, err := getPlainSSMParametersByPaths(localPlainPaths, strict, recursive)
		if err != nil {
			log.WithError(err).WithFields(
				log.Fields{"plain_paths": localPaths},
			).Fatal("Can't get plain parameters by paths")
		}
		for _, parameter := range parametersFromPlainPaths {
			parameters = append(parameters, parameter)
		}
	}

	if len(localPlainNames) > 0 {
		parametersFromPlainNames, err := getPlainSSMParameters(localPlainNames, strict)
		if err != nil {
			log.WithError(err).WithFields(
				log.Fields{"plain_names": localPlainNames},
			).Fatal("Can't get plain parameters by names")
		}
		for _, parameter := range parametersFromPlainNames {
			parameters = append(parameters, parameter)
		}
	}
	return
}
