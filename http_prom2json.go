package main

import (
	"net/http"
	"runtime"
	"strings"
	"errors"
	"regexp"
	"github.com/labstack/echo"

	dto "github.com/prometheus/client_model/go"
)

func prom2jsonHandler(c echo.Context) error {
	url := c.QueryParams()["url"]
	if url == nil {
		return BadRequestJSON(c, "`url` parameter for prometheus text-formatted output is required.")
	}

	predicates := c.QueryParams()["filter"]

	runtime.GOMAXPROCS(2)

	mfChan := make(chan *dto.MetricFamily, 1024)
	errChan := make(chan error)
	go fetchMetricFamilies(url[0], mfChan, errChan)

	result := []*metricFamily{}
	for mf := range mfChan {
		passed, err := filter(predicates, mf)
		if err != nil {
			return BadRequestJSON(c, err.Error())
		}

		if !passed {
			continue
		}

		result = append(result, newMetricFamily(mf))
	}

	select {
	case err := <-errChan:
		return BadRequestJSON(c, err.Error())
	default:
	}

	return c.JSON(
		http.StatusOK, &JsonResponse{
			Content: result,
		},
	)
}

// TBD
// Predicates are in format key|<regexp>
func filter(predicates []string, elem *dto.MetricFamily) (bool, error) {
	for _, predicate := range predicates {
		split := strings.Split(predicate, "|")
		if len(split) != 2 {
			return false, errors.New("Bad format of filter predicate. Should be: filter=key|<regexp>")
		}

		// Get json name tag.
		//field, ok := reflect.TypeOf(dto.MetricFamily).Elem().FieldByName("name")
		//tag = string(field.Tag)
		if split[0] != "name" {
			return false, errors.New("Other keys than 'name' not implemented yet.")
		}

		reg, err := regexp.Compile(split[1])
		if err != nil {
			return false, err
		}

		if !reg.MatchString(*elem.Name) {
			return false, nil
		}
	}

	return true, nil
}