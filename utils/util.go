package utils

// NOT threadsafe
func ChainErr(lambdas ...func() error) error {
	for _, lambda := range lambdas {
		var err = lambda()
		if err != nil {
			return err
		}
	}

	return nil
}
