package fb

import (
	"net/http"

	"github.com/facebookgo/fbapi"
	"github.com/pkg/errors"
)

type (
	paginate struct {
		req     *http.Request
		hasNext bool
		result  struct {
			Data   interface{} `json:"data"`
			Paging struct {
				Next string `json:"next"`
				Prev string `json:"prev"`
			} `json:"paging"`
		}
		client *fbapi.Client
	}
)

func newPaginate(req *http.Request, client *fbapi.Client) *paginate {
	var p paginate
	p.req = req
	p.client = client
	p.hasNext = true
	return &p
}

func (p *paginate) Next() bool {
	return p.hasNext
}

func (p *paginate) Data(res interface{}) error {
	p.result.Data = res
	_, err := p.client.Do(p.req, &p.result)
	if err != nil {
		return errors.WithStack(err)
	}
	if next := p.result.Paging.Next; next != "" {
		p.req, err = http.NewRequest(http.MethodGet, next, nil)
		if err != nil {
			return errors.WithStack(err)
		}
		p.hasNext = true
	} else {
		p.hasNext = false
	}
	return nil
}

func (r *repository) paginate(req *http.Request) *paginate {
	return newPaginate(req, r.client)
}
