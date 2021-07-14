package app

import (
	"github.com/gofiber/fiber/v2"
	"schneider.vip/problem"
)

// ProblemJSONErrorHandle send error handle
func ProblemJSONErrorHandle(ctx *fiber.Ctx, err error) error {
	// Statuscode defaults to 500
	code := fiber.StatusInternalServerError

	// Retreive the custom statuscode if it's an fiber.*Error
	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
	}

	// Send custom error page
	err = ctx.Status(code).JSON(problem.Of(code).Append(problem.Detail(err.Error())))
	if err != nil {
		// In case the SendFile fails
		return ctx.Status(500).JSON(problem.Of(500).Append(problem.Detail(err.Error())))
	}

	// Return from handler
	return nil
}
