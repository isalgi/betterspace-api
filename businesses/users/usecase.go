package users

import (
	"backend/app/middlewares"
	"strconv"
)

type UserUsecase struct {
	userRepository Repository
	jwtAuth        *middlewares.ConfigJWT
}

func NewUserUsecase(ur Repository, jwtAuth *middlewares.ConfigJWT) Usecase {
	return &UserUsecase{
		userRepository: ur,
		jwtAuth:        jwtAuth,
	}
}

func (uu *UserUsecase) Register(userDomain *Domain) Domain {
	return uu.userRepository.Register(userDomain)
}

func (uu *UserUsecase) Login(userDomain *LoginDomain) string {
	user := uu.userRepository.GetByEmail(userDomain)

	if user.ID == 0 {
		return ""
	}

	token := uu.jwtAuth.GenerateToken(strconv.Itoa(int(user.ID)), user.Roles)

	return token
}

func (uu *UserUsecase) GetAll() []Domain {
	return uu.userRepository.GetAll()
}

func (uu *UserUsecase) GetByID(id string) Domain {
	return uu.userRepository.GetByID(id)
}

func (uu *UserUsecase) Delete(id string) bool {
	return uu.userRepository.Delete(id)
}

func (uu *UserUsecase) UpdateProfilePhoto(id string, userDomain *PhotoDomain) bool {
	return uu.userRepository.InsertURLtoUser(id, userDomain)
}