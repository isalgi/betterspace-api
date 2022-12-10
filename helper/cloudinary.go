package helper

import (
	_utils "backend/utils"
	"context"
	"fmt"
	"log"
	"mime/multipart"
	"os"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

var (
	cloudinaryCloud string = os.Getenv("CLOUDINARY_CLOUD")
	cloudinaryKey string = os.Getenv("CLOUDINARY_KEY")
	cloudinarySecret string = os.Getenv("CLOUDINARY_SECRET")
)

func CloudinaryUpload(ctx context.Context, source multipart.File, userId string) (string, error) {
	cld, _ := cloudinary.NewFromParams(cloudinaryCloud, cloudinaryKey, cloudinarySecret)

	// Upload image and set the PublicID to userId.
	resp, err := cld.Upload.Upload(
		ctx,
		source,
		uploader.UploadParams{
			PublicID: fmt.Sprintf("user-%s", userId),
			Format:   "jpg",
			Folder:   "better-space/staging/photo",
		},
	)

	url := resp.SecureURL

	return url, err
}

func CloudinaryUploadOfficeImgs(files []*multipart.FileHeader) ([]string, error) {
	ctx := context.Background()

	cld, _ := cloudinary.NewFromParams(cloudinaryCloud, cloudinaryKey, cloudinarySecret)
	
	var imageURLs []string
	var err error

	for i := len(files) - 1; i >= 0; i-- {
		src, err := files[i].Open()
		
		if err != nil {
			log.Println(err)
			return imageURLs, err
		}

		fileName := _utils.RandomString(25)

		// upload image and set the PublicID to fileName.
		resp, err := cld.Upload.Upload(
			ctx,
			src,
			uploader.UploadParams{
				PublicID: fileName,
				Format:   "jpg",
				Folder:   "better-space/staging/office-image",
			},
		)

		if err != nil {
			log.Println(err)
			return imageURLs, err
		}

		url := resp.SecureURL

		imageURLs = append(imageURLs, url)

		defer src.Close()
	}
	
	return imageURLs, err
}