package leadtime_test

import (
	"time"

	"github.com/concourse/concourse/atc"
	"github.com/concourse/concourse/go-concourse/concourse"
	leadtimefakes "github.com/concourse/leadtime/lead-timefakes"
	"github.com/concourse/leadtime/leadtime"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 github.com/concourse/concourse/go-concourse/concourse.Team

var _ = Describe("Repository", func() {
	var fakeClient *leadtimefakes.FakeClient
	var fakeTeam *leadtimefakes.FakeTeam
	var node leadtime.Node
	var now time.Time

	BeforeEach(func() {
		now = time.Now()
		fakeClient = new(leadtimefakes.FakeClient)
		fakeTeam = new(leadtimefakes.FakeTeam)
		fakeClient.TeamReturns(fakeTeam)
		fakeTeam.JobBuildsReturns(
			[]atc.Build{
				atc.Build{
					ID:           1,
					TeamName:     "main",
					Name:         "0",
					Status:       "failed",
					JobName:      "publish-binaries",
					APIURL:       "/api/v1/builds/1",
					PipelineName: "concourse",
					StartTime:    0,
					EndTime:      1000,
					ReapTime:     0,
				},
				atc.Build{
					ID:           0,
					TeamName:     "main",
					Name:         "0",
					Status:       "succeeded",
					JobName:      "publish-binaries",
					APIURL:       "/api/v1/builds/0",
					PipelineName: "concourse",
					StartTime:    0,
					EndTime:      now.Unix(),
					ReapTime:     0,
				},
			},
			concourse.Pagination{},
			true,
			nil,
		)

		node = leadtime.NewRepository(fakeClient).FindEnd()
	})

	It("gets 'main' team", func() {
		Expect(fakeClient.TeamCallCount()).To(Equal(1))
		Expect(fakeClient.TeamArgsForCall(0)).To(Equal("main"))
	})

	It("gets latest builds of 'publish-binaries' job in 'concourse' pipeline", func() {
		Expect(fakeTeam.JobBuildsCallCount()).To(Equal(1))
		pipelineName, jobName, _ := fakeTeam.JobBuildsArgsForCall(0)
		Expect(pipelineName).To(Equal("concourse"), "wrong pipeline name")
		Expect(jobName).To(Equal("publish-binaries"), "wrong job name")
	})

	It("returns successful build from API", func() {
		Expect(node.Done()).To(Equal(now.Truncate(1 * time.Second)))
	})
	// TODO what if the JobBuilds call errors?
	// TODO what if zero builds are returned?
})
