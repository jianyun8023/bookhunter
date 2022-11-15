package fetcher

import (
	"errors"
	"fmt"
	"io"

	"github.com/bookstairs/bookhunter/internal/driver"
)

// service is a real implementation for fetcher.
type service interface {
	// size is the total download amount of this given service.
	size() (int64, error)

	// formats will query the available downloadable file formats.
	formats(int64) (map[Format]driver.Share, error)

	// fetch the given book ID.
	fetch(int64, Format, driver.Share, io.Writer) error
}

// newService is the endpoint for creating all the supported download service.
func newService(c *Config) (service, error) {
	switch c.Category {
	case Talebook:
		return newTalebookService(c)
	case SanQiu:
		return newSanqiuService(c)
	case SoBooks:
		// TODO We are working on this feature now.
		return nil, errors.New("we don't support sobooks now")
	case TianLang:
		// TODO We are working on this feature now.
		return nil, errors.New("we don't support tianlang now")
	default:
		return nil, fmt.Errorf("no such fetcher service [%s] supported", c.Category)
	}
}
