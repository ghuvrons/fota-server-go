package controllers

import (
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/ghuvrons/fota-server-go/models"
	"github.com/ghuvrons/fota-server-go/utils"

	"github.com/google/uuid"
)

func UploadFirmware(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	if err := r.ParseMultipartForm(1024); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	version := r.FormValue("model_number")
	uploadedFile, _, err := r.FormFile("firmware")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer uploadedFile.Close()

	versionUint := utils.Version_StrToUint(version)
	if versionUint == 0 {
		http.Error(w, "Version is not valid", http.StatusInternalServerError)
		return
	}

	vm := models.VehicleModelFind(r.Context(), map[string]interface{}{
		"model_number": versionUint,
	})
	if vm == nil {
		http.Error(w, "Version is not valid", http.StatusInternalServerError)
		return
	}
	fmt.Println("vehicle model", vm)

	myFilePath := ""
	if err, myFilePath = saveFirmware(uploadedFile, version); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	vf := &models.VehicleFirmware{
		ModelId:    vm.Id,
		BinaryPath: myFilePath,
	}

	vf.Save(r.Context())

	fmt.Println("uploaded :", myFilePath)
	w.Write([]byte("done"))
}

func saveFirmware(uploadedFile multipart.File, fileVersion string) (error, string) {
	basePath, err := os.Getwd()
	if err != nil {
		return err, ""
	}

	// path naming
	uuidWithHyphen := uuid.New()
	myFilePath := filepath.Join("/_firmwares", fileVersion, uuidWithHyphen.String()+".bin")
	myFilePath = strings.ReplaceAll(myFilePath, "\\", "/")
	fileLocation := filepath.Join(basePath, myFilePath)

	// check directory and create if not exist
	fileInfo, err := os.Stat(filepath.Dir(fileLocation))
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			if err := os.Mkdir(filepath.Dir(fileLocation), 0755); err != nil {
				return err, ""
			}
		} else {
			return err, ""
		}
	} else if !fileInfo.IsDir() {
		return os.ErrInvalid, ""
	}

	// create file
	targetFile, err := os.OpenFile(fileLocation, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return err, ""
	}
	defer targetFile.Close()

	// write file
	if _, err := io.Copy(targetFile, uploadedFile); err != nil {
		return err, ""
	}

	return nil, myFilePath
}
