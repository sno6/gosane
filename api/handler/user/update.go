package user

import (
	"github.com/gin-gonic/gin"
	"github.com/sno6/gosane/api/dto"
	"github.com/sno6/gosane/api/handler/user/dto/request"
	"github.com/sno6/gosane/ent"
	"github.com/sno6/gosane/internal/http"
	"github.com/sno6/gosane/internal/validator"
	"github.com/sno6/gosane/middleware"
	"github.com/sno6/gosane/service/user"
)

type UpdateHandler struct {
	validator   *validator.Validator
	userService *user.Service
}

func NewUpdateHandler(
	validator *validator.Validator,
	userService *user.Service,
) *UpdateHandler {
	return &UpdateHandler{
		validator:   validator,
		userService: userService,
	}
}

func (*UpdateHandler) Path() string {
	return ""
}

func (*UpdateHandler) Method() string {
	return http.MethodPut
}

func (uh *UpdateHandler) HandleFunc(c *gin.Context) {
	var body request.UpdateUserBody
	if err := uh.validator.ValidateJSON(c.Request.Body, &body); err != nil {
		c.Error(http.BadRequest).SetMeta(err)
		return
	}

	u, err := middleware.UserFromContext(c)
	if err != nil {
		c.Error(http.Unauthorized).SetMeta(err)
		return
	}

	user, err := uh.userService.UpdateByUUID(c, u.UUID, &ent.User{
		FirstName: body.FirstName,
		LastName:  body.LastName,
	})
	if err != nil {
		c.Error(http.Internal).SetMeta(err)
		return
	}

	c.JSON(200, dto.NewFromUser(user))
}
