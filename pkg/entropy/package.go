/*

Copyright (C) 2017-2020  Daniele Rondina <geaaru@sabayonlinux.org>

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program. If not, see <http://www.gnu.org/licenses/>.

*/
package entropy

import (
	"database/sql"
	"errors"
	"fmt"
	"sort"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

type EntropyPackageDetail struct {
	Id           int
	Package      EntropyPackage
	Dependencies []*EntropyPackage
}

func retrievePackageDependencies(db *sql.DB, ans *EntropyPackageDetail) error {
	// Retrieve dependencies
	getDeps := fmt.Sprintf(`
SELECT DISTINCT d2.dependency
FROM dependencies d,
     dependenciesreference d2
WHERE d2.iddependency = d.iddependency
AND d.idpackage = %d`, ans.Id)

	rows, err := db.Query(getDeps)
	if err != nil {
		return errors.New("Error on retrieve dependencies: " + err.Error())
	}
	defer rows.Close()

	mdeps := make(map[string]bool)
	var dep string
	for rows.Next() {
		err = rows.Scan(&dep)
		if err != nil {
			return errors.New("Error on parse row for retrieve idpackage: " + err.Error())
		}
		idx := strings.Index(dep, ";")
		if idx > 0 {
			conditions := strings.Split(dep, ";")
			// Take last deps version to avoid use of major deps not availble in entropy
			dep = conditions[len(conditions)-1]
		}

		idx = strings.Index(dep, "[")
		if idx > 0 {
			dep = dep[0:idx]
		}

		// Drop slot. I will retrieve it later
		idx = strings.Index(dep, ":")
		if idx > 0 {
			dep = dep[0:idx]
		}

		// Fix =* not complaint to gentoo syntax
		if strings.HasPrefix(dep, "=*") {
			dep = "=" + dep[2:] + "*"
		}

		// Drop ?
		dep = strings.ReplaceAll(dep, "?", "")

		mdeps[dep] = true
	}

	dArr := make([]string, 0)
	for dkey, _ := range mdeps {
		dArr = append(dArr, dkey)
	}
	sort.Strings(dArr)
	for _, dkey := range dArr {
		d, err := NewEntropyPackage(dkey)
		if err != nil {
			return errors.New("Error on parse dependency " + dkey + ": " + err.Error())
		}
		ans.Dependencies = append(ans.Dependencies, d)
	}

	err = rows.Err()
	if err != nil {
		return errors.New("Error on parse row for retrieve idpackage: " + err.Error())
	}

	return nil
}

func getPackageDepDetail(db *sql.DB, pkg *EntropyPackage) error {

	if pkg.Version != "" {

		// Try to retrieve data with atom

		tmpD := &EntropyPackageDetail{
			Package:      *pkg,
			Dependencies: make([]*EntropyPackage, 0),
		}

		err := getPackageDataByAtom(db, tmpD)
		if err == nil {
			pkg.License = tmpD.Package.License
			pkg.Slot = tmpD.Package.Slot
			return nil
		}

	}

	getDepData := fmt.Sprintf(`SELECT slot,license
FROM baseinfo
WHERE name = '%s' and category = '%s'
ORDER BY version DESC LIMIT 1`, pkg.Name, pkg.Category)

	rows, err := db.Query(getDepData)
	if err != nil {
		return errors.New("Error on retrieve dependency data: " + err.Error())
	}
	defer rows.Close()
	rows.Next()

	var slot, license string
	err = rows.Scan(&slot, &license)
	if err != nil {
		return errors.New(
			fmt.Sprintf("Error on parse row for dependency detail for pkg %s/%s: ",
				pkg.Category, pkg.Name) + err.Error())
	}
	err = rows.Err()
	if err != nil {
		return errors.New("Error on fetch rows: " + err.Error())
	}

	pkg.License = license
	pkg.Slot = slot

	return nil
}

func getPackageUseFlags(db *sql.DB, pkg *EntropyPackage, id int) error {
	// Retrieve use flags
	getUserFlags := fmt.Sprintf(`
SELECT ur.flagname
FROM useflags u,
     useflagsreference ur
WHERE u.idflag = ur.idflag
AND u.idpackage = %d`, id)

	rows, err := db.Query(getUserFlags)
	if err != nil {
		return errors.New("Error on retrieve useflags: " + err.Error())
	}
	defer rows.Close()

	var flag string
	for rows.Next() {
		err = rows.Scan(&flag)
		if err != nil {
			return errors.New("Error on parse row for retrieve flags: " + err.Error())
		}

		pkg.UseFlags = append(pkg.UseFlags, flag)
	}

	err = rows.Err()
	if err != nil {
		return errors.New("Error on parse row for retrieve use flags: " + err.Error())
	}

	return nil
}

func getPackageDataByAtom(db *sql.DB, ans *EntropyPackageDetail) error {

	var idPackage int
	var name, version, slot, license string
	// Retrieve id package of the selected package as atom
	getIdPackageByAtom := fmt.Sprintf(
		"SELECT idpackage,name,version,slot,license FROM baseinfo WHERE atom = '%s/%s-%s'",
		ans.Package.Category, ans.Package.Name,
		fmt.Sprintf("%s%s", ans.Package.Version, ans.Package.VersionSuffix),
	)

	rows, err := db.Query(getIdPackageByAtom)
	if err != nil {
		return errors.New("Error on retrieve idpackage: " + err.Error())
	}
	defer rows.Close()

	rows.Next()
	err = rows.Scan(&idPackage, &name, &version, &slot, &license)
	if err != nil {
		return errors.New("Error on parse row for retrieve idpackage: " + err.Error())
	}
	err = rows.Err()
	if err != nil {
		return errors.New("Error on parse row for retrieve idpackage: " + err.Error())
	}

	ans.Id = idPackage
	ans.Package.License = license
	ans.Package.Slot = slot

	return nil
}

func getPackageDataByCatName(db *sql.DB, ans *EntropyPackageDetail) error {

	var idPackage int
	var name, version, slot, license string
	// Retrieve id package of the selected package as atom
	getIdPackageByCatName := fmt.Sprintf(`
SELECT idpackage,name,version,slot,license
FROM baseinfo
WHERE category = '%s' AND name = '%s'
ORDER BY version DESC LIMIT 1`, ans.Package.Category, ans.Package.Name)

	rows, err := db.Query(getIdPackageByCatName)
	if err != nil {
		return errors.New("Error on retrieve idpackage: " + err.Error())
	}
	defer rows.Close()

	rows.Next()
	err = rows.Scan(&idPackage, &name, &version, &slot, &license)
	if err != nil {
		return errors.New("Error on parse row for retrieve idpackage: " + err.Error())
	}
	err = rows.Err()
	if err != nil {
		return errors.New("Error on parse row for retrieve idpackage: " + err.Error())
	}

	ans.Id = idPackage
	ans.Package.License = license
	ans.Package.Slot = slot

	tmpPkg, err := NewEntropyPackage(
		fmt.Sprintf("%s/%s-%s", ans.Package.Category, ans.Package.Name, version))
	if err != nil {
		return errors.New("Unexpected error on elaborate version: " + err.Error())
	}

	ans.Package.Version = tmpPkg.Version
	ans.Package.VersionSuffix = tmpPkg.VersionSuffix

	return nil
}

func RetrievePackageData(pkg *EntropyPackage, dbpath string) (*EntropyPackageDetail, error) {
	ans := &EntropyPackageDetail{
		Package:      *pkg,
		Dependencies: make([]*EntropyPackage, 0),
	}

	// Open the connection
	db, err := sql.Open("sqlite3", dbpath)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	db.SetMaxOpenConns(1)

	if pkg.Version != "" {
		err := getPackageDataByAtom(db, ans)
		if err != nil {
			return nil, err
		}
	} else {
		err := getPackageDataByCatName(db, ans)
		if err != nil {
			return nil, err
		}
	}

	// Retrieve dependencies
	err = retrievePackageDependencies(db, ans)
	if err != nil {
		return nil, err
	}

	// Retrieve dependencies data
	for _, d := range ans.Dependencies {
		_ = getPackageDepDetail(db, d)
		// Ignore error for package with missing slot
	}

	// Retrieve use flags of package
	err = getPackageUseFlags(db, &ans.Package, ans.Id)

	return ans, nil
}

func RetrieveRepoPackages(dbpath string) ([]*EntropyPackage, error) {
	ans := []*EntropyPackage{}

	// Open the connection
	db, err := sql.Open("sqlite3", dbpath)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	db.SetMaxOpenConns(1)

	listDbPkgQuery := `
SELECT atom,slot
FROM baseinfo`

	rows, err := db.Query(listDbPkgQuery)
	if err != nil {
		return ans, errors.New("Error on retrieve pkg list: " + err.Error())
	}
	defer rows.Close()

	var atom, slot string
	for rows.Next() {
		err = rows.Scan(&atom, &slot)
		if err != nil {
			return ans, errors.New("Error on parse row for retrieve data: " + err.Error())
		}

		pkg, err := NewEntropyPackage(atom)
		if err != nil {
			return ans, errors.New("Error on parse atom " + atom + ": " + err.Error())
		}
		pkg.Slot = slot

		ans = append(ans, pkg)
	}

	err = rows.Err()
	if err != nil {
		return ans, errors.New("Error on parse row for retrieve data: " + err.Error())
	}

	return ans, nil
}
