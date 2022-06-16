package models

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

type VehicleModel struct {
	Id          int
	Name        string
	ModelNumber uint32
}

func VehicleModelFind(ctx context.Context, conditions map[string]interface{}) *VehicleModel {
	db, isOK := ctx.Value("DB_MYSQL").(*sql.DB)
	if !isOK {
		return nil
	}

	_where := ""
	for key, val := range conditions {
		if _where != "" {
			_where += " AND"
		}
		valStr, isString := val.(string)
		if isString {
			_where += fmt.Sprint(" ", key, " = ", "\"", valStr, "\"")
		} else {
			_where += fmt.Sprint(" ", key, " = ", val)
		}
	}

	rows, err := db.Query("SELECT * FROM vehicle_model where " + _where)
	if err != nil {
		return nil
	}
	defer rows.Close()

	if rows.Next() {
		var vm = &VehicleModel{}
		var err = rows.Scan(&vm.Id, &vm.Name, &vm.ModelNumber)
		if err != nil {
			return nil
		}

		return vm
	}
	return nil
}

func VehicleModelGetLatestFirmware(ctx context.Context, model_number uint32, firmwareType VehicleFirmwareType) *VehicleFirmware {
	// TODO :
	// 1 : check cache
	// 		return from cache if exist
	// 2 : get from DB
	// 3 : check filepath info by GetFirmwareInfo()
	// 4 : save in cache

	var firmware *VehicleFirmware

	firmwares, isModelFound := firmwareCache[model_number]
	if isModelFound {
		var isOK bool
		firmware, isOK = firmwares[firmwareType]
		if isOK {
			return firmware
		}
	}

	db, isOK := ctx.Value("DB_MYSQL").(*sql.DB)
	if !isOK {
		var err error
		db, err = CreateDBConnection()
		if err != nil {
			return nil
		}
		defer func() {
			db.Close()
		}()
	}

	// get model
	vm := VehicleModelFind(ctx, map[string]interface{}{
		"model_number": model_number,
	})

	if vm == nil {
		return nil
	}

	// get firmware
	vf := VehicleFirmwareGetLatest(ctx, vm.Id)

	if !isModelFound {
		firmwareCache[model_number] = map[VehicleFirmwareType]*VehicleFirmware{}
		firmwares = firmwareCache[model_number]
	}

	firmwares[firmwareType] = vf

	return vf
}
