package offices

import (
	"backend/businesses/offices"
	"fmt"
	"sort"
	"strconv"
	"strings"

	"gorm.io/gorm"
)

type officeRepository struct {
	conn *gorm.DB
}

func NewMySQLRepository(conn *gorm.DB) offices.Repository {
	return &officeRepository{
		conn: conn,
	}
}

func (or *officeRepository) GetAll() []offices.Domain {
	var rec []Office

	or.conn.Find(&rec)
	
	var imgsUrlPerID []imgs

	queryGetImgs := "SELECT `offices`.`id`, " + 
		"GROUP_CONCAT(office_images.url ORDER BY office_images.id SEPARATOR ' , ') AS images " + 
		"FROM offices " + 
		"INNER JOIN office_images on offices.id = office_images.office_id " + 
		"GROUP BY offices.id"
	or.conn.Raw(queryGetImgs).Scan(&imgsUrlPerID)

	var officeFacilitiesPerID []facilities

	queryGetFacilities := "SELECT `offices`.`id`, " + 
		"GROUP_CONCAT(`office_facilities`.`facilities_id` ORDER BY `office_facilities`.`facilities_id` SEPARATOR ' , ') AS f_id, " + 
		"GROUP_CONCAT(`facilities`.`description` ORDER BY `office_facilities`.`facilities_id` SEPARATOR ' , ') AS f_desc, " +
		"GROUP_CONCAT(`facilities`.`slug` ORDER BY `office_facilities`.`facilities_id` SEPARATOR ' , ') AS f_slug " + 
		"FROM `offices` " + 
		"INNER JOIN `office_facilities` ON `offices`.`id`=`office_facilities`.`office_id` " + 
		"INNER JOIN `facilities` ON `office_facilities`.`facilities_id`=`facilities`.`id` " + 
		"GROUP BY `offices`.`id`"
	or.conn.Raw(queryGetFacilities).Scan(&officeFacilitiesPerID)

	var totalBooked []totalbooked
	queryGetTotalBooked := "SELECT `office_id`, COUNT(*) AS total_booked FROM `transactions` WHERE `status` NOT IN ('rejected', 'cancelled') GROUP BY `office_id`"
	or.conn.Raw(queryGetTotalBooked).Scan(&totalBooked)

	var rateScore []ratescore
	queryGetRateScore := "SELECT `office_id`, ROUND(AVG(`score`), 1) FROM `reviews` GROUP BY `office_id`;"
	or.conn.Raw(queryGetRateScore).Scan(&rateScore)


	officeDomain := []offices.Domain{}
	
	for _, office := range rec {
		for _, v := range imgsUrlPerID {
			if strconv.Itoa(int(office.ID)) == v.Id {
				url := v.Images
				img := strings.Split(url, " , ")
				office.Images = img
			}
		}

		for _, fac := range officeFacilitiesPerID {
			if strconv.Itoa(int(office.ID)) == fac.Id {
				f_id := fac.F_id
				facilitesId := strings.Split(f_id, " , ")
				f_desc := fac.F_desc
				facilitesDesc := strings.Split(f_desc, " , ")
				f_slug := fac.F_slug
				facilitiesSlug := strings.Split(f_slug, " , ")

				office.FacilitiesId =  facilitesId
				office.FacilitiesDesc = facilitesDesc
				office.FacilitesSlug = facilitiesSlug
			}
		}

		for _, b := range totalBooked {
			if strconv.Itoa(int(office.ID)) == b.OfficeId {
				office.TotalBooked = b.TotalBooked
			}
		}

		for _, r := range rateScore {
			if strconv.Itoa(int(office.ID)) == r.OfficeId {
				office.Rate = r.Score
			}
		}

		officeDomain = append(officeDomain, office.ToDomain())
	}

	return officeDomain
}

