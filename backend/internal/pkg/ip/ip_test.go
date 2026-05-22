package ip

import (
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func TestGetClientCountryCode(t *testing.T) {
	gin.SetMode(gin.TestMode)

	req, err := http.NewRequest(http.MethodGet, "/", nil)
	require.NoError(t, err)
	req.Header.Set("CF-IPCountry", " cn ")

	c := &gin.Context{Request: req}
	require.Equal(t, "CN", GetClientCountryCode(c))
}

func TestGetClientCountryCodeFallbackHeaders(t *testing.T) {
	gin.SetMode(gin.TestMode)

	req, err := http.NewRequest(http.MethodGet, "/", nil)
	require.NoError(t, err)
	req.Header.Set("X-GeoIP-Country-Code", "US")
	req.Header.Set("X-Forwarded-Country", "CN")

	c := &gin.Context{Request: req}
	require.Equal(t, "US", GetClientCountryCode(c))
}

func TestIsMainlandChinaCountryCode(t *testing.T) {
	require.True(t, IsMainlandChinaCountryCode("CN"))
	require.True(t, IsMainlandChinaCountryCode("chn"))
	require.True(t, IsMainlandChinaCountryCode("mainland-china"))
	require.True(t, IsMainlandChinaCountryCode("中国大陆"))

	require.False(t, IsMainlandChinaCountryCode("HK"))
	require.False(t, IsMainlandChinaCountryCode("TW"))
	require.False(t, IsMainlandChinaCountryCode(""))
}
