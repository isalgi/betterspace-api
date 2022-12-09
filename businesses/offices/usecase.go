package offices

type OfficeUsecase struct {
	officeRepository Repository
}

func NewOfficeUsecase(or Repository) Usecase {
	return &OfficeUsecase{
		officeRepository: or,
	}
}

func (ou *OfficeUsecase) GetAll() []Domain {
	return ou.officeRepository.GetAll()
}

func (ou *OfficeUsecase) GetByID(id string) Domain {
	return ou.officeRepository.GetByID(id)
}

func (ou *OfficeUsecase) Create(officeDomain *Domain) Domain {
	return ou.officeRepository.Create(officeDomain)
}

func (ou *OfficeUsecase) Update(id string, officeDomain *Domain) Domain {
	return ou.officeRepository.Update(id, officeDomain)
}

func (ou *OfficeUsecase) Delete(id string) bool {
	return ou.officeRepository.Delete(id)
}

func (ou *OfficeUsecase) SearchByCity(city string) []Domain {
	return ou.officeRepository.SearchByCity(city)
}

func (ou *OfficeUsecase) SearchByRate(rate string) []Domain {
	return ou.officeRepository.SearchByRate(rate)
}

func (ou *OfficeUsecase) SearchByTitle(title string) []Domain {
	return ou.officeRepository.SearchByTitle(title)
}

func (ou *OfficeUsecase) GetCoworkingSpace() []Domain {
	return ou.officeRepository.GetCoworkingSpace()
}

func (ou *OfficeUsecase) GetOffices() []Domain {
	return ou.officeRepository.GetOffices()
}

func (ou *OfficeUsecase) GetMeetingRooms() []Domain {
	return ou.officeRepository.GetMeetingRooms()
}

func (ou *OfficeUsecase) GetRecommendation() []Domain {
	return ou.officeRepository.GetRecommendation()
}

func (ou *OfficeUsecase) GetNearest(lat string, long string) []Domain {
	return ou.officeRepository.GetNearest(lat, long)
}