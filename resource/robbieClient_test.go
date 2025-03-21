package resource

import (
	"encoding/json"
	"github.com/coreos/etcd/pkg/testutil"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_signOffPost(t *testing.T) {
	var digest = "digest123"
	var signOffReq = robbieSignOffRequest{
		KcdName:      "my_kcd",
		KcdNameSpace: "my_namespace",
		KcdLables: map[string]string{
			"a": "a1",
			"b": "b1",
		},
		KcdTag:       "my_tag",
		KcdImageRepo: "my_repo",
		Versions:     []string{"123", "456", "789"},
		Digest:       &digest,
	}

	var server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b, err := io.ReadAll(r.Body)
		if err != nil {
			panic(err)
		}

		signOffReq := robbieSignOffRequest{}
		_ = json.Unmarshal(b, &signOffReq)

		testutil.AssertEqual(t, "my_kcd", signOffReq.KcdName)
		testutil.AssertEqual(t, "my_namespace", signOffReq.KcdNameSpace)
		testutil.AssertEqual(t, "my_tag", signOffReq.KcdTag)
		testutil.AssertEqual(t, "my_repo", signOffReq.KcdImageRepo)
		testutil.AssertEqual(t, "digest123", *(signOffReq.Digest))
		testutil.AssertEqual(t, 2, len(signOffReq.KcdLables))
		testutil.AssertEqual(t, 3, len(signOffReq.Versions))

		var review = signOffReview{
			Result: true,
			Uuid:   signOffReq.KcdName + "-" + signOffReq.KcdImageRepo,
		}
		var responseBody, _ = json.Marshal(&review)
		io.WriteString(w, string(responseBody))
		// mock here
	}))

	var review, err = signOffPost(&signOffReq, server.URL)
	testutil.AssertNil(t, err)
	testutil.AssertEqual(t, "my_kcd-my_repo", review.Uuid)
	testutil.AssertTrue(t, review.Result)
}

func Test_signOffPostRetryable_Good(t *testing.T) {
	var digest = "digest123"
	var signOffReq = robbieSignOffRequest{
		KcdName:      "my_kcd",
		KcdNameSpace: "my_namespace",
		KcdLables: map[string]string{
			"a": "a1",
			"b": "b1",
		},
		KcdTag:       "my_tag",
		KcdImageRepo: "my_repo",
		Versions:     []string{"123", "456", "789"},
		Digest:       &digest,
	}

	var server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b, err := io.ReadAll(r.Body)
		if err != nil {
			panic(err)
		}

		signOffReq := robbieSignOffRequest{}
		_ = json.Unmarshal(b, &signOffReq)

		testutil.AssertEqual(t, "my_kcd", signOffReq.KcdName)
		testutil.AssertEqual(t, "my_namespace", signOffReq.KcdNameSpace)
		testutil.AssertEqual(t, "my_tag", signOffReq.KcdTag)
		testutil.AssertEqual(t, "my_repo", signOffReq.KcdImageRepo)
		testutil.AssertEqual(t, "digest123", *(signOffReq.Digest))
		testutil.AssertEqual(t, 2, len(signOffReq.KcdLables))
		testutil.AssertEqual(t, 3, len(signOffReq.Versions))

		var review = signOffReview{
			Result: true,
			Uuid:   signOffReq.KcdName + "-" + signOffReq.KcdImageRepo,
		}
		var responseBody, _ = json.Marshal(&review)
		io.WriteString(w, string(responseBody))
		// mock here
	}))

	var res, err = signOffRetryalbe(&signOffReq, server.URL)
	testutil.AssertNil(t, err)
	testutil.AssertTrue(t, res)
}
