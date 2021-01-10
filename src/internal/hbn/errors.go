package hbn

import (
		"fmt"
		"errors"
)

func ErrMethodDoesNotAllowed(method string) error {
		err := errors.New(fmt.Sprintf("method %s does not allowed to use", method))
		return err
}

func ErrUseragentsWasNotFound(path string) error {
		err := errors.New(fmt.Sprintf("useragents wasn't fount on this path: "+path))
		return err
}

func ErrConfigWasNotFound(path string) error {
		err := errors.New(fmt.Sprintf("config file wasn't found on this path: "+path))
		return err
}

func ErrReadFile(file_name string) error {
		err := errors.New(fmt.Sprintf("error while reading %s", file_name))
		return err
}

func ErrConfigRead() error {
		err := errors.New(fmt.Sprintf("error while reading config, maybe your config isn't valid"))
		return err
}
