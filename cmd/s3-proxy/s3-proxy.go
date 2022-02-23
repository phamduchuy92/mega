package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/gabriel-vasile/mimetype"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gosimple/slug"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/confmap"
	"github.com/knadh/koanf/providers/file"
	minio "github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"gitlab.com/emi2/mega/internal/app"
)

var mc *minio.Client

func setupRoutes(app *fiber.App) {
	// + eptw regions endpoints
	app.Post("api/file-upload", FileUpload)
	app.Post("api/file-upload/:bucket", FileUpload)

	app.Delete("api/file-upload/:name", FileRemove)
	app.Delete("api/file-upload/:bucket/:name", FileRemove)

	app.Get("api/statics/:name", FileDownload)
	app.Get("api/statics/:bucket/:name", FileDownload)

	app.Get("api/file-info/:bucket/:name", FileCheck)
	app.Get("api/file-info/:name", FileCheck)

	app.Get("api/file-list/:bucket", FileList)
}

// configure application runtime
func configure() {
	// koanf defautl values
	app.Config.Load(confmap.Provider(map[string]interface{}{
		"http.listen":        ":3008",
		"s3.endpoint":        "localhost:9091",
		"s3.accessKeyID":     "Q3AM3UQ867SPQQA43P2F",
		"s3.secretAccessKey": "zuf+tfteSlswRu7BJ86wekitnifILbZam1KYY3TG",
		"s3.useSSL":          true,
		"s3.bucket":          "my-bucket",
	}, "."), nil)
	// override configuration with YAML
	app.Config.Load(file.Provider("configs/s3-proxy.yaml"), yaml.Parser())
}

// main function
func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	configure()
	srv := fiber.New(fiber.Config{
		ErrorHandler: app.ProblemJSONErrorHandle,
	})
	srv.Use(logger.New())
	// Initialize minio client object.
	minioClient, err := minio.New(app.Config.String("s3.endpoint"), &minio.Options{
		Creds:  credentials.NewStaticV4(app.Config.String("s3.accessKeyID"), app.Config.String("s3.secretAccessKey"), ""),
		Secure: app.Config.Bool("s3.useSSL"),
	})
	if err != nil {
		log.Fatalln(err)
	}
	mc = minioClient
	setupRoutes(srv)

	log.Fatal(srv.Listen(app.Config.String("http.listen")))
}

// FileUpload upload a file into minio bucket
func FileUpload(c *fiber.Ctx) error {
	bucketName := c.Params("bucket", app.Config.String("s3.bucket"))
	if exists, ok := mc.BucketExists(c.Context(), bucketName); ok == nil {
		if !exists {
			err := mc.MakeBucket(c.Context(), bucketName, minio.MakeBucketOptions{})
			if err != nil {
				return err
			}
		}
	} else {
		return fiber.ErrBadGateway
	}
	results := make([]minio.UploadInfo, 0)
	// Parse the multipart form:
	form, err := c.MultipartForm()
	if err != nil {
		return err
	}
	// => *multipart.Form

	// Get all files from "documents" key:
	files := form.File["files"]
	// => []*multipart.FileHeader

	// Loop through files:
	for _, uploadedFile := range files {
		fmt.Println(uploadedFile.Filename, uploadedFile.Size, uploadedFile.Header["Content-Type"][0])
		// => "tutorial.pdf" 360641 "application/pdf"

		// Save the files to disk:
		fileName := c.Get("name", uploadedFile.Filename)
		filePath := os.TempDir() + "/" + fileName
		err := c.SaveFile(uploadedFile, filePath)

		// Check for errors
		if err != nil {
			return err
		}

		// + upload file to bucket
		file, err := os.Open(filePath)
		if err != nil {
			return err
		}
		defer file.Close()

		fileStat, err := file.Stat()
		if err != nil {
			return err
		}

		fileMime := uploadedFile.Header.Get("Content-Type")
		if fileMime == "" {
			mime, err := mimetype.DetectReader(file)
			if err != nil {
				fileMime = mime.String()
			}
		}

		// TODO: Check if object exists, slug the basepath, and append suffix
		ext := filepath.Ext(fileName)
		baseName := slug.Make(strings.TrimSuffix(filepath.Base(fileName), ext))
		suffix := 0
		newFileName := fmt.Sprintf("%s%s", baseName, ext)
		for {
			_, err := mc.StatObject(context.Background(), bucketName, newFileName, minio.StatObjectOptions{})
			if err != nil {
				uploadInfo, err := mc.PutObject(context.Background(), bucketName, newFileName, file, fileStat.Size(), minio.PutObjectOptions{ContentType: fileMime})
				if err == nil {
					results = append(results, uploadInfo)
					break
				}
			}
			suffix++
			newFileName = fmt.Sprintf("%s-%d%s", baseName, suffix, ext)
		}
	}
	return c.JSON(results)
	// // Get first file from form field "document":
	// uploadedFile, err := c.FormFile("file")
	// if err != nil {
	// 	return err
	// }
	// log.Printf("Got file header %+v", uploadedFile.Header)
	// fileName := c.Get("name", uploadedFile.Filename)
	// filePath := os.TempDir() + "/" + fileName
	// err = c.SaveFile(uploadedFile, filePath)
	// if err != nil {
	// 	return err
	// }
	// // + upload file to bucket
	// file, err := os.Open(filePath)
	// if err != nil {
	// 	return err
	// }
	// defer file.Close()

	// fileStat, err := file.Stat()
	// if err != nil {
	// 	return err
	// }

	// fileMime := uploadedFile.Header.Get("Content-Type")
	// if fileMime == "" {
	// 	mime, err := mimetype.DetectReader(file)
	// 	if err != nil {
	// 		fileMime = mime.String()
	// 	}
	// }

	// // TODO: Check if object exists, slug the basepath, and append suffix
	// ext := filepath.Ext(fileName)
	// baseName := slug.Make(strings.TrimSuffix(filepath.Base(fileName), ext))
	// suffix := 0
	// newFileName := fmt.Sprintf("%s%s", baseName, ext)
	// for {
	// 	_, err := mc.StatObject(context.Background(), bucketName, newFileName, minio.StatObjectOptions{})
	// 	if err != nil {
	// 		uploadInfo, err := mc.PutObject(context.Background(), bucketName, newFileName, file, fileStat.Size(), minio.PutObjectOptions{ContentType: fileMime})
	// 		if err == nil {
	// 			return c.JSON(uploadInfo)
	// 		}
	// 	}
	// 	suffix++
	// 	newFileName = fmt.Sprintf("%s-%d%s", baseName, suffix, ext)
	// }
}

