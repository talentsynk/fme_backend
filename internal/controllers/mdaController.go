package controllers

import (
	"fme_backend/internal/config"
	"fme_backend/internal/responses"
	"fme_backend/internal/schemas"
	"fme_backend/internal/servers"
	"fme_backend/internal/utilities"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

var mdaServer servers.MdaServer

type mdaController struct{}

func (m *mdaController) CreateMda(c *fiber.Ctx) error {
	var body schemas.MdaCreateSchema
	if c.BodyParser(&body) != nil {
		return responses.ErrorResponse(c, responses.BAD_DATA, 400)
	}
	if err := mdaServer.CreateMda(config.DB, body); err != nil {
		return responses.ErrorResponse(c, err.Error(), 400)
	}
	return responses.SuccessResponse(c, "message", responses.USER_CREATED, 200)
}

func (m *mdaController) GetAllMdas(c *fiber.Ctx) error {
	limitStr := c.Query("limit")
	pageStr := c.Query("page")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 100
	}
	page, err := strconv.Atoi(pageStr)
	if err != nil || page <= 0 {
		page = 1
	}
	offset := (page - 1) * limit
	resp, err := mdaServer.GetAllMdas(config.DB, limit, offset)
	if err != nil {
		return responses.ErrorResponse(c, err.Error(), 400)
	}
	return responses.SuccessResponse(c, "mdas", resp, 200)
}

func (m *mdaController) GetMdaByID(c *fiber.Ctx) error {
	id_str := c.Params("id")
	id, err := strconv.Atoi(id_str)
	if err != nil {
		return responses.ErrorResponse(c, responses.INCOMPLETE_DATA, 400)
	}

	resp, err := mdaServer.GetMdaById(config.DB, id)
	if err != nil {
		return responses.ErrorResponse(c, err.Error(), 400)
	}
	return responses.CustomResponse(c, 200, resp)
}

func (m *mdaController) UpdateMda(c *fiber.Ctx) error {
	id_str := c.Params("id")
	id, err := strconv.Atoi(id_str)
	if err != nil {
		return responses.ErrorResponse(c, responses.INCOMPLETE_DATA, 400)
	}
	var body schemas.UpdateMdaSchema
	if c.BodyParser(&body) != nil {
		return responses.ErrorResponse(c, responses.BAD_DATA, 400)
	}
	body.Id = id
	resp, err := mdaServer.UpdateMda(config.DB, body)
	if err != nil {
		return responses.ErrorResponse(c, err.Error(), 400)
	}
	return responses.CustomResponse(c, 200, resp)
}

func (m *mdaController) SearchMda(c *fiber.Ctx) error {
	query := c.Query("query")
	if utilities.IsEmpty(query) {
		return responses.ErrorResponse(c, responses.INCOMPLETE_DATA, 400)
	}
	resp, err := mdaServer.SearchMda(config.DB, query)
	if err != nil {
		return responses.ErrorResponse(c, err.Error(), 400)
	}
	return responses.CustomResponse(c, 200, resp)
}

func (m *mdaController) SuspendMda(c *fiber.Ctx) error {
	id := c.Params("id")
	if utilities.IsEmpty(id) {
		return responses.ErrorResponse(c, responses.INCOMPLETE_DATA, 400)
	}
	if err := mdaServer.SuspendMda(config.DB, id); err != nil {
		return responses.ErrorResponse(c, err.Error(), 400)
	}
	return responses.SuccessResponse(c, "suspend", "Mda Suspended", 200)
}

func (m *mdaController) ActivateMda(c *fiber.Ctx) error {
	id := c.Params("id")
	if utilities.IsEmpty(id) {
		return responses.ErrorResponse(c, responses.INCOMPLETE_DATA, 400)
	}
	if err := mdaServer.ActivateMda(config.DB, id); err != nil {
		return responses.ErrorResponse(c, err.Error(), 400)
	}
	return responses.SuccessResponse(c, "activate", "Mda activated", 200)
}

func (m *mdaController) MdaTotal(c *fiber.Ctx) error {
	resp, err := mdaServer.MdaTotal(config.DB)
	if err != nil {
		return responses.ErrorResponse(c, err.Error(), 400)
	}
	return responses.CustomResponse(c, 200, resp)
}

func (m *mdaController) GetMdaProfile(c *fiber.Ctx) error {
	mda_id := c.Get("mdaID")
	id, _ := strconv.Atoi(mda_id)
	if utilities.IsEmpty(id) {
		return responses.ErrorResponse(c, responses.UNAUTHORIZED_ACCESS, 401)
	}
	resp, err := mdaServer.GetMdaProfile(config.DB, id)
	if err != nil {
		return responses.ErrorResponse(c, err.Error(), 400)
	}
	return responses.CustomResponse(c, 200, resp)
}

// TODO: fix this
func (m *mdaController) DownloadMdaCsv(c *fiber.Ctx) error {
	return nil
}

func (m *mdaController) EditMdaData(c *fiber.Ctx) error {
	var body schemas.MdaCreateSchema
	if c.BodyParser(&body) != nil {
		return responses.ErrorResponse(c, responses.BAD_DATA, 400)
	}
	idStr := c.Params("id")
	id, _ := strconv.Atoi(idStr)
	_, valid_state := utilities.ValidateState(body.StateOfOperation)
	if utilities.IsEmpty(id) || (!utilities.IsEmpty(body.StateOfOperation) && !valid_state) {
		return responses.ErrorResponse(c, responses.INCOMPLETE_DATA, 400)
	}
	if err := mdaServer.EditMdaData(config.DB, body, id); err != nil {
		return responses.ErrorResponse(c, err.Error(), 400)
	}
	return responses.SuccessResponse(c, "message", "successfully updated the record", 200)
}
