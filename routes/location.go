package routes

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"

	"github.com/kataras/iris/v12"
)

func Autocomplete(ctx iris.Context) {
	limit := "10"
	location := ctx.URLParam("location")

	limitQuery := ctx.URLParam("limit")

	if limitQuery != "" {
		limit = limitQuery
	}

	apikey := os.Getenv("LOCATION_IQ_TOKEN")
	// 	url := "https://api.locationiq.com/v1/autocomplete?key=" + apikey + "&q=" + location + "&limit=" + limit
	baseURL := "https://api.locationiq.com/v1/autocomplete"
	params := url.Values{}
	params.Add("key", apikey)
	params.Add("q", location)
	params.Add("limit", limit)

	url := fmt.Sprintf("%s?%s", baseURL, params.Encode())

	fetchLocations(url, ctx)

}

func Search(ctx iris.Context) {
	location := ctx.URLParam("location")

	apikey := os.Getenv("LOCATION_IQ_TOKEN")

	baseURL := "https://api.locationiq.com/v1/search"

	params := url.Values{}
	params.Add("key", apikey)
	params.Add("q", location)
	params.Add("format", "json")

	url := fmt.Sprintf("%s?%s", baseURL, params.Encode())

	fetchLocations(url, ctx)
}

func fetchLocations(url string, ctx iris.Context) {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)

	res, LocationErr := client.Do(req)

	if LocationErr != nil {
		ctx.StopWithProblem(iris.StatusInternalServerError,
			iris.NewProblem().
				Title("internal Server Error").
				Detail("autocomplete Internal server Error"))
		return
	}

	defer res.Body.Close()

	body, bodyErr := io.ReadAll(res.Body)

	if bodyErr != nil {
		ctx.StopWithProblem(iris.StatusInternalServerError,
			iris.NewProblem().
				Title("internal Server Error").
				Detail("Response Internal server Error"))
		return
	}

	var objMap []map[string]interface{}
	jsonErr := json.Unmarshal(body, &objMap)

	if jsonErr != nil {
		ctx.StopWithProblem(iris.StatusInternalServerError,
			iris.NewProblem().
				Title("internal Server Error").
				Detail("Json Internal server Error"))
		return
	}

	ctx.JSON(objMap)
}