// FileDownload download a file from minio
func FileDownload(c *fiber.Ctx) error {
	bucketName := c.Params("bucket", c.Get("bucket", app.Config.String("s3.bucket")))
	fileName := c.Params("name", c.Get("name"))
	if fileName == "" {
		return fiber.ErrBadRequest
	}
	obj, err := mc.GetObject(context.Background(), bucketName, fileName, minio.GetObjectOptions{})
	if err != nil {
		return err
	}
	fileStat, err := obj.Stat()
	if err != nil {
		return err
	}
	c.Set(fiber.HeaderContentType, fileStat.ContentType)
	return c.SendStream(obj, int(fileStat.Size))
}

// FileCheck check if one object exists
func FileCheck(c *fiber.Ctx) error {
	bucketName := c.Params("bucket", c.Get("bucket", app.Config.String("s3.bucket")))
	fileName := c.Params("name", c.Get("name"))
	if fileName == "" {
		return fiber.ErrBadRequest
	}
	objInfo, err := mc.StatObject(context.Background(), bucketName, fileName, minio.StatObjectOptions{})
	if err != nil {
		return err
	}
	return c.JSON(objInfo)
}

// FileRemove remove one file from specified bucket
func FileRemove(c *fiber.Ctx) error {
	bucketName := c.Params("bucket", c.Get("bucket", app.Config.String("s3.bucket")))
	fileName := c.Params("name", c.Get("name"))
	if fileName == "" {
		return fiber.ErrBadRequest
	}
	err := mc.RemoveObject(c.Context(), bucketName, fileName, minio.RemoveObjectOptions{})
	if err != nil {
		return err
	}
	return c.SendStatus(fiber.StatusNoContent)
}

func FileList(c *fiber.Ctx) error {
	bucketName := c.Params("bucket", c.Get("bucket", app.Config.String("s3.bucket")))
	prefix := c.Query("prefix")
	objectList := make([]minio.ObjectInfo, 0)
	if bucketName == "" {
		return fiber.ErrBadRequest
	}
	objectCh := mc.ListObjects(c.Context(), bucketName, minio.ListObjectsOptions{
		Prefix:    prefix,
		Recursive: true,
	})
	for object := range objectCh {
		if object.Err != nil {
			return object.Err
		}
		objectList = append(objectList, object)
	}
	c.Set("X-Total-Count", fmt.Sprint(len(objectList)))
	return c.JSON(objectList)
}
