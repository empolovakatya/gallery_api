package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/disintegration/imaging"
	"github.com/gorilla/mux"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

type Photo struct {
	ID      int    `json:"id"`
	Image   string `json:"image"`
	Preview string `json:"preview"`
}

func Contains(a []string, x string) bool {
	for _, n := range a {
		if x == n {
			return true
		}
	}
	return false
}
func deleteFile(file string) {
	os.Remove(file)
}

func getPhotos(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	db, err := sql.Open("sqlite3", "./pkg/data.db")
	if err != nil {
		d := map[string]string{"status": err.Error()}
		mapVar2, _ := json.Marshal(d)
		http.Error(w, string(mapVar2), http.StatusInternalServerError)
		return
	}
	defer db.Close()
	rows, err := db.Query("select * from Photos")
	if err != nil {
		d := map[string]string{"status": err.Error()}
		mapVar2, _ := json.Marshal(d)
		http.Error(w, string(mapVar2), http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	var photos []Photo
	for rows.Next() {
		p := Photo{}
		err := rows.Scan(&p.ID, &p.Image, &p.Preview)
		if err != nil {
			fmt.Println(err)
			continue
		}
		photos = append(photos, p)
	}

	json.NewEncoder(w).Encode(photos)

}

func getPhoto(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	db, err := sql.Open("sqlite3", "./pkg/data.db")
	if err != nil {
		d := map[string]string{"status": err.Error()}
		mapVar2, _ := json.Marshal(d)
		http.Error(w, string(mapVar2), http.StatusInternalServerError)
		return
	}
	defer db.Close()
	params := mux.Vars(r)
	id := params["id"]
	row := db.QueryRow("select * from Photos where id = $1", id)
	p := Photo{}
	err = row.Scan(&p.ID, &p.Image, &p.Preview)
	if err != nil {
		d := map[string]string{"status": err.Error()}
		mapVar2, _ := json.Marshal(d)
		http.Error(w, string(mapVar2), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(p)
}

func uploadPhoto(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	file, fileHeader, err := r.FormFile("image")
	if err != nil {
		d := map[string]string{"status": err.Error()}
		mapVar2, _ := json.Marshal(d)
		http.Error(w, string(mapVar2), http.StatusBadRequest)
		return
	}

	defer file.Close()

	err = os.MkdirAll("./ui/static/img/uploads", os.ModePerm)
	if err != nil {
		d := map[string]string{"status": err.Error()}
		mapVar2, _ := json.Marshal(d)
		http.Error(w, string(mapVar2), http.StatusInternalServerError)
		return
	}
	err = os.MkdirAll("./ui/static/img/previews", os.ModePerm)
	if err != nil {
		d := map[string]string{"status": err.Error()}
		mapVar2, _ := json.Marshal(d)
		http.Error(w, string(mapVar2), http.StatusInternalServerError)
		return
	}

	fileName := time.Now().UnixNano()
	fileType := filepath.Ext(fileHeader.Filename)

	image_types := []string{".bmp", ".ecw", ".gif", ".ico", ".ilbm", ".jpeg", ".jpg", ".png", ".psd", ".tga", ".tiff"}

	if Contains(image_types, fileType) {
		dst, err := os.Create(fmt.Sprintf("./ui/static/img/uploads/%d%s", fileName, fileType))
		if err != nil {
			d := map[string]string{"status": err.Error()}
			mapVar2, _ := json.Marshal(d)
			http.Error(w, string(mapVar2), http.StatusInternalServerError)
			return
		}
		defer dst.Close()

		_, err = io.Copy(dst, file)
		if err != nil {
			d := map[string]string{"status": err.Error()}
			mapVar2, _ := json.Marshal(d)
			http.Error(w, string(mapVar2), http.StatusInternalServerError)
			return
		}
		oldImg, err := imaging.Open(fmt.Sprintf("./ui/static/img/uploads/%d%s", fileName, fileType))
		fileNameApi := "./ui/static/img/uploads/" + strconv.FormatInt(fileName, 10) + fileType
		imagePreview := imaging.Resize(oldImg, 128, 0, imaging.Lanczos)

		err = imaging.Save(imagePreview, fmt.Sprintf("./ui/static/img/previews/%d_preview.jpg", fileName))
		fileNamePreviewApi := "./ui/static/img/previews/" + strconv.FormatInt(fileName, 10) + "_preview.jpg"
		if err != nil {
			log.Fatalf("failed to save image: %v", err)
		}

		db, err := sql.Open("sqlite3", "./pkg/data.db")
		if err != nil {
			d := map[string]string{"status": err.Error()}
			mapVar2, _ := json.Marshal(d)
			http.Error(w, string(mapVar2), http.StatusInternalServerError)
			return
		}
		defer db.Close()
		photo, err := db.Exec("insert into photos (image, preview) values ($1, $2)",
			fileNameApi, fileNamePreviewApi)
		log.Println(photo)
		if err != nil {
			d := map[string]string{"status": err.Error()}
			mapVar2, _ := json.Marshal(d)
			http.Error(w, string(mapVar2), http.StatusInternalServerError)
			return
		}
		d := map[string]string{"status": "upload successful"}
		mapVar2, _ := json.Marshal(d)

		fmt.Fprintf(w, string(mapVar2))
	} else {
		d := map[string]string{"status": "wrong file format"}
		mapVar2, _ := json.Marshal(d)
		w.Header().Set("Content-Type", "application/json")
		http.Error(w, string(mapVar2), http.StatusInternalServerError)
		return
	}
}

func deletePhoto(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	id := params["id"]
	db, err := sql.Open("sqlite3", "./pkg/data.db")
	if err != nil {
		d := map[string]string{"status": err.Error()}
		mapVar2, _ := json.Marshal(d)
		http.Error(w, string(mapVar2), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	row := db.QueryRow("select * from Photos where id = $1", id)
	p := Photo{}
	err = row.Scan(&p.ID, &p.Image, &p.Preview)
	if err != nil {
		d := map[string]string{"status": err.Error()}
		mapVar2, _ := json.Marshal(d)
		http.Error(w, string(mapVar2), http.StatusInternalServerError)
		return
	}
	deleteFile(p.Image)
	deleteFile(p.Preview)
	result, err := db.Exec("delete from photos where id = $1", id)
	if err != nil {
		d := map[string]string{"status": err.Error()}
		mapVar2, _ := json.Marshal(d)
		http.Error(w, string(mapVar2), http.StatusInternalServerError)
		return
	}
	log.Println(result)
	d := map[string]string{"status": "delete successful"}
	mapVar2, _ := json.Marshal(d)
	fmt.Fprintf(w, string(mapVar2))
}
