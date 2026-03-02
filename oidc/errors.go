package oidc

type RegistryError struct {
	Message string
}

func (e *RegistryError) Error() string {
	return e.Message
}

type SyncError struct {
	Err error
}

func (e *SyncError) Error() string {
	return e.Err.Error()
}

func (e *SyncError) Unwrap() error {
	return e.Err
}

type ExchangeError struct {
	Err error
}

func (e *ExchangeError) Error() string {
	return e.Err.Error()
}

func (e *ExchangeError) Unwrap() error {
	return e.Err
}
