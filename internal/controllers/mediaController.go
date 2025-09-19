package controllers

import (
	"net/http"

	"github.com/PC-Core/pc-core-backend/internal/errors"
	"github.com/PC-Core/pc-core-backend/internal/static"
	"github.com/gin-gonic/gin"
	"github.com/go-http-utils/headers"
)

type Locations struct {
	Locations []string `json:"locations"`
}

type MediaController struct {
	engine *gin.Engine
	static static.StaticDataController
}

func NewStaticController(engine *gin.Engine, static static.StaticDataController) *MediaController {
	return &MediaController{
		engine, static,
	}
}

func (c *MediaController) ApplyRoutes() {
	gr := c.engine.Group("/media")
	{
		gr.POST("/upload", c.uploadFile)
	}
}

// Upload media
// @Summary      Upload media
// @Tags         media
// @Accept       json
// @Produce      json
// @Param 		 uploads[] formData		file	true	"files to upload"
// @Success      200  {array}  Locations
// @Failure      400  {object}  errors.PublicPCCError
// @Router       /media/upload [post]
func (c *MediaController) uploadFile(ctx *gin.Context) {
	form, err := ctx.MultipartForm()

	if err != nil {
		// TODO: Error type
		CheckErrorAndWriteBadRequest(ctx, errors.NewInternalSecretError())
		return
	}

	mfiles := form.File["upload[]"]

	files := make([]static.StaticFile, 0, len(mfiles))

	for _, file := range mfiles {
		var err error
		rfile, err := file.Open()

		if err != nil {
			// TODO: Error type
			CheckErrorAndWriteBadRequest(ctx, errors.NewInternalSecretError())
			return
		}

		files = append(files, *static.NewStaticFile(rfile, file.Filename, file.Header.Get(headers.ContentType)))
	}

	locs, uerr := c.static.UploadFiles(files)

	if CheckErrorAndWriteBadRequest(ctx, uerr) {
		return
	}

	ctx.JSON(http.StatusOK, Locations{
		locs,
	})
}
