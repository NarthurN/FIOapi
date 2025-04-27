package user

import (
	"net/url"
	"strconv"
)

func parseFilters(params url.Values) *UserFilter {
	filter := UserFilter{}
	// чтобы не были равны нулю по умолчанию
	filter.AgeFrom = -1
	filter.AgeTo = -1

	if name := params.Get("name"); name != "" {
		filter.Name = name
	}

	if surname := params.Get("surname"); surname != "" {
		filter.Surname = surname
	}

	if ageFrom, err := strconv.Atoi(params.Get("age_from")); err == nil {
		filter.AgeFrom = ageFrom
	}

	if ageTo, err := strconv.Atoi(params.Get("age_to")); err == nil {
		filter.AgeTo = ageTo
	}

	if sex := params.Get("sex"); sex != "" {
		filter.Sex = sex
	}

	if nationality := params.Get("nationality"); nationality != "" {
		filter.Nationality = nationality
	}

	return &filter
}

func parsePagination(params url.Values) *Pagination {
	page := 1
	if p, err := strconv.Atoi(params.Get("page")); err == nil && p > 0 {
		page = p
	}

	perPage := 10
	if pp, err := strconv.Atoi(params.Get("per_page")); err == nil && pp > 0 {
		perPage = pp
	}

	return &Pagination{
		Page:    page,
		PerPage: perPage,
	}
}
