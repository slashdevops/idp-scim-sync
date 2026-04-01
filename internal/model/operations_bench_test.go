package model

import (
	"fmt"
	"testing"
)

// buildBenchmarkData creates IDP and SCIM group members for benchmarking.
// nGroups groups, each with nMembers members.
// Half the members exist in both (equal), half only in IDP (create).
func buildBenchmarkData(nGroups, nMembers int) (idp, scim []*GroupMembers) {
	idp = make([]*GroupMembers, nGroups)
	scim = make([]*GroupMembers, nGroups)

	for g := range nGroups {
		idpMembers := make([]*Member, nMembers)
		scimMembers := make([]*Member, nMembers/2)

		for m := range nMembers {
			email := fmt.Sprintf("user%d@group%d.com", m, g)
			idpMembers[m] = MemberBuilder().
				WithIPID(fmt.Sprintf("ipid-%d-%d", g, m)).
				WithEmail(email).
				WithStatus("ACTIVE").
				Build()
		}

		// SCIM has only the first half of members (the "equal" set)
		for m := range nMembers / 2 {
			email := fmt.Sprintf("user%d@group%d.com", m, g)
			scimMembers[m] = MemberBuilder().
				WithIPID(fmt.Sprintf("ipid-%d-%d", g, m)).
				WithSCIMID(fmt.Sprintf("scimid-%d-%d", g, m)).
				WithEmail(email).
				WithStatus("ACTIVE").
				Build()
		}

		group := GroupBuilder().
			WithIPID(fmt.Sprintf("gipid-%d", g)).
			WithSCIMID(fmt.Sprintf("gscimid-%d", g)).
			WithName(fmt.Sprintf("group-%d", g)).
			WithEmail(fmt.Sprintf("group%d@test.com", g)).
			Build()

		idp[g] = GroupMembersBuilder().WithGroup(group).WithResources(idpMembers).Build()
		scim[g] = GroupMembersBuilder().WithGroup(group).WithResources(scimMembers).Build()
	}
	return
}

func BenchmarkMembersDataSets(b *testing.B) {
	benchmarks := []struct {
		groups  int
		members int
	}{
		{10, 20},
		{50, 50},
		{50, 120},
		{100, 100},
	}

	for _, bm := range benchmarks {
		idp, scim := buildBenchmarkData(bm.groups, bm.members)
		b.Run(fmt.Sprintf("groups=%d_members=%d", bm.groups, bm.members), func(b *testing.B) {
			for range b.N {
				membersDataSets(idp, scim)
			}
		})
	}
}

func BenchmarkMergeGroupsMembersResult(b *testing.B) {
	benchmarks := []struct {
		groups  int
		members int
	}{
		{10, 20},
		{50, 50},
		{50, 120},
	}

	for _, bm := range benchmarks {
		idp, scim := buildBenchmarkData(bm.groups, bm.members)
		create, equal, _ := membersDataSets(idp, scim)

		created := &GroupsMembersResult{Items: len(create), Resources: create}
		created.SetHashCode()
		equaled := &GroupsMembersResult{Items: len(equal), Resources: equal}
		equaled.SetHashCode()

		b.Run(fmt.Sprintf("groups=%d_members=%d", bm.groups, bm.members), func(b *testing.B) {
			for range b.N {
				MergeGroupsMembersResult(created, equaled)
			}
		})
	}
}
