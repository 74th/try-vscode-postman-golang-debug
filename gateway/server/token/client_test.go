package token_test

import (
	"context"
	"testing"

	"github.com/74th/vscode-book-r2-golang/gateway/server/token"
)

func TestClient(t *testing.T) {
	client := token.New("https://a274ebfe-ee67-4c03-a2db-f278c2535a83.mock.pstmn.io")
	ctx := context.Background()

	ok, err := client.Validate(ctx, "token")
	if err != nil {
		t.Error("エラーが返らないこと")
		return
	}

	if !ok {
		t.Error("OKが返ること")
	}
}
