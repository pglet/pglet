package model

import (
	"errors"
	"fmt"
	"strings"

	"github.com/pglet/pglet/internal/utils"
)

const (
	publicNamespace = "public"
)

type PageName struct {
	Namespace string
	Name      string
}

func ParsePageName(pageName string) (*PageName, error) {

	p := &PageName{}
	p.Name = strings.ToLower(strings.Trim(strings.ReplaceAll(pageName, "\\", "/"), "/"))

	if strings.Count(p.Name, "/") > 1 {
		return nil, errors.New("Page name must be in format {page} or {namespace}/{page}")
	}

	if strings.Count(p.Name, "/") == 1 {
		// namespace specified
		parts := strings.Split(p.Name, "/")
		p.Namespace = parts[0]
		p.Name = parts[1]
	} else {
		p.Namespace = publicNamespace
	}

	rndText, err := utils.GenerateRandomString(12)
	if err != nil {
		return nil, err
	}

	p.Name = strings.ReplaceAll(p.Name, "*", rndText)

	return p, nil
}

func (pn *PageName) String() string {
	return fmt.Sprintf("%s/%s", pn.Namespace, pn.Name)
}
