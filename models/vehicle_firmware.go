package models

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/ghuvrons/fota-server-go/utils"
)

type VehicleFirmwareType int

const (
	FIRMWARE_VCU VehicleFirmwareType = iota
	FIRMWARE_HMI
)

type VehicleFirmware struct {
	Id           int64
	ModelId      int
	FirmwareType VehicleFirmwareType
	Version      uint32
	Crc          uint32
	BinaryPath   string
	CreatedAt    time.Time
	inserted     bool
	FileLength   uint32
}

func VehicleFirmwareFind(ctx context.Context, conditions map[string]interface{}) *VehicleFirmware {
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

	rows, err := db.Query("SELECT * FROM vehicle_firmware where " + _where)
	if err != nil {
		return nil
	}
	defer rows.Close()

	if rows.Next() {
		var vf = &VehicleFirmware{}
		if vf.scanAll(rows) != nil {
			return nil
		}

		return vf
	}
	return nil
}

func VehicleFirmwareGetLatest(ctx context.Context, model_id int) *VehicleFirmware {
	db, isOK := ctx.Value("DB_MYSQL").(*sql.DB)
	if !isOK {
		return nil
	}

	query := fmt.Sprint("SELECT * FROM vehicle_firmware where `model_id` = ", model_id, " ORDER BY created_at DESC")

	rows, err := db.Query(query)
	if err != nil {
		return nil
	}
	defer rows.Close()

	if rows.Next() {
		var vf = &VehicleFirmware{}
		if vf.scanAll(rows) != nil {
			return nil
		}

		err, realPath, length := utils.GetFirmwareInfo(vf.BinaryPath)
		if err != nil {
			return nil
		}

		vf.BinaryPath = realPath
		vf.FileLength = length

		return vf
	}
	return nil
}

func (vf *VehicleFirmware) Save(ctx context.Context) error {
	db, isOK := ctx.Value("DB_MYSQL").(*sql.DB)
	if !isOK {
		return fmt.Errorf("DB Not Connected")
	}
	var fType string = "VCU"
	if vf.FirmwareType == FIRMWARE_HMI {
		fType = "HMI"
	}

	if !vf.inserted {
		vf.CreatedAt = time.Now()
		err, realPath, _ := utils.GetFirmwareInfo(vf.BinaryPath)
		if err != nil {
			return nil
		}
		vf.Crc = utils.CalculateCrcFile(realPath)

		result, err := db.Exec("INSERT INTO `vehicle_firmware` (`model_id`, `type`, `version`, `crc`, `bin_path`, `created_at`) VALUES (?, ?, ?, ?, ?, ?);",
			vf.ModelId, fType, vf.Version, vf.Crc, vf.BinaryPath, vf.CreatedAt.Format("2006-01-02 15:04:05"))
		if err != nil {
			return err
		}
		vf.Id, err = result.LastInsertId()
		vf.inserted = true
	}

	return nil
}

func (vf *VehicleFirmware) scanAll(rows *sql.Rows) error {
	var fType string
	err := rows.Scan(&vf.Id, &vf.ModelId, &fType, &vf.Version, &vf.Crc, &vf.BinaryPath, &vf.CreatedAt)
	if err == nil {
		vf.inserted = true
		vf.FirmwareType = FIRMWARE_VCU
		if fType == "HMI" {
			vf.FirmwareType = FIRMWARE_HMI
		}
	}

	return err
}
