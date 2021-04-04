package http

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	DefaultLimit  = 10
	DefaultOffset = 0
)

type Direction string

const (
	Ascending  Direction = "asc"
	Descending Direction = "desc"
)

type SortPair struct {
	Attribute string
	Direction Direction
}

type PaginationRequest struct {
	Limit     int        `json:"limit"`
	Offset    int        `json:"offset"`
	SortPairs []SortPair `json:"sort"`
}

func PaginationFromContext(c *gin.Context) (*PaginationRequest, error) {
	limitStr := c.Query("limit")
	offsetStr := c.Query("offset")
	sort := c.Query("sort")

	var pg PaginationRequest

	if limitStr != "" && offsetStr != "" {
		limit, err := strconv.Atoi(limitStr)
		if err != nil {
			return nil, err
		}

		offset, err := strconv.Atoi(offsetStr)
		if err != nil {
			return nil, err
		}

		if limit < 0 || offset < 0 {
			return nil, errors.New("limit and offset should be greater than 0")
		}

		pg = PaginationRequest{
			Limit:  limit,
			Offset: offset,
		}
	} else {
		pg.Limit = DefaultLimit
		pg.Offset = DefaultOffset
	}

	if sort != "" {
		pairs := strings.Split(sort, ",")
		for i, pair := range pairs {
			ps := strings.Split(pair, ":")
			if len(ps) != 2 {
				return nil, fmt.Errorf("invalid sort string formatting at index [%d]", i)
			}

			direction, err := parseDirection(ps[1])
			if err != nil {
				return nil, err
			}

			pg.SortPairs = append(pg.SortPairs, SortPair{
				Attribute: strings.TrimSpace(ps[0]),
				Direction: direction,
			})
		}
	}

	return &pg, nil
}

func parseDirection(dir string) (Direction, error) {
	switch strings.TrimSpace(strings.ToLower(dir)) {
	case string(Ascending):
		return Ascending, nil
	case string(Descending):
		return Descending, nil
	}

	return "", errors.New("invalid sort parameter: invalid direction value, needs to be in [asc, desc]")

}
