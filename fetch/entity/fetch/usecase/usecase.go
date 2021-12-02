package usecase

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/pajarwp/auth-fetch-app/entity/fetch"
	repo "github.com/pajarwp/auth-fetch-app/entity/fetch/repository"
	"github.com/patrickmn/go-cache"
)

type FetchUsecase interface {
	GetClaims(token string) (fetch.UserClaims, error)
	FetchResource(token string) ([]fetch.Resource, error)
	AggregateResource(token string) (map[string][]fetch.Aggregate, error)
}

type fetchUsecase struct {
	repo.FetchRepository
	*cache.Cache
}

func NewFetchUsecase(c *cache.Cache, repo repo.FetchRepository) FetchUsecase {
	return fetchUsecase{
		repo,
		c,
	}
}

func (f fetchUsecase) GetClaims(token string) (fetch.UserClaims, error) {
	return f.FetchRepository.GetClaims(token)
}

func (f fetchUsecase) FetchResource(token string) ([]fetch.Resource, error) {
	var listresource []fetch.Resource
	_, err := f.FetchRepository.GetClaims(token)
	if err != nil {
		return listresource, err
	}

	resp, err := f.FetchRepository.FetchResource(token)
	if err != nil {
		return listresource, err
	}

	converter, err := f.getCurrencyConverter()
	if err != nil {
		return listresource, err
	}

	for _, val := range resp {
		var resource fetch.Resource
		resource.UUID = val["uuid"]
		resource.Komoditas = val["komoditas"]
		resource.AreaProvinsi = val["area_provinsi"]
		resource.AreaKota = val["area_kota"]
		resource.Size = val["size"]
		resource.Price = val["price"]
		resource.TanggalParsed = val["tgl_parsed"]
		resource.Timestamp = val["timestamp"]
		if resource.Price != "" {
			price, err := strconv.ParseFloat(resource.Price, 64)
			if err != nil {
				return listresource, err
			}
			priceUSD := price * converter
			resource.PriceUSD = fmt.Sprintf("%0f", priceUSD)
		}
		listresource = append(listresource, resource)
	}
	return listresource, err
}

func (f fetchUsecase) AggregateResource(token string) (map[string][]fetch.Aggregate, error) {
	grouping := make(map[string][]fetch.Aggregate)
	claims, err := f.FetchRepository.GetClaims(token)
	if err != nil {
		return grouping, err
	}
	if strings.ToLower(claims.Role) != "admin" {
		return grouping, fmt.Errorf("Forbidden Access for role " + claims.Role)
	}

	response, err := f.FetchRepository.FetchResource(token)
	if err != nil {
		return grouping, err
	}

	var resp []map[string]string

	// remove data with empty provinsi and timestamp
	for _, v := range response {
		if v["area_provinsi"] != "" && v["timestamp"] != "" {
			resp = append(resp, v)
		}
	}

	// sort data by timestamp
	sort.Slice(resp, func(i, j int) bool {
		time1 := f.parseTimestamp(resp[i]["timestamp"])
		time2 := f.parseTimestamp(resp[j]["timestamp"])
		return time1.Before(time2)
	})

	for i := 0; i < len(resp); i++ {
		provinsi := strings.ToUpper(resp[i]["area_provinsi"])
		// check if provinsi already grouped, if true skip
		if _, ok := grouping[provinsi]; ok {
			continue
		} else {
			tempAggregate := f.assignTempAggregate(resp[i])
			grouping[provinsi] = []fetch.Aggregate{}

			for j := i + 1; j < len(resp); j++ {
				itemProvinsi := strings.ToUpper(resp[j]["area_provinsi"])
				// check if provinsi is the same with current provincy and timestamp is not null
				if itemProvinsi == provinsi && resp[j]["timestamp"] != "" {
					timestamp := f.parseTimestamp(resp[j]["timestamp"])
					// check if timestamp in range current week
					if timestamp.After(tempAggregate.WeekStart) && timestamp.Before(tempAggregate.WeekEnd) {
						price, _ := strconv.ParseFloat(resp[j]["price"], 64)
						size, _ := strconv.ParseFloat(resp[j]["size"], 64)
						tempAggregate.TotalPrice += price
						tempAggregate.TotalSize += size
						tempAggregate.ListPrice = append(tempAggregate.ListPrice, price)
						tempAggregate.ListSize = append(tempAggregate.ListSize, size)
					} else {
						// if timestamp not in range current week, aggregate current week data and reset tempAggregate with new data
						aggregate := f.getAggregate(tempAggregate)
						grouping[provinsi] = append(grouping[provinsi], aggregate)
						tempAggregate = f.assignTempAggregate(resp[j])
					}
				}
				if j == len(resp)-1 {
					aggregate := f.getAggregate(tempAggregate)
					grouping[provinsi] = append(grouping[provinsi], aggregate)
				}
			}
			if i == len(resp)-1 {
				aggregate := f.getAggregate(tempAggregate)
				grouping[provinsi] = append(grouping[provinsi], aggregate)
			}

		}
	}
	return grouping, err
}

