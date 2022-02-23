// Package diffmanager is used for cloning repo
package diffmanager

import (
	"context"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/LambdaTest/synapse/pkg/core"
	"github.com/LambdaTest/synapse/pkg/errs"
	"github.com/LambdaTest/synapse/pkg/global"
	"github.com/LambdaTest/synapse/testutils"
)

func Test_updateWithOr(t *testing.T) {
	check := func(t *testing.T) {
		dm := &diffManager{}
		m := make(map[string]int)
		key := "key"
		val := rand.Intn(1000)
		dm.updateWithOr(m, key, val)
		if ans, exists := m[key]; !exists || ans != val {
			t.Errorf("Expected: %v, received: %v", val, m[key])
		}
		newVal := rand.Intn(1000)
		dm.updateWithOr(m, key, newVal)
		if ans, exists := m[key]; !exists || ans != (val|newVal) {
			t.Errorf("Expected: %v, received: %v", val|newVal, m[key])
		}
	}
	t.Run("Test_updateWithOr", func(t *testing.T) {
		check(t)
	})
}

func Test_diffManager_GetChangedFiles_PRDiff(t *testing.T) {
	server := httptest.NewServer( // mock server
		http.FileServer(http.Dir("../../testutils")), // mock data stored at testutils/testdata
	)
	defer server.Close()

	logger, err := testutils.GetLogger()
	if err != nil {
		t.Errorf("Can't get logger, received: %s", err)
	}
	config, err := testutils.GetConfig()
	if err != nil {
		t.Errorf("Can't get logger, received: %s", err)
	}

	dm := NewDiffManager(config, logger)
	type args struct {
		ctx        context.Context
		payload    *core.Payload
		cloneToken string
	}
	tests := []struct {
		name string
		args args
		want map[string]int
	}{
		{"Test GetChangedFile for PRdiff for github gitprovider", args{ctx: context.TODO(), payload: &core.Payload{RepoSlug: "/testdata", RepoLink: server.URL + "/testdata", BuildTargetCommit: "abcd", BuildBaseCommit: "efgh", TargetCommit: "ijkl", BaseCommit: "mnop", TaskID: "taskid", BranchName: "branchname", BuildID: "buildid", RepoID: "repoid", OrgID: "orgid", GitProvider: "github", PrivateRepo: false, EventType: "pull-request", Diff: "xyz", PullRequestNumber: 2}, cloneToken: ""}, map[string]int{}},

		{"Test GetChangedFile for PRdiff for gitlab gitprovider", args{ctx: context.TODO(), payload: &core.Payload{RepoSlug: "/testdata", RepoLink: server.URL + "/testdata", BuildTargetCommit: "abcd", BuildBaseCommit: "efgh", TargetCommit: "ijkl", BaseCommit: "mnop", TaskID: "taskid", BranchName: "branchname", BuildID: "buildid", RepoID: "repoid", OrgID: "orgid", GitProvider: "gitlab", PrivateRepo: false, EventType: "pull-request", Diff: "xyz", PullRequestNumber: 2}, cloneToken: ""}, map[string]int{}},

		{"Test GetChangedFile for PRdiff for unsupported gitprovider", args{ctx: context.TODO(), payload: &core.Payload{RepoSlug: "/testdata", RepoLink: server.URL + "/testdata", BuildTargetCommit: "abcd", BuildBaseCommit: "efgh", TargetCommit: "ijkl", BaseCommit: "mnop", TaskID: "taskid", BranchName: "branchname", BuildID: "buildid", RepoID: "repoid", OrgID: "orgid", GitProvider: "", PrivateRepo: false, EventType: "pull-request", Diff: "xyz", PullRequestNumber: 2}, cloneToken: ""}, map[string]int{}},

		{"Test GetChangedFile for PRdiff for non 200 status", args{ctx: context.TODO(), payload: &core.Payload{RepoSlug: "/notfound/", RepoLink: server.URL + "/testdata", BuildTargetCommit: "abcd", BuildBaseCommit: "efgh", TargetCommit: "ijkl", BaseCommit: "mnop", TaskID: "taskid", BranchName: "branchname", BuildID: "buildid", RepoID: "repoid", OrgID: "orgid", GitProvider: "", PrivateRepo: false, EventType: "pull-request", Diff: "xyz", PullRequestNumber: 2}, cloneToken: ""}, map[string]int{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			global.APIHostURLMap[tt.args.payload.GitProvider] = server.URL
			resp, err := dm.GetChangedFiles(tt.args.ctx, tt.args.payload, tt.args.cloneToken)
			if tt.args.payload.GitProvider == "" {
				if err == nil {
					t.Errorf("Expected error: 'unsupoorted git provider', received: %v\n", err)
				}
				return
			}
			if tt.args.payload.RepoSlug == "/notfound/" {
				expErr := errs.New("non 200 response")
				if err == nil {
					t.Errorf("Expected error: %s, received error: %s", expErr, err)
				}
				return
			}
			expResp := testutils.GetGitDiff()
			if err != nil {
				t.Errorf("error in getting changed files, error %v", err.Error())
			} else if tt.args.payload.GitProvider == "github" && !reflect.DeepEqual(resp, expResp) {
				t.Errorf("Expected: %+v, received: %+v", expResp, resp)
			} else if tt.args.payload.GitProvider == "gitlab" && len(resp) != 17 {
				t.Errorf("Expected map entries: 17, received: %v, received map: %v", len(resp), resp)
			}
			return
		})
	}
}

