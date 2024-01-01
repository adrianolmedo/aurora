package main

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func router(storage *Storage) *fiber.App {
	s := NewService(storage)
	f := fiber.New()

	g := f.Group("/v1/users")
	g.Post("", signUpUser(s))
	g.Get("", listUsers(s))
	g.Get("/:id", findUser(s))
	g.Delete("/:id", deleteUser(s))
	return f
}

// signUpUser handler POST: /users
func signUpUser(s *Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		u := &User{}
		err := c.BodyParser(u)
		if err != nil {
			resp := respJSON(msgError, "the JSON structure is not correct", nil)
			return c.Status(http.StatusBadRequest).JSON(resp)
		}

		err = s.User.SignUp(u)
		if err != nil {
			resp := respJSON(msgError, err.Error(), nil)
			return c.Status(http.StatusInternalServerError).JSON(resp)
		}

		resp := respJSON(msgOK, "user created", u)
		return c.Status(http.StatusCreated).JSON(resp)
	}
}

// listUsers handler GET: /users
func listUsers(s *Service) fiber.Handler {
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

func getFilter(c *fiber.Ctx) (*Filter, error) {
	f := NewFilter(10)

	err := f.SetLimit(c.QueryInt("limit"))
	if err != nil {
		return nil, err
	}

	err = f.SetPage(c.QueryInt("page"))
	if err != nil {
		return nil, err
	}

	f.SetSort(c.Query("sort", "created_at"))
	f.SetDirection(c.Query("direction"))
	return f, nil
}

// findUser handler GET: /users/:id
func findUser(s *Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := strconv.Atoi(c.Params("id"))
		if id < 0 || err != nil {
			resp := respJSON(msgError, "positive number expected for ID user", nil)
			return c.Status(http.StatusBadRequest).JSON(resp)
		}

		user, err := s.User.Find(id)
		if errors.Is(err, ErrUserNotFound) {
			resp := respJSON(msgError, err.Error(), nil)
			return c.Status(http.StatusNotFound).JSON(resp)
		}

		if err != nil {
			resp := respJSON(msgError, err.Error(), nil)
			return c.Status(http.StatusBadRequest).JSON(resp)
		}

		resp := respJSON(msgOK, "", user)
		return c.Status(http.StatusOK).JSON(resp)
	}
}

// deleteUser handler DELETE: /users/:id
func deleteUser(s *Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := strconv.Atoi(c.Params("id"))
		if id < 0 || err != nil {
			resp := respJSON(msgError, "positive number expected for ID user", nil)
			return c.Status(http.StatusBadRequest).JSON(resp)
		}

		err = s.User.Remove(id)

		if errors.Is(err, ErrUserNotFound) {
			resp := respJSON(msgError, err.Error(), nil)
			return c.Status(http.StatusNoContent).JSON(resp)
		}

		if err != nil {
			resp := respJSON(msgError, err.Error(), nil)
			return c.Status(http.StatusInternalServerError).JSON(resp)
		}

		resp := respJSON(msgOK, "user deleted", nil)
		return c.Status(http.StatusOK).JSON(resp)
	}
}
