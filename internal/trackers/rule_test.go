package trackers_test

import (
	"net/url"
	"testing"

	"github.com/carnivorecookies/Clipit/internal/trackers"
	"github.com/stretchr/testify/assert"
)

var amazonRule = &trackers.Rule{
	UrlPattern: `^https?:\/\/(?:[a-z0-9-]+\.)*?amazon(?:\.[a-z]{2,}){1,}`,
	Blocked:    false,
	QueryRules: []string{
		"p[fd]_rd_[a-z]*",
		"qid",
		"srs?",
		"__mk_[a-z]{1,3}_[a-z]{1,3}",
		"spIA",
		"ms3_c",
		"[a-z%0-9]*ie",
		"refRID",
		"colii?d",
		"[^a-z%0-9]adId",
		"qualifier",
		"_encoding",
		"smid",
		"field-lbr_brands_browse-bin",
		"ref_?",
		"th",
		"sprefix",
		"crid",
		"keywords",
		"cv_ct_[a-z]+",
		"linkCode",
		"creativeASIN",
		"ascsubtag",
		"aaxitk",
		"hsa_cr_id",
		"sb-ci-[a-z]+",
		"rnid",
		"dchild",
		"camp",
		"creative",
		"content-id",
		"dib",
		"dib_tag",
		"social_share",
		"starsLeft",
		"skipTwisterOG",
	},
	RawURLRules: []string{
		`\/ref=[^/?]*`,
	},
	Exceptions: []string{
		`^https?:\/\/(?:[a-z0-9-]+\.)*?amazon(?:\.[a-z]{2,}){1,}\/gp\/.*?(?:redirector.html|cart\/ajax-update.html|video\/api\/)`,
		`^https?:\/\/(?:[a-z0-9-]+\.)*?amazon(?:\.[a-z]{2,}){1,}\/(?:hz\/reviews-render\/ajax\/|message-us\?|s\?)`,
	},
	Redirections: nil,
}

func TestRule(t *testing.T) {
	t.Run("amazon", func(t *testing.T) {
		assert := assert.New(t)
		link, _ := url.Parse("https://www.amazon.com/dp/exampleProduct/ref=sxin_0_pb?__mk_de_DE=ÅMÅŽÕÑ&keywords=tea&pd_rd_i=exampleProduct&pd_rd_r=8d39e4cd-1e4f-43db-b6e7-72e969a84aa5&pd_rd_w=1pcKM&pd_rd_wg=hYrNl&pf_rd_p=50bbfd25-5ef7-41a2-68d6-74d854b30e30&pf_rd_r=0GMWD0YYKA7XFGX55ADP&qid=1517757263&rnid=2914120011")

		removed := amazonRule.Compile().RemoveTrackersFrom(link)
		assert.NotNil(removed)
		assert.Empty(removed.Query())
		assert.Equal(removed.String(), "https://www.amazon.com/dp/exampleProduct")
	})

	t.Run("blocked", func(t *testing.T) {
		assert := assert.New(t)

		r := &trackers.Rule{UrlPattern: ".*", Blocked: true}
		link, _ := url.Parse("https://foo.bar")
		removed := r.Compile().RemoveTrackersFrom(link)
		assert.Nil(removed)
	})

	// TODO: test redirection
}