func Test_diffManager_GetChangedFiles_CommitDiff_Github(t *testing.T) {
	server := httptest.NewServer( // mock server
		http.FileServer(http.Dir("../../testutils")), // mock data stored at testutils/testdata
	)
	defer server.Close()

	logger, err := testutils.GetLogger()
	if err != nil {
		t.Errorf("Can't get logger, received: %s", err)
	}
	config, err := testutils.GetConfig()
	if err != nil {
		t.Errorf("Can't get logger, received: %s", err)
	}

	dm := NewDiffManager(config, logger)
	type args struct {
		ctx        context.Context
		payload    *core.Payload
		cloneToken string
	}
	tests := []struct {
		name string
		args args
		want map[string]int
	}{
		{"Test GetChangedFile for CommitDiff for github gitprovider", args{ctx: context.TODO(), payload: &core.Payload{RepoSlug: "/testdata", RepoLink: server.URL + "/testdata", BuildTargetCommit: "abcd", BuildBaseCommit: "efgh", TargetCommit: "xyz", BaseCommit: "abc", TaskID: "taskid", BranchName: "branchname", BuildID: "buildid", RepoID: "repoid", OrgID: "orgid", GitProvider: "github", PrivateRepo: false, EventType: "push", Diff: "xyz", PullRequestNumber: 2}, cloneToken: ""}, map[string]int{}},

		{"Test GetChangedFile for CommitDiff for github provider and empty base commit", args{ctx: context.TODO(), payload: &core.Payload{RepoSlug: "/testdata", RepoLink: server.URL + "/testdata", BuildTargetCommit: "abcd", BuildBaseCommit: "efgh", TargetCommit: "", BaseCommit: "", TaskID: "taskid", BranchName: "branchname", BuildID: "buildid", RepoID: "repoid", OrgID: "orgid", GitProvider: "gitlab", PrivateRepo: false, EventType: "push", Diff: "xyz", PullRequestNumber: 2}, cloneToken: ""}, map[string]int{}},

		{"Test GetChangedFile for CommitDiff for github provider for non 200 response", args{ctx: context.TODO(), payload: &core.Payload{RepoSlug: "/notfound/", RepoLink: server.URL + "/notfound/", BuildTargetCommit: "abcd", BuildBaseCommit: "efgh", TargetCommit: "xyz", BaseCommit: "abc", TaskID: "taskid", BranchName: "branchname", BuildID: "buildid", RepoID: "repoid", OrgID: "orgid", GitProvider: "gitlab", PrivateRepo: false, EventType: "push", Diff: "xyz", PullRequestNumber: 2}, cloneToken: ""}, map[string]int{}},

		{"Test GetChangedFile for CommitDiff for non github and gitlab gitprovider", args{ctx: context.TODO(), payload: &core.Payload{RepoSlug: "/notfound/", RepoLink: server.URL + "/notfound/", BuildTargetCommit: "abcd", BuildBaseCommit: "efgh", TargetCommit: "xyz", BaseCommit: "abc", TaskID: "taskid", BranchName: "branchname", BuildID: "buildid", RepoID: "repoid", OrgID: "orgid", GitProvider: "gittest", PrivateRepo: false, EventType: "push", Diff: "xyz", PullRequestNumber: 2}, cloneToken: ""}, map[string]int{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			global.APIHostURLMap[tt.args.payload.GitProvider] = server.URL
			resp, err := dm.GetChangedFiles(tt.args.ctx, tt.args.payload, tt.args.cloneToken)
			if tt.args.payload.GitProvider == "gittest" {
				if resp != nil || err == nil {
					t.Errorf("Expected error: 'unsupoorted git provider', received: %v\nexpected response: nil, received: %v", err, resp)
				}
				return
			}
			if tt.args.payload.BaseCommit == "" || tt.args.payload.RepoSlug == "/notfound/" {
				if err != nil {
					t.Errorf("Received error: %v, response: %v", err, resp)
				}
				return
			}
			expResp := make(map[string]int)
			if err != nil {
				t.Errorf("error in getting changed files, error %v", err.Error())
			} else if !reflect.DeepEqual(resp, expResp) {
				t.Errorf("Expected: %+v, received: %+v", expResp, resp)
			}
		})
	}
}

