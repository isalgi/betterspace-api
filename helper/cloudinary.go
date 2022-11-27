package helper

import (
	_util "backend/utils"
	"context"
	"fmt"
	"mime/multipart"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

func CloudinaryUpload(ctx context.Context, source multipart.File, userId string) (string, error) {
	cloudinaryCloud := _util.GetConfig("CLOUDINARY_CLOUD")
	cloudinaryKey := _util.GetConfig("CLOUDINARY_KEY")
	cloudinarySecret := _util.GetConfig("CLOUDINARY_SECRET")

	cld, _ := cloudinary.NewFromParams(cloudinaryCloud, cloudinaryKey, cloudinarySecret)

	// Upload image and set the PublicID to userId.
	resp, err := cld.Upload.Upload(
		ctx,
		source,
		uploader.UploadParams{
			PublicID: fmt.Sprintf("user-%s", userId),
			Format:   "jpg",
			Folder:   "office-booking-profile-photo-user",
		},
	)

	url := resp.SecureURL

	return url, err
}