func (or *officeRepository) GetByID(id string) offices.Domain {
	var office Office

	or.conn.First(&office, "id = ?", id)
	
	var imagesString string
	
	// get office images
	querySQL := fmt.Sprintf("SELECT GROUP_CONCAT(office_images.url ORDER BY office_images.id SEPARATOR ' , ') AS images " + 
		"FROM offices " + 
		"INNER JOIN office_images on offices.id = office_images.office_id " + 
		"WHERE `offices`.`id` = %s " + 
		"GROUP BY offices.id", id)

	or.conn.Raw(querySQL).Scan(&imagesString)

	img := strings.Split(imagesString, " , ")
	office.Images = img

	var fac facilities

	querySQL = fmt.Sprintf("SELECT `offices`.`id`, " + 
		"GROUP_CONCAT(`office_facilities`.`facilities_id` ORDER BY `office_facilities`.`facilities_id` SEPARATOR ' , ') AS f_id, " + 
		"GROUP_CONCAT(`facilities`.`description` ORDER BY `office_facilities`.`facilities_id` SEPARATOR ' , ') AS f_desc, " + 
		"GROUP_CONCAT(`facilities`.`slug` ORDER BY `office_facilities`.`facilities_id` SEPARATOR ' , ') AS f_slug " + 
		"FROM `offices` " + 
		"INNER JOIN `office_facilities` ON `offices`.`id`=`office_facilities`.`office_id` " + 
		"INNER JOIN `facilities` ON `office_facilities`.`facilities_id`=`facilities`.`id` " + 
		"WHERE `offices`.`id` = %s", id)

	or.conn.Raw(querySQL).Scan(&fac)

	f_id := fac.F_id
	facilitesId := strings.Split(f_id, " , ")
	f_desc := fac.F_desc
	facilitesDesc := strings.Split(f_desc, " , ")
	f_slug := fac.F_slug
	facilitiesSlug := strings.Split(f_slug, " , ")

	office.FacilitiesId =  facilitesId
	office.FacilitiesDesc = facilitesDesc
	office.FacilitesSlug = facilitiesSlug

	var count int64

	or.conn.Table("transactions").Not(map[string]interface{}{"status": []string{"rejected", "cancelled"}}).Where("office_id = ?", office.ID).Count(&count)

	office.TotalBooked = count

	var rate_score float64

	or.conn.Table("reviews").Where("office_id = ?", office.ID).Select("round(avg(`score`), 1)").Scan(&rate_score)

	office.Rate = rate_score

	return office.ToDomain()
}