func Test_diffManager_GetChangedFiles_CommitDiff_Gitlab(t *testing.T) {
	data, err := testutils.GetGitlabCommitDiff()
	if err != nil {
		t.Errorf("Received error in getting test gitlab commit diff, error: %v", err)
	}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/testdata/repository/compare" {
			t.Errorf("Expected to request, got: %v", r.URL.Path)
		}
		w.WriteHeader(200)
		w.Header().Set("Content-Type", "application/json")
		w.Write(data)
	}))
	defer server.Close()

	logger, err := testutils.GetLogger()
	if err != nil {
		t.Errorf("Can't get logger, received: %s", err)
	}
	config, err := testutils.GetConfig()
	if err != nil {
		t.Errorf("Can't get logger, received: %s", err)
	}

	dm := NewDiffManager(config, logger)
	type args struct {
		ctx        context.Context
		payload    *core.Payload
		cloneToken string
	}
	tests := []struct {
		name string
		args args
		want map[string]int
	}{
		{"Test GetChangedFile for CommitDiff for gitlab gitprovider", args{ctx: context.TODO(), payload: &core.Payload{RepoSlug: "/testdata", RepoLink: server.URL + "/testdata", BuildTargetCommit: "abcd", BuildBaseCommit: "efgh", TargetCommit: "xyz", BaseCommit: "abc", TaskID: "taskid", BranchName: "branchname", BuildID: "buildid", RepoID: "repoid", OrgID: "orgid", GitProvider: "gitlab", PrivateRepo: false, EventType: "push", Diff: "xyz", PullRequestNumber: 2}, cloneToken: ""}, map[string]int{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			global.APIHostURLMap[tt.args.payload.GitProvider] = server.URL
			resp, err := dm.GetChangedFiles(tt.args.ctx, tt.args.payload, tt.args.cloneToken)
			if tt.args.payload.BaseCommit == "" || tt.args.payload.RepoSlug == "/notfound/" {
				if err != nil || resp != nil {
					t.Errorf("Received error: %v, response: %v", err, resp)
				}
				return
			}
			if err != nil {
				t.Errorf("error in getting changed files, error %v", err.Error())
			} else if len(resp) != 202 {
				t.Errorf("Expected map length: 202, received: %v\nreceived map: %v", len(resp), resp)
			}
		})
	}
}
