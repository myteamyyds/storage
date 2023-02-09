package objects

import (
	"fmt"
	"io"
	"storage/objectStream"
	"storage/part3/apiServer/locate"
)

func getStream(object string) (io.Reader, error) {
	server := locate.Locate(object)
	if server == "" {
		return nil, fmt.Errorf("object %s locate fail", object)
	}
	return objectStream.NewGetStream(server, object)
}