func (or *officeRepository) Create(officeDomain *offices.Domain) offices.Domain {
	var result *gorm.DB

	rec := FromDomain(officeDomain)

	facilitiesIdList := []int{}

	for _, v := range rec.FacilitiesId {
		id, _ := strconv.Atoi(v)
		facilitiesIdList = append(facilitiesIdList, id)
	}

	sort.Sort(sort.Reverse(sort.IntSlice(facilitiesIdList)))

	for _, v := range facilitiesIdList {
		if err := or.conn.Exec(fmt.Sprintf("SELECT * FROM `facilities` WHERE `id` = %d", v)).Error; err != nil {
			return rec.ToDomain()
		}
	}

	err := or.conn.Transaction(func(tx *gorm.DB) error {
		result = tx.Create(&rec)
		result.Last(&rec)
		
		// insert to pivot table `office_images`
		for _, v := range rec.Images {
			querySQL := fmt.Sprintf("INSERT INTO `office_images`(`url`, `office_id`) VALUES ('%s', '%s')", v, strconv.Itoa(int(rec.ID)))
			
			if err := tx.Table("office_images").Exec(querySQL).Error; err != nil {
				return err
			}
		}

		// insert to pivot table `office_facilities`
		for _, v := range facilitiesIdList {
			querySQL := fmt.Sprintf("INSERT INTO `office_facilities`(`facilities_id`, `office_id`) VALUES ('%d','%d')", v, rec.ID)
			if err := tx.Table("office_facilities").Exec(querySQL).Error; err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		rec.ID = 0
		return rec.ToDomain()
	}

	return rec.ToDomain()
}

func (or *officeRepository) Update(id string, officeDomain *offices.Domain) offices.Domain {
	var office offices.Domain = or.GetByID(id)

	if office.ID == 0 {
		return office
	}

	updatedOffice := FromDomain(&office)

	facilitiesIdList := []int{}

	for _, v := range officeDomain.FacilitiesId {
		id, _ := strconv.Atoi(v)
		facilitiesIdList = append(facilitiesIdList, id)
	}

	sort.Sort(sort.Reverse(sort.IntSlice(facilitiesIdList)))

	for _, v := range facilitiesIdList {
		if err := or.conn.Exec(fmt.Sprintf("SELECT * FROM `facilities` WHERE `id` = %d", v)).Error; err != nil {
			office.ID = 0
			return office
		}
	}

	err := or.conn.Transaction(func(tx *gorm.DB) error {
		if len(officeDomain.Images) != 0 {
			queryDeleteImgs := fmt.Sprintf("DELETE FROM `office_images` WHERE `office_id` = %d", office.ID)
	
			or.conn.Table("office_images").Exec(queryDeleteImgs)
		
			// insert to pivot table `office_images`
			for _, v := range officeDomain.Images {
				querySQL := fmt.Sprintf("INSERT INTO `office_images`(`url`, `office_id`) VALUES ('%s', '%d')", v, office.ID)

				if err := or.conn.Table("office_images").Exec(querySQL).Error; err != nil {
					return err
				}
			}
		}

		queryDeleteFacs := fmt.Sprintf("DELETE FROM `office_facilities` WHERE `office_id` = %d", office.ID)
	
		if err := or.conn.Table("office_images").Exec(queryDeleteFacs).Error; err != nil {
			return err
		}

		// insert to pivot table `office_facilities`
		for _, v := range facilitiesIdList {
			querySQL := fmt.Sprintf("INSERT INTO `office_facilities`(`facilities_id`, `office_id`) VALUES ('%d','%d')", v, office.ID)
			if err := tx.Table("office_facilities").Exec(querySQL).Error; err != nil {
				return err
			}
		}

		updatedOffice.Title = officeDomain.Title
		updatedOffice.Description = officeDomain.Description
		updatedOffice.OfficeType = officeDomain.OfficeType
		updatedOffice.OfficeLength = officeDomain.OfficeLength
		updatedOffice.Price = officeDomain.Price
		updatedOffice.OpenHour = officeDomain.OpenHour
		updatedOffice.CloseHour = officeDomain.CloseHour
		updatedOffice.Lat = officeDomain.Lat
		updatedOffice.Lng = officeDomain.Lng
		updatedOffice.Accommodate = officeDomain.Accommodate
		updatedOffice.WorkingDesk = officeDomain.WorkingDesk
		updatedOffice.MeetingRoom = officeDomain.MeetingRoom
		updatedOffice.PrivateRoom = officeDomain.PrivateRoom
		updatedOffice.City = officeDomain.City
		updatedOffice.District = officeDomain.District
		updatedOffice.Address = officeDomain.Address
		
		tx.Save(&updatedOffice)

		return nil
	})

	if err != nil {
		updatedOffice.ID = 0
		return updatedOffice.ToDomain()
	}

	return updatedOffice.ToDomain()
}

func (or *officeRepository) Delete(id string) bool {
	var office offices.Domain = or.GetByID(id)

	if office.ID == 0 {
		return false
	}

	deletedOffice := FromDomain(&office)
	result := or.conn.Delete(&deletedOffice)

	if result.RowsAffected == 0 {
		return false
	}

	queryDeleteImgs := fmt.Sprintf("DELETE FROM `office_images` WHERE `office_id` = '%d'", deletedOffice.ID)
	or.conn.Table("office_images").Exec(queryDeleteImgs)

	queryDeleteFac := fmt.Sprintf("DELETE FROM `office_facilities` WHERE `office_id` = '%d'", deletedOffice.ID)
	or.conn.Table("office_facilities").Exec(queryDeleteFac)

	return true
}

func (or *officeRepository) SearchByCity(city string) []offices.Domain {
	var rec []Office

	or.conn.Find(&rec, "city = ?", city)

	var imgsUrlPerID []imgs

	queryGetImgs := "SELECT `offices`.`id`, GROUP_CONCAT(office_images.url ORDER BY office_images.id SEPARATOR ' , ') AS images " + 
		"FROM offices " + 
		"INNER JOIN office_images on offices.id = office_images.office_id " + 
		"GROUP BY offices.id"
	or.conn.Raw(queryGetImgs).Scan(&imgsUrlPerID)

	var officeFacilitiesPerID []facilities
	queryGetFacilities := "SELECT `offices`.`id`, " + 
		"GROUP_CONCAT(`office_facilities`.`facilities_id` ORDER BY `office_facilities`.`facilities_id` SEPARATOR ' , ') AS f_id, " + 
		"GROUP_CONCAT(`facilities`.`description` ORDER BY `office_facilities`.`facilities_id` SEPARATOR ' , ') AS f_desc, " + 
		"GROUP_CONCAT(`facilities`.`slug` ORDER BY `office_facilities`.`facilities_id` SEPARATOR ' , ') AS f_slug " + 
		"FROM `offices` " + 
		"INNER JOIN `office_facilities` ON `offices`.`id`=`office_facilities`.`office_id` " + 
		"INNER JOIN `facilities` ON `office_facilities`.`facilities_id`=`facilities`.`id` " + 
		"GROUP BY `offices`.`id`"
	or.conn.Raw(queryGetFacilities).Scan(&officeFacilitiesPerID)

	officeDomain := []offices.Domain{}

	for _, office := range rec {
		for _, v := range imgsUrlPerID {
			if strconv.Itoa(int(office.ID)) == v.Id {
				url := v.Images
				img := strings.Split(url, " , ")
				office.Images = img
			}
		}

		for _, fac := range officeFacilitiesPerID {
			if strconv.Itoa(int(office.ID)) == fac.Id {
				f_id := fac.F_id
				facilitesId := strings.Split(f_id, " , ")
				f_desc := fac.F_desc
				facilitesDesc := strings.Split(f_desc, " , ")
				f_slug := fac.F_slug
				facilitiesSlug := strings.Split(f_slug, " , ")

				office.FacilitiesId =  facilitesId
				office.FacilitiesDesc = facilitesDesc
				office.FacilitesSlug = facilitiesSlug
			}
		}

		var count int64

		or.conn.Table("transactions").Not(map[string]interface{}{"status": []string{"rejected", "cancelled"}}).Where("office_id = ?", office.ID).Count(&count)

		office.TotalBooked = count

		var rate_score float64

		or.conn.Table("reviews").Where("office_id = ?", office.ID).Select("round(avg(`score`), 1)").Scan(&rate_score)

		office.Rate = rate_score

		officeDomain = append(officeDomain, office.ToDomain())
	}

	return officeDomain
}

func (or *officeRepository) SearchByRate(rate string) []offices.Domain {
	rec  := or.GetAll()
	var officeDomain []offices.Domain
	intRate, _ := strconv.Atoi(rate)

	for _, v := range rec {
		office := FromDomain(&v)

		switch intRate {
			case 5:
				if office.Rate == 5 {
					officeDomain = append(officeDomain, office.ToDomain())
				}
			case 0:
				if office.Rate == 0 {
					officeDomain = append(officeDomain, office.ToDomain())
				}
			default:
				if office.Rate >= float64(intRate) && office.Rate < (float64(intRate) + 1) {
					officeDomain = append(officeDomain, office.ToDomain())
				}
		}
	}

	return officeDomain
}

func (or *officeRepository) SearchByTitle(title string) []offices.Domain {
	var rec []Office

	or.conn.Find(&rec, "title = ?", title)

	var imgsUrlPerID []imgs

	queryGetImgs := "SELECT `offices`.`id`, " + 
		"GROUP_CONCAT(office_images.url ORDER BY office_images.id SEPARATOR ' , ') AS images " + 
		"FROM offices " + 
		"INNER JOIN office_images on offices.id = office_images.office_id " + 
		"GROUP BY offices.id"
	or.conn.Raw(queryGetImgs).Scan(&imgsUrlPerID)

	var officeFacilitiesPerID []facilities
	queryGetFacilities := "SELECT `offices`.`id`, " + 
		"GROUP_CONCAT(`office_facilities`.`facilities_id` ORDER BY `office_facilities`.`facilities_id` SEPARATOR ' , ') AS f_id, " + 
		"GROUP_CONCAT(`facilities`.`description` ORDER BY `office_facilities`.`facilities_id` SEPARATOR ' , ') AS f_desc, " + 
		"GROUP_CONCAT(`facilities`.`slug` ORDER BY `office_facilities`.`facilities_id` SEPARATOR ' , ') AS f_slug " + 
		"FROM `offices` " + 
		"INNER JOIN `office_facilities` ON `offices`.`id`=`office_facilities`.`office_id` " + 
		"INNER JOIN `facilities` ON `office_facilities`.`facilities_id`=`facilities`.`id` " + 
		"GROUP BY `offices`.`id`"
	or.conn.Raw(queryGetFacilities).Scan(&officeFacilitiesPerID)

	officeDomain := []offices.Domain{}

	for _, office := range rec {
		for _, v := range imgsUrlPerID {
			if strconv.Itoa(int(office.ID)) == v.Id {
				url := v.Images
				img := strings.Split(url, " , ")
				office.Images = img
			}
		}

		for _, fac := range officeFacilitiesPerID {
			if strconv.Itoa(int(office.ID)) == fac.Id {
				f_id := fac.F_id
				facilitesId := strings.Split(f_id, " , ")
				f_desc := fac.F_desc
				facilitesDesc := strings.Split(f_desc, " , ")
				f_slug := fac.F_slug
				facilitiesSlug := strings.Split(f_slug, " , ")

				office.FacilitiesId =  facilitesId
				office.FacilitiesDesc = facilitesDesc
				office.FacilitesSlug = facilitiesSlug
			}
		}

		var count int64

		or.conn.Table("transactions").Not(map[string]interface{}{"status": []string{"rejected", "cancelled"}}).Where("office_id = ?", office.ID).Count(&count)

		office.TotalBooked = count

		var rate_score float64

		or.conn.Table("reviews").Where("office_id = ?", office.ID).Select("round(avg(`score`), 1)").Scan(&rate_score)

		office.Rate = rate_score

		officeDomain = append(officeDomain, office.ToDomain())
	}

	return officeDomain
}

func (or *officeRepository) GetOffices() []offices.Domain {
	var rec []Office

	or.conn.Find(&rec, "office_type = ?", "office")

	var imgsUrlPerID []imgs

	queryGetImgs := "SELECT `offices`.`id`, " + 
		"GROUP_CONCAT( office_images.url ORDER BY office_images.id SEPARATOR ' , ') AS images " + 
		"FROM offices " + 
		"INNER JOIN office_images on offices.id = office_images.office_id " + 
		"GROUP BY offices.id"
	or.conn.Raw(queryGetImgs).Scan(&imgsUrlPerID)

	var officeFacilitiesPerID []facilities
	queryGetFacilities := "SELECT `offices`.`id`, " + 
		"GROUP_CONCAT(`office_facilities`.`facilities_id` ORDER BY `office_facilities`.`facilities_id` SEPARATOR ' , ') AS f_id, " + 
		"GROUP_CONCAT(`facilities`.`description` ORDER BY `office_facilities`.`facilities_id` SEPARATOR ' , ') AS f_desc, " + 
		"GROUP_CONCAT(`facilities`.`slug` ORDER BY `office_facilities`.`facilities_id` SEPARATOR ' , ') AS f_slug " + 
		"FROM `offices` " + 
		"INNER JOIN `office_facilities` ON `offices`.`id`=`office_facilities`.`office_id` " + 
		"INNER JOIN `facilities` ON `office_facilities`.`facilities_id`=`facilities`.`id` " + 
		"GROUP BY `offices`.`id`"
	or.conn.Raw(queryGetFacilities).Scan(&officeFacilitiesPerID)

	officeDomain := []offices.Domain{}

	for _, office := range rec {
		for _, v := range imgsUrlPerID {
			if strconv.Itoa(int(office.ID)) == v.Id {
				url := v.Images
				img := strings.Split(url, " , ")
				office.Images = img
			}
		}

		for _, fac := range officeFacilitiesPerID {
			if strconv.Itoa(int(office.ID)) == fac.Id {
				f_id := fac.F_id
				facilitesId := strings.Split(f_id, " , ")
				f_desc := fac.F_desc
				facilitesDesc := strings.Split(f_desc, " , ")
				f_slug := fac.F_slug
				facilitiesSlug := strings.Split(f_slug, " , ")

				office.FacilitiesId =  facilitesId
				office.FacilitiesDesc = facilitesDesc
				office.FacilitesSlug = facilitiesSlug
			}
		}

		var count int64

		or.conn.Table("transactions").Not(map[string]interface{}{"status": []string{"rejected", "cancelled"}}).Where("office_id = ?", office.ID).Count(&count)

		office.TotalBooked = count

		var rate_score float64

		or.conn.Table("reviews").Where("office_id = ?", office.ID).Select("round(avg(`score`), 1)").Scan(&rate_score)

		office.Rate = rate_score

		officeDomain = append(officeDomain, office.ToDomain())
	}

	return officeDomain
}

func (or *officeRepository) GetCoworkingSpace() []offices.Domain {
	var rec []Office

	or.conn.Find(&rec, "office_type = ?", "coworking space")

	var imgsUrlPerID []imgs

	queryGetImgs := "SELECT `offices`.`id`, " + 
		"GROUP_CONCAT( office_images.url ORDER BY office_images.id SEPARATOR ' , ') AS images " + 
		"FROM offices " + 
		"INNER JOIN office_images on offices.id = office_images.office_id " + 
		"GROUP BY offices.id"
	or.conn.Raw(queryGetImgs).Scan(&imgsUrlPerID)

	var officeFacilitiesPerID []facilities
	queryGetFacilities := "SELECT `offices`.`id`, " + 
		"GROUP_CONCAT(`office_facilities`.`facilities_id` ORDER BY `office_facilities`.`facilities_id` SEPARATOR ' , ') AS f_id, " + 
		"GROUP_CONCAT(`facilities`.`description` ORDER BY `office_facilities`.`facilities_id` SEPARATOR ' , ') AS f_desc, " + 
		"GROUP_CONCAT(`facilities`.`slug` ORDER BY `office_facilities`.`facilities_id` SEPARATOR ' , ') AS f_slug " + 
		"FROM `offices` " + 
		"INNER JOIN `office_facilities` ON `offices`.`id`=`office_facilities`.`office_id` " + 
		"INNER JOIN `facilities` ON `office_facilities`.`facilities_id`=`facilities`.`id` " + 
		"GROUP BY `offices`.`id`"
	or.conn.Raw(queryGetFacilities).Scan(&officeFacilitiesPerID)

	officeDomain := []offices.Domain{}

	for _, office := range rec {
		for _, v := range imgsUrlPerID {
			if strconv.Itoa(int(office.ID)) == v.Id {
				url := v.Images
				img := strings.Split(url, " , ")
				office.Images = img
			}
		}

		for _, fac := range officeFacilitiesPerID {
			if strconv.Itoa(int(office.ID)) == fac.Id {
				f_id := fac.F_id
				facilitesId := strings.Split(f_id, " , ")
				f_desc := fac.F_desc
				facilitesDesc := strings.Split(f_desc, " , ")
				f_slug := fac.F_slug
				facilitiesSlug := strings.Split(f_slug, " , ")

				office.FacilitiesId =  facilitesId
				office.FacilitiesDesc = facilitesDesc
				office.FacilitesSlug = facilitiesSlug
			}
		}

		var count int64

		or.conn.Table("transactions").Not(map[string]interface{}{"status": []string{"rejected", "cancelled"}}).Where("office_id = ?", office.ID).Count(&count)

		office.TotalBooked = count

		var rate_score float64

		or.conn.Table("reviews").Where("office_id = ?", office.ID).Select("round(avg(`score`), 1)").Scan(&rate_score)

		office.Rate = rate_score

		officeDomain = append(officeDomain, office.ToDomain())
	}

	return officeDomain
}

func (or *officeRepository) GetMeetingRooms() []offices.Domain {
	var rec []Office

	or.conn.Find(&rec, "office_type = ?", "meeting room")

	var imgsUrlPerID []imgs

	queryGetImgs := "SELECT `offices`.`id`, " + 
		"GROUP_CONCAT( office_images.url ORDER BY office_images.id SEPARATOR ' , ') AS images " + 
		"FROM offices " + 
		"INNER JOIN office_images on offices.id = office_images.office_id " + 
		"GROUP BY offices.id"
	or.conn.Raw(queryGetImgs).Scan(&imgsUrlPerID)

	var officeFacilitiesPerID []facilities
	queryGetFacilities := "SELECT `offices`.`id`, " + 
		"GROUP_CONCAT(`office_facilities`.`facilities_id` ORDER BY `office_facilities`.`facilities_id` SEPARATOR ' , ') AS f_id, " + 
		"GROUP_CONCAT(`facilities`.`description` ORDER BY `office_facilities`.`facilities_id` SEPARATOR ' , ') AS f_desc, " + 
		"GROUP_CONCAT(`facilities`.`slug` ORDER BY `office_facilities`.`facilities_id` SEPARATOR ' , ') AS f_slug " + 
		"FROM `offices` " + 
		"INNER JOIN `office_facilities` ON `offices`.`id`=`office_facilities`.`office_id` " + 
		"INNER JOIN `facilities` ON `office_facilities`.`facilities_id`=`facilities`.`id` " + 
		"GROUP BY `offices`.`id`"
	or.conn.Raw(queryGetFacilities).Scan(&officeFacilitiesPerID)

	officeDomain := []offices.Domain{}

	for _, office := range rec {
		for _, v := range imgsUrlPerID {
			if strconv.Itoa(int(office.ID)) == v.Id {
				url := v.Images
				img := strings.Split(url, " , ")
				office.Images = img
			}
		}

		for _, fac := range officeFacilitiesPerID {
			if strconv.Itoa(int(office.ID)) == fac.Id {
				f_id := fac.F_id
				facilitesId := strings.Split(f_id, " , ")
				f_desc := fac.F_desc
				facilitesDesc := strings.Split(f_desc, " , ")
				f_slug := fac.F_slug
				facilitiesSlug := strings.Split(f_slug, " , ")

				office.FacilitiesId =  facilitesId
				office.FacilitiesDesc = facilitesDesc
				office.FacilitesSlug = facilitiesSlug
			}
		}

		var count int64

		or.conn.Table("transactions").Not(map[string]interface{}{"status": []string{"rejected", "cancelled"}}).Where("office_id = ?", office.ID).Count(&count)

		office.TotalBooked = count

		var rate_score float64

		or.conn.Table("reviews").Where("office_id = ?", office.ID).Select("round(avg(`score`), 1)").Scan(&rate_score)

		office.Rate = rate_score

		officeDomain = append(officeDomain, office.ToDomain())
	}

	return officeDomain
}

func (or *officeRepository) GetRecommendation() []offices.Domain {
	var rec []Office

	or.conn.Order("rate desc, title, description").Find(&rec)
	
	var imgsUrlPerID []imgs

	queryGetImgs := "SELECT `offices`.`id`, " + 
		"GROUP_CONCAT(office_images.url ORDER BY office_images.id SEPARATOR ' , ') AS images " + 
		"FROM offices " + 
		"INNER JOIN office_images on offices.id = office_images.office_id " + 
		"GROUP BY offices.id"
	or.conn.Raw(queryGetImgs).Scan(&imgsUrlPerID)

	var officeFacilitiesPerID []facilities
	queryGetFacilities := "SELECT `offices`.`id`, " + 
		"GROUP_CONCAT(`office_facilities`.`facilities_id` ORDER BY `office_facilities`.`facilities_id` SEPARATOR ' , ') AS f_id, " + 
		"GROUP_CONCAT(`facilities`.`description` ORDER BY `office_facilities`.`facilities_id` SEPARATOR ' , ') AS f_desc, " + 
		"GROUP_CONCAT(`facilities`.`slug` ORDER BY `office_facilities`.`facilities_id` SEPARATOR ' , ') AS f_slug " + 
		"FROM `offices` " + 
		"INNER JOIN `office_facilities` ON `offices`.`id`=`office_facilities`.`office_id` " + 
		"INNER JOIN `facilities` ON `office_facilities`.`facilities_id`=`facilities`.`id` " + 
		"GROUP BY `offices`.`id`"
	or.conn.Raw(queryGetFacilities).Scan(&officeFacilitiesPerID)

	officeDomain := []offices.Domain{}
	
	for _, office := range rec {
		for _, v := range imgsUrlPerID {
			if strconv.Itoa(int(office.ID)) == v.Id {
				url := v.Images
				img := strings.Split(url, " , ")
				office.Images = img
			}
		}

		for _, fac := range officeFacilitiesPerID {
			if strconv.Itoa(int(office.ID)) == fac.Id {
				f_id := fac.F_id
				facilitesId := strings.Split(f_id, " , ")
				f_desc := fac.F_desc
				facilitesDesc := strings.Split(f_desc, " , ")
				f_slug := fac.F_slug
				facilitiesSlug := strings.Split(f_slug, " , ")

				office.FacilitiesId =  facilitesId
				office.FacilitiesDesc = facilitesDesc
				office.FacilitesSlug = facilitiesSlug
			}
		}

		var count int64

		or.conn.Table("transactions").Not(map[string]interface{}{"status": []string{"rejected", "cancelled"}}).Where("office_id = ?", office.ID).Count(&count)

		office.TotalBooked = count

		var rate_score float64

		or.conn.Table("reviews").Where("office_id = ?", office.ID).Select("round(avg(`score`), 1)").Scan(&rate_score)

		office.Rate = rate_score

		officeDomain = append(officeDomain, office.ToDomain())
	}

	return officeDomain
}

func (or *officeRepository) GetNearest(lat string, long string) []offices.Domain {
	var rec []Office

	or.conn.Find(&rec)
	
	var imgsUrlPerID []imgs

	queryGetImgs := "SELECT `offices`.`id`, " + 
		"GROUP_CONCAT( office_images.url ORDER BY office_images.id SEPARATOR ' , ') AS images " + 
		"FROM offices " + 
		"INNER JOIN office_images on offices.id = office_images.office_id " + 
		"GROUP BY offices.id"
	or.conn.Raw(queryGetImgs).Scan(&imgsUrlPerID)

	var officeFacilitiesPerID []facilities
	queryGetFacilities := "SELECT `offices`.`id`, " + 
		"GROUP_CONCAT(`office_facilities`.`facilities_id` ORDER BY `office_facilities`.`facilities_id` SEPARATOR ' , ') AS f_id, " + 
		"GROUP_CONCAT(`facilities`.`description` ORDER BY `office_facilities`.`facilities_id` SEPARATOR ' , ') AS f_desc, " + 
		"GROUP_CONCAT(`facilities`.`slug` ORDER BY `office_facilities`.`facilities_id` SEPARATOR ' , ') AS f_slug " + 
		"FROM `offices` " + 
		"INNER JOIN `office_facilities` ON `offices`.`id`=`office_facilities`.`office_id` " + 
		"INNER JOIN `facilities` ON `office_facilities`.`facilities_id`=`facilities`.`id` " + 
		"GROUP BY `offices`.`id`"
	or.conn.Raw(queryGetFacilities).Scan(&officeFacilitiesPerID)
	
	// find nearest logic here
	var distance []distance

	// distance is in kilometer
	queryGetDistance := fmt.Sprintf("SELECT `offices`.`id`, " + 
		"CAST(" + 
			"SQRT(" + 
				"POW(69.1 * (`offices`.`lat` - '%s'), 2) + " + 
				"POW(69.1 * ('%s' - `offices`.`lng`) * " + 
				"COS(`offices`.`lng` / 57.3), 2)) * 1.60934 " + 
			"AS DECIMAL(16,2)) " + 
		"AS distance " + 
		"FROM `offices` " + 
		"HAVING distance < 40 " + 
		"ORDER BY distance", lat, long)

	or.conn.Raw(queryGetDistance).Scan(&distance)

	officeDomain := []offices.Domain{}

	for _, d := range distance {
		for _, office := range rec {
			for _, v := range imgsUrlPerID {
				if strconv.Itoa(int(office.ID)) == v.Id {
					url := v.Images
					img := strings.Split(url, " , ")
					office.Images = img
				}
			}
	
			for _, fac := range officeFacilitiesPerID {
				if strconv.Itoa(int(office.ID)) == fac.Id {
					f_id := fac.F_id
					facilitesId := strings.Split(f_id, " , ")
					f_desc := fac.F_desc
					facilitesDesc := strings.Split(f_desc, " , ")
					f_slug := fac.F_slug
					facilitiesSlug := strings.Split(f_slug, " , ")
	
					office.FacilitiesId =  facilitesId
					office.FacilitiesDesc = facilitesDesc
					office.FacilitesSlug = facilitiesSlug
				}
			}

			var count int64

			or.conn.Table("transactions").Not(map[string]interface{}{"status": []string{"rejected", "cancelled"}}).Where("office_id = ?", office.ID).Count(&count)

			office.TotalBooked = count

			var rate_score float64

			or.conn.Table("reviews").Where("office_id = ?", office.ID).Select("round(avg(`score`), 1)").Scan(&rate_score)

			office.Rate = rate_score

			if strconv.Itoa(int(office.ID)) == d.Id {
				office.Distance = d.Distance
				officeDomain = append(officeDomain, office.ToDomain())
			}
		}
	}

	return officeDomain
}