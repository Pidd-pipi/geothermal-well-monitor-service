package main

import (
	"context"
	"testing"
)

func TestOpsStoreListLabelsIsolated(t *testing.T) {
	svc := newOpsService(nil)
	ctx := context.Background()
	rec, err := svc.Create(ctx, OpsRecord{Subject: "井口检查", Owner: "liang", Priority: OpsPriorityHigh, Labels: map[string]string{"site": "north"}})
	if err != nil {
		t.Fatal(err)
	}
	page, err := svc.Search(ctx, OpsQuery{Page: 1, PageSize: 25})
	if err != nil {
		t.Fatal(err)
	}
	if len(page.Items) != 1 {
		t.Fatalf("items=%d", len(page.Items))
	}
	page.Items[0].Labels["note"] = "polluted"
	got, err := svc.Get(ctx, rec.ID)
	if err != nil {
		t.Fatal(err)
	}
	if got.Labels["note"] == "polluted" {
		t.Fatalf("list mutation leaked into store labels: %v", got.Labels)
	}
}

func TestOpsStoreGetLabelsIsolated(t *testing.T) {
	svc := newOpsService(nil)
	ctx := context.Background()
	rec, err := svc.Create(ctx, OpsRecord{Subject: "阀门巡检", Owner: "zhao", Priority: OpsPriorityNormal, Labels: map[string]string{"site": "east"}})
	if err != nil {
		t.Fatal(err)
	}
	first, err := svc.Get(ctx, rec.ID)
	if err != nil {
		t.Fatal(err)
	}
	first.Labels["note"] = "polluted"
	second, err := svc.Get(ctx, rec.ID)
	if err != nil {
		t.Fatal(err)
	}
	if second.Labels["note"] == "polluted" {
		t.Fatalf("get mutation leaked into store labels: %v", second.Labels)
	}
}

func TestOpsStorePutIsolatesCallerLabels(t *testing.T) {
	svc := newOpsService(nil)
	ctx := context.Background()
	labels := map[string]string{"site": "west"}
	rec := OpsRecord{Subject: "渗漏复查", Owner: "liang", Priority: OpsPriorityHigh, Labels: labels}
	created, err := svc.Create(ctx, rec)
	if err != nil {
		t.Fatal(err)
	}
	labels["note"] = "polluted"
	got, err := svc.Get(ctx, created.ID)
	if err != nil {
		t.Fatal(err)
	}
	if got.Labels["note"] == "polluted" {
		t.Fatalf("caller label mutation leaked into store: %v", got.Labels)
	}
}

func TestOpsRecordCloneDeepCopiesLabels(t *testing.T) {
	original := OpsRecord{Subject: "s", Owner: "o", Labels: map[string]string{"site": "north"}}
	clone := original.Clone()
	clone.Labels["site"] = "south"
	if original.Labels["site"] != "north" {
		t.Fatalf("clone shared labels map: %v", original.Labels)
	}
}

func TestOpsSearchHugePageEmpty(t *testing.T) {
	svc := newOpsService(nil)
	ctx := context.Background()
	for _, subject := range []string{"a", "b", "c"} {
		if _, err := svc.Create(ctx, OpsRecord{Subject: subject, Owner: "liang", Priority: OpsPriorityNormal, Labels: map[string]string{"site": "north"}}); err != nil {
			t.Fatal(err)
		}
	}
	page, err := svc.Search(ctx, OpsQuery{Page: 99999, PageSize: 5})
	if err != nil {
		t.Fatal(err)
	}
	if len(page.Items) != 0 {
		t.Fatalf("huge page returned %d items (want empty)", len(page.Items))
	}
}

func TestOpsSearchOwnerExactMatch(t *testing.T) {
	svc := newOpsService(nil)
	ctx := context.Background()
	if _, err := svc.Create(ctx, OpsRecord{Subject: "zhang 的记录", Owner: "zhang", Priority: OpsPriorityNormal, Labels: map[string]string{"site": "north"}}); err != nil {
		t.Fatal(err)
	}
	if _, err := svc.Create(ctx, OpsRecord{Subject: "li 的记录", Owner: "li", Priority: OpsPriorityNormal, Labels: map[string]string{"site": "east"}}); err != nil {
		t.Fatal(err)
	}
	page, err := svc.Search(ctx, OpsQuery{Owner: "an", Page: 1, PageSize: 25})
	if err != nil {
		t.Fatal(err)
	}
	if len(page.Items) != 0 {
		t.Fatalf("owner substring match leaked %d records: %+v", len(page.Items), page.Items)
	}
}

func TestOpsStoreUpdateLabelsIsolated(t *testing.T) {
	svc := newOpsService(nil)
	ctx := context.Background()
	rec, err := svc.Create(ctx, OpsRecord{Subject: "压力校验", Owner: "liang", Priority: OpsPriorityHigh, Labels: map[string]string{"site": "north"}})
	if err != nil {
		t.Fatal(err)
	}
	updated, err := svc.Transition(ctx, rec.ID, rec.Revision, OpsStatusActive, "liang")
	if err != nil {
		t.Fatal(err)
	}
	updated.Labels["note"] = "polluted"
	got, err := svc.Get(ctx, rec.ID)
	if err != nil {
		t.Fatal(err)
	}
	if got.Labels["note"] == "polluted" {
		t.Fatalf("update result mutation leaked into store labels: %v", got.Labels)
	}
}
