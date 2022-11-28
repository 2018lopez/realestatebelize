//Filename: internal/data/report.go

package data

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

type TopAgents struct {
	AgentName         string  `json:"agent_name"`
	TotalPropertySold int64   `json:"total_property_sold"`
	TotalSales        float64 `json:"total_sales"`
}

type ListingsStatus struct {
	SoldLeased int64 `json:"sold_leased"`
	Available  int64 `json:"available"`
}

type TotalSales struct {
	TotalSales float64 `json:"total_sales"`
}

// Define a ReportModel which wrap a sql.DB connection pool
type ReportModel struct {
	DB *sql.DB
}

func (m ReportModel) GetTopAgents() ([]*TopAgents, error) {
	//construct query

	query := fmt.Sprintf(`
	select  u.fullname, count(l.id), sum(l.price) from users u inner join userproperties up on u.id = up.userid
	inner join listing l on l.id = up.listingid
	inner join propertystatus ps on ps.id = l.propertystatusid
	where ps.name='Sold'  group by u.fullname, l.price 
	order by l.price desc 
	limit 5
		`)
	//CREATE a 3 sec timeout context
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	//execute
	rows, err := m.DB.QueryContext(ctx, query)

	if err != nil {
		return nil, err
	}

	//close the result set
	defer rows.Close()

	//Initialize an empty slice to hold topAgents data
	topagents := []*TopAgents{}

	//iterate over the rows in the result set

	for rows.Next() {
		var topagent TopAgents
		//scan the values from row into  topagent struct
		err := rows.Scan(
			&topagent.AgentName,
			&topagent.TotalPropertySold,
			&topagent.TotalSales,
		)

		if err != nil {
			return nil, err
		}

		//add the topagent to our slice
		topagents = append(topagents, &topagent)

	}

	//Check for errors after looping through the result set

	if err = rows.Err(); err != nil {
		return nil, err
	}

	// return slice of topagents
	return topagents, nil

}

func (m ReportModel) GetListingStatus() ([]*ListingsStatus, error) {
	//construct query

	query := fmt.Sprintf(`
	select (select count(lg.id) from listing lg inner join propertystatus ps on lg.propertystatusid=ps.id where ps.name='Sold' or ps.Name='Leased') as SoldLeased,  
	(select count(lg.id) from listing lg inner join propertystatus ps on lg.propertystatusid=ps.id where ps.name='Available') as Available

		`)
	//CREATE a 3 sec timeout context
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	//execute
	rows, err := m.DB.QueryContext(ctx, query)

	if err != nil {
		return nil, err
	}

	//close the result set
	defer rows.Close()

	//Initialize an empty slice to hold topAgents data
	listingsstatus := []*ListingsStatus{}

	//iterate over the rows in the result set

	for rows.Next() {
		var listing ListingsStatus
		//scan the values from row into  topagent struct
		err := rows.Scan(
			&listing.SoldLeased,
			&listing.Available,
		)

		if err != nil {
			return nil, err
		}

		//add the topagent to our slice
		listingsstatus = append(listingsstatus, &listing)

	}

	//Check for errors after looping through the result set

	if err = rows.Err(); err != nil {
		return nil, err
	}

	// return slice of topagents
	return listingsstatus, nil

}

// Get Total Sales of all properties Sold/leased
func (m ReportModel) GetTotalSales() ([]*TotalSales, error) {
	//construct query

	query := fmt.Sprintf(`
	SELECT sum(l.price) as totalSales from listing l inner join propertystatus ps on l.propertystatusid=ps.id 
	where ps.name ='Sold' OR ps.name='Leased'

		`)
	//CREATE a 3 sec timeout context
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	//execute
	rows, err := m.DB.QueryContext(ctx, query)

	if err != nil {
		return nil, err
	}

	//close the result set
	defer rows.Close()

	//Initialize an empty slice to hold topAgents data
	totalsales := []*TotalSales{}

	//iterate over the rows in the result set

	for rows.Next() {
		var totalsale TotalSales
		//scan the values from row into  topagent struct
		err := rows.Scan(
			&totalsale.TotalSales,
		)

		if err != nil {
			return nil, err
		}

		//add the topagent to our slice
		totalsales = append(totalsales, &totalsale)

	}

	//Check for errors after looping through the result set

	if err = rows.Err(); err != nil {
		return nil, err
	}

	// return slice of totalsales
	return totalsales, nil

}
