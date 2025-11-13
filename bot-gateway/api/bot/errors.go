package bot

import (
	"fmt"
)

func errCantCreate(err error) error {
	return fmt.Errorf("cannot create bot: %s", err.Error())
}
