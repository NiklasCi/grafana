package clients

import (
	"context"

	"github.com/grafana/grafana/pkg/services/authn"
	"github.com/grafana/grafana/pkg/util/errutil"
	"github.com/grafana/grafana/pkg/web"
)

var (
	errBadForm = errutil.NewBase(errutil.StatusBadRequest, "form-auth.invalid", errutil.WithPublicMessage("bad login data"))
)

var _ authn.Client = new(Form)

func ProvideForm(client authn.PasswordClient) *Form {
	return &Form{client}
}

type Form struct {
	client authn.PasswordClient
}

type loginForm struct {
	Username string `json:"user" binding:"Required"`
	Password string `json:"password" binding:"Required"`
}

func (f *Form) Authenticate(ctx context.Context, r *authn.Request) (*authn.Identity, error) {
	form := loginForm{}
	if err := web.Bind(r.HTTPRequest, &form); err != nil {
		return nil, errBadForm.Errorf("failed to parse request: %w", err)
	}
	return f.client.AuthenticatePassword(ctx, r, form.Username, form.Password)
}

func (f *Form) Test(ctx context.Context, r *authn.Request) bool {
	// FIXME: How should we detect this??
	// Maybe create client test interface and not all clients has to implement this??
	return true
}
