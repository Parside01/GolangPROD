package postgres

import (
	"slices"
	"solution/models"
)

var (
	getCountriesByRegion = "SELECT name, alpha2, alpha3, region FROM countries WHERE alpha2=$1"
	getAll               = "SELECT name, alpha2, alpha3, region FROM countries"
	isExist              = "SELECT EXISTS(SELECT 1 FROM countries WHERE alpha2 = $1)"
	byregion             = "SELECT name, alpha2, alpha3, region FROM countries WHERE region = $1"
)

func (p *PostgresDB) GetCountries(filter string) (*models.Country, error) {
	c := new(models.Country)
	err := p.db.QueryRow(getCountriesByRegion, filter).Scan(&c.Name, &c.Alpha2, &c.Alpha3, &c.Region)
	if err != nil {
		p.logger.Error("PostgresDB.GetCountries: rows.Scan", "error", err)
		return nil, err
	}

	return c, nil
}

func (p *PostgresDB) GetCountriesByRegion(region string) ([]*models.Country, error) {
	if region == "" {
		res, err := p.GetAllCountries()
		if err != nil {
			return nil, err
		}
		slices.SortFunc(res, func(a, b *models.Country) int {
			if a.Alpha2 < b.Alpha2 {
				return -1
			}
			if a.Alpha2 > b.Alpha2 {
				return 1
			}
			return 0
		})
		return res, nil
	}

	rows, err := p.db.Query(byregion, region)
	if err != nil {
		return nil, err
	}
	res := []*models.Country{}
	for rows.Next() {
		c := new(models.Country)
		err := rows.Scan(&c.Name, &c.Alpha2, &c.Alpha3, &c.Region)
		if err != nil {
			p.logger.Error("PostgresDB.GetCountries: rows.Scan", "error", err)
			return nil, err
		}
		res = append(res, c)
	}

	slices.SortFunc(res, func(a, b *models.Country) int {
		if a.Alpha2 < b.Alpha2 {
			return -1
		}
		if a.Alpha2 > b.Alpha2 {
			return 1
		}
		return 0
	})

	return res, nil
}

// В @alpha2 прилетает только 2 буквы, что-то другое не пойдет.
func (p *PostgresDB) CountryIsExist(alpha2 string) (bool, error) {
	if alpha2 == "" {
		return true, nil
	}
	row := p.db.QueryRow(isExist, alpha2)
	var exist bool
	err := row.Scan(&exist)
	if err != nil {
		p.logger.Error("PostgresDB.IsExist: row.Scan", "error", err)
		return false, err
	}
	return exist, nil
}

func (p *PostgresDB) GetAllCountries() ([]*models.Country, error) {
	rows, err := p.db.Query(getAll)
	if err != nil {
		return nil, err
	}

	res := []*models.Country{}
	for rows.Next() {
		c := new(models.Country)
		err := rows.Scan(&c.Name, &c.Alpha2, &c.Alpha3, &c.Region)
		if err != nil {
			p.logger.Error("PostgresDB.GetCountries: rows.Scan", "error", err)
			return nil, err
		}
		res = append(res, c)
	}

	defer rows.Close()

	return res, nil
}
