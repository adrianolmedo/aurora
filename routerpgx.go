package main

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
)

func routerPXG(storage *StoragePGX) *fiber.App {
	s := NewServicePGX(storage)
	f := fiber.New()

	g := f.Group("/v2/users")
	g.Get("", listUsersPGX(s))
	return f
}

// listUsers handler GET: /users
func listUsersPGX(s *ServicePGX) fiber.Handler {
	return func(c *fiber.Ctx) error {
		f, err := getFilter(c)
		if err != nil {
			resp := respJSON(msgError, err.Error(), nil)
			return c.Status(http.StatusBadRequest).JSON(resp) // 400
		}

		fr, err := s.User.List(f)
		if err != nil {
			resp := respJSON(msgError, err.Error(), nil)
			return c.Status(http.StatusInternalServerError).JSON(resp)
		}

		users, ok := fr.Rows.(Users)
		if !ok {
			resp := respJSON(msgError, "error in data assertion", nil)
			return c.Status(http.StatusInternalServerError).JSON(resp)
		}

		if users.IsEmpty() {
			resp := respJSON(msgOK, "there are not users", nil)
			return c.Status(http.StatusOK).JSON(resp)
		}

		ls := f.GenLinksResp(c.Path(), fr.TotalPages)
		resp := respJSON(msgOK, "", users).setLinks(ls).setMeta(fr)
		return c.Status(http.StatusOK).JSON(resp)
	}
}
