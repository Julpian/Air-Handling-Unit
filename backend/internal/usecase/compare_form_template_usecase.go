package usecase

import (
	"ahu-backend/internal/domain"
	"ahu-backend/internal/repository"
)

type CompareFormTemplateUsecase struct {
	formRepo repository.FormRepository
}

func NewCompareFormTemplateUsecase(
	formRepo repository.FormRepository,
) *CompareFormTemplateUsecase {
	return &CompareFormTemplateUsecase{formRepo}
}

func (uc *CompareFormTemplateUsecase) Execute(
	fromID string,
	toID string,
) (*domain.FormTemplateDiff, error) {

	from, err := uc.formRepo.GetTemplateByID(fromID)
	if err != nil {
		return nil, err
	}

	to, err := uc.formRepo.GetTemplateByID(toID)
	if err != nil {
		return nil, err
	}

	diff := &domain.FormTemplateDiff{
		From: domain.FormTemplateMeta{
			ID:      from.ID,
			Version: from.Version,
		},
		To: domain.FormTemplateMeta{
			ID:      to.ID,
			Version: to.Version,
		},
	}

	// map sections by code
	fromSections := map[string]domain.FormSection{}
	toSections := map[string]domain.FormSection{}

	for _, s := range from.Sections {
		fromSections[s.Code] = s
	}
	for _, s := range to.Sections {
		toSections[s.Code] = s
	}

	visited := map[string]bool{}

	// compare existing
	for _, fromSec := range from.Sections {
		code := fromSec.Code
		visited[code] = true
		toSec, ok := toSections[code]

		if !ok {
			diff.Sections = append(diff.Sections, domain.SectionDiff{
				Code:   fromSec.Code,
				Title:  fromSec.Title,
				Status: "removed",
			})
			continue
		}

		sectionDiff := domain.SectionDiff{
			Code:  fromSec.Code,
			Title: fromSec.Title,
		}

		itemDiffs := compareItems(fromSec.Items, toSec.Items)

		sectionDiff.Items = itemDiffs
		sectionDiff.Status = sectionStatus(itemDiffs)

		diff.Sections = append(diff.Sections, sectionDiff)
	}

	// added sections
	for code, toSec := range toSections {
		if visited[code] {
			continue
		}
		diff.Sections = append(diff.Sections, domain.SectionDiff{
			Code:   toSec.Code,
			Title:  toSec.Title,
			Status: "added",
		})
	}

	return diff, nil
}

func compareItems(
	from []domain.FormItem,
	to []domain.FormItem,
) []domain.ItemDiff {

	toMap := map[string]domain.FormItem{}

	for _, i := range to {
		toMap[i.Label] = i
	}

	visited := map[string]bool{}
	var diffs []domain.ItemDiff

	// 🔥 LOOP FROM SLICE (bukan map)
	for _, f := range from {
		t, ok := toMap[f.Label]
		visited[f.Label] = true

		if !ok {
			diffs = append(diffs, domain.ItemDiff{
				Label:  f.Label,
				Status: "removed",
				From:   &f,
			})
			continue
		}

		if f.InputType != t.InputType || f.Required != t.Required {
			diffs = append(diffs, domain.ItemDiff{
				Label:  f.Label,
				Status: "modified",
				From:   &f,
				To:     &t,
			})
		} else {
			diffs = append(diffs, domain.ItemDiff{
				Label:  f.Label,
				Status: "unchanged",
			})
		}
	}

	// 🔥 added items → tetap urut sesuai "to"
	for _, t := range to {
		if visited[t.Label] {
			continue
		}

		diffs = append(diffs, domain.ItemDiff{
			Label:  t.Label,
			Status: "added",
			To:     &t,
		})
	}

	return diffs
}

func sectionStatus(items []domain.ItemDiff) string {
	for _, i := range items {
		if i.Status != "unchanged" {
			return "modified"
		}
	}
	return "unchanged"
}
