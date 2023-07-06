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
		users, err := s.User.List()
		if err != nil {
			resp := respJSON(msgError, "", nil)
			return c.Status(http.StatusInternalServerError).JSON(resp)
		}

		if users.IsEmpty() {
			resp := respJSON(msgOK, "there are not users", nil)
			return c.Status(http.StatusOK).JSON(resp)
		}

		resp := respJSON(msgOK, "", users)
		return c.Status(http.StatusCreated).JSON(resp)
	}
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
		return c.Status(http.StatusCreated).JSON(resp)
	}
}