func (f fetchUsecase) getCurrencyConverter() (float64, error) {
	converterCache, found := f.Cache.Get("IDR_USD")
	if found {
		return converterCache.(float64), nil
	}
	converter, err := f.FetchRepository.GetCurrencyConverter()
	if err != nil {
		return 0, err
	}
	f.Cache.Set("IDR_USD", converter["IDR_USD"], 1*time.Hour)
	return converter["IDR_USD"], nil
}

func (f fetchUsecase) parseTimestamp(timestamp string) time.Time {
	var parsedTime time.Time
	parse, _ := strconv.ParseInt(timestamp, 10, 64)
	if len(timestamp) == 10 {
		parsedTime = time.Unix(parse, 0)
	} else {
		parsedTime = time.Unix(0, parse*int64(time.Millisecond))
	}
	return parsedTime
}

func (f fetchUsecase) getWeekRange(currentDate string) (time.Time, time.Time) {
	parsedCurrentDate := f.parseTimestamp(currentDate)
	day := int(parsedCurrentDate.Weekday())
	weekStart := parsedCurrentDate.AddDate(0, 0, -day+1)
	weekEnd := parsedCurrentDate.AddDate(0, 0, 7-day)
	return weekStart, weekEnd
}

func (f fetchUsecase) getMedian(list []float64) float64 {
	middle := len(list) / 2
	var median float64
	if len(list)%2 == 0 {
		median = (list[middle] + list[middle-1]) / 2
	} else {
		median = list[middle]
	}
	return median
}

func (f fetchUsecase) getAggregate(tempAggregate fetch.TempAggregate) fetch.Aggregate {
	var aggregate fetch.Aggregate

	aggregate.AvgPrice = tempAggregate.TotalPrice / float64(len(tempAggregate.ListPrice))
	aggregate.AvgSize = tempAggregate.TotalSize / float64(len(tempAggregate.ListSize))

	sort.Float64s(tempAggregate.ListPrice)
	sort.Float64s(tempAggregate.ListSize)
	aggregate.MaxPrice = tempAggregate.ListPrice[len(tempAggregate.ListPrice)-1]
	aggregate.MinPrice = tempAggregate.ListPrice[0]
	aggregate.MaxSize = tempAggregate.ListSize[len(tempAggregate.ListSize)-1]
	aggregate.MinSize = tempAggregate.ListSize[0]
	aggregate.MedianPrice = f.getMedian(tempAggregate.ListPrice)
	aggregate.MedianSize = f.getMedian(tempAggregate.ListSize)
	startPeriod := tempAggregate.WeekStart.Format("2006-01-02")
	endPeriod := tempAggregate.WeekEnd.Format("2006-01-02")
	aggregate.Period = startPeriod + " - " + endPeriod
	aggregate.ListPrice = tempAggregate.ListPrice
	aggregate.ListSize = tempAggregate.ListSize
	return aggregate
}

func (f fetchUsecase) assignTempAggregate(resp map[string]string) fetch.TempAggregate {
	var tempAggregate fetch.TempAggregate
	intPrice, _ := strconv.ParseFloat(resp["price"], 64)
	intSize, _ := strconv.ParseFloat(resp["size"], 64)
	tempAggregate.TotalPrice = intPrice
	tempAggregate.ListPrice = []float64{intPrice}
	tempAggregate.TotalSize = intSize
	tempAggregate.ListSize = []float64{intSize}
	tempAggregate.WeekStart, tempAggregate.WeekEnd = f.getWeekRange(resp["timestamp"])
	return tempAggregate
}
