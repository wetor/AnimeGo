package test_test

import (
	"context"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/wetor/AnimeGo/internal/constant"
	"github.com/wetor/AnimeGo/internal/pkg/request"
	"github.com/wetor/AnimeGo/test"
)

func TestMockMikanStart(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	host := test.MockMikanStart(ctx)
	request.Init(&request.Options{
		Host: map[string]*request.HostOptions{
			constant.MikanHost: {
				Redirect: host,
				Cookie: map[string]string{
					constant.MikanAuthCookie: "MikanAuthCookie",
				},
				Params: map[string]string{
					"testdata": "mikan",
				},
			},
		},
	})

	res, err := request.GetString(constant.MikanHost + "/Home/bangumi/228")
	if err != nil {
		log.Fatal(err)
	}
	data, _ := test.GetData("mikan", "228")
	assert.Equal(t, res, string(data))
	cancel()
}
