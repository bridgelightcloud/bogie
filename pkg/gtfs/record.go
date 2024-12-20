package gtfs

import (
	"fmt"
	"io"

	"github.com/bridgelightcloud/bogie/pkg/csvmum"
)

type record interface {
	key() string
	validate() errorList
}

func parse[T record](f io.Reader, records map[string]T, errors *errorList) {
	csvm, err := csvmum.NewUnmarshaler[T](f)
	if err != nil {
		errors.add(fmt.Errorf("error creating unmarshaler for file: %w", err))
		return
	}

	for {
		var r T

		err = csvm.Unmarshal(&r)
		if err == io.EOF {
			break
		}
		if err != nil {
			errors.add(fmt.Errorf("error unmarshalling file: %w", err))
			break
		}

		errs := r.validate()
		if errs != nil {
			fmt.Println("errors", errs)
			for _, e := range errs {
				errors.add(fmt.Errorf("invalid record: %w", e))
			}
			continue
		}

		if _, ok := records[r.key()]; ok {
			errors.add(fmt.Errorf("duplicate key: %s", r.key()))
			continue
		}

		records[r.key()] = r
	}
}
