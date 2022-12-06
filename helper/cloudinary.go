package helper

import (
	"context"
	"fmt"
	"mime/multipart"
	"os"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

func CloudinaryUpload(ctx context.Context, source multipart.File, userId string) (string, error) {
	cloudinaryCloud := os.Getenv("CLOUDINARY_CLOUD")
	cloudinaryKey := os.Getenv("CLOUDINARY_KEY")
	cloudinarySecret := os.Getenv("CLOUDINARY_SECRET")

	cld, _ := cloudinary.NewFromParams(cloudinaryCloud, cloudinaryKey, cloudinarySecret)

	// Upload image and set the PublicID to userId.
	resp, err := cld.Upload.Upload(
		ctx,
		source,
		uploader.UploadParams{
			PublicID: fmt.Sprintf("user-%s", userId),
			Format:   "jpg",
			Folder:   "better-space/testing/photo",
		},
	)

	url := resp.SecureURL

	return url, err
}