package utils_test

import (
	"testing"

	"github.com/venturarome/DaftWatch/internal/model"
	"github.com/venturarome/DaftWatch/internal/utils"
)

func TestDiffSlice(t *testing.T) {
	p1 := model.Property{ListingId: "01"}
	p2 := model.Property{ListingId: "02"}
	p3 := model.Property{ListingId: "03"}
	p4 := model.Property{ListingId: "04"}

	fromSlice := []model.Property{p3, p4}
	compareSlice := []model.Property{p1, p2, p3}

	diffSlice := utils.DiffSlice(
		fromSlice,
		compareSlice,
		func(pA model.Property, pB model.Property) bool {
			return pA.ListingId == pB.ListingId
		},
	)

	if len(diffSlice) != 1 {
		t.Errorf("diffSlice should have length %d but has length %d", 1, len(diffSlice))
	}
	if diffSlice[0].ListingId != p4.ListingId {
		t.Errorf("diffSlice[0].ListingId should have value %s but has length %s", p4.ListingId, diffSlice[0].ListingId)
	}
}

func TestDiffSliceWithEmptyFromSlice(t *testing.T) {
	p1 := model.Property{ListingId: "01"}
	p2 := model.Property{ListingId: "02"}
	p3 := model.Property{ListingId: "03"}

	fromSlice := []model.Property{}
	compareSlice := []model.Property{p1, p2, p3}

	diffSlice := utils.DiffSlice(
		fromSlice,
		compareSlice,
		func(pA model.Property, pB model.Property) bool {
			return pA.ListingId == pB.ListingId
		},
	)

	if len(diffSlice) != 0 {
		t.Errorf("diffSlice should have length %d but has length %d", 0, len(diffSlice))
	}
}

func TestDiffSliceWithEmptyCompareSlice(t *testing.T) {
	p1 := model.Property{ListingId: "01"}
	p2 := model.Property{ListingId: "02"}
	p3 := model.Property{ListingId: "03"}

	fromSlice := []model.Property{p1, p2, p3}
	compareSlice := []model.Property{}

	diffSlice := utils.DiffSlice(
		fromSlice,
		compareSlice,
		func(pA model.Property, pB model.Property) bool {
			return pA.ListingId == pB.ListingId
		},
	)

	if len(diffSlice) != 3 {
		t.Errorf("diffSlice should have length %d but has length %d", 1, len(diffSlice))
	}
	for i := range []int{0, 1, 2} {
		if diffSlice[i].ListingId != fromSlice[i].ListingId {
			t.Errorf("diffSlice[%d].ListingId should have value %s but has length %s", i, fromSlice[i].ListingId, diffSlice[i].ListingId)
		}
	}

}
