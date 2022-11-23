package users

import "backend/app/middlewares"

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

	if user.ID == "" {
		return ""
	}

	if user.Roles {
		token := uu.jwtAuth.GenerateAdminToken(user.ID)
		return token
	}

	token := uu.jwtAuth.GenerateToken(user.ID)

	return token
}