package model

import "errors"

const defaultLimitUsers = 10

type ListUsersParams struct {
	Limit  uint
	Query  string
	Page   uint
	SortBy ListUsersParamsSortBy
}

type ListUsersParamsSortBy string

type ListUsersParamsOption func(pms *ListUsersParams) error

func WithLimit(limit uint) ListUsersParamsOption {
	return func(pms *ListUsersParams) error {
		if limit < 1 {
			return errors.New("limit parameter must be greater or equal 1")
		}
		// TODO: add max limit
		pms.Limit = limit
		return nil
	}
}

func WithPage(page uint) ListUsersParamsOption {
	return func(pms *ListUsersParams) error {
		if page < 1 {
			return errors.New("page parameter must be greater or equal 1")
		}
		pms.Page = page
		return nil
	}
}

func WithQuery(query string) ListUsersParamsOption {
	return func(pms *ListUsersParams) error {
		pms.Query = query
		return nil
	}
}

func SortBy(by ListUsersParamsSortBy) ListUsersParamsOption {
	return func(pms *ListUsersParams) error {
		switch by {
		case ListUsersParamsSortByAlphabet, ListUsersParamsSortByDepartment:
			pms.SortBy = by
		default:
			return errors.New("unknown sort by parameter")
		}
		return nil
	}
}

const (
	ListUsersParamsSortByAlphabet   ListUsersParamsSortBy = "alphabet"
	ListUsersParamsSortByDepartment ListUsersParamsSortBy = "department"
)

func NewListUsersParams(opts ...ListUsersParamsOption) (ListUsersParams, error) {
	pms := &ListUsersParams{
		Limit:  defaultLimitUsers,
		Query:  "",
		Page:   1,
		SortBy: ListUsersParamsSortByAlphabet,
	}

	for _, opt := range opts {
		if err := opt(pms); err != nil {
			return *pms, err
		}
	}
	return *pms, nil
}